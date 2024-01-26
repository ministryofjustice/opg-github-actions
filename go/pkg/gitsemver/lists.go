package gitsemver

import (
	"opg-github-actions/pkg/semver"

	"github.com/go-git/go-git/v5/plumbing"
)

func Releases(tags []*plumbing.Reference) (releases []*plumbing.Reference) {
	releases = []*plumbing.Reference{}

	for _, t := range tags {
		s, e := semver.New(t.Name().Short())
		if e == nil && !s.IsPrerelease() {
			releases = append(releases, t)
		}
	}
	return
}

func Prereleases(tags []*plumbing.Reference) (prereleases []*plumbing.Reference) {
	prereleases = []*plumbing.Reference{}

	for _, t := range tags {
		s, e := semver.New(t.Name().Short())
		if e == nil && s.IsPrerelease() {
			prereleases = append(prereleases, t)
		}
	}
	return
}

func MatchingPrereleases(match string, tags []*plumbing.Reference) (matches []*plumbing.Reference) {
	matches = []*plumbing.Reference{}

	for _, t := range tags {
		s, e := semver.New(t.Name().Short())
		if e == nil && s.IsPrereleaseMatch(match) {
			matches = append(matches, t)
		}
	}

	return
}
