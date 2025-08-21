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
	"github.com/google/go-github/v74/github"
)

const ErrNoBranchName string = "branch-name is required, but not found."

type Options struct {
	RepositoryDirectory    string // Directory where the dit repo is
	Prerelease             bool   // if this is a prerelease of a full release
	PrereleaseSuffixLength int    // length of the prerelease suffix
	DefaultBranch          string // default branch name - generally main, used to compare commits against
	BranchName             string // branch name is used as the prerelease suffix
	DefaultBump            string // what to increment the semver by (major, minor, patch)
	EventContentFile       string // content from pull request title / body where their might be extra #major content
	WithoutPrefix          bool
	TestMode               bool
}

func (self *Options) SafeSuffix() (safeAndShort string) {
	safeAndShort, _ = strs.Safe(self.BranchName, self.PrereleaseSuffixLength)
	return
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
		EventContentFile:       "",
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
		if in.EventContentFile != "" {
			opts.EventContentFile = in.EventContentFile
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
	if gittags, err = tags.All(lg, repository); err != nil {
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
		use = semver.Prerelease(lg, semvers, bump, options.SafeSuffix())
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
		remotes              []*git.Remote
		errFailedToCreateTag string = "error: failed to create tag [%s]"
		errFailedToPush      string = "error: failed to push tags to remote [tag: %s]"
		tagName              string = use.String()
		auth                        = &http.BasicAuth{
			Username: "opg-github-actions",
			Password: token,
		}
	)
	lg = lg.With("operation", "createAndPushTag", "semver", use.String())

	// we do nothing if this is in test mode or we're using no bumping
	if options.TestMode || bump == semver.NO_BUMP {
		lg.Debug("returning, test mode / no increment enabled", "test", options.TestMode, "bump", string(bump))
		return
	}
	// fetch the remotes of the repo
	remotes, err = repository.Remotes()
	if err != nil {
		lg.Error("error getting remotes on repository")
		return
	}

	lg.Debug("creating tag ... ")
	// try to create the tag locally
	createdTag, err = tags.Create(repository, tagName, use.GitRef.Hash())
	if err != nil {
		tagErr := fmt.Errorf(errFailedToCreateTag, tagName)
		err = errors.Join(tagErr, err)
		lg.Error("failed to create tag ... ", "err", err.Error())
		return
	}

	// if we have some remotes, push
	if len(remotes) > 0 {
		lg.Debug("pushing tags ... ")
		err = tags.Push(repository, auth)
		if err != nil {
			tagErr := fmt.Errorf(errFailedToPush, tagName)
			err = errors.Join(tagErr, err)
			lg.Error("failed to push tag ... ", "err", err.Error())
			return
		}
	}

	return
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

// getContentFromEventFile reads and parses the event file
// that might be present at this path
// Doing it this way as the event content contains special
// characters that dont escape very well
func getContentFromEventFile(lg *slog.Logger, file string) (content string) {
	var (
		err       error
		bytes     []byte
		pr        *github.PullRequest      = nil
		prEvent   *github.PullRequestEvent = &github.PullRequestEvent{}
		pushEvent *github.PushEvent        = &github.PushEvent{}
		raw       map[string]interface{}   = map[string]interface{}{}
	)
	// ignore empty of missing files
	if file == "" || !fileExists(file) {
		return
	}

	content = ""
	if bytes, err = os.ReadFile(file); err != nil {
		lg.Error("err with reading file", "err", err.Error())
		return
	}

	// first, unmarshal into a map to test
	err = json.Unmarshal(bytes, &raw)
	if err != nil {
		lg.Error("err with unmarshal", "err", err.Error())
		return
	}
	// handle push event
	if _, ok := raw["commits"]; ok {
		json.Unmarshal(bytes, &pushEvent)
		for _, c := range pushEvent.Commits {
			content += *c.Message
		}
	}
	// look if its a pull request, and if so, parse as a
	// pr and generate the content from that
	if _, ok := raw["pull_request"]; ok {
		err = json.Unmarshal(bytes, &prEvent)
		if err != nil {
			lg.Error("error unmarshaling pull request event", "err", err.Error())
			return
		}
		// If the pr body or title are empty, then their pointer is nil, so add
		// handling for that
		if prEvent != nil && prEvent.PullRequest != nil {
			pr = prEvent.PullRequest
		}
		if pr != nil && pr.Title != nil {
			content += fmt.Sprintf("%s\n", *pr.Title)
		}
		if pr != nil && pr.Body != nil {
			content += fmt.Sprintf(" %s", *pr.Body)
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
		lastRelease   *semver.Semver                                              // the last release or default branch
		baseRef       *plumbing.Reference                                         // git ref for previous ref (main / or last release)
		currentCommit *plumbing.Reference                                         // git sha / ref for where the git repo currently is
		createdTag    *plumbing.Reference                                         // the new semver tag thats been created
		newCommits    []*object.Commit                                            // all commits that exist in the ref
		bump          semver.Increment    = semver.Increment(options.DefaultBump) // default increment
		basePoint     string              = ""                                    // either ref of last release or the default branch
		token         string              = os.Getenv("GH_TOKEN")                 // github auth token for pushing to the remote
		bumpCommit    string              = ""                                    // commit message the cump wasa found within
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

	// output all semvers for debuggin
	lg.Debug("found semvers ... ", "semvers", semvers)
	// get the last release for the comparison point
	lastRelease = semver.GetLastRelease(lg, semvers)
	lg.Debug("last release found ... ", "lastRelease", lastRelease)

	if lastRelease != nil {
		basePoint = lastRelease.Stringy(true)
	} else {
		basePoint = options.DefaultBranch
	}

	// get the default branch info
	if baseRef, err = commits.FindReference(lg, repository, basePoint); err != nil {
		lg.Error("error getting git reference for previous comparison point.", "err", err.Error(), "basePoint", basePoint)
		return
	}
	// get info on the current commit
	if currentCommit, err = commits.FindReference(lg, repository, options.BranchName); err != nil {
		lg.Error("error getting git reference for branch.", "err", err.Error(), "branch", options.BranchName)
		return
	}
	// find new commits between the baseRef (main / last release) and the current commit
	if newCommits, err = commits.DiffBetween(lg, repository, baseRef.Hash(), currentCommit.Hash()); err != nil {
		lg.Error("error getting commits between references.", "err", err.Error(), "base", baseRef.Hash().String(), "head", currentCommit.Hash().String())
		return
	}

	// add content to the commit list from the event file
	if extra := getContentFromEventFile(lg, options.EventContentFile); len(extra) > 0 {
		newCommits = append(newCommits, &object.Commit{Hash: plumbing.ZeroHash, Message: extra})
	}

	lg.Info("commits found and event file used ... ", "commits", len(newCommits), "event-file", options.EventContentFile)

	// dump the commits for debugging one at a time (messages can be long)
	for _, c := range newCommits {
		lg.Debug("commit", "message", c.Message, "hash", c.Hash)
	}

	// look for bump in the commits,
	foundBump, bumpCommit := semver.GetBumpFromCommits(lg, newCommits, bump)
	if len(newCommits) > 0 && foundBump != "" {
		bump = foundBump
	}

	use = getSemverToUse(lg, semvers, bump, options)
	lg.Info("got semver ... ", "use", use)
	// set the git ref to the current place
	use.GitRef = currentCommit
	// create and try to push tags
	createdTag, err = createAndPushTag(lg, repository, use, bump, token, options)

	result = map[string]string{
		"tag":     use.String(),
		"hash":    use.GitRef.Hash().String(),
		"branch":  options.SafeSuffix(),
		"test":    fmt.Sprintf("%t", options.TestMode),
		"created": fmt.Sprintf("%t", (createdTag != nil)),
		"bump":    string(bump),
		"commit":  bumpCommit,
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
	flag.StringVar(&runOptions.EventContentFile, "event-content-file", runOptions.EventContentFile, "The github event file that contains extra content")
}

func main() {
	var lg *slog.Logger = logger.New("info", "text")
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
