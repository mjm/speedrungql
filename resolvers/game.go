package resolvers

import (
	"context"
	"encoding/json"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

func (v *Viewer) Games(ctx context.Context, args struct {
	Filter *struct {
		Name *string
	}
	Order *struct {
		Field     *string
		Direction *string
	}
	First *int32
	After *Cursor
}) (*GameConnection, error) {
	u := v.r.baseURL + "/games"
	res, err := v.r.httpClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp speedrungql.GamesResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &GameConnection{res: &resp}, nil
}

type GameConnection struct {
	res *speedrungql.GamesResponse
}

func (gc *GameConnection) Edges() []*GameEdge {
	var edges []*GameEdge
	for _, g := range gc.res.Data {
		edges = append(edges, &GameEdge{
			Node: &Game{g},
		})
	}
	return edges
}

func (gc *GameConnection) Nodes() []*Game {
	var nodes []*Game
	for _, g := range gc.res.Data {
		nodes = append(nodes, &Game{g})
	}
	return nodes
}

func (gc *GameConnection) PageInfo() PageInfo {
	return PageInfo{}
}

type GameEdge struct {
	Node *Game
}

func (*GameEdge) Cursor() *Cursor {
	return nil
}

type Game struct {
	speedrungql.Game
}

func (g *Game) ID() graphql.ID {
	return relay.MarshalID("game", g.Game.ID)
}

func (g *Game) Names() *GameNames {
	return &GameNames{g.Game.Names}
}

func (g *Game) Abbreviation() *string {
	if g.Game.Abbreviation == "" {
		return nil
	}
	return &g.Game.Abbreviation
}

func (g *Game) Platforms(ctx context.Context) []*Platform {
	return nil
}

type GameNames struct {
	speedrungql.GameNames
}

func (gn *GameNames) Japanese() *string {
	if gn.GameNames.Japanese == "" {
		return nil
	}
	return &gn.GameNames.Japanese
}

func (gn *GameNames) Twitch() *string {
	if gn.GameNames.Twitch == "" {
		return nil
	}
	return &gn.GameNames.Twitch
}
