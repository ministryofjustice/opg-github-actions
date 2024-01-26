package latesttag

import (
	"errors"
	"fmt"
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
	p := safestrings.Safestring(*prerelease)
	if *prerelease == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "prerelease")
	} else if _, e := p.AsBool(); e != nil {
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
