package commits

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// logHashes returns the hashes starting from the commit passed
func logHashes(repository *git.Repository, commit *object.Commit) (log map[string]bool, err error) {
	log = map[string]bool{}
	iter, err := repository.Log(&git.LogOptions{From: commit.Hash})
	err = iter.ForEach(func(c *object.Commit) error {
		hash := c.Hash.String()
		log[hash] = true
		return nil
	})

	return
}

// DiffBetween gets the commits on base and head and then finds those that are present in head,
// but not within bases history, then returns the commit objects for those
//
// Intention is act in similar fashion to `git log main..my-branch` to return new commits
// which are then used to look for trigger strings for semver
//
//   - base => main
//   - head => feature-branch
func DiffBetween(repository *git.Repository, base plumbing.Hash, head plumbing.Hash) (commits []*object.Commit, err error) {
	var (
		baseCommit *object.Commit
		baseLog    map[string]bool
		headCommit *object.Commit
		headLog    map[string]bool
	)
	commits = []*object.Commit{}
	baseCommit, err = repository.CommitObject(base)
	if err != nil {
		return
	}
	headCommit, err = repository.CommitObject(head)
	if err != nil {
		return
	}

	baseLog, err = logHashes(repository, baseCommit)
	if err != nil {
		return
	}
	headLog, err = logHashes(repository, headCommit)
	if err != nil {
		return
	}

	for hash, _ := range headLog {
		// if not found, then add to set of missing commits
		if _, ok := baseLog[hash]; !ok {
			commit, _ := repository.CommitObject(plumbing.NewHash(hash))
			commits = append(commits, commit)
		}
	}

	return
}
