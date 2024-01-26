package semver

import (
	"fmt"
	"log/slog"

	"facette.io/natsort"
)

// Sort orders semvers in ascending order as a string
//
// Prefixes are not removed, so they will be sorted
// based on value including prefix
func Sort(versions []*Semver) (sorted []*Semver) {
	slog.Info(fmt.Sprintf("semver Sort: sorting [%d] versions", len(versions)))
	sorted = []*Semver{}
	// convert to strings
	toSort := []string{}
	for _, v := range versions {
		if Valid(v.String()) {
			toSort = append(toSort, v.String())
		}
	}
	// sort
	natsort.Sort(toSort)
	// back to semvers
	for _, s := range toSort {
		sorted = append(sorted, Must(New(s)))
	}
	slog.Debug(fmt.Sprintf("semver Sort: returned [%d] versions", len(sorted)))
	return
}
