package tags

import (
	"opg-github-actions/pkg/testlib"
	"os"
	"slices"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

// TestTagsList creates a test repo and checks :
// - number tags found matches number of tags created
// - value of tags found matches value of tags originally created
func TestTagsList(t *testing.T) {
	testlib.Testlogger(nil)
	directory, r, branch := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(directory)
	// create some tags in the repository and then fetch to make sure they are found
	tags := []string{"tag-test-01", "tag-test02", "1.0.0-beta.0+bA1"}
	for _, tag := range tags {
		rev := plumbing.Revision(branch.Hash().String())
		testlib.TestRepositoryCreateTag(r, tag, &rev)
	}

	tagset, _ := New(directory)
	all, _ := tagset.All()
	// number of tags should match
	if len(all) != len(tags) {
		t.Errorf("error: incorrect number of tags found. Expected [%d] Actual [%d]", len(tags), len(all))
	}

	matched := 0
	for _, t := range all {
		if slices.Contains(tags, t.Name().Short()) {
			matched++
		}
	}
	if matched != len(tags) {
		t.Errorf("error: expected [%d] matches would be found, actual [%d]", len(tags), matched)
	}
}

func TestTagsAt(t *testing.T) {
	testlib.Testlogger(nil)
	directory, r, _ := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(directory)
	//
	shouldHave := map[plumbing.Hash]int{}
	// make some new commits and track the locations
	commits := []string{"test 1", "test 2", "test 3", "test 4"}
	created := map[plumbing.Hash]string{}
	for _, c := range commits {
		h, _ := testlib.TestRepositoryCommit(r, c)
		created[h] = c
	}
	// now make some tags at each of them
	testTags := [][]string{
		{"test-tag-1", "test-2", "0.0.1"},
		{"dummy-tag"},
		{"v1.0.0", "v1.0.0-beta.01"},
		{"v2.0.0"},
	}
	i := 0
	for hash, _ := range created {
		tags := testTags[i]
		rev := plumbing.Revision(hash.String())
		for _, tagName := range tags {
			testlib.TestRepositoryCreateTag(r, tagName, &rev)
		}
		shouldHave[hash] = len(tags)
		i++
	}
	if len(shouldHave) <= 0 {
		t.Errorf("No tags created")
	}
	// Now check the tags are created correctly
	tagset, _ := New(directory)
	for hash, count := range shouldHave {
		at, _ := tagset.At(&hash)
		if len(at) != count {
			t.Errorf("error: Failed to find correct number of tags at [%s] Expected [%d] Actual [%d]", hash.String(), count, len(at))
		}
	}

}
