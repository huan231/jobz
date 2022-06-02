package job

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAlreadyExists    = errors.New("already exists")
	ErrAlreadyCompleted = errors.New("already completed")
)

type Service interface {
	Create(ctx context.Context, id string, cronJobID string, t time.Time) (*Job, error)
	Succeed(ctx context.Context, id string, t time.Time) (*Job, error)
	Fail(ctx context.Context, id string, t time.Time) (*Job, error)
	List(ctx context.Context) ([]Job, error)
	Delete(ctx context.Context, cronJobID string) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(ctx context.Context, id string, cronJobID string, t time.Time) (*Job, error) {
	job := Job{ID: id, CronJobID: cronJobID, Status: Running, CreatedAt: t, UpdatedAt: t}

	ok, err := s.r.Add(ctx, job)

	if err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrAlreadyExists
	}

	return &job, nil
}

func (s *service) Succeed(ctx context.Context, id string, t time.Time) (*Job, error) {
	return s.complete(ctx, id, Succeeded, t)
}

func (s *service) Fail(ctx context.Context, id string, t time.Time) (*Job, error) {
	return s.complete(ctx, id, Failed, t)
}

func (s *service) complete(ctx context.Context, id string, status Status, t time.Time) (*Job, error) {
	ok, err := s.r.CompleteByID(ctx, id, status, t)

	if err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrAlreadyCompleted
	}

	return s.r.FindByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]Job, error) {
	return s.r.List(ctx)
}

func (s *service) Delete(ctx context.Context, cronJobID string) error {
	return s.r.DeleteByCronJobID(ctx, cronJobID)
}
