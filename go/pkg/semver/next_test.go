package semver

import (
	"opg-github-actions/pkg/testlib"
	"testing"
)

type nxFixture struct {
	Error    bool
	Expected string
	LP       string // last prerelease
	LR       string // last release
	Pre      bool
	Suffix   string
	C        IncrementCounters
}

func TestSemverNextTag(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []nxFixture{
		// patches
		{
			Expected: "0.0.1",
			LP:       "", LR: "",
			Pre: false, Suffix: "",
			C: IncrementCounters{Major: 0, Minor: 0, Patch: 1},
		},
		{
			Expected: "0.0.1-beta.0",
			LP:       "", LR: "",
			Pre: true, Suffix: "beta",
			C: IncrementCounters{Major: 0, Minor: 0, Patch: 1},
		},
		{
			Expected: "0.0.1-beta.1",
			LP:       "0.0.1-beta.0",
			LR:       "0.0.1",
			Pre:      true, Suffix: "beta",
			C: IncrementCounters{Major: 0, Minor: 0, Patch: 1},
		},
		// minors
		{
			Expected: "0.1.0",
			LP:       "", LR: "",
			Pre: false, Suffix: "",
			C: IncrementCounters{Major: 0, Minor: 1, Patch: 0},
		},
		{
			Expected: "0.1.0-beta.0",
			LP:       "", LR: "",
			Pre: true, Suffix: "beta",
			C: IncrementCounters{Major: 0, Minor: 1, Patch: 1},
		},
		{
			Expected: "0.1.0-beta.1",
			LP:       "0.1.0-beta.0",
			LR:       "0.0.1",
			Pre:      true, Suffix: "beta",
			C: IncrementCounters{Major: 0, Minor: 1, Patch: 1},
		},
		// majors
		{
			Expected: "1.0.0",
			LP:       "", LR: "",
			Pre: false, Suffix: "",
			C: IncrementCounters{Major: 1, Minor: 1, Patch: 0},
		},
		{
			Expected: "1.0.0-beta.0",
			LP:       "", LR: "",
			Pre: true, Suffix: "beta",
			C: IncrementCounters{Major: 1, Minor: 1, Patch: 1},
		},
		{
			Expected: "1.0.0-beta.1",
			LP:       "1.0.0-beta.0",
			LR:       "0.0.1",
			Pre:      true, Suffix: "beta",
			C: IncrementCounters{Major: 1, Minor: 1, Patch: 1},
		},
	}

	for _, f := range fixtures {

		lp, _ := New(f.LP)
		lr := Must(New(f.LR))

		n, err := Next(lp, lr, f.Pre, f.Suffix, f.C)

		if f.Error {
			if err == nil {
				t.Errorf("error: expected an error")
			}
		} else {
			if err != nil {
				t.Errorf("error: unexpected error")
				t.Error(err)
			}

			if n.String() != f.Expected {
				t.Errorf("error: expected [%s] actual [%v]", f.Expected, n.String())

			}

		}

	}
}
