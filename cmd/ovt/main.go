package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/devops-works/openconnect-vyos-telemetry/internal/fetcher"
	"github.com/devops-works/openconnect-vyos-telemetry/internal/logger"
	"github.com/oklog/run"

	opcfetch "github.com/devops-works/openconnect-vyos-telemetry/internal/fetcher/openconnect-metrics"
	influxlogger "github.com/devops-works/openconnect-vyos-telemetry/internal/logger/influxdb"
	flags "github.com/jessevdk/go-flags"
)

var (
	// Version of current binary
	Version string
	// BuildDate of current binary
	BuildDate string
)

func main() {
	var (
		g    run.Group
		opts struct {
			LokiURL        string `short:"L" long:"loki.url" description:"Loki server URL" required:"true"`
			LokiUsername   string `short:"U" long:"loki.user" description:"Loki basic auth username" required:"false"`
			LokiPassword   string `short:"P" long:"loki.pass" description:"Loki basic auth password" required:"false"`
			LokiOrgID      string `short:"O" long:"loki.orgid" description:"Loki X-Scope-OrgID header to add" required:"false"`
			InfluxURL      string `short:"I" long:"influx.url" description:"InfluxDB server URL" required:"true"`
			InfluxUsername string `short:"V" long:"influx.user" description:"InfluxDB basic auth username" required:"false"`
			InfluxPassword string `short:"Q" long:"influx.pass" description:"InfluxDB basic auth password" required:"false"`
			InfluxDatabase string `short:"D" long:"influx.db" description:"InfluxDB database name" required:"true"`

			MetricsDelay int    `short:"d" long:"metrics.delay" description:"Delay between metrics collection in seconds" required:"false" default:"2"`
			MetricsCmd   string `short:"M" long:"metrics.cmd" description:"Command to run to fetch metrics" required:"false" default:"sudo occtl -s /run/ocserv/occtl.socket -j show users"`

			EventsCmd string `short:"E" long:"events.cmd" description:"Command to run to fetch events" required:"false" default:"sudo occtl -s /run/ocserv/occtl.socket show events"`

			Debug  bool `long:"debug" description:"Enable debug logging" required:"false"`
			DryRun bool `short:"n" long:"dry-run" description:"Dry run, do not send any data to InfluxDB" required:"false"`

			Version func() `short:"v" long:"version" description:"displays versions"`
		}
	)

	opts.Version = func() {
		fmt.Fprintf(os.Stderr, "ovs (openconnect-vyos-telemetry) version %s (built %s)\n", Version, BuildDate)
		os.Exit(0)
	}

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	influxfetcher := &opcfetch.Fetcher{}
	influxlogger := influxlogger.New(opts.InfluxURL, opts.InfluxDatabase, opts.InfluxUsername, opts.InfluxPassword, opts.DryRun)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan struct{})
	g.Add(func() error {
		return signalHandler(stop)
	}, func(error) { close(stop) })

	g.Add(func() error {
		return collectMetrics(ctx, influxfetcher, influxlogger, opts.MetricsCmd, opts.MetricsDelay)
	}, func(error) { cancel() })

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}

func signalHandler(stop chan struct{}) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	select {
	case <-signals: // received signal
		return fmt.Errorf("received interrupt")
	case <-stop: // stopped by rungroup
		return nil
	}
}

func collectMetrics(ctx context.Context, f fetcher.MetricsFetcher, l logger.MetricsLogger, cmd string, delay int) error {
	for {
		splitted := strings.Split(cmd, " ")
		out, err := exec.Command(splitted[0], splitted[1:]...).Output()
		if err != nil {
			log.Fatal(err)
		}

		m, err := f.Fetch(ctx, out)
		if err != nil {
			return err
		}

		if err := l.Add(m); err != nil {
			return err
		}

		// flush to influxdb
		err = l.Flush()
		if err != nil {
			log.Printf("error flushing metrics: %s", err)
		}

		// Sleep with cancel
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		case <-time.After(time.Duration(delay) * time.Second): // sleep before next fetch
		}
	}
}
