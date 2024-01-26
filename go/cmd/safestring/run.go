package safestring

import (
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/safestrings"
)

func process(
	original string, suffix string, length int,
	conditionalMatch string, conditionalValue string) (output map[string]string, err error) {

	var (
		data safestrings.Safestring = safestrings.Safestring(original)
		full string                 = string(*data.Safe())
		safe string
	)
	output = map[string]string{}
	if conditionalMatch == original {
		full = conditionalValue
		safe = conditionalValue
	} else if suffix != "" && length > 0 {
		l := length - len(suffix)
		safe = fmt.Sprintf("%s%s", *data.SafeAndShort(l), suffix)
	} else if length > 0 {
		safe = string(*data.SafeAndShort(length))
	} else {
		safe = fmt.Sprintf("%s%s", full, suffix)
	}

	output = map[string]string{
		"original":          original,
		"suffix":            suffix,
		"length":            fmt.Sprintf("%v", length),
		"conditional_match": conditionalMatch,
		"conditional_value": conditionalValue,
		"safe":              safe,
		"full_length":       full,
	}
	return
}

func Run(args []string) (output map[string]string, err error) {
	slog.Info("[" + Name + "] Run")

	// parse command arguments
	FlagSet.Parse(args)
	err = parseArgs()
	if err != nil {
		return
	}
	output, err = process(*original, *suffix, *length, *conditionalMatch, *conditionalValue)
	return
}
