package commits

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestCommitsDiffBetween(t *testing.T) {

	var dir = t.TempDir()
	repo, defBranch := randomRepository(dir)

	// get the default branch
	base, err := FindReference(repo, defBranch.Name().String())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	// get a test branch
	other, err := FindReference(repo, "my-branch-1")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	commits, err := DiffBetween(repo, base.Hash(), other.Hash())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(commits) == 0 {
		t.Errorf("should have found commits between different reference points")
	}

	// compare same commit points
	commits, err = DiffBetween(repo, base.Hash(), base.Hash())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(commits) != 0 {
		t.Errorf("should not have found commits between same ref points")
	}

	// compare same commit points
	commits, err = DiffBetween(repo, other.Hash(), other.Hash())
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if len(commits) != 0 {
		t.Errorf("should not have found commits between same ref points")
	}

}

func TestCommitsFindReference(t *testing.T) {

	var dir = t.TempDir()
	repo, _ := randomRepository(dir)

	// this tag should always exist on the main branch
	lookFor := "v1.0.1-thing"
	found, err := FindReference(repo, lookFor)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if found.Name().Short() != lookFor {
		t.Errorf("error finding the reference, expected [%s], got [%s]", lookFor, found.Name().Short())
	}

	// look for a set tag on the end of a branch
	lookFor = "v1.1.1-pre"
	found, err = FindReference(repo, lookFor)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if found.Name().Short() != lookFor {
		t.Errorf("error finding the reference, expected [%s], got [%s]", lookFor, found.Name().Short())
	}

	// look for a known branch name
	lookFor = "my-branch-1"
	found, err = FindReference(repo, lookFor)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if found.Name().Short() != lookFor {
		t.Errorf("error finding the reference, expected [%s], got [%s]", lookFor, found.Name().Short())
	}

}

// randomRepository makes a repo with a mix of commits at that base
// then some tags at various points on the base branch, then some branches,
// each with a few commits and a tag
func randomRepository(dir string) (r *git.Repository, defaultBranch *plumbing.Reference) {
	var (
		commitsN = rand.Intn(50) + 10 // somewhere between 10-50 commits
		preTags  = rand.Intn(15) + 5  // 5-15 tags
		relTags  = rand.Intn(10) + 2  // 2-10 tags
		branches = rand.Intn(2) + 5   // 2-5 branches
		hashes   = []plumbing.Hash{}
	)
	// create the repository locally
	r, _ = git.PlainInit(dir, false)
	w, _ := r.Worktree()

	// create some commits on the base
	for i := 0; i < commitsN; i++ {
		msg := fmt.Sprintf("commit %d", i)
		hash, e := w.Commit(msg, &git.CommitOptions{
			AllowEmptyCommits: true,
			Author:            &object.Signature{Name: "go test", Email: "test@example.com"},
		})
		if e == nil {
			hashes = append(hashes, hash)
		}
	}
	// make some random tags with v1.0.x-thing format
	for i := 0; i < preTags; i++ {
		var randI = rand.Intn(len(hashes))
		var hash = hashes[randI]
		var tag = "v1.0." + strconv.Itoa(i+1) + "-thing"
		var rev = plumbing.Revision(hash.String())

		revHash, _ := r.ResolveRevision(rev)
		r.CreateTag(tag, *revHash, nil)
	}

	// make some random tags with v1.0.x-thing format
	for i := 0; i < relTags; i++ {
		var randI = rand.Intn(len(hashes))
		var hash = hashes[randI]
		var tag = "v1.0." + strconv.Itoa(i+1)
		var rev = plumbing.Revision(hash.String())

		revHash, _ := r.ResolveRevision(rev)
		r.CreateTag(tag, *revHash, nil)
	}

	defaultBranch, _ = r.Head()
	// checkout to default branch
	w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  true,
		Branch: defaultBranch.Name(),
	})

	// now make a few branches and generate some commits on each branch, then create a tag on the end of that branch
	for i := 0; i < branches; i++ {
		var branch = fmt.Sprintf("refs/heads/my-branch-%d", i+1)
		var commitCount = rand.Intn(5) + 1 // 1-5 commits
		var tag = "v1.1." + strconv.Itoa(i+1) + "-pre"
		var commit plumbing.Hash

		branchRef := plumbing.ReferenceName(branch)
		w.Checkout(&git.CheckoutOptions{
			Create: true,
			Force:  true,
			Branch: branchRef,
		})

		for x := 0; x < commitCount; x++ {
			commit, _ = w.Commit(fmt.Sprintf("commit %d.%d #minor", i, x), &git.CommitOptions{
				AllowEmptyCommits: true,
				Author:            &object.Signature{Name: "go test", Email: "test@example.com"},
			})
		}
		r.CreateTag(tag, commit, nil)
		// checkout to default branch
		w.Checkout(&git.CheckoutOptions{
			Create: false,
			Force:  true,
			Branch: defaultBranch.Name(),
		})

	}

	// checkout to default branch
	w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  true,
		Branch: defaultBranch.Name(),
	})

	return
}
