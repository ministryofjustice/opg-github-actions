package semver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
)

// Regex for splitting and validating prerelease semver segments
var (
	PrereleaseSplitRegex string = `(?m)(?P<pprefix>.*)\.(?P<pbuild>[0-9]+)(?P<pextra>.*)`
	PrereleaseValidRegex string = `(?m)((?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
)

// Prerelease tracks the prerelease segment of the semver
// and provides functions to get specific chunks
// Uses regexes to form the data
type Prerelease struct {
	Complete string `json:"poriginal,omitempty"`
	Prefix   string `json:"pprefix,omitempty"`
	Build    string `json:"pbuild,omitempty"`
	Extra    string `json:"pextra,omitempty"`
	loaded   bool
}

// BuildNumber returns either uint64 or nil when it cant
// find a valid build number in the string
//
// The buildNumber is classed as the last block of integers
// within the string, so
//
//	'beta.0' => 0
//	'beta--test.2024.01.01.9' => 9
func (p *Prerelease) BuildNumber() *uint64 {
	if p == nil {
		return nil
	}
	if p.Build != "" {
		slog.Debug(fmt.Sprintf("prerelease BuildNumber: parsing [%s]", p.Build))
		return intOrNil(p.Build)
	}
	return nil
}

// Bump will try to increment the build number for
// this prerelease.
// If there is no buildnumber, then an error is
// returned
func (p *Prerelease) Bump() error {
	buildNumber := p.BuildNumber()
	if buildNumber == nil {
		return fmt.Errorf(ErrorPrereleaseNoBuildNumber)
	}
	i := *buildNumber
	i = i + 1
	p.Build = fmt.Sprintf("%d", i)
	p.Refresh()
	return nil
}

func (p *Prerelease) MustBump(prefix string) {
	err := p.Bump()
	if err != nil {
		if p.Prefix == "" {
			p.Prefix = prefix
		}
		p.Build = "0"
		p.Refresh()
	}
}

// Refresh updates the internal structure allowing
// fields to be changed but properties like 'Complete'
// to be updated
// Used after updating elements via bump functions etc
func (p *Prerelease) Refresh() (err error) {
	// create the map
	m, err := p.Map()
	if err != nil {
		return
	}
	// make it json
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return
	}
	// unmarshal against self
	err = json.Unmarshal(jsonStr, p)
	p.loaded = true
	return
}

// Map will generate a map of string using the regex group
// name as the indexes; allowing for easy updating via
// json marshalling
//
// On first run, uses the 'Complete' property which is set
// on creation, otherwise uses the String() func
func (p *Prerelease) Map() (m map[string]string, err error) {

	m = map[string]string{}
	str := p.Complete
	if p.loaded {
		str = p.String()
	}
	exp := regexp.MustCompile(PrereleaseSplitRegex)
	match := exp.FindStringSubmatch(str)
	slog.Debug(fmt.Sprintf("prerelease Map: converting [%s]...", p))

	if len(match) > 0 {
		slog.Debug("prerelease Map: found regex segment matches")
		for i, name := range exp.SubexpNames() {
			if name == "" {
				m["poriginal"] = match[i]
			} else {
				m[name] = match[i]
			}
		}
	} else {
		m["poriginal"] = str
		m["pprefix"] = str
	}

	return
}

// Return a string version of the struct
// - can return nil when the struct has no data
func (p *Prerelease) String() (str string) {
	var build = p.Build
	if p == nil {
		return ""
	}
	if len(build) > 0 {
		build = "." + build
	}
	if p.loaded {
		return fmt.Sprintf("%s%s%s", p.Prefix, build, p.Extra)
	}
	return p.Complete
}

// -- Validation

func ValidPrerelease(s string) (valid bool) {
	reg := regexp.MustCompile(PrereleaseValidRegex)
	valid = reg.MatchString(s)
	slog.Info(fmt.Sprintf("prerelease Valid: [%s] = [%t]", s, valid))
	return
}

// -- Creation

func NewPrerelease(str string) (p *Prerelease, err error) {
	if !ValidPrerelease(str) {
		return nil, fmt.Errorf(ErrorInvalidPrerelease, str)
	}
	p = &Prerelease{Complete: str}
	err = p.Refresh()

	return
}

func MustPrerelease(p *Prerelease, err error) *Prerelease {
	if err != nil {
		slog.Error("error creating prerelease:" + err.Error())
		return &Prerelease{}
	}
	return p
}
