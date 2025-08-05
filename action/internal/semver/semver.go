package semver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/maruel/natural"
)

type SortOrder string

const (
	SORT_ASC  SortOrder = "asc"
	SORT_DESC SortOrder = "desc"
)

type Increment string

func (self Increment) Stringy() string {
	return fmt.Sprintf("#%s", strings.ToLower(string(self)))
}

const (
	NO_BUMP Increment = "none"
	MAJOR   Increment = "major"
	MINOR   Increment = "minor"
	PATCH   Increment = "patch"
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
	GitRef          *plumbing.Reference `json:"-"`      // the git reference this semver relates to if set
	Original        string              `json:"o"`      // original short reference name
	Valid           bool                `json:"v"`      // if the semver is valid or not
	Prefix          string              `json:"prefix"` // prefix part of the semver (typically `v`)
	Major           string              `json:"major"`
	Minor           string              `json:"minor"`
	Patch           string              `json:"patch"`
	PreleaseName    string              `json:"prerelease"`
	PrereleaseBuild string              `json:"prereleasebuild"`
	BuildMetadata   string              `json:"buildmetadata"`
}

func (self *Semver) IsPrerelease() bool {
	return (self.PreleaseName != "")
}
func (self *Semver) IsRelease() bool {
	return (self.PreleaseName == "")
}

// String returns a json friendly version of self
func (self *Semver) String() string {
	return self.Stringy(true)
}

// Stringy is used instead of String in places where we want to toggle
// the inclusion of the semver prefix
func (self *Semver) Stringy(includePrefix bool) string {

	return format(self, &strOpts{
		PrereleaseName:  true,
		PrereleaseBuild: true,
		BuildMetadata:   true,
		Prefix:          includePrefix,
	})

}

type strOpts struct {
	Prefix          bool
	PrereleaseName  bool
	PrereleaseBuild bool
	BuildMetadata   bool
}

// format used internally to format the semver string in various ways
// for partial or complete matching in various filters
func format(s *Semver, opts *strOpts) (str string) {
	var (
		prefix          string = ""
		version         string = ""
		prereleaseName  string = ""
		prereleaseBuild string = ""
		prerelease      string = ""
		buildMetadata   string = ""
	)
	if s == nil {
		return ""
	}
	version = fmt.Sprintf("%s.%s.%s", s.Major, s.Minor, s.Patch)

	if opts.Prefix {
		prefix = s.Prefix
	}
	if opts.PrereleaseName && s.PreleaseName != "" {
		prereleaseName = fmt.Sprintf("%s", s.PreleaseName)
	}
	if opts.PrereleaseBuild && s.PrereleaseBuild != "" {
		prereleaseBuild = fmt.Sprintf(".%s", s.PrereleaseBuild)
	}
	if prereleaseName != "" || prereleaseBuild != "" {
		prerelease = fmt.Sprintf("-%s%s", s.PreleaseName, prereleaseBuild)
	}
	if opts.BuildMetadata && s.BuildMetadata != "" {
		buildMetadata = fmt.Sprintf("+%s", s.BuildMetadata)
	}

	str = fmt.Sprintf("%s%s%s%s",
		prefix,
		version,
		prerelease,
		buildMetadata,
	)

	return
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

// Strings takes semvers and returns string version of those; generally used for sorting
func Strings(versions []*Semver, prefix bool) (strs []string) {
	strs = []string{}
	for _, v := range versions {
		if v != nil {
			strs = append(strs, v.Stringy(prefix))
		}
	}
	return
}

// Sort orders semvers in order based on their .String() values
//
// Removes duplicates and invalid semvers
func Sort(lg *slog.Logger, versions []*Semver, order SortOrder, prefixes bool) (sorted []*Semver) {
	var toSort []string
	sorted = []*Semver{}

	lg = lg.With("operation", "Sort", "order", string(order))

	lg.Debug("getting string versions for sorting ... ")
	toSort = Strings(versions, prefixes)

	lg.Debug("sorting and removing duplicates ... ")
	// sort & remove duplicates
	slices.Sort(toSort)
	toSort = slices.Compact(toSort)
	// change sort order to requested
	if order == SORT_DESC {
		sort.Sort(sort.Reverse(natural.StringSlice(toSort)))
	} else {
		sort.Sort(natural.StringSlice(toSort))
	}

	lg.Debug("adding values based on sort order and ignoring duplicates ... ")
	// loop over all the sorted version and add the first semver that matches
	// breaking the inner loop when one is found to avoid duplicates
	for _, key := range toSort {
		for _, sem := range versions {
			if key == sem.Stringy(prefixes) {
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

// FromString uses the value passed as the original value
// which is then split into segments (major, minor, patch
// etc)
//
// If the ref name is invalid or the tag does not parse correctly
// then nil is returned
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

// FromStrings uses the values passed as the original values
// which is then split into segments (major, minor, patch
// etc)
//
// If the name is invalid or does not parse correctly
// then it is skipped and not returned
func FromStrings(refs ...string) (semvers []*Semver, err error) {
	semvers = []*Semver{}

	for _, ref := range refs {
		if sv := FromString(ref); sv != nil {
			semvers = append(semvers, sv)
		}
	}
	return
}

// FromGitRefs uses the git tags / references and generates a semver
// for each using the ref.Name.Short() (refs/tags/v4.1.1 => v4.1.1) as
// the starting point and the parses out the segments (major, minor,
// patch etc)
//
// If the ref name is invalid or the tag does not parse correctly
// then its skipped
//
// A `Must` style pattern to chain with `tags.All(dir)`
func FromGitRefs(refs []*plumbing.Reference) (semvers []*Semver, err error) {
	semvers = []*Semver{}

	for _, ref := range refs {
		if sv := New(ref); sv != nil {
			semvers = append(semvers, sv)
		}
	}
	return
}

// GetPrereleases filters down the list of semvers to only prereleases
func GetPrereleases(existing []*Semver) (prereleases []*Semver) {
	prereleases = []*Semver{}

	for _, exists := range existing {
		if exists != nil && exists.IsPrerelease() {
			prereleases = append(prereleases, exists)
		}
	}

	return
}

// GetReleases filters down the list of semvers to only releases
func GetReleases(existing []*Semver) (releases []*Semver) {
	releases = []*Semver{}

	for _, exists := range existing {
		if exists != nil && exists.IsRelease() {
			releases = append(releases, exists)
		}
	}

	return
}

// inc is internal function that handles incrementing a string by 1 and
// returning string value (used to bump along semver segments)
func inc(s string) (i string) {
	n := atoi(s) + 1
	i = strconv.Itoa(n)

	return
}

func atoi(s string) (i int) {
	i = 0
	if n, err := strconv.Atoi(s); err == nil {
		i = n
	}
	return
}

// Release runs over the existing Semvers, finds that largest (naturally sorted) version
// and increments that value by bump.
//
// If `bump` is not one of `MAJOR`, `MINOR`, `PATCH` then the semver is not updated.
// By using `NONE` the last release version is returned instead.
//
// If no releases are found, `0.0.0` is used instead.
func Release(lg *slog.Logger, existing []*Semver, bump Increment) (next *Semver) {
	var (
		last     *Semver
		releases []*Semver
	)
	lg = lg.With("operation", "Release", "bump", string(bump))
	lg.Debug("sorting releases ... ")
	// get releases only and sort them descending order
	releases = Sort(lg, GetReleases(existing), SORT_DESC, false)
	// If there are no releases, then use 0 as base
	// Otherwise, use the last release
	if len(releases) == 0 {
		last = FromString("0.0.0")
	} else {
		last = releases[0]
	}

	lg.Debug("last release ... ", "last", last.Stringy(true))
	next = last
	// if we are bumping just patch, update and return
	// if bump is minor or major then patch is reset
	// if bump is major, then reset patch & minor
	if bump == PATCH {
		next.Patch = inc(next.Patch)
	} else if bump == MINOR {
		next.Patch = "0"
		next.Minor = inc(next.Minor)
	} else if bump == MAJOR {
		next.Patch = "0"
		next.Minor = "0"
		next.Major = inc(next.Major)
	}
	lg.Debug("next release ... ", "next", next.Stringy(true))
	return
}

// Prerelease looks at all the existing semvers, finds that last release and uses that with the
// suffix value passed to generate a prerealease version with a build counter (v1.0.1-suffix.1)
//
// It finds the last release by calling Release and using that as the base line.
//
// If gets all prereleases from the existing set and matches those with the same
// `MAJOR.MINOR.PATCH-suffix.buildNumber` pattern, then increments the buildNumber
func Prerelease(lg *slog.Logger, existing []*Semver, bump Increment, suffix string) (next *Semver) {
	var (
		partial string
		build   int      = 0
		opts    *strOpts = &strOpts{
			PrereleaseName:  true,
			Prefix:          false,
			PrereleaseBuild: false,
			BuildMetadata:   false,
		}
	)
	lg = lg.With("operation", "Prerelease", "bump", string(bump), "suffix", suffix)
	// get the last release by passing along -1, so an increment is never triggered
	next = Release(lg, existing, bump)
	// now setup the prefixes for this being a prerelease
	next.PreleaseName = suffix
	next.PrereleaseBuild = "0"
	lg.Debug("release semver", "release", next.Stringy(true))
	// grab most of the semver signature, ignore the build number & build metadata
	partial = format(next, opts)
	lg.Debug("partial string to match ... ", "partial", partial)
	// now look for any prereleases that have the same partial signature
	// and if we find them, we use the latest values and increment the build number
	pres := GetPrereleases(existing)
	for _, pre := range pres {
		// grab most of this semver signature, ignore the build number & build metadata
		compare := format(pre, opts)
		lg.Debug("compare == partial ... ", "partial", partial, "compare", compare)
		if compare == partial {
			// if the build number is higher than one we've seen before, use this prerelease as a bench mark
			// and track the build number of compare
			if bn := atoi(pre.PrereleaseBuild); bn >= build {
				next = pre
				build = bn
			}
		}

	}
	next.PrereleaseBuild = inc(next.PrereleaseBuild)
	lg.Debug("prerelease generated ...  ", "next", next.Stringy(true))
	return
}

// GetBumpFromCommits scans the commit messages and looks for triggers that
// would increment the semver (#major|#minor|#patch) and returns a counter for each type
//
// If no triggers are found then the counter that matches 'fallback' param will be
// incremented instead.
//
// Calls `GetBump` underneath
func GetBumpFromCommits(lg *slog.Logger, commits []*object.Commit, defaultBump Increment) (bump Increment) {

	lg.Debug("generating message strings from commits", "operation", "GetBumpFromCommits", "defaultBump", string(defaultBump))

	var messages = []string{}
	for _, commit := range commits {
		messages = append(messages, commit.Message)
	}

	bump = GetBump(lg, messages, defaultBump)
	return
}

// GetBump scans the strings (commit messages) and looks for triggers that
// would increment the semver (#major|#minor|#patch) and returns a counter for each type
//
// If no triggers are found then the counter that matches 'fallback' param will be
// incremented instead.
func GetBump(lg *slog.Logger, commitMessages []string, defaultBump Increment) (bump Increment) {
	lg = lg.With("operation", "GetBump", "defaultBump", string(defaultBump))

	bump = ""

	// if there are any commits, then should at lease be a patch bump
	if len(commitMessages) > 0 {
		lg.Debug("commits were found, so setting base increment to patch")
		bump = PATCH
	}

	lg.Debug("checking commit messages ... ")
	for _, content := range commitMessages {
		// if we find any major, then return
		// if the bump isnt a major, and we find a minor, then set to minor
		// if the bump isnt major or minor and we find patch, set to patch
		if strings.Contains(content, MAJOR.Stringy()) {
			bump = MAJOR
			lg.Debug("found major ... ")
			return
		} else if bump != MAJOR && strings.Contains(content, MINOR.Stringy()) {
			lg.Debug("found minor ... ")
			bump = MINOR
		} else if bump != MAJOR && bump != MINOR && strings.Contains(content, PATCH.Stringy()) {
			lg.Debug("found patch ... ")
			bump = PATCH
		}
	}

	if bump == "" {
		bump = defaultBump
	}

	lg.Debug("calculated bump", "dump", string(bump))
	return
}

// New uses the git tag / reference and generates a semver struct
// using the ref.Name.Short() (refs/tags/v4.1.1 => v4.1.1) as the
// starting point and the parses out the segments (major, minor,
// patch etc)
//
// If the ref name is invalid or the tag does not parse correctly
// then nil is returned
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
