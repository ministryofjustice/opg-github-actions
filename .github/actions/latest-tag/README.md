# Latest Tag Composite Action

Tries to find the last semver styled release tag and the latest tag thats relevant to the `prerelease_suffix` passed.

For releases, the code uses `git tags --list`, checks each as being a valid semver and then filters to those that are not a prerelease.

For prerelease, `git tags --list` data is then filtered for prereleases that look like `1.1.1-${prerelease_suffix}.${counter}`.


## Usage

Within you github workflow job you can place a step such as:

```yaml
    - name: "Find latest tag"
      id: latest_tag
      uses: 'ministryofjustice/opg-github-actions/.github/actions/latest-tag@v2.3.1'
      with:
          branch_name: "my-feature"
          prerelease: "true"
          prerelease_suffix: "myfeature"
```

## Inputs and Outputs

Inputs:
- `prerelease`
- `prerelease_suffix`
- `branch_name` (default: "beta")
- `release_branch` (default: "main")

Outputs:
- `prerelease`
- `prerelease_suffix`
- **`last_release`**
- **`last_prerelease`**

### Inputs

#### `prerelease`
A boolean-ish value, that when true ("true", "True", true etc) will look for existing tags that use the `prerelease_suffix`. This value can be overridden in the code if the `branch_name` passed is within the `release_branch` - typically if its flagged as being a prerelease, but its on `main` branch, it should be considered a release.

#### `prerelease_suffix`
A tag safe version of `branch_name`. This is used to find existing tags for this branch by looking for following pattern against prerelease tags: `${prerelease_suffix}.[0-9]+$"`.

#### `branch_name` (default: "beta")
The branch being used that you want to find the latest tag against.

#### `release_branches` (default: "main,master")
These are used as a configurable item incase your release is to something like "production" and "main" is actually a working / dev tree.


### Outputs

#### `prerelease`
A boolean value; this is the calculated version of the inputted `prerelease`, after `release_branch` has been compared to the `branch_name` and shows how the code thinks this tag should be classified.

#### `prerelease_suffix`
The inputted suffix value.

#### `last_release`
The semver-ish tag of the last release tag in the repository.

#### `last_prerelease`
The last prerelease tag created with the prerelease suffix passed in.
