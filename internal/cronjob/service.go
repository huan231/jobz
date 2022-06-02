package cronjob

import (
	"k8s.io/apimachinery/pkg/labels"
	batchv1 "k8s.io/client-go/listers/batch/v1"
)

type Service interface {
	List() ([]CronJob, error)
}

type service struct {
	lister batchv1.CronJobLister
}

func NewService(lister batchv1.CronJobLister) Service {
	return &service{lister}
}

func (s *service) List() ([]CronJob, error) {
	ret, err := s.lister.List(labels.Everything())

	if err != nil {
		return nil, err
	}

	cronJobs := make([]CronJob, len(ret))

	for i, cronJob := range ret {
		cronJobs[i] = *NewCronJob(cronJob)
	}

	return cronJobs, nil
}
