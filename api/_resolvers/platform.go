package resolvers

import (
	"context"
	"errors"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

func (v *Viewer) Platforms(ctx context.Context, args struct {
	Order *struct {
		Field     *string
		Direction *speedrun.OrderDirection
	}
	First *int32
	After *Cursor
}) (*PlatformConnection, error) {
	var opts []speedrun.FetchOption
	if args.Order != nil {
		opts = append(opts, speedrun.WithOrder(args.Order.Field, args.Order.Direction))
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

	plats, pageInfo, err := v.client.ListPlatforms(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &PlatformConnection{v.client, plats, pageInfo}, nil
}

type PlatformConnection struct {
	client    *speedrun.Client
	platforms []*speedrun.Platform
	pageInfo  *speedrun.PageInfo
}

func (pc *PlatformConnection) Edges() []*PlatformEdge {
	var edges []*PlatformEdge
	for _, p := range pc.platforms {
		edges = append(edges, &PlatformEdge{
			Node: &Platform{*p, pc.client},
		})
	}
	return edges
}

func (pc *PlatformConnection) Nodes() []*Platform {
	var nodes []*Platform
	for _, p := range pc.platforms {
		nodes = append(nodes, &Platform{*p, pc.client})
	}
	return nodes
}

func (pc *PlatformConnection) PageInfo() *PageInfo {
	return &PageInfo{pc.pageInfo}
}

type PlatformEdge struct {
	Node *Platform
}

func (pe *PlatformEdge) Cursor() *Cursor {
	return nil
}

type Platform struct {
	speedrun.Platform
	client *speedrun.Client
}

func (p *Platform) ID() graphql.ID {
	return relay.MarshalID("platform", p.Platform.ID)
}

func (p *Platform) RawID() string {
	return p.Platform.ID
}

func (p *Platform) Games(ctx context.Context, args FetchGamesArgs) (*GameConnection, error) {
	if args.Filter != nil && args.Filter.Platform != nil {
		return nil, errors.New("cannot filter games by platform when reading from a specific platform")
	}

	return fetchGameConnection(ctx, p.client, args, speedrun.WithFilter("platform", p.Platform.ID))
}
