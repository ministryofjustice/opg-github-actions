package repo

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var auth = &http.BasicAuth{
	Username: "opg-github-actions",
	Password: os.Getenv("GITHUB_TOKEN"),
}

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
		URL:  url,
		Auth: auth,
	})

	if err != nil {
		return
	}
	fetch(r)
	return
}

func fetch(r *git.Repository) (err error) {
	slog.Info("fetching remotes ...")
	var refs []*plumbing.Reference
	w, _ := r.Worktree()
	remotes, err := r.Remotes()
	for _, remote := range remotes {
		// fetch branches and tags for this remote
		name := remote.Config().Name
		slog.Info("fetching remote data for: " + name)

		specs := []config.RefSpec{
			// "refs/*:refs/*",
			"+refs/tags/*:refs/tags/*",
			"+refs/heads/*:refs/remotes/origin/*",
			// "HEAD:refs/heads/HEAD",
		}
		for _, spec := range specs {
			slog.Info("ref spec: " + spec.String())
			err = r.Fetch(&git.FetchOptions{
				RemoteName: name,
				RefSpecs:   []config.RefSpec{spec},
			})
			if err != nil {
				slog.Error("Error fetching from origin: ")
				slog.Error(err.Error())
				// return
			}
		}

		if err != nil {
			return
		}

		refs, err = remote.List(&git.ListOptions{Auth: auth})
		if err != nil {
			slog.Error("Error fetching remote list")
			slog.Error(err.Error())
			return
		}
		for _, rf := range refs {
			rfName := rf.Name()
			slog.Debug(remote.Config().Name + " -> " + rfName.Short())
			if rf.Name().IsBranch() {
				slog.Info(remote.Config().Name + " found branch: " + rfName.Short())

				w.Checkout(&git.CheckoutOptions{Create: false, Force: true, Branch: rfName})
			}
		}

	}
	return
}
