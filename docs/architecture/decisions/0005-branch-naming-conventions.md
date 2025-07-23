# 5. Branch Naming Conventions

Date: 2025-07-22

## Status

Accepted

## Context

When running a pull request or similar event within GitHub and other CI tooling we need to make use of a sensible identifier to align with non-production version of docker image / infrastructure. The action in this repository (`branch-name`) is generally used to provide that value.

There are certain restrictions on this value due to where it is used. For example, using a forward slash in some contexts, like a git tag, would be valid, it would not be accepted as a docker tag on ECR.

## Decision

The action will generate an alphanumeric only string that is reasonably short to stay within tagging / naming character limits while also being relatively unique.


## Consequences
