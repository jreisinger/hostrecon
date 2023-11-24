package dns

import (
	"errors"
	"net"
	"strings"

	"github.com/jreisinger/recon"
)

func errNotFound(err error) bool {
	var dnsErr *net.DNSError
	return errors.As(err, &dnsErr) && dnsErr.IsNotFound
}

// ---

type cname struct{}

func Cname() recon.Reconnoiterer { return cname{} }

func (cname) Recon(target string) recon.Report {
	report := recon.Report{Target: target, Info: "cname"}
	cname, err := net.LookupCNAME(target)
	if err != nil {
		report.Err = err
		return report
	}
	cname, _ = strings.CutSuffix(cname, ".")
	if cname != target {
		report.Results = append(report.Results, cname)
	}
	return report
}

// ---

type ipaddr struct{}

func IPAddr() recon.Reconnoiterer { return ipaddr{} }

func (ipaddr) Recon(target string) recon.Report {
	report := recon.Report{Target: target, Info: "ip addresses"}
	addrs, err := net.LookupHost(target)
	if err != nil {
		report.Err = err
		return report
	}
	report.Results = append(report.Results, addrs...)
	return report
}

// ---

type mx struct{}

func MX() recon.Reconnoiterer { return mx{} }

func (mx) Recon(target string) recon.Report {
	report := recon.Report{Target: target, Info: "mail servers"}
	mxs, err := net.LookupMX(target)
	if err != nil {
		report.Err = err
		return report
	}
	for _, mx := range mxs {
		s, _ := strings.CutSuffix(mx.Host, ".")
		if s == "" {
			continue
		}
		report.Results = append(report.Results, s)
	}
	return report
}

// ---

type ns struct{}

func NS() recon.Reconnoiterer { return ns{} }

func (ns) Recon(target string) recon.Report {
	report := recon.Report{Target: target, Info: "name servers"}
	nss, err := net.LookupNS(target)
	if err != nil {
		report.Err = err
		return report
	}
	for _, ns := range nss {
		n, _ := strings.CutSuffix(ns.Host, ".")
		report.Results = append(report.Results, n)
	}
	return report
}

// ---

type txt struct{}

func TXT() recon.Reconnoiterer { return txt{} }

func (txt) Recon(target string) recon.Report {
	report := recon.Report{Target: target, Info: "txt records"}
	records, err := net.LookupTXT(target)
	if err != nil {
		report.Err = err
		return report
	}
	report.Results = append(report.Results, records...)
	return report
}
