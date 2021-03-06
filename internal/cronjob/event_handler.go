package cronjob

import (
	"github.com/huan231/jobz/pkg/events"
	batchv1 "k8s.io/api/batch/v1"
)

type EventHandler interface {
	OnCronJobAdd(cronJob *batchv1.CronJob)
	OnCronJobUpdate(oldCronJob, newCronJob *batchv1.CronJob)
	OnCronJobDelete(cronJob *batchv1.CronJob)
}

type eventHandler struct {
	h events.Hub
}

func NewEventHandler(h events.Hub) EventHandler {
	return &eventHandler{h}
}

func (e *eventHandler) OnCronJobAdd(cronJob *batchv1.CronJob) {
	e.h.Publish(events.Event{Type: "cronjobadd", Payload: NewCronJob(cronJob)})
}

func (e *eventHandler) OnCronJobUpdate(oldCronJob, newCronJob *batchv1.CronJob) {
	if oldCronJob.Spec.Schedule == newCronJob.Spec.Schedule {
		return
	}

	e.h.Publish(events.Event{Type: "cronjobupdate", Payload: NewCronJob(newCronJob)})
}

type cronJobDeletePayload struct {
	ID string `json:"id"`
}

func (e *eventHandler) OnCronJobDelete(cronJob *batchv1.CronJob) {
	e.h.Publish(events.Event{Type: "cronjobdelete", Payload: cronJobDeletePayload{ID: string(cronJob.UID)}})
}
