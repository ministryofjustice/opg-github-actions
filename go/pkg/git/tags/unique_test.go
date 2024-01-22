package tags

import (
	"opg-github-actions/pkg/testlib"
	"os"
	"strings"
	"testing"
)

type uFixture struct {
	Test     string
	Error    bool
	Contains string
	Equals   bool
	Expected string
	TagSet   []string
}

func TestUniqueGeneration(t *testing.T) {
	testlib.Testlogger(nil)
	dir, r, _ := testlib.TestRepositorySkeleton()
	defer os.RemoveAll(dir)

	fixtures := []uFixture{
		// a non semver tag that does exist should
		// generate a new tag that contains the old one
		{
			Test:     "non-semver-tag-that-exists",
			TagSet:   []string{"non-semver-tag-that-exists", "others", "1.0.0"},
			Error:    false,
			Contains: "non-semver-tag-that-exists",
			Equals:   false,
		},
		// a new tag that doesnt exist so should be the same
		{
			Test:   "1.0.0-beta.0+b1",
			TagSet: []string{"non-semver-tag-that-exists", "others", "1.0.0"},
			Error:  false,
			Equals: true,
		},
		// a release tag that already exists, this should
		// increment the patch number only
		{
			Test:     "1.0.0",
			TagSet:   []string{"non-semver-tag-that-exists", "others", "1.0.0"},
			Error:    false,
			Equals:   false,
			Expected: "1.0.1",
		},
		// a prerelease that already exists, this should change
		// the prerelease prefix segment of the semver
		{
			Test:     "1.0.0-beta.0",
			TagSet:   []string{"non-semver-tag-that-exists", "others", "1.0.0-beta.0", "1.0.0-beta.1"},
			Error:    false,
			Contains: "1.0.0-",
			Equals:   false,
		},
	}

	for _, f := range fixtures {
		tagset := &Tags{repository: r, Directory: dir}
		n, e := tagset.Unique(f.Test, f.TagSet)

		if f.Error {
			if e == nil {
				t.Errorf("error: expected an error, none recieved.")
			}
		} else {
			if e != nil {
				t.Errorf("error: unexpected error")
				t.Error(e)
			}

			if len(f.Contains) > 0 {
				if !strings.HasPrefix(n, f.Contains) {
					t.Errorf("error: expected [%s] result to contain [%s], actual [%s]", f.Test, f.Contains, n)
				}
			} else if f.Equals {
				if n != f.Test {
					t.Errorf("error: expected [%s], actual [%s]", f.Test, n)
				}
			} else if len(f.Expected) > 0 {
				if n != f.Expected {
					t.Errorf("error: [%s] expected [%s], actual [%s]", f.Test, f.Expected, n)
				}
			}

		}

	}
}
