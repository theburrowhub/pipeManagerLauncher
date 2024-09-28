// Package httpServer contains the implementation of the HTTP server that listens for incoming webhook requests.
package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/logging"
)

// Job represents an HTTP request to be processed
// It contains the request URL, method, path, arguments, headers, body, and a result channel
type Job struct {
	RequestURL string              `json:"requestURL"`
	Method     string              `json:"method"`
	Path       string              `json:"path"`
	Args       map[string][]string `json:"args"`
	Headers    map[string][]string `json:"headers"`
	Body       json.RawMessage     `json:"body"`
	ResultChan chan JobResult      `json:"-"`
	RequestID  string              `json:"-"`
}

// JobResult represents the result of processing a job
// It contains the status code and a message
type JobResult struct {
	StatusCode int
	Message    string
}

// jobQueue is a channel for incoming jobs
var jobQueue chan Job

// HttpServer starts the HTTP server to listen for incoming webhook requests
// It processes the requests and sends them to the worker pool
// It also captures termination signals to stop the server and workers
// It returns an error if the server fails to start
func HttpServer(listenAddr string) error {
	// Setup
	maxWorkers := config.Webhook.Data.Workers
	jobQueue = make(chan Job, maxWorkers)

	// Register routes
	routes()

	// Setup termination signals
	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Start the worker pool
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(&wg, i, done)
	}

	// Start the HTTP server
	server := &http.Server{Addr: listenAddr}
	go func() {
		logging.Logger.Info("Starting HTTP server", "listenAddr", listenAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Logger.Error("HTTP server error", "error", fmt.Sprintf("%v", err))
			panic(err)
		}
	}()

	// Wait for a termination signal
	<-sigChan
	logging.Logger.Info("Received shutdown signal, stopping workers and server...")
	close(done)

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	// Wait for all workers to finish
	wg.Wait()

	logging.Logger.Info("Server and workers stopped successfully")
	return nil
}

// routes registers the allowed routes for the HTTP server
func routes() {
	// Health check endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Each allowed route from config file
	for _, route := range config.Webhook.Data.Routes {
		http.HandleFunc(fmt.Sprintf("POST %s", route.Path), webhookHandler)
	}
}

// worker is a function that processes HTTP requests
// It reads jobs from the jobQueue channel and processes them
// It stops when the done channel is closed
func worker(wg *sync.WaitGroup, id int, done <-chan struct{}) {
	defer wg.Done()

	logging.AddAttribute("workerID", fmt.Sprintf("%d", id))
	logging.Logger.Info("Worker started")

	for {
		select {
		case <-done:
			logging.Logger.Info("Worker stopping")
			return
		case job := <-jobQueue:
			requestID := uuid.New().String()
			logging.AddAttribute("requestID", requestID)
			job.RequestID = requestID
			err := processJob(job)
			if err != nil {
				logging.Logger.Error("Error processing job", "error", fmt.Sprintf("%v", err))
				job.ResultChan <- JobResult{
					StatusCode: http.StatusInternalServerError,
					Message:    fmt.Sprintf("Error processing job: %v\n", err),
				}
			} else {
				job.ResultChan <- JobResult{
					StatusCode: http.StatusOK,
					Message:    "Job processed successfully\n",
				}
			}
			logging.RemoveAttribute("requestID")
		}
	}
}
