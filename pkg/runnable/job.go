package runnable

import (
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"
)

type State string

// All the possible states of a job.
const (
	NotStarted State = "NotStarted"
	Running          = "Running"
	Failed           = "Failed"
	Completed        = "Completed"
	Stopped          = "Stopped"
)

type Job struct {
	ID        string
	Cmd       *exec.Cmd
	OwnerID   string
	startTime time.Time
	endTime   time.Time
	state     State
	exitCode  int
	logWriter io.Writer
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

	return &Job{
		ID:       jobID,
		Cmd:      cmd,
		OwnerID:  ownerID,
		state:    NotStarted,
		exitCode: -1,
	}, nil

}

func (job *Job) SetLogWriter(writer io.Writer) {
	job.logWriter = writer
}

// get job's state
func (job *Job) State() State {
	job.lock.RLock()
	defer job.lock.RUnlock()
	return job.state
}

// get job's start time
func (job *Job) StartTime() time.Time {
	job.lock.RLock()
	defer job.lock.RUnlock()
	return job.startTime
}

// get job's end time
func (job *Job) EndTime() time.Time {
	job.lock.RLock()
	defer job.lock.RUnlock()
	return job.endTime
}

// get job's exit code
func (job *Job) ExitCode() int {
	job.lock.RLock()
	defer job.lock.RUnlock()
	return job.exitCode
}

// Runs the job by calling Cmd.Run()
// A goroutine will wait for the job to finish executing.
// Returns InvalidStateError if the job is not in a NotStarted state.
func (job *Job) Start() error {
	job.lock.Lock()
	defer job.lock.Unlock()

	op := "JobService.Start"
	if job.state != NotStarted {
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

	job.state = Running
	job.startTime = time.Now()

	_, err = io.Copy(job.logWriter, logOutput)
	if err != nil {
		return err
	}

	go job.wait()

	return nil
}

// Waits for the wrapped process to finish before updating exit code (else there will be a race on Cmd)
// Also updates the job state and end time.
func (job *Job) wait() {
	err := job.Cmd.Wait()

	job.lock.Lock()
	defer job.lock.Unlock()

	job.exitCode = job.Cmd.ProcessState.ExitCode()

	// update the job's state based on if it succeeded or failed
	// however, we should check first if it wasnt already stopped (by the user)
	if job.state == Stopped {
		return
	}

	var state State
	if err != nil {
		state = Failed
	} else {
		state = Completed
	}

	job.state = state
	job.endTime = time.Now()
}

// Stops the job.
// Returns InvalidStateError if the job is not currently Running.
func (job *Job) Stop() error {
	op := "JobService.Stop"
	if job.State() != Running {
		return &Error{
			Code:    EINVALID,
			Op:      op,
			Message: "Job is not in a Running state.",
		}
	}

	err := job.Cmd.Process.Kill()
	job.lock.Lock()
	defer job.lock.Unlock()

	if err != nil {
		return &Error{
			Code:    EINTERNAL,
			Op:      op,
			Message: "Failed to stop job.",
			Err:     err,
		}
	}

	job.endTime = time.Now()
	job.state = Stopped

	return nil
}
