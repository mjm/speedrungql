package resolvers

import (
	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

type Category struct {
	speedrungql.Category
}

func (c *Category) ID() graphql.ID {
	return relay.MarshalID("category", c.Category.ID)
}

func (c *Category) RawID() string {
	return c.Category.ID
}

func (c *Category) Type() string {
	switch c.Category.Type {
	case "per-game":
		return "PER_GAME"
	case "per-level":
		return "PER_LEVEL"
	default:
		return ""
	}
}

func (c *Category) Players() *CategoryPlayers {
	return &CategoryPlayers{c.Category.Players}
}

type CategoryPlayers struct {
	speedrungql.CategoryPlayers
}

func (c *CategoryPlayers) Type() string {
	switch c.CategoryPlayers.Type {
	case "exactly":
		return "EXACTLY"
	case "up-to":
		return "UP_TO"
	default:
		return ""
	}
}

func (c *CategoryPlayers) Value() int32 {
	return int32(c.CategoryPlayers.Value)
}
