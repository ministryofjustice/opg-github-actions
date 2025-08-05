# Branch Name Composite Action

Provides a consistent `branch_name` for the active repository and versions that can be using for tagging. The outputs of this action are alphanumeric only and therefore remove special characters and common seperators like `/` and `-` that would be added by tools like dependabot and renovate.


## Usage

Within you github workflow job you can place a step such as, which will try to determine branch name from github context variables:

```yaml
- name: "Generate safe branch name"
  id: branch_name
  uses: 'ministryofjustice/opg-github-actions/actions/branch-name@v1.2.3'
```

However, if you are using this action outside of `pull_request` and `push` events, you can directly specify the original value to work from like:

```yaml
- name: "Generate safe branch name"
  id: branch_name
  uses: 'ministryofjustice/opg-github-actions/actions/branch-name@v1.2.3'
  with:
    name: feature/my-branch-1
```



## Inputs and Outputs

Inputs:
- `name`
- `length` (default: 14)

Outputs:
- `branch_name`
- `full_length`
- **`safe`**

### Inputs

#### `name`

If you are using this action outside of a pull_request or push event, you can pass a value her to use as the base for sanitisation.

#### `length`

The maximum length the `safe` value should be - defaults to 12. This is the length check done after being converted to a alphanumeric string only.

### Outputs

#### `branch_name`

In a `pull_request` workflow, this is `github.head_ref` value, and will be the name of the branch being worked on - eg `my-feature-1`.

For a `push` workflow, this is `github.ref_name` value, typically the branch where the code has been pushed into - eg `main`.

This value will retain any special charaters or seperaters from the branch name, so can also look like `dependabot/package/update`.

#### `full_length`

This is the `branch_name` value as only alphanumeric characters, so `branch_name` of `dependabot/update/12` this would be `dependabotupdate12`.

#### `safe`

This is a truncated form of `full_length`, limited to `length` (12 by default) characters. This is what you would typically use within git tags.
