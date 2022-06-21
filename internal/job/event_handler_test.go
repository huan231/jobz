package job

import (
	"context"
	"fmt"
	"github.com/huan231/jobz/pkg/events"
	"github.com/stretchr/testify/mock"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/batch/v1"
	"testing"
	"time"
)

func TestOnEventAddInvolvedObjectNotJob(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddJobNotFound(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(nil, fmt.Errorf("not found"))

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddJobIDMismatch(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(&batchv1.Job{}, nil)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{UID: "858558f7-305a-4266-9775-29e7143e6f99", Kind: "Job"}})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddJobNotCreatedByCronJob(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(&batchv1.Job{}, nil)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddCronJobNotFound(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(nil, fmt.Errorf("not found"))

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddCronJobIDMismatch(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{ObjectMeta: metav1.ObjectMeta{UID: "858558f7-305a-4266-9775-29e7143e6f99"}}, nil)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddSuccessfulCreateAlreadyExists(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Create").Return(nil, ErrAlreadyExists)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "SuccessfulCreate"})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddSuccessfulCreate(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Create").Return(&Job{}, nil)
	h.On("Publish", mock.Anything).Return()

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "SuccessfulCreate"})

	h.AssertCalled(t, "Publish", mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.Payload.(*Job)

		return e.Type == "jobadd" && ok
	}))
}

func TestOnEventAddCompletedAlreadyCompleted(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Succeed").Return(nil, ErrAlreadyCompleted)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "Completed"})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddCompleted(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Succeed").Return(&Job{}, nil)
	h.On("Publish", mock.Anything).Return()

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "Completed"})

	h.AssertCalled(t, "Publish", mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.Payload.(*Job)

		return e.Type == "jobcomplete" && ok
	}))
}

func TestOnEventAddBackoffLimitExceededAlreadyCompleted(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Fail").Return(nil, ErrAlreadyCompleted)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "BackoffLimitExceeded"})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddDeadlineExceededAlreadyCompleted(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Fail").Return(nil, ErrAlreadyCompleted)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "DeadlineExceeded"})

	h.AssertNotCalled(t, "Publish")
}

func TestOnEventAddBackoffLimitExceeded(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Fail").Return(&Job{}, nil)
	h.On("Publish", mock.Anything).Return()

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "BackoffLimitExceeded"})

	h.AssertCalled(t, "Publish", mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.Payload.(*Job)

		return e.Type == "jobcomplete" && ok
	}))
}

func TestOnEventAddDeadlineExceeded(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)
	s.On("Fail").Return(&Job{}, nil)
	h.On("Publish", mock.Anything).Return()

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}, Reason: "DeadlineExceeded"})

	h.AssertCalled(t, "Publish", mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.Payload.(*Job)

		return e.Type == "jobcomplete" && ok
	}))
}

func TestOnEventAddUnknownReason(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}
	jobNamespaceLister := &jobNamespaceListerStub{}
	cronJobNamespaceLister := &cronJobNamespaceListerStub{}

	controller := true

	jobLister.On("Jobs").Return(jobNamespaceLister)
	jobNamespaceLister.On("Get").Return(
		&batchv1.Job{ObjectMeta: metav1.ObjectMeta{OwnerReferences: []metav1.OwnerReference{{Kind: "CronJob", Controller: &controller}}}},
		nil,
	)
	cronJobLister.On("CronJobs").Return(cronJobNamespaceLister)
	cronJobNamespaceLister.On("Get").Return(&batchv1.CronJob{}, nil)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnEventAdd(&corev1.Event{InvolvedObject: corev1.ObjectReference{Kind: "Job"}})

	h.AssertNotCalled(t, "Publish")
}

func TestOnCronJobDelete(t *testing.T) {
	s := &serviceMock{}
	h := &hubMock{}
	cronJobLister := &cronJobListerStub{}
	jobLister := &jobListerStub{}

	s.On("Delete", mock.Anything).Return(nil)

	e := NewEventHandler(s, h, cronJobLister, jobLister)

	e.OnCronJobDelete(&batchv1.CronJob{})

	s.AssertCalled(t, "Delete")
}

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) Create(context.Context, string, string, time.Time) (*Job, error) {
	args := s.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Job), args.Error(1)
}

func (s *serviceMock) Succeed(context.Context, string, time.Time) (*Job, error) {
	args := s.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Job), args.Error(1)
}

func (s *serviceMock) Fail(context.Context, string, time.Time) (*Job, error) {
	args := s.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Job), args.Error(1)
}

func (s *serviceMock) List(context.Context) ([]Job, error) {
	return nil, nil
}

func (s *serviceMock) Delete(context.Context, string) error {
	args := s.Called()

	return args.Error(0)
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

type jobListerStub struct {
	mock.Mock
}

func (j *jobListerStub) List(labels.Selector) ([]*batchv1.Job, error) {
	return nil, nil
}

func (j *jobListerStub) Jobs(string) v1.JobNamespaceLister {
	args := j.Called()

	return args.Get(0).(v1.JobNamespaceLister)
}

func (j *jobListerStub) GetPodJobs(*corev1.Pod) ([]batchv1.Job, error) {
	return nil, nil
}

type jobNamespaceListerStub struct {
	mock.Mock
}

func (j *jobNamespaceListerStub) List(labels.Selector) ([]*batchv1.Job, error) {
	return nil, nil
}

func (j *jobNamespaceListerStub) Get(string) (*batchv1.Job, error) {
	args := j.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*batchv1.Job), args.Error(1)
}

type cronJobListerStub struct {
	mock.Mock
}

func (c *cronJobListerStub) List(labels.Selector) ([]*batchv1.CronJob, error) {
	return nil, nil
}

func (c *cronJobListerStub) CronJobs(string) v1.CronJobNamespaceLister {
	args := c.Called()

	return args.Get(0).(v1.CronJobNamespaceLister)
}

type cronJobNamespaceListerStub struct {
	mock.Mock
}

func (c *cronJobNamespaceListerStub) List(labels.Selector) ([]*batchv1.CronJob, error) {
	return nil, nil
}

func (c *cronJobNamespaceListerStub) Get(string) (*batchv1.CronJob, error) {
	args := c.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*batchv1.CronJob), args.Error(1)
}
