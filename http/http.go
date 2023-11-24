package http

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jreisinger/recon"
)

type conn struct {
	port    string
	timeout time.Duration
}

type option func(t *conn)

func WithPort(port string) option {
	return func(c *conn) {
		c.port = port
	}
}

func WithTimeout(timeout time.Duration) option {
	return func(c *conn) {
		c.timeout = timeout
	}
}

type version conn

func Version(opts ...option) recon.Reconnoiterer {
	c := &conn{
		port:    "443",
		timeout: 3 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	return version(*c)
}

func (t version) Recon(target string) recon.Report {
	report := recon.Report{Target: target, Info: "http version"}
	addr := net.JoinHostPort(target, t.port)
	ver, err := getVersion(addr, t.timeout)
	if err != nil {
		report.Err = err
		return report
	}
	report.Results = append(report.Results, ver)
	return report
}

func getVersion(host string, timeout time.Duration) (string, error) {
	url := fmt.Sprintf("https://%s", host)
	c := http.Client{Timeout: timeout}
	resp, err := c.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return resp.Proto, nil
}
