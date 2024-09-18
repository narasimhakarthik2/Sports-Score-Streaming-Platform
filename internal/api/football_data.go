package api

import (
	"Sports-Score-Streaming-Platform/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FootballDataClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewFootballDataClient(apiKey string) *FootballDataClient {
	return &FootballDataClient{
		baseURL: "http://api.football-data.org/v4",
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *FootballDataClient) FetchMatches() ([]model.Match, error) {
	url := fmt.Sprintf("%s/matches", c.baseURL)
	return c.fetchMatches(url)
}

func (c *FootballDataClient) FetchMatchesByCompetition(competitionID string, dateFrom, dateTo time.Time) ([]model.Match, error) {
	url := fmt.Sprintf("%s/competitions/%s/matches?dateFrom=%s&dateTo=%s",
		c.baseURL, competitionID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02"))
	return c.fetchMatches(url)
}

func (c *FootballDataClient) FetchMatchesByTeam(teamID string, dateFrom, dateTo time.Time) ([]model.Match, error) {
	url := fmt.Sprintf("%s/teams/%s/matches?dateFrom=%s&dateTo=%s",
		c.baseURL, teamID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02"))
	return c.fetchMatches(url)
}

func (c *FootballDataClient) fetchMatches(url string) ([]model.Match, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		Matches []struct {
			ID          int `json:"id"`
			Competition struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"competition"`
			Season struct {
				ID int `json:"id"`
			} `json:"season"`
			UtcDate  string `json:"utcDate"`
			Status   string `json:"status"`
			Matchday int    `json:"matchday"`
			Stage    string `json:"stage"`
			Group    string `json:"group"`
			HomeTeam struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"homeTeam"`
			AwayTeam struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"awayTeam"`
			Score struct {
				Winner   string `json:"winner"`
				Duration string `json:"duration"`
				FullTime struct {
					Home int `json:"home"`
					Away int `json:"away"`
				} `json:"fullTime"`
				HalfTime struct {
					Home int `json:"home"`
					Away int `json:"away"`
				} `json:"halfTime"`
			} `json:"score"`
		} `json:"matches"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var matches []model.Match
	for _, m := range result.Matches {
		startTime, _ := time.Parse(time.RFC3339, m.UtcDate)
		matches = append(matches, model.Match{
			ID:     fmt.Sprintf("football-data-%d", m.ID),
			Sport:  "Soccer",
			League: m.Competition.Name,
			Season: m.Season.ID,
			HomeTeam: model.Team{
				ID:   fmt.Sprintf("%d", m.HomeTeam.ID),
				Name: m.HomeTeam.Name,
			},
			AwayTeam: model.Team{
				ID:   fmt.Sprintf("%d", m.AwayTeam.ID),
				Name: m.AwayTeam.Name,
			},
			StartTime: startTime,
			Status:    m.Status,
			Matchday:  m.Matchday,
			Stage:     m.Stage,
			Group:     m.Group,
			Score: model.Score{
				Winner:   m.Score.Winner,
				Duration: m.Score.Duration,
				FullTime: model.ScoreDetail{
					Home: m.Score.FullTime.Home,
					Away: m.Score.FullTime.Away,
				},
				HalfTime: model.ScoreDetail{
					Home: m.Score.HalfTime.Home,
					Away: m.Score.HalfTime.Away,
				},
			},
		})
	}

	return matches, nil
}

func (c *FootballDataClient) FetchCompetitions() ([]model.Competition, error) {
	url := fmt.Sprintf("%s/competitions", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		Competitions []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
			Type string `json:"type"`
			Area struct {
				Name string `json:"name"`
			} `json:"area"`
		} `json:"competitions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var competitions []model.Competition
	for _, c := range result.Competitions {
		competitions = append(competitions, model.Competition{
			ID:   fmt.Sprintf("%d", c.ID),
			Name: c.Name,
			Code: c.Code,
			Type: c.Type,
			Area: c.Area.Name,
		})
	}

	return competitions, nil
}
