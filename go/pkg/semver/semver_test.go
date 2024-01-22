package semver

import (
	"fmt"
	"opg-github-actions/pkg/testlib"
	"testing"
)

type mFixture struct {
	Test       string // semver string to test
	Invalid    bool   // if the semver string is valid
	ParseError bool   // will the semver parse ok
	Expected   map[string]string
}

func TestSemverMap(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []mFixture{
		{Test: "1.0.0", Invalid: false, ParseError: false, Expected: map[string]string{"prefix": "", "major": "1", "minor": "0", "patch": "0"}},
		{Test: "1.0.0-beta.0", Invalid: false, ParseError: false, Expected: map[string]string{"prefix": "", "major": "1", "minor": "0", "patch": "0", "prerelease": "beta.0"}},
		{Test: "1.0.0-beta.0+b1", Invalid: false, ParseError: false, Expected: map[string]string{"prefix": "", "major": "1", "minor": "0", "patch": "0", "prerelease": "beta.0", "buildmetadata": "b1"}},
		{Test: "", Invalid: true, ParseError: true},
		{Test: "test-not-semver", Invalid: true, ParseError: false},
	}

	for _, f := range fixtures {
		s, err := New(f.Test)

		if f.Invalid && err == nil {
			t.Errorf("error: expected validity error, did not happen")
		} else if !f.Invalid && err != nil {
			t.Error("error: unexpcted error")
			t.Error(err)
		}

		if s != nil {
			m, err := s.Map()

			if f.ParseError && err == nil {
				t.Errorf("error: expected parse error, did not happen")
			} else if !f.ParseError && err != nil {
				t.Error("error: unexpcted error")
				t.Error(err)
			}

			for k, v := range f.Expected {
				if m[k] != v {
					t.Errorf("error: expected [%s] to be [%s]; actual [%v]", k, v, m[k])
				}
			}
		}

	}

}

func TestSemverValidConversions(t *testing.T) {
	testlib.Testlogger(nil)
	fixtures := []string{
		"0.0.4",
		"1.2.3",
		"10.20.30",
		"1.1.2-prerelease+meta",
		"1.1.2+meta",
		"1.1.2+meta-valid",
		"1.0.0-alpha",
		"1.0.0-beta",
		"1.0.0-alpha.beta",
		"1.0.0-alpha.beta.1",
		"1.0.0-alpha.1",
		"1.0.0-alpha0.valid",
		"1.0.0-alpha.0valid",
		"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
		"1.0.0-rc.1+build.1",
		"2.0.0-rc.1+build.123",
		"1.2.3-beta",
		"10.2.3-DEV-SNAPSHOT",
		"1.2.3-SNAPSHOT-123",
		"1.0.0",
		"2.0.0",
		"1.1.7",
		"2.0.0+build.1848",
		"2.0.1-alpha.1227",
		"1.0.0-alpha+beta",
		"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
		"1.2.3----R-S.12.9.1--.12+meta",
		"1.2.3----RC-SNAPSHOT.12.9.1--.12",
		"1.0.0+0.build.1-rc.10000aaa-kk-0.1",
		"99999999999999999999999.999999999999999999.99999999999999999",
		"1.0.0-0A.is.legal",
	}

	for _, f := range fixtures {
		sv, err := New(f)
		if err != nil {
			t.Errorf("error: unexpected error")
			t.Error(err)
		}
		if sv.String() != f {
			t.Errorf("error: expected [%s] actual [%v]", f, sv.String())
		}
	}

}

func TestSemverInvalidConversions(t *testing.T) {
	testlib.Testlogger(nil)
	fixtures := []string{
		"1",
		"1.2",
		"1.2.3-0123",
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
		"1.2.3.DEV",
		"1.2-SNAPSHOT",
		"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
		"1.2-RC-SNAPSHOT",
		"-1.0.3-gamma+b7718",
		"+justmeta",
		"9.8.7+meta+meta",
		"9.8.7-whatever+meta+meta",
		"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
	}

	for _, f := range fixtures {
		sem, e := New(f)
		if e == nil {
			t.Errorf("expected an error for %s", sem)
		}

	}

}

type cFixture struct {
	Test     string
	Expected string
}

func TestSemverChangestoMajor(t *testing.T) {
	testlib.Testlogger(nil)
	fixtures := []cFixture{
		{Test: "1.0.0", Expected: "2.0.0"},
		{Test: "1.0.0-beta", Expected: "2.0.0-beta"},
		{Test: "1.0.0-beta.0", Expected: "2.0.0-beta.0"},
		{Test: "1.0.0-beta.0+b1", Expected: "2.0.0-beta.0+b1"},
	}

	for _, f := range fixtures {
		s, _ := New(f.Test)
		s.BumpMajor()
		if s.String() != f.Expected {
			fmt.Println(s)
			t.Errorf("error: expected [%s] actual [%s]", f.Expected, s.String())
		}
	}

}

type pbFixture struct {
	Test      string
	Error     bool
	Expected  *uint64
	ExpectedS string
}

func TestPrereleaseBumping(t *testing.T) {
	testlib.Testlogger(nil)
	i := uint64(10)
	fixtures := []pbFixture{
		{Test: "1.0.0", Error: true, Expected: nil},
		{Test: "1.0.0-beta.9", Error: false, Expected: &i, ExpectedS: "1.0.0-beta.10"},
		{Test: "1.2.3----RC-SNAPSHOT.12.9.1--.9+788", Error: false, Expected: &i, ExpectedS: "1.2.3----RC-SNAPSHOT.12.9.1--.10+788"},
		{Test: "1.0.0-beta+5.12", Error: true, Expected: nil},
		{Test: "1.0.0-beta-01-12", Error: true, Expected: nil},
	}

	for _, f := range fixtures {
		s, _ := New(f.Test)

		err := s.BumpPrerelease()
		if f.Error && err == nil {
			t.Errorf("error: expected an error for [%s]", f.Test)
		} else if !f.Error && err != nil {
			t.Errorf("error: unexpected an error for [%s]", f.Test)
			t.Error(err)
		}

		actual := s.BuildNumber()

		if f.Expected == nil && actual != nil {
			t.Errorf("error: expected nil, actual [%v]", *actual)
		} else if actual != nil && *actual != *f.Expected {
			t.Errorf("error: [%s] expected [%v] actual [%v]", f.Test, *f.Expected, *actual)
		}

		if f.ExpectedS != "" && f.ExpectedS != s.String() {
			t.Errorf("error: expected generated string to be [%s], actual [%v]", f.ExpectedS, s.String())
		}

	}
}

func TestSemverSetPrereleasePrefix(t *testing.T) {
	testlib.Testlogger(nil)

	sv, e := New("1.0.0-beta.0+ba1")

	if e != nil {
		t.Errorf("error: unexpected error")
		t.Error(e)
	}

	sv.SetPrereleasePrefix("test")

	if sv.PrereleasePrefix() != "test" {
		t.Errorf("error: failed to updated prerelease prefix")
	}
}
