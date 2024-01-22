package repo

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v5"
)

// InitRepo will trigger a plain repository from an empty folder path at 'directoty'
func InitRepo(directory string) (r *git.Repository, err error) {
	slog.Debug(fmt.Sprintf("intialising repo at [%s]", directory))
	return git.PlainInit(directory, false)
}

// OpenRepo returns a Repository pointer based on an already existing git repository
// folder structure whose root directory is found at 'directory'
func OpenRepo(directory string) (r *git.Repository, err error) {
	slog.Debug(fmt.Sprintf("opening repo at [%s]", directory))
	return git.PlainOpen(directory)
}
