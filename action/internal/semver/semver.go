package semver

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
)

// Regex patterns for validation matching
// see:
//   - https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
//   - https://regex101.com/r/Ly7O1x/3/
const (
	regexPrefix        string = `^(?P<prefix>[A-Za-z]{0,1})`
	regexMajor         string = `(?P<major>0|[1-9]\d*)`
	regexMinor         string = `(?P<minor>0|[1-9]\d*)`
	regexPatch         string = `(?P<patch>0|[1-9]\d*)`
	regexPrerelease    string = `(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?`
	regexBuildMetadata string = `(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`
)

type Semver struct {
	GitRef          *plumbing.Reference `json:"-"` // the git reference this semver relates to if set
	Original        string              `json:"-"` // original short reference name
	Valid           bool                `json:"-"`
	Prefix          string              `json:"prefix"`
	Major           string              `json:"major"`
	Minor           string              `json:"minor"`
	Patch           string              `json:"patch"`
	PreleaseName    string              `json:"prerelease"`
	PrereleaseBuild string              `json:"prereleasebuild"`
	BuildMetadata   string              `json:"buildmetadata"`
}

func (self *Semver) String() string {
	var prelease string = ""
	var buildmeta string = ""

	if self.PreleaseName != "" {
		prelease += fmt.Sprintf("-%s", self.PreleaseName)
	}
	if self.PrereleaseBuild != "" {
		prelease += fmt.Sprintf(".%s", self.PrereleaseBuild)
	}
	if self.BuildMetadata != "" {
		buildmeta = fmt.Sprintf("+%s", self.BuildMetadata)
	}

	return fmt.Sprintf("%s%s.%s.%s%s%s",
		self.Prefix,
		self.Major,
		self.Minor,
		self.Patch,
		prelease,
		buildmeta,
	)
}

// Equal compares the core values of the Semver structs to check if they
// have the same values
func Equal(a, b *Semver) (equal bool) {
	equal = true

	if a.Valid != b.Valid {
		return false
	}
	if a.Prefix != b.Prefix {
		return false
	}
	if a.Major != b.Major {
		return false
	}
	if a.Minor != b.Minor {
		return false
	}
	if a.Patch != b.Patch {
		return false
	}
	if a.PreleaseName != b.PreleaseName {
		return false
	}
	if a.PrereleaseBuild != b.PrereleaseBuild {
		return false
	}
	if a.BuildMetadata != b.BuildMetadata {
		return false
	}
	return
}

// Regex returns the constructed regex pattern to use for semvar parsing
func Regex() string {
	return fmt.Sprintf(`(?m)%s%s\.%s\.%s%s%s$`,
		regexPrefix,
		regexMajor,
		regexMinor,
		regexPatch,
		regexPrerelease,
		regexBuildMetadata,
	)
}

// Valid uses regex to determine if the string passed
// is valid
func Valid(s string) (valid bool) {
	reg := regexp.MustCompile(Regex())
	valid = reg.MatchString(s)
	return
}

// convert takes original struct of T and by marshaling and then unmarshaling applied its
// content to destination R
func convert[T any, R any](source T, destination R) (err error) {
	var bytes []byte
	bytes, err = json.MarshalIndent(source, "", "  ")
	if err == nil {
		err = json.Unmarshal(bytes, destination)
	}
	return
}

// parse takes the semver and runs the regex against its original form to determine
// the consituient parts
func parse(s *Semver) {
	var (
		asMap   map[string]string = map[string]string{}
		matches []string          = []string{}
		str     string            = s.Original
		exp     *regexp.Regexp    = regexp.MustCompile(Regex())
	)
	s.Valid = Valid(s.Original)
	// return if its not valid
	if !s.Valid {
		fmt.Printf("[%s] is not valid", s.Original)
		return
	}

	matches = exp.FindStringSubmatch(str)
	for i, name := range exp.SubexpNames() {
		asMap[name] = matches[i]
	}
	convert(asMap, &s)

	// if prerelease is set, then split the prerelease and build number up
	if s.PreleaseName != "" && strings.LastIndex(s.PreleaseName, ".") > 0 {
		i := strings.LastIndex(s.PreleaseName, ".")
		s.PrereleaseBuild = s.PreleaseName[i+1:]
		s.PreleaseName = s.PreleaseName[:i]
	}
}

func FromString(ref string) (s *Semver) {
	s = &Semver{Original: ref}
	parse(s)
	return
}

func New(ref *plumbing.Reference) (s *Semver) {

	s = &Semver{
		GitRef:   ref,
		Original: ref.Name().Short(),
	}
	parse(s)

	return
}
