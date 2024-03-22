package nexttag

import (
	"errors"
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/commonstrings"
	"opg-github-actions/pkg/safestrings"
	"os"
)

// parseArgs is the internal function to handle the validation and verification
// of the command arguments passed in to 'next-tag'.
// When any conditions are not met, an error is returned stating what
func parseArgs() error {

	if *repoDir == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "repository")
	} else if _, err := os.Stat(*repoDir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(commonstrings.ErrorArgumentDirNotExist, "repository", *repoDir)
	}

	if *headRef == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "head")
	}

	if *baseRef == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "base")
	}

	if *prerelease == "" {
		slog.Warn("--prerelease was empty, setting to 'false'")
		s := "false"
		prerelease = &s
	} else if _, e := safestrings.ToBool(*prerelease); e != nil {
		return fmt.Errorf(commonstrings.ErrorArumentNotBoolean, "prerelease", *prerelease)
	}

	if *withV != "" {
		if _, wve := safestrings.ToBool(*withV); wve != nil {
			return fmt.Errorf(commonstrings.ErrorArumentNotBoolean, "with-v", *withV)
		}
	}

	if *prereleaseSuffix == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "prerelease-suffix")
	}

	if *defaultBump == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "default-bump")
	} else if *defaultBump != "major" && *defaultBump != "minor" && *defaultBump != "patch" {
		return fmt.Errorf(commonstrings.ErrorArgumentInvalidChoice, "default-bump", "major,minor,patch", *defaultBump)
	}
	return nil
}
