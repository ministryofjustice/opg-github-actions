# Create Tag Composite Action

Provides a method to create a lightweight git tag a set location specified via a git commit-ish reference.

## Usage

Within you github workflow job you can place a step such as:

```yaml
    - name: "Create tag"
      id: create_tag
      uses: 'ministryofjustice/opg-github-actions/.github/actions/create-tag@v2.3.1'
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
- **`created`**
- `requested`
- `success`
- `regenerated`

### Inputs

#### `commitish`
A git reference to where the `tag_name` should be created. This can be a commit hash, branch name or object name.

#### `tag_name`
The desired tag you want to create. If the tag already exists, then an alternative is generated.

If you are using semver notation - a prerelease tag with have its prerelease segment adjusted, otherwise a major version bump will be triggered.

#### `test`
When true ("true" or true), the tag will be created, but is **NOT** pushed to the remote and therefore will not persist.

### Outputs

#### `requested`
The tag name originally requested to be created. Will always match `tag_name`.

#### `created`
The value of the tag that was created. It will normally match the `tag_name`, but in cases where that tag already existed a new tag will be generated using its information and the generated version.

#### `success`
Boolean value to show if the creation was successful. Determined by being able to find the request tag in the set of *local* tags after the create.

#### `regenerated`
Boolean flag to state if the `tag_name` had to be changed due to a clash.
