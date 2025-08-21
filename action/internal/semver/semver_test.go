package semver

import (
	"fmt"
	"opg-github-actions/action/internal/logger"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

type tBump struct {
	Content        []string
	Default        Increment
	Expected       Increment
	ExpectedCommit string
}

func TestSemverBumpCount(t *testing.T) {
	var lg = logger.New("error", "text")
	var tests = []*tBump{
		{
			Default:        PATCH,
			Expected:       MINOR,
			ExpectedCommit: "just a #minor",
			Content: []string{
				"captain, my commit mesage has nothing in",
				"just a #minor",
				"what a lovely #patch",
			},
		},
		{
			Default:        MAJOR,
			Expected:       MINOR,
			ExpectedCommit: "just a #minor",
			Content: []string{
				"major, my commit mesage has nothing in",
				"just a #minor",
				"what a lovely #patch",
			},
		},
		{
			Default:        MINOR,
			Expected:       MAJOR,
			ExpectedCommit: "breaking something with a #major",
			Content: []string{
				"my commit mesage has nothing in",
				"just a #minor",
				"what a lovely #patch",
				"breaking something with a #major",
			},
		},
		{
			Default:        PATCH,
			Expected:       MAJOR,
			ExpectedCommit: "breaking something with a #major",
			Content: []string{
				"major minor patch other",
				"breaking something with a #major",
				"a long bit of content that has many words and letters and all of those fun things without actually saying anything useful. #test",
			},
		},
		{
			Default:        PATCH,
			Expected:       MINOR,
			ExpectedCommit: "this is really only a !minor",
			Content: []string{
				"major minor patch other",
				"breaking something with a #major",
				"a long bit of content that has many words and letters and all of those fun things without actually saying anything useful. #test",
				"this is really only a !minor",
			},
		},
	}

	for i, test := range tests {
		actual, commit := GetBump(lg, test.Content, test.Default)

		if actual != test.Expected {
			t.Errorf("[%d] bump did not match, expected [%s] actual [%s]", i, test.Expected, actual)
		}
		if commit != test.ExpectedCommit {
			t.Errorf("[%d] commit did not match, expected [%s] actual [%s]", i, test.ExpectedCommit, commit)
		}
	}

}

type tSemverNextPre struct {
	Versions []*Semver
	Expected *Semver
	Suffix   string
	Bump     Increment
}

func TestSemverPrerelease(t *testing.T) {
	var lg = logger.New("error", "text")
	var tests = []*tSemverNextPre{
		{
			Bump:   MAJOR,
			Suffix: "beta",
			Versions: []*Semver{
				FromString("not-a-semver"),
				FromString("0.1.0"),
				FromString("1.0.0-beta.0+b1"),
				FromString("1.0.0-beta.1"),
				FromString("1.0.0-beta.2"),
				FromString("1.0.0-beta.10"),
				FromString("1.0.0-beta.101"),
				FromString("1.0.0-beta.0"),
			},
			Expected: FromString("1.0.0-beta.102"),
		},
		{
			Bump:   MINOR,
			Suffix: "beta",
			Versions: []*Semver{
				FromString("1.0.0"),
				FromString("10.0.0"),
				FromString("10.0.0-alpha.0"),
				FromString("not-a-semver"),
				FromString("1.0.0-beta.0+b1"),
				FromString("1.0.0-beta.0"),
			},
			Expected: FromString("10.1.0-beta.1"),
		},
		{
			Bump:   PATCH,
			Suffix: "---RC-SNAPSHOT.12.9.1--",
			Versions: []*Semver{
				FromString("1.2.2"),
				FromString("not-a-semver"),
				FromString("1.0.0-beta.0+b1"),
				FromString("1.2.3----RC-SNAPSHOT.12.9.1--.9+788"),
			},
			Expected: FromString("1.2.3----RC-SNAPSHOT.12.9.1--.10+788"),
		},
	}

	for _, test := range tests {
		actual := Prerelease(lg, test.Versions, test.Bump, test.Suffix)
		if actual.Stringy(true) != test.Expected.Stringy(true) {
			t.Errorf("error with next prerelease - expected [%s] actual [%s]", test.Expected, actual)
		}
	}
}

type tSemverNext struct {
	Versions []*Semver
	Bump     Increment
	Expected *Semver
}

func TestSemverNextRelease(t *testing.T) {
	var lg = logger.New("error", "text")
	var tests = []*tSemverNext{
		{
			Bump: MAJOR,
			Versions: []*Semver{
				FromString("v9.5.0"),
				FromString("tag-test-01"),
				FromString("v10.1.0"),
				FromString("tag-test-02"),
				FromString("v1.0.0-beta.1+bA1"),
				FromString("v1.0.0-beta.0+bA1"),
				FromString("v1.0.0"),
				FromString("v10.1.0"),
				FromString("v1.0.0-beta.0+bA2"),
			},
			Expected: FromString("v11.0.0"),
		},
		{
			Bump: MAJOR,
			Versions: []*Semver{
				FromString("v9.5.0"),
				FromString("tag-test-01"),
				FromString("10.1.0"),
				FromString("tag-test-02"),
				FromString("1.0.0-beta.1+bA1"),
				FromString("1.0.0-beta.0+bA1"),
				FromString("1.0.0"),
				FromString("10.1.0"),
				FromString("1.0.0-beta.0+bA2"),
			},
			Expected: FromString("11.0.0"),
		},
		{
			Bump: MINOR,
			Versions: []*Semver{
				FromString("v9.5.0"),
				FromString("tag-test-01"),
				FromString("10.1.0"),
				FromString("tag-test-02"),
				FromString("1.0.0-beta.1+bA1"),
				FromString("1.0.0-beta.0+bA1"),
				FromString("1.0.0"),
				FromString("10.1.0"),
				FromString("1.0.0-beta.0+bA2"),
			},
			Expected: FromString("10.2.0"),
		},
		{
			Bump: PATCH,
			Versions: []*Semver{
				FromString("v9.5.0"),
				FromString("tag-test-01"),
				FromString("tag.test-03"),
				FromString("10.1.0"),
				FromString("tag-test-02"),
				FromString("1.0.0-beta.1+bA1"),
				FromString("1.0.0-beta.0+bA1"),
				FromString("1.0.0"),
				FromString("10.1.0"),
				FromString("1.0.0-beta.0+bA2"),
			},
			Expected: FromString("10.1.1"),
		},
		{
			Bump:     PATCH,
			Versions: []*Semver{},
			Expected: FromString("0.0.1"),
		},
		{
			Bump:     MINOR,
			Versions: []*Semver{},
			Expected: FromString("0.1.0"),
		},
		{
			Bump:     MAJOR,
			Versions: []*Semver{},
			Expected: FromString("1.0.0"),
		},
		{
			Bump: MAJOR,
			Versions: []*Semver{
				FromString("1.0.0-beta.0+b1"),
				FromString("1.0.0-beta.0"),
			},
			Expected: FromString("1.0.0"),
		},
		{
			Bump: MINOR,
			Versions: []*Semver{
				FromString("1.0.0-beta.0+b1"),
				FromString("1.0.0-beta.0"),
			},
			Expected: FromString("0.1.0"),
		},
		{
			Bump: PATCH,
			Versions: []*Semver{
				FromString("1.0.0-beta.0+b1"),
				FromString("1.0.0-beta.0"),
			},
			Expected: FromString("0.0.1"),
		},
	}

	for _, test := range tests {
		actual := Release(lg, test.Versions, test.Bump)
		if actual.Stringy(true) != test.Expected.Stringy(true) {
			t.Errorf("error with next release - expected [%s] actual [%s]", test.Expected, actual)
		}
	}

}

type tSemverSort struct {
	Data     []*Semver
	Expected []*Semver
}

func TestSemverSort(t *testing.T) {
	var lg = logger.New("error", "text")
	// tests that include invalid semvers that wont be returned in the sorting
	// and duplicates that will be removed
	var tests = []*tSemverSort{
		{
			Data: []*Semver{
				FromString("v100.1.1"),
				FromString("9.5.0"),
				FromString("tag-test-01"),
				FromString("10.1.0"),
				FromString("tag-test-02"),
				FromString("1.0.0-beta.1+bA1"),
				FromString("1.0.0-beta.0+bA1"),
				FromString("10.1.0"),
				FromString("1.0.0-beta.0+bA2"),
				FromString("100.1.0"),
			},
			Expected: []*Semver{
				FromString("1.0.0-beta.0+bA1"),
				FromString("1.0.0-beta.0+bA2"),
				FromString("1.0.0-beta.1+bA1"),
				FromString("9.5.0"),
				FromString("10.1.0"),
				FromString("100.1.0"),
				FromString("v100.1.1"),
			},
		},
	}

	for i, test := range tests {
		sorted := Sort(lg, test.Data, SORT_ASC, true)

		if len(sorted) != len(test.Expected) {
			t.Errorf("semver sort test [%d] - mismatch length", i)
		} else {
			for idx, expected := range test.Expected {
				actual := sorted[idx]
				if expected.Stringy(true) != actual.Stringy(true) {
					t.Errorf("semver order not as expected in set [%d:%d], expected [%v] actual [%v]", i, idx, expected, actual)
				}

			}
		}

	}
}

type tFromStr struct {
	Ref      string
	Expected string
}

func TestSemverFromStringSuccess(t *testing.T) {
	var tests = []*tFromStr{
		{Ref: "0.0.4", Expected: "0.0.4"},
		{Ref: "1.2.3", Expected: "1.2.3"},
		{Ref: "10.20.30", Expected: "10.20.30"},
		{Ref: "1.1.2-prerelease+meta", Expected: "1.1.2-prerelease+meta"},
		{Ref: "1.1.2+meta", Expected: "1.1.2+meta"},
		{Ref: "1.1.2+meta-valid", Expected: "1.1.2+meta-valid"},
		{Ref: "1.0.0-alpha", Expected: "1.0.0-alpha"},
		{Ref: "1.0.0-beta", Expected: "1.0.0-beta"},
		{Ref: "1.0.0-alpha.beta", Expected: "1.0.0-alpha.beta"},
		{Ref: "1.0.0-alpha.beta.1", Expected: "1.0.0-alpha.beta.1"},
		{Ref: "1.0.0-alpha.1", Expected: "1.0.0-alpha.1"},
		{Ref: "1.0.0-alpha0.valid", Expected: "1.0.0-alpha0.valid"},
		{Ref: "1.0.0-alpha.0valid", Expected: "1.0.0-alpha.0valid"},
		{Ref: "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", Expected: "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay"},
		{Ref: "1.0.0-rc.1+build.1", Expected: "1.0.0-rc.1+build.1"},
		{Ref: "2.0.0-rc.1+build.123", Expected: "2.0.0-rc.1+build.123"},
		{Ref: "1.2.3-beta", Expected: "1.2.3-beta"},
		{Ref: "10.2.3-DEV-SNAPSHOT", Expected: "10.2.3-DEV-SNAPSHOT"},
		{Ref: "1.2.3-SNAPSHOT-123", Expected: "1.2.3-SNAPSHOT-123"},
		{Ref: "1.0.0", Expected: "1.0.0"},
		{Ref: "2.0.0", Expected: "2.0.0"},
		{Ref: "1.1.7", Expected: "1.1.7"},
		{Ref: "2.0.0+build.1848", Expected: "2.0.0+build.1848"},
		{Ref: "2.0.1-alpha.1227", Expected: "2.0.1-alpha.1227"},
		{Ref: "1.0.0-alpha+beta", Expected: "1.0.0-alpha+beta"},
		{Ref: "1.0.0+0.build.1-rc.10000aaa-kk-0.1", Expected: "1.0.0+0.build.1-rc.10000aaa-kk-0.1"},
		{Ref: "99999999999999999999999.999999999999999999.99999999999999999", Expected: "99999999999999999999999.999999999999999999.99999999999999999"},
	}

	for _, test := range tests {
		actual := FromString(test.Ref)
		if actual.Stringy(true) != test.Expected {
			t.Errorf("error: expected [%s] to be [%s] actual [%s]", test.Ref, test.Expected, actual.Stringy(true))
			fmt.Printf("%+v\n", actual)
		}
	}
}

func TestSemverFromStringFailures(t *testing.T) {
	var tests = []string{
		"1",
		"1.2",
		"1.2.3-0123.0123",
		"1.1.2+.123",
		"+invalid",
		"-invalid",
		"-invalid+invalid",
		"-invalid.01",
		"alpha",
		"alpha.beta",
		"alpha.beta.1",
		"alpha.1",
		"alpha+beta",
		"alpha_beta",
		"alpha.",
		"alpha..",
		"beta",
		"1.0.0-alpha_beta",
		"-alpha.",
		"1.0.0-alpha..",
		"1.0.0-alpha..1",
		"1.0.0-alpha...1",
		"1.0.0-alpha....1",
		"1.0.0-alpha.....1",
		"1.0.0-alpha......1",
		"1.0.0-alpha.......1",
		"01.1.1",
		"1.01.1",
		"1.1.01",
		"1.2",
		"1.2-SNAPSHOT",
		"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
		"1.2-RC-SNAPSHOT",
		"-1.0.3-gamma+b7718",
		"+justmeta",
		"9.8.7+meta+meta",
		"9.8.7-whatever+meta+meta",
		"1.2.3.DEV",
		"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
	}

	for _, test := range tests {
		if Valid(test) {
			t.Errorf("expected [%s] to be invalid", test)
		}
	}
}

type tNewSemver struct {
	Ref      *plumbing.Reference
	Expected *Semver
	Valid    bool
}

func TestSemverNew(t *testing.T) {
	var tests = []*tNewSemver{
		{
			Ref:      plumbing.NewReferenceFromStrings("refs/heads/v4", "6ecf0af"),
			Expected: &Semver{Valid: false, Original: "v4"},
			Valid:    false,
		},
		{
			Ref:      plumbing.NewReferenceFromStrings("refs/heads/v4.0.0", "6ecf0af"),
			Expected: &Semver{Valid: true, Prefix: "v", Major: "4", Minor: "0", Patch: "0"},
			Valid:    true,
		},
		{
			Ref:      plumbing.NewReferenceFromStrings("refs/heads/4.0.0", "6ecf0af"),
			Expected: &Semver{Valid: true, Prefix: "", Major: "4", Minor: "0", Patch: "0"},
			Valid:    true,
		},
		{
			Ref:      plumbing.NewReferenceFromStrings("refs/heads/s4.1.1", "6ecf0ab"),
			Expected: &Semver{Valid: true, Prefix: "s", Major: "4", Minor: "1", Patch: "1"},
			Valid:    true,
		},
		{
			Ref:      plumbing.NewReferenceFromStrings("refs/heads/4.1.1-beta.0+b1", "6ecf0ab"),
			Expected: &Semver{Valid: true, Prefix: "", Major: "4", Minor: "1", Patch: "1", PreleaseName: "beta", PrereleaseBuild: "0", BuildMetadata: "b1"},
			Valid:    true,
		},
		{
			Ref:      plumbing.NewReferenceFromStrings("refs/heads/1.0.0-alpha.beta.1+b1", "6ecf0ab"),
			Expected: &Semver{Valid: true, Prefix: "", Major: "1", Minor: "0", Patch: "0", PreleaseName: "alpha.beta", PrereleaseBuild: "1", BuildMetadata: "b1"},
			Valid:    true,
		},
	}

	for _, test := range tests {
		short := test.Ref.Name().Short()
		actual := New(test.Ref)
		if test.Valid && actual != nil {
			// check values match
			if !Equal(test.Expected, actual) {
				t.Errorf("semver equal failure with ref [%s] expected [%+v] actual [%+v]", short, test.Expected, actual)
			}
		} else if actual != nil && !test.Valid && actual.Valid {
			t.Errorf("error with ref [%s], expected to be invalid, returned a valid", short)
		}

	}

}

type tValid struct {
	Value    string
	Expected bool
}

func TestSemverValid(t *testing.T) {
	var tests = []*tValid{
		{
			Value:    "1.0.1",
			Expected: true,
		},
		{
			Value:    "1.0.0-beta.0",
			Expected: true,
		},
		{
			Value:    "1.0.0-beta.0+b1",
			Expected: true,
		},
		{
			Value:    "v1.0.0-beta.0+b1",
			Expected: true,
		},
		{
			Value:    "test-not-semver",
			Expected: false,
		},
		{
			Value:    "v1.not-semver",
			Expected: false,
		},
		{
			Value:    "v1.not.semver",
			Expected: false,
		},
		{
			Value:    "",
			Expected: false,
		},
	}

	for _, test := range tests {
		actual := Valid(test.Value)
		if actual != test.Expected {
			t.Errorf("error for [%s] expected [%v] actual [%v]", test.Value, test.Expected, actual)
		}
	}

}
