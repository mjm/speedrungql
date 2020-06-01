package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

type Category struct {
	speedrun.Category
	client *speedrun.Client
}

func (c *Category) ID() graphql.ID {
	return relay.MarshalID("category", c.Category.ID)
}

func (c *Category) RawID() string {
	return c.Category.ID
}

func (c *Category) Game(ctx context.Context) (*Game, error) {
	gameURI := speedrun.FindLink(c.Links, "game")
	if gameURI == "" {
		return nil, nil
	}

	game, err := c.client.GetGame(ctx, gameURI)
	if err != nil {
		return nil, err
	}

	if game == nil {
		return nil, nil
	}
	return &Game{*game, c.client}, nil
}

func (c *Category) Type() CategoryType {
	return CategoryType(c.Category.Type)
}

func (c *Category) Players() *CategoryPlayers {
	return &CategoryPlayers{c.Category.Players}
}

func (c *Category) Variables(ctx context.Context) ([]*Variable, error) {
	vs, err := c.client.ListCategoryVariables(ctx, c.Category.ID)
	if err != nil {
		return nil, err
	}

	var res []*Variable
	for _, v := range vs {
		res = append(res, &Variable{*v, c.client})
	}
	return res, nil
}

func (c *Category) Runs(ctx context.Context, args struct {
	Filter *struct {
		User     *graphql.ID     `filter:"user"`
		Guest    *string         `filter:"guest"`
		Examiner *graphql.ID     `filter:"examiner"`
		Game     *graphql.ID     `filter:"game"`
		Level    *graphql.ID     `filter:"level"`
		Category *graphql.ID     `filter:"category"`
		Platform *graphql.ID     `filter:"platform"`
		Region   *graphql.ID     `filter:"region"`
		Emulated *bool           `filter:"emulated"`
		Status   *RunStatusValue `filter:"status"`
	}
	Order *struct {
		Field     *RunOrderField
		Direction *speedrun.OrderDirection
	}
	First *int32
	After *Cursor
}) (*RunConnection, error) {
	if args.Filter != nil && args.Filter.Category != nil {
		return nil, errors.New("cannot filter runs by category when reading from a specific category")
	}

	opts := []speedrun.FetchOption{
		speedrun.WithFilter("category", c.Category.ID),
	}

	if args.Order != nil {
		opts = append(opts, speedrun.WithOrder((*string)(args.Order.Field), args.Order.Direction))
	}
	if args.Filter != nil {
		opts = append(opts, speedrun.WithFilters(*args.Filter))
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

	runs, pageInfo, err := c.client.ListRuns(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &RunConnection{
		client:   c.client,
		runs:     runs,
		pageInfo: pageInfo,
	}, nil
}

type CategoryType speedrun.CategoryType

func (CategoryType) ImplementsGraphQLType(name string) bool {
	return name == "CategoryType"
}

func (v CategoryType) String() string {
	switch speedrun.CategoryType(v) {
	case speedrun.CategoryPerGame:
		return "PER_GAME"
	case speedrun.CategoryPerLevel:
		return "PER_LEVEL"
	default:
		return ""
	}
}

func (v *CategoryType) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("CategoryType value was not a string")
	}

	switch s {
	case "PER_GAME":
		*v = CategoryType(speedrun.CategoryPerGame)
	case "PER_LEVEL":
		*v = CategoryType(speedrun.CategoryPerLevel)
	default:
		return fmt.Errorf("unknown CategoryType value %q", s)
	}

	return nil
}

type CategoryPlayers struct {
	speedrun.CategoryPlayers
}

func (c *CategoryPlayers) Type() CategoryPlayersType {
	return CategoryPlayersType(c.CategoryPlayers.Type)
}

func (c *CategoryPlayers) Value() int32 {
	return int32(c.CategoryPlayers.Value)
}

type CategoryPlayersType speedrun.CategoryPlayersType

func (CategoryPlayersType) ImplementsGraphQLType(name string) bool {
	return name == "CategoryPlayersType"
}

func (v CategoryPlayersType) String() string {
	switch speedrun.CategoryPlayersType(v) {
	case speedrun.PlayersExactly:
		return "EXACTLY"
	case speedrun.PlayersUpTo:
		return "UP_TO"
	default:
		return ""
	}
}

func (v *CategoryPlayersType) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("CategoryPlayersType value was not a string")
	}

	switch s {
	case "EXACTLY":
		*v = CategoryPlayersType(speedrun.PlayersExactly)
	case "UP_TO":
		*v = CategoryPlayersType(speedrun.PlayersUpTo)
	default:
		return fmt.Errorf("unknown CategoryPlayersType value %q", s)
	}

	return nil
}
