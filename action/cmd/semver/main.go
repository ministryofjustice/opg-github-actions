package main

import (
	"flag"
	"log/slog"
	"opg-github-actions/action/internal/logger"
	"os"
)

var (
	prereleaseSuffixLength int    = 12
	prereleaseSuffix       string = ""
)

const ErrMissingPrereleaseSuffix string = "error: --source argument not passed."

// init does the setup of args
func init() {
	flag.IntVar(&prereleaseSuffixLength, "prerelease_suffix_length", prereleaseSuffixLength, "Set the max length to use for tag suffixes")
	flag.StringVar(&prereleaseSuffix, "prerelease_suffix", prereleaseSuffix, "The value to convert into prerelease tag.")
}

func main() {
	var (
		err error
		res map[string]string
		lg  *slog.Logger = logger.New("INFO", "TEXT")
	)
	// process the arguments and fetch the fallback value from environment values
	flag.Parse()
	// run the command
	// res, err := Run(lg, original, maxLength)
	if err != nil {
		lg.Error(err.Error())
		os.Exit(1)
	}
	logger.Result(lg, res)

}
