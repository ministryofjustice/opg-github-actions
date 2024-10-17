package branchname

import (
	"opg-github-actions/pkg/testlib"
	"testing"
)

// inFixture are the input arguments
type testFixture struct {
	Event    string
	Length   int
	Content  []byte
	Error    bool
	Expected map[string]string
}

func TestBranchNameCleaned(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []testFixture{
		{
			Length:   DefaultMaxLength,
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "my-feature-branch", "title", "body"),
			Expected: map[string]string{"full_length": "myfeaturebranch", "safe": "myfeaturebra", "branch_name": "my-feature-branch"},
		},
		{
			Length:   DefaultMaxLength,
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "my-feature??", "title", "body"),
			Expected: map[string]string{"full_length": "myfeature", "safe": "myfeature", "branch_name": "my-feature??"},
		},
		{
			Length:   DefaultMaxLength,
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "long/string/with/slashes", "title", "body"),
			Expected: map[string]string{"full_length": "longstringwithslashes", "safe": "longstringwi"},
		},
		{
			Length:   DefaultMaxLength,
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "long/string-with-others?!.><~@#", "title", "body"),
			Expected: map[string]string{"full_length": "longstringwithothers", "safe": "longstringwi"},
		},
		{
			Length:   DefaultMaxLength,
			Event:    "push",
			Content:  testlib.TestEventPush("long/string-with-others?!.><~@#♥", "commit123", "commit1234"),
			Expected: map[string]string{"full_length": "longstringwithothers", "safe": "longstringwi", "branch_name": "long/string-with-others?!.><~@#♥"},
		},
		{
			Length:   6,
			Event:    "push",
			Content:  testlib.TestEventPush("long/string", "commit123", "commit1234"),
			Expected: map[string]string{"full_length": "longstring", "safe": "longst", "branch_name": "long/string"},
		},
		{
			Length:   14,
			Event:    "push",
			Content:  testlib.TestEventPush("renovate/my-feature-thingy-update", "commit123", "commit1234"),
			Expected: map[string]string{"full_length": "renovatemyfeaturethingyupdate", "safe": "renovatemyfeat", "branch_name": "renovate/my-feature-thingy-update"},
		},
	}
	for _, f := range fixtures {
		out, e := process(f.Event, f.Length, f.Content)

		if e != nil {
			t.Errorf("error: unexpected error")
			t.Error(e)
		}

		for k, v := range f.Expected {
			if out[k] != v {
				t.Errorf("error: [%s] expected [%s] actual [%s]", k, v, out[k])
			}
		}
	}

}
