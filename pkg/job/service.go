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

func (jobSvc *JobService) Start(ownerID string, command string, args ...string) (jobID string, err error) {
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

	err = job.Start()
	if err != nil {
		// TODO delete the job from the job store and delete log file
		// similar rationale to above
		return "", err
	}

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

	err := job.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (jobSvc *JobService) Get(ownerID string, jobID string) (job *runnable.Job, err error) {
	job, exists := jobSvc.jobStoreSvc.Get(jobID)
	if !exists {
		return nil, &runnable.Error{
			Code:    runnable.ENOTFOUND,
			Op:      "JobService.Get",
			Message: "Job does not exist.",
		}
	}

	return job, nil
}

func (jobSvc *JobService) GetLogs(ownerID string, jobID string) (logs *string, err error) {
	_, exists := jobSvc.jobStoreSvc.Get(jobID)
	if !exists {
		return nil, &runnable.Error{
			Code:    runnable.ENOTFOUND,
			Op:      "JobService.GetLogs",
			Message: "Job does not exist.",
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
