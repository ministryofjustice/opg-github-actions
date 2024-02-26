package safestring

import (
	"opg-github-actions/pkg/testlib"
	"testing"
)

// inFixture are the input arguments
type testFixture struct {
	Original         string
	Length           int
	Suffix           string
	ConditionalMatch string
	ConditionalValue string
	Expected         map[string]string
}

func TestSafeString(t *testing.T) {
	testlib.Testlogger(nil)

	fixtures := []testFixture{
		{Original: "long/string/with/slashes", Expected: map[string]string{"full_length": "longstringwithslashes", "safe": "longstringwithslashes"}},
		{Original: "long/string/with/slashes", Length: 4, Expected: map[string]string{"full_length": "longstringwithslashes", "safe": "long"}},
		{Original: "1).what-is-this?", Expected: map[string]string{"full_length": "1whatisthis", "safe": "1whatisthis"}},
		{Original: "is this ok? #; ? <> + = \\", Expected: map[string]string{"full_length": "isthisok", "safe": "isthisok"}},
	}
	for _, f := range fixtures {
		out, e := process(f.Original, f.Suffix, f.Length, f.ConditionalMatch, f.ConditionalValue)

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
