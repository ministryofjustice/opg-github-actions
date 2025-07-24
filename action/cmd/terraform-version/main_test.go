package main

import "testing"

type tvFixture struct {
	Test     string
	Expected string
}

func TestTerraformVersion(t *testing.T) {

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
	}

	for i, f := range fixtures {
		out, e := versionData(f.Test)
		if e != nil {
			t.Errorf("error: unexpected error for test [%d]: %s", i, e.Error())
		}

		if out["version"] != f.Expected {
			t.Errorf("error: expected [%s] actual [%s]", f.Expected, out["version"])
		}
	}
}
