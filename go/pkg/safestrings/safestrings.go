// Package safestrings provides a couple of simple wrappers to make
// a string git tag / branch name friendly, removing special
// characters and returning only lowercase alphaunmerics
//
// Also provides a 'true' => true bool helper for param
// parsing
package safestrings

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	safeStringPattern = "[^a-zA-Z0-9]+"
)

type Safestring string

func (s *Safestring) Safe() *Safestring {
	source := strings.ToLower(string(*s))
	reg, _ := regexp.Compile(safeStringPattern)
	clean := reg.ReplaceAllString(source, "")
	safe := Safestring(clean)
	return &safe
}

func (s *Safestring) Short(maxLength int) *Safestring {
	source := string(*s)
	if len(source) > maxLength {
		runes := []rune(source)
		source = string(runes[0:maxLength])
	}
	short := Safestring(source)
	return &short
}

func (s *Safestring) SafeAndShort(maxLength int) *Safestring {
	return s.Safe().Short(maxLength)
}

func (s *Safestring) AsBool() (b bool, err error) {
	return strconv.ParseBool(string(*s))
}

func ToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
