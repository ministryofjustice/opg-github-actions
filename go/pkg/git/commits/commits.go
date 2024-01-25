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

	return
}

func (c *Commits) StrToReference(str string) (ref *plumbing.Reference, err error) {
	return StrToRef(c.repository, str)
}

func StrToRef(r *git.Repository, str string) (ref *plumbing.Reference, err error) {
	refName := str
	// if the string doesnt start with "refs/", presume its a short form and look for a match
	// comparing the last segment of the full reference
	if !strings.Contains(refName, "refs/") {
		refs, _ := r.References()
		slog.Debug("commits: StrToReference looking for matching ref..")
		refs.ForEach(func(ref *plumbing.Reference) error {
			name := ref.Name().String()
			end := strings.HasSuffix(name, "/"+str)
			suf := "❌"
			if end {
				refName = name
				suf = "✅"
			}
			slog.Debug(fmt.Sprintf("[%s] commits: StrToReference - str [%s] match ref [%s] ? [%t]", suf, str, name, end))
			return nil
		})
	}
	rev := plumbing.Revision(refName)
	hash, err := r.ResolveRevision(rev)
	slog.Info(fmt.Sprintf("commits: StrToReference [%s:%s] => revision [%s] => hash [%s]", str, refName, rev, hash))
	if err != nil {
		return
	}
	ref = plumbing.NewReferenceFromStrings(str, hash.String())
	return
}
