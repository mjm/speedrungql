package resolvers

import (
	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

type Region struct {
	speedrun.Region
}

func (r *Region) ID() graphql.ID {
	return relay.MarshalID("region", r.Region.ID)
}

func (r *Region) RawID() string {
	return r.Region.ID
}
