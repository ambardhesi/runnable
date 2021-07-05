package mock

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"github.com/ambardhesi/runnable/pkg/runnable"
)

type MockLogFileService struct {
	Buf   *bytes.Buffer
	JobID string
}

func (mockLfs *MockLogFileService) CreateLogFile(jobID string) (writer io.Writer, err error) {
	mockLfs.JobID = jobID
	return mockLfs.Buf, nil
}

func (mockLfs *MockLogFileService) GetLogFile(jobID string) (readCloser io.ReadCloser, err error) {
	if mockLfs.JobID == jobID {
		return ioutil.NopCloser(mockLfs.Buf), nil
	} else {
		return nil, errors.New("error")
	}
}

func (mockLfs *MockLogFileService) DeleteAllLogFiles() error {
	return nil
}

type MockJobStoreService struct {
	Job *runnable.Job
}

func (mockJss *MockJobStoreService) Store(job *runnable.Job) error {
	mockJss.Job = job
	return nil
}

func (mockJss *MockJobStoreService) Get(jobID string) (job *runnable.Job, exists bool) {
	if mockJss.Job == nil {
		return nil, false
	} else if mockJss.Job.ID == jobID {
		return mockJss.Job, true
	} else {
		return nil, false
	}
}
