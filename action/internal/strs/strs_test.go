package strs

import "testing"

type cleanFixture struct {
	Test     string
	Expected string
}

func TestCleanStrings(t *testing.T) {

	var tests = []*cleanFixture{
		{Test: "dependabot/loves/shashes", Expected: "dependabotlovesshashes"},
		{Test: "no-more-hyphens", Expected: "nomorehyphens"},
		{Test: "no_serperation-any more", Expected: "noserperationanymore"},
		{Test: "not just letters! 123 () ? > * & . \\ #", Expected: "notjustletters123"},
		{Test: "exactlythesame", Expected: "exactlythesame"},
		{Test: "", Expected: ""},
	}

	for _, test := range tests {
		actual := Clean(test.Test)
		if actual != test.Expected {
			t.Errorf("error cleaning string [%s], expected [%s] actual [%s]", test.Test, test.Expected, actual)
		}
	}

}

type truncateFixture struct {
	Test     string
	Length   int
	Expected string
}

func TestTruncatedStrings(t *testing.T) {

	var tests = []*truncateFixture{
		{Test: "dependabot/loves/shashes", Length: 6, Expected: "depend"},
		{Test: "no-more-hyphens", Length: 6, Expected: "no-mor"},
		{Test: "no_seperation-any more", Length: 6, Expected: "no_sep"},
		{Test: "not just letters! 123 () ? > * & . \\ #", Length: 6, Expected: "not ju"},
		{Test: "! 123 () ? > * & . \\ # not just letters", Length: 3, Expected: "! 1"},
		{Test: "exactlythesame", Length: 6, Expected: "exactl"},
		{Test: "", Length: 6, Expected: ""},
		{Test: "---massive, extremely long, infeasbile and quite frankly ridiculus length of a string to use", Length: 3, Expected: "---"},
	}

	for _, test := range tests {
		actual := Truncate(test.Test, test.Length)
		if actual != test.Expected {
			t.Errorf("error truncating string [%s], expected [%s] actual [%s]", test.Test, test.Expected, actual)
		}
	}

}

type safeFixture struct {
	Test     string
	Length   int
	Expected string
}

func TestSafeStrings(t *testing.T) {

	var tests = []*safeFixture{
		{Test: "dependabot/loves/shashes", Length: 6, Expected: "depend"},
		{Test: "no-more-hyphens", Length: 6, Expected: "nomore"},
		{Test: "no_seperation-any more", Length: 6, Expected: "nosepe"},
		{Test: "not just letters! 123 () ? > * & . \\ #", Length: 6, Expected: "notjus"},
		{Test: "! 123 () ? > * & . \\ # not just letters", Length: 3, Expected: "123"},
		{Test: "exactlythesame", Length: 6, Expected: "exactl"},
		{Test: "", Length: 6, Expected: ""},
		{Test: "---massive, extremely long, infeasbile and quite frankly ridiculus length of a string to use", Length: 3, Expected: "mas"},
	}

	for _, test := range tests {
		actual, _ := Safe(test.Test, test.Length)
		if actual != test.Expected {
			t.Errorf("error creating safe string string [%s], expected [%s] actual [%s]", test.Test, test.Expected, actual)
		}
	}

}
