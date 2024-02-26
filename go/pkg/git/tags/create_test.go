package tags

import (
	"opg-github-actions/pkg/testlib"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

func TestCreateTag(t *testing.T) {
	testlib.Testlogger(nil)
	directory, r, branch := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(directory)

	tags := []string{"tag-ext-01", "1.0.0-beta.0+bA1"}
	for _, tag := range tags {
		rev := plumbing.Revision(branch.Hash().String())
		testlib.TestRepositoryCreateTag(r, tag, &rev)
	}

	hash, _ := testlib.TestRepositoryCommit(r, "test commit")
	tagset, _ := New(directory)

	createTag := "test-tag"
	create, _ := tagset.CreateAt(createTag, &hash)

	// should be 3 tags in total
	all, _ := tagset.All()
	if len(all) != len(tags)+1 {
		t.Errorf("error: should have more tags than found")
	}
	// should be one tag at this location
	h := create.Hash().String()
	here := plumbing.NewHash(h)
	at, _ := tagset.At(&here)

	if at[0].Hash().String() != h {
		t.Errorf("error: commit hash and tag hash do not match")
	}

	if at[0].Name().Short() != createTag {
		t.Errorf("error: tag name created does not match expected.")
	}
}
