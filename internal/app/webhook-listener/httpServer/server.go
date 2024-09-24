// Package httpServer contains the implementation of the HTTP server that listens for incoming webhook requests.
package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/databuilder"
	"github.com/sergiotejon/pipeManager/internal/app/webhook-listener/pipeline"
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
func HttpServer(listenPort int) error {
	// Setup
	maxWorkers := config.Webhook.Data.Workers
	jobQueue = make(chan Job, maxWorkers)

	// Each allowed route from config file
	for _, route := range config.Webhook.Data.Routes {
		http.HandleFunc(fmt.Sprintf("POST %s", route.Path), webhookHandler)
	}

	// TODO:
	// Add route to health check endpoint
	// Move webhookHandler to a separate file

	// Capture termination signals
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
	server := &http.Server{Addr: fmt.Sprintf(":%d", listenPort)}
	go func() {
		logging.Logger.Info("Starting HTTP server", "port", listenPort)
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

// webhookHandler is the function that handles incoming webhook requests
// It reads the request body and headers, creates a job, and sends it to the worker pool
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Read the body as a json.RawMessage
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Logger.Error("Error reading request body", "error", fmt.Sprintf("%v", err))
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var jsonCheck interface{}
	err = json.Unmarshal(bodyBytes, &jsonCheck)
	if err != nil {
		logging.Logger.Info("Error validating body as JSON", "error", fmt.Sprintf("%v", err))
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Read headers as a map
	headers := make(map[string][]string)
	for name, values := range r.Header {
		headers[name] = values
	}

	// TODO: Add additional validation here if needed
	// Optional: Verify the HMAC of the webhook for added security.
	// token := "my_webhook_secret"
	// sig := r.Header.Get("X-Hub-Signature-256")
	// if !verifySignature(body, sig, token) {
	//     http.Error(w, "Invalid signature", http.StatusUnauthorized)
	//     return
	// }

	// Create a job
	resultChan := make(chan JobResult)
	job := Job{
		RequestURL: r.URL.String(),
		Method:     r.Method,
		Path:       r.URL.Path,
		Args:       r.URL.Query(),
		Headers:    headers,
		Body:       json.RawMessage(bodyBytes),
		ResultChan: resultChan,
	}
	jobQueue <- job // Send the job to the worker queue

	// Send a response
	result := <-resultChan
	w.WriteHeader(result.StatusCode)
	_, err = w.Write([]byte(result.Message))
	if err != nil {
		logging.Logger.Error("Error writing response", "error", fmt.Sprintf("%v", err))
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

// processJob is the function that processes the incoming HTTP request
// It creates a PipelineData object and launches a job
// It returns an error if the job fails to launch
func processJob(job Job) error {
	var err error

	logging.Logger.Info("Request received", "method", job.Method, "path", job.Path)
	logging.Logger.Debug("Payload", "job", job)

	var jsonData []byte
	jsonData, err = json.MarshalIndent(job, "", "  ")
	if err != nil {
		return err
	}

	var pipelineData *databuilder.PipelineData
	pipelineData, err = databuilder.Run(jsonData, job.Path, config.Webhook.Data.Routes)
	if err != nil {
		return err
	}

	logging.Logger.Debug("Pipeline", "data", pipelineData)

	err = pipeline.LaunchJob(job.RequestID, pipelineData)
	if err != nil {
		return err
	}

	return nil
}
