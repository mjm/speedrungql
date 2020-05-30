package resolvers

import (
	"context"

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

func (r *Run) Game(ctx context.Context) (*Game, error) {
	return nil, nil
}

func (r *Run) Category(ctx context.Context) (*Category, error) {
	return nil, nil
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

type RunStatus struct {
	speedrungql.RunStatus
	client *speedrungql.Client
}

func (rs *RunStatus) Status() string {
	switch rs.RunStatus.Status {
	case "new":
		return "NEW"
	case "verified":
		return "VERIFIED"
	case "rejected":
		return "REJECTED"
	default:
		return ""
	}
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
