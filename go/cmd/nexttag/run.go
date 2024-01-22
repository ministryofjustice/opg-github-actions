package nexttag

import (
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/git/commits"
	"opg-github-actions/pkg/safestrings"
	"opg-github-actions/pkg/semver"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func process(
	commitDiff []*object.Commit,
	prerelease bool,
	prereleaseSuffix string,
	lastPrerelease string,
	lastRelease string,
	defaultBump semver.Increment,
	withV bool,
) (output map[string]string, err error) {

	var (
		latestRelease    *semver.Semver = nil
		latestPrerelease *semver.Semver = nil
		next             *semver.Semver = nil
		by               semver.Increment
	)
	output = map[string]string{}

	// get counters and what to bump by
	counter := commits.VersionBumpsInCommits(commitDiff, defaultBump)
	if counter.Major > 0 {
		by = semver.Major
	} else if counter.Minor > 0 {
		by = semver.Minor
	} else if counter.Patch > 0 {
		by = semver.Patch
	}

	latestRelease = semver.Must(semver.New(lastRelease))
	latestPrerelease, _ = semver.New(lastPrerelease)

	next, err = semver.Next(
		latestPrerelease, latestRelease,
		prerelease, prereleaseSuffix,
		*counter,
	)

	if withV {
		next.SetPrefix('v')
	} else {
		next.RemovePrefix()
	}

	output["bumped_by"] = string(by)
	output["last_release"] = fmt.Sprintf("%s", latestRelease.String())
	output["last_prerelease"] = fmt.Sprintf("%s", latestPrerelease.String())
	output["majors"] = fmt.Sprintf("%d", counter.Major)
	output["minors"] = fmt.Sprintf("%d", counter.Minor)
	output["patches"] = fmt.Sprintf("%d", counter.Patch)
	output["next_tag"] = fmt.Sprintf("%s", next.String())
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

	commitSet, err := commits.New(*repoDir)
	if err != nil {
		slog.Error("commitSet failed.")
		return
	}

	base, err := commitSet.StrToReference(*baseRef)
	if err != nil {
		slog.Error(fmt.Sprintf("base reference failed [%s]", *baseRef))
		return
	}
	head, err := commitSet.StrToReference(*headRef)
	if err != nil {
		slog.Error(fmt.Sprintf("head reference failed [%s]", *headRef))
		return
	}

	diff, err := commitSet.DiffBetween(base.Hash(), head.Hash())
	if err != nil {
		slog.Error("commit diff failed")
		return
	}

	isPre, _ := safestrings.ToBool(*prerelease)
	prefix, _ := safestrings.ToBool(*withV)
	bump := semver.MustIncrement(semver.NewIncrement(*defaultBump))

	output, err = process(
		diff,
		isPre, *prereleaseSuffix,
		*lastPrerelease, *lastRelease,
		bump,
		prefix,
	)
	if err != nil {
		return
	}
	output["directory"] = *repoDir
	output["head"] = *headRef
	output["base"] = *baseRef
	output["prerelease"] = fmt.Sprintf("%t", isPre)
	output["prerelease_suffix"] = *prereleaseSuffix
	output["with_v"] = fmt.Sprintf("%t", prefix)
	output["default_bump"] = *defaultBump

	return
}
