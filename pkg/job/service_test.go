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

func TestEndToEnd(t *testing.T) {
	buf := &bytes.Buffer{}

	lfs := &mock.MockLogFileService{
		Buf: buf,
	}
	jss := &mock.MockJobStoreService{}
	js := job.NewJobService(jss, lfs)

	// job sleeps for 2 seconds to give us time to do some assertions and to stop it
	jobID, err := js.Start("ownerID", "sleep", "2")
	if err != nil {
		t.Errorf("Did not expect to get an error starting job")
	}

	job, err := js.Get("ownerID", jobID)
	if err != nil {
		t.Errorf("Did not expect to get an error fetching job with id %v", jobID)
	}

	var state runnable.State
	for state = job.Status().State; state != runnable.Running; state = job.Status().State {
	}

	err = js.Stop("ownerID", jobID)
	time.Sleep(200 * time.Millisecond)

	if err != nil {
		t.Errorf("expected no errors, got %v", err)
	}

	if job.Status().State != runnable.Stopped {
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
