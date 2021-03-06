package speedrun

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

type EnginesResponse struct {
	Data       []*Engine `json:"data"`
	Pagination *PageInfo `json:"pagination"`
}

type Engine struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GamesResponse struct {
	Data       []*Game   `json:"data"`
	Pagination *PageInfo `json:"pagination"`
}

type Game struct {
	ID           string                       `json:"id"`
	Names        GameNames                    `json:"names"`
	Abbreviation string                       `json:"abbreviation"`
	Weblink      string                       `json:"weblink"`
	ReleaseDate  string                       `json:"release-date"`
	Ruleset      GameRuleset                  `json:"ruleset"`
	Platforms    []string                     `json:"platforms"`
	Regions      []string                     `json:"regions"`
	Genres       []string                     `json:"genres"`
	Engines      []string                     `json:"engines"`
	Moderators   map[string]GameModeratorRole `json:"moderators"`
	Assets       map[GameAssetKind]*GameAsset `json:"assets"`
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

type GameModeratorRole string

const (
	Moderator      GameModeratorRole = "moderator"
	SuperModerator GameModeratorRole = "super-moderator"
)

type GameAssetKind string

const (
	AssetLogo         GameAssetKind = "logo"
	AssetCoverTiny    GameAssetKind = "cover-tiny"
	AssetCoverSmall   GameAssetKind = "cover-small"
	AssetCoverMedium  GameAssetKind = "cover-medium"
	AssetCoverLarge   GameAssetKind = "cover-large"
	AssetIcon         GameAssetKind = "icon"
	AssetTrophyFirst  GameAssetKind = "trophy-1st"
	AssetTrophySecond GameAssetKind = "trophy-2nd"
	AssetTrophyThird  GameAssetKind = "trophy-3rd"
	AssetTrophyFourth GameAssetKind = "trophy-4th"
	AssetBackground   GameAssetKind = "background"
	AssetForeground   GameAssetKind = "foreground"
)

type GameAsset struct {
	URI    string `json:"uri"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

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
	Links         []Link          `json:"links"`
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

type GenresResponse struct {
	Data       []*Genre  `json:"data"`
	Pagination *PageInfo `json:"pagination"`
}

type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LevelsResponse struct {
	Data []*Level `json:"data"`
}

type Level struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Weblink string `json:"weblink"`
	Rules   string `json:"rules"`
	Links   []Link `json:"link"`
}

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
	LevelID    string      `json:"level"`
	Timing     GameRunTime `json:"timing"`
	Runs       []PlacedRun `json:"runs"`
}

type PlacedRunsResponse struct {
	Data []PlacedRun `json:"data"`
}

type PlacedRun struct {
	Place int  `json:"place"`
	Run   *Run `json:"run"`
}

type RunsResponse struct {
	Data       []*Run    `json:"data"`
	Pagination *PageInfo `json:"pagination"`
}

type Run struct {
	ID         string            `json:"id"`
	GameID     string            `json:"game"`
	CategoryID string            `json:"category"`
	LevelID    string            `json:"level"`
	Videos     *RunVideos        `json:"videos"`
	Comment    string            `json:"comment"`
	Status     RunStatus         `json:"status"`
	Date       string            `json:"date"`
	Submitted  string            `json:"submitted"`
	Players    []RunPlayer       `json:"players"`
	Times      RunTimes          `json:"times"`
	Splits     *Link             `json:"splits"`
	Values     map[string]string `json:"values"`
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

type RunPlayer struct {
	Rel  RunPlayerRel `json:"rel"`
	ID   string       `json:"id"`
	Name string       `json:"name"`
	URI  string       `json:"uri"`
}

type RunPlayerRel string

const (
	PlayerUser  RunPlayerRel = "user"
	PlayerGuest RunPlayerRel = "guest"
)

type RunTimes struct {
	Primary         float64 `json:"primary_t"`
	RealTime        float64 `json:"realtime_t"`
	RealTimeNoLoads float64 `json:"realtime_noloads_t"`
	InGame          float64 `json:"ingame_t"`
}

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UsersResponse struct {
	Data       []*User   `json:"data"`
	Pagination *PageInfo `json:"pagination"`
}

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

type VariablesResponse struct {
	Data []*Variable `json:"data"`
}

type Variable struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	CategoryID    string         `json:"category"`
	Scope         VariableScope  `json:"scope"`
	Mandatory     bool           `json:"mandatory"`
	UserDefined   bool           `json:"user-defined"`
	Obsoletes     bool           `json:"obsoletes"`
	Values        VariableValues `json:"values"`
	IsSubcategory bool           `json:"is-subcategory"`
	Links         []Link         `json:"links"`
}

type VariableScope struct {
	Type    VariableScopeType `json:"type"`
	LevelID string            `json:"level"`
}

type VariableScopeType string

const (
	ScopeGlobal      VariableScopeType = "global"
	ScopeFullGame    VariableScopeType = "full-game"
	ScopeAllLevels   VariableScopeType = "all-levels"
	ScopeSingleLevel VariableScopeType = "single-level"
)

type VariableValues struct {
	Values  map[string]VariableValue `json:"values"`
	Default string                   `json:"default"`
}

type VariableValue struct {
	Label string              `json:"label"`
	Rules string              `json:"rules"`
	Flags *VariableValueFlags `json:"flags"`
}

type VariableValueFlags struct {
	Miscellaneous bool `json:"miscellaneous"`
}
