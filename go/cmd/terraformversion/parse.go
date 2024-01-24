package terraformversion

import (
	"errors"
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/commonstrings"
	"os"
	"path/filepath"
)

// parseArgs handles the validation and verification of the arguments for this command
func parseArgs() error {

	if *directory == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "directory")
	} else if _, err := os.Stat(*directory); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(commonstrings.ErrorArgumentFileNotExist, "directory", *directory)
	}
	d := *directory
	path := filepath.Join(d, *versionsFile)
	slog.Debug(fmt.Sprintf("args: checking path: [%s][%s]=>[%s]", *directory, *versionsFile, path))

	if *versionsFile == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "versions-file")
	} else if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(commonstrings.ErrorArgumentFileNotExist, "versions-file", path)
	}

	return nil
}
