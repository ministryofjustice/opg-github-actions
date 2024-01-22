/*
nexttag determines the next semver tag to use based on information passed
Requires an actual git repository to be present at repository_root to fetch existing tags.

Find commits that are present in `--commitish_b` history, but not in `--commitish_a`s and uses those to look for semver version bump triggers.

Version bumps are found in the message property and are - '#major', '#minor' or '#patch'.

If no version bump is found, the `--default_bump` is applied

Usage:

	next-tag [flags]

The flags are:

	--repository 		(required)
		Root of the git repository to fetch tag information from.
		This directory must exist and be an active repository already.
	--base 		(required)
		Base git reference (typically branch being merged into - 'main' or 'master' or the 'before' commit).
		Commits between this point and `--commitish_b` are diff'd and commit messages checked for version bump triggers.
			- destination_commitish from branch-name
	--head	(required)
		Head git reference (typically the active branch being merged or the 'after' commit).
		Commits that exist on at this reference, but not witin `--commitish_a` are used to find semver version bump triggers.
			- source_commitish from branch-name
	--prerelease		(default: 'false')
		String representation of a boolean.
		Determines if the next flag should be considered as a prerelease or not
	--prerelease-suffix	(default: 'beta')
		Contains the prerelease segment to use in the tag (when `--prerelease` is 'true').
	--last-prerelease
		String representation of a semver tag ('1.0.0-beta.1+build01').
		Used with `--last-release` and `--prerelease` to determine the base tag to increment on.
	--last-release
		String representation of a semver tag ('1.0.0-beta.1+build01').
		Used with `--last-prerelease` and `--prerelease` to determine the base tag to increment on.
	--with-v			(default: 'false')
		String representation of a boolean.
		Determines if the generated next_tag string will have a 'v' prefix or not.
	--default_bump		(default: 'patch', options: 'major', 'minor', 'patch')
		If there are no trigger tags found in the commits, then bump the next tag along by this type.
	--event_data_file
		Extra event data that can be passed to caputure triggers within the raised pull_request

# Example Usage

	next-tag \
	    --repository="./tmp" \
		--base="main" \
		--head="my-feature" \
	    --prerelease="true" \
	    --prerelease-suffix="myfeature" \
	    --last-prerelease="1.0.1-myfeature.0" \
		--last-release="1.0.0"
*/
package nexttag

import (
	"flag"
)

var Name = "next-tag"                                 // Command name (exposed)
var FlagSet = flag.NewFlagSet(Name, flag.ExitOnError) // FlagSet is the group of commands (exposed)

// Input arguments
var (
	repoDir          = FlagSet.String("repository", "", "Root directory of the repository to use.")                                                                  // path to the git repo root
	baseRef          = FlagSet.String("base", "", "Used for commit comparisons (typically 'main')")                                                                  // git ref locations for comparison
	headRef          = FlagSet.String("head", "", "Used for commit comparisons (typically the feature-branch)")                                                      // git ref locations for comparison
	prerelease       = FlagSet.String("prerelease", "false", "If set, looks for pre-release tag patterns (v1.1.1-${suffix}.${count})")                               // determines if this is a release
	prereleaseSuffix = FlagSet.String("prerelease-suffix", "beta", "If prerelease is set, this string is used as the ${suffix} in the tag pattern. (Default: beta)") // used to search and generate a prerelease tag
	lastPrerelease   = FlagSet.String("last-prerelease", "", "Last tag for this prerelease branch. Used as a basis for working out the next")                        // semver last prerelease
	lastRelease      = FlagSet.String("last-release", "", "Last production release semver tag")                                                                      // the last production release tag
	withV            = FlagSet.String("with-v", "false", "If set, forces adding a prefix to the generated tag")                                                      // decides if the new tag has a v prefix
	defaultBump      = FlagSet.String("default-bump", "patch", "Default version trigger")                                                                            // semver bump
)
