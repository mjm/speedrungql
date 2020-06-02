package resolvers

import (
	"context"
	"errors"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

type Region struct {
	speedrun.Region
	client *speedrun.Client
}

func (r *Region) ID() graphql.ID {
	return relay.MarshalID("region", r.Region.ID)
}

func (r *Region) RawID() string {
	return r.Region.ID
}

func (r *Region) Games(ctx context.Context, args FetchGamesArgs) (*GameConnection, error) {
	if args.Filter != nil && args.Filter.Region != nil {
		return nil, errors.New("cannot filter games by region when reading from a specific region")
	}

	return fetchGameConnection(ctx, r.client, args, speedrun.WithFilter("region", r.Region.ID))
}
