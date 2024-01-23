package commits

import (
	"errors"
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/git/repo"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
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
	t = &Commits{
		Directory:  directory,
		repository: r,
	}

	// we need to fetch all branch info from the remotes
	slog.Debug("fetching remotes ...")
	remotes, err := r.Remotes()

	for _, remote := range remotes {
		// fetch branches and tags for this remote
		slog.Debug("fetching remote data for :" + remote.Config().Name)
		r.Fetch(&git.FetchOptions{
			RemoteName: remote.Config().Name,
			RefSpecs:   []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD", "+refs/tags/*:refs/tags/*"},
		})

	}

	return
}

func (c *Commits) StrToReference(str string) (ref *plumbing.Reference, err error) {
	rev := plumbing.Revision(str)
	hash, err := c.repository.ResolveRevision(rev)
	if err != nil {
		return
	}
	ref = plumbing.NewReferenceFromStrings(str, hash.String())
	return
}
