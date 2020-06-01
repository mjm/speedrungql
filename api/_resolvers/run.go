package resolvers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql/speedrun"
)

func (v *Viewer) Runs(ctx context.Context, args struct {
	Filter *struct {
		User     *graphql.ID     `filter:"user"`
		Guest    *string         `filter:"guest"`
		Examiner *graphql.ID     `filter:"examiner"`
		Game     *graphql.ID     `filter:"game"`
		Level    *graphql.ID     `filter:"level"`
		Category *graphql.ID     `filter:"category"`
		Platform *graphql.ID     `filter:"platform"`
		Region   *graphql.ID     `filter:"region"`
		Emulated *bool           `filter:"emulated"`
		Status   *RunStatusValue `filter:"status"`
	}
	Order *struct {
		Field     *RunOrderField
		Direction *speedrun.OrderDirection
	}
	First *int32
	After *Cursor
}) (*RunConnection, error) {
	var opts []speedrun.FetchOption
	if args.Order != nil {
		opts = append(opts, speedrun.WithOrder((*string)(args.Order.Field), args.Order.Direction))
	}
	if args.Filter != nil {
		opts = append(opts, speedrun.WithFilters(*args.Filter))
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

	runs, pageInfo, err := v.client.ListRuns(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &RunConnection{
		client:   v.client,
		runs:     runs,
		pageInfo: pageInfo,
	}, nil
}

type RunOrderField string

func (RunOrderField) ImplementsGraphQLType(name string) bool {
	return name == "RunOrderField"
}

func (v *RunOrderField) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("RunOrderField value was not a string")
	}

	switch s {
	case "VERIFY_DATE":
		*v = "verify-date"
	default:
		*v = RunOrderField(strings.ToLower(s))
	}

	return nil
}

type RunConnection struct {
	client   *speedrun.Client
	runs     []*speedrun.Run
	pageInfo *speedrun.PageInfo
}

func (rc *RunConnection) Edges() []*RunEdge {
	var edges []*RunEdge
	for _, r := range rc.runs {
		edges = append(edges, &RunEdge{
			Node: &Run{*r, rc.client},
		})
	}
	return edges
}

func (rc *RunConnection) Nodes() []*Run {
	var nodes []*Run
	for _, r := range rc.runs {
		nodes = append(nodes, &Run{*r, rc.client})
	}
	return nodes
}

func (rc *RunConnection) PageInfo() *PageInfo {
	return &PageInfo{rc.pageInfo}
}

type RunEdge struct {
	Node *Run
}

func (e *RunEdge) Cursor() *Cursor {
	return nil
}

type Run struct {
	speedrun.Run
	client *speedrun.Client
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

	return &Category{*c, r.client}, nil
}

func (r *Run) Level(ctx context.Context) (*Level, error) {
	if r.LevelID == "" {
		return nil, nil
	}

	l, err := r.client.GetLevel(ctx, r.LevelID)
	if err != nil {
		return nil, err
	}

	if l == nil {
		return nil, nil
	}

	return &Level{*l, r.client}, nil
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

func (r *Run) Splits() *Link {
	return wrapLink(r.Run.Splits)
}

func (r *Run) Time(args struct {
	Timing *GameRunTime
}) *float64 {
	var t float64

	if args.Timing == nil {
		t = r.Times.Primary
	} else {
		switch speedrun.GameRunTime(*args.Timing) {
		case speedrun.RealTime:
			t = r.Times.RealTime
		case speedrun.RealTimeNoLoads:
			t = r.Times.RealTimeNoLoads
		case speedrun.InGame:
			t = r.Times.InGame
		}
	}

	if t == 0 {
		return nil
	}
	return &t
}

func (r *Run) Values(ctx context.Context) ([]*VariableValue, error) {
	var vals []*VariableValue

	for varID, valID := range r.Run.Values {
		v, err := r.client.GetVariable(ctx, varID)
		if err != nil {
			return nil, err
		}

		if v == nil {
			continue
		}

		varResolver := &Variable{*v, r.client}
		val, ok := v.Values.Values[valID]
		if !ok {
			continue
		}

		vals = append(vals, &VariableValue{val, valID, varResolver})
	}

	return vals, nil
}

func (r *Run) Value(ctx context.Context, args struct {
	VariableID graphql.ID
}) (*VariableValue, error) {
	varID := string(args.VariableID)
	valID, ok := r.Run.Values[varID]
	if !ok {
		return nil, nil
	}

	v, err := r.client.GetVariable(ctx, varID)
	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, nil
	}

	varResolver := &Variable{*v, r.client}
	val, ok := v.Values.Values[valID]
	if !ok {
		return nil, nil
	}

	return &VariableValue{val, valID, varResolver}, nil
}

type RunStatus struct {
	speedrun.RunStatus
	client *speedrun.Client
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

type RunStatusValue speedrun.RunStatusValue

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
		*v = RunStatusValue(speedrun.RunNew)
	case "VERIFIED":
		*v = RunStatusValue(speedrun.RunVerified)
	case "REJECTED":
		*v = RunStatusValue(speedrun.RunRejected)
	default:
		return fmt.Errorf("unknown RunStatusValue value %q", s)
	}

	return nil
}

type RunVideos struct {
	speedrun.RunVideos
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
	speedrun.RunPlayer
	client *speedrun.Client
}

func (rp *RunPlayer) ToUserRunPlayer() (*UserRunPlayer, bool) {
	if rp.Rel != speedrun.PlayerUser {
		return nil, false
	}

	return &UserRunPlayer{rp.RunPlayer, rp.client}, true
}

func (rp *RunPlayer) ToGuestRunPlayer() (*GuestRunPlayer, bool) {
	if rp.Rel != speedrun.PlayerGuest {
		return nil, false
	}

	return &GuestRunPlayer{rp.RunPlayer}, true
}

type UserRunPlayer struct {
	speedrun.RunPlayer
	client *speedrun.Client
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
	speedrun.RunPlayer
}
