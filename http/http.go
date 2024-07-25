package http

import (
	"fmt"
	"hostrecon"
	"net"
	"net/http"
	"time"
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

type Version struct {
	Port    string // e.g. "443"
	Timeout time.Duration
}

func (v Version) Recon(target string) hostrecon.Info {
	report := hostrecon.Info{Host: target, Kind: "http version"}
	addr := net.JoinHostPort(target, v.Port)
	ver, err := getVersion(addr, v.Timeout)
	if err != nil {
		report.Err = err
		return report
	}
	report.Info = append(report.Info, ver)
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
