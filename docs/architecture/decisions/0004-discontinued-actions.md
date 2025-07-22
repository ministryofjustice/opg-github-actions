# 4. Discontinued Actions

Date: 2025-07-22

## Status

Accepted

## Context

In previous versions we maintained two similar, but slightly different actions (`safe-strings` and `branch-name`). There is very little difference between the them, with safe-string providing an override mechanism that is no longer required.

The only usage of this action is within `opg-github-workflows` which is mostly defunct as well.

## Decision

We will only support `branch-name` moving forward.

## Consequences

Any location requiring `safe-strings` can either adjust or stay on current versions.
