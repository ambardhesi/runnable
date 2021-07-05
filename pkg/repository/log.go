package repository

import (
	"io"
	"os"
	"path"

	"github.com/ambardhesi/runnable/pkg/runnable"
)

const (
	userReadWriteExecutePermission = 0700
)

// implementation of runnable.LogFileService
type LocalFileSystem struct {
	dir string
}

func NewLocalFileSystem(dir string) (*LocalFileSystem, error) {
	// create a dir for storing logs, with relevant permissions
	err := checkAndCreateLogsDir(dir)
	if err != nil {
		return nil, &runnable.Error{
			Message: "Failed to create log dir.",
			Err:     err,
		}
	}

	return &LocalFileSystem{
		dir: dir,
	}, nil
}

// if the log dir doesn't exist, create it
func checkAndCreateLogsDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, userReadWriteExecutePermission)
	}
	return nil
}

func (lfs *LocalFileSystem) CreateLogFile(jobID string) (io.Writer, error) {
	checkAndCreateLogsDir(lfs.dir)

	logFile, err := os.Create(path.Join(lfs.dir, jobID))
	if err != nil {
		return nil, &runnable.Error{
			Code:    runnable.EINTERNAL,
			Op:      "LocalFileSystem.CreateLogFile",
			Message: "Failed to create log file.",
			Err:     err,
		}
	}

	return logFile, nil
}

func (lfs *LocalFileSystem) GetLogFile(jobID string) (io.ReadCloser, error) {
	logFile, err := os.Open(path.Join(lfs.dir, jobID))
	if err != nil {
		return nil, &runnable.Error{
			Code:    runnable.EINTERNAL,
			Op:      "LocalFileSystem.GetLogFile",
			Message: "Failed to open log file",
			Err:     err,
		}
	}

	return logFile, nil
}

func (lfs *LocalFileSystem) DeleteAllLogs() error {
	err := os.RemoveAll(lfs.dir)
	if err != nil {
		return &runnable.Error{
			Code:    runnable.EINTERNAL,
			Op:      "LocalFileSystem.DeleteAllLogs",
			Message: "Failed to delete log dir",
			Err:     err,
		}
	}

	return nil
}
