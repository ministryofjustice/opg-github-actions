package semver

import (
	"opg-github-actions/pkg/testlib"
	"testing"
)

type pFixture struct {
	Test     string
	Expected bool
}

func TestPrereleaseValidation(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []pFixture{
		{Test: "prerelease", Expected: true},
		{Test: "prerelease.0", Expected: true},
		{Test: "alpha-a.b-c-somethinglong", Expected: true},
		{Test: "---RC-SNAPSHOT.0.9.1--.12", Expected: true},
		{Test: "alpha..", Expected: true},
	}

	for _, f := range fixtures {
		actual := ValidPrerelease(f.Test)
		if actual != f.Expected {
			t.Errorf("error: [%s] should be [%t], actual [%v]", f.Test, f.Expected, actual)
		}
	}
}

type bnFixture struct {
	Test     string // prerelease string to parse and get buildnumber from
	Valid    bool
	Expected *uint64
}

func TestPrereleaseBuildNumberAndBump(t *testing.T) {
	testlib.Testlogger(nil)
	i := uint64(10)
	fixtures := []bnFixture{
		{Test: "prerelease", Valid: true, Expected: nil},
		{Test: "prerelease.10", Valid: true, Expected: &i},
		{Test: "alpha-a.b-c-somethinglong", Valid: true, Expected: nil},
		{Test: "---RC-SNAPSHOT.1.2.3--.10", Valid: true, Expected: &i},
	}

	for _, f := range fixtures {
		p, e := NewPrerelease(f.Test)
		if f.Valid && e != nil {
			t.Errorf("error: unexpected error")
			t.Error(e)
		} else if !f.Valid && e == nil {
			t.Errorf("error: expected an error for [%s], did not happen", f.Test)
		}

		bn := p.BuildNumber()
		if f.Expected == nil && bn != nil {
			t.Errorf("error: expected nil result for [%s]; actual [%v]", f.Test, *bn)
		} else if bn != nil && f.Expected != nil && *bn != *f.Expected {
			t.Errorf("error: [%s] expected [%v], actual [%v]", f.Test, *f.Expected, *bn)
		}

		if bn != nil {
			i := *bn
			p.Bump()
			b := p.BuildNumber()
			if *b != (i + 1) {
				t.Errorf("failed to bump version")
			}
		}
	}
}

type mbFixture struct {
	Test      string
	Prefix    string
	ExpectedN *uint64
	ExpectedS string
}

func TestPrereleaseMustBump(t *testing.T) {
	testlib.Testlogger(nil)
	a := uint64(0)
	b := uint64(1)
	fixtures := []mbFixture{
		{Test: "beta", Prefix: "beta", ExpectedN: &a, ExpectedS: "beta.0"},
		{Test: "test.0", Prefix: "beta", ExpectedN: &b, ExpectedS: "test.1"},
		{Test: "alpha-a.b-c-somethinglong", Prefix: "beta", ExpectedN: &a, ExpectedS: "alpha-a.b-c-somethinglong.0"},
		{Test: "", Prefix: "beta", ExpectedN: &a, ExpectedS: "beta.0"},
	}

	for _, f := range fixtures {
		p, _ := NewPrerelease(f.Test)
		p.MustBump(f.Prefix)

		bn := p.BuildNumber()
		if *bn != *f.ExpectedN {
			t.Errorf("error: expected [%s] to be [%v] actual [%v]", f.Test, *f.ExpectedN, *bn)
		}
		if p.String() != f.ExpectedS {
			t.Errorf("error: expected [%s] to be [%v] actual [%v]", f.Test, f.ExpectedS, p.String())
		}
	}
}
