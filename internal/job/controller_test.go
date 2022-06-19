package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListSuccessResponse(t *testing.T) {
	s := &serviceStub{jobs: []Job{{Status: Running}}}

	c := NewController(s)

	r := gin.New()
	r.GET("/jobs", c.List)

	req, _ := http.NewRequest("GET", "/jobs", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var jobs []Job

	json.Unmarshal(w.Body.Bytes(), &jobs)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, jobs)
}

func TestListErrorResponse(t *testing.T) {
	err := fmt.Errorf("unexpected error")

	s := &serviceStub{err: err}

	c := NewController(s)

	r := gin.New()
	r.GET("/jobs", c.List)

	req, _ := http.NewRequest("GET", "/jobs", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var body any

	json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.True(t, assert.ObjectsAreEqualValues(map[string]any{"message": "unexpected error"}, body))
}

type serviceStub struct {
	err  error
	jobs []Job
}

func (s *serviceStub) Create(context.Context, string, string, time.Time) (*Job, error) {
	return nil, nil
}

func (s *serviceStub) Succeed(context.Context, string, time.Time) (*Job, error) {
	return nil, nil
}

func (s *serviceStub) Fail(context.Context, string, time.Time) (*Job, error) {
	return nil, nil
}

func (s *serviceStub) List(context.Context) ([]Job, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.jobs, nil
}

func (s *serviceStub) Delete(context.Context, string) error {
	return nil
}
