package strs

import (
	"regexp"
	"strings"
)

const allowedCharacters string = "[^a-zA-Z0-9]+"

// Clean converts a string to only include alpha numerics characters, so removes
// forward slashes and so on
func Clean(s string) (safe string) {
	var (
		err   error
		exp   *regexp.Regexp
		lower string = strings.ToLower(s)
	)
	if exp, err = regexp.Compile(allowedCharacters); err != nil {
		return
	}
	safe = exp.ReplaceAllString(lower, "")
	return
}

// Truncate the string to at most maxlength characters
func Truncate(s string, maxLength int) (short string) {
	if len(s) > maxLength {
		s = s[0:maxLength]
	}
	short = s
	return
}

// Safe cleans the string to only allowed characters and then truncates
// that result to desired length - returning both values
func Safe(s string, maxLength int) (safeAndShort string, safe string) {
	safe = Clean(s)
	safeAndShort = Truncate(safe, maxLength)
	return

}
