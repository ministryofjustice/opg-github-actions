package repo

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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

func CloneRepo(directory string, url string) (r *git.Repository, err error) {
	slog.Debug(fmt.Sprintf("opening repo at [%s]", directory))
	return git.PlainClone(directory, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: "opg-github-actions",
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})
}
