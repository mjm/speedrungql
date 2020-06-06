package resolvers

import (
	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

type Engine struct {
	speedrun.Engine
	client *speedrun.Client
}

func (e *Engine) ID() graphql.ID {
	return relay.MarshalID("engine", e.Engine.ID)
}

func (e *Engine) RawID() string {
	return e.Engine.ID
}
