package cronjob

import (
	batchv1 "k8s.io/api/batch/v1"
)

type CronJob struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Schedule  string `json:"schedule"`
}

func NewCronJob(cronJob *batchv1.CronJob) *CronJob {
	return &CronJob{ID: string(cronJob.UID), Namespace: cronJob.Namespace, Name: cronJob.Name, Schedule: cronJob.Spec.Schedule}
}
