package logger

import (
	"github.com/devops-works/openconnect-vyos-telemetry/internal/fetcher"
)

// EventLogger must be implemented to log events
type EventLogger interface {
	Log(message string, labels map[string]interface{}) error
}

// MetricsLogger must be implemented to log metrics
type MetricsLogger interface {
	Add(fetcher.Metrics) error
	Flush() error
}
