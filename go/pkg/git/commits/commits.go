package commits

import (
	"errors"
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/git/repo"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Commits struct {
	Directory  string
	repository *git.Repository
}

func New(directory string) (t *Commits, err error) {
	if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("Repostitory directory not found [%s]", directory)
	}

	r, err := repo.OpenRepo(directory)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	w, err := r.Worktree()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defaultBranch, err := r.Head()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Error("branch:" + defaultBranch.Name().String())
	err = w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  true,
		Branch: defaultBranch.Name(),
	})

	t = &Commits{
		Directory:  directory,
		repository: r,
	}
	return
}

func (c *Commits) StrToReference(str string) (ref *plumbing.Reference, err error) {
	rev := plumbing.Revision(str)
	hash, err := c.repository.ResolveRevision(rev)
	if err != nil {
		// _, e := c.repository.Worktree()
		// slog.Error("worktree:" + e.Error())
		// _, e = c.repository.Head()
		// slog.Error("head:" + e.Error())
		// slog.Error("str: " + str)
		slog.Error("rev: " + rev.String())
		return
	}
	ref = plumbing.NewReferenceFromStrings(str, hash.String())
	return
}
