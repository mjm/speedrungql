package resolvers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/speedrungql"
)

type User struct {
	speedrungql.User
}

func (u *User) ID() graphql.ID {
	return relay.MarshalID("user", u.User.ID)
}

func (u *User) RawID() string {
	return u.User.ID
}

func (u *User) Names() *UserNames {
	return &UserNames{u.User.Names}
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

type UserNames struct {
	speedrungql.UserNames
}

func (un *UserNames) Japanese() *string {
	if un.UserNames.Japanese == "" {
		return nil
	}
	return &un.UserNames.Japanese
}

type UserNameStyle struct {
	speedrungql.UserNameStyle
}

func (uns *UserNameStyle) ToSolidUserNameStyle() (*SolidUserNameStyle, bool) {
	if uns.Style != speedrungql.StyleSolid {
		return nil, false
	}

	return &SolidUserNameStyle{uns.UserNameStyle}, true
}

func (uns *UserNameStyle) ToGradientUserNameStyle() (*GradientUserNameStyle, bool) {
	if uns.Style != speedrungql.StyleGradient {
		return nil, false
	}

	return &GradientUserNameStyle{uns.UserNameStyle}, true
}

type SolidUserNameStyle struct {
	speedrungql.UserNameStyle
}

type GradientUserNameStyle struct {
	speedrungql.UserNameStyle
}

type UserRole speedrungql.UserRole

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
		*v = UserRole(speedrungql.RoleUser)
	case "BANNED":
		*v = UserRole(speedrungql.RoleBanned)
	case "TRUSTED":
		*v = UserRole(speedrungql.RoleTrusted)
	case "MODERATOR":
		*v = UserRole(speedrungql.RoleModerator)
	case "ADMIN":
		*v = UserRole(speedrungql.RoleAdmin)
	case "PROGRAMMER":
		*v = UserRole(speedrungql.RoleProgrammer)
	default:
		return fmt.Errorf("unknown UserRole value %q", s)
	}

	return nil
}
