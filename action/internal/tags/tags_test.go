package tags

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var gauth = &http.BasicAuth{
	Username: "",
	Password: os.Getenv("GITHUB_TOKEN"),
}

type tTagStrings struct {
	Refs     []*plumbing.Reference
	Error    error
	Expected []string
}

func TestTagsStrings(t *testing.T) {

	var tests = []*tTagStrings{
		{
			Error: nil,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("refs/heads/v4", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("refs/tags/v1.1.0-pre.0", "1ecf0ef2c2dffb79603ae5a02219af86ec6484e4"),
			},
			Expected: []string{
				"refs/heads/v4 6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
				"refs/tags/v1.1.0-pre.0 1ecf0ef2c2dffb79603ae5a02219af86ec6484e4",
			},
		},
	}

	for _, test := range tests {
		actual := Strings(test.Refs, test.Error)

		for _, expected := range test.Expected {
			found := false

			for _, act := range actual {
				if act == expected {
					found = true
				}
			}
			if !found {
				t.Errorf("failed to find match, expected [%s] in set :\n%v\n", expected, actual)
			}
		}

	}

}

func TestTagsRefs(t *testing.T) {

	var tests = []*tTagStrings{
		{
			Error: nil,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("refs/heads/v4", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("refs/tags/v1.1.0-pre.0", "1ecf0ef2c2dffb79603ae5a02219af86ec6484e4"),
			},
			Expected: []string{
				"refs/heads/v4",
				"refs/tags/v1.1.0-pre.0",
			},
		},
	}

	for _, test := range tests {
		actual := Refs(test.Refs, test.Error)

		for _, expected := range test.Expected {
			found := false

			for _, act := range actual {
				if act == expected {
					found = true
				}
			}
			if !found {
				t.Errorf("failed to find match, expected [%s] in set :\n%v\n", expected, actual)
			}
		}

	}

}

func TestTagsShortRefs(t *testing.T) {

	var tests = []*tTagStrings{
		{
			Error: nil,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("refs/heads/v4", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("refs/tags/v1.1.0-pre.0", "1ecf0ef2c2dffb79603ae5a02219af86ec6484e4"),
			},
			Expected: []string{
				"v4",
				"v1.1.0-pre.0",
			},
		},
	}

	for _, test := range tests {
		actual := ShortRefs(test.Refs, test.Error)

		for _, expected := range test.Expected {
			found := false

			for _, act := range actual {
				if act == expected {
					found = true
				}
			}
			if !found {
				t.Errorf("failed to find match, expected [%s] in set :\n%v\n", expected, actual)
			}
		}

	}

}

// func TestTagsAllWithAShallowRepo(t *testing.T) {
// 	if os.Getenv("GITHUB_TOKEN") == "" {
// 		t.Skip()
// 	}

// 	var (
// 		err error
// 		r   *git.Repository
// 		dir = t.TempDir()
// 	)

// 	r, err = repo.ShallowClone(dir, "https://github.com/ministryofjustice/opg-github-actions.git", gauth)
// 	if err != nil {
// 		t.Errorf("err: %s", err.Error())
// 		t.FailNow()
// 	}

// 	_, err = All(r)
// 	if err == nil {
// 		t.Errorf("expected error about being a shallow clone ")
// 	}

// }

// func TestTagsAllWithANormalRepo(t *testing.T) {
// 	if os.Getenv("GITHUB_TOKEN") == "" {
// 		t.Skip()
// 	}

// 	var (
// 		err error
// 		r   *git.Repository
// 		// tags       []*plumbing.Reference
// 		stringTags []string
// 		dir        string = "./repo-test" //t.TempDir()
// 	)
// 	os.RemoveAll(dir)
// 	os.MkdirAll(dir, os.ModePerm)

// 	r, err = repo.Clone(dir, "https://github.com/ministryofjustice/opg-github-actions.git", gauth, nil)
// 	if err != nil {
// 		t.Errorf("err: %s", err.Error())
// 		t.FailNow()
// 	}

// 	stringTags = Strings(All(r))

// 	for _, tg := range stringTags {
// 		fmt.Println(tg)
// 	}

// 	t.FailNow()

// }
