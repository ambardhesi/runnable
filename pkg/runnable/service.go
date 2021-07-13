package runnable

import (
	"io"
)

type JobService interface {
	Start(ownerID string, command string, args ...string) (string, error)
	Stop(ownerID string, jobID string) error
	Get(ownerID string, jobID string) (*Job, error)
	GetLogs(ownerID string, jobID string) (*string, error)
}

type JobStoreService interface {
	Store(job *Job) error
	Get(jobID string) (*Job, bool)
}

type LogFileService interface {
	CreateLogFile(jobID string) (io.WriteCloser, error)
	GetLogFile(jobID string) (io.ReadCloser, error)
	DeleteAllLogFiles() error
}
