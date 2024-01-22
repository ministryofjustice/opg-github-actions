package tags

import (
	"strings"

	"facette.io/natsort"
	"github.com/go-git/go-git/v5/plumbing"
)

// Sort will take set of tags references, convert to strings,
// use natsort to order them and then convert back to references
func Sort(tags []*plumbing.Reference) (sorted []*plumbing.Reference) {
	sorted = []*plumbing.Reference{}

	toSort, _ := RefsToStrings(tags)
	natsort.Sort(toSort)

	for _, str := range toSort {
		info := strings.Split(str, " ")
		ref := plumbing.NewReferenceFromStrings(info[0], info[1])
		sorted = append(sorted, ref)
	}
	return
}

func Join(tags []*plumbing.Reference) (joined string) {
	strs, _ := RefsToShortNames(tags)
	joined = strings.Join(strs, ", ")
	return
}
