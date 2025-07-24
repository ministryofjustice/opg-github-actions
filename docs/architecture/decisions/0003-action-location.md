# 3. Action Source Locations

Date: 2025-07-22

## Status

Accepted

## Context

Originally, there was a limitation on where the source for a github action could be located for GitHub to use it (under `./github/actions/`). This is no longer the case and you can use any folder path within the repository both locally and via `uses` as long as there is an `action.yml` file in the path.

Currently, the folder structure puts all availably actions underneath the `.github` path, which makes it hard to tell what should be used publically and what is used to reduce / reuse code for just this repository.


## Decision

All external actions will be defined underneath a top level `./actions` folder for clarity on internal / external usage.

There will be two top level folders with similar names, `actions` for github and `action` which will be the root of the go code base. This will need to be clearly explained in the project README.


## Consequences

When updating to the new versions of our github actions, the path will need to change as well as the version.
