package testlib

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v58/github"
)

func TestEventPullRequest(base string, head string, title string, body string) []byte {
	var (
		action   string = "open"
		id       int64  = 100
		addtions int    = 4
		number   int    = 10
	)

	pr := &github.PullRequestEvent{
		Action: &action,
		PullRequest: &github.PullRequest{
			ID:        &id,
			Number:    &number,
			Additions: &addtions,
			Body:      &body,
			Title:     &title,
			Base: &github.PullRequestBranch{
				Ref: &base,
			},
			Head: &github.PullRequestBranch{
				Ref: &head,
			},
		},
	}

	b, _ := json.Marshal(pr)

	return b
}

func TestEventPush(head string, before string, after string) []byte {
	var (
		action string = "open"
		id     int64  = 100
		ref    string = fmt.Sprintf("refs/heads/%s", head)
	)

	p := &github.PushEvent{
		Action: &action,
		PushID: &id,
		Ref:    &ref,
		Before: &before,
		After:  &after,
	}

	b, _ := json.Marshal(p)

	return b
}
