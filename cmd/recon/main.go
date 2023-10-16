package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jreisinger/recon"
	"github.com/jreisinger/recon/dns"
	"github.com/jreisinger/recon/geo"
	"github.com/jreisinger/recon/http"
	"github.com/jreisinger/recon/tcp"
	"github.com/jreisinger/recon/tls"
)

var usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s - gather info about network hosts\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "options\n")
	flag.PrintDefaults()
}

var names = Names{
	"cname":   dns.Cname(),
	"ips":     dns.IPAddr(),
	"mx":      dns.MX(),
	"ns":      dns.NS(),
	"txt":     dns.TXT(),
	"geo":     geo.DBip(),
	"httpver": http.Version(),
	"ports":   tcp.OpenPorts(),
	"ca":      tls.CA(),
	"iss":     tls.Issuer(),
	"tlsver":  tls.Version(),
}

var (
	c = flag.Int("c", 10, "max hosts being reconned concurrently")
	j = flag.Bool("j", false, "json output")
	r = flag.String(
		"r",
		"",
		fmt.Sprintf(
			"run just one of the reconnoiterers: %s",
			strings.Join(names.GetAll(), " ")))
)

func main() {
	flag.Usage = usage
	flag.Parse()

	var reconnoiterers []recon.Reconnoiterer

	if *r == "" {
		reconnoiterers = names.GetAllReconnoiterers()
	} else {
		switch *r {
		case "cname":
			reconnoiterers = append(reconnoiterers, dns.Cname())
		case "ips":
			reconnoiterers = append(reconnoiterers, dns.IPAddr())
		case "mx":
			reconnoiterers = append(reconnoiterers, dns.MX())
		case "ns":
			reconnoiterers = append(reconnoiterers, dns.NS())
		case "txt":
			reconnoiterers = append(reconnoiterers, dns.TXT())
		case "geo":
			reconnoiterers = append(reconnoiterers, geo.DBip())
		case "httpver":
			reconnoiterers = append(reconnoiterers, http.Version())
		case "ports":
			reconnoiterers = append(reconnoiterers, tcp.OpenPorts())
		case "ca":
			reconnoiterers = append(reconnoiterers, tls.CA())
		case "iss":
			reconnoiterers = append(reconnoiterers, tls.Issuer())
		case "tlsver":
			reconnoiterers = append(reconnoiterers, tls.Version())
		default:
			fmt.Fprintf(
				os.Stderr,
				"recon: unknown reconnoiterer: %s: pick one from: %s\n",
				*r, strings.Join(names.GetAll(), " "))
			os.Exit(1)
		}
	}

	recon.NewRunner(
		recon.WithGoroutines(*c),
		recon.WithJsonOutput(*j),
		recon.WithTargets(flag.Args()...),
	).Run(reconnoiterers)
}
