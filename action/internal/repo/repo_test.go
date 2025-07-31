package repo

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func TestRepoClone(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip()
	}
	var (
		dir  = t.TempDir()
		auth = &http.BasicAuth{
			Username: "",
			Password: os.Getenv("GITHUB_TOKEN"),
		}
	)

	_, err := Clone(dir, "https://github.com/ministryofjustice/opg-github-actions.git", auth, nil)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}

}

func TestRepoShallowClone(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip()
	}

	var (
		dir  = t.TempDir()
		auth = &http.BasicAuth{
			Username: "",
			Password: os.Getenv("GITHUB_TOKEN"),
		}
	)

	_, err := ShallowClone(dir, "https://github.com/ministryofjustice/opg-github-actions.git", auth)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	}

}
