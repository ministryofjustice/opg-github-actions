package tags

import (
	"fmt"
	"opg-github-actions/pkg/semver"
	"strings"

	"github.com/dchest/uniuri"
)

func (t *Tags) Unique(tagName string, allTags []string) (newTag string, err error) {
	short := 5
	// check semver
	isSemver := semver.Valid(tagName)
	exists := t.ExistsIn(tagName, allTags)
	newTag = tagName
	// now handle the adjustments
	if isSemver {
		semv, _ := semver.New(tagName)
		preRel := semv.IsPrerelease()

		for exists == true {
			if preRel {
				rand := strings.ToLower(uniuri.NewLen(short))
				semv.SetPrereleasePrefix(rand)
			} else {
				semv.BumpPatch()
			}
			newTag = semv.String()
			exists = t.ExistsIn(newTag, allTags)
		}
	} else {
		for exists == true {
			rand := strings.ToLower(uniuri.NewLen(short))
			newTag = fmt.Sprintf("%s.%s", tagName, rand)
			exists = t.ExistsIn(newTag, allTags)
		}
	}

	return
}
