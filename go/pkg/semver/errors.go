package semver

var (
	ErrorInvalidSemver           string = "error: [%v] is not a valid semver."
	ErrorInvalidPrerelease       string = "error: [%v] is not a valid semver prerelease segment."
	ErrorConversionAsMapNoMatch  string = "error: failed to convert. String [%s] did not match against semver regex."
	ErrorInvalidIncrement        string = "error: [%v] is not a valid increment value."
	ErrorPrereleaseNoBuildNumber string = "error: item has no buildnumber."
)
