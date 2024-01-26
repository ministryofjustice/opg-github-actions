package tags

import (
	"slices"
)

// Exists compares the tagName passed to all tags within the
// repository.
func (t *Tags) Exists(tagName string) (exists bool) {
	exists = false
	strs, _ := t.AllAsStrings()
	exists = slices.Contains(strs, tagName)
	return
}

// ExistsIn looks for the tagName within the set of tags passed
// - similar to exists, but allows pass what tags to compare
func (t *Tags) ExistsIn(tagName string, tags []string) (exists bool) {
	exists = slices.Contains(tags, tagName)
	return
}
