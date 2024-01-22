package terraformversion

import (
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/safestrings"
	"os"
	"path/filepath"
	"regexp"
)

var RequiredVersionPattern string = `(?m).*required_version.*=.*"(?P<version>.*)"(.*)$`

func process(versionsFileContent string, simple bool) (output map[string]string, err error) {

	output = map[string]string{}
	// for a simple file, return the content directly
	if simple {
		output["version"] = versionsFileContent
	} else {
		// otherwise, compare the file content to the regex
		// 	- when matched, use the named groups to populate the output data
		//	- when not matched, set an error
		re := regexp.MustCompile(RequiredVersionPattern)
		match := re.FindStringSubmatch(versionsFileContent)
		if len(match) > 0 {
			for i, name := range re.SubexpNames() {
				if name != "" {
					output[name] = match[i]
				}
			}
		} else {
			err = fmt.Errorf(ErrorParsingVersionFile)
		}
	}

	return
}

func Run(args []string) (output map[string]string, err error) {
	slog.Info("[" + Name + "] Run")
	FlagSet.Parse(args)

	// parse command arguments
	err = parseArgs()
	if err != nil {
		return
	}

	simple := safestrings.Safestring(*isSimple)
	simpleBool, err := simple.AsBool()
	if err != nil {
		return
	}

	filePath := filepath.Join(*directory, *versionsFile)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	output, err = process(string(content), simpleBool)

	return
}
