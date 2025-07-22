# Using actions in a mono-repo structure

When uses a github action there are some useful context variables present to determine how this action is being used - they are as below:

```
github.action_path
github.action_repository
github.action_ref
```

Both `github.action_repository` and `github.action_ref` have to be referenced within `env` to be readable

You can read more about them in the [github reference doc](https://docs.github.com/en/actions/reference/contexts-reference#github-context).

## Used within its own repository

When being used within its own repostiory via a relative path include like this:

```
    - id: "local_action"
      name: "Local"
      uses: ./.github/actions/test-action
```

## Used externally by via version tag

When being used within from another repository like this:

```
    - id: "remote_action"
      name: "Remote action"
      uses: ministryofjustice/opg-github-actions/.github/actions/test-action@v-restructured-test
```


## Used externally by via pinned sha

When being used within from another repository like this:

```
    - id: "remote_action"
      name: "Remote action"
      uses: ministryofjustice/opg-github-actions/.github/actions/test-action@56dc197f27b49f50f5807e97d6b79e9c39536dab #v-restructured-test
```
