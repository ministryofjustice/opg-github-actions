# Create Tag Composite Action

Provides a method to create a lightweight git tag a set location specified via a git commit-ish reference.

## Usage

Within you github workflow job you can place a step such as:

```yaml
    - name: "Create tag"
      id: create_tag
      uses: 'ministryofjustice/opg-github-actions/.github/actions/create-tag@v2.1.3'
      with:
          commitish: "${git_reference}"
          tag_name: "1.0.0-myfeature.2"
```

## Inputs and Outputs

Inputs:
- `commitish`
- `tag_name`
- `test`

Outputs:
- `commitish`
- `original_tag_name`
- `requested`
- `test`
- `latest`
- `created`

### Inputs

#### `commitish`
A git reference to where the `tag_name` should be created. This can be a commit hash, branch name or object name.

#### `tag_name`
The desired tag you want to create. If the tag already exists, then an alternative is generated.

If you are using semver notation - a prerelease tag with have its prerelease segment adjusted, otherwise a major version bump will be triggered.

#### `test`
When true ("True", "true" or true), the tag will be created, but not pushed to the remote and therefore will not persist.

### Outputs

#### `commitish`
The inputted `commitish` value.

#### `original_tag_name` and `requested`
Contains the `tag_name` value.

#### `test`
Contains the inputted `test` value.

#### `latest`
Latest contains the last tag that was created at the `commitish` value passed in. As the sorting of the value is based on string sort, this will be largest, but not necessarily the requested tag.

For example, if at commit `ae034f` there are 3 tags already (`1.0.0`, `10.0.1` and `9.0.5`) and this creates `2.0.1` the latest value would be `10.0.1`

#### `created`
The value of the tag that was created. It will normally match the `original_tag_name`, but in cases where that tag already existed a new tag will be generated and this will contain that value.
