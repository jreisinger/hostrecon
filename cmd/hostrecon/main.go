package main

import (
	"flag"
	"time"

	"github.com/jreisinger/hostrecon"
	"github.com/jreisinger/hostrecon/dns"
	"github.com/jreisinger/hostrecon/geo"
	"github.com/jreisinger/hostrecon/http"
	"github.com/jreisinger/hostrecon/subdomains"
	"github.com/jreisinger/hostrecon/tls"
)

func main() {
	c := flag.Int("c", 10, "number of recons performed concurrently")
	j := flag.Bool("j", false, "json output")
	p := flag.String("p", "443", "http server port")
	t := flag.Duration("t", 3*time.Second, "http connection timeout")
	x := flag.Int("x", 3, "max number of info elements; negative means all")
	flag.Parse()

	recon := hostrecon.New(*c, *j)
	if len(flag.Args()) > 0 {
		recon.Hosts(flag.Args())
	}
	recon.MaxInfoElems(*x)
	recon.Perform([]hostrecon.Scout{
		dns.Cname{},
		dns.IpAddr{},
		dns.Mx{},
		dns.Ns{},
		dns.Txt{},
		geo.DbIpCom{},
		http.Version{Port: *p, Timeout: *t},
		tls.CA{Port: *p, Timeout: *t},
		tls.Issuer{Port: *p, Timeout: *t},
		tls.Version{Port: *p, Timeout: *t},
		subdomains.Subdomains{},
	})
}
