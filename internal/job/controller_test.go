package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListSuccessResponse(t *testing.T) {
	s := &serviceStub{}
	s.On("List").Return([]Job{{Status: Running}}, nil)

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
	s := &serviceStub{}
	s.On("List").Return(nil, fmt.Errorf("unexpected error"))

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
	mock.Mock
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
	args := s.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Job), args.Error(1)
}

func (s *serviceStub) Delete(context.Context, string) error {
	return nil
}
