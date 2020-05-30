package resolvers

import (
	"github.com/mjm/speedrungql"
)

type Link struct {
	speedrungql.Link
}

func (l *Link) Rel() *string {
	if l.Link.Rel == "" {
		return nil
	}
	return &l.Link.Rel
}
