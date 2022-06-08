package openconnectmetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/devops-works/openconnect-vyos-telemetry/internal/fetcher"
)

// Fetcher is a metrics fetcher implementation for OpenConnect
type Fetcher struct{}

// Fetch metrics from occtl output
func (f *Fetcher) Fetch(ctx context.Context, b []byte) (fetcher.Metrics, error) {
	e := Entries{}
	metrics := fetcher.Metrics{}
	err := json.Unmarshal(b, &e)
	if err != nil {
		return metrics, err
	}

	for _, v := range e {
		// user is currently connecting to the VPN
		if v.Username == "(none)" {
			continue
		}

		m := fetcher.Metric{
			Measurement: "openconnect",
			Labels: map[string]interface{}{
				"username":  v.Username,
				"state":     v.State,
				"remote_ip": v.RemoteIP,
			},
			Values:    map[string]float64{},
			Timestamp: time.Now().UTC(),
		}

		f, err := strconv.ParseFloat(v.Rx, 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing Rx: %w", err)
		}
		m.Values["received_bytes"] = f

		f, err = strconv.ParseFloat(v.Tx, 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing Tx: %w", err)
		}
		m.Values["transmitted_bytes"] = f

		rx := strings.Split(v.AverageRX, "")
		f, err = strconv.ParseFloat(rx[0], 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing AverageRX: %w", err)
		}
		m.Values["avg_received_bytespersec"] = f

		tx := strings.Split(v.AverageTX, "")
		f, err = strconv.ParseFloat(tx[0], 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing AverageTX: %w", err)
		}
		m.Values["avg_transmitted_bytespersec"] = f

		t, err := time.Parse("2006-01-02 15:04", v.ConnectedAt)
		if err != nil {
			return nil, err
		}
		duration := time.Now().UTC().Sub(t)
		m.Values["connection_age"] = duration.Seconds()

		metrics = append(metrics, m)
	}

	return metrics, err
}
