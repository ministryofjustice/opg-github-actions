package logger

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	workspace string = "GITHUB_WORKSPACE"
	githubout string = "GITHUB_OUTPUT"
)

// Result provides a wrapper to echo results information out to the stdout for github
func Result(logger *slog.Logger, results map[string]string) {
	var outFile string
	var f *os.File
	// github output via os.environ['GITHUB_OUTPUT']
	if os.Getenv(workspace) != "" {
		logger.Info("github workspace found, writting to GITHUB_OUTPUT")
		outFile = os.Getenv(githubout)
		f, _ = os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
	} else {
		logger.Info("github workspace not found...")
	}

	for k, v := range results {
		str := fmt.Sprintf("%s=%s\n", k, v)
		fmt.Printf(str)
		if outFile != "" && f != nil {
			logger.Debug("writting to GITHUB_OUTPUT")
			f.WriteString(str)
		}
	}
}
