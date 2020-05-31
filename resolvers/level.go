package resolvers

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

type Level struct {
	speedrungql.Level
	client *speedrungql.Client
}

func (l *Level) ID() graphql.ID {
	return relay.MarshalID("level", l.Level.ID)
}

func (l *Level) RawID() string {
	return l.Level.ID
}

func (l *Level) Game(ctx context.Context) (*Game, error) {
	gameURI := speedrungql.FindLink(l.Links, "game")
	if gameURI == "" {
		return nil, nil
	}

	game, err := l.client.GetGame(ctx, gameURI)
	if err != nil {
		return nil, err
	}

	if game == nil {
		return nil, nil
	}
	return &Game{*game, l.client}, nil
}
