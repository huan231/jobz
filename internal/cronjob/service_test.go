package cronjob

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/batch/v1"
	"testing"
)

func TestList(t *testing.T) {
	l := &listerMock{}
	l.On("List", mock.Anything).Return([]*batchv1.CronJob{{}}, nil)

	s := NewService(l)

	cronJobs, _ := s.List()

	l.AssertCalled(t, "List", labels.Everything())

	assert.NotEmpty(t, cronJobs)
}

func TestListError(t *testing.T) {
	l := &listerMock{}
	l.On("List", mock.Anything).Return(nil, fmt.Errorf("unexpected error"))

	s := NewService(l)

	_, err := s.List()

	assert.NotNil(t, err)
}

type listerMock struct {
	mock.Mock
}

func (l *listerMock) List(selector labels.Selector) ([]*batchv1.CronJob, error) {
	args := l.Called(selector)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*batchv1.CronJob), args.Error(1)
}

func (l *listerMock) CronJobs(string) v1.CronJobNamespaceLister {
	return nil
}
