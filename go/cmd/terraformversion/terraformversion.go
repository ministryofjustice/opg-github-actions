/*
terraformversion parses the file passed to determine the terraform-version.

The 'versions-file' is presumed to contain the terraform block with 'required_version'
present.

If 'simple' is set, it presumes the content of the 'versions-file' is only a string
relating to the version and nothing else.

Usage:

	terraform-version [flags]

The flags are:

	--directory			(required, default: ./)
	--versions-file		(required, default: ./versions.tf)
	--simple
*/
package terraformversion

import "flag"

var (
	Name         = "terraform-version"                     // Command name
	FlagSet      = flag.NewFlagSet(Name, flag.ExitOnError) // Argument group
	directory    = FlagSet.String("directory", "", "Directory to look for versions.tf")
	versionsFile = FlagSet.String("versions-file", "versions.tf", "Name of the versions.tf file")
	isSimple     = FlagSet.String("simple", "", "When set, presumes a simple file with just a version number string")
)
