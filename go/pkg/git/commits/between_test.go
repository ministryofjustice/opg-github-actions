package commits

import (
	"fmt"
	"opg-github-actions/pkg/testlib"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

// TestDiffBetween creates a test repo, a new branch
// and series of commits on the new branch and then
// finds the difference between the two, checking they
// match
func TestDiffBetween(t *testing.T) {
	testlib.Testlogger(nil)
	directory, r, branch := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(directory)

	baseHash := branch.Hash()
	commitSet, _ := New(directory)

	newBranch := "test-commit-diffs"
	testlib.TestRepositoryCreateBranch(r, newBranch)
	l := 10
	for i := 0; i < l; i++ {
		testlib.TestRepositoryCommit(r, fmt.Sprintf("commit %d", i))
	}

	headRef := plumbing.Revision(newBranch)
	headHash, _ := r.ResolveRevision(headRef)
	head := plumbing.NewReferenceFromStrings(headHash.String(), headHash.String())

	diffs, _ := commitSet.DiffBetween(baseHash, head.Hash())

	if len(diffs) != l {
		t.Errorf("error: did not find correct amount of diffs")
	}

}
