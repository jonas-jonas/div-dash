package model

import "time"

type APIError struct {
	Message   string    `json:"message"`
	Status    int       `json:"status"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

func NewApiError(message string, status int, path string) APIError {
	return APIError{
		Message:   message,
		Status:    status,
		Path:      path,
		Timestamp: time.Now(),
	}
}
