package latesttag

import (
	"errors"
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/commonstrings"
	"opg-github-actions/pkg/safestrings"
	"os"
)

func parseArgs() error {

	if *repositoryDir == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "repository")
	} else if _, err := os.Stat(*repositoryDir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(commonstrings.ErrorArgumentDirNotExist, "repository", *repositoryDir)
	}
	// prerelease
	if *prerelease == "" {
		slog.Warn("--prerelease was empty, setting to 'false'")
		s := "false"
		prerelease = &s
	} else if _, e := safestrings.ToBool(*prerelease); e != nil {
		return fmt.Errorf(commonstrings.ErrorArumentNotBoolean, "prerelease", *prerelease)
	}
	// prerelease_suffix
	if *prereleaseSuffix == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "prerelease-suffix")
	}
	// branch_name
	if *branchName == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "branch")
	}
	// release_branches
	if *releaseBranches == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "release-branches")
	}

	return nil
}
