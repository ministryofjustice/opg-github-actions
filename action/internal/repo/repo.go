package repo

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const ErrDirectoryNotFound string = "Directory not found: [%s]"

// directoryExists checks if the path exists and is a directory
func directoryExists(path string) (exists bool) {
	exists = false
	info, err := os.Stat(path)

	if err == nil && info.IsDir() {
		exists = true
	}
	return
}

// FromDir creates a git.Repository from an existing checked out git repo
// at `localDirectory` and returns
//
// If the directory does not exist or the `PlainOpen` command fails then
// an error is returned
func FromDir(localDirectory string) (repo *git.Repository, err error) {

	if !directoryExists(localDirectory) {
		err = fmt.Errorf(ErrDirectoryNotFound, localDirectory)
		return
	}

	repo, err = git.PlainOpen(localDirectory)
	return
}

// Clone will clone the git repository at `remoteUrl` into the `localDirecory` path
// and return a git.Repository for it
//
// the `auth` param presumes a basic auth like:
//
//	&http.BasicAuth{
//		Username: "username",
//		Password: os.Getenv("GITHUB_TOKEN"),
//	}
func Clone(localDirectory string, remoteUrl string, auth *http.BasicAuth, opts *git.CloneOptions) (r *git.Repository, err error) {

	if opts == nil {
		opts = &git.CloneOptions{URL: remoteUrl}
	}

	if auth != nil {
		opts.Auth = auth
	}

	r, err = git.PlainClone(localDirectory, false, opts)

	return
}

// ShallowClone will checkout a repo but without branches or tags
//
// Mimics the actions/checkout behaviour when `depth:` is left as default (1)
func ShallowClone(localDirectory string, remoteUrl string, auth *http.BasicAuth) (r *git.Repository, err error) {
	var opts = &git.CloneOptions{
		URL:               remoteUrl,
		ShallowSubmodules: true,
		Depth:             1,
		Tags:              git.NoTags,
	}

	r, err = Clone(localDirectory, remoteUrl, auth, opts)

	return
}

// Init create a new git repository at the localDirectory
func Init(localDirectory string) (r *git.Repository, err error) {
	r, err = git.PlainInit(localDirectory, false)
	if err != nil {
		return
	}
	return
}

// fetch might not be needed - added for when
// repo is shallow and doesnt have all the refs
// when then causes a failure on branch look up
func Fetch(lg *slog.Logger, r *git.Repository, auth *http.BasicAuth) (err error) {
	lg = lg.With("operation", "Fetch")

	lg.Info("fetching updates from remotes ...")

	remotes, err := r.Remotes()
	specs := []config.RefSpec{
		"refs/*:refs/*",
		"HEAD:refs/heads/HEAD",
		"+refs/tags/*:refs/tags/*",
		"+refs/heads/*:refs/remotes/origin/*",
	}
	for _, remote := range remotes {
		// fetch branches and tags for this remote
		name := remote.Config().Name
		lg.Debug("fetching data for remote ", "remote", name)

		err = r.Fetch(&git.FetchOptions{
			RemoteName: name,
			RefSpecs:   specs,
			Auth:       auth,
		})
		// this isnt an error, so handle and ignore it
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			lg.Warn("repository up to date", "warning", err.Error())
			err = nil
		}

	}
	return
}
