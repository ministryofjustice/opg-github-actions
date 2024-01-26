package tags

import (
	"fmt"
	"opg-github-actions/pkg/testlib"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

func TestTagsSortWithCreate(t *testing.T) {
	testlib.Testlogger(nil)
	directory, r, branch := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(directory)
	// create some tags in the repository and then fetch to make sure they are found
	tags := []string{"9.5.0", "10.1.0", "tag-test-01", "tag-test02", "1.0.0-beta.0+bA1"}
	expected := []string{"1.0.0-beta.0+bA1", "9.5.0", "10.1.0", "tag-test02", "tag-test-01"}
	for _, tag := range tags {
		rev := plumbing.Revision(branch.Hash().String())
		testlib.TestRepositoryCreateTag(r, tag, &rev)
	}

	tagset, _ := New(directory)
	all, _ := tagset.All()

	sorted := Sort(all)

	if len(sorted) != len(expected) {
		t.Fatalf("error: mismatch length of sort results")
	}

	for i, exp := range expected {
		act := sorted[i].Name().Short()
		if act != exp {
			t.Errorf("error: expected [%s] at [%d] actual [%s]", exp, i, act)
		}
	}

}

func TestTagsSortDirect(t *testing.T) {
	testlib.Testlogger(nil)
	directory, _, _ := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(directory)
	// tags to mock creating
	tags := []string{"9.5.0", "10.1.0", "tag-test-01", "tag-test02", "1.0.0-beta.0+bA1"}
	expected := []string{"1.0.0-beta.0+bA1", "9.5.0", "10.1.0", "tag-test02", "tag-test-01"}

	all := []*plumbing.Reference{}
	for _, tag := range tags {
		name := fmt.Sprintf("refs/tags/%s", tag)
		hash := "abc0123"
		ref := plumbing.NewReferenceFromStrings(name, hash)
		all = append(all, ref)
	}

	sorted := Sort(all)

	if len(sorted) != len(expected) {
		t.Fatalf("error: mismatch length of sort results")
	}

	for i, exp := range expected {
		act := sorted[i].Name().Short()
		if act != exp {
			t.Errorf("error: expected [%s] at [%d] actual [%s]", exp, i, act)
		}
	}

}
