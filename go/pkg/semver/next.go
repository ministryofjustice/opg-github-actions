package semver

func Next(
	lastPrerelease *Semver, lastRelease *Semver,
	prerelease bool, prereleaseSuffix string,
	counters IncrementCounters,
) (next *Semver, err error) {

	var tag *Semver
	lastPrereleaseValid := Valid(lastPrerelease.String())

	if prerelease && lastPrereleaseValid {
		tag, _ = New(lastPrerelease.String())
		// tag = *lastPrerelease
	} else {
		tag, _ = New(lastRelease.String())
		// tag = *lastRelease
	}

	if counters.Major > 0 {
		if prerelease && *tag.Major() <= *lastRelease.Major() {
			tag.BumpMajor()
			tag.SetMinor(0)
			tag.SetPatch(0)
			if !tag.IsPrerelease() {
				tag.MustBumpPrerelease(prereleaseSuffix)
			}
		} else if prerelease {
			tag.MustBumpPrerelease(prereleaseSuffix)
		} else {
			tag.BumpMajor()
			tag.SetPrerelease("")
		}
	} else if counters.Minor > 0 {
		if prerelease && lastPrereleaseValid {
			// Last release of v1.4.0 + is a prerelease + has a minor flag + latest_tag of v1.5.0-beta.0
			// => v1.5.0-beta.1
			tag.MustBumpPrerelease(prereleaseSuffix)
		} else if prerelease {
			// Last release of v1.4.0 + is a prerelease + has a minor flag
			// based on the release tag
			// => v1.5.0-beta.0
			tag.BumpMinor()
			tag.SetPatch(0)
			tag.MustBumpPrerelease(prereleaseSuffix)
		} else {
			// Last release of v1.4.0 + has a minor flag
			// based on the release tag
			// => v1.5.0
			tag.BumpMinor()
			tag.SetPrerelease("")
		}
	} else if counters.Patch > 0 {
		if prerelease && lastPrereleaseValid {
			// Last release of v1.4.0 + is a prerelease + has a minor flag + latest_tag of v1.4.1-beta.0
			// => v1.4.1-beta.1
			tag.MustBumpPrerelease(prereleaseSuffix)
		} else if prerelease {
			// Last release of v1.4.0 + is a prerelease + has a minor flag
			// based on the release tag
			// => v1.4.1-beta.0
			tag.BumpPatch()
			tag.MustBumpPrerelease(prereleaseSuffix)
		} else {
			// Last release of v1.4.0 + has a minor flag
			// based on the release tag
			// => v1.4.1
			tag.BumpPatch()
			tag.SetPrerelease("")
		}
		if *tag.Patch() == 0 {
			tag.BumpPatch()
		}
	}
	next = tag

	return
}
