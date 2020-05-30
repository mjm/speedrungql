package resolvers

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

type Resolvers struct {
	client *speedrungql.Client
}

func New(baseURL string) *Resolvers {
	return &Resolvers{
		client: speedrungql.NewClient(baseURL),
	}
}

func (r *Resolvers) Viewer() *Viewer {
	return &Viewer{client: r.client}
}

type Viewer struct {
	client *speedrungql.Client
}

func (v *Viewer) Platforms(ctx context.Context, args struct {
	Order *struct {
		Field     *string
		Direction *speedrungql.OrderDirection
	}
	First *int32
	After *Cursor
}) (*PlatformConnection, error) {
	opts := []speedrungql.FetchOption{
		speedrungql.WithOrder(args.Order.Field, args.Order.Direction),
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

	plats, pageInfo, err := v.client.ListPlatforms(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &PlatformConnection{plats, pageInfo}, nil
}

type PlatformConnection struct {
	platforms []*speedrungql.Platform
	pageInfo  *speedrungql.PageInfo
}

func (pc *PlatformConnection) Edges() []*PlatformEdge {
	var edges []*PlatformEdge
	for _, p := range pc.platforms {
		edges = append(edges, &PlatformEdge{
			Node: &Platform{*p},
		})
	}
	return edges
}

func (pc *PlatformConnection) Nodes() []*Platform {
	var nodes []*Platform
	for _, p := range pc.platforms {
		nodes = append(nodes, &Platform{*p})
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
	speedrungql.Platform
}

func (p *Platform) ID() graphql.ID {
	return relay.MarshalID("platform", p.Platform.ID)
}
