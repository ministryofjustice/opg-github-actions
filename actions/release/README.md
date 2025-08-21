# Release Composite Action

Will create a new release using `gh` cli tool based on the tag value passed along.

Limited to one release per tag.


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


## Inputs and Outputs

Common inputs:
- `tag` (required)
- `prerelease` (default: "true")

Rarely used inputs:
- `github_token`
- `latest`
- `release_artifact`
- `release_notes_flag` (default: "--notes-from-tag")


### Inputs

#### `tag` (required)
Git tag to associated with this release. Normally provided by the `semver` action.

#### `prerelease` (default: "true")
Flag to say if this is a prerelease or not, defaults to true. There is no logic to change this, so make sure it is accurate.

#### `github_token`
By default, the action uses the `github.token` value to push to the repository, but if you need a different scope of auth, then pass along your own token in this variable.

#### `latest`
By default, the `--latest` flag is set when `prerelease` is `false` (as in a production release), but you can use this field to force a value.

#### `release_artifact`
Pattern or file path for artifacts you want to attach to this release, such as built binaries. Runs from the `github.workspace` directory.

#### `release_notes_flag` (default: "--notes-from-tag")
When creating a release with the `gh` cli tool there are two two methods for generating notes, this lets you swap between them.

### Outputs

#### `url`
Link to the url of the release for the tag passed.
