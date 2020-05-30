package resolvers

import (
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

func (u *User) Names() *UserNames {
	return &UserNames{u.User.Names}
}

func (u *User) NameStyle() *UserNameStyle {
	return &UserNameStyle{u.User.NameStyle}
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
	if uns.Style != "solid" {
		return nil, false
	}

	return &SolidUserNameStyle{uns.UserNameStyle}, true
}

func (uns *UserNameStyle) ToGradientUserNameStyle() (*GradientUserNameStyle, bool) {
	if uns.Style != "gradient" {
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
