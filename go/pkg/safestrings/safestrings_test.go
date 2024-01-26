package safestrings

import "testing"

type safeFixture struct {
	Test     string
	Expected string
}

func TestSafe(t *testing.T) {

	fixtures := []safeFixture{
		{Test: "dependabot/loves/shashes", Expected: "dependabotlovesshashes"},
		{Test: "no-more-hyphens", Expected: "nomorehyphens"},
		{Test: "no_serperation-any more", Expected: "noserperationanymore"},
		{Test: "not just letters! 123 () ? > * & . \\ #", Expected: "notjustletters123"},
		{Test: "exactlythesame", Expected: "exactlythesame"},
		{Test: "", Expected: ""},
	}

	for _, fixture := range fixtures {
		s := Safestring(fixture.Test)
		res := s.Safe()
		actual := string(*res)
		if fixture.Expected != actual {
			t.Errorf("Expected [%s] Actual [%s]", fixture.Expected, actual)
		}
	}

}

func TestShort(t *testing.T) {

	fixtures := []safeFixture{
		{Test: "dependabot/loves/shashes", Expected: "depend"},
		{Test: "no-more-hyphens", Expected: "no-mor"},
		{Test: "no_seperation-any more", Expected: "no_sep"},
		{Test: "not just letters! 123 () ? > * & . \\ #", Expected: "not ju"},
		{Test: "exactlythesame", Expected: "exactl"},
		{Test: "", Expected: ""},
	}

	for _, fixture := range fixtures {
		s := Safestring(fixture.Test)
		res := s.Short(6)
		actual := string(*res)
		if fixture.Expected != actual {
			t.Errorf("Expected [%s] Actual [%s]", fixture.Expected, actual)
		}
	}

}

func TestSafeShort(t *testing.T) {

	fixtures := []safeFixture{
		{Test: "dependabot/loves/shashes", Expected: "depend"},
		{Test: "no-more-hyphens", Expected: "nomore"},
		{Test: "no_seperation-any more", Expected: "nosepe"},
		{Test: "not just letters! 123 () ? > * & . \\ #", Expected: "notjus"},
		{Test: "exactlythesame", Expected: "exactl"},
		{Test: "", Expected: ""},
	}

	for _, fixture := range fixtures {
		s := Safestring(fixture.Test)
		res := s.SafeAndShort(6)
		actual := string(*res)
		if fixture.Expected != actual {
			t.Errorf("Expected [%s] Actual [%s]", fixture.Expected, actual)
		}
	}

}
