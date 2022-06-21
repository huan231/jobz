package cronjob

import (
	"github.com/huan231/jobz/pkg/events"
	"github.com/stretchr/testify/mock"
	batchv1 "k8s.io/api/batch/v1"
	"testing"
)

func TestOnCronJobAdd(t *testing.T) {
	h := &hubMock{}
	h.On("Publish", mock.Anything).Return()

	e := NewEventHandler(h)

	e.OnCronJobAdd(&batchv1.CronJob{})

	h.AssertCalled(t, "Publish", mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.Payload.(*CronJob)

		return e.Type == "cronjobadd" && ok
	}))
}

func TestOnCronJobUpdateScheduleUpdated(t *testing.T) {
	h := &hubMock{}
	h.On("Publish", mock.Anything).Return()

	e := NewEventHandler(h)

	e.OnCronJobUpdate(
		&batchv1.CronJob{Spec: batchv1.CronJobSpec{Schedule: "* * * * *"}},
		&batchv1.CronJob{Spec: batchv1.CronJobSpec{Schedule: "0 * * * *"}},
	)

	h.AssertCalled(t, "Publish", mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.Payload.(*CronJob)

		return e.Type == "cronjobupdate" && ok
	}))
}

func TestOnCronJobUpdateScheduleNotUpdated(t *testing.T) {
	h := &hubMock{}

	e := NewEventHandler(h)

	e.OnCronJobUpdate(&batchv1.CronJob{}, &batchv1.CronJob{})

	h.AssertNotCalled(t, "Publish")
}

func TestOnCronJobDelete(t *testing.T) {
	h := &hubMock{}
	h.On("Publish", mock.Anything).Return()

	e := NewEventHandler(h)

	e.OnCronJobDelete(&batchv1.CronJob{})

	h.AssertCalled(t, "Publish", mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.Payload.(cronJobDeletePayload)

		return e.Type == "cronjobdelete" && ok
	}))
}

type hubMock struct {
	mock.Mock
}

func (h *hubMock) Register(events.Subscriber) {
}

func (h *hubMock) Unregister(events.Subscriber) {
}

func (h *hubMock) Close() {
}

func (h *hubMock) Publish(e events.Event) {
	h.Called(e)
}
