# Detect Changes Composite Action

Detects whether files in specified paths have changed between commits. Requires `actions/checkout` with `fetch-depth: 0` in the calling job.

## Usage

```yaml
- name: Detect changes
  id: detect-changes
  uses: 'ministryofjustice/opg-github-actions/actions/detect-changes@<SHA> # <version>'
  with:
    paths: |
      service-admin
      terraform
      **/*.md
```

## Inputs

#### `paths` (required)

List of paths to check for changes.

## Outputs

#### `has_changes`
`true` if any files in the specified paths changed.

#### `only_changed`
`true` if the specified paths are the only things that changed

#### `files`
List of all changed files in the specified paths.

## Event support

| Event | Comparison |
| --- | --- |
| `pull_request` | Merge base of branch against base ref (`origin/<base>...HEAD`) |
| `push` | Previous commit against pushed commit |
| `merge_group` | Merge group base SHA against HEAD |
| `workflow_dispatch` | `origin/main` against HEAD |

Any other event type will cause the action to exit with an error.