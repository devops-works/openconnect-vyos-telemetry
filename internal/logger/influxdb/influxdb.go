package influxdb

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/devops-works/openconnect-vyos-telemetry/internal/fetcher"
)

type point struct {
	measurement string
	tags        string
	values      map[string]float64
	timestamp   time.Time
}

func (p point) String() string {
	influxString := p.measurement
	if p.tags != "" {
		influxString = fmt.Sprintf("%s,%s", influxString, p.tags)
	}

	values := []string{}
	for k, v := range p.values {
		values = append(values, fmt.Sprintf("%s=%g", k, v))
	}
	influxString = fmt.Sprintf("%s %s %d", influxString, strings.Join(values, ","), p.timestamp.UnixNano())

	return influxString
}

// Logger implements the Logger interface for influxdb
type Logger struct {
	url       string
	username  string
	password  string
	hostname  string
	database  string
	timeout   int
	retries   int
	points    []point
	maxpoints int
	dryrun    bool
	mu        sync.Mutex
}

// New returns a new influxdb logger
func New(url, db, user, pass string, dryrun bool) *Logger {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	return &Logger{
		url:       url,
		database:  db,
		username:  user,
		password:  pass,
		hostname:  host,
		timeout:   2000,
		retries:   2,
		maxpoints: 1000,
		dryrun:    dryrun,
		points:    []point{},
	}
}

// Add adds a new point to the bag
func (i *Logger) Add(metrics fetcher.Metrics) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, m := range metrics {
		if len(i.points) >= i.maxpoints {
			return fmt.Errorf("influx max cached points reached")
		}

		tags := ""

		for k, v := range m.Labels {
			tags += fmt.Sprintf("%s=%s,", k, v)
		}

		tags += "host=" + i.hostname

		i.points = append(i.points, point{
			measurement: m.Measurement,
			tags:        tags,
			values:      m.Values,
			timestamp:   m.Timestamp,
		})
	}
	return nil
}

// Flush the points to the influxdb server
func (i *Logger) Flush() error {
	return i.batchLogInfluxDB()
}

func (i *Logger) logInfluxDB(lines []byte) error {
	params := url.Values{}
	url := fmt.Sprintf("%s/write?", i.url)

	if i.url[len(i.url)-1] == '/' {
		url = fmt.Sprintf("%swrite?", i.url)
	}

	params.Set("db", i.database)
	params.Set("u", i.username)
	params.Set("p", i.password)

	url += params.Encode()

	client := &http.Client{
		Timeout: time.Duration(i.timeout) * time.Millisecond,
	}

	var (
		err  error
		resp *http.Response
	)

	r := bytes.NewReader(lines)

	for try := 0; try < i.retries; try++ {
		resp, err = client.Post(url, "application/x-www-form-urlencoded", r)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("unable to write to influxdb server %s, got response: %s", i.url, resp.Status)
	}
	return nil
}

func (i *Logger) dumpInfluxDB(lines []byte, w io.Writer) error {
	var uri string
	if i.url[len(i.url)-1] == '/' {
		uri = fmt.Sprintf("%swrite?db=%s", i.url, i.database)
	} else {
		uri = fmt.Sprintf("%s/write?db=%s", i.url, i.database)
	}

	if i.username != "" {
		uri += fmt.Sprintf("&u=%s&p=%s", i.username, i.password)
	}

	fmt.Fprintln(w, uri)
	fmt.Fprintln(w, string(lines))
	return nil
}

func (i *Logger) batchLogInfluxDB() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if len(i.points) == 0 {
		return nil
	}

	lines := []byte{}
	for c, pt := range i.points {
		pt.measurement = "openconnect"

		lines = append(lines, []byte(pt.String())...)
		lines = append(lines, '\n')
		if (c+1)%500 == 0 {
			err := i.logInfluxDB(lines)
			if err != nil {
				return err
			}
			lines = []byte{}
		}
	}

	if i.dryrun {
		err := i.dumpInfluxDB(lines, os.Stdout)
		i.points = []point{}

		return err
	}

	err := i.logInfluxDB(lines)
	if err != nil {
		return err
	}

	i.points = []point{}

	return nil
}
