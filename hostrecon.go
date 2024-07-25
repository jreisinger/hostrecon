// Package hostrecon obtains information about network hosts.
package hostrecon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// Scout obtains some kind of information about a network host.
type Scout interface {
	Recon(host string) Info
}

// Information about a network host.
type Info struct {
	Host string   `json:"host"`
	Kind string   `json:"kind"`
	Info []string `json:"info"`
	Err  error    `json:"-"`
}

type Recon struct {
	concurrency  int
	jsonOutput   bool
	hosts        io.Reader
	out          io.Writer
	err          io.Writer
	maxInfoElems int
}

func New(concurrency int, jsonOutput bool) *Recon {
	r := &Recon{
		concurrency:  concurrency,
		jsonOutput:   jsonOutput,
		hosts:        os.Stdin,
		out:          os.Stdout,
		err:          os.Stderr,
		maxInfoElems: -1,
	}
	return r
}

// Err defines where errors should go. Default is os.Stderr.
func (r *Recon) Err(w io.Writer) {
	r.err = w
}

// Hosts defines the hosts to perform recon on. Default is os.Stdin.
func (r *Recon) Hosts(hosts []string) {
	s := strings.Join(hosts, "\n")
	r.hosts = strings.NewReader(s)
}

// Out defines where recon info should go. Default is os.Stdout.
func (r *Recon) Out(w io.Writer) {
	r.out = w
}

// MaxInfoElems defines the maximum number elements in the Info.Info string
// slice. Negative number means all.
func (r *Recon) MaxInfoElems(n int) {
	r.maxInfoElems = n
}

func (recon *Recon) Perform(scouts []Scout) {
	var wg sync.WaitGroup

	in := make(chan string)

	wg.Add(1)
	go func() {
		host := bufio.NewScanner(recon.hosts)
		for host.Scan() {
			in <- host.Text()
		}
		if err := host.Err(); err != nil {
			fmt.Fprintf(recon.err, "hostrecon: %s", err)
		}
		close(in)
		wg.Done()
	}()

	out := make(chan Info)

	for i := 0; i < recon.concurrency; i++ {
		wg.Add(1)
		go func() {
			for host := range in {
				for _, s := range scouts {
					out <- s.Recon(host)
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for report := range out {
		recon.write(report)
	}
}

// write writes the report to the output writer.
func (recon *Recon) write(info Info) {
	if info.Err != nil {
		fmt.Fprintf(recon.err, "hostrecon: %s: %s\n", info.Kind, info.Err)
		return
	}

	elems, remaining := maxElems(info.Info, recon.maxInfoElems)
	if remaining > 0 {
		elems = append(elems, fmt.Sprintf("... %d more", remaining))
	}
	info.Info = elems

	if recon.jsonOutput {
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Fprint(recon.err, err)
		}
		fmt.Fprintf(recon.out, "%s\n", data)
	} else {
		fmt.Fprintf(recon.out, "%s: %s: %s\n", info.Host, info.Kind, strings.Join(info.Info, ", "))
	}
}

func maxElems(ss []string, n int) (elems []string, remaining int) {
	if n < 0 || n > len(ss) {
		return ss, 0
	}
	return ss[:n], len(ss) - n
}
