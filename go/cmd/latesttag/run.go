package latesttag

import (
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/git/tags"
	"opg-github-actions/pkg/gitsemver"
	"opg-github-actions/pkg/safestrings"
	"opg-github-actions/pkg/semver"

	"github.com/go-git/go-git/v5/plumbing"
)

func process(allTags []*plumbing.Reference, prerelease bool, prereleaseSuffix string, branch string) (output map[string]string, err error) {
	var (
		releases              = []*plumbing.Reference{}
		prereleases           = []*plumbing.Reference{}
		allPrereleases        = []*plumbing.Reference{}
		lastPre        string = ""
		lastRelease    string = ""
	)
	output = map[string]string{}

	// get all releases
	releases = gitsemver.Releases(allTags)
	releases = tags.Sort(releases)
	// find the last release
	if len(releases) > 0 {
		last := releases[len(releases)-1]
		lastRelease = last.Name().Short()
	}
	allPrereleases = gitsemver.Prereleases(allTags)
	// if the prerelease is set as well as the suffix,
	// then find all matching prereleases and
	// set the last one
	// - last determined using natural sort
	if prerelease && len(prereleaseSuffix) > 0 {
		prereleases = gitsemver.MatchingPrereleases(prereleaseSuffix, allTags)
		prereleases = tags.Sort(prereleases)
		if len(prereleases) > 0 {
			last := prereleases[len(prereleases)-1]
			lastPre = last.Name().Short()
		}
	}
	// if either the prerelease or a release has a prefix, then flag as true
	hasPrefix := false
	if semver.HasPrefix(lastPre) || semver.HasPrefix(lastRelease) {
		hasPrefix = true
	}

	output["all_releases"] = tags.Join(releases)
	output["all_prereleases"] = tags.Join(allPrereleases)
	output["relevent_prereleases"] = tags.Join(prereleases)
	output["with_v"] = fmt.Sprintf("%t", hasPrefix)

	output["last_release"] = lastRelease
	output["last_prerelease"] = lastPre
	return
}

func Run(args []string) (output map[string]string, err error) {
	slog.Info("[" + Name + "] Run")
	FlagSet.Parse(args)

	// parse command arguments
	err = parseArgs()
	if err != nil {
		return
	}

	pre := safestrings.Safestring(*prerelease)
	isPre, err := pre.AsBool()
	if err != nil {
		return
	}

	tagset, err := tags.New(*repositoryDir)
	if err != nil {
		return
	}
	tags, err := tagset.All()
	if err != nil {
		return
	}
	output, err = process(tags, isPre, *prereleaseSuffix, *branchName)
	if err != nil {
		return
	}

	output["directory"] = *repositoryDir
	output["branch_name"] = *branchName
	output["prerelease"] = fmt.Sprintf("%t", isPre)
	output["prerelease_suffix"] = *prereleaseSuffix

	return
}
