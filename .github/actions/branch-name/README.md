# Branch Name Composite Action

Provides a consistent `branch_name` for the active repository and versions that can be using for tagging. The *tag safe* versions are alphanumeric only and therefore remove special characters and common seperators like `/` and `-` that would be added by tools like dependabot and renovate.

**Works with `pull_request` and `push` workflows only**.

## Usage

Within you github workflow job you can place a step such as:

```yaml
- name: "Generate safe branch name"
      id: branch_name
      uses: 'ministryofjustice/opg-github-actions/.github/actions/branch-name@v2.1.3'
```

### Using the data

Example of using the data in another step

```yaml
    - name: "Create tag"
      id: create_tag
      uses: './.github/actions/create-tag'
      with:
          commitish: ${{steps.branch_name.outputs.branch_name}}
          tag_name: ${{steps.next_tag.outputs.next_tag}}
```

## Inputs and Outputs

It does not require any inputs and will return the following data:

- `branch_name`
- `full_length`
- `source_commitish`
- `destination_commitish`
- **`safe`**


#### `branch_name`

In a `pull_request` workflow, this is `github.pull_request.head.ref` value, and will be the name of the branch being worked on - eg `my-feature-1`.
For a `push` workflow, this is `github.ref` value, typically the branch where the code has been pushed into - eg `main`.
Both of these then have `refs/heads` removed from their value.

This value will retain any special charaters or seperaters from the branch name, so can also look like `dependabot/package/update`.

#### `full_length`

This is the `branch_name` value as only alphanumeric characters, so `branch_name` of `dependabot/update/12` this would be `dependabotupdate12`.

#### `safe`

This is a truncated form of `full_length`, limited to 12 characters. This is what you would typically use within git tags.

#### `source_commitish` and `destination_commitish`

Useful for semver and similar processes, these variables contain a `commit-ish` reference for comparison between the points in history. On a `pull_request` workflow, they will contain the branch names being merged (so `my-feature-1` and `main`), and for a `push` workflow they contain the before and after commit references.
