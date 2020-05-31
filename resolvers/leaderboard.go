package resolvers

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

func (v *Viewer) Leaderboard(ctx context.Context, args struct {
	Game      graphql.ID
	Category  graphql.ID
	Level     *graphql.ID
	Variables *[]struct {
		ID    graphql.ID
		Value graphql.ID
	}
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

	var opts []speedrungql.FetchOption
	if args.Variables != nil {
		for _, v := range *args.Variables {
			var varID string
			if err := relay.UnmarshalSpec(v.ID, &varID); err != nil {
				return nil, err
			}

			opts = append(opts, speedrungql.WithFilter("var-"+varID, string(v.Value)))
		}
	}

	lb, err := v.client.GetLeaderboard(ctx, gameID, categoryID, levelID, opts...)
	if err != nil {
		return nil, err
	}

	if lb == nil {
		return nil, nil
	}

	return &Leaderboard{*lb, v.client}, nil
}

type Leaderboard struct {
	speedrungql.Leaderboard
	client *speedrungql.Client
}

func (l *Leaderboard) Game(ctx context.Context) (*Game, error) {
	g, err := l.client.GetGame(ctx, l.GameID)
	if err != nil {
		return nil, err
	}

	if g == nil {
		return nil, nil
	}

	return &Game{*g, l.client}, nil
}

func (l *Leaderboard) Category(ctx context.Context) (*Category, error) {
	c, err := l.client.GetCategory(ctx, l.CategoryID)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	return &Category{*c, l.client}, nil
}

func (l *Leaderboard) Level(ctx context.Context) (*Level, error) {
	if l.LevelID == "" {
		return nil, nil
	}

	lev, err := l.client.GetLevel(ctx, l.LevelID)
	if err != nil {
		return nil, err
	}

	if lev == nil {
		return nil, nil
	}

	return &Level{*lev, l.client}, nil
}

func (l *Leaderboard) Timing() GameRunTime {
	return GameRunTime(l.Leaderboard.Timing)
}

func (l *Leaderboard) Runs(args struct {
	First int32
}) []*PlacedRun {
	max := int32(len(l.Leaderboard.Runs))
	if max > args.First {
		max = args.First
	}

	var runs []*PlacedRun
	for _, r := range l.Leaderboard.Runs[:max] {
		runs = append(runs, &PlacedRun{r, l.client})
	}
	return runs
}

type PlacedRun struct {
	speedrungql.PlacedRun
	client *speedrungql.Client
}

func (pr *PlacedRun) Place() int32 {
	return int32(pr.PlacedRun.Place)
}

func (pr *PlacedRun) Run() *Run {
	return &Run{*pr.PlacedRun.Run, pr.client}
}
