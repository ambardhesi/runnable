package runnable

import (
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"
)

type State string

// All the possible statuses of a job.
const (
	NotStarted State = "NotStarted"
	Running          = "Running"
	Failed           = "Failed"
	Completed        = "Completed"
	Stopped          = "Stopped"
)

type Status struct {
	State     State
	StartTime time.Time
	EndTime   time.Time
	ExitCode  int
}

type Job struct {
	ID        string
	Cmd       *exec.Cmd
	OwnerID   string
	status    Status
	logWriter io.WriteCloser
	lock      sync.RWMutex
}

// Creates a new job for a given command and owner ID.
// Job will have a state of NotStarted and a new UUID as its ID.
func NewJob(ownerID string, command string, args ...string) (*Job, error) {
	if command == "" {
		return nil, &Error{
			Code:    EINVALID,
			Op:      "Job.NewJob",
			Message: "Command is empty.",
		}
	}

	cmd := exec.Command(command, args...)
	jobID := uuid.NewString()

	status := Status{
		State:    NotStarted,
		ExitCode: -1,
	}

	return &Job{
		ID:      jobID,
		Cmd:     cmd,
		OwnerID: ownerID,
		status:  status,
	}, nil

}

func (job *Job) SetLogWriter(wc io.WriteCloser) {
	job.logWriter = wc
}

// get job's state
func (job *Job) Status() Status {
	job.lock.RLock()
	defer job.lock.RUnlock()
	return job.status
}

// Runs the job by calling Cmd.Run()
// A goroutine will wait for the job to finish executing.
// Returns InvalidStateError if the job is not in a NotStarted state.
func (job *Job) Start() error {
	op := "JobService.Start"
	if job.Status().State != NotStarted {
		return &Error{
			Code:    EINVALID,
			Op:      op,
			Message: "Job is not in a NotStarted state.",
		}
	}

	stdout, err := job.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := job.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	logOutput := io.MultiReader(stdout, stderr)

	err = job.Cmd.Start()

	if err != nil {
		return &Error{
			Code:    EINTERNAL,
			Op:      op,
			Message: "Failed to start job.",
			Err:     err,
		}
	}

	job.lock.Lock()
	job.status.State = Running
	job.status.StartTime = time.Now()
	job.lock.Unlock()

	_, err = io.Copy(job.logWriter, logOutput)
	if err != nil {
		return err
	}

	go func() {
		err := job.wait()
		if err != nil {
			// task failed, set status to Failed
			job.lock.Lock()
			defer job.lock.Unlock()

			job.status.State = Failed
			job.status.EndTime = time.Now()
		}
	}()

	return nil
}

// Waits for the wrapped process to finish before updating exit code (else there will be a race on Cmd)
// Also updates the job state and end time.
func (job *Job) wait() error {
	defer job.logWriter.Close()

	var exitCode int

	switch err := job.Cmd.Wait().(type) {
	case nil:
		// job completed successfully
		exitCode = job.Cmd.ProcessState.ExitCode()

	case *exec.ExitError:
		// job exited with an exit code
		exitCode = err.ProcessState.ExitCode()

	default:
		// job failed
		return err
	}

	job.lock.Lock()
	defer job.lock.Unlock()

	// if job has already failed or completed, do nothing as we should have already updated status before
	if job.status.State == Failed || job.status.State == Completed {
		return nil
	}

	job.status.ExitCode = exitCode
	job.status.EndTime = time.Now()
	// Update the job's state and end time.
	// However, we should check first if it wasnt already stopped (by the user).
	if job.status.State != Stopped {
		job.status.State = Completed
	}

	return nil
}

// Stops the job.
// Returns InvalidStateError if the job is not currently Running.
func (job *Job) Stop() error {
	job.lock.Lock()
	defer job.lock.Unlock()

	op := "JobService.Stop"
	if job.status.State != Running {
		return &Error{
			Code:    EINVALID,
			Op:      op,
			Message: "Job is not in a Running state.",
		}
	}

	err := job.Cmd.Process.Kill()

	if err != nil {
		return &Error{
			Code:    EINTERNAL,
			Op:      op,
			Message: "Failed to stop job.",
			Err:     err,
		}
	}

	job.status.State = Stopped

	return nil
}
