package resolvers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

type Run struct {
	speedrungql.Run
	client *speedrungql.Client
}

func (r *Run) ID() graphql.ID {
	return relay.MarshalID("run", r.Run.ID)
}

func (r *Run) RawID() string {
	return r.Run.ID
}

func (r *Run) Game(ctx context.Context) (*Game, error) {
	g, err := r.client.GetGame(ctx, r.GameID)
	if err != nil {
		return nil, err
	}

	if g == nil {
		return nil, nil
	}

	return &Game{*g, r.client}, nil
}

func (r *Run) Category(ctx context.Context) (*Category, error) {
	c, err := r.client.GetCategory(ctx, r.CategoryID)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	return &Category{*c}, nil
}

func (r *Run) Videos() *RunVideos {
	if r.Run.Videos == nil {
		return nil
	}

	return &RunVideos{*r.Run.Videos}
}

func (r *Run) Status() *RunStatus {
	return &RunStatus{r.Run.Status, r.client}
}

func (r *Run) Date() *string {
	if r.Run.Date == "" {
		return nil
	}
	return &r.Run.Date
}

func (r *Run) Submitted() *string {
	if r.Run.Submitted == "" {
		return nil
	}
	return &r.Run.Submitted
}

func (r *Run) Players() []*RunPlayer {
	var rps []*RunPlayer
	for _, rp := range r.Run.Players {
		rps = append(rps, &RunPlayer{rp, r.client})
	}
	return rps
}

type RunStatus struct {
	speedrungql.RunStatus
	client *speedrungql.Client
}

func (rs *RunStatus) Status() RunStatusValue {
	return RunStatusValue(rs.RunStatus.Status)
}

func (rs *RunStatus) Examiner(ctx context.Context) (*User, error) {
	if rs.ExaminerID == "" {
		return nil, nil
	}

	user, err := rs.client.GetUser(ctx, rs.ExaminerID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &User{*user}, nil
}

func (rs *RunStatus) VerifyDate() *graphql.Time {
	if rs.RunStatus.VerifyDate == nil {
		return nil
	}

	return &graphql.Time{Time: *rs.RunStatus.VerifyDate}
}

func (rs *RunStatus) Reason() *string {
	if rs.RunStatus.Reason == "" {
		return nil
	}
	return &rs.RunStatus.Reason
}

type RunStatusValue speedrungql.RunStatusValue

func (RunStatusValue) ImplementsGraphQLType(name string) bool {
	return name == "RunStatusValue"
}

func (v RunStatusValue) String() string {
	return strings.ToUpper(string(v))
}

func (v *RunStatusValue) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("RunStatusValue value was not a string")
	}

	switch s {
	case "NEW":
		*v = RunStatusValue(speedrungql.RunNew)
	case "VERIFIED":
		*v = RunStatusValue(speedrungql.RunVerified)
	case "REJECTED":
		*v = RunStatusValue(speedrungql.RunRejected)
	default:
		return fmt.Errorf("unknown RunStatusValue value %q", s)
	}

	return nil
}

type RunVideos struct {
	speedrungql.RunVideos
}

func (rv *RunVideos) Text() *string {
	if rv.RunVideos.Text == "" {
		return nil
	}
	return &rv.RunVideos.Text
}

func (rv *RunVideos) Links() []*Link {
	var links []*Link
	for _, l := range rv.RunVideos.Links {
		links = append(links, &Link{l})
	}
	return links
}

type RunPlayer struct {
	speedrungql.RunPlayer
	client *speedrungql.Client
}

func (rp *RunPlayer) ToUserRunPlayer() (*UserRunPlayer, bool) {
	if rp.Rel != speedrungql.PlayerUser {
		return nil, false
	}

	return &UserRunPlayer{rp.RunPlayer, rp.client}, true
}

func (rp *RunPlayer) ToGuestRunPlayer() (*GuestRunPlayer, bool) {
	if rp.Rel != speedrungql.PlayerGuest {
		return nil, false
	}

	return &GuestRunPlayer{rp.RunPlayer}, true
}

type UserRunPlayer struct {
	speedrungql.RunPlayer
	client *speedrungql.Client
}

func (urp *UserRunPlayer) User(ctx context.Context) (*User, error) {
	user, err := urp.client.GetUser(ctx, urp.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &User{*user}, nil
}

type GuestRunPlayer struct {
	speedrungql.RunPlayer
}
