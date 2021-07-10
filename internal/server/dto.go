package server

import (
	"time"

	"github.com/ambardhesi/runnable/pkg/runnable"
)

type StartJobRequest struct {
	Command string `json:"command" binding:"required"`
}

type StartJobResponse struct {
	JobID string `json:"jobID"`
}

type GetJobResponse struct {
	State     string    `json:"state"`
	ExitCode  int       `json:"exitCode"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func FromJob(job *runnable.Job) GetJobResponse {
	status := job.Status()
	return GetJobResponse{
		State:     string(status.State),
		ExitCode:  status.ExitCode,
		StartTime: status.StartTime,
		EndTime:   status.EndTime,
	}
}
