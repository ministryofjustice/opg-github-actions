package repo

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// InitRepo will trigger a plain repository from an empty folder path at 'directoty'
func InitRepo(directory string) (r *git.Repository, err error) {
	slog.Debug(fmt.Sprintf("intialising repo at [%s]", directory))
	r, err = git.PlainInit(directory, false)
	if err != nil {
		return
	}
	return
}

// OpenRepo returns a Repository pointer based on an already existing git repository
// folder structure whose root directory is found at 'directory'
func OpenRepo(directory string) (r *git.Repository, err error) {
	slog.Debug(fmt.Sprintf("opening repo at [%s]", directory))

	r, err = git.PlainOpen(directory)
	if err != nil {
		return
	}
	// we need to fetch all branch info from the remotes
	fetch(r)
	return
}

func CloneRepo(directory string, url string) (r *git.Repository, err error) {
	slog.Debug(fmt.Sprintf("opening repo at [%s]", directory))
	r, err = git.PlainClone(directory, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: "opg-github-actions",
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})

	if err != nil {
		return
	}
	fetch(r)
	return
}

func fetch(r *git.Repository) (err error) {
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
