package createtag

import (
	"errors"
	"fmt"
	"opg-github-actions/pkg/commonstrings"
	"opg-github-actions/pkg/safestrings"
	"os"
)

// parseArgs is the internal function to handle the validation and verification
// When any conditions are not met, an error is returned stating what
func parseArgs() error {

	if *repoDir == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "repository")
	} else if _, err := os.Stat(*repoDir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(commonstrings.ErrorArgumentDirNotExist, "repository", *repoDir)
	}

	if *commitish == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "commitish")
	}

	if *tagName == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "tag-name")
	}

	if *regenTag == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "regen")
	} else if _, e := safestrings.ToBool(*regenTag); e != nil {
		return fmt.Errorf(commonstrings.ErrorArumentNotBoolean, "regen", *regenTag)
	}

	if *pushToRemote == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "push")
	} else if _, e := safestrings.ToBool(*pushToRemote); e != nil {
		return fmt.Errorf(commonstrings.ErrorArumentNotBoolean, "push", *pushToRemote)
	}

	return nil
}
