# opg-github-actions

Collection of re-usable, composite actions to share between teams foused on small, common tasks that we do in many places.

### Branch Name

Intended to provide a consistent method for getting branch data out of the active repository. Primarily provides the current branch name and a *tag safe* version of the branch name thats truncated to **12** characters.

**Works with `pull_request` and `push` workflows only**.

[More Details](./.github/actions/branch-name/README.md)


### Create Tag

Provides a method to create a lightweight git tag at the commit-ish value passed. If the tag name that is requested to be created already exists a different tag name is used instead.

[More Details](./.github/actions/create-tag/README.md)
