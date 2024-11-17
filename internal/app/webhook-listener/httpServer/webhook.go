package httpServer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sergiotejon/pipeManagerLauncher/internal/pkg/logging"
)

// webhookHandler is the function that handles incoming webhook requests
// It reads the request body and headers, creates a job, and sends it to the worker pool
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Defer closing the request body to prevent resource leaks
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Logger.Error("Error closing request body", "error", fmt.Sprintf("%v", err))
		}
	}(r.Body)

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
