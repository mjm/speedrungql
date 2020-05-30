package resolvers

import (
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
