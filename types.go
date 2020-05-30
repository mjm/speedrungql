package speedrungql

import (
	"encoding/json"
)

type EnvelopeResponse struct {
	Data json.RawMessage `json:"data"`
}

type GamesResponse struct {
	Data []Game `json:"data"`
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

type PlatformsResponse struct {
	Data []Platform `json:"data"`
}

type PlatformResponse struct {
	Data Platform `json:"data"`
}

type Platform struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Released int32  `json:"released"`
}
