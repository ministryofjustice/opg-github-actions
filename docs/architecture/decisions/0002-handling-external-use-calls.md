# 2. Handling External Use Calls

Date: 2025-07-22

## Status

Accepted

## Context

This repository provides github actions used in many other repositories via github actions `uses` keyword. We follow recommended approach of using pinned versions in our github workflows to improve supply chain security, so we generally specify the exact git commit, rather than a movable tag.

Therefore, when you call an action with the sha you do not have an immediate way to find built binaries as a release is based on names & git tags only.

There are roughly five approches for this:

1. __Use bash only__. This would be fastest, but maintainability would suffer as would testing capabilities.
2. __Use docker container action__. Based on what github docs say, this would be slower and we'd have to maintain images and use a public registry
3. __Use Node / Python__. We have limited Node / Javascript knowledge in the team and Python requires a lot of setup and install - making it slower.
4. __Commit binaries__. Manually commiting we'd risk forgetting or commiting bad / out-of-sync versions; Automated we would have to resolve pipeline issues like enforce signed commits and recursive triggers.
5. __Always build__. This would mean every call has ~1 min build time.


## Decision

We've decided to build on every request. This is mostly what happens now so limited change and we can focus on smaller / faster build processes to reduce that time for everyone.

Using GitHubs docker registery (GitHub Packages) or Docker Hub might be a future route, but it needs more investigation. There may be rate limiting and other issues from using a public registry.


## Consequences

Every call to `opg-github-action` action will trigger a go build.
