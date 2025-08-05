# 8. Semver Action changes

Date: 2025-08-05

## Status

Accepted

## Context

The previous version of semver action was actually an amalgamation of multiple commands (`branch-name`, `latest-tag`, `next-tag` & `create-tag`) and was built as a single command. This does add have a small overhead on build time (as every action built all go code) and it meant the action.yml was large and complex with a lot of extra inputs and outputs that were never used most of the time.

The previous action did not support attaching files to releases and relied on a 3rd party action to create the release, which is less than ideal.

It did not support a way to return an existing semver when no changes have been made. Typically used when a workflow is re-run against main.


## Decision

Replace the existing `semver-tag` with new action called `semver` that removes reliance on third party code for releases, expands release capability to allow for attachments and provides a way to return the same semver tag for a release.


## Consequences

Location, input and output values of the action have been changed, so upgrades will require some adjustments to use the new version that renovate / dependabot won't pick up on. See the actions README for details of new values.
