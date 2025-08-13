package main

import (
	"fmt"
	"math/rand"
	"opg-github-actions/action/internal/logger"
	"opg-github-actions/action/internal/semver"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type tSemTestCommit struct {
	Message      string
	Branch       string
	Tag          string
	ChildCommits []string
}

type tSemTest struct {
	Commits        []*tSemTestCommit
	ExpectedTag    string
	ExpectedBump   string
	ExpectedBranch string
	Input          *Options
	CreateRelease  bool
	ShouldError    bool
}

// Tests that align to what happens on push - so commits on the default branch
// after a release that should trigger a semver bump
func TestMainPush(t *testing.T) {
	var lg = logger.New("ERROR", "TEXT")
	var tests = []*tSemTest{
		// make a commit with #minor commit on the top of master (default branch)
		{
			ExpectedTag:    "v1.1.0",
			ExpectedBump:   string(semver.MINOR),
			ExpectedBranch: "master",
			ShouldError:    false,
			CreateRelease:  true,
			Input: &Options{
				Prerelease:    false,
				DefaultBranch: "master",
				BranchName:    "master",
			},
			Commits: []*tSemTestCommit{
				{Message: "just one #minor", Branch: "master"},
			},
		},
		// make commits with #major in there on the top of master (default branch)
		{
			ExpectedTag:   "v2.0.0",
			ExpectedBump:  string(semver.MAJOR),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease:    false,
				DefaultBranch: "master",
				BranchName:    "master",
			},
			Commits: []*tSemTestCommit{
				{Message: "foobar", Branch: "master"},
				{Message: "just one #minor", Branch: "master"},
				{Message: "this ones #major", Branch: "master"},
				{Message: "so is this one #major", Branch: "master"},
			},
		},
	}

	// var dir = "./test-repo"
	// os.RemoveAll(dir)
	// os.MkdirAll(dir, os.ModePerm)
	// r, defBranch := randomRepository(dir, true)
	// w, _ := r.Worktree()

	for i, test := range tests {
		var (
			dir          = t.TempDir()
			r, defBranch = randomRepository(dir, test.CreateRelease)
			w, _         = r.Worktree()
		)

		err := testSetup(test, r, w, defBranch)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// setup the options to run
		opts := newRunOptions(test.Input)
		if opts.RepositoryDirectory == "" {
			opts.RepositoryDirectory = dir
		}
		if opts.DefaultBranch == "" {
			opts.DefaultBranch = defBranch.Name().Short()
		}
		// now run the command and compare
		res, err := Run(lg, opts)
		// check error states
		if !test.ShouldError && err != nil {
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		} else if test.ShouldError && err == nil {
			t.Errorf("[%d] expected an error, but did not get one", i)
		}
		// check values
		if res["tag"] != test.ExpectedTag {
			t.Errorf("[%d] expected tag [%s] actual [%s]", i, test.ExpectedTag, res["tag"])
			debug(res)
		}
		if res["bump"] != test.ExpectedBump {
			t.Errorf("[%d] expected bump [%s] actual [%s]", i, test.ExpectedBump, res["bump"])
			debug(res)
		}

		if test.ExpectedBranch != "" && res["branch"] != test.ExpectedBranch {
			t.Errorf("[%d] expected branch [%s] actual [%s]", i, test.ExpectedBranch, res["branch"])
			debug(res)
		}
		// debug(res)

	}

	// t.FailNow()

}

// Test generating various semvers in setup that looks like pull requests
func TestMainPR(t *testing.T) {

	var lg = logger.New("ERROR", "TEXT")
	var tests = []*tSemTest{
		// test a prerelease tag that clashes with a similar branch and
		// tag that triggers a patch
		{
			ExpectedTag:   "v1.0.1-renovatefeatur.2",
			ExpectedBump:  string(semver.PATCH),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease: true,
				BranchName: "renovate-feature-b",
			},
			Commits: []*tSemTestCommit{
				{
					Message: "different, but similar commit that should trigger patch only",
					Branch:  "renovate-feature-a",
					Tag:     "v1.0.1-renovatefeatur.1",
				},
				{
					Message: "single commit thats not even a patch but using defaults",
					Branch:  "renovate-feature-b",
				},
			},
		},
		// trying to get the last semver, but there is a commit on a the branch
		// being merged, so it should be incremented by patch (form commit)
		{
			ExpectedTag:   "v1.1.0",
			ExpectedBump:  string(semver.MINOR),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease:    false,
				DefaultBranch: "master",
				BranchName:    "mybranch-testA",
				DefaultBump:   string(semver.NO_BUMP),
			},
			Commits: []*tSemTestCommit{
				{Message: "just one #minor", Branch: "mybranch-testA"},
			},
		},
		// test requesting no bump, but there being a commit between the last tag
		// and head of branch, so a patch should be returned
		//  - some bits here where we are comparing the tag v1.0.0 to master
		{
			ExpectedTag:   "v1.0.1",
			ExpectedBump:  string(semver.PATCH),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease:    false,
				DefaultBranch: "v1.0.0",
				BranchName:    "master",
				DefaultBump:   string(semver.NO_BUMP),
			},
			Commits: []*tSemTestCommit{
				{Message: "just one commit"},
			},
		},
		// test a prerelease tag that clashes with a similar branch and
		// tag that triggers a minor
		{
			ExpectedTag:   "v1.1.0-renovatefeatur.2",
			ExpectedBump:  string(semver.MINOR),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease: true,
				BranchName: "renovate-feature-b",
			},
			Commits: []*tSemTestCommit{
				{
					Message: "different, but similar commit thats #minor",
					Branch:  "renovate-feature-a",
					Tag:     "v1.1.0-renovatefeatur.1",
				},
				{
					Message: "single commit thats #minor",
					Branch:  "renovate-feature-b",
				},
			},
		},
		// test a prerelease tag that clashes with a similar branch and
		// tag that triggers a major
		{
			ExpectedTag:   "v2.0.0-renovatefeatur.2",
			ExpectedBump:  string(semver.MAJOR),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease:  true,
				BranchName:  "renovate-feature-b",
				DefaultBump: string(semver.PATCH),
			},
			Commits: []*tSemTestCommit{
				{
					Message: "different, but similar commit thats #major",
					Branch:  "renovate-feature-a",
					Tag:     "v2.0.0-renovatefeatur.1",
				},
				{
					Message: "single commit thats #major",
					Branch:  "renovate-feature-b",
				},
			},
		},
		// test a release tag when there are no prior releases
		// based on commit message and disable prefix
		{
			ExpectedTag:   "1.0.0",
			ExpectedBump:  string(semver.MAJOR),
			ShouldError:   false,
			CreateRelease: false,
			Input: &Options{
				Prerelease:    false,
				BranchName:    "foobar",
				DefaultBump:   string(semver.PATCH),
				WithoutPrefix: true,
			},
			Commits: []*tSemTestCommit{
				{
					Message: "single commit thats #major",
					Branch:  "foobar",
				},
			},
		},
		// test a release tag when there are no prior releases
		// based on commit message
		{
			ExpectedTag:   "v1.0.0",
			ExpectedBump:  string(semver.MAJOR),
			ShouldError:   false,
			CreateRelease: false,
			Input: &Options{
				Prerelease:  false,
				BranchName:  "foobar",
				DefaultBump: string(semver.PATCH),
			},
			Commits: []*tSemTestCommit{
				{
					Message: "single commit thats #major",
					Branch:  "foobar",
				},
			},
		},
		// test a patch prerelease when there are no prior releases
		// and its based on the default bump only
		{
			ExpectedTag:   "v0.0.1-foobar.1",
			ExpectedBump:  string(semver.PATCH),
			ShouldError:   false,
			CreateRelease: false,
			Input: &Options{
				Prerelease:  true,
				BranchName:  "foobar",
				DefaultBump: string(semver.PATCH),
			},
			Commits: []*tSemTestCommit{
				{
					Message: "single commit with no actual tag, but some close things like minor and major",
					Branch:  "foobar",
				},
			},
		},
		// test a release tag when there are no prior releases
		{
			ExpectedTag:   "v0.0.1",
			ExpectedBump:  string(semver.PATCH),
			ShouldError:   false,
			CreateRelease: false,
			Input: &Options{
				Prerelease:  false,
				BranchName:  "foobar",
				DefaultBump: string(semver.PATCH),
			},
			Commits: []*tSemTestCommit{
				{
					Message: "single commit with no actual tag, but some close things like minor and major",
					Branch:  "foobar",
				},
			},
		},
		// test a patch release version increment thats based on the
		// default bump
		{
			ExpectedTag:   "v1.0.1",
			ExpectedBump:  string(semver.PATCH),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease:  false,
				BranchName:  "test-branch-b",
				DefaultBump: string(semver.PATCH),
			},
			Commits: []*tSemTestCommit{
				{
					Message: "single commit with no actual tag, but some close things like minor and major",
					Branch:  "test-branch-b",
				},
			},
		},
		// test returning the last semver release tag as there
		// are no new commits and default bump is set to 'none'
		{
			ExpectedTag:   "v1.0.0",
			ExpectedBump:  string(semver.NO_BUMP),
			ShouldError:   false,
			CreateRelease: true,
			Input: &Options{
				Prerelease:  false,
				BranchName:  "master",
				DefaultBump: string(semver.NO_BUMP),
			},
			Commits: []*tSemTestCommit{},
		},
		// test a series of commits on a branch that should generate
		// a prerelease tag
		{
			ExpectedTag:    "v1.1.0-testbrancha.1",
			ExpectedBump:   string(semver.MINOR),
			ExpectedBranch: "testbrancha",
			ShouldError:    false,
			CreateRelease:  true,
			Input: &Options{
				Prerelease: true,
				BranchName: "test-branch-a",
			},
			Commits: []*tSemTestCommit{
				{
					Message: "my test commit without anything",
					Branch:  "test-branch-a",
					ChildCommits: []string{
						"this commit is not really a patch or minor and really not major",
						"clearly a change #patch",
						"a bigger change #minor",
					},
				},
			},
		},
		// test a series commits with multi lines and special chars
		// that should create a minor
		{
			ExpectedTag:    "v2.0.0-testbranchutf.1",
			ExpectedBump:   string(semver.MAJOR),
			ExpectedBranch: "testbranchutf",
			ShouldError:    false,
			CreateRelease:  true,
			Input: &Options{
				PrereleaseSuffixLength: 15,
				Prerelease:             true,
				BranchName:             "test-branch-utf",
			},
			Commits: []*tSemTestCommit{
				{
					Message: "my test commit without anything",
					Branch:  "test-branch-utf",
					ChildCommits: []string{
						`?? % & # >< : \ = - + ♥

				end here`,
						"this commit is not really a patch or minor and really not major",
						"a little change",
						"a bigger change #minor",
						"a massive change #major",
					},
				},
			},
		},
	}

	// var dir = "./test-repo"
	// os.RemoveAll(dir)
	// os.MkdirAll(dir, os.ModePerm)
	// r, defBranch := randomRepository(dir, true)
	// w, _ := r.Worktree()

	for i, test := range tests {
		var (
			dir          = t.TempDir()
			r, defBranch = randomRepository(dir, test.CreateRelease)
			w, _         = r.Worktree()
		)

		err := testSetup(test, r, w, defBranch)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		// setup the options to run
		opts := newRunOptions(test.Input)
		if opts.RepositoryDirectory == "" {
			opts.RepositoryDirectory = dir
		}
		if opts.DefaultBranch == "" {
			opts.DefaultBranch = defBranch.Name().Short()
		}
		// now run the command and compare
		res, err := Run(lg, opts)
		// check error states
		if !test.ShouldError && err != nil {
			t.Errorf("[%d] unexpected error: %s", i, err.Error())
		} else if test.ShouldError && err == nil {
			t.Errorf("[%d] expected an error, but did not get one", i)
		}
		// check values
		if res["tag"] != test.ExpectedTag {
			t.Errorf("[%d] expected tag [%s] actual [%s]", i, test.ExpectedTag, res["tag"])
			debug(res)
		}
		if res["bump"] != test.ExpectedBump {
			t.Errorf("[%d] expected bump [%s] actual [%s]", i, test.ExpectedBump, res["bump"])
			debug(res)
		}
		if test.ExpectedBranch != "" && res["branch"] != test.ExpectedBranch {
			t.Errorf("[%d] expected branch [%s] actual [%s]", i, test.ExpectedBranch, res["branch"])
			debug(res)
		}

		// debug(res)

	}

	// t.FailNow()

}

// testSetup generates some base commits / branches / tags to use in test scenarios
func testSetup(test *tSemTest, r *git.Repository, w *git.Worktree, defBranch *plumbing.Reference) (err error) {
	var author = &object.Signature{Name: "go test", Email: "test@example.com"}
	// now create the test commits
	for _, commit := range test.Commits {
		var err error
		var hash plumbing.Hash
		var createBranch = true
		var branch = defBranch.Name()
		// if theres a branch name, use that instead of default
		if commit.Branch != "" {
			branch = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", commit.Branch))
		}
		// if the branch being asked for is the default branch, then dont create it
		if branch == defBranch.Name() {
			createBranch = false
		}
		// checkout to the branch we want
		err = w.Checkout(&git.CheckoutOptions{Create: createBranch, Force: true, Branch: branch})
		if err != nil {
			return fmt.Errorf("checkout unexpected error [%s]: %s", branch, err.Error())
		}
		// create the commit
		hash, err = w.Commit(commit.Message, &git.CommitOptions{AllowEmptyCommits: true, Author: author})
		if err != nil {
			return fmt.Errorf("commit unexpected error: %s", err.Error())
		}
		// create a tag at this commit if configured
		if commit.Tag != "" {
			rev := plumbing.Revision(hash.String())
			sha, _ := r.ResolveRevision(rev)
			r.CreateTag(commit.Tag, *sha, nil)
		}

		// now create any child commits on this branch
		for _, child := range commit.ChildCommits {
			_, err = w.Commit(child, &git.CommitOptions{AllowEmptyCommits: true, Author: author})
			if err != nil {
				return fmt.Errorf("commit unexpected error: %s", err.Error())
			}
		}
	}
	return
}

// randomRepository makes a repo with a mix of and a v1 release is asked
func randomRepository(dir string, createRelease bool) (r *git.Repository, defaultBranch *plumbing.Reference) {
	var (
		hash     plumbing.Hash
		commitsN = rand.Intn(100) + 30 // somewhere between 30-100 commits
		hashes   = []plumbing.Hash{}
	)
	// create the repository locally
	r, _ = git.PlainInit(dir, false)
	w, _ := r.Worktree()

	// create some commits on the base
	for i := 0; i < commitsN; i++ {
		var e error
		msg := fmt.Sprintf("commit %d", i)
		hash, e = w.Commit(msg, &git.CommitOptions{
			AllowEmptyCommits: true,
			Author:            &object.Signature{Name: "go test", Email: "test@example.com"},
		})
		if e == nil {
			hashes = append(hashes, hash)
		}
	}

	if createRelease {
		rev := plumbing.Revision(hash.String())
		sha, _ := r.ResolveRevision(rev)
		r.CreateTag("v1.0.0", *sha, nil)

	}

	defaultBranch, _ = r.Head()
	// checkout to default branch
	w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  true,
		Branch: defaultBranch.Name(),
	})

	return
}
