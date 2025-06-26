package cron

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
)

// CronJobHandlerInterface defines the contract for cron job handlers.
type CronJobHandlerInterface interface {
	TriggerCleanupJob(c *gin.Context)
	TriggerDataSyncJob(c *gin.Context)
}

// CronJobHandler struct
type CronJobHandler struct {
	// Simulate job state tracking
	mu   sync.Mutex
	jobs map[int]string // jobID -> state ("queued", "started", "completed")
}

// NewCronJobHandler creates a new cron job handler.
func NewCronJobHandler() CronJobHandlerInterface {
	return &CronJobHandler{
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
func (h *CronJobHandler) TriggerCleanupJob(c *gin.Context) {
	// Start a span for the API request
	ctx, span := tracer.StartSpan(c.Request.Context(), "TriggerDataSyncJob")
	defer span.End()
	logger.WithContext(ctx).Infof("Cleanup job triggered via API")
	// Cleanup jobs that are not yet started or still queued
	h.mu.Lock()
	cleaned := []int{}
	for id, state := range h.jobs {
		if state == "queued" {
			delete(h.jobs, id)
			cleaned = append(cleaned, id)
		}
	}
	h.mu.Unlock()
	logger.WithContext(ctx).Infof("Cleaned up jobs: %v", cleaned)
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Cleanup job has been triggered. Cleaned jobs: %v", cleaned),
	})
}

// TriggerDataSyncJob handles POST /cron/sync
// @Summary Trigger a data sync job
// @Description Triggers a background data synchronization job.
// @Tags Cron
// @Produce json
// @Success 202 {object} map[string]interface{}
// @Router /cron/sync [post]
func (h *CronJobHandler) TriggerDataSyncJob(c *gin.Context) {
	ctx, span := tracer.StartSpan(c.Request.Context(), "TriggerDataSyncJob")
	defer span.End()
	logger.WithContext(ctx).Infof("Data sync job triggered via API")

	// Trigger jobs with delays 1,2,3,4,5 seconds
	h.jobs[99] = "queued" // Initialize a dummy job to show that jobs are being queued
	for i := 1; i <= 5; i++ {
		h.mu.Lock()
		sec := 6 + rand.Intn(5) // random seconds between 6 and 10
		logger.WithContext(ctx).Infof("Queuing data sync job %d with delay %d seconds", i, sec)
		h.jobs[i] = "queued"
		h.mu.Unlock()
		// Pass context with span to goroutine for child span
		go h.runCronJob(ctx, i, sec)
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "Data sync jobs (1-5) have been triggered.",
	})
}

// runCronJob simulates running a cron job with a delay and logs completion.
// Note: Mutex is locked twice to avoid race conditions as multiple goroutines
// can try to access the same shared object (h.jobs) concurrently.
func (h *CronJobHandler) runCronJob(parentCtx context.Context, jobID int, duration int) {
	// Start a child span for the background job
	ctx, span := tracer.StartSpan(parentCtx, "runCronJob") // You can add attributes here if needed
	defer span.End()

	h.mu.Lock()
	h.jobs[jobID] = "started"
	h.mu.Unlock()

	span.AddEvent("Job started")
	time.Sleep(time.Duration(duration) * time.Second)
	logger.WithContext(ctx).Infof("job %d executed successfully in %d seconds", jobID, duration)
	span.AddEvent("Job completed")

	h.mu.Lock()
	h.jobs[jobID] = "completed"
	h.mu.Unlock()
}
