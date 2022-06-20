package cronjob

import (
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewCronJob(t *testing.T) {
	cronJob := NewCronJob(&batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{UID: "31c51d6c-f422-4a7e-bf4c-1e5ed93fea45", Namespace: "default", Name: "hello"},
		Spec:       batchv1.CronJobSpec{Schedule: "* * * * *"},
	})

	assert.Equal(t, "31c51d6c-f422-4a7e-bf4c-1e5ed93fea45", cronJob.ID)
	assert.Equal(t, "default", cronJob.Namespace)
	assert.Equal(t, "hello", cronJob.Name)
	assert.Equal(t, "* * * * *", cronJob.Schedule)
}
