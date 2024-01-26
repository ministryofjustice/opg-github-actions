# Semver Tag Composite Action

Uses multiple steps to mimic the sequence of `branch-name`, `latest-tag`, `next-tag` & `create-tag` and calls the same function code underneath to create a semver tag.

With `test` enabled, the created tag will not be pushed to the remote and only kept locally.

If the `release_branch` value matches the active branch a release is triggered.

A set of collated information is sent to `${GITHUB_STEP_SUMMARY}` as a markdown table at the end of the run.

## Usage

Within you github workflow job you can place a step such as:

```yaml
    - name: "Semver tag"
      id: semver_tag
      uses: 'ministryofjustice/opg-github-actions/.github/actions/semver-tag@v2.3.1'
      with:
          prerelease: "true"
          with_v: ""
          github_token: ${{ secrets.token }}
```

## Inputs and Outputs

Inputs:
- `prerelease`
- `release_branch` (default: "main")
- `with_v` (default: "true")
- `show_verbose_summary` (default: "")
- `test` (default: "")
- `releases_enabled` (default: 'true')
- `draft_release` (default: 'false')

Outputs:
- **`created_tag`**
- `release_id`
- `release_url`
- Outputs from sub-actions, please see their respective README's
  - [`branch_original`](../branch-name/README.md)
  - [`branch_full_length`](../branch-name/README.md)
  - [`branch_safe`](../branch-name/README.md)
  - [`latest_tag_last_prerelease`](../latest-tag/README.md)
  - [`latest_tag_last_release`](../latest-tag/README.md)
  - [`next_tag`](../next-tag/README.md)
  - [`next_tag_commitish_a`](../latest-tag/README.md)
  - [`next_tag_commitish_b`](../latest-tag/README.md)
  - [`create_tag_latest`](../create-tag/README.md)
  - [`create_tag_success`](../create-tag/README.md)


### Inputs

#### `prerelease`
Flag to say if this is a prerelease or not. Can be overridden by logic within the code if the active branch matches a release_branch.

#### `release_branch` (default: "main")
The branch that should be considered a release when being pushed to.  If the active branch matches this value then a release is triggered.

#### `with_v` (default: "true")
Determines if the semver tags are created with a `v` prefix.

#### `show_verbose_summary` (default: "")
If this is "true" then the larger collated information will be outputed to the `${GITHUB_STEP_SUMMARY}`

#### `test` (default: "")
When true ("True", "true" or true), the tag will be created, but not pushed to the remote and therefore will not persist.

#### `releases_enabled` (default: "true")
When this is true and this is on a releaase branch and not a test, then a release will be created.

#### `draft_release` (default: "false")
If true, will mark and release created as being in draft.


### Outputs

#### `created_tag`
Contains the tag that has been created.

#### `release_id` and `release_url`
If a release waas triggered (by the active branch matching release_branch) then these will contain release id and url that can be outputed elsewhere.
