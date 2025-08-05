package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"opg-github-actions/action/internal/commits"
	"opg-github-actions/action/internal/logger"
	"opg-github-actions/action/internal/repo"
	"opg-github-actions/action/internal/semver"
	"opg-github-actions/action/internal/strs"
	"opg-github-actions/action/internal/tags"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const ErrNoBranchName string = "branch-name is required, but not found."

type Options struct {
	RepositoryDirectory    string // Directory where the dit repo is
	Prerelease             bool   // if this is a prerelease of a full release
	PrereleaseSuffixLength int    // length of the prerelease suffix
	DefaultBranch          string // default branch name - generally main, used to compare commits against
	BranchName             string // branch name is used as the prerelease suffix
	DefaultBump            string // what to increment the semver by (major, minor, patch)
	ExtraContent           string // content from pull request title / body where their might be extra #major content
	WithoutPrefix          bool
	TestMode               bool
}

var runOptions *Options = newRunOptions(&Options{DefaultBranch: "main"})

// newRunOptions helper to return default options merged with
// overwrites
func newRunOptions(in *Options) (opts *Options) {
	opts = &Options{
		RepositoryDirectory:    "",
		Prerelease:             false,
		PrereleaseSuffixLength: 14,
		BranchName:             "",
		DefaultBranch:          "",
		DefaultBump:            string(semver.PATCH),
		ExtraContent:           "",
		WithoutPrefix:          false,
		TestMode:               true,
	}
	if in != nil {
		if in.RepositoryDirectory != "" {
			opts.RepositoryDirectory = in.RepositoryDirectory
		}

		if in.PrereleaseSuffixLength > 0 {
			opts.PrereleaseSuffixLength = in.PrereleaseSuffixLength
		}
		if in.BranchName != "" {
			opts.BranchName = in.BranchName
		}
		if in.DefaultBranch != "" {
			opts.DefaultBranch = in.DefaultBranch
		}
		if in.DefaultBump != "" {
			opts.DefaultBump = in.DefaultBump
		}
		if in.ExtraContent != "" {
			opts.ExtraContent = in.ExtraContent
		}
		opts.Prerelease = in.Prerelease
		opts.TestMode = in.TestMode
		opts.WithoutPrefix = in.WithoutPrefix
	}

	return
}

// getExistingSemvers fetches all git tags and converts valid entries into semvers
// which are then used to work out the next value
func getExistingSemvers(lg *slog.Logger, repository *git.Repository) (semvers []*semver.Semver, err error) {

	var gittags []*plumbing.Reference // all tags in the repo
	// get the tags
	if gittags, err = tags.All(repository); err != nil {
		lg.Error("error getting tags from repository", "err", err.Error())
		return
	}
	// get the semvers from the tags
	if semvers, err = semver.FromGitRefs(gittags); err != nil {
		lg.Error("error getting semvers from tags", "err", err.Error())
		return
	}
	return
}

// getSemverToUse looks at the semvers and the options passed along and determines if we should be used prerelease or release
// semver tag and handles prefix usage.
func getSemverToUse(lg *slog.Logger, semvers []*semver.Semver, bump semver.Increment, options *Options) (use *semver.Semver) {
	use = &semver.Semver{}
	// decide if we do prerelease or not based on input
	if options.Prerelease {
		// run safe on the branch name for prerelease usage
		suffix, _ := strs.Safe(options.BranchName, options.PrereleaseSuffixLength)
		use = semver.Prerelease(lg, semvers, bump, suffix)
	} else {
		use = semver.Release(lg, semvers, bump)
	}

	// setup the prefix
	if options.WithoutPrefix {
		use.Prefix = ""
	} else {
		use.Prefix = "v"
	}

	return
}

// createAndPushTag handles the logic of creating and then pushing tags.
//
// If we are in test mode, or there is no semver increment to do, then
// we immediately return and createdTag will be nil
//
// If there is no remote on the repository (ie a locally created repo)
// then the tag is created, but not pushed
func createAndPushTag(
	lg *slog.Logger,
	repository *git.Repository,
	use *semver.Semver,
	bump semver.Increment,
	token string,
	options *Options) (createdTag *plumbing.Reference, err error) {

	var (
		remotes, _ = repository.Remotes()
		auth       = &http.BasicAuth{
			Username: "opg-github-actions",
			Password: token,
		}
	)

	// we do nothing if this is in test mode or we're using no bumping
	if options.TestMode || bump == semver.NO_BUMP {
		return
	}

	// try to create the tag locally
	createdTag, err = tags.Create(repository, use.String(), use.GitRef.Hash())
	if err != nil {
		err = errors.Join(fmt.Errorf("failed to create a tag"), err)
		return
	}

	// if we have some remotes, push
	if len(remotes) > 0 {
		err = tags.Push(repository, auth)
		if err != nil {
			err = errors.Join(fmt.Errorf("failed to push tags"), err)
			return
		}
	}

	return
}

// Run handles gluing together the process of creating a new semver tag from the git repository and outputting the created
// values.
//
//   - Generate a repository object from the directory path arguments (or returns error)
//   - Finds all existing semver formatted git tags in the repository
//   - Finds the git sha / hash reference for the configured default branch and the currently checked out location
//   - Finds all commits that exist in the currently checked out location, but not in default branch tree - these are the new ones
//     -- Merges the extra-content argument into this data (pull request details)
//   - Looks at the new commits for #major|minor|patch content to determine the semver increment
//   - Works out the new new tag
//   - Creates and pushes the tag
//   - Outputs data
//
// If you have enabled test mode (via `--test`) the tag will not be created or pushed.
// If there are no new commits, or no commits with #major|minor|patch then the default increment (`--bump`)
// will be used.
func Run(lg *slog.Logger, options *Options) (result map[string]string, err error) {
	var (
		repository    *git.Repository                                             // the object for this repo
		semvers       []*semver.Semver                                            // all valid semver tags in the repo
		use           *semver.Semver                                              // the semver to use for the prerelease / release
		defaultBranch *plumbing.Reference                                         // git ref for the default branch
		currentCommit *plumbing.Reference                                         // git sha / ref for where the git repo currently is
		createdTag    *plumbing.Reference                                         // the new semver tag thats been created
		newCommits    []*object.Commit                                            // all commits that exist in the ref
		bump          semver.Increment    = semver.Increment(options.DefaultBump) // default increment
		token         string              = os.Getenv("GH_TOKEN")                 // github auth token for pushing to the remote
	)
	result = map[string]string{}

	if options.Prerelease && options.BranchName == "" {
		err = fmt.Errorf(ErrNoBranchName)
		return
	}
	// generate a repo
	if repository, err = repo.FromDir(options.RepositoryDirectory); err != nil {
		lg.Error("error creating repository from directory", "err", err.Error(), "dir", options.RepositoryDirectory)
		return
	}

	// get the semvers from the tags
	semvers, err = getExistingSemvers(lg, repository)
	if err != nil {
		return
	}
	// get the default branch info
	if defaultBranch, err = commits.FindReference(lg, repository, options.DefaultBranch); err != nil {
		lg.Error("error getting git reference for default branch", "err", err.Error(), "default_branch", options.DefaultBranch)
		return
	}
	// get info on the current commit
	if currentCommit, err = commits.FindReference(lg, repository, options.BranchName); err != nil {
		lg.Error("error getting git reference for branch", "err", err.Error(), "branch", options.BranchName)
		return
	}
	// find new commits that are in the current tree, but not
	if newCommits, err = commits.DiffBetween(lg, repository, defaultBranch.Hash(), currentCommit.Hash()); err != nil {
		lg.Error("error commits between references", "err", err.Error(), "base", defaultBranch.Hash().String(), "head", currentCommit.Hash().String())
		return
	}

	lg.Debug("found commits", "len", len(newCommits))

	// add content to the commit list
	if options.ExtraContent != "" {
		newCommits = append(newCommits, &object.Commit{Hash: plumbing.ZeroHash, Message: options.ExtraContent})
	}

	// look for bump in the commits,
	foundBump := semver.GetBumpFromCommits(newCommits, bump)
	if len(newCommits) > 0 && foundBump != "" {
		bump = foundBump
	}

	use = getSemverToUse(lg, semvers, bump, options)
	// set the git ref to the current place
	use.GitRef = currentCommit
	// create and try to push tags
	createdTag, err = createAndPushTag(lg, repository, use, bump, token, options)

	result = map[string]string{
		"tag":     use.String(),
		"hash":    use.GitRef.Hash().String(),
		"test":    fmt.Sprintf("%t", options.TestMode),
		"created": fmt.Sprintf("%t", (createdTag != nil)),
		"bump":    string(bump),
	}

	return
}

// init does the setup of args
func init() {
	flag.StringVar(&runOptions.RepositoryDirectory, "directory", runOptions.RepositoryDirectory, "The directory path of the git repository.")
	// branch details
	flag.StringVar(&runOptions.BranchName, "branch", runOptions.BranchName, "The current branch name to use for prerelease suffixes")
	flag.StringVar(&runOptions.DefaultBranch, "default-branch", runOptions.DefaultBranch, "The default branch name for this repo - used for commit comparisons")
	// prerelease related options
	flag.BoolVar(&runOptions.Prerelease, "prerelease", runOptions.Prerelease, "Set to true to generate a prerelease version.")
	flag.IntVar(&runOptions.PrereleaseSuffixLength, "prerelease-suffix-length", runOptions.PrereleaseSuffixLength, "Set the max length to use for tag suffixes")
	// Semver increments
	flag.StringVar(&runOptions.DefaultBump, "bump", runOptions.DefaultBump, "The default value to increment semver by if no comment if found. (default: patch)")
	// use a prefix?
	flag.BoolVar(&runOptions.WithoutPrefix, "without-prefix", runOptions.WithoutPrefix, "Use to disable prefix usage.")
	// test mode - disables creating tags
	flag.BoolVar(&runOptions.TestMode, "test", runOptions.TestMode, "Set to true to disable creating tag.")
	//
	flag.StringVar(&runOptions.ExtraContent, "extra-content", runOptions.ExtraContent, "Additional content that might also contain # references")
}

func main() {
	var lg *slog.Logger = logger.New("INFO", "TEXT")
	// process the arguments and fetch the fallback value from environment values
	flag.Parse()

	// run the command
	res, err := Run(lg, runOptions)
	if err != nil {
		lg.Error(err.Error())
		os.Exit(1)
	}
	logger.Result(lg, res)

}

func debug[T any](item T) {
	bytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Printf("%+v\n", string(bytes))
}
