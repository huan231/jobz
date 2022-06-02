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
	s Service
	h events.Hub
}

func NewEventHandler(s Service, h events.Hub) EventHandler {
	return &eventHandler{s, h}
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

func (e *eventHandler) OnCronJobDelete(cronJob *batchv1.CronJob) {
	e.h.Publish(events.Event{Type: "cronjobdelete", Payload: map[string]any{"id": string(cronJob.UID)}})
}
