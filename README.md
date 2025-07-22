# opg-github-actions



# Composite actions information

Usefuly information found from testing the setup of actions

When using an action either the code is checked out into a known location (`github.action_path`). When you are pulling the action remotely (via `uses`) additional information is also set - the repository name and reference point - within the github context as `github.action_repository` and `github.action_ref` respectively.

When checking out the action into the `action_path` only the source code is included, no binaries or other attachments on the release - this makes sense when the `sha` being used might not actually be a release, it is purely related to the git tag / commit.

When pinning to a `sha` there is no direct way to find the release for that, you have to inspect git history.

If we commit the built binaries is it then visible?


## Paths

Actions being within `./github/actions` is now only an requirement for local actions, calling from `uses` can be any folder path.
