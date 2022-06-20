package job

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateSuccess(t *testing.T) {
	id, cronJobID := "a4757ca2-f45c-44ee-9db3-cb258d606b77", "99013750-5341-4dc5-bea7-6e3ae807a475"

	r := &repositoryStruct{ok: true}

	s := NewService(r)

	job, _ := s.Create(context.Background(), id, cronJobID, time.Time{})

	assert.Equal(t, id, job.ID)
	assert.Equal(t, cronJobID, job.CronJobID)
	assert.Equal(t, Running, job.Status)
	assert.Equal(t, time.Time{}, job.CreatedAt)
	assert.Equal(t, time.Time{}, job.UpdatedAt)
	assert.Nil(t, job.CompletedAt)
}

func TestCreateAlreadyExists(t *testing.T) {
	r := &repositoryStruct{ok: false}

	s := NewService(r)

	_, err := s.Create(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", "99013750-5341-4dc5-bea7-6e3ae807a475", time.Time{})

	assert.Equal(t, ErrAlreadyExists, err)
}

func TestCreateError(t *testing.T) {
	r := &repositoryStruct{err: fmt.Errorf("unexpected error")}

	s := NewService(r)

	_, err := s.Create(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", "99013750-5341-4dc5-bea7-6e3ae807a475", time.Time{})

	assert.NotNil(t, err)
}

func TestSucceedSuccess(t *testing.T) {
	r := &repositoryStruct{ok: true, job: &Job{}}

	s := NewService(r)

	job, _ := s.Succeed(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, job)
}

func TestSucceedAlreadyCompleted(t *testing.T) {
	r := &repositoryStruct{ok: false}

	s := NewService(r)

	_, err := s.Succeed(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.Equal(t, ErrAlreadyCompleted, err)
}

func TestSucceedError(t *testing.T) {
	r := &repositoryStruct{err: fmt.Errorf("unexpected error")}

	s := NewService(r)

	_, err := s.Succeed(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, err)
}

func TestFailSuccess(t *testing.T) {
	r := &repositoryStruct{ok: true, job: &Job{}}

	s := NewService(r)

	job, _ := s.Fail(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, job)
}

func TestFailAlreadyCompleted(t *testing.T) {
	r := &repositoryStruct{ok: false}

	s := NewService(r)

	_, err := s.Fail(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.Equal(t, ErrAlreadyCompleted, err)
}

func TestFailError(t *testing.T) {
	r := &repositoryStruct{err: fmt.Errorf("unexpected error")}

	s := NewService(r)

	_, err := s.Fail(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, err)
}

func TestListSuccess(t *testing.T) {
	r := &repositoryStruct{jobs: []Job{{}}}

	s := NewService(r)

	jobs, _ := s.List(context.Background())

	assert.NotEmpty(t, jobs)
}

func TestListError(t *testing.T) {
	r := &repositoryStruct{err: fmt.Errorf("unexpected error")}

	s := NewService(r)

	_, err := s.List(context.Background())

	assert.NotNil(t, err)
}

func TestDeleteSuccess(t *testing.T) {
	r := &repositoryStruct{}

	s := NewService(r)

	err := s.Delete(context.Background(), "99013750-5341-4dc5-bea7-6e3ae807a475")

	assert.Nil(t, err)
}

func TestDeleteError(t *testing.T) {
	r := &repositoryStruct{err: fmt.Errorf("unexpected error")}

	s := NewService(r)

	err := s.Delete(context.Background(), "99013750-5341-4dc5-bea7-6e3ae807a475")

	assert.NotNil(t, err)
}

type repositoryStruct struct {
	ok   bool
	job  *Job
	jobs []Job
	err  error
}

func (r *repositoryStruct) Add(context.Context, Job) (bool, error) {
	return r.ok, r.err
}

func (r *repositoryStruct) List(context.Context) ([]Job, error) {
	return r.jobs, r.err
}

func (r *repositoryStruct) CompleteByID(context.Context, string, Status, time.Time) (bool, error) {
	return r.ok, r.err
}

func (r *repositoryStruct) DeleteByCronJobID(context.Context, string) error {
	return r.err
}

func (r *repositoryStruct) FindByID(context.Context, string) (*Job, error) {
	return r.job, r.err
}
