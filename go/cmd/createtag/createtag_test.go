package createtag

import (
	"opg-github-actions/pkg/git/tags"
	"opg-github-actions/pkg/testlib"
	"os"
	"strings"
	"testing"
)

type ctFixture struct {
	Test     string
	Tags     []string
	GitRef   string
	Regen    bool
	Error    bool
	Expected map[string]string
	Contains bool
}

func TestCmdCreateTag(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []ctFixture{
		// No regen with a tag that already exists
		// will trigger an error
		{
			Test:   "1.0.0",
			Tags:   []string{"not-a-semver", "1.0.0", "others"},
			GitRef: "master", // head of the repo
			Regen:  false,
			Error:  true,
		},
		// Regen enabled, release tag that already exists, so should
		// generate a patch bump version
		{
			Test:     "1.0.0",
			Tags:     []string{"not-a-semver", "1.0.0", "others"},
			GitRef:   "master", // head of the repo
			Regen:    true,
			Expected: map[string]string{"requested_tag": "1.0.0", "created_tag": "1.0.1", "regenerated": "true"},
		},
		// Regen enabled with a prerelaese tag
		// creates a new tag that contains part of the existing tag
		{
			Test:     "1.0.0-beta.0+b1",
			Tags:     []string{"not-a-semver", "1.0.0-beta.0+b1", "others"},
			GitRef:   "master", // head of the repo
			Regen:    true,
			Expected: map[string]string{"created_tag": "1.0.0-", "regenerated": "true"},
			Contains: true,
		},
		// Creates the tag as requested, no issues
		{
			Test:     "1.0.0-beta.1+b1",
			Tags:     []string{"not-a-semver", "1.0.0-beta.0+b1", "others"},
			GitRef:   "master", // head of the repo
			Regen:    true,
			Expected: map[string]string{"created_tag": "1.0.0-beta.1+b1", "regenerated": "false"},
		},
		// Creates non-semver tag
		{
			Test:     "another-tag",
			Tags:     []string{"not-a-semver", "1.0.0-beta.0+b1", "others"},
			GitRef:   "master", // head of the repo
			Regen:    true,
			Expected: map[string]string{"created_tag": "another-tag", "regenerated": "false"},
		},
		// Creates non-semver tag that is regenerated
		{
			Test:     "not-a-semver",
			Tags:     []string{"not-a-semver", "1.0.0-beta.0+b1", "others"},
			GitRef:   "master", // head of the repo
			Regen:    true,
			Expected: map[string]string{"created_tag": "not-a-semver", "regenerated": "true"},
			Contains: true,
		},
	}

	dir, _, _ := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(dir)

	for _, f := range fixtures {
		tagSet, _ := tags.New(dir)

		res, err := process(tagSet, f.Test, f.Tags, f.GitRef, f.Regen, false)

		if f.Error {
			if err == nil {
				t.Errorf("error: expected an error")
			}
		} else {
			if err != nil {
				t.Errorf("error: unexpected error")
				t.Error(err)
			}

			for k, v := range f.Expected {
				if f.Contains {
					if !strings.Contains(res[k], v) {
						t.Errorf("error: [%s] expected to contain [%s] actual [%s]", k, v, res[k])
					}
				} else if res[k] != v {
					t.Errorf("error: [%s] expected [%s] actual [%s]", k, v, res[k])
				}
			}
		}

	}

}
