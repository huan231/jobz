package job

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreateSuccess(t *testing.T) {
	id, cronJobID := "a4757ca2-f45c-44ee-9db3-cb258d606b77", "99013750-5341-4dc5-bea7-6e3ae807a475"

	r := &repositoryStruct{}
	r.On("Add").Return(true, nil)

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
	r := &repositoryStruct{}
	r.On("Add").Return(false, nil)

	s := NewService(r)

	_, err := s.Create(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", "99013750-5341-4dc5-bea7-6e3ae807a475", time.Time{})

	assert.Equal(t, ErrAlreadyExists, err)
}

func TestCreateError(t *testing.T) {
	r := &repositoryStruct{}
	r.On("Add").Return(false, fmt.Errorf("unexpected error"))

	s := NewService(r)

	_, err := s.Create(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", "99013750-5341-4dc5-bea7-6e3ae807a475", time.Time{})

	assert.NotNil(t, err)
}

func TestSucceedSuccess(t *testing.T) {
	r := &repositoryStruct{}
	r.On("CompleteByID").Return(true, nil)
	r.On("FindByID").Return(&Job{}, nil)

	s := NewService(r)

	job, _ := s.Succeed(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, job)
}

func TestSucceedAlreadyCompleted(t *testing.T) {
	r := &repositoryStruct{}
	r.On("CompleteByID").Return(false, nil)

	s := NewService(r)

	_, err := s.Succeed(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.Equal(t, ErrAlreadyCompleted, err)
}

func TestSucceedError(t *testing.T) {
	r := &repositoryStruct{}
	r.On("CompleteByID").Return(false, fmt.Errorf("unexpected error"))

	s := NewService(r)

	_, err := s.Succeed(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, err)
}

func TestFailSuccess(t *testing.T) {
	r := &repositoryStruct{}
	r.On("CompleteByID").Return(true, nil)
	r.On("FindByID").Return(&Job{}, nil)

	s := NewService(r)

	job, _ := s.Fail(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, job)
}

func TestFailAlreadyCompleted(t *testing.T) {
	r := &repositoryStruct{}
	r.On("CompleteByID").Return(false, nil)

	s := NewService(r)

	_, err := s.Fail(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.Equal(t, ErrAlreadyCompleted, err)
}

func TestFailError(t *testing.T) {
	r := &repositoryStruct{}
	r.On("CompleteByID").Return(false, fmt.Errorf("unexpected error"))

	s := NewService(r)

	_, err := s.Fail(context.Background(), "a4757ca2-f45c-44ee-9db3-cb258d606b77", time.Time{})

	assert.NotNil(t, err)
}

func TestListSuccess(t *testing.T) {
	r := &repositoryStruct{}
	r.On("List").Return([]Job{{}}, nil)

	s := NewService(r)

	jobs, _ := s.List(context.Background())

	assert.NotEmpty(t, jobs)
}

func TestListError(t *testing.T) {
	r := &repositoryStruct{}
	r.On("List").Return(nil, fmt.Errorf("unexpected error"))

	s := NewService(r)

	_, err := s.List(context.Background())

	assert.NotNil(t, err)
}

func TestDeleteSuccess(t *testing.T) {
	r := &repositoryStruct{}
	r.On("DeleteByCronJobID").Return(nil)

	s := NewService(r)

	err := s.Delete(context.Background(), "99013750-5341-4dc5-bea7-6e3ae807a475")

	assert.Nil(t, err)
}

func TestDeleteError(t *testing.T) {
	r := &repositoryStruct{}
	r.On("DeleteByCronJobID").Return(fmt.Errorf("unexpected error"))

	s := NewService(r)

	err := s.Delete(context.Background(), "99013750-5341-4dc5-bea7-6e3ae807a475")

	assert.NotNil(t, err)
}

type repositoryStruct struct {
	mock.Mock
}

func (r *repositoryStruct) Add(context.Context, Job) (bool, error) {
	args := r.Called()

	return args.Bool(0), args.Error(1)
}

func (r *repositoryStruct) List(context.Context) ([]Job, error) {
	args := r.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Job), args.Error(1)
}

func (r *repositoryStruct) CompleteByID(context.Context, string, Status, time.Time) (bool, error) {
	args := r.Called()

	return args.Bool(0), args.Error(1)
}

func (r *repositoryStruct) DeleteByCronJobID(context.Context, string) error {
	args := r.Called()

	return args.Error(0)
}

func (r *repositoryStruct) FindByID(context.Context, string) (*Job, error) {
	args := r.Called()

	return args.Get(0).(*Job), args.Error(1)
}
