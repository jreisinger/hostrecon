package subdomains

import (
	"hostrecon"
	"os/exec"
	"strings"
)

type Subdomains struct{}

func (s Subdomains) Recon(host string) hostrecon.Info {
	info := hostrecon.Info{Host: host, Kind: "subdomains"}
	out, err := exec.Command("subfinder", "--silent", "-d", host).Output()
	if err != nil {
		info.Err = err
		return info
	}
	var subdomains []string
	for _, s := range strings.Split(string(out), "\n") {
		if s != "" {
			subdomains = append(subdomains, s)
		}
	}
	info.Info = subdomains
	return info
}
