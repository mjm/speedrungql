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

func (v *Viewer) Users(ctx context.Context, args struct {
	Filter *struct {
		Lookup        *string `filter:"lookup"`
		Name          *string `filter:"name"`
		Twitch        *string `filter:"twitch"`
		Hitbox        *string `filter:"hitbox"`
		Twitter       *string `filter:"twitter"`
		SpeedRunsLive *string `filter:"speedrunslive"`
	}
	Order *struct {
		Field     *UserOrderField
		Direction *speedrun.OrderDirection
	}
	First *int32
	After *Cursor
}) (*UserConnection, error) {
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

	users, pageInfo, err := v.client.ListUsers(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &UserConnection{
		client:   v.client,
		users:    users,
		pageInfo: pageInfo,
	}, nil
}

type UserOrderField string

func (UserOrderField) ImplementsGraphQLType(name string) bool {
	return name == "UserOrderField"
}

func (v *UserOrderField) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("UserOrderField value was not a string")
	}

	switch s {
	case "NAME_INT":
		*v = "name.int"
	case "NAME_JAP":
		*v = "name.jap"
	default:
		*v = UserOrderField(strings.ToLower(s))
	}

	return nil
}

type UserConnection struct {
	client   *speedrun.Client
	users    []*speedrun.User
	pageInfo *speedrun.PageInfo
}

func (uc *UserConnection) Edges() []*UserEdge {
	var edges []*UserEdge
	for _, user := range uc.users {
		edges = append(edges, &UserEdge{
			Node: &User{*user, uc.client},
		})
	}
	return edges
}

func (uc *UserConnection) Nodes() []*User {
	var nodes []*User
	for _, user := range uc.users {
		nodes = append(nodes, &User{*user, uc.client})
	}
	return nodes
}

func (uc *UserConnection) PageInfo() *PageInfo {
	return &PageInfo{uc.pageInfo}
}

type UserEdge struct {
	Node *User
}

func (*UserEdge) Cursor() *Cursor {
	return nil
}

type User struct {
	speedrun.User
	client *speedrun.Client
}

func (u *User) ID() graphql.ID {
	return relay.MarshalID("user", u.User.ID)
}

func (u *User) RawID() string {
	return u.User.ID
}

func (u *User) Name(args struct {
	Variant string
}) *string {
	var s string

	switch args.Variant {
	case "INTERNATIONAL":
		s = u.Names.International
	case "JAPANESE":
		s = u.Names.Japanese
	}

	if s == "" {
		return nil
	}
	return &s
}

func (u *User) NameStyle() *UserNameStyle {
	return &UserNameStyle{u.User.NameStyle}
}

func (u *User) Role() UserRole {
	return UserRole(u.User.Role)
}

func (u *User) Signup() *graphql.Time {
	if u.User.Signup == nil {
		return nil
	}

	return &graphql.Time{Time: *u.User.Signup}
}

func (u *User) Twitch() *Link {
	return wrapLink(u.User.Twitch)
}

func (u *User) Hitbox() *Link {
	return wrapLink(u.User.Hitbox)
}

func (u *User) YouTube() *Link {
	return wrapLink(u.User.YouTube)
}

func (u *User) Twitter() *Link {
	return wrapLink(u.User.Twitter)
}

func (u *User) SpeedRunsLive() *Link {
	return wrapLink(u.User.SpeedRunsLive)
}

func (u *User) Runs(ctx context.Context, args FetchRunsArgs) (*RunConnection, error) {
	if args.Filter != nil && args.Filter.User != nil {
		return nil, errors.New("cannot filter runs by user when reading from a specific user")
	}

	return fetchRunConnection(ctx, u.client, args, speedrun.WithFilter("user", u.User.ID))
}

func (u *User) PersonalBests(ctx context.Context) ([]*PlacedRun, error) {
	bests, err := u.client.ListUserPersonalBests(ctx, u.User.ID)
	if err != nil {
		return nil, err
	}

	var res []*PlacedRun
	for _, run := range bests {
		res = append(res, &PlacedRun{run, u.client})
	}
	return res, nil
}

func (u *User) ModeratedGames(ctx context.Context, args FetchGamesArgs) (*GameConnection, error) {
	if args.Filter != nil && args.Filter.Moderator != nil {
		return nil, errors.New("cannot filter games by moderator when reading from a specific user")
	}

	return fetchGameConnection(ctx, u.client, args, speedrun.WithFilter("moderator", u.User.ID))
}

type UserNames struct {
	speedrun.UserNames
}

func (un *UserNames) Japanese() *string {
	if un.UserNames.Japanese == "" {
		return nil
	}
	return &un.UserNames.Japanese
}

type UserNameStyle struct {
	speedrun.UserNameStyle
}

func (uns *UserNameStyle) ToSolidUserNameStyle() (*SolidUserNameStyle, bool) {
	if uns.Style != speedrun.StyleSolid {
		return nil, false
	}

	return &SolidUserNameStyle{uns.UserNameStyle}, true
}

func (uns *UserNameStyle) ToGradientUserNameStyle() (*GradientUserNameStyle, bool) {
	if uns.Style != speedrun.StyleGradient {
		return nil, false
	}

	return &GradientUserNameStyle{uns.UserNameStyle}, true
}

type SolidUserNameStyle struct {
	speedrun.UserNameStyle
}

type GradientUserNameStyle struct {
	speedrun.UserNameStyle
}

type UserRole speedrun.UserRole

func (UserRole) ImplementsGraphQLType(name string) bool {
	return name == "UserRole"
}

func (v UserRole) String() string {
	return strings.ToUpper(string(v))
}

func (v *UserRole) UnmarshalGraphQL(input interface{}) error {
	s, ok := input.(string)
	if !ok {
		return errors.New("UserRole value was not a string")
	}

	switch s {
	case "USER":
		*v = UserRole(speedrun.RoleUser)
	case "BANNED":
		*v = UserRole(speedrun.RoleBanned)
	case "TRUSTED":
		*v = UserRole(speedrun.RoleTrusted)
	case "MODERATOR":
		*v = UserRole(speedrun.RoleModerator)
	case "ADMIN":
		*v = UserRole(speedrun.RoleAdmin)
	case "PROGRAMMER":
		*v = UserRole(speedrun.RoleProgrammer)
	default:
		return fmt.Errorf("unknown UserRole value %q", s)
	}

	return nil
}
