package semver

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/maruel/natural"
)

type SortOrder bool

const (
	SORT_ASC  SortOrder = true
	SORT_DESC SortOrder = false
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

// String returns the string format of a semver
func (self *Semver) String() string {
	var prelease string = ""
	var buildmeta string = ""

	if self == nil {
		return ""
	}

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

func Strings(versions []*Semver) (strs []string) {
	strs = []string{}

	for _, v := range versions {
		if v != nil {
			strs = append(strs, v.String())
		}
	}
	return
}

// Sort orders semvers in order based on their .String() values
//
// Removes duplicates and invalid semvers
func Sort(versions []*Semver, order SortOrder) (sorted []*Semver) {
	var toSort []string = Strings(versions)

	sorted = []*Semver{}
	// sort & remove duplicates
	slices.Sort(toSort)
	toSort = slices.Compact(toSort)
	// change sort order to requested
	if order == SORT_DESC {
		sort.Sort(sort.Reverse(natural.StringSlice(toSort)))
	} else {
		sort.Sort(natural.StringSlice(toSort))
	}

	// loop over all the sorted version and add the first semver that matches
	// breaking the inner loop when one is found to avoid duplicates
	for _, key := range toSort {
		for _, sem := range versions {
			if key == sem.String() {
				sorted = append(sorted, sem)
				break
			}
		}
	}

	return
}

// parse takes the semver and runs the regex against its original form to determine
// the consituient parts
func parse(s *Semver) (err error) {
	var (
		asMap   map[string]string = map[string]string{}
		matches []string          = []string{}
		str     string            = s.Original
		exp     *regexp.Regexp    = regexp.MustCompile(Regex())
	)

	matches = exp.FindStringSubmatch(str)
	for i, name := range exp.SubexpNames() {
		asMap[name] = matches[i]
	}

	if err = convert(asMap, &s); err != nil {
		return
	}

	// if prerelease is set, then split the prerelease and build number up
	if s.PreleaseName != "" && strings.LastIndex(s.PreleaseName, ".") > 0 {
		i := strings.LastIndex(s.PreleaseName, ".")
		s.PrereleaseBuild = s.PreleaseName[i+1:]
		s.PreleaseName = s.PreleaseName[:i]
	}
	return
}

func New(ref *plumbing.Reference) (s *Semver) {

	s = &Semver{
		GitRef:   ref,
		Original: ref.Name().Short(),
		Valid:    true,
	}

	if !Valid(s.Original) {
		return nil
	}

	if err := parse(s); err != nil {
		return nil
	}

	return
}

func FromString(ref string) (s *Semver) {
	s = &Semver{
		Original: ref,
		Valid:    true,
	}

	if !Valid(s.Original) {
		return nil
	}

	if err := parse(s); err != nil {
		return nil
	}

	return
}

func FromStrings(refs ...string) (semvers []*Semver) {
	semvers = []*Semver{}

	for _, ref := range refs {
		if sv := FromString(ref); sv != nil {
			semvers = append(semvers, sv)
		}
	}
	return
}

func FromGitRefs(refs ...*plumbing.Reference) (semvers []*Semver) {
	semvers = []*Semver{}

	for _, ref := range refs {
		if sv := New(ref); sv != nil {
			semvers = append(semvers, sv)
		}
	}
	return
}
