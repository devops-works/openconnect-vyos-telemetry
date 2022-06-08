package fetcher

import (
	"context"
	"io"
	"time"
)

// Event contains a log event
type Event struct {
	Message string
	Labels  map[string]interface{}
	Error   error
}

// Metric contains a single metric
type Metric struct {
	Measurement string
	Labels      map[string]interface{}
	Values      map[string]float64
	Timestamp   time.Time
	Error       error
}

// Metrics contains a list of metrics
type Metrics []Metric

// EventsFetcher is an interface for fetching events
type EventsFetcher interface {
	Fetch(stream io.Reader) (chan Event, error)
}

// MetricsFetcher is an interface for fetching metrics
type MetricsFetcher interface {
	Fetch(context.Context, []byte) (Metrics, error)
}
