package tcp

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type PortsToScan []int

func (ports *PortsToScan) String() string {
	var out []string
	for _, f := range *ports {
		out = append(out, strconv.Itoa(f))
	}
	return strings.Join(out, ",")
}

func (ports *PortsToScan) Set(value string) error {
	*ports = PortsToScan{} // remove defaults
	for _, v := range strings.Split(value, ",") {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("can't turn %v into int: %v", v, err)
		}
		*ports = append(*ports, n)
	}
	return nil
}

func PortsToScanFlag(name string, value []int, usage string) *PortsToScan {
	p := PortsToScan(value)
	flag.Var(&p, name, usage)
	return &p
}
