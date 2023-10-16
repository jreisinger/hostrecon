package main

import (
	"slices"

	"github.com/jreisinger/recon"
)

type Names map[string]recon.Reconnoiterer

func (names Names) GetAll() []string {
	var all []string
	for k := range names {
		all = append(all, k)
	}
	slices.Sort(all)
	return all
}

func (names Names) GetAllReconnoiterers() []recon.Reconnoiterer {
	var all []recon.Reconnoiterer
	for _, v := range names {
		all = append(all, v)
	}
	return all
}
