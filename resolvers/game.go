package resolvers

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
		Direction *speedrungql.OrderDirection
	}
	First *int32
	After *Cursor
}) (*GameConnection, error) {
	var opts []speedrungql.FetchOption
	if args.Order != nil {
		opts = append(opts, speedrungql.WithOrder(args.Order.Field, args.Order.Direction))
	}
	if args.Filter != nil {
		if args.Filter.Name != nil {
			opts = append(opts, speedrungql.WithFilter("name", *args.Filter.Name))
		}
	}
	if args.First != nil {
		opts = append(opts, speedrungql.WithLimit(int(*args.First)))
	}
	if args.After != nil {
		offset, err := args.After.GetOffset()
		if err != nil {
			return nil, err
		}
		opts = append(opts, speedrungql.WithOffset(offset))
	}

	games, pageInfo, err := v.client.ListGames(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &GameConnection{
		client:   v.client,
		games:    games,
		pageInfo: pageInfo,
	}, nil
}

type GameConnection struct {
	client   *speedrungql.Client
	games    []*speedrungql.Game
	pageInfo *speedrungql.PageInfo
}

func (gc *GameConnection) Edges() []*GameEdge {
	var edges []*GameEdge
	for _, g := range gc.games {
		edges = append(edges, &GameEdge{
			Node: &Game{*g, gc.client},
		})
	}
	return edges
}

func (gc *GameConnection) Nodes() []*Game {
	var nodes []*Game
	for _, g := range gc.games {
		nodes = append(nodes, &Game{*g, gc.client})
	}
	return nodes
}

func (gc *GameConnection) PageInfo() *PageInfo {
	return &PageInfo{gc.pageInfo}
}

type GameEdge struct {
	Node *Game
}

func (*GameEdge) Cursor() *Cursor {
	return nil
}

type Game struct {
	speedrungql.Game
	client *speedrungql.Client
}

func (g *Game) ID() graphql.ID {
	return relay.MarshalID("game", g.Game.ID)
}

func (g *Game) RawID() string {
	return g.Game.ID
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

func (g *Game) Ruleset() *GameRuleset {
	return &GameRuleset{g.Game.Ruleset}
}

func (g *Game) Platforms(ctx context.Context) ([]*Platform, error) {
	plats, err := g.client.GetPlatforms(ctx, g.Game.Platforms)
	if err != nil {
		return nil, err
	}

	var res []*Platform
	for _, plat := range plats {
		res = append(res, &Platform{*plat})
	}

	return res, nil
}

func (g *Game) Categories(ctx context.Context) ([]*Category, error) {
	cats, err := g.client.ListGameCategories(ctx, g.Game.ID)
	if err != nil {
		return nil, err
	}

	var res []*Category
	for _, cat := range cats {
		res = append(res, &Category{*cat})
	}
	return res, nil
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

type GameRuleset struct {
	speedrungql.GameRuleset
}

func (gr *GameRuleset) RunTimes() []GameRunTime {
	var rts []GameRunTime
	for _, rt := range gr.GameRuleset.RunTimes {
		rts = append(rts, GameRunTime(rt))
	}
	return rts
}

func (gr *GameRuleset) DefaultRunTime() GameRunTime {
	return GameRunTime(gr.GameRuleset.DefaultRunTime)
}

type GameRunTime speedrungql.GameRunTime

func (GameRunTime) ImplementsGraphQLType(name string) bool {
	return name == "GameRunTime"
}

func (v GameRunTime) String() string {
	return strings.ToUpper(string(v))
}

func (v *GameRunTime) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("GameRunTime value was not a string")
	}

	switch s {
	case "REALTIME":
		*v = GameRunTime(speedrungql.RealTime)
	case "REALTIME_NOLOADS":
		*v = GameRunTime(speedrungql.RealTimeNoLoads)
	case "INGAME":
		*v = GameRunTime(speedrungql.InGame)
	default:
		return fmt.Errorf("unknown GameRunTime value %q", s)
	}

	return nil
}
