package commonstrings

var (
	ErrorArgumentMissing       string = "Argument [--%s] is required, but missing"
	ErrorArgumentInvalidChoice string = "Argument [--%s] is invalid, value should be one of [%s], actual [%s]"
	ErrorArgumentFileNotExist  string = "Argument [--%s] is invalid, value refers to a file that does not exist: [%s]"
	ErrorArgumentDirNotExist   string = "Argument [--%s] is invalid, value refers to a directory that does not exist: [%s]"
	ErrorArumentNotBoolean     string = "Argument [--%s] is not a valid boolean [%v]"
)
