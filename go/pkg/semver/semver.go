package semver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
)

// Regex patterns for validation matching
// see:
//   - https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
//   - https://regex101.com/r/Ly7O1x/3/
var (
	regexOperator         string = `(?m)`
	officialRegex         string = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
	prefixRegexGroup      string = `^(?P<prefix>[A-Za-z]{0,1})`
	SemverWithPrefixRegex string = fmt.Sprintf(`%s%s%s`, regexOperator, prefixRegexGroup, officialRegex[1:]) // generate the regex from parts
)

type semver struct {
	Complete      string `json:"original,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	MajorStr      string `json:"major,omitempty"`
	MinorStr      string `json:"minor,omitempty"`
	PatchStr      string `json:"patch,omitempty"`
	PrereleaseStr string `json:"prerelease"`
	Buildmetadata string `json:"buildmetadata,omitempty"`
	prerelease    *Prerelease
}

// Prerelease will return either the current prerelease
// item or a new one based on the PrereleaseStr
func (sv *semver) Prerelease() *Prerelease {
	if sv.prerelease != nil {
		return sv.prerelease
	}
	return MustPrerelease(NewPrerelease(sv.PrereleaseStr))
}

// String converts struct back to a single string
func (sv *semver) String() string {
	if sv == nil {
		return ""
	}
	pre := sv.Prerelease().String()
	meta := sv.Buildmetadata

	if len(pre) > 0 {
		pre = "-" + pre
	}
	if len(meta) > 0 {
		meta = "+" + meta
	}
	return fmt.Sprintf(
		"%s%s.%s.%s%s%s",
		sv.Prefix, sv.MajorStr, sv.MinorStr, sv.PatchStr,
		pre, meta,
	)
}

// Semver is the main struct used in this package.
// Represents a semver formatted string that expands on the
// officalt regex to allow the inclusion of a single character
// prefix
//
// The internal `data` property stores all the raw semver
// string data. This prviate semver type is configured
// to allow json marshaling directly from the regex split,
// making conversion between string, struct etc cleaner.
//
// As json marshalling requires the properties to be visible
// make it a private part of this struct so we dont expose
// the raw strings outside of the package to avoid direct
// changes
type Semver struct {
	data   *semver
	loaded bool
}

// -- getters and setters
// Prefix returns the semver prefix if one has been
// set. This is typically a single char like 'v'
func (s *Semver) Prefix() string {
	return s.data.Prefix
}

// SetPrefix allows adding of a single char
// prefix to this semver entity, this would
// typically be a 'v'
func (s *Semver) SetPrefix(prefix rune) {
	s.data.Prefix = string(prefix)
	s.Refresh()
}

// RemovePrefix clears the prefix string
// so it wont have a 'v' at the start etc
func (s *Semver) RemovePrefix() {
	s.data.Prefix = ""
	s.Refresh()
}

// Major returns unsigned int (or nil) version
// of the string
func (s *Semver) Major() *uint64 {
	return intOrNil(s.data.MajorStr)
}

// SetMajor updates the MajorStr to this int
// and then triggers a refresh of the data
// object
func (s *Semver) SetMajor(i uint64) {
	s.data.MajorStr = fmt.Sprintf("%d", i)
	s.Refresh()
}

// BumpMajor increments the Major segment by 1
func (s *Semver) BumpMajor() {
	i := *s.Major()
	s.SetMajor(i + 1)
}

// Minor returns int or nil version of the
// string
func (s *Semver) Minor() *uint64 {
	return intOrNil(s.data.MinorStr)
}

// SetMinor updates the MinorStr to this int
// and then triggers a refresh of the data
// object
func (s *Semver) SetMinor(i uint64) {
	s.data.MinorStr = fmt.Sprintf("%d", i)
	s.Refresh()
}

// BumpMinor increments the Minor segment by 1
func (s *Semver) BumpMinor() {
	i := *s.Minor()
	s.SetMinor(i + 1)
}

// Patch gets the int | nil version of the
// patch string
func (s *Semver) Patch() *uint64 {
	return intOrNil(s.data.PatchStr)
}

// SetPatch updates the PatchStr to this int
// and then triggers a refresh of the data
// object
func (s *Semver) SetPatch(i uint64) {
	s.data.PatchStr = fmt.Sprintf("%d", i)
	s.Refresh()
}

// BumpPatch increments the Patch segment by 1
func (s *Semver) BumpPatch() {
	i := *s.Patch()
	s.SetPatch(i + 1)
}

// -- prerelease helpers

// IsPrerelease returns true|false depending on if
// the parsed data has any form of prerelease string
func (s *Semver) IsPrerelease() bool {
	return len(s.data.PrereleaseStr) > 0
}

// IsPrereleaseMatch checks the prerelease segment of
// the semver for a match against the passed string
//
// "1.0.0-beta.1"
//   - prerelease = "beta.1"
//   - prerelease prefix = "beta"
func (s *Semver) IsPrereleaseMatch(str string) bool {
	if len(s.data.PrereleaseStr) > 0 {
		return (s.data.Prerelease().Prefix == str)
	}
	return false
}

// Prerelrease simply returns the private prerelease item
func (s *Semver) Prerelease() *Prerelease {
	return s.data.Prerelease()
}

// SetPrerelease provides a way to full overwrite the
// prerelease data so can be removed or changed
func (s *Semver) SetPrerelease(str string) {
	s.data.PrereleaseStr = str
	s.Refresh()
}

// BumpPrerelease will try to bump the prerelease
// build number. If the prerelreease does not contain
// a build number an error is returned
func (s *Semver) BumpPrerelease() error {
	e := s.Prerelease().Bump()
	s.Refresh()
	return e
}

// MustBumpPrerelease calls the underpinning MustBump
// This will create a new prerelease based on the prefix
// with a build number of 0 if a build number cannot be
// found
func (s *Semver) MustBumpPrerelease(prefix string) {
	s.Prerelease().MustBump(prefix)
	s.Refresh()
}

// PrereleaseString returns the prerelease data
// as a formatted string
func (s *Semver) PrereleaseString() string {
	return s.data.Prerelease().String()
}

// PrerereleasePrefix returns only the first
// segment of the prerelease.
//
// If the semver was '1.0.0-beta.1', this would
// return 'beta;'
func (s *Semver) PrereleasePrefix() string {
	return s.data.Prerelease().Prefix
}

// SetPrereleasePrefix directly changes the prerelease
// prefix segment and updates all structures
//
// If your semver was '1.0.0-beta.01+b1' and you
// passed 'test' to this method, the updated
// semver would be '1.0.0-test.01+b1'
func (s *Semver) SetPrereleasePrefix(str string) {
	pre := s.data.prerelease
	pre.Prefix = str
	pre.Refresh()
	s.Refresh()
}

// BuildNumber gets the builfnumber from a prerelease
// string that matches:
//
//	"1.0.0-beta.01+b1"
func (s *Semver) BuildNumber() *uint64 {
	return s.data.Prerelease().BuildNumber()
}

// Refresh take a map of this object and update its
// internal data properties to match
//
// This is so when a version segment (major,minor,patch)
// or prerelease is changed it calls this to ensure
// all properties are insync
//
// Sets the private `loaded` flag to true to avoid
// cyclical logic
func (s *Semver) Refresh() (err error) {
	// create the map
	m, err := s.Map()
	if err != nil {
		return
	}
	// make it json
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return
	}
	// unmarshal against self
	err = json.Unmarshal(jsonStr, s.data)
	// regen the prerelease
	s.data.prerelease = MustPrerelease(NewPrerelease(s.data.PrereleaseStr))
	s.loaded = true
	return
}

// --- CONVERSIONS

// Map converts this semver to a map of strings.
//
// On the first run, this is uses the `Complete` property that
// is set as part of New function. Afterwards it will use
// the generated String() call which will providate latest info
//
// Map is created by splitting the source string based on the
// offical semver regex (with modification for prefix) and
// grouped accordingly. Those group names then become the keys
// to the map with the empty group being the complete matching
// string
func (s *Semver) Map() (m map[string]string, err error) {
	// if this is the first run, use the raw string
	// otherwise used the calculated one
	str := s.data.Complete
	if s.loaded {
		str = s.data.String()
	}
	exp := regexp.MustCompile(SemverWithPrefixRegex)
	match := exp.FindStringSubmatch(str)
	slog.Debug(fmt.Sprintf("semver Map: converting [%s] to a map...", s.String()))

	if len(match) > 0 {
		m = map[string]string{}
		slog.Debug("semver Map: found regex segment matches")
		for i, name := range exp.SubexpNames() {
			if name == "" {
				m["original"] = match[i]
			} else {
				m[name] = match[i]
			}
		}
	} else {
		err = fmt.Errorf(ErrorConversionAsMapNoMatch, str)
		slog.Error(err.Error())

	}
	return
}

// Returns a string version of this struct.
// Where possible, returns the String() of data
func (s *Semver) String() (str string) {
	if s != nil && s.data != nil && s.data.Complete != "" {
		return s.data.String()
	}
	return ""
}

// -- Validation

// Valid uses regex to determine if the string passed
// is valid
func Valid(s string) (valid bool) {
	reg := regexp.MustCompile(SemverWithPrefixRegex)
	valid = reg.MatchString(s)
	slog.Info(fmt.Sprintf("semver Valid: [%s] = [%t]", s, valid))
	return
}

// HasPrefix checks to see if the string passed contains
// a prefix (like 'v')
//
// Will check is the string is a valid semver, returns
// false if its not
//
// If an error occurs in creation or Map() conversion then
// this will return false
func HasPrefix(s string) bool {
	if !Valid(s) {
		return false
	}
	sm, err := New(s)
	if err != nil {
		return false
	}
	m, err := sm.Map()
	if p, ok := m["prefix"]; ok {
		return (len(p) > 0)
	}
	return false

}

// --- CREATION
// Default creates a new semver with "0.0.0"
func Default() (s *Semver) {
	s, _ = New(defaultSemver)
	return
}

// New will create a fresh semver if it can, or return an
// error.
//
// An error will be generated if 'str' does not pass
// validation regex
func New(str string) (s *Semver, err error) {
	if !Valid(str) {
		return nil, fmt.Errorf(ErrorInvalidSemver, str)
	}
	s = &Semver{data: &semver{Complete: str}}
	s.Refresh()
	return
}

// Must wraps the new, so if new fails it will return
// a semver via Default()
func Must(sv *Semver, err error) *Semver {
	if err != nil {
		return Default()
	}
	return sv
}
