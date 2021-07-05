package job_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/ambardhesi/runnable/pkg/job"
	"github.com/ambardhesi/runnable/pkg/runnable"
	"github.com/ambardhesi/runnable/pkg/runnable/mock"
)

func TestStart(t *testing.T) {
	buf := &bytes.Buffer{}

	lfs := &mock.MockLogFileService{
		Buf: buf,
	}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	jobID, err := js.Start("ownerID", "echo", "hello world")

	time.Sleep(200 * time.Millisecond)

	if err != nil {
		t.Errorf("expected no errors, got %v", err)
	}

	if jobID == "" {
		t.Errorf("expected non empty jobID")
	}

	if jss.Job.State() != runnable.Running && jss.Job.State() != runnable.Completed {
		t.Errorf("expected state to be running/completed, got %v", jss.Job.State())
	}
}

func TestStop(t *testing.T) {
	buf := &bytes.Buffer{}

	lfs := &mock.MockLogFileService{
		Buf: buf,
	}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	jobID, _ := js.Start("ownerID", "sleep", "5")
	err := js.Stop("ownerID", jobID)

	if err != nil {
		t.Errorf("expected no errors, got %v", err)
	}

	if jss.Job.State() != runnable.Stopped {
		t.Errorf("expected job to have stopped")
	}
}

func TestStopJobDoesNotExist(t *testing.T) {
	lfs := &mock.MockLogFileService{}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	err := js.Stop("ownerID", "jobID")

	if runnable.ErrorCode(err) != runnable.ENOTFOUND {
		t.Errorf("expected error type %v, got error%v", runnable.ENOTFOUND, err)
	}

}

func TestGet(t *testing.T) {
	buf := &bytes.Buffer{}

	lfs := &mock.MockLogFileService{
		Buf: buf,
	}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	jobID, _ := js.Start("ownerID", "echo", "hello world")

	if job, _ := js.Get("ownerID", jobID); job != jss.Job {
		t.Errorf("expected job %v, got %v", jss.Job, job)
	}
}

func TestGetJobDoesNotExist(t *testing.T) {
	lfs := &mock.MockLogFileService{}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	_, err := js.Get("ownerID", "jobID")

	if runnable.ErrorCode(err) != runnable.ENOTFOUND {
		t.Errorf("expected error type %v, got error%v", runnable.ENOTFOUND, err)
	}
}

func TestGetLogs(t *testing.T) {
	buf := &bytes.Buffer{}

	lfs := &mock.MockLogFileService{
		Buf: buf,
	}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	jobID, _ := js.Start("ownerID", "echo", "hello world")
	time.Sleep(200 * time.Millisecond)

	logs, err := js.GetLogs("ownerID", jobID)

	if err != nil {
		t.Errorf("expected no errors, got %v", err)
	}

	if !strings.Contains(*logs, "hello world") {
		t.Errorf("expected logs to contain %v, got %v", "hello world", logs)
	}
}

func TestGetLogsJobDoesNotExist(t *testing.T) {
	lfs := &mock.MockLogFileService{}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	_, err := js.GetLogs("ownerID", "jobID")

	if runnable.ErrorCode(err) != runnable.ENOTFOUND {
		t.Errorf("expected error type %v, got error%v", runnable.ENOTFOUND, err)
	}

}
