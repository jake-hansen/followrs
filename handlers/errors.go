package handlers

import "fmt"

// APIError represents an error that occurred during an operation on an endpoint.
type APIError struct {
	Status  int    `json:"status"`  // HTTP status returned to client.
	Err     error  `json:"error"`   // Error that occurred.
	Message string `json:"message"` // Additional information about error returned to client.
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s (HTTP %d): %s", e.Message, e.Status, e.Err.Error())
}

func (e *APIError) Unwrap() error {
	return e.Err
}
