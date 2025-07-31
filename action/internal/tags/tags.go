package tags

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/maruel/natural"
)

const (
	ErrLikelyShallowRepository string = "No tags or other branches found; repository is probably a shallow clone."
)

type SortOrder bool

const (
	SORT_ASC  SortOrder = true
	SORT_DESC SortOrder = false
)

// branchCount counts all branches for the repo
func branchCount(repo *git.Repository) (count int) {
	var err error
	var iter storer.ReferenceIter
	count = 0

	iter, err = repo.Branches()
	if err != nil {
		return
	}

	iter.ForEach(func(ref *plumbing.Reference) error {
		count++
		return nil
	})

	return
}

// Return all tags for a repository
//
// Will check if the repository has some tags and more than one branch
// and will flag an error - so the first tag on brand new repo may cause
// unexpected error. TODO: how to improve shallow clone detection
func All(repo *git.Repository) (tags []*plumbing.Reference, err error) {
	var (
		iter storer.ReferenceIter
	)
	tags = []*plumbing.Reference{}

	iter, err = repo.Tags()
	if err != nil {
		return
	}

	iter.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, ref)
		return nil
	})

	// if no tags are found, check branches and if there is only one, this is likely a
	// shallow clone
	if len(tags) == 0 && branchCount(repo) == 1 {
		err = fmt.Errorf(ErrLikelyShallowRepository)
		return
	}

	return
}

// Strings converts list of tags / git references into a list of strings in the format:
//
// `refs/tags/tag-name f2bd01600ded17a4f3f2c8348443b33cd48c8902`
//
// Designed to be used with `All(path)` like `Strings(All(path))` to simplify
// usage - like the Must pattern.
//
// If there is an error passed, then return empty slice
func Strings(tags []*plumbing.Reference, err error) (refs []string) {

	refs = []string{}
	if err != nil {
		return
	}

	for _, tag := range tags {
		var nameAndHash = fmt.Sprintf("%s %s", tag.Name().String(), tag.Hash().String())
		refs = append(refs, nameAndHash)
	}

	return
}

// Refs converts list of tags / git references into a list of just the full reference in the format:
//
// `refs/tags/tag-name`
//
// Designed to be used with `All(path)` like `Refs(All(path))` to simplify usage - like
// the Must pattern.
//
// If there is an error passed, then return empty slice
func Refs(tags []*plumbing.Reference, err error) (refs []string) {

	refs = []string{}
	if err != nil {
		return
	}

	for _, tag := range tags {
		var name = tag.Name().String()
		refs = append(refs, name)
	}

	return
}

// ShortRefs converts list of tags / git references into a list of just the full reference in the format:
//
// `tag-name`
//
// Designed to be used with `All(path)` like `ShortRefs(All(path))` to simplify usage - like
// the Must pattern.
//
// If there is an error passed, then return empty slice
func ShortRefs(tags []*plumbing.Reference, err error) (shortRefs []string) {

	shortRefs = []string{}
	if err != nil {
		return
	}

	for _, tag := range tags {
		var name = tag.Name().Short()
		shortRefs = append(shortRefs, name)
	}

	return
}

// Sort will take set of tags references, convert to strings,
// use natural sorting to order them and then convert back to references
func Sort(tags []*plumbing.Reference, order SortOrder) (sorted []*plumbing.Reference) {
	sorted = []*plumbing.Reference{}

	toSort := Strings(tags, nil)
	if order == SORT_DESC {
		sort.Sort(sort.Reverse(natural.StringSlice(toSort)))
	} else {
		sort.Sort(natural.StringSlice(toSort))
	}
	// remove dups
	toSort = slices.Compact(toSort)

	for _, str := range toSort {
		info := strings.Split(str, " ")
		ref := plumbing.NewReferenceFromStrings(info[0], info[1])
		sorted = append(sorted, ref)
	}
	return
}
