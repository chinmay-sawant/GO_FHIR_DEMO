// Package cron provides Gin handlers for cron jobs.
package cron

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"sync"
	"time"

	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
)

const (
	stateQueued    = "queued"
	stateStarted   = "started"
	stateCompleted = "completed"
)

// JobHandler defines the contract for cron job handlers.
type JobHandler interface {
	TriggerCleanupJob(c *gin.Context)
	TriggerDataSyncJob(c *gin.Context)
}


// JobHandlerImpl struct
type JobHandlerImpl struct {
	// Simulate job state tracking
	mu   sync.Mutex
	jobs map[int]string // jobID -> state ("queued", "started", "completed")
}

// NewJobHandler creates a new cron job handler.
func NewJobHandler() JobHandler {
	return &JobHandlerImpl{
		jobs: make(map[int]string),
	}
}

// TriggerCleanupJob handles POST /cron/cleanup
// @Summary Trigger a cleanup job
// @Description Triggers a background cleanup job.
// @Tags Cron
// @Produce json
// @Success 202 {object} map[string]interface{}
// @Router /cron/cleanup [post]
func (h *JobHandlerImpl) TriggerCleanupJob(c *gin.Context) {
	// Start a span for the API request
	ctx, span := tracer.StartSpan(c.Request.Context(), "TriggerDataSyncJob")
	defer span.End()
	logger.WithContext(ctx).Infof("Cleanup job triggered via API")
	// Cleanup jobs that are not yet started or still queued
	cleaned := make([]int, 0, len(h.jobs))
	h.mu.Lock()
	for id, state := range h.jobs {
		if state == stateQueued {
			delete(h.jobs, id)
			cleaned = append(cleaned, id)
		}
	}
	h.mu.Unlock()
	logger.WithContext(ctx).Infof("Cleaned up jobs: %v", cleaned)
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
	})
}

// TriggerDataSyncJob handles POST /cron/sync
// @Summary Trigger a data sync job
// @Description Triggers a background data synchronization job.
// @Tags Cron
// @Produce json
// @Success 202 {object} map[string]interface{}
// @Router /cron/sync [post]
func (h *JobHandlerImpl) TriggerDataSyncJob(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "TriggerDataSyncJob")
	defer span.End()
	logger.WithContext(ctx).Infof("Data sync job triggered via API")

	// Trigger jobs with delays 1,2,3,4,5 seconds
	h.mu.Lock()
	h.jobs[99] = stateQueued // Initialize a dummy job

	// Use a WaitGroup to coordinate the spawned goroutines
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 1; i <= 5; i++ {
		rndVal, _ := rand.Int(rand.Reader, big.NewInt(5))
		sec := 6 + int(rndVal.Int64()) // random seconds between 6 and 10
		logger.WithContext(ctx).Infof("Queuing data sync job %d with delay %d seconds", i, sec)
		h.jobs[i] = stateQueued

		// Pass context with span to goroutine for child span
		go h.runCronJobAsync(ctx, &wg, i, sec)
	}
	h.mu.Unlock()

	// We don't wait for wg here to return 202 quickly, but we could
	// launch another goroutine to wait and log if needed.
	go func() {
		wg.Wait()
		logger.WithContext(ctx).Info("All triggered data sync jobs completed")
	}()
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "Data sync jobs (1-5) have been triggered.",
	})
}

func (h *JobHandlerImpl) runCronJobAsync(parentCtx context.Context, wg *sync.WaitGroup, jobID int, duration int) {
	defer wg.Done()
	h.runCronJob(parentCtx, jobID, duration)
}

// runCronJob simulates running a cron job with a delay and logs completion.
// Note: Mutex is locked twice to avoid race conditions as multiple goroutines
// can try to access the same shared object (h.jobs) concurrently.
func (h *JobHandlerImpl) runCronJob(parentCtx context.Context, jobID int, duration int) {
	// Start a child span for the background job
	ctx, span := tracer.StartSpan(parentCtx, "runCronJob") // You can add attributes here if needed
	defer span.End()

	h.mu.Lock()
	h.jobs[jobID] = stateStarted
	h.mu.Unlock()

	span.AddEvent("Job started")
	time.Sleep(time.Duration(duration) * time.Second)
	logger.WithContext(ctx).Infof("job %d executed successfully in %d seconds", jobID, duration)
	span.AddEvent("Job completed")

	h.mu.Lock()
	h.jobs[jobID] = stateCompleted
	h.mu.Unlock()
}
