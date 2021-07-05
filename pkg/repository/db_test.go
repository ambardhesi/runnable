package repository_test

import (
	"testing"

	"github.com/ambardhesi/runnable/pkg/repository"
	"github.com/ambardhesi/runnable/pkg/runnable"
)

func TestStore(t *testing.T) {
	db := repository.NewInMemoryDB()
	jobID := "id"
	job := &runnable.Job{
		ID: jobID,
	}

	if err := db.Store(job); err != nil {
		t.Errorf("did not expect an error")
	}
}

func TestGetItemExists(t *testing.T) {
	db := repository.NewInMemoryDB()
	jobID := "id"
	job := &runnable.Job{
		ID: jobID,
	}

	if err := db.Store(job); err != nil {
		t.Errorf("did not expect an error")
	}

	fetchedJob, ok := db.Get(jobID)
	if !ok {
		t.Errorf("expected item to exist in DB")
	}

	if fetchedJob.ID != jobID {
		t.Errorf("expected job to have id %v", jobID)
	}
}

func TestGetItemDoesNotExist(t *testing.T) {
	db := repository.NewInMemoryDB()
	jobID := "id"

	_, ok := db.Get(jobID)
	if ok {
		t.Errorf("expected item to not exist in DB")
	}
}
