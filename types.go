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
	ShowMilliseconds    bool     `json:"show-milliseconds"`
	RequireVerification bool     `json:"require-verification"`
	RequireVideo        bool     `json:"require-video"`
	RunTimes            []string `json:"run-times"`
	DefaultRunTime      string   `json:"default-time"`
	EmulatorsAllowed    bool     `json:"emulators-allowed"`
}

type CategoriesResponse struct {
	Data []*Category `json:"data"`
}

type Category struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Weblink       string          `json:"weblink"`
	Type          string          `json:"type"`
	Rules         string          `json:"rules"`
	Players       CategoryPlayers `json:"players"`
	Miscellaneous bool            `json:"miscellaneous"`
}

type CategoryPlayers struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
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
	Timing     string      `json:"timing"`
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
	Status     string     `json:"status"`
	ExaminerID string     `json:"examiner"`
	VerifyDate *time.Time `json:"verify-date"`
	Reason     string     `json:"reason"`
}
