package nexttag

import (
	"opg-github-actions/pkg/semver"
	"opg-github-actions/pkg/testlib"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type ntFixture struct {
	LastPrerelease   string
	LastRelease      string
	Prerelease       bool
	PrereleaseSuffix string
	DefaultBump      semver.Increment
	WithV            bool
	Commits          []*object.Commit
	Expected         map[string]string
}

func TestNextTagProcess(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []ntFixture{
		// release - everything empty, so the outcome should be the default bump
		{
			LastRelease:      "",
			LastPrerelease:   "",
			Prerelease:       false,
			PrereleaseSuffix: "",
			DefaultBump:      semver.Patch,
			WithV:            false,
			Commits:          []*object.Commit{},
			Expected:         map[string]string{"next_tag": "0.0.1", "patches": "1"},
		},
		// release - 1 commit with minor and nothing else
		{
			LastRelease:      "",
			LastPrerelease:   "",
			Prerelease:       false,
			PrereleaseSuffix: "",
			DefaultBump:      semver.Patch,
			WithV:            false,
			Commits: []*object.Commit{
				{Message: "not a patch"},
				{Message: "a #minor change"},
				{Message: `?? % & # >< : \ = - + ♥

				end here`},
			},
			Expected: map[string]string{"next_tag": "0.1.0", "minors": "1"},
		},
		// release - with an existing release
		{
			LastRelease:      "1.0.0",
			LastPrerelease:   "1.1.0-beta.0",
			Prerelease:       false,
			PrereleaseSuffix: "",
			DefaultBump:      semver.Patch,
			WithV:            false,
			Commits: []*object.Commit{
				{Message: "#patch"},
				{Message: "a #minor change"},
				{Message: `?? % & # >< : \ = - + ♥

				end here`},
			},
			Expected: map[string]string{"next_tag": "1.1.0", "minors": "1", "patches": "1"},
		},
		// pre-release
		{
			LastRelease:      "",
			LastPrerelease:   "",
			Prerelease:       true,
			PrereleaseSuffix: "beta",
			DefaultBump:      semver.Patch,
			WithV:            false,
			Commits: []*object.Commit{
				{Message: "not a patch"},
				{Message: "a #minor change"},
				{Message: `?? % & # >< : \ = - + ♥

				end here`},
			},
			Expected: map[string]string{"next_tag": "0.1.0-beta.0", "minors": "1"},
		},
		{
			LastRelease:      "0.0.1",
			LastPrerelease:   "0.1.0-beta.0",
			Prerelease:       true,
			PrereleaseSuffix: "beta",
			DefaultBump:      semver.Patch,
			WithV:            false,
			Commits: []*object.Commit{
				{Message: "not a patch"},
				{Message: "a #minor change"},
				{Message: `?? % & # >< : \ = - + ♥
				
				end here`},
			},
			Expected: map[string]string{"next_tag": "0.1.0-beta.1", "minors": "1"},
		},
		{
			LastRelease:      "0.0.1",
			LastPrerelease:   "0.1.0-beta.0",
			Prerelease:       true,
			PrereleaseSuffix: "beta",
			DefaultBump:      semver.Patch,
			WithV:            false,
			Commits: []*object.Commit{
				{Message: "not a patch"},
				{Message: "a #major change"},
				{Message: `?? % & # >< : \ = - + ♥
				
				end here`},
			},
			Expected: map[string]string{"next_tag": "1.0.0-beta.0", "majors": "1"},
		},
		{
			LastRelease:      "0.0.1",
			LastPrerelease:   "1.0.0-beta.0",
			Prerelease:       true,
			PrereleaseSuffix: "beta",
			DefaultBump:      semver.Patch,
			WithV:            true,
			Commits: []*object.Commit{
				{Message: "not a patch"},
				{Message: "a #major change"},
				{Message: `?? % & # >< : \ = - + ♥
				
				end here`},
			},
			Expected: map[string]string{"next_tag": "v1.0.0-beta.1", "majors": "1"},
		},

		{
			LastRelease:      "v1.0.1",
			LastPrerelease:   "v2.0.0-beta.0",
			Prerelease:       true,
			PrereleaseSuffix: "beta",
			DefaultBump:      semver.Patch,
			WithV:            false,
			Commits: []*object.Commit{
				{Message: "not a patch"},
				{Message: "a #major change"},
				{Message: `?? % & # >< : \ = - + ♥
				
				end here`},
			},
			Expected: map[string]string{"next_tag": "2.0.0-beta.1", "majors": "1"},
		},
	}

	for _, f := range fixtures {
		out, err := process(f.Commits, f.Prerelease, f.PrereleaseSuffix, f.LastPrerelease, f.LastRelease, f.DefaultBump, f.WithV)
		if err != nil {
			t.Errorf("error: unexpected error")
			t.Error(err)
		}
		for k, v := range f.Expected {
			if out[k] != v {
				t.Errorf("error: expected [%s] to be [%s] actual [%v]", k, v, out[k])

			}
		}
		// pp.Println(f.Expected)
		// pp.Println(out)
		// println("---")

	}
}
