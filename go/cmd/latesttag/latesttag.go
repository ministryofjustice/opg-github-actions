/*
latesttag finds the latest tag relevant to the details passed and the last release tag. Typically used to find prerelease versions of the branch.
Requires an actual git repository to be present at --repository_root to fetch existing tags.

Usage:

	latest-tag [flags]

The flags are:

	--repository 		(required)
		Root of the git repository to fetch tag information from.
		This directory must exist and be an active repository already.
	--branch-name 		(required)
		Active branch name being used
	--prerelease		(default: 'false')
		String representation of a boolean.
		Determines if looking for prerelease branches.
		Can be overridden if branch_name matches a release_branch
	--prerelease-suffix		(default: 'beta')
		Contains the prerelease segment to use in the tag (when `--prerelease` is 'true').
	--release-branches		(default: 'main,master')
		List of branches that are considered to be a release

# Example Usage

	latest-tag \
	    --repository="./tmp" \
	    --prerelease="false" \
	    --prerelease-suffix="test" \
	    --branch="test" \
	    --release-branches="main"
*/
package latesttag

import (
	"flag"
)

var (
	Name             = "latest-tag"                            // Command name
	FlagSet          = flag.NewFlagSet(Name, flag.ExitOnError) // Grouped set of arguments
	repositoryDir    = FlagSet.String("repository", "", "Root directory of the repository to use.")
	prerelease       = FlagSet.String("prerelease", "false", "If set, looks for pre-release tag patterns (v1.1.1-${suffix}.${count})")
	prereleaseSuffix = FlagSet.String("prerelease-suffix", "beta", "If prerelease is set, this string is used as the ${suffix} in the tag pattern. (Default: beta)")
	branchName       = FlagSet.String("branch", "", "Current git branch name.")
	releaseBranches  = FlagSet.String("release-branches", "main,master", "Branches that would trigger this as a release.")
)
