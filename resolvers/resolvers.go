package resolvers

import (
	"context"
	"encoding/json"

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
		Direction *string
	}
	First *int32
	After *Cursor
}) (*PlatformConnection, error) {
	u := v.client.BaseURL + "/platforms"
	res, err := v.client.HTTPClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp speedrungql.PlatformsResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &PlatformConnection{res: &resp}, nil
}

type PlatformConnection struct {
	res *speedrungql.PlatformsResponse
}

func (pc *PlatformConnection) Edges() []*PlatformEdge {
	var edges []*PlatformEdge
	for _, p := range pc.res.Data {
		edges = append(edges, &PlatformEdge{
			Node: &Platform{p},
		})
	}
	return edges
}

func (pc *PlatformConnection) Nodes() []*Platform {
	var nodes []*Platform
	for _, p := range pc.res.Data {
		nodes = append(nodes, &Platform{p})
	}
	return nodes
}

func (pc *PlatformConnection) PageInfo() PageInfo {
	return PageInfo{}
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
