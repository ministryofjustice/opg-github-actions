package commits

import (
	"opg-github-actions/pkg/semver"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// VersionBumpsInCommits looks for version bump triggers (#major, #minor, #patch) in the commit objects
//
// Converts the commits to strings and then uses semver.VersionBumpCount
func VersionBumpsInCommits(commits []*object.Commit, defaultInc semver.Increment) (counters *semver.IncrementCounters) {

	strs := []string{}
	for _, c := range commits {
		msg := strings.ToLower(c.Message)
		strs = append(strs, msg)
	}

	counters = semver.VersionBumpCount(strs, defaultInc)
	return
}
