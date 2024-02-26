package safestring

import (
	"fmt"
	"opg-github-actions/pkg/commonstrings"
)

// parseArgs handles the validation and verification of the arguments for this command
func parseArgs() error {

	if *original == "" {
		return fmt.Errorf(commonstrings.ErrorArgumentMissing, "string")
	}

	return nil
}
