package tags

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"opg-github-actions/action/internal/logger"
	"strconv"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type tTagSorted struct {
	Refs     []*plumbing.Reference
	Expected []*plumbing.Reference
}

func TestTagsSort(t *testing.T) {
	var lg = logger.New("error", "text")
	var tests = []*tTagSorted{
		{
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("9.5.0", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("tag-test-01", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("10.1.0", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("tag-test-02", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("1.0.0-beta.1+bA1", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("1.0.0-beta.0+bA1", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("1.0.0-beta.0+bA2", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
			},
			Expected: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("1.0.0-beta.0+bA1", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("1.0.0-beta.0+bA2", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("1.0.0-beta.1+bA1", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("9.5.0", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("10.1.0", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("tag-test-01", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("tag-test-02", "6ecf0ef2c2dffaa96033e5a02219af86ec6584e5"),
			},
		},
	}

	for i, test := range tests {
		sorted := Sort(lg, test.Refs, SORT_ASC)

		for idx, actual := range sorted {
			if actual.Name().Short() != test.Expected[idx].Name().Short() {
				t.Errorf("order not as expected in set [%d:%d], expected [%v] actual [%v]", i, idx, test.Expected[idx], actual)
			}
		}

	}
}

type tTagStrings struct {
	Refs     []*plumbing.Reference
	Error    error
	Expected []string
}

func TestTagsStrings(t *testing.T) {
	var lg = logger.New("error", "text")
	var tests = []*tTagStrings{
		{
			Error: nil,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("refs/heads/v4", "6ecf0ef2c2000000000000000000000000000000"),
				plumbing.NewReferenceFromStrings("refs/tags/v1.1.0-pre.0", "6ecf0ef2c2000000000000000000000000000000"),
			},
			Expected: []string{
				"refs/heads/v4 6ecf0ef2c2000000000000000000000000000000",
				"refs/tags/v1.1.0-pre.0 6ecf0ef2c2000000000000000000000000000000",
			},
		},
	}

	for _, test := range tests {
		actual := Strings(lg, test.Refs, test.Error)

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
	var lg = logger.New("error", "text")
	var tests = []*tTagStrings{
		{
			Error: nil,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("refs/heads/v4", "6ecf0ef2c2daaaa96033a5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("refs/tags/v1.1.0-pre.0", "1ecf0ef2c2dffb79603ae5a02219af86ec6484e4"),
			},
			Expected: []string{
				"refs/heads/v4",
				"refs/tags/v1.1.0-pre.0",
			},
		},
	}

	for _, test := range tests {
		actual := Refs(lg, test.Refs, test.Error)

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
	var lg = logger.New("error", "text")
	var tests = []*tTagStrings{
		{
			Error: nil,
			Refs: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("refs/heads/v4", "6ecf0ef2c2daaaa96033a5a02219af86ec6584e5"),
				plumbing.NewReferenceFromStrings("refs/tags/v1.1.0-pre.0", "1ecf0ef2c2dffb79603ae5a02219af86ec6484e4"),
			},
			Expected: []string{
				"v4",
				"v1.1.0-pre.0",
			},
		},
	}

	for _, test := range tests {
		actual := ShortRefs(lg, test.Refs, test.Error)

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

func TestTagCreationSimple(t *testing.T) {
	var (
		lg      = logger.New("error", "text")
		dir     = t.TempDir()
		tagName = "test-tag"
	)

	repo, head := randomRepository(dir)
	tags1, _ := All(lg, repo)

	Create(lg, repo, tagName, head.Hash())

	tags2, _ := All(lg, repo)

	if len(tags1)+1 != len(tags2) {
		t.Errorf("tag count not as expected")
	}

	found := false
	for _, tg := range tags2 {
		var nm = tg.Name().Short()
		if nm == tagName {
			found = true
		}
	}
	if !found {
		t.Errorf("did not find newly created tag")
	}

}

// randomRepository make a repo with a mix of commits and various tags
func randomRepository(dir string) (r *git.Repository, defaultBranch *plumbing.Reference) {
	var (
		commitsN = rand.Intn(50) + 10 // somewhere between 10-50 commits
		preTags  = rand.Intn(15) + 5  // 5-15 tags
		relTags  = rand.Intn(10) + 2  // 2-10 tags
		hashes   = []plumbing.Hash{}
	)
	// create the repository locally
	r, _ = git.PlainInit(dir, false)
	w, _ := r.Worktree()

	// create some commits on the base
	for i := 0; i < commitsN; i++ {
		msg := fmt.Sprintf("commit %d", i)
		hash, e := w.Commit(msg, &git.CommitOptions{
			AllowEmptyCommits: true,
			Author:            &object.Signature{Name: "go test", Email: "test@example.com"},
		})
		if e == nil {
			hashes = append(hashes, hash)
		}
	}
	// make some random tags with v1.0.x-thing format
	for i := 0; i < preTags; i++ {
		var randI = rand.Intn(len(hashes))
		var hash = hashes[randI]
		var tag = "v1.0." + strconv.Itoa(i+1) + "-thing"
		var rev = plumbing.Revision(hash.String())

		revHash, _ := r.ResolveRevision(rev)
		r.CreateTag(tag, *revHash, nil)
	}

	// make some random tags with v1.0.x-thing format
	for i := 0; i < relTags; i++ {
		var randI = rand.Intn(len(hashes))
		var hash = hashes[randI]
		var tag = "v1.0." + strconv.Itoa(i+1)
		var rev = plumbing.Revision(hash.String())

		revHash, _ := r.ResolveRevision(rev)
		r.CreateTag(tag, *revHash, nil)
	}

	defaultBranch, _ = r.Head()
	// checkout to default branch
	w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  true,
		Branch: defaultBranch.Name(),
	})

	return
}

// Debug is a helper function that runs printf against a json
// string version of the item passed.
// Used for testing only.
func debug[T any](item T) {
	bytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Printf("%+v\n", string(bytes))
}
