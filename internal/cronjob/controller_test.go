package cronjob

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSuccessResponse(t *testing.T) {
	s := &serviceStub{}
	s.On("List").Return([]CronJob{{}}, nil)

	c := NewController(s)

	r := gin.New()
	r.GET("/cronjobs", c.List)

	req, _ := http.NewRequest("GET", "/cronjobs", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var cronJobs []CronJob

	json.Unmarshal(w.Body.Bytes(), &cronJobs)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, cronJobs)
}

func TestListErrorResponse(t *testing.T) {
	s := &serviceStub{}
	s.On("List").Return(nil, fmt.Errorf("unexpected error"))

	c := NewController(s)

	r := gin.New()
	r.GET("/cronjobs", c.List)

	req, _ := http.NewRequest("GET", "/cronjobs", nil)
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

func (s *serviceStub) List() ([]CronJob, error) {
	args := s.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]CronJob), args.Error(1)
}
