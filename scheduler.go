package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type JobStatus int

const (
	Pending JobStatus = iota
	Running
	Completed
	Failed
)

type JobQueue struct {
	maxConcurrency   int
	currentScheduled int
	jobs             map[uuid.UUID]*Job
	jobsChan         chan uuid.UUID
	mu               sync.RWMutex
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
}

type Job struct {
	ID        uuid.UUID   `json:"id"`
	ClientId  uuid.UUID   `json:"client_id"`
	Data      string      `json:"data"`
	Status    JobStatus   `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Result    interface{} `json:"result,omitempty"`
	Error     string      `json:"error,omitempty"`
}

type ScheduleRequest struct {
	ClientId string `json:"client_id"`
	Data     string `json:"data"`
}

type ScheduleResponse struct {
	JobID   uuid.UUID `json:"job_id"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
}

type CapacityResponse struct {
	Current   int `json:"current"`
	Max       int `json:"max"`
	Available int `json:"available"`
}

type Server struct {
	jobQueue *JobQueue
	router   *mux.Router
}

func NewJobQueue(maxConcurrentJobs int) *JobQueue {
	ctx, cancel := context.WithCancel(context.Background())

	jq := &JobQueue{
		maxConcurrency: maxConcurrentJobs,
		jobs:           make(map[uuid.UUID]*Job),
		jobsChan:       make(chan uuid.UUID, 1000),
		ctx:            ctx,
		cancel:         cancel,
	}

	for i := 0; i < maxConcurrentJobs; i++ {
		jq.wg.Add(1)
		go jq.worker()
	}

	return jq
}

func (jq *JobQueue) TrySchedule(clientId uuid.UUID, data string) (*Job, error) {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	if jq.currentScheduled >= jq.maxConcurrency {
		return nil, errors.New("job queue at capacity")
	}

	job := &Job{
		ID:        uuid.New(),
		ClientId:  clientId,
		Data:      data,
		Status:    Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jq.jobs[job.ID] = job
	jq.currentScheduled++
	jq.jobsChan <- job.ID

	return job, nil
}

func (jq *JobQueue) GetCapacity() (current, max int) {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	return jq.currentScheduled, jq.maxConcurrency
}

func (jq *JobQueue) GetJob(jobID uuid.UUID) (*Job, bool) {
	jq.mu.RLock()
	defer jq.mu.RUnlock()
	job, exists := jq.jobs[jobID]
	return job, exists
}

func (jq *JobQueue) worker() {
	defer jq.wg.Done()

	for {
		select {
		case jobID := <-jq.jobsChan:
			jq.processJob(jobID)
		case <-jq.ctx.Done():
			return
		}
	}
}

func (jq *JobQueue) processJob(jobID uuid.UUID) {
	jq.mu.Lock()
	job, exists := jq.jobs[jobID]
	if !exists {
		jq.mu.Unlock()
		return
	}
	job.Status = Running
	job.UpdatedAt = time.Now()
	jq.mu.Unlock()

	result, err := jq.executeJob(job)

	jq.mu.Lock()
	if err != nil {
		job.Status = Failed
		job.Error = err.Error()
	} else {
		job.Status = Completed
		job.Result = result
	}
	job.UpdatedAt = time.Now()
	jq.currentScheduled--
	jq.mu.Unlock()
}

func (jq *JobQueue) executeJob(job *Job) (interface{}, error) {
	time.Sleep(10 * time.Second)

	return map[string]interface{}{
		"processed_data": job.Data,
		"timestamp":      time.Now(),
	}, nil
}

func (jq *JobQueue) Shutdown() {
	jq.cancel()
	jq.wg.Wait()
}

func NewServer(maxJobs int) *Server {
	jobQueue := NewJobQueue(maxJobs)
	router := mux.NewRouter()

	server := &Server{
		jobQueue: jobQueue,
		router:   router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/schedule", s.scheduleJob).Methods("POST")
	s.router.HandleFunc("/jobs/{id}", s.getJob).Methods("GET")
	s.router.HandleFunc("/capacity", s.getCapacity).Methods("GET")
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")
}

func (s *Server) scheduleJob(w http.ResponseWriter, r *http.Request) {
	var req ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	clientID, err := uuid.Parse(req.ClientId)
	if err != nil {
		http.Error(w, "Invalid client_id format", http.StatusBadRequest)
		return
	}

	job, err := s.jobQueue.TrySchedule(clientID, req.Data)
	if err != nil {
		response := ScheduleResponse{
			Status:  "rejected",
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ScheduleResponse{
		JobID:   job.ID,
		Status:  "accepted",
		Message: "Job scheduled successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid job ID format", http.StatusBadRequest)
		return
	}

	job, exists := s.jobQueue.GetJob(jobID)
	if !exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (s *Server) getCapacity(w http.ResponseWriter, r *http.Request) {
	current, max := s.jobQueue.GetCapacity()
	response := CapacityResponse{
		Current:   current,
		Max:       max,
		Available: max - current,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	current, max := s.jobQueue.GetCapacity()
	status := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"jobs":      s.jobQueue.jobs,
		"capacity": map[string]int{
			"current": current,
			"max":     max,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (s *Server) Start(port string) error {
	fmt.Println("Starting server on port", port)
	return http.ListenAndServe(":"+port, s.router)
}

func (s *Server) Shutdown() {
	s.jobQueue.Shutdown()
}

func main() {
	server := NewServer(3)
	defer server.Shutdown()

	server.Start("3000")
}
