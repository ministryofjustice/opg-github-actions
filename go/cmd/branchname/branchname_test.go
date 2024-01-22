package branchname

import (
	"opg-github-actions/pkg/testlib"
	"testing"
)

// inFixture are the input arguments
type testFixture struct {
	Event    string
	Content  []byte
	Error    bool
	Expected map[string]string
}

func TestBranchNameCleaned(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []testFixture{
		{
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "my-feature-branch", "title", "body"),
			Expected: map[string]string{"full_length": "myfeaturebranch", "safe": "myfeaturebra", "branch_name": "my-feature-branch"},
		},
		{
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "my-feature??", "title", "body"),
			Expected: map[string]string{"full_length": "myfeature", "safe": "myfeature", "branch_name": "my-feature??"},
		},
		{
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "long/string/with/slashes", "title", "body"),
			Expected: map[string]string{"full_length": "longstringwithslashes", "safe": "longstringwi"},
		},
		{
			Event:    "pull_request",
			Content:  testlib.TestEventPullRequest("main", "long/string-with-others?!.><~@#", "title", "body"),
			Expected: map[string]string{"full_length": "longstringwithothers", "safe": "longstringwi"},
		},
		{
			Event:    "push",
			Content:  testlib.TestEventPush("long/string-with-others?!.><~@#♥", "commit123", "commit1234"),
			Expected: map[string]string{"full_length": "longstringwithothers", "safe": "longstringwi", "branch_name": "long/string-with-others?!.><~@#♥"},
		},
	}
	for _, f := range fixtures {
		out, e := process(f.Event, f.Content)

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
