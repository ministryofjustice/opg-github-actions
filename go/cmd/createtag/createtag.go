/*
createtag tries to make a tag at the commit references passed in.

Requires an actual git repository to be present at directory to
fetch existing tags.

If `--tag-name` already exists then a similar tag will be generated and created instead.

Usage:

	create-tag [flags]

The flags are:

	--repository		(required)
		Root of the git repository to fetch tag information from.
		This directory must exist and be an active repository already.
	--commitish		(required)
		Git reference of where we'll generate the tag.
	--tag-name		(required)
		Name of the tag to create
	--regen			(default: 'true')
		When true, if the tag-name passed clashes with an existing a tag a new, similar tag will be created.
			- On a clashing prerelease semver tag, the prerelease segment will be adjusted
			- On a clashing release tag, the patch version will be bumped
			- For a non-semver tag, a random string is append to the end
	--push			(default: 'false')
		When true, the created tag will then be pushed to the remote.
		To do this, it uses the GITHUB_TOKEN environment variable directly.

# Example Usage

	create-tag \
	    --repository="./tmp" \
		--commitish="my-feature-branch" \
		--tag_name="1.0.1-beta.0"
*/
package createtag

import (
	"flag"
)

var Name = "create-tag"                               // Command name (exposed)
var FlagSet = flag.NewFlagSet(Name, flag.ExitOnError) // FlagSet is the group of commands (exposed)
// Input arguments
var (
	repoDir      = FlagSet.String("repository", "", "Root directory of the repository to use.") // repository_directory is path to the git repo root
	commitish    = FlagSet.String("commitish", "", "Git reference to create tag_name at")
	tagName      = FlagSet.String("tag-name", "", "Tag to create")
	regenTag     = FlagSet.String("regen", "true", "Flag to determine if when tag name exists we should create a similar one and use that.")
	pushToRemote = FlagSet.String("push", "false", "When true, push the created tag to the remote.")
)
