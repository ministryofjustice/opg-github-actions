package createtag

import (
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/git/tags"
	"opg-github-actions/pkg/safestrings"
	"os"
	"slices"
	"strings"
)

func process(
	tagSet *tags.Tags, tagName string, allTags []string,
	commitish string,
	regen bool, push bool,
) (output map[string]string, err error) {

	var (
		originalTag    string = tagName
		createdSuccess bool   = false
	)

	exists := tagSet.ExistsIn(tagName, allTags)
	regenerated := false

	// decide if need to regenerate or fail if this exists
	if exists && !regen {
		err = fmt.Errorf(ErrorTagExists, tagName)
		return
	} else if exists && regen {
		regenerated = true
		tagName, err = tagSet.Unique(tagName, allTags)
	}

	// generate at hash from the string
	hash, err := tagSet.StringToHash(commitish)
	if err != nil {
		return
	}

	// create the new tag
	at, err := tagSet.CreateAt(tagName, hash)
	if err != nil {
		return
	}

	// check if created successfully
	foundAt, err := tags.RefStringify(tagSet.At(hash))
	if err != nil {
		return
	}
	// if the foundAt contains the new tagName at the same location
	// assumed its created ok
	if slices.Contains(foundAt, tagName) {
		createdSuccess = true
	}

	isTest := (len(os.Getenv("RUN_AS_TEST")) > 0)
	// if true, push to the remote
	if push && !isTest {
		slog.Info(fmt.Sprintf("pushing created tag [%s] to origin.", tagName))
		//err = tagSet.Push()
	}
	output = map[string]string{
		"test":          fmt.Sprintf("%t", isTest),
		"tags_at":       strings.Join(foundAt, ", "),
		"success":       fmt.Sprintf("%t", createdSuccess),
		"requested_tag": originalTag,
		"created_tag":   tagName,
		"created_at":    at.String(),
		"regenerated":   fmt.Sprintf("%t", regenerated),
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

	regen, err := safestrings.ToBool(*regenTag)
	if err != nil {
		return
	}

	push, err := safestrings.ToBool(*pushToRemote)
	if err != nil {
		return
	}

	tagset, err := tags.New(*repoDir)
	if err != nil {
		return
	}
	allTags, err := tagset.AllAsStrings()
	if err != nil {
		return
	}
	output, err = process(tagset, *tagName, allTags, *commitish, regen, push)

	output["directory"] = *repoDir
	output["push"] = fmt.Sprintf("%t", push)
	output["regen"] = fmt.Sprintf("%t", regen)
	output["commitish"] = *commitish

	return
}
