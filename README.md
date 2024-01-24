Recon is a small CLI tool (and a Go package) for gathering public information about network hosts. It's meant to be fast, easy to use, and easy to extend (by adding reconnoiterers).

Install

```
# optional; to install inside an ephemeral container
docker run --rm -it golang /bin/bash

go install github.com/jreisinger/recon/cmd/recon@latest
```

Run all reconnoiterers against one host

```
$ recon example.com
example.com: txt records: v=spf1 -all, wgyf8z8cgvm2qmxpnbnldrcltvk4xqfn
example.com: geolocation: 93.184.216.34: London - GB, 2606:2800:220:1:248:1893:25c8:1946: New York - US
example.com: http version: HTTP/2.0
example.com: open tcp ports: 80, 443
example.com: certificate issuer: DigiCert Inc
example.com: tls version: TLS 1.3
example.com: ip addresses: 93.184.216.34, 2606:2800:220:1:248:1893:25c8:1946
example.com: name servers: b.iana-servers.net, a.iana-servers.net
example.com: certificate authority: DigiCert Inc
```

Run just one of the reconnoiterers (`-r`) against multiple hosts

```
$ recon -r ips example.com golang.org
example.com: ip addresses: 93.184.216.34, 2606:2800:220:1:248:1893:25c8:1946
golang.org: ip addresses: 142.251.36.81, 2a00:1450:4014:80a::2011
```

Output JSON (`-j`)

```
$ recon -r ips -j example.com golang.org
{"target":"example.com","info":"ip addresses","results":["93.184.216.34","2606:2800:220:1:248:1893:25c8:1946"]}
{"target":"golang.org","info":"ip addresses","results":["142.251.36.81","2a00:1450:4014:80a::2011"]}
```

Embed within a pipeline

```
$ subfinder --silent -d example.net | recon -r ips -j 2> /dev/null | jq '.results[]' -r | checkip 2> /dev/null                                            
--- 2606:2800:220:1:248:1893:25c8:1946 ---
db-ip.com       New York, United States
iptoasn.com     EDGECAST
is on AWS       false
malicious       0% (0/5) ✅
--- 93.184.216.34 ---
db-ip.com       London, United Kingdom
iptoasn.com     EDGECAST
is on AWS       false
tls             TLS 1.3, exp. 2024/02/13, www.example.org, example.net, example.edu, example.com, example.org, www.example.com, www.example.edu, www.example.net
malicious       14% (1/7) ✅
```

- [subfinder](https://github.com/projectdiscovery/subfinder)
- [jq](https://jqlang.github.io/jq/)
- [checkip](https://github.com/jreisinger/checkip)
