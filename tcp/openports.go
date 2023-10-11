package tcp

import (
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/jreisinger/recon"
)

type openPorts struct {
	ports    []int
	scanners int
	timeout  time.Duration // per port
}

type option func(r *openPorts)

func WithPorts(ports []int) option {
	return func(r *openPorts) {
		r.ports = ports
	}
}

func WithScanners(n int) option {
	return func(r *openPorts) {
		r.scanners = n
	}
}

func WithTimeout(timeout time.Duration) option {
	return func(r *openPorts) {
		r.timeout = timeout
	}
}

func OpenPorts(opts ...option) recon.Reconnoiterer {
	op := &openPorts{
		ports:    []int{22, 80, 443},
		scanners: 3,
		timeout:  3 * time.Second,
	}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

func (p openPorts) Recon(target string) recon.Report {
	report := recon.Report{Target: target, Desc: "open tcp ports"}
	openports := openports(target, p.ports, p.scanners, p.timeout)
	for _, port := range openports {
		report.Results = append(report.Results, strconv.Itoa(port))
	}
	return report
}

func openports(target string, ports []int, scanners int, timeout time.Duration) []int {
	in := make(chan int, scanners)
	out := make(chan int)

	for i := 0; i < cap(in); i++ {
		go scanner(target, timeout, in, out)
	}

	go func() {
		for _, port := range ports {
			in <- port
		}
	}()

	var openports []int
	for range ports {
		port := <-out
		if port != 0 {
			openports = append(openports, port)
		}
	}
	sort.Ints(openports)
	return openports
}

func scanner(host string, timeout time.Duration, in, out chan int) {
	for port := range in {
		addr := net.JoinHostPort(host, strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			out <- 0
			continue
		}
		conn.Close()
		out <- port
	}
}
