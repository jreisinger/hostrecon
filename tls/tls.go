package tls

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/jreisinger/hostrecon"
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

// ---

type CA struct {
	Port    string // e.g. 443
	Timeout time.Duration
}

func (ca CA) Recon(host string) hostrecon.Info {
	recon := hostrecon.Info{Host: host, Kind: "tls ca"}
	addr := net.JoinHostPort(host, ca.Port)
	conn, err := getConn(addr, ca.Timeout, true)
	if err != nil {
		recon.Err = err
		return recon
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	certsButLast := certs[len(certs)-1]
	recon.Info = append(recon.Info, certsButLast.Issuer.Organization...)
	return recon
}

// ---

type Issuer struct {
	Port    string // e.g. 443
	Timeout time.Duration
}

func (i Issuer) Recon(host string) hostrecon.Info {
	info := hostrecon.Info{Host: host, Kind: "tls cert issuer"}
	addr := net.JoinHostPort(host, i.Port)
	conn, err := getConn(addr, i.Timeout, true)
	if err != nil {
		info.Err = err
		return info
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	leaf := certs[0]
	info.Info = append(info.Info, leaf.Issuer.Organization...)
	return info
}

// ---

type Version struct {
	Port    string // e.g. 443
	Timeout time.Duration
}

func (v Version) Recon(target string) hostrecon.Info {
	info := hostrecon.Info{Host: target, Kind: "tls version"}
	addr := net.JoinHostPort(target, v.Port)
	conn, err := getConn(addr, v.Timeout, true)
	if err != nil {
		info.Err = err
		return info
	}
	defer conn.Close()
	ver := tls.VersionName(conn.ConnectionState().Version)
	info.Info = append(info.Info, ver)
	return info
}
