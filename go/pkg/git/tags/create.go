package tags

import (
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

// CreateAt will make a tag on this repository at the commit hash passed
func (t *Tags) CreateAt(tagName string, ref *plumbing.Hash) (*plumbing.Reference, error) {
	return t.repository.CreateTag(tagName, *ref, nil)
}

// Push send a push of tags from local to remote.
// Uses the environment variable 'GITHUB_TOKEN' for authentication
func (t *Tags) Push() (err error) {
	err = t.repository.Push(
		&git.PushOptions{
			RemoteName: "origin",
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
			Auth:       auth,
		},
	)
	if err != nil {
		slog.Error("error with push:")
		slog.Error(err.Error())
	}
	return
}
