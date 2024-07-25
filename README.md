Hostrecon is a CLI tool (and a Go package) for obtaining information about network hosts. It's meant to be fast, easy to use, and easy to extend.

Installation

```
go install github.com/jreisinger/hostrecon/cmd/hostrecon@latest
```

Usage

```
❯ hostrecon example.com
example.com: cname: 
example.com: ip addresses: 93.184.215.14, 2606:2800:21f:cb07:6820:80da:af6b:8b2c
example.com: mail servers: 
example.com: name servers: b.iana-servers.net, a.iana-servers.net
example.com: txt records: wgyf8z8cgvm2qmxpnbnldrcltvk4xqfn, v=spf1 -all
example.com: db-ip.com: 93.184.215.14: London GB, 2606:2800:21f:cb07:6820:80da:af6b:8b2c: Los Angeles (Playa Vista) US
example.com: http version: HTTP/2.0
example.com: tls ca: DigiCert Inc
example.com: tls cert issuer: DigiCert Inc
example.com: tls version: TLS 1.3
example.com: subdomains: yeskelco1.example.com, vm105673.example.com, a25462769725.example.com, ... 22246 more
```

```
❯ echo -e "example.com\nexample.net" | hostrecon -x -1 -j | jq -r 'select(.kind=="ip addresses")'
{
  "host": "example.net",
  "kind": "ip addresses",
  "info": [
    "93.184.215.14",
    "2606:2800:21f:cb07:6820:80da:af6b:8b2c"
  ]
}
{
  "host": "example.com",
  "kind": "ip addresses",
  "info": [
    "93.184.215.14",
    "2606:2800:21f:cb07:6820:80da:af6b:8b2c"
  ]
}
```

NOTE: you might need to install some 3rd party commands for some Scouts to start working.
