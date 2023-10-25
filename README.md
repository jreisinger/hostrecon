Recon is a small CLI tool (and a Go package) for gathering public information about network hosts. It's meant to be fast, easy to use, and easy to extend.

Installation

```
# optional; to install inside an ephemeral container
docker run --rm -it golang /bin/bash

go install github.com/jreisinger/recon/cmd/recon@latest
```

Usage

```
# run all reconnoiterers against one host
$ recon example.com
example.com: tls version: TLS 1.3
example.com: ip addresses: 93.184.216.34, 2606:2800:220:1:248:1893:25c8:1946
example.com: geolocation: 93.184.216.34: New York - US, 2606:2800:220:1:248:1893:25c8:1946: New York - US
example.com: name servers: a.iana-servers.net, b.iana-servers.net
example.com: txt records: wgyf8z8cgvm2qmxpnbnldrcltvk4xqfn, v=spf1 -all
example.com: http version: HTTP/2.0
example.com: open tcp ports: 80, 443
example.com: certificate authority: DigiCert Inc
example.com: certificate issuer: DigiCert Inc

# run just one of the reconnoiterers against multiple hosts and output json
$ recon -r ips -j perl.org golang.org
{"target":"golang.org","desc":"ip addresses","results":["142.251.36.81","2a00:1450:4014:80a::2011"]}
{"target":"perl.org","desc":"ip addresses","results":["139.178.67.96"]}

# embed within a pipeline
$ subfinder --silent -d ibm.com | recon -r tlsver 2> /dev/null | grep -v 'TLS 1.3'
ns190.name.cloud.ibm.com: tls version: TLS 1.2
ns027.name.cloud.ibm.com: tls version: TLS 1.2
admin-api.3c10-1ca316c5.eu-de.apiconnect.cloud.ibm.com: tls version: TLS 1.2
<...>
```
