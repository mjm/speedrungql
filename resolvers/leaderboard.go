package resolvers

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

func (v *Viewer) Leaderboard(ctx context.Context, args struct {
	Game     graphql.ID
	Category graphql.ID
	Level    *graphql.ID
}) (*Leaderboard, error) {
	var gameID string
	if err := relay.UnmarshalSpec(args.Game, &gameID); err != nil {
		return nil, err
	}
	var categoryID string
	if err := relay.UnmarshalSpec(args.Category, &categoryID); err != nil {
		return nil, err
	}
	var levelID *string
	if args.Level != nil {
		levelID = new(string)
		if err := relay.UnmarshalSpec(*args.Level, levelID); err != nil {
			return nil, err
		}
	}

	lb, err := v.client.GetLeaderboard(ctx, gameID, categoryID, levelID)
	if err != nil {
		return nil, err
	}

	if lb == nil {
		return nil, nil
	}

	return &Leaderboard{*lb}, nil
}

type Leaderboard struct {
	speedrungql.Leaderboard
}

func (l *Leaderboard) Game(ctx context.Context) (*Game, error) {
	return nil, nil
}

func (l *Leaderboard) Category(ctx context.Context) (*Category, error) {
	return nil, nil
}

func (l *Leaderboard) Runs() []*PlacedRun {
	var runs []*PlacedRun
	for _, r := range l.Leaderboard.Runs {
		runs = append(runs, &PlacedRun{r})
	}
	return runs
}

type PlacedRun struct {
	speedrungql.PlacedRun
}

func (pr *PlacedRun) Place() int32 {
	return int32(pr.PlacedRun.Place)
}

func (pr *PlacedRun) Run() *Run {
	return &Run{*pr.PlacedRun.Run}
}
