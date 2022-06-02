package job

import (
	"encoding/json"
	"errors"
	"time"
)

type Status int

const (
	Running Status = iota + 1
	Succeeded
	Failed
)

func (s Status) MarshalJSON() ([]byte, error) {
	var v string

	switch s {
	case Running:
		v = "running"
	case Succeeded:
		v = "succeeded"
	case Failed:
		v = "failed"
	default:
		return nil, errors.New("invalid value")
	}

	return json.Marshal(v)
}

type Job struct {
	ID          string     `json:"id"`
	CronJobID   string     `json:"cronJobId"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}
