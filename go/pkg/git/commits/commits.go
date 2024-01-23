package commits

import (
	"errors"
	"fmt"
	"log/slog"
	"opg-github-actions/pkg/git/repo"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/k0kubun/pp"
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
	// fetch branches
	_, err = r.Branches()
	return
}

func (c *Commits) StrToReference(str string) (ref *plumbing.Reference, err error) {
	rev := plumbing.Revision(str)
	hash, err := c.repository.ResolveRevision(rev)
	if err != nil {
		branchIter, e := c.repository.Branches()
		pp.Println(branchIter)
		branches := []string{}
		branchIter.ForEach(func(ref *plumbing.Reference) error {
			branches = append(branches, ref.Name().Short())
			return nil
		})
		slog.Error(fmt.Sprintf("branches: [%s]", strings.Join(branches, " ")))
		if e != nil {
			slog.Error(e.Error())
		}
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
