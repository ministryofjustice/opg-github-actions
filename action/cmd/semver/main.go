package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"opg-github-actions/action/internal/logger"
	"opg-github-actions/action/internal/repo"
	"opg-github-actions/action/internal/semver"
	"opg-github-actions/action/internal/tags"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const ErrNoBranchName string = "branch-name is required, but not found."

type Options struct {
	RepositoryDirectory    string
	Prerelease             bool
	PrereleaseSuffixLength int
	BranchName             string
}

var runOptions *Options = &Options{
	RepositoryDirectory:    "",
	Prerelease:             false,
	PrereleaseSuffixLength: 12,
	BranchName:             "",
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

	// prerelease related options
	flag.StringVar(&runOptions.BranchName, "branch", runOptions.BranchName, "The current branch name to use for prerelease suffixes")
	flag.BoolVar(&runOptions.Prerelease, "prerelease", runOptions.Prerelease, "Set to true to generate a prerelease version.")
	flag.IntVar(&runOptions.PrereleaseSuffixLength, "prerelease-suffix-length", runOptions.PrereleaseSuffixLength, "Set the max length to use for tag suffixes")
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

func Run(lg *slog.Logger, options *Options) (result map[string]string, err error) {
	var (
		repository *git.Repository
		gittags    []*plumbing.Reference
		semvers    []*semver.Semver
		use        *semver.Semver
	)
	result = map[string]string{}

	if options.Prerelease && options.BranchName == "" {
		err = fmt.Errorf(ErrNoBranchName)
		return
	}

	if repository, err = repo.FromDir(options.RepositoryDirectory); err != nil {
		lg.Error("error creating repository from directory", "err", err.Error())
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

	if options.Prerelease {
		use = semver.Prerelease(lg, semvers, semver.PATCH, options.BranchName)
	} else {
		use = semver.Release(lg, semvers, semver.PATCH)
	}

	debug(use.String())

	return
}
