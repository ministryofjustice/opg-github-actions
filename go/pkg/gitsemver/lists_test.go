package gitsemver

import (
	"opg-github-actions/pkg/testlib"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

type tgFixture struct {
	Refs     []*plumbing.Reference
	Expected int
}

func TestGitSemverPrereleases(t *testing.T) {
	testlib.Testlogger(nil)
	fixtures := []tgFixture{
		{
			Expected: 1,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-beta.1", "commit"),
				plumbing.NewReferenceFromStrings("test-normal-string", "commit"),
			},
		},
		{
			Expected: 0,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-beta.01", "commit"),
				plumbing.NewReferenceFromStrings("test-normal-string", "commit"),
			},
		},
		{
			Expected: 3,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-beta.0+b1", "commit"),
				plumbing.NewReferenceFromStrings("test-normal-string", "commit"),
				plumbing.NewReferenceFromStrings("10.1.1-beta.0+b1", "commit"),
				plumbing.NewReferenceFromStrings("10.1.1-beta-with-extras--", "commit"),
			},
		},
	}

	for _, f := range fixtures {
		actual := Prereleases(f.Refs)

		if len(actual) != f.Expected {
			t.Errorf("error: expected [%d] results, actual [%d]", f.Expected, len(actual))
		}
	}
}

func TestGitSemverMatchingPrereleases(t *testing.T) {
	testlib.Testlogger(nil)
	fixtures := []tgFixture{
		{
			Expected: 2,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-beta.1", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-beta.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-test.0", "commit"),
				plumbing.NewReferenceFromStrings("test-normal-string", "commit"),
			},
		},
		{
			Expected: 0,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1----beta.1", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-notbeta.1", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-notbeta.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-test.0", "commit"),
				plumbing.NewReferenceFromStrings("test-normal-string", "commit"),
			},
		},
		{
			Expected: 0,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.0", "commit"),
				plumbing.NewReferenceFromStrings("11.0.1-notbeta.0", "commit"),
				plumbing.NewReferenceFromStrings("11.0.1-beta--what.0", "commit"),
				plumbing.NewReferenceFromStrings("v1.0.1-test.0", "commit"),
				plumbing.NewReferenceFromStrings("test-normal-string", "commit"),
			},
		},
	}

	for _, f := range fixtures {
		actual := MatchingPrereleases("beta", f.Refs)
		if len(actual) != f.Expected {
			t.Errorf("error: expected [%d] results, actual [%d]", f.Expected, len(actual))
		}
	}
}
