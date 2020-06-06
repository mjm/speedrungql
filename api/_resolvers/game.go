package resolvers

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

func (v *Viewer) Games(ctx context.Context, args FetchGamesArgs) (*GameConnection, error) {
	return fetchGameConnection(ctx, v.client, args)
}

type FetchGamesArgs struct {
	Filter *struct {
		Name         *string     `filter:"name"`
		Abbreviation *string     `filter:"abbreviation"`
		Released     *int32      `filter:"released"`
		GameType     *graphql.ID `filter:"gametype"`
		Platform     *graphql.ID `filter:"platform"`
		Region       *graphql.ID `filter:"region"`
		Genre        *graphql.ID `filter:"genre"`
		Engine       *graphql.ID `filter:"engine"`
		Developer    *graphql.ID `filter:"developer"`
		Publisher    *graphql.ID `filter:"publisher"`
		Moderator    *graphql.ID `filter:"moderator"`
	}
	Order *struct {
		Field     *GameOrderField
		Direction *speedrun.OrderDirection
	}
	First *int32
	After *Cursor
}

func fetchGameConnection(ctx context.Context, c *speedrun.Client, args FetchGamesArgs, extraOpts ...speedrun.FetchOption) (*GameConnection, error) {
	var opts []speedrun.FetchOption
	opts = append(opts, extraOpts...)
	if args.Order != nil {
		opts = append(opts, speedrun.WithOrder((*string)(args.Order.Field), args.Order.Direction))
	}
	if args.Filter != nil {
		opts = append(opts, speedrun.WithFilters(args.Filter))
	}
	if args.First != nil {
		opts = append(opts, speedrun.WithLimit(int(*args.First)))
	}
	if args.After != nil {
		offset, err := args.After.GetOffset()
		if err != nil {
			return nil, err
		}
		opts = append(opts, speedrun.WithOffset(offset))
	}

	games, pageInfo, err := c.ListGames(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &GameConnection{
		client:   c,
		games:    games,
		pageInfo: pageInfo,
	}, nil
}

type GameOrderField string

func (GameOrderField) ImplementsGraphQLType(name string) bool {
	return name == "GameOrderField"
}

func (v *GameOrderField) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("GameOrderField value was not a string")
	}

	switch s {
	case "NAME_INT":
		*v = "name.int"
	case "NAME_JAP":
		*v = "name.jap"
	default:
		*v = GameOrderField(strings.ToLower(s))
	}
	return nil
}

type GameConnection struct {
	client   *speedrun.Client
	games    []*speedrun.Game
	pageInfo *speedrun.PageInfo
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
	speedrun.Game
	client *speedrun.Client
}

func (g *Game) ID() graphql.ID {
	return relay.MarshalID("game", g.Game.ID)
}

func (g *Game) RawID() string {
	return g.Game.ID
}

func (g *Game) Name(args struct {
	Variant string
}) *string {
	var s string

	switch args.Variant {
	case "INTERNATIONAL":
		s = g.Names.International
	case "JAPANESE":
		s = g.Names.Japanese
	case "TWITCH":
		s = g.Names.Twitch
	}

	if s == "" {
		return nil
	}
	return &s
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
		res = append(res, &Platform{*plat, g.client})
	}

	return res, nil
}

func (g *Game) Regions(ctx context.Context) ([]*Region, error) {
	regs, err := g.client.GetRegions(ctx, g.Game.Regions)
	if err != nil {
		return nil, err
	}

	var res []*Region
	for _, reg := range regs {
		res = append(res, &Region{*reg, g.client})
	}

	return res, nil
}

func (g *Game) Genres(ctx context.Context) ([]*Genre, error) {
	gens, err := g.client.GetGenres(ctx, g.Game.Genres)
	if err != nil {
		return nil, err
	}

	var res []*Genre
	for _, gen := range gens {
		res = append(res, &Genre{*gen, g.client})
	}
	return res, nil
}

func (g *Game) Engines(ctx context.Context) ([]*Engine, error) {
	engs, err := g.client.GetEngines(ctx, g.Game.Engines)
	if err != nil {
		return nil, err
	}

	var res []*Engine
	for _, eng := range engs {
		res = append(res, &Engine{*eng, g.client})
	}
	return res, nil
}

func (g *Game) Moderators() []*GameModerator {
	var gms []*GameModerator
	for userID, role := range g.Game.Moderators {
		gms = append(gms, &GameModerator{
			userID: userID,
			role:   role,
			client: g.client,
		})
	}

	sort.Slice(gms, func(i, j int) bool {
		return gms[i].userID < gms[j].userID
	})
	return gms
}

func (g *Game) Assets() []*GameAsset {
	var assets []*GameAsset
	for kind, asset := range g.Game.Assets {
		if asset == nil {
			continue
		}

		assets = append(assets, &GameAsset{*asset, GameAssetKind(kind)})
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Kind < assets[j].Kind
	})
	return assets
}

func (g *Game) Asset(args struct {
	Kind GameAssetKind
}) *GameAsset {
	asset := g.Game.Assets[speedrun.GameAssetKind(args.Kind)]
	if asset == nil {
		return nil
	}

	return &GameAsset{*asset, args.Kind}
}

func (g *Game) Categories(ctx context.Context) ([]*Category, error) {
	cats, err := g.client.ListGameCategories(ctx, g.Game.ID)
	if err != nil {
		return nil, err
	}

	var res []*Category
	for _, cat := range cats {
		res = append(res, &Category{*cat, g.client})
	}
	return res, nil
}

func (g *Game) Levels(ctx context.Context) ([]*Level, error) {
	levs, err := g.client.ListGameLevels(ctx, g.Game.ID)
	if err != nil {
		return nil, err
	}

	var res []*Level
	for _, lev := range levs {
		res = append(res, &Level{*lev, g.client})
	}
	return res, nil
}

func (g *Game) Variables(ctx context.Context) ([]*Variable, error) {
	vs, err := g.client.ListGameVariables(ctx, g.Game.ID)
	if err != nil {
		return nil, err
	}

	var res []*Variable
	for _, v := range vs {
		res = append(res, &Variable{*v, g.client})
	}
	return res, nil
}

func (g *Game) Runs(ctx context.Context, args FetchRunsArgs) (*RunConnection, error) {
	if args.Filter != nil && args.Filter.Game != nil {
		return nil, errors.New("cannot filter runs by game when reading from a specific game")
	}

	return fetchRunConnection(ctx, g.client, args, speedrun.WithFilter("game", g.Game.ID))
}

type GameRuleset struct {
	speedrun.GameRuleset
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

type GameRunTime speedrun.GameRunTime

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
		*v = GameRunTime(speedrun.RealTime)
	case "REALTIME_NOLOADS":
		*v = GameRunTime(speedrun.RealTimeNoLoads)
	case "INGAME":
		*v = GameRunTime(speedrun.InGame)
	default:
		return fmt.Errorf("unknown GameRunTime value %q", s)
	}

	return nil
}

type GameModerator struct {
	userID string
	role   speedrun.GameModeratorRole
	client *speedrun.Client
}

func (gm *GameModerator) User(ctx context.Context) (*User, error) {
	user, err := gm.client.GetUser(ctx, gm.userID)
	if err != nil {
		return nil, err
	}

	return &User{*user, gm.client}, nil
}

func (gm *GameModerator) Role() GameModeratorRole {
	return GameModeratorRole(gm.role)
}

type GameModeratorRole speedrun.GameModeratorRole

func (GameModeratorRole) ImplementsGraphQLType(name string) bool {
	return name == "GameModeratorRole"
}

func (v GameModeratorRole) String() string {
	switch speedrun.GameModeratorRole(v) {
	case speedrun.Moderator:
		return "MODERATOR"
	case speedrun.SuperModerator:
		return "SUPER_MODERATOR"
	default:
		return ""
	}
}

func (v *GameModeratorRole) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("GameModeratorRole value was not a string")
	}

	switch s {
	case "MODERATOR":
		*v = GameModeratorRole(speedrun.Moderator)
	case "SUPER_MODERATOR":
		*v = GameModeratorRole(speedrun.SuperModerator)
	default:
		return fmt.Errorf("unknown GameModeratorRole value %q", s)
	}

	return nil
}

type GameAssetKind speedrun.GameAssetKind

func (GameAssetKind) ImplementsGraphQLType(name string) bool {
	return name == "GameAssetKind"
}

func (v GameAssetKind) String() string {
	switch speedrun.GameAssetKind(v) {
	case speedrun.AssetLogo:
		return "LOGO"
	case speedrun.AssetCoverTiny:
		return "COVER_TINY"
	case speedrun.AssetCoverSmall:
		return "COVER_SMALL"
	case speedrun.AssetCoverMedium:
		return "COVER_MEDIUM"
	case speedrun.AssetCoverLarge:
		return "COVER_LARGE"
	case speedrun.AssetIcon:
		return "ICON"
	case speedrun.AssetTrophyFirst:
		return "TROPHY_1ST"
	case speedrun.AssetTrophySecond:
		return "TROPHY_2ND"
	case speedrun.AssetTrophyThird:
		return "TROPHY_3RD"
	case speedrun.AssetTrophyFourth:
		return "TROPHY_4TH"
	case speedrun.AssetBackground:
		return "BACKGROUND"
	case speedrun.AssetForeground:
		return "FOREGROUND"
	default:
		return ""
	}
}

func (v *GameAssetKind) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("GameAssetKind value was not a string")
	}

	switch s {
	case "LOGO":
		*v = GameAssetKind(speedrun.AssetLogo)
	case "COVER_TINY":
		*v = GameAssetKind(speedrun.AssetCoverTiny)
	case "COVER_SMALL":
		*v = GameAssetKind(speedrun.AssetCoverSmall)
	case "COVER_MEDIUM":
		*v = GameAssetKind(speedrun.AssetCoverMedium)
	case "COVER_LARGE":
		*v = GameAssetKind(speedrun.AssetCoverLarge)
	case "ICON":
		*v = GameAssetKind(speedrun.AssetIcon)
	case "TROPHY_1ST":
		*v = GameAssetKind(speedrun.AssetTrophyFirst)
	case "TROPHY_2ND":
		*v = GameAssetKind(speedrun.AssetTrophySecond)
	case "TROPHY_3RD":
		*v = GameAssetKind(speedrun.AssetTrophyThird)
	case "TROPHY_4TH":
		*v = GameAssetKind(speedrun.AssetTrophyFourth)
	case "BACKGROUND":
		*v = GameAssetKind(speedrun.AssetBackground)
	case "FOREGROUND":
		*v = GameAssetKind(speedrun.AssetForeground)
	default:
		return fmt.Errorf("unknown GameAssetKind value %q", s)
	}

	return nil
}

type GameAsset struct {
	speedrun.GameAsset
	Kind GameAssetKind
}
