package resolvers

import (
	"github.com/mjm/speedrungql/speedrun"
)

type Link struct {
	speedrun.Link
}

func (l *Link) Rel() *string {
	if l.Link.Rel == "" {
		return nil
	}
	return &l.Link.Rel
}

func wrapLink(link *speedrun.Link) *Link {
	if link == nil {
		return nil
	}

	return &Link{*link}
}
