package commits

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// logHashes
func (c *Commits) logHashes(commit *object.Commit) (log map[string]bool, err error) {
	log = map[string]bool{}
	iter, err := c.repository.Log(&git.LogOptions{From: commit.Hash})
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
func (c *Commits) DiffBetween(base plumbing.Hash, head plumbing.Hash) (commits []*object.Commit, err error) {
	var (
		baseCommit *object.Commit
		baseLog    map[string]bool
		headCommit *object.Commit
		headLog    map[string]bool
	)
	commits = []*object.Commit{}
	slog.Info(fmt.Sprintf("commits DiffBetween [%s] [%s]", base.String(), head.String()))
	baseCommit, err = c.repository.CommitObject(base)
	if err != nil {
		return
	}
	headCommit, err = c.repository.CommitObject(head)
	if err != nil {
		return
	}

	baseLog, err = c.logHashes(baseCommit)
	if err != nil {
		return
	}
	headLog, err = c.logHashes(headCommit)
	if err != nil {
		return
	}

	for hash, _ := range headLog {
		// if not found, then add to set of missing commits
		if _, ok := baseLog[hash]; !ok {
			commit, _ := c.repository.CommitObject(plumbing.NewHash(hash))
			commits = append(commits, commit)
		}
	}

	return
}
