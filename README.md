Recon is a CLI tool (and Go package) for gathering public information about network hosts. It's fast and easy to use, and easy to extend.

```
$ go install ./cmd/recon.go

$ recon example.com
example.com: ip addreses: 93.184.216.34, 2606:2800:220:1:248:1893:25c8:1946
example.com: name servers: b.iana-servers.net, a.iana-servers.net
example.com: txt records: wgyf8z8cgvm2qmxpnbnldrcltvk4xqfn, v=spf1 -all
example.com: geolocation: 93.184.216.34: New York - US, 2606:2800:220:1:248:1893:25c8:1946: New York - US
example.com: open tcp ports: 80, 443
example.com: certificate authority: DigiCert Inc
example.com: certificate issuer: DigiCert Inc
example.com: tls version: TLS 1.3

$ recon -j example.com | jq 'select(.desc == "geolocation")'
{
  "target": "example.com",
  "desc": "geolocation",
  "results": [
    "93.184.216.34: New York - US",
    "2606:2800:220:1:248:1893:25c8:1946: New York - US"
  ]
}
```
