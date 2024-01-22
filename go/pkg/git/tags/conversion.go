package tags

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
)

// RefsToStrings converts a slice of *plumbing.Reference into a slice of strings.
// Each string is in the format 'ReferenceName CommitHash' with space as seperator
func RefsToStrings(refs []*plumbing.Reference) (strs []string, err error) {
	strs = []string{}
	for _, ref := range refs {
		strs = append(strs, fmt.Sprintf("%s %s", ref.Name().String(), ref.Hash().String()))
	}
	return
}

// RefsToShortNames converts a slice of *plumbing.Reference into a slice of strings.
// For this case, we use the short form of the ref name.
// - Uses RefStringify under the hood
func RefsToShortNames(refs []*plumbing.Reference) (strs []string, err error) {
	return RefStringify(refs, nil)
}

// RefStringify converts a slice of *plumbing.Reference into a slice of strings.
// For this case, we use the short form of the ref name.
//
//   - Invisigned to be used as a quick wrapper around .At() and similar
//   - If error is passed then this is returned immediately
func RefStringify(refs []*plumbing.Reference, e error) (strs []string, err error) {
	if e != nil {
		err = e
		return
	}
	strs = []string{}
	for _, ref := range refs {
		strs = append(strs, ref.Name().Short())
	}
	return
}

// StringToHash tries to convert into a plumbing.Hash or fails and returns error
func (t *Tags) StringToHash(str string) (*plumbing.Hash, error) {
	return t.repository.ResolveRevision(plumbing.Revision(str))
}
