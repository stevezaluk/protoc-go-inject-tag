package inject

import (
	"fmt"
	"strings"
)

type TagItem struct {
	Key   string
	Value string
}

func TagFromComment(comment string) (tag string) {
	match := GetRegex(CommentRegex).FindStringSubmatch(comment)
	if len(match) == 2 {
		tag = match[1]
	}
	return
}

type tagItems []TagItem

func (ti tagItems) format() string {
	tags := []string{}
	for _, item := range ti {
		tags = append(tags, fmt.Sprintf(`%s:%s`, item.Key, item.Value))
	}
	return strings.Join(tags, " ")
}

func (ti tagItems) override(nti tagItems) tagItems {
	overrided := []TagItem{}
	for i := range ti {
		dup := -1
		for j := range nti {
			if ti[i].Key == nti[j].Key {
				dup = j
				break
			}
		}
		if dup == -1 {
			overrided = append(overrided, ti[i])
		} else {
			overrided = append(overrided, nti[dup])
			nti = append(nti[:dup], nti[dup+1:]...)
		}
	}
	return append(overrided, nti...)
}

func newTagItems(tag string) tagItems {
	items := []TagItem{}
	splitted := GetRegex(TagsRegex).FindAllString(tag, -1)

	for _, t := range splitted {
		sepPos := strings.Index(t, ":")
		items = append(items, TagItem{
			Key:   t[:sepPos],
			Value: t[sepPos+1:],
		})
	}
	return items
}
