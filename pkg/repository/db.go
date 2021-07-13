package repository

import (
	"sync"

	"github.com/ambardhesi/runnable/pkg/runnable"
)

// Implementation of runnable.JobStoreService
type InMemoryDB struct {
	jobs map[string]*runnable.Job
	lock sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		jobs: make(map[string]*runnable.Job),
	}
}

func (db *InMemoryDB) Store(job *runnable.Job) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.jobs[job.ID] = job
	return nil
}

func (db *InMemoryDB) Get(jobID string) (*runnable.Job, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	job, exists := db.jobs[jobID]
	return job, exists
}
