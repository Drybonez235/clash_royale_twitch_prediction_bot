package clash_royale_api

import()

type Match_25 struct {
	Matches []Match
}

type Match struct {
	Type               string  `json:"type"`
	BattleTime         string  `json:"battleTime"`
	IsLadderTournament bool    `json:"isLadderTournament"`
	Arena             Arena   `json:"arena"`
	GameMode          GameMode `json:"gameMode"`
	DeckSelection     string  `json:"deckSelection"`
	Team             []Player `json:"team"`
	Opponent         []Player `json:"opponent"`
	IsHostedMatch     bool    `json:"isHostedMatch"`
	LeagueNumber      int     `json:"leagueNumber"`
}

type Arena struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GameMode struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Player struct {
	Tag                   string   `json:"tag"`
	Name                  string   `json:"name"`
	StartingTrophies      int      `json:"startingTrophies"`
	TrophyChange         int      `json:"trophyChange"`
	Crowns               int      `json:"crowns"`
	KingTowerHitPoints   int      `json:"kingTowerHitPoints"`
	PrincessTowersHitPoints []int `json:"princessTowersHitPoints"`
	Clan                 Clan     `json:"clan"`
	Cards               []Card    `json:"cards"`
	SupportCards        []Card    `json:"supportCards"`
	GlobalRank          *int      `json:"globalRank"`
	ElixirLeaked        float64   `json:"elixirLeaked"`
}

type Clan struct {
	Tag     string `json:"tag"`
	Name    string `json:"name"`
	BadgeID int    `json:"badgeId"`
}

type Card struct {
	Name             string `json:"name"`
	ID               int    `json:"id"`
	Level            int    `json:"level"`
	StarLevel        *int   `json:"starLevel,omitempty"`
	EvolutionLevel   *int   `json:"evolutionLevel,omitempty"`
	MaxLevel         int    `json:"maxLevel"`
	MaxEvolutionLevel *int  `json:"maxEvolutionLevel,omitempty"`
	Rarity           string `json:"rarity"`
	ElixirCost       int    `json:"elixirCost"`
	IconUrls         IconUrls `json:"iconUrls"`
}

type IconUrls struct {
	Medium         string `json:"medium"`
	EvolutionMedium *string `json:"evolutionMedium,omitempty"`
}

