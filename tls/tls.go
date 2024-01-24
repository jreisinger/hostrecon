package tls

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/jreisinger/recon"
)

func getConn(addr string, timeout time.Duration, insecureSkipVerify bool) (*tls.Conn, error) {
	dialer := net.Dialer{
		Timeout: timeout,
	}
	conn, err := tls.DialWithDialer(
		&dialer, "tcp",
		addr,
		&tls.Config{InsecureSkipVerify: insecureSkipVerify})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

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

// ---

type ca conn

func CA(opts ...option) recon.Reconnoiterer {
	c := &conn{
		port:    "443",
		timeout: 3 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	return ca(*c)
}

func (c ca) Recon(target string) recon.Report {
	recon := recon.Report{Host: target, Area: "certificate authority"}
	addr := net.JoinHostPort(target, c.port)
	conn, err := getConn(addr, c.timeout, true)
	if err != nil {
		recon.Err = err
		return recon
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	ca := certs[len(certs)-1]
	recon.Data = append(recon.Data, ca.Issuer.Organization...)
	return recon
}

// ---

type issuer conn

func Issuer(opts ...option) recon.Reconnoiterer {
	c := &conn{
		port:    "443",
		timeout: 3 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	return issuer(*c)
}

func (t issuer) Recon(target string) recon.Report {
	report := recon.Report{Host: target, Area: "certificate issuer"}
	addr := net.JoinHostPort(target, t.port)
	conn, err := getConn(addr, t.timeout, true)
	if err != nil {
		report.Err = err
		return report
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	leaf := certs[0]
	report.Data = append(report.Data, leaf.Issuer.Organization...)
	return report
}

// ---

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
	report := recon.Report{Host: target, Area: "tls version"}
	addr := net.JoinHostPort(target, t.port)
	conn, err := getConn(addr, t.timeout, true)
	if err != nil {
		report.Err = err
		return report
	}
	defer conn.Close()
	ver := tls.VersionName(conn.ConnectionState().Version)
	report.Data = append(report.Data, ver)
	return report
}
