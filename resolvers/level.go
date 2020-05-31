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

func (l *Level) Categories(ctx context.Context) ([]*Category, error) {
	cats, err := l.client.ListLevelCategories(ctx, l.Level.ID)
	if err != nil {
		return nil, err
	}

	var res []*Category
	for _, cat := range cats {
		res = append(res, &Category{*cat, l.client})
	}
	return res, nil
}

func (l *Level) Variables(ctx context.Context) ([]*Variable, error) {
	vars, err := l.client.ListLevelVariables(ctx, l.Level.ID)
	if err != nil {
		return nil, err
	}

	var res []*Variable
	for _, v := range vars {
		res = append(res, &Variable{*v, l.client})
	}
	return res, nil
}
