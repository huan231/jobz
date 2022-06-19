package cronjob

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSuccessResponse(t *testing.T) {
	s := &serviceStub{cronJobs: []CronJob{{}}}

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
	err := fmt.Errorf("unexpected error")

	s := &serviceStub{err: err}

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
	err      error
	cronJobs []CronJob
}

func (s *serviceStub) List() ([]CronJob, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.cronJobs, nil
}
