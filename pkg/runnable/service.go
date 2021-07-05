package runnable

import (
	"io"
)

type JobService interface {
	Start(ownerID string, command string, args ...string) (jobID string, err error)
	Stop(ownerID string, jobID string) error
	Get(ownerID string, jobID string) (job *Job, err error)
	GetLogs(ownerID string, jobID string) (logs *string, err error)
}

type JobStoreService interface {
	Store(job *Job) error
	Get(jobID string) (job *Job, exists bool)
}

type LogFileService interface {
	CreateLogFile(jobID string) (writer io.Writer, err error)
	GetLogFile(jobID string) (readCloser io.ReadCloser, err error)
	DeleteAllLogFiles() error
}
