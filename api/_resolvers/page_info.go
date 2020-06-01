package resolvers

import (
	"fmt"

	"github.com/mjm/speedrungql/speedrun"
)

type PageInfo struct {
	pi *speedrun.PageInfo
}

func (pi *PageInfo) StartCursor() *Cursor {
	return nil
}

func (pi *PageInfo) EndCursor() *Cursor {
	if !pi.HasNextPage() {
		return nil
	}

	nextCursor := pi.pi.Max + pi.pi.Offset
	c := Cursor(fmt.Sprintf("%d", nextCursor))
	return &c
}

func (pi *PageInfo) HasNextPage() bool {
	return pi.pi.Max == pi.pi.Size
}

func (pi *PageInfo) HasPreviousPage() bool {
	return false
}
