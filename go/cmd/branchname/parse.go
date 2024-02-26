package branchname

import (
	"errors"
	"fmt"
	"opg-github-actions/pkg/commonstrings"
	"os"
	"slices"
	"strings"
)

// parseArgs handles the validation and verification of the arguments for this command
func parseArgs() error {

	if *eventName == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "event-name")
	} else if !slices.Contains(eventNameChoices, *eventName) {
		return fmt.Errorf(commonstrings.ErrorArgumentInvalidChoice, "event-name", strings.Join(eventNameChoices, ", "), *eventName)
	}

	if *eventDataFile == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "event-data-file")
	} else if _, err := os.Stat(*eventDataFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(commonstrings.ErrorArgumentFileNotExist, "event-data-file", *eventDataFile)
	}

	return nil
}
