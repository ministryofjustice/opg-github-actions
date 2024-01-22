package tags

import (
	"errors"
	"fmt"
	"opg-github-actions/pkg/git/repo"
	"os"

	"github.com/go-git/go-git/v5"
)

type Tags struct {
	Directory  string
	repository *git.Repository
}

// New opens an existing repo from directory and returns the Tags type
func New(directory string) (t *Tags, err error) {
	if _, err := os.Stat(directory); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("Repostitory directory not found [%s]", directory)
	}

	r, err := repo.OpenRepo(directory)
	t = &Tags{
		Directory:  directory,
		repository: r,
	}
	return
}
