/*
branchname uses the event data passed to determine the correct short branch name.

Both `--event_name` and `--event_data_file` aret typically fed from github event variables, but can be any valid option.

Returns an alphanumeric only version of the branch name as 'full_length' and a shorterned (12 chars) version as 'safe', which is intended to be used for creating tags.

Will also return 'source_commitish' and 'destination_commitish' for comparison. For a pull request these are head & base barnches, for push these are before and after refs.

Usage:

	branch-name [flags]

The flags are:

	--event-name			(required, options: 'pull_request', 'push')
		Name of the event type that has happened, normally prefilled via github.event_name
		Determines which information is used from the `--event-data-file`.
	--event-data-file			(required)
		File path to a json file that contains all the current event data.
		Prefilled from github.event_path.
*/
package branchname

import "flag"

var (
	Name          = "branch-name"                           // Command name
	FlagSet       = flag.NewFlagSet(Name, flag.ExitOnError) // Argument group
	Length        = 12                                      // Max length
	eventName     = FlagSet.String("event-name", "", "Name of the event: [pull_request|push]")
	eventDataFile = FlagSet.String("event-data-file", "", "File where event environment data is stored")
)

var eventNameChoices = []string{"pull_request", "push"}

var ErrorIncorrectEventType string = "error: expected event to be [%s], recieved [%s]"
