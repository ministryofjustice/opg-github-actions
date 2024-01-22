package main

import (
	"fmt"
	"opg-github-actions/cmd/branchname"
	"opg-github-actions/cmd/createtag"
	"opg-github-actions/cmd/latesttag"
	"opg-github-actions/cmd/nexttag"
	"opg-github-actions/pkg/semver"
	"opg-github-actions/pkg/testlib"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// fixture for handling all the end to end tests
type fixture struct {
	RepoSetup  repoSetup
	EventSetup eventSetup

	Prerelease    bool
	BranchTest    branchFixture
	LatestTagTest latestTagFixture
	NextTagTest   nextTagFixture
	CreateTagTest createTagFixture
}

type eventSetup struct {
	Event string
	Data  []byte
}

type repoSetup struct {
	Btc []branchCommitTags
}

type branchCommitTags struct {
	Branch  string
	Capture bool // when true, the commit hashes from this sequence will be tracked
	TagCom  []tagCommit
}
type tagCommit struct {
	Msg string
	Tag string
}

type branchFixture struct {
	Event    string
	Expected map[string]string
	Error    bool
}
type latestTagFixture struct {
	Expected map[string]string
	Error    bool
}

type nextTagFixture struct {
	Expected map[string]string
	Error    bool
}
type createTagFixture struct {
	Expected map[string]string
	Error    bool
}

var (
	// END TO END semver tests
	// create spec for a repository and events,
	// add expected results for each step (branch-name -> latest-tegs -> next-tag -> create-tag)
	semverScenarios = []fixture{
		// Testing a prerelease pull request from beta -> master with 3 commits, 1 major
		// starting with no previous releases or prereleases
		// => 1.0.0-myveryverylo.0
		{
			Prerelease: true,
			EventSetup: eventSetup{
				Event: "pull_request",
				Data:  testlib.TestEventPullRequest("master", "my-very-very-long-branch-name", "beta merge", "beta merge body #patch"),
			},
			// simple repo, 1 branch, 1 prerelease tag, 3 commits with 1 major
			RepoSetup: repoSetup{
				Btc: []branchCommitTags{
					{
						Branch: "my-very-very-long-branch-name",
						TagCom: []tagCommit{
							{Msg: "release"},
							{Msg: "commit test"},
							{Msg: "commit #major"},
						},
					},
				},
			},
			// branch name test setup
			BranchTest: branchFixture{
				Event:    "pull_request",
				Expected: map[string]string{"full_length": "myveryverylongbranchname", "safe": "myveryverylo", "branch_name": "my-very-very-long-branch-name"},
			},
			LatestTagTest: latestTagFixture{
				Expected: map[string]string{"last_release": "", "last_prerelease": ""},
			},
			NextTagTest: nextTagFixture{
				Expected: map[string]string{"next_tag": "1.0.0-myveryverylo.0"},
			},
			CreateTagTest: createTagFixture{
				Expected: map[string]string{"created_tag": "1.0.0-myveryverylo.0", "regenerated": "false"},
			},
		},
		// Prerelease from my/branch/name/has-slashes -> master with 1 minor
		// has a previous prerelease
		// with a prefix
		// => v2.1.0-mybranchname.1
		{
			Prerelease: true,
			EventSetup: eventSetup{
				Event: "pull_request",
				Data:  testlib.TestEventPullRequest("master", "my/branch/name/has-slashes", "merge", "merge body"),
			},
			// simple repo, 1 branch, 1 prerelease tag, 3 commits with 1 major
			RepoSetup: repoSetup{
				Btc: []branchCommitTags{
					{
						Branch: "my/branch/name/has-slashes",
						TagCom: []tagCommit{
							{Msg: "release", Tag: "v2.0.1"},
							{Msg: "commit test"},
							{Msg: "commit #minor", Tag: "v2.1.0-mybranchname.0"},
						},
					},
				},
			},
			// branch name test setup
			BranchTest: branchFixture{
				Event:    "pull_request",
				Expected: map[string]string{"full_length": "mybranchnamehasslashes", "safe": "mybranchname", "branch_name": "my/branch/name/has-slashes"},
			},
			LatestTagTest: latestTagFixture{
				Expected: map[string]string{"last_release": "v2.0.1", "last_prerelease": "v2.1.0-mybranchname.0"},
			},
			NextTagTest: nextTagFixture{
				Expected: map[string]string{"next_tag": "v2.1.0-mybranchname.1"},
			},
			CreateTagTest: createTagFixture{
				Expected: map[string]string{"created_tag": "v2.1.0-mybranchname.1", "regenerated": "false"},
			},
		},
		// Testing a prerelease pull request from beta -> master with 3 commits, 1 major 1 patch
		// starting 1.0.0 => 2.0.0-beta.0
		{
			Prerelease: true,
			EventSetup: eventSetup{
				Event: "pull_request",
				Data:  testlib.TestEventPullRequest("master", "beta", "beta merge", "beta merge body #patch"),
			},
			// simple repo, 1 branch, 1 prerelease tag, 3 commits with 1 major
			RepoSetup: repoSetup{
				Btc: []branchCommitTags{
					{
						Branch: "beta",
						TagCom: []tagCommit{
							{Msg: "release", Tag: "1.0.0"},
							{Msg: "commit test"},
							{Msg: "commit #major", Tag: "1.0.1-beta.0"},
						},
					},
				},
			},
			// branch name test setup
			BranchTest: branchFixture{
				Event:    "pull_request",
				Expected: map[string]string{"full_length": "beta", "safe": "beta", "branch_name": "beta"},
			},
			LatestTagTest: latestTagFixture{
				Expected: map[string]string{"last_release": "1.0.0", "last_prerelease": "1.0.1-beta.0"},
			},
			NextTagTest: nextTagFixture{
				Expected: map[string]string{"next_tag": "2.0.0-beta.0"},
			},
			CreateTagTest: createTagFixture{
				Expected: map[string]string{"created_tag": "2.0.0-beta.0", "regenerated": "false"},
			},
		},
		// Testing a release
		// No previous tags at all, but 1 #major in the commit history
		// => 1.0.0
		{
			Prerelease: false,
			EventSetup: eventSetup{
				Event: "push",
			},
			// simple repo, 1 branch, 1 prerelease tag, 3 commits with 1 major
			RepoSetup: repoSetup{
				Btc: []branchCommitTags{
					{
						Capture: true,
						Branch:  "master",
						TagCom: []tagCommit{
							{Msg: "release"},
							{Msg: "commit test"},
							{Msg: "commit #major"},
						},
					},
				},
			},
			// branch name test setup
			BranchTest: branchFixture{
				Event:    "push",
				Expected: map[string]string{"full_length": "master", "safe": "master", "branch_name": "master"},
			},
			LatestTagTest: latestTagFixture{
				Expected: map[string]string{"last_release": "", "last_prerelease": ""},
			},
			NextTagTest: nextTagFixture{
				Expected: map[string]string{"next_tag": "1.0.0"},
			},
			CreateTagTest: createTagFixture{
				Expected: map[string]string{"created_tag": "1.0.0", "regenerated": "false"},
			},
		},
		// Test a release
		// v1.9.0 and 1.10.0 releases, testing that the sort returns 1.10.0 as the latest
		// test that prerelease from a different branch is found, but not used
		// only a #minor
		// => v1.11.0
		{
			Prerelease: false,
			EventSetup: eventSetup{
				Event: "push",
			},
			RepoSetup: repoSetup{
				Btc: []branchCommitTags{
					{
						Branch: "alpha",
						TagCom: []tagCommit{
							{Msg: "commit", Tag: "v0.0.1-beta.0"},
							{Msg: "this commit should not be used #patch"},
						},
					},
					{
						Capture: true,
						Branch:  "master",
						TagCom: []tagCommit{
							{Msg: "release", Tag: "v1.10.0"},
							{Msg: "commit test", Tag: "v1.9.0"},
							{Msg: "commit #minor"},
						},
					},
					{
						Branch: "beta",
						TagCom: []tagCommit{
							{Msg: "commit", Tag: "v2.0.0-beta.0"},
							{Msg: "this commit should not be used either #major"},
						},
					},
				},
			},
			// branch name test setup
			BranchTest: branchFixture{
				Event:    "push",
				Expected: map[string]string{"full_length": "master", "safe": "master", "branch_name": "master"},
			},
			LatestTagTest: latestTagFixture{
				Expected: map[string]string{"last_release": "v1.10.0", "last_prerelease": "", "all_releases": "v1.9.0, v1.10.0", "all_prereleases": "v0.0.1-beta.0, v2.0.0-beta.0"},
			},
			NextTagTest: nextTagFixture{
				Expected: map[string]string{"next_tag": "v1.11.0"},
			},
			CreateTagTest: createTagFixture{
				Expected: map[string]string{"created_tag": "v1.11.0", "regenerated": "false"},
			},
		},
	}
)

// doEventSetup creates a dummy event file at a known location
// and a directory to store it in
func doEventSetup(es eventSetup) (dir string) {
	dir, _ = os.MkdirTemp("./", "test_event_")
	path := filepath.Join(dir, "./event.json")
	os.WriteFile(path, es.Data, 0777)
	return
}

// doRepoSetup creates a test repo, branches, commits and tags
// to setup a test scenario
func doRepoSetup(s repoSetup, r *git.Repository, defBranch *plumbing.Reference, dir string) (base plumbing.Hash, head plumbing.Hash) {
	w, _ := r.Worktree()
	b, _ := r.Head()
	base = b.Hash()

	for _, btc := range s.Btc {
		branchIter, _ := r.Branches()
		branches := []string{}
		branchIter.ForEach(func(ref *plumbing.Reference) error {
			branches = append(branches, ref.Name().Short())
			return nil
		})
		// see if the branch already exists
		if slices.Contains(branches, btc.Branch) {
			// checkout to the branch
			w.Checkout(&git.CheckoutOptions{Create: false, Force: true, Branch: plumbing.ReferenceName(btc.Branch)})
		} else {
			// create the branch
			testlib.TestRepositoryCreateBranch(r, btc.Branch)
		}
		// now create commits on the branch
		for _, tc := range btc.TagCom {
			commit, _ := testlib.TestRepositoryCommit(r, tc.Msg)
			if btc.Capture {
				head = commit
			}
			// create tag at this commit
			if len(tc.Tag) > 0 {
				rev := plumbing.Revision(commit.String())
				testlib.TestRepositoryCreateTag(r, tc.Tag, &rev)
			}

		}
		// check out to the default
		w.Checkout(&git.CheckoutOptions{
			Create: false,
			Force:  true,
			Branch: defBranch.Name(),
		})
	}
	// check out to the default
	w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  true,
		Branch: defBranch.Name(),
	})
	return
}

func TestSemverEndToEnd(t *testing.T) {
	testlib.Testlogger(nil)

	base, _ := os.Getwd()
	// loop over all the scenerios as test chaining of outputs -> inputs etc
	for i, f := range semverScenarios {
		// run repo setup
		repoDir, r, branch := testlib.TestRepositorySkeleton()
		baseH, headH := doRepoSetup(f.RepoSetup, r, branch, repoDir)
		// if this is a push event, have to setup data here, after the commits
		// so can use the hashes
		if f.EventSetup.Event == "push" {
			f.EventSetup.Data = testlib.TestEventPush("master", baseH.String(), headH.String())
		}
		// run event setup
		eventDir := doEventSetup(f.EventSetup)

		defer os.RemoveAll(eventDir)
		defer os.RemoveAll(repoDir)

		//----- BRANCH NAME
		branchResult, e := branchname.Run([]string{
			fmt.Sprintf(`--event-name=%s`, f.EventSetup.Event),
			fmt.Sprintf(`--event-data-file=%s/%s/%s`, base, eventDir, "event.json"),
		})

		if e != nil {
			t.Errorf("error: unexpected (%s:%d) error:", branchname.Name, i)
			t.Error(e)
		}
		for k, v := range f.BranchTest.Expected {
			if branchResult[k] != v {
				t.Errorf("error: (%s:%d) expected [%s] to be [%s] actual [%v]", branchname.Name, i, k, v, branchResult[k])
			}
		}

		//----- LATEST TAG
		latestTagResult, e := latesttag.Run([]string{
			fmt.Sprintf(`--repository=%s`, repoDir),
			fmt.Sprintf(`--branch=%s`, branchResult["branch_name"]),
			fmt.Sprintf(`--prerelease=%t`, f.Prerelease),
			fmt.Sprintf(`--prerelease-suffix=%s`, branchResult["safe"]),
		})

		if e != nil {
			t.Errorf("error: unexpected (%s:%d) error:", latesttag.Name, i)
			t.Error(e)
		}
		for k, v := range f.LatestTagTest.Expected {
			if latestTagResult[k] != v {
				t.Errorf("error: (%s:%d) expected [%s] to be [%s] actual [%v]", latesttag.Name, i, k, v, latestTagResult[k])
			}
		}

		//----- NEXT TAG
		nextTagResult, e := nexttag.Run([]string{
			fmt.Sprintf(`--repository=%s`, repoDir),
			fmt.Sprintf(`--base=%s`, branchResult["base_commitish"]),
			fmt.Sprintf(`--head=%s`, branchResult["head_commitish"]),
			fmt.Sprintf(`--prerelease=%s`, latestTagResult["prerelease"]),
			fmt.Sprintf(`--prerelease-suffix=%s`, latestTagResult["prerelease_suffix"]),
			fmt.Sprintf(`--last-release=%s`, latestTagResult["last_release"]),
			fmt.Sprintf(`--last-prerelease=%s`, latestTagResult["last_prerelease"]),
			fmt.Sprintf(`--with-v=%s`, latestTagResult["with_v"]),
			fmt.Sprintf(`--default-bump=%s`, string(semver.Patch)[1:]),
		})
		if e != nil {
			t.Errorf("error: unexpected (%s:%d) error:", nexttag.Name, i)
			t.Error(e)
		}

		for k, v := range f.NextTagTest.Expected {
			if nextTagResult[k] != v {
				t.Errorf("error: (%s:%d) expected [%s] to be [%s] actual [%v]", nexttag.Name, i, k, v, nextTagResult[k])
				fmt.Println(nextTagResult)
			}
		}

		//----- CREATE TAG
		createTagResult, e := createtag.Run([]string{
			fmt.Sprintf(`--repository=%s`, repoDir),
			fmt.Sprintf(`--commitish=%s`, branchResult["branch_name"]),
			fmt.Sprintf(`--tag-name=%s`, nextTagResult["next_tag"]),
			fmt.Sprintf(`--regen=%t`, true), // force true so it will always try
			fmt.Sprintf(`--push=%t`, false), // force false so we dont try to push to non-existant remote
		})
		if e != nil {
			t.Errorf("error: unexpected (%s:%d) error:", createtag.Name, i)
			t.Error(e)
		}

		for k, v := range f.CreateTagTest.Expected {
			if createTagResult[k] != v {
				t.Errorf("error: (%s:%d) expected [%s] to be [%s] actual [%v]", createtag.Name, i, k, v, createTagResult[k])
				fmt.Println(nextTagResult)
				fmt.Println(createTagResult)
			}
		}

	}

}
