package cronjob

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/batch/v1"
	"testing"
)

func TestList(t *testing.T) {
	l := &listerMock{t: t, cronJobs: []*batchv1.CronJob{{}}}

	s := NewService(l)

	cronJobs, _ := s.List()

	assert.NotEmpty(t, cronJobs)
}

func TestListError(t *testing.T) {
	l := &listerMock{err: fmt.Errorf("unexpected error")}

	s := NewService(l)

	_, err := s.List()

	assert.NotNil(t, err)
}

type listerMock struct {
	err      error
	t        *testing.T
	cronJobs []*batchv1.CronJob
}

func (l *listerMock) List(selector labels.Selector) ([]*batchv1.CronJob, error) {
	if l.err != nil {
		return nil, l.err
	}

	assert.Equal(l.t, labels.Everything(), selector)

	return l.cronJobs, nil
}

func (l *listerMock) CronJobs(string) v1.CronJobNamespaceLister {
	return nil
}
