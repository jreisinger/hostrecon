// Package recon gathers information about targets by running various reconnoiterers.
package recon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Reconnoiterer interface {
	Recon(target string) Report
}

type Report struct {
	Target  string   `json:"target"`
	Desc    string   `json:"desc"`
	Results []string `json:"results"`
	Err     error    `json:"-"`
}

type Runner struct {
	concurrentTargets int
	jsonOutput        bool
	targets           io.Reader
	output            io.Writer
	err               io.Writer
}

func NewRunner(opts ...option) *Runner {
	r := &Runner{
		concurrentTargets: 10,
		targets:           os.Stdin,
		output:            os.Stdout,
		err:               os.Stderr,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type option func(r *Runner)

func WithGoroutines(n int) option {
	return func(r *Runner) {
		r.concurrentTargets = n
	}
}

func WithJsonOutput(b bool) option {
	return func(r *Runner) {
		r.jsonOutput = b
	}
}

// If args is empty targets are read from STDIN.
func WithTargets(args ...string) option {
	return func(r *Runner) {
		if len(args) > 0 {
			s := strings.Join(args, "\n")
			r.targets = strings.NewReader(s)
		}
	}
}

func WithOutput(w io.Writer) option {
	return func(r *Runner) {
		r.output = w
	}
}

func WithErr(w io.Writer) option {
	return func(r *Runner) {
		r.output = w
	}
}

func (r *Runner) Run(rs []Reconnoiterer) {
	var wg sync.WaitGroup

	in := make(chan string)

	wg.Add(1)
	go func() {
		host := bufio.NewScanner(r.targets)
		for host.Scan() {
			in <- host.Text()
		}
		if err := host.Err(); err != nil {
			fmt.Fprintf(r.err, "recon: %s", err)
		}
		close(in)
		wg.Done()
	}()

	out := make(chan Report)

	for i := 0; i < r.concurrentTargets; i++ {
		wg.Add(1)
		go func() {
			for host := range in {
				wg.Add(1)
				go func(host string) {
					for _, r := range rs {
						out <- r.Recon(host)
					}
					wg.Done()
				}(host)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for report := range out {
		r.write(report)
	}
}

func (run *Runner) write(rep Report) {
	if rep.Err != nil {
		fmt.Fprintf(run.err, "recon: %s: %s: %s\n", rep.Target, rep.Desc, rep.Err)
		return
	}
	if len(rep.Results) > 0 {
		if run.jsonOutput {
			data, err := json.Marshal(rep)
			if err != nil {
				fmt.Fprint(run.err, err)
			}
			fmt.Fprintf(run.output, "%s\n", data)
		} else {
			fmt.Fprintf(run.output, "%s: %s: %s\n", rep.Target, rep.Desc, strings.Join(rep.Results, ", "))
		}
	}
}
