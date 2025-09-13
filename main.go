package main

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus int

const (
	Pending JobStatus = iota
	Running
	Completed
	Failed
	Rejected
)

type JobQueue struct {
	// Mut  sync.Mutex
	// Jobs map[uuid.UUID]Job
}

// func NewJobQueue(maxConcurrentJobs int) *JobQueue {
// 	// return &JobQueue{
// 	// 	Jobs: make([]Job, 0),
// 	// }
// }

type Job struct {
	ID        uuid.UUID
	ClientId  uuid.UUID
	Data      string
	Status    JobStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	Result    interface{}
	Error     error
}

func main() {

}
