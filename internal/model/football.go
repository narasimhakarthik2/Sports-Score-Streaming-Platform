package model

import "time"

type Match struct {
	ID        string    `json:"id"`
	Sport     string    `json:"sport"`
	League    string    `json:"league"`
	Season    int       `json:"season"`
	HomeTeam  Team      `json:"home_team"`
	AwayTeam  Team      `json:"away_team"`
	StartTime time.Time `json:"start_time"`
	Status    string    `json:"status"`
	Matchday  int       `json:"matchday"`
	Stage     string    `json:"stage"`
	Group     string    `json:"group"`
	Score     Score     `json:"score"`
}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Score struct {
	Winner   string      `json:"winner"`
	Duration string      `json:"duration"`
	FullTime ScoreDetail `json:"full_time"`
	HalfTime ScoreDetail `json:"half_time"`
}

type ScoreDetail struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

type Competition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Type string `json:"type"`
	Area string `json:"area"`
}
