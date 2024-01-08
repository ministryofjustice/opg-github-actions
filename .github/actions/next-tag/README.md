# Next Tag Composite Action

Used to generate the next semver-ish tag suitable for the repository based on the last release, the branch name and the commit history between the points specified.

Looks for `#major|minor|patch` within the subject, body and notes of each commit message and will increment the `latest_tag` suitably.

For prereleases, the prerelease segment is updated and if a major version bump is changed, that will be updated as well. This means your feature branch might start as `1.1.0-myfeature.2`, but if you then do a commit with `#major` in the next tag will be `2.0.0-myfeature.0`.

**Squashed commits** - With a squashed commit you may get an inconsistent version number compared to the pull_request version that is generated. As squashed commits flattern the commit history into a singlular commit the new commit message may not contain the same `#major|minor|patch` triggers as the original history.


## Usage

Within you github workflow job you can place a step such as this for finding prerelease versions:

```yaml
    - name: "Find next tag"
      id: next_tag
      uses: 'ministryofjustice/opg-github-actions/.github/actions/next-tag@v2.1.3'
      with:
          prerelease: "true"
          prerelease_suffix: "myfeature"
          latest_tag: "1.1.0-myfeature.0"
          last_release: "1.0.1"
```
or for release versions:

```yaml
    - name: "Find next tag"
      id: next_tag
      uses: 'ministryofjustice/opg-github-actions/.github/actions/next-tag@v2.1.3'
      with:
          prerelease_suffix: "myfeature"
          last_release: "1.0.1"
          commitish_a: "my-feature"
          commitish_b: "main"

```

## Inputs and Outputs

Inputs:
- `prerelease`
- `prerelease_suffix`
- `latest_tag`
- `last_release`
- `commitish_a`
- `commitish_b`
- `default_bump` (default; "patch")
- `with_v` (default: ""|False)


Outputs:
- `prerelease`
- `prerelease_suffix`
- `latest_tag`
- `last_release`
- `commitish_a`
- `commitish_b`
- `default_bump` (default; "patch")
- `with_v` (default: ""|False)
- `majors`
- `minors`
- `patches`
- **`next_tag`**

### Inputs

#### `prerelease`
A boolean-ish value, that when true ("true", "True", true etc) will look for existing tags that use the `prerelease_suffix`.

#### `prerelease_suffix`
A tag safe version of a branch name. This is used to find existing tags for this branch by looking for following pattern against prerelease tags: `${prerelease_suffix}.[0-9]+$"`.

#### `latest_tag`
The latest tag created with the prerelease suffix. This can be found by using [`latest-tag` action](../latest-tag/README.md)

#### `last_release`
The semver-ish tag of the last release version in the repository.

#### `commitish_a` and `commitish_b`
The two points in git commit history to use as comparisions and look for the #major | #minor | #patch string which will then determine any version increments.

This can be found by using [`branch-name` action](../branch-name/README.md)

#### `default_bump` (default: "patch")
If there are no version bump triggers found within the commits between `commitish_a` and `commitish_b` then this value will be used as the default increment for a version number.
In the case of prereleases, the prerelease counter is increased instead (1.1.0-myfeature.2 => 1.1.0-myfeature.3)

#### `with_v` (default: "" | False)
If enabled, the next_tag generated will start with a `v` prefix, such as `v1.1.0-myfeature.3`


### Outputs

#### `prerelease`
Mirror of the inputted value.

#### `prerelease_suffix`
Mirror of the inputted value.

#### `latest_tag`
Mirror of the inputted value.

#### `last_release`
Mirror of the inputted value.

#### `commitish_a` and `commitish_b`
Mirror of the inputted value.

#### `default_bump` (default: "patch")
Mirror of the inputted value.

#### `with_v` (default: "" | False)
Mirror of the inputted value.

#### `majors`, `minors` and `patches`
These are counters showing how many of each trigger was found within the commits found between `commitish_a` and `commitish_b`.

#### `next_tag`
The next_tag that should be used based on the commits and config passed in. This will be semver-ish and may contain a `v` prefix if that has been enabled.
