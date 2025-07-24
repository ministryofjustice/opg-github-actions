package main

import (
	"flag"
	"fmt"
	"log/slog"
	"opg-github-actions/action/internal/logger"
	"opg-github-actions/action/internal/strs"
	"os"
	"strings"
)

var (
	maxLength int    = 12
	original  string = ""
)

const ErrMissingValues string = "error: --source argument not passed."

// Run processes the input and returns the values
func Run(lg *slog.Logger, source string, length int) (result map[string]string, err error) {
	var (
		safe         string
		safeAndShort string
	)
	if source == "" {
		err = fmt.Errorf(ErrMissingValues)
		return
	}
	// remove any head references for fully formed branch values
	source = strings.TrimPrefix(source, "refs/heads/")

	safeAndShort, safe = strs.Safe(source, length)

	result = map[string]string{
		"branch_name": source,
		"safe":        safeAndShort,
		"full_length": safe,
	}

	return
}

// init does the setup of args
func init() {
	flag.IntVar(&maxLength, "length", maxLength, "Set the max length of the safe string to return")
	flag.StringVar(&original, "source", original, "The value to convert into a safe branch name.")
}

func main() {
	var lg *slog.Logger = logger.New("INFO", "TEXT")
	// process the arguments and fetch the fallback value from environment values
	flag.Parse()
	// run the command
	res, err := Run(lg, original, maxLength)
	if err != nil {
		lg.Error(err.Error())
		os.Exit(1)
	}
	logger.Result(lg, res)

}
