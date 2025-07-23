# 6. Branch Name Action Changes

Date: 2025-07-22

## Status

Accepted

## Context

The previous version of `branch-name` only works on `pull_request` and `push` GitHub event types reads directly from the provided json file. This has meant the action is tightly coupled to the `github` context object, so it cannot run at other times.

## Decision

Update the action to take an option input of branch_name, then when left empty is worked out from the same github context values, but from within the action.yml and not in the go code.

- If not empty, use the input variable, otherwise..
- for `pull_requests` check `github.head_ref` property
- for `push`, check `github.ref_name` property
- if still empty, throw an error

See event docs on the [github context](https://docs.github.com/en/actions/reference/contexts-reference#github-context) and example [push event contents](https://docs.github.com/en/actions/reference/contexts-reference#example-contents-of-the-github-context).

The go code should take two simple parameters - starting branch name and maximum length. Fully formed branch names should be reduced (so `refs/heads/` removed from the start of the string)

If should return the following:

- The original branch name as `branch_name`
- The full length but tag safe value as `full_length`
- The tag safe and shortened value as `safe`

## Consequences

Some minor compatabilty changes and internal re-working should mean the majority of uses are not effected
