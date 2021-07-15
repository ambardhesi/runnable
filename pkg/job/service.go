package job

import (
	"io"

	"github.com/ambardhesi/runnable/pkg/runnable"
)

// implementation of runnable.JobService
type JobService struct {
	jobStoreSvc runnable.JobStoreService
	logFileSvc  runnable.LogFileService
}

func NewJobService(jobStoreSvc runnable.JobStoreService, logFileSvc runnable.LogFileService) *JobService {
	return &JobService{
		jobStoreSvc: jobStoreSvc,
		logFileSvc:  logFileSvc,
	}
}

func (jobSvc *JobService) Start(ownerID string, command string, args ...string) (string, error) {
	job, err := runnable.NewJob(ownerID, command, args...)
	if err != nil {
		return "", err
	}

	logFile, err := jobSvc.logFileSvc.CreateLogFile(job.ID)
	if err != nil {
		return "", err
	}

	job.SetLogWriter(logFile)

	err = jobSvc.jobStoreSvc.Store(job)
	if err != nil {
		// TODO delete the created log file for this job
		// all log files will be deleted on server shutdown, so skipping this for now
		return "", err
	}

	go func() {
		job.Start()
	}()

	return job.ID, nil

}

func (jobSvc *JobService) Stop(ownerID string, jobID string) error {
	job, exists := jobSvc.jobStoreSvc.Get(jobID)
	if !exists {
		return &runnable.Error{
			Code:    runnable.ENOTFOUND,
			Op:      "JobService.Stop",
			Message: "Job does not exist.",
		}
	}

	// TODO move auth into its own service
	// something that takes in a jobID and ownerID and looks up whatever it needs to to make a decision
	// this could be something like a DB call (via job service), or an external service call
	if job.OwnerID != ownerID {
		return &runnable.Error{
			Code:    runnable.EUNAUTHORIZED,
			Op:      "JobService.Stop",
			Message: "User is unauthorized.",
		}
	}

	err := job.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (jobSvc *JobService) Get(ownerID string, jobID string) (*runnable.Job, error) {
	job, exists := jobSvc.jobStoreSvc.Get(jobID)
	if !exists {
		return nil, &runnable.Error{
			Code:    runnable.ENOTFOUND,
			Op:      "JobService.Get",
			Message: "Job does not exist.",
		}
	}

	// TODO move auth into its own service
	// something that takes in a jobID and ownerID and looks up whatever it needs to to make a decision
	// this could be something like a DB call (via job service), or an external service call
	if job.OwnerID != ownerID {
		return nil, &runnable.Error{
			Code:    runnable.EUNAUTHORIZED,
			Op:      "JobService.Get",
			Message: "User is unauthorized.",
		}
	}

	return job, nil
}

func (jobSvc *JobService) GetLogs(ownerID string, jobID string) (*string, error) {
	job, exists := jobSvc.jobStoreSvc.Get(jobID)
	if !exists {
		return nil, &runnable.Error{
			Code:    runnable.ENOTFOUND,
			Op:      "JobService.GetLogs",
			Message: "Job does not exist.",
		}
	}
	// TODO move auth into its own service
	// something that takes in a jobID and ownerID and looks up whatever it needs to to make a decision
	// this could be something like a DB call (via job service), or an external service call
	if job.OwnerID != ownerID {
		return nil, &runnable.Error{
			Code:    runnable.EUNAUTHORIZED,
			Op:      "JobService.GetLogs",
			Message: "User is unauthorized.",
		}
	}

	file, err := jobSvc.logFileSvc.GetLogFile(jobID)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// TODO we assume here that the file can fit in memory.
	// Ideally we should read in chunks, and implement streaming

	b, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	content := string(b)
	return &content, nil
}
