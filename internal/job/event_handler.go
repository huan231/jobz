package job

import (
	"context"
	"github.com/huan231/jobz/pkg/events"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/listers/batch/v1"
)

type EventHandler interface {
	OnEventAdd(event *corev1.Event)
	OnCronJobDelete(cronJob *batchv1.CronJob)
}

type eventHandler struct {
	s             Service
	h             events.Hub
	cronJobLister v1.CronJobLister
	jobLister     v1.JobLister
}

func NewEventHandler(s Service, h events.Hub, cronJobLister v1.CronJobLister, jobLister v1.JobLister) EventHandler {
	return &eventHandler{s, h, cronJobLister, jobLister}
}

func (e *eventHandler) complete(job *Job, err error) {
	if err != nil {
		if err == ErrAlreadyCompleted {
			return
		}

		panic(err)
	}

	e.h.Publish(events.Event{Type: "jobcomplete", Payload: job})
}

func (e *eventHandler) OnEventAdd(event *corev1.Event) {
	if event.InvolvedObject.Kind != "Job" {
		return
	}

	job, err := e.jobLister.Jobs(event.Namespace).Get(event.InvolvedObject.Name)

	if err != nil || event.InvolvedObject.UID != job.UID {
		return
	}

	controllerRef := metav1.GetControllerOf(job)

	if controllerRef == nil || controllerRef.Kind != "CronJob" {
		return
	}

	cronJob, err := e.cronJobLister.CronJobs(event.Namespace).Get(controllerRef.Name)

	if err != nil || cronJob.UID != controllerRef.UID {
		return
	}

	ctx := context.Background()

	switch event.Reason {
	case "SuccessfulCreate":
		j, err := e.s.Create(ctx, string(job.UID), string(cronJob.UID), event.LastTimestamp.UTC())

		if err != nil {
			if err == ErrAlreadyExists {
				return
			}

			panic(err)
		}

		e.h.Publish(events.Event{Type: "jobadd", Payload: j})
	case "Completed":
		e.complete(e.s.Succeed(ctx, string(job.UID), event.LastTimestamp.UTC()))
	case "BackoffLimitExceeded", "DeadlineExceeded":
		e.complete(e.s.Fail(ctx, string(job.UID), event.LastTimestamp.UTC()))
	}
}

func (e *eventHandler) OnCronJobDelete(cronJob *batchv1.CronJob) {
	err := e.s.Delete(context.Background(), string(cronJob.UID))

	if err != nil {
		panic(err)
	}
}
