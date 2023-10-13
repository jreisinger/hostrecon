package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jreisinger/recon"
	"github.com/jreisinger/recon/dns"
	"github.com/jreisinger/recon/geo"
	"github.com/jreisinger/recon/http"
	"github.com/jreisinger/recon/tcp"
	"github.com/jreisinger/recon/tls"
)

var (
	c = flag.Int("c", 10, "concurrency")
	j = flag.Bool("j", false, "json output")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s - gather info about network hosts\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "options\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	recon.NewRunner(
		recon.WithGoroutines(*c),
		recon.WithJsonOutput(*j),
		recon.WithTargets(flag.Args()...),
	).Run([]recon.Reconnoiterer{
		dns.Cname(),
		dns.IPAddr(),
		dns.MX(),
		dns.NS(),
		dns.TXT(),
		geo.DBip(),
		http.Version(),
		tcp.OpenPorts(),
		tls.CA(),
		tls.Issuer(),
		tls.Version(),
	})
}
