# Semver Composite Action

Will generate a new (or return latest) git tag formatted as a semver, typically used for releases, docker image tags and so on.

Commit messages can contain `#major, #minor, #patch` to trigger semver increment increase.

You can toggle the use of a `v` prefix on or off depending on your needs.

A set of collated information is sent to `${GITHUB_STEP_SUMMARY}` as a markdown table at the end of the run.

## Usage

This action is often used in conjuction with the release action to allow tag creation to be split from the release like this:

```yaml
    - name: "Create Semver tag"
      id: semver
      uses: 'ministryofjustice/opg-github-actions/actions/semver@v4.2.0'
      with:
        prerelease: ${{ github.ref != 'refs/heads/main' }}
        create_release: false
        github_token: ${{ github.token }}
    ... other steps ...
    - name: "Create release"
      id: release
      uses: 'ministryofjustice/opg-github-actions/actions/release@v4.2.0'
      with:
        tag: ${{ steps.semver.outputs.tag }}
        prerelease: ${{ github.ref != 'refs/heads/main' }}
        github_token: ${{ github.token }}
```


To generate a prerelease (`v1.2.0-mybranch.1`) style tag within your workflow which also creates a release:

```yaml
    - name: "Create Semver tag"
      id: semver
      uses: 'ministryofjustice/opg-github-actions/actions/semver@v4.2.0'
```

To generate a release tag for use in a workflow that determines if its a prerelease or not dynamically:

```yaml
    - name: "Create Semver tag"
      id: semver
      uses: 'ministryofjustice/opg-github-actions/actions/semver@v4.2.0'
      with:
        prerelease: ${{ github.ref != 'refs/heads/main' }}
```

To return a the latest release tag that can be re-used:

```yaml
    - name: "Create Semver tag"
      id: semver
      uses: 'ministryofjustice/opg-github-actions/actions/semver@v4.2.0'
      with:
        prerelease: false
        bump: "none"
```


## Inputs and Outputs

Common inputs:
- `prerelease` (default: "true")
- `create_release` (default: "true")
- `default_bump` (default: "patch")
- `release_artifact` (default: "")

Rarely used inputs:
- `prelease_suffix_length` (default: "14")
- `branch_name`
- `without_prefix` (default: "false")
- `github_token`
- `release_notes_flag` (default: "--notes-from-tag")
- `test` (default: "false")

Outputs:
- **`tag`**
- `hash`
- `created`
- `bump`
- `test`

### Inputs

#### `prerelease` (default: "true")
Flag to say if this is a prerelease or not, defaults to true. There is no logic to change this, so make sure it is accurate.

#### `create_release` (default: "true")
Boolean to determine if a github release is also created for this tag. This makes use of the `gh` tool to do this and you will need suitable permissions on the workflow.

#### `default_bump` (default: "patch")
If there are no commits found, or no commit messages containing `#major, #minor, #patch` then then created semver will be based on this increment.

#### `release_artifact` (default: "")
Pattern or file path for artifacts you want to attach to this release, such as built binaries. Runs from the `github.workspace` directory.

#### `prelease_suffix_length` (default: "14")
Length of the suffix to use in creating a prerelease tag

#### `branch_name`
The branch name is used for generating the prerelease suffix (`v1.2.3-$suffix.1`) and is generally determined from the github context values (`github.head_ref` & `github.ref_name`), but you can pass in a value here to overwrite that.

#### `without_prefix` (default: "false")
By default, the semver tag is created with a `v` prefix at the start - if this is set to true, then it will be removed.

#### `github_token`
By default, the action uses the `github.token` value to push to the repository, but if you need a different scope of auth, then pass along your own token in this variable.

#### `release_notes_flag` (default: "--notes-from-tag")
When creating a release with the `gh` cli tool there are two two methods for generating notes, this lets you swap between them.

#### `test` (default: "false")
When set to `true`, the semver tag is no actually created, allows you try the workflow with another tool.

### Outputs

#### `tag`
Contains the tag that has been created.

#### `hash`
The git hash / sha that the tag was created at.

#### `created`
A boolean shoing if the tag was actually created - helpful for `test` and `none` increment usages.

#### `bump`
What bump was used to create

#### `test`
Boolean mirroring the `test` input.
