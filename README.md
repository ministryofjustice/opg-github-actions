# opg-github-actions

Collection of re-usable, composite actions to share between teams foused on small, common tasks that we do in many places.

## Details

All actions have a github workflow that provides basic tests for their functionality, these are then run as part of the pr & release pipeline of this repository ensuring they pass.

Most actions are now `go` based within a shared [shared library](./go/) with all related tests. There are end-to-end semver tests which mimic how `action.yml` run the commands within [`main_test.go`](./go/main_test.go) and chain them together as the semver action would.

There is a [Makefile](./Makefile) to help build and create release versions of the go code and it also contains target to call tun tests.


## Available Actions

### Branch Name

Intended to provide a consistent method for getting branch data out of the active repository. Primarily provides the current branch name and a *tag safe* version of the branch name thats truncated to **12** characters.

**Works with `pull_request` and `push` workflows only**.

[More Details](./.github/actions/branch-name/README.md)


### Create Tag

Provides a method to create a lightweight git tag at the commit-ish value passed. If the tag name that is requested to be created already exists a different tag name is used instead.

[More Details](./.github/actions/create-tag/README.md)


### Latest Tag

Tries to find the last semver styled release tag and the latest tag thats relevant to the `prerelease_suffix` passed.

[More Details](./.github/actions/latest-tag/README.md)


### Next Tag

Used to generate the next semver-ish tag suitable for the repository based on the last release, the branch name and the commits in the history between the points specified.

Looks for `#major|minor|patch` within the subject, body and notes of each commit message and will increment the `latest_tag` suitably.

[More Details](./.github/actions/next-tag/README.md)



### Safe Strings

A helper action to convert a string to only contain alphanumeric values and the option to limit its overall length. Intented to be used for tagging and similar activities where other characters (like `/` and `|`) cause errors.

String is converted to lowercase version.


[More Details](./.github/actions/safe-strings/README.md)


### Semver Tag

With `test` enabled, the created tag will not be pushed to the remote and only kept locally.

If the `release_branch` value matches the active branch a release is triggered.

A set of collated information is sent to `${GITHUB_STEP_SUMMARY}` as a markdown table at the end of the run.

[More Details](./.github/actions/semver-tag/README.md)


### Terraform Version

Used to find and then switch terraform versions in projects where multiple versions maybe in use (monorepos with many state files etc).

Parse the terraform versions file from the directory passed and return the required terraform version string.

[More Details](./.github/actions/terraform-version/README.md)

### Terraform Workspace Manager

A shared terraform workspace tool to track and list workspaces protected from deletion and a time that should be protected for.

[More Details](./.github/actions/terraform-workspace-manager/README.md)

## Tests

Tests can be run using commands with the `Makefile` from the root of this directory.

### Using the `Makefile`

Run all tests:

```bash
make tests
```

Run a specific test:

```bash
make test test-name="TestRemoteRepo"
```

Run a series of tests based on a patter:

```bash
make test test-name="TestSemver*"
```

### Running tests without make

All `go` code can be tested using :

```bash
cd ./go
go test -v ./...
```

By default, the tests disable and minimise the log data, but this can be overridden with environment variables:

```bash
cd ./go
env LOG_LEVEL="info" go test -v ./...
```
