package terraformversion

import "testing"

type tvFixture struct {
	Test     string
	Expected string
	Simple   bool
}

func TestTerraformVersionProcess(t *testing.T) {

	fixtures := []tvFixture{
		{
			Test: `terraform {
				required_version = "1.6.5"
			}`,
			Expected: "1.6.5",
		},
		{
			Test:     `terraform {required_version = "1.6.0"}`,
			Expected: "1.6.0",
		},
		{
			Test:     `terraform {required_version = ">= 1.1.0"}`,
			Expected: ">= 1.1.0",
		},
		{
			Test:     `1.0.0`,
			Expected: "1.0.0",
			Simple:   true,
		},
	}

	for _, f := range fixtures {
		out, e := process(f.Test, f.Simple)
		if e != nil {
			t.Errorf("error: unexpected error")
			t.Error(e)
		}

		if out["version"] != f.Expected {
			t.Errorf("error: expected [%s] actual [%s]", f.Expected, out["version"])
		}
	}
}
