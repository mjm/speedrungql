package speedrungql

import (
	"encoding/json"
	"time"
)

type EnvelopeResponse struct {
	Data       json.RawMessage `json:"data"`
	Pagination *PageInfo       `json:"pagination"`
}

type PageInfo struct {
	Offset int    `json:"offset"`
	Max    int    `json:"max"`
	Size   int    `json:"size"`
	Links  []Link `json:"links"`
}

type Link struct {
	Rel string `json:"rel"`
	URI string `json:"uri"`
}

type GamesResponse struct {
	Data       []*Game   `json:"data"`
	Pagination *PageInfo `json:"pagination"`
}

type Game struct {
	ID           string      `json:"id"`
	Names        GameNames   `json:"names"`
	Abbreviation string      `json:"abbreviation"`
	Weblink      string      `json:"weblink"`
	ReleaseDate  string      `json:"release-date"`
	Ruleset      GameRuleset `json:"ruleset"`
	Platforms    []string    `json:"platforms"`
}

type GameNames struct {
	International string `json:"international"`
	Japanese      string `json:"japanese"`
	Twitch        string `json:"twitch"`
}

type GameRuleset struct {
	ShowMilliseconds    bool          `json:"show-milliseconds"`
	RequireVerification bool          `json:"require-verification"`
	RequireVideo        bool          `json:"require-video"`
	RunTimes            []GameRunTime `json:"run-times"`
	DefaultRunTime      GameRunTime   `json:"default-time"`
	EmulatorsAllowed    bool          `json:"emulators-allowed"`
}

type GameRunTime string

const (
	RealTime        GameRunTime = "realtime"
	RealTimeNoLoads GameRunTime = "realtime_noloads"
	InGame          GameRunTime = "ingame"
)

type CategoriesResponse struct {
	Data []*Category `json:"data"`
}

type Category struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Weblink       string          `json:"weblink"`
	Type          CategoryType    `json:"type"`
	Rules         string          `json:"rules"`
	Players       CategoryPlayers `json:"players"`
	Miscellaneous bool            `json:"miscellaneous"`
}

type CategoryType string

const (
	CategoryPerGame  CategoryType = "per-game"
	CategoryPerLevel CategoryType = "per-level"
)

type CategoryPlayers struct {
	Type  CategoryPlayersType `json:"type"`
	Value int                 `json:"value"`
}

type CategoryPlayersType string

const (
	PlayersExactly CategoryPlayersType = "exactly"
	PlayersUpTo    CategoryPlayersType = "up-to"
)

type PlatformsResponse struct {
	Data       []*Platform `json:"data"`
	Pagination *PageInfo   `json:"pagination"`
}

type PlatformResponse struct {
	Data Platform `json:"data"`
}

type Platform struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Released int32  `json:"released"`
}

type LeaderboardResponse struct {
	Data *Leaderboard `json:"data"`
}

type Leaderboard struct {
	GameID     string      `json:"game"`
	CategoryID string      `json:"category"`
	Timing     GameRunTime `json:"timing"`
	Runs       []PlacedRun `json:"runs"`
}

type PlacedRun struct {
	Place int  `json:"place"`
	Run   *Run `json:"run"`
}

type Run struct {
	ID         string     `json:"id"`
	GameID     string     `json:"game"`
	CategoryID string     `json:"category"`
	Videos     *RunVideos `json:"videos"`
	Comment    string     `json:"comment"`
	Status     RunStatus  `json:"status"`
	Date       string     `json:"date"`
	Submitted  string     `json:"submitted"`
}

type RunVideos struct {
	Text  string `json:"text"`
	Links []Link `json:"links"`
}

type RunStatus struct {
	Status     RunStatusValue `json:"status"`
	ExaminerID string         `json:"examiner"`
	VerifyDate *time.Time     `json:"verify-date"`
	Reason     string         `json:"reason"`
}

type RunStatusValue string

const (
	RunNew      RunStatusValue = "new"
	RunVerified RunStatusValue = "verified"
	RunRejected RunStatusValue = "rejected"
)

type User struct {
	ID            string        `json:"id"`
	Names         UserNames     `json:"names"`
	Weblink       string        `json:"weblink"`
	NameStyle     UserNameStyle `json:"name-style"`
	Role          UserRole      `json:"role"`
	Signup        *time.Time    `json:"signup"`
	Twitch        *Link         `json:"twitch"`
	Hitbox        *Link         `json:"hitbox"`
	YouTube       *Link         `json:"youtube"`
	Twitter       *Link         `json:"twitter"`
	SpeedRunsLive *Link         `json:"speedrunslive"`
}

type UserNames struct {
	International string `json:"international"`
	Japanese      string `json:"japanese"`
}

type UserNameStyle struct {
	Style     UserNameStyleValue `json:"style"`
	Color     *Color             `json:"color"`
	FromColor *Color             `json:"color-from"`
	ToColor   *Color             `json:"color-to"`
}

type UserNameStyleValue string

const (
	StyleSolid    UserNameStyleValue = "solid"
	StyleGradient UserNameStyleValue = "gradient"
)

type Color struct {
	Light string `json:"light"`
	Dark  string `json:"dark"`
}

type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleBanned     UserRole = "banned"
	RoleTrusted    UserRole = "trusted"
	RoleModerator  UserRole = "moderator"
	RoleAdmin      UserRole = "admin"
	RoleProgrammer UserRole = "programmer"
)
