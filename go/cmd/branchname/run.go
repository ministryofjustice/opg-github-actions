package branchname

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"opg-github-actions/pkg/safestrings"
	"strings"

	"github.com/google/go-github/v58/github"
)

func process(eventType string, content []byte) (output map[string]string, err error) {
	var (
		branch  string                 = ""
		headRef string                 = ""
		baseRef string                 = ""
		clean   safestrings.Safestring = ""
	)
	output = map[string]string{}

	if eventType == "pull_request" {
		event := &github.PullRequestEvent{}
		err = json.Unmarshal(content, &event)
		headRef = *event.PullRequest.Head.Ref
		baseRef = *event.PullRequest.Base.Ref
		branch = headRef

	} else if eventType == "push" {
		event := &github.PushEvent{}
		err = json.Unmarshal(content, &event)
		headRef = *event.After
		baseRef = *event.Before
		branch = strings.ReplaceAll(*event.Ref, "refs/heads/", "")
	} else {
		err = fmt.Errorf(ErrorIncorrectEventType, strings.Join(eventNameChoices, ", "), eventType)
		return
	}

	clean = safestrings.Safestring(branch)

	output = map[string]string{
		"event_name":     eventType,
		"branch_name":    branch,
		"head_commitish": headRef,
		"base_commitish": baseRef,
		"safe":           string(*clean.SafeAndShort(Length)),
		"full_length":    string(*clean.Safe()),
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

	content, err := ioutil.ReadFile(*eventDataFile)
	if err != nil {
		return
	}

	output, err = process(*eventName, content)
	output["event_data_file"] = *eventDataFile
	return
}
