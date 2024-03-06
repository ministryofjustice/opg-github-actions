package latesttag

import (
	"opg-github-actions/pkg/testlib"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

type ltFixture struct {
	ExpectedPre string
	ExpectedRe  string
	Tags        []*plumbing.Reference
	Prerelease  bool
}

func TestLatestTag(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []ltFixture{
		{
			ExpectedPre: "v1.0.0-beta.0",
			ExpectedRe:  "1.0.0",
			Prerelease:  true,
			Tags: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("0.0.0", ""),
				plumbing.NewReferenceFromStrings("1.0.0", ""),
				plumbing.NewReferenceFromStrings("0.1.0", ""),
				plumbing.NewReferenceFromStrings("v1.0.0-beta.0", ""),
			},
		},
		{
			ExpectedPre: "",
			ExpectedRe:  "1.0.0",
			Prerelease:  true,
			Tags: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("0.0.0", ""),
				plumbing.NewReferenceFromStrings("1.0.0", ""),
				plumbing.NewReferenceFromStrings("0.1.0", ""),
				plumbing.NewReferenceFromStrings("v1.0.0-test.0", ""),
			},
		},
		{
			ExpectedPre: "10.0.1-beta.10",
			ExpectedRe:  "1.0.0",
			Prerelease:  true,
			Tags: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("0.0.0", ""),
				plumbing.NewReferenceFromStrings("1.0.0", ""),
				plumbing.NewReferenceFromStrings("9.1.1-beta.0", ""),
				plumbing.NewReferenceFromStrings("10.0.1-beta.0", ""),
				plumbing.NewReferenceFromStrings("10.0.1-beta.9+b1", ""),
				plumbing.NewReferenceFromStrings("10.0.1-beta.10", ""),
			},
		},
		{
			ExpectedPre: "",
			ExpectedRe:  "10.10.1",
			Prerelease:  false,
			Tags: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("0.0.0", ""),
				plumbing.NewReferenceFromStrings("1.0.0", ""),
				plumbing.NewReferenceFromStrings("9.1.1", ""),
				plumbing.NewReferenceFromStrings("10.0.1-beta.0", ""),
				plumbing.NewReferenceFromStrings("10.0.1-beta.9+b1", ""),
				plumbing.NewReferenceFromStrings("10.0.1-beta.10", ""),
				plumbing.NewReferenceFromStrings("10.0.1", ""),
				plumbing.NewReferenceFromStrings("10.9.1", ""),
				plumbing.NewReferenceFromStrings("10.10.1", ""),
			},
		},
		{
			ExpectedPre: "",
			ExpectedRe:  "v1.0.0", // v prefix makes it high / latest than a number
			Prerelease:  false,
			Tags: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("0.0.0", ""),
				plumbing.NewReferenceFromStrings("v1.0.0", ""),
				plumbing.NewReferenceFromStrings("9.1.1", ""),
				plumbing.NewReferenceFromStrings("10.10.1", ""),
			},
		},
	}
	for _, f := range fixtures {
		actual, e := process(f.Tags, f.Prerelease, "beta", "beta", "master,main")
		if e != nil {
			t.Errorf("unexpected error")
			t.Error(e)
		}
		if actual["last_prerelease"] != f.ExpectedPre {
			t.Errorf("error: expected [%s] actual [%s]", f.ExpectedPre, actual["last_prerelease"])
		}
		if actual["last_release"] != f.ExpectedRe {
			t.Errorf("error: expected [%s] actual [%s]", f.ExpectedRe, actual["last_release"])
		}
	}

}
