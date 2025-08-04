package main

import (
	"encoding/json"
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
)

const ErrNoBranchName string = "branch-name is required, but not found."

type Options struct {
	RepositoryDirectory    string // Directory where the dit repo is
	Prerelease             bool   // if this is a prerelease of a full release
	PrereleaseSuffixLength int    // length of the prerelease suffix
	DefaultBranch          string // default branch name - generally main, used to compare commits against
	BranchName             string // branch name is used as the prerelease suffix
	DefaultBump            string // what to increment the semver by (major, minor, patch)
	WithPrefix             bool
	TestMode               bool
}

var runOptions *Options = &Options{
	RepositoryDirectory:    "",
	Prerelease:             false,
	PrereleaseSuffixLength: 12,
	BranchName:             "",
	DefaultBranch:          "main",
	DefaultBump:            string(semver.PATCH),
	WithPrefix:             true,
	TestMode:               true,
}

func Run(lg *slog.Logger, options *Options) (result map[string]string, err error) {
	var (
		repository    *git.Repository
		gittags       []*plumbing.Reference
		semvers       []*semver.Semver
		use           *semver.Semver
		defaultBranch *plumbing.Reference
		refBranch     *plumbing.Reference
		newCommits    []*object.Commit
		bump          semver.Increment = semver.Increment(options.DefaultBump)
	)
	result = map[string]string{}

	if options.Prerelease && options.BranchName == "" {
		err = fmt.Errorf(ErrNoBranchName)
		return
	}

	if repository, err = repo.FromDir(options.RepositoryDirectory); err != nil {
		lg.Error("error creating repository from directory", "err", err.Error(), "dir", options.RepositoryDirectory)
		return
	}

	if gittags, err = tags.All(repository); err != nil {
		lg.Error("error getting tags from repository", "err", err.Error())
		return
	}

	if semvers, err = semver.FromGitRefs(gittags); err != nil {
		lg.Error("error getting semvers from tags", "err", err.Error())
		return
	}

	// TODO: find #bump within commits
	// 		- get commits between main & current git pointer
	if defaultBranch, err = commits.FindReference(repository, options.DefaultBranch); err != nil {
		lg.Error("error getting git reference for default branch", "err", err.Error(), "default_branch", options.DefaultBranch)
		return
	}

	if refBranch, err = commits.FindReference(repository, options.BranchName); err != nil {
		lg.Error("error getting git reference for branch", "err", err.Error(), "branch", options.BranchName)
		return
	}

	if newCommits, err = commits.DiffBetween(repository, defaultBranch.Hash(), refBranch.Hash()); err != nil {
		lg.Error("error commits between references", "err", err.Error(), "base", defaultBranch.Hash().String(), "head", refBranch.Hash().String())
		return
	}

	debug(len(newCommits))

	if options.Prerelease {
		// run safe on the branch name for prerelease usage
		suffix, _ := strs.Safe(options.BranchName, options.PrereleaseSuffixLength)
		use = semver.Prerelease(lg, semvers, bump, suffix)
	} else {
		use = semver.Release(lg, semvers, bump)
	}

	// setup the prefix
	if options.WithPrefix {
		use.Prefix = "v"
	} else {
		use.Prefix = ""
	}

	// TODO: CREATE TAG
	if !options.TestMode {
		tag, err := tags.Create(repository, use.String(), refBranch.Hash())
	}

	debug(use.String())

	return
}

// Debug is a helper function that runs printf against a json
// string version of the item passed.
// Used for testing only.
func debug[T any](item T) {
	bytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Printf("%+v\n", string(bytes))
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
	flag.BoolVar(&runOptions.WithPrefix, "with-prefix", runOptions.WithPrefix, "Should we use a prefix. (default: true - will use v)")
	// test mode - disables creating tags
	flag.BoolVar(&runOptions.TestMode, "test", runOptions.TestMode, "Set to true to disable creating tag.")
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
