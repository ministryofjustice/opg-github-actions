package tags

import (
	"log/slog"

	"github.com/go-git/go-git/v5/plumbing"
)

// All returns all tag references found in the repository as a plumbing reference
func (t *Tags) All() (tags []*plumbing.Reference, err error) {
	slog.Info("Getting all tags for repository")

	tags = []*plumbing.Reference{}
	iter, _ := t.repository.Tags()

	iter.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, ref)
		return nil
	})

	return
}

// AllAsStrings provides a wrapper to return the tag data as strings directly
// rather than havign to pass in via RefsToShortNames in many places
func (t *Tags) AllAsStrings() (tags []string, err error) {
	all, err := t.All()
	if err != nil {
		return
	}
	tags, err = RefsToShortNames(all)
	return
}

// At limits the tags returned to be those with a commit hash matching the value
// passed
func (t *Tags) At(reference *plumbing.Hash) (tags []*plumbing.Reference, err error) {
	tags = []*plumbing.Reference{}
	all, err := t.All()

	if err != nil {
		return
	}

	for _, ref := range all {
		if ref.Hash().String() == reference.String() {
			tags = append(tags, ref)
		}
	}

	return
}
