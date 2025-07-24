package main

import (
	"opg-github-actions/action/internal/logger"
	"testing"
)

type cmdFixture struct {
	Length   int
	Source   string
	Error    bool
	Expected map[string]string
}

func TestBranchNameCommandWorking(t *testing.T) {
	lg := logger.New("ERROR", "TEXT")
	tests := []cmdFixture{
		{

			Source:   "my-feature-branch",
			Length:   20,
			Expected: map[string]string{"full_length": "myfeaturebranch", "safe": "myfeaturebranch", "branch_name": "my-feature-branch"},
		},
		{
			Length:   12,
			Source:   "my-feature??",
			Expected: map[string]string{"full_length": "myfeature", "safe": "myfeature", "branch_name": "my-feature??"},
		},
		{
			Length:   maxLength,
			Source:   "long/string/with/slashes",
			Expected: map[string]string{"full_length": "longstringwithslashes", "safe": "longstringwi"},
		},
		{
			Length:   maxLength,
			Source:   "long/string-with-others?!.><~@#",
			Expected: map[string]string{"full_length": "longstringwithothers", "safe": "longstringwi"},
		},
		{
			Length:   maxLength,
			Source:   "long/string-with-others?!.><~@#♥",
			Expected: map[string]string{"full_length": "longstringwithothers", "safe": "longstringwi", "branch_name": "long/string-with-others?!.><~@#♥"},
		},
		{
			Length:   6,
			Source:   "long/string",
			Expected: map[string]string{"full_length": "longstring", "safe": "longst", "branch_name": "long/string"},
		},
		{
			Length:   14,
			Source:   "renovate/my-feature-thingy-update",
			Expected: map[string]string{"full_length": "renovatemyfeaturethingyupdate", "safe": "renovatemyfeat", "branch_name": "renovate/my-feature-thingy-update"},
		},
	}

	for _, test := range tests {
		actual, err := Run(lg, test.Source, test.Length)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		for k, v := range test.Expected {
			if actual[k] != v {
				t.Errorf("error: [%s] expected [%s] actual [%s]", k, v, actual[k])
			}
		}
	}
}
