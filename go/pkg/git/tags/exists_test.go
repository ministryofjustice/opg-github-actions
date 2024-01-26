package tags

import (
	"opg-github-actions/pkg/testlib"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

func TestTagExists(t *testing.T) {
	testlib.Testlogger(nil)
	directory, r, branch := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(directory)

	tags := []string{"tag-ext-01", "tag-ext", "1.0.0-beta.0+bA1"}
	for _, tag := range tags {
		rev := plumbing.Revision(branch.Hash().String())
		testlib.TestRepositoryCreateTag(r, tag, &rev)
	}
	tagset, _ := New(directory)
	// valid tests
	for _, tag := range tags {
		if !tagset.Exists(tag) {
			t.Errorf("error: tag [%s] should exist", tag)
		}
	}

	tag := "this-does-not-exist"
	if tagset.Exists(tag) {
		t.Errorf("error: tag [%s] should not exist", tag)
	}

}
