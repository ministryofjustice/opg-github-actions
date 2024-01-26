/*
safestring uses the string passed to generate a tag safe version ('safe') that is alphanumeric only and
limited to 'length'.

Suffix parameter will be append to the end of the resulting short string and length reduced accordingly,
so:

	--string="my/branch/name" --length="3" --suffix="!" => "my!"

The 'conditional-match' and 'conditional-value' operate as a pair. If the valud of 'string' is the same
as the 'conditional-match' then all returned strings (safe & full) will become the value of
'conditional-value'. The intention here is to allow swapping of things like branch main being 'main'
so you want the value 'production' returned for use.

Usage:

	safe-string [flags]

The flags are:

	--string			(required)
	--length			(required)
	--suffix
	--conditional-match
	--conditional-value
*/
package safestring

import "flag"

var (
	Name             string = "safe-string"
	FlagSet                 = flag.NewFlagSet(Name, flag.ExitOnError) // Argument group
	original                = FlagSet.String("string", "", "String to make tag safe")
	length                  = FlagSet.Int("length", -1, "Max length of the string")
	suffix                  = FlagSet.String("suffix", "", "Optional suffix to append")
	conditionalMatch        = FlagSet.String("conditional-match", "", "If the original string matches this value, then use the conditional_value directly.")
	conditionalValue        = FlagSet.String("conditional-value", "", "When original matches conditional_match use this value for all other outputs directly.")
)
