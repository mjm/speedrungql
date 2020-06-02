package resolvers

import (
	"context"
	"errors"
	"strings"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

func (v *Viewer) Genres(ctx context.Context, args struct {
	Order *struct {
		Field     *GenreOrderField
		Direction *speedrun.OrderDirection
	}
	First *int32
	After *Cursor
}) (*GenreConnection, error) {
	var opts []speedrun.FetchOption
	if args.Order != nil {
		opts = append(opts, speedrun.WithOrder((*string)(args.Order.Field), args.Order.Direction))
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

	genres, pageInfo, err := v.client.ListGenres(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &GenreConnection{
		client:   v.client,
		genres:   genres,
		pageInfo: pageInfo,
	}, nil
}

type GenreOrderField string

func (GenreOrderField) ImplementsGraphQLType(name string) bool {
	return name == "GenreOrderField"
}

func (v *GenreOrderField) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("GenreOrderField value was not a string")
	}

	*v = GenreOrderField(strings.ToLower(s))
	return nil
}

type GenreConnection struct {
	client   *speedrun.Client
	genres   []*speedrun.Genre
	pageInfo *speedrun.PageInfo
}

func (gc *GenreConnection) Edges() []*GenreEdge {
	var edges []*GenreEdge
	for _, g := range gc.genres {
		edges = append(edges, &GenreEdge{
			Node: &Genre{*g, gc.client},
		})
	}
	return edges
}

func (gc *GenreConnection) Nodes() []*Genre {
	var nodes []*Genre
	for _, g := range gc.genres {
		nodes = append(nodes, &Genre{*g, gc.client})
	}
	return nodes
}

func (gc *GenreConnection) PageInfo() *PageInfo {
	return &PageInfo{gc.pageInfo}
}

type GenreEdge struct {
	Node *Genre
}

func (GenreEdge) Cursor() *Cursor {
	return nil
}

type Genre struct {
	speedrun.Genre
	client *speedrun.Client
}

func (g *Genre) ID() graphql.ID {
	return relay.MarshalID("genre", g.Genre.ID)
}

func (g *Genre) RawID() string {
	return g.Genre.ID
}

func (g *Genre) Games(ctx context.Context, args FetchGamesArgs) (*GameConnection, error) {
	if args.Filter != nil && args.Filter.Genre != nil {
		return nil, errors.New("cannot filter games by genre when reading from a specific genre")
	}

	return fetchGameConnection(ctx, g.client, args, speedrun.WithFilter("genre", g.Genre.ID))
}
