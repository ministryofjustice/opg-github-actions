package testlib

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func TestRepositorySkeleton() (dir string, r *git.Repository, defaultBranch *plumbing.Reference) {
	dir, _ = os.MkdirTemp("./", "test_repository_")
	// create the repository locally
	r, _ = git.PlainInit(dir, false)
	w, _ := r.Worktree()

	// create some commits on the base
	n := TestRandInRange(10, 150)
	for i := 0; i < n; i++ {
		msg := fmt.Sprintf("commit %d", i)
		w.Commit(msg, &git.CommitOptions{
			AllowEmptyCommits: true,
		})
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

func TestRepositoryCreateBranch(r *git.Repository, newBranchName string) (plumbing.ReferenceName, error) {
	w, _ := r.Worktree()
	branch := fmt.Sprintf("refs/heads/%s", newBranchName)
	branchRef := plumbing.ReferenceName(branch)
	return branchRef, w.Checkout(&git.CheckoutOptions{
		Create: true,
		Force:  true,
		Branch: branchRef,
	})
}

func TestRepositoryCommit(r *git.Repository, commitMsg string) (plumbing.Hash, error) {
	w, _ := r.Worktree()

	return w.Commit(commitMsg, &git.CommitOptions{
		AllowEmptyCommits: true,
	})
}

func TestRepositoryCreateTag(r *git.Repository, tagName string, ref *plumbing.Revision) (*plumbing.Reference, error) {
	rev, _ := r.ResolveRevision(*ref)
	return r.CreateTag(tagName, *rev, nil)
}
