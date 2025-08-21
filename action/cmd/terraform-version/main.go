package main

import (
	"flag"
	"fmt"
	"log/slog"
	"opg-github-actions/action/internal/logger"
	"os"
	"path/filepath"
	"regexp"
)

// Errors
const (
	ErrFileNotFound       string = "error: file [%s] was not found."
	ErrParsingVersionFile string = "error: failed to parse versions file via regex"
)

// fixed pattern to find values from the terraform file
const requiredVersionPattern string = `(?m).*required_version.*=.*"(?P<version>.*)"(.*)$`

// input arguments used via flag
var (
	directory string = ""
	file      string = "versions.tf"
)

// FileExists checks if the file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	// if there is an error, or the filepath doesnt exist, return false
	if err != nil || os.IsNotExist(err) {
		return false
	}
	// return false for directories - its not a file
	if info.IsDir() {
		return false
	}

	return true
}

// VersionData expects the content of the version file and runs regex against it
// to find the values
//
// Return example:
//
//	map[string]string{ "version": "1.1.0"}
func versionData(content string) (data map[string]string, err error) {

	data = map[string]string{}

	// compare the file content to the regex
	// 	- when matched, use the named groups to populate the output data
	//	- when not matched, set an error
	re := regexp.MustCompile(requiredVersionPattern)
	match := re.FindStringSubmatch(content)
	if len(match) > 0 {
		for i, name := range re.SubexpNames() {
			if name != "" {
				data[name] = match[i]
			}
		}
	} else {
		err = fmt.Errorf(ErrParsingVersionFile)
	}
	return
}

// Run processes the input and returns the values
// `result` should contain a `version` property - which we set a defauly version
func Run(lg *slog.Logger, dir string, fp string) (result map[string]string, err error) {
	var (
		content []byte
		file    string = filepath.Join(dir, fp)
	)
	result = map[string]string{"version": ""}

	if !fileExists(file) {
		err = fmt.Errorf(ErrFileNotFound, file)
		return
	}
	// read the file
	content, err = os.ReadFile(file)
	if err != nil {
		return
	}
	// pass into the check and return values
	result, err = versionData(string(content))

	return
}

// init does the setup of args
func init() {
	flag.StringVar(&directory, "directory", directory, "Directory path to operate from.")
	flag.StringVar(&file, "file", file, "The terraform file that contains the required_version property.")
}

func main() {
	var lg *slog.Logger = logger.New("info", "text")
	// process the arguments and fetch the fallback value from environment values
	flag.Parse()
	// run the command
	res, err := Run(lg, directory, file)
	if err != nil {
		lg.Error(err.Error())
		os.Exit(1)
	}
	logger.Result(lg, res)

}
