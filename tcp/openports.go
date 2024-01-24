package tcp

import (
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/jreisinger/recon"
)

type openPorts struct {
	portsToScan []int
	timeout     time.Duration // per port
}

type option func(r *openPorts)

func WithPortsToScan(ports []int) option {
	return func(r *openPorts) {
		r.portsToScan = ports
	}
}

func WithTimeout(timeout time.Duration) option {
	return func(r *openPorts) {
		r.timeout = timeout
	}
}

func OpenPorts(opts ...option) recon.Reconnoiterer {
	op := &openPorts{
		portsToScan: []int{22, 80, 443},
		timeout:     3 * time.Second,
	}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

func (p openPorts) Recon(target string) recon.Report {
	report := recon.Report{Host: target, Area: "open tcp ports"}
	openports := openports(target, p.portsToScan, p.timeout)
	for _, port := range openports {
		report.Info = append(report.Info, strconv.Itoa(port))
	}
	return report
}

func openports(target string, ports []int, timeout time.Duration) []int {
	results := make(chan int)
	for _, port := range ports {
		go func(host string, port int) {
			addr := net.JoinHostPort(host, strconv.Itoa(port))
			conn, err := net.DialTimeout("tcp", addr, timeout)
			if err != nil {
				results <- 0
				return
			}
			conn.Close()
			results <- port
		}(target, port)
	}
	var openports []int
	for range ports {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	sort.Ints(openports)
	return openports
}
