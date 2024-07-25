package dns

import (
	"hostrecon"
	"net"
	"strings"
)

type Cname struct{}

func (Cname) Recon(host string) hostrecon.Info {
	info := hostrecon.Info{Host: host, Kind: "cname"}
	cname, err := net.LookupCNAME(host)
	if err != nil {
		info.Err = err
		return info
	}
	cname, _ = strings.CutSuffix(cname, ".")
	if cname != host {
		info.Info = append(info.Info, cname)
	}
	return info
}

// ---

type IpAddr struct{}

func (IpAddr) Recon(host string) hostrecon.Info {
	info := hostrecon.Info{Host: host, Kind: "ip addresses"}
	addrs, err := net.LookupHost(host)
	if err != nil {
		info.Err = err
		return info
	}
	info.Info = append(info.Info, addrs...)
	return info
}

// ---

type Mx struct{}

func (Mx) Recon(host string) hostrecon.Info {
	info := hostrecon.Info{Host: host, Kind: "mail servers"}
	mxs, err := net.LookupMX(host)
	if err != nil {
		info.Err = err
		return info
	}
	for _, mx := range mxs {
		s, _ := strings.CutSuffix(mx.Host, ".")
		if s == "" {
			continue
		}
		info.Info = append(info.Info, s)
	}
	return info
}

// ---

type Ns struct{}

func (Ns) Recon(host string) hostrecon.Info {
	info := hostrecon.Info{Host: host, Kind: "name servers"}
	nss, err := net.LookupNS(host)
	if err != nil {
		info.Err = err
		return info
	}
	for _, ns := range nss {
		n, _ := strings.CutSuffix(ns.Host, ".")
		info.Info = append(info.Info, n)
	}
	return info
}

// ---

type Txt struct{}

func (t Txt) Recon(host string) hostrecon.Info {
	info := hostrecon.Info{Host: host, Kind: "txt records"}
	records, err := net.LookupTXT(host)
	if err != nil {
		info.Err = err
		return info
	}
	info.Info = records
	return info
}
