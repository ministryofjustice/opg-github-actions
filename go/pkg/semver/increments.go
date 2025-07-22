package semver

import (
	"fmt"
	"strings"
)

type Increment string

const (
	Major Increment = "#major"
	Minor Increment = "#minor"
	Patch Increment = "#patch"
	Pre   Increment = "prerelease"
)

type IncrementCounters struct {
	Major int
	Minor int
	Patch int
}

func MustIncrement(inc Increment, err error) Increment {
	if err != nil {
		return Patch
	}
	return inc
}

func NewIncrement(str string) (inc Increment, err error) {
	switch str {
	case string(Major):
		inc = Major
	case string(Minor):
		inc = Minor
	case string(Patch):
		inc = Patch
	default:
		err = fmt.Errorf(ErrorInvalidIncrement, str)
	}
	return
}

// VersionBumpCount scans the strings (commit messages) and looks for triggers that
// would increment the semver (#major|#minor|#patch) and returns a counter for each type
//
// If no triggers are found then the counter that matches 'fallback' param will be
// incremented instead.
func VersionBumpCount(data []string, defaultIncrement Increment) (counter *IncrementCounters) {
	// we cant bump by prerelease
	if defaultIncrement == Pre {
		return nil
	}
	counter = &IncrementCounters{}
	// check within the strings
	for _, str := range data {
		if strings.Contains(str, string(Major)) {
			counter.Major++
		}
		if strings.Contains(str, string(Minor)) {
			counter.Minor++
		}
		if strings.Contains(str, string(Patch)) {
			counter.Patch++
		}
	}
	// deal with default element
	if counter.Major == 0 && counter.Minor == 0 && counter.Patch == 0 {
		switch defaultIncrement {
		case Major:
			counter.Major++
		case Minor:
			counter.Minor++
		case Patch:
			counter.Patch++
		}
	}
	return
}
