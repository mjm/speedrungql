package resolvers

import (
	"github.com/mjm/speedrungql/speedrun"
)

type Resolvers struct {
	client *speedrun.Client
}

func New(baseURL string) *Resolvers {
	return &Resolvers{
		client: speedrun.NewClient(baseURL),
	}
}

func (r *Resolvers) Viewer() *Viewer {
	return &Viewer{client: r.client}
}

type Viewer struct {
	client *speedrun.Client
}
