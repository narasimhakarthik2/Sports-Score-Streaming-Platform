package api

import (
	"Sports-Score-Streaming-Platform/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type NFLDataClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewNFLDataClient(apiKey string) *NFLDataClient {
	return &NFLDataClient{
		baseURL: "https://nfl-api-data.p.rapidapi.com",
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Fetch the current NFL week
func (c *NFLDataClient) FetchWeekForDate(date time.Time) (int, error) {
	url := fmt.Sprintf("%s/nfl-whitelist", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("x-rapidapi-key", c.apiKey)
	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	sections, ok := data["sections"].([]interface{})
	if !ok || len(sections) < 2 {
		return 0, fmt.Errorf("unexpected sections structure")
	}

	weeks := sections[1].(map[string]interface{})["entries"].([]interface{})

	for _, week := range weeks {
		weekEntry := week.(map[string]interface{})

		startDate, err := time.Parse(time.RFC3339Nano, weekEntry["startDate"].(string))
		if err != nil {
			log.Printf("Error parsing startDate: %v", err)
			continue
		}

		endDate, err := time.Parse(time.RFC3339Nano, weekEntry["endDate"].(string))
		if err != nil {
			log.Printf("Error parsing endDate: %v", err)
			continue
		}

		if date.After(startDate) && date.Before(endDate) {
			return int(weekEntry["value"].(float64)), nil
		}
	}

	return 0, fmt.Errorf("week not found for the given date")
}

func (c *NFLDataClient) FetchGamesForCurrentWeek() ([]model.NFLMatch, error) {
	currentWeek := 2 // This should be determined dynamically in a real scenario
	url := fmt.Sprintf("%s/nfl-weeks-events?year=%d&week=%d&type=2", c.baseURL, time.Now().Year(), currentWeek)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-rapidapi-key", c.apiKey)
	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []struct {
			EventID string `json:"eventid"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var matches []model.NFLMatch
	for _, item := range result.Items {
		match, err := c.FetchGameDetails(item.EventID)
		if err != nil {
			log.Printf("Error fetching game details for event %s: %v", item.EventID, err)
			continue
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (c *NFLDataClient) FetchGameDetails(eventID string) (model.NFLMatch, error) {
	url := fmt.Sprintf("%s/nfl-single-events?id=%s", c.baseURL, eventID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return model.NFLMatch{}, err
	}

	req.Header.Set("x-rapidapi-key", c.apiKey)
	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

	res, err := c.client.Do(req)
	if err != nil {
		return model.NFLMatch{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return model.NFLMatch{}, fmt.Errorf("API request failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.NFLMatch{}, err
	}

	var rawMatch struct {
		ID           string `json:"id"`
		UID          string `json:"uid"`
		Name         string `json:"name"`
		ShortName    string `json:"shortName"`
		Date         string `json:"date"`
		Competitions []struct {
			ID         string `json:"id"`
			Attendance int    `json:"attendance"`
			Type       struct {
				ID           string `json:"id"`
				Text         string `json:"text"`
				Abbreviation string `json:"abbreviation"`
				Slug         string `json:"slug"`
			} `json:"type"`
			NeutralSite bool `json:"neutralSite"`
			Competitors []struct {
				ID       string `json:"id"`
				HomeAway string `json:"homeAway"`
				Winner   bool   `json:"winner"`
				Team     struct {
					ID          string `json:"id"`
					DisplayName string `json:"displayName"`
					Logos       string `json:"logos"`
				} `json:"team"`
			} `json:"competitors"`
			Situation struct {
				Down         int  `json:"down"`
				YardLine     int  `json:"yardLine"`
				Distance     int  `json:"distance"`
				IsRedZone    bool `json:"isRedZone"`
				HomeTimeouts int  `json:"homeTimeouts"`
				AwayTimeouts int  `json:"awayTimeouts"`
			} `json:"situation"`
			HasDefensiveStats bool `json:"hasDefensiveStats"`
		} `json:"competitions"`
		Venue struct {
			ID       string `json:"id"`
			FullName string `json:"fullName"`
			Address  struct {
				City    string `json:"city"`
				State   string `json:"state"`
				ZipCode string `json:"zipCode"`
			} `json:"address"`
			Grass  bool `json:"grass"`
			Indoor bool `json:"indoor"`
			Images []struct {
				Href   string   `json:"href"`
				Width  int      `json:"width"`
				Height int      `json:"height"`
				Alt    string   `json:"alt"`
				Rel    []string `json:"rel"`
			} `json:"images"`
		} `json:"venue"`
		Status struct {
			Clock        int    `json:"clock"`
			DisplayClock string `json:"displayClock"`
			Period       int    `json:"period"`
			Type         struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				State       string `json:"state"`
				Completed   bool   `json:"completed"`
				Description string `json:"description"`
				Detail      string `json:"detail"`
				ShortDetail string `json:"shortDetail"`
			} `json:"type"`
		} `json:"status"`
	}

	if err := json.Unmarshal(body, &rawMatch); err != nil {
		return model.NFLMatch{}, err
	}

	date, err := time.Parse("2006-01-02T15:04Z", rawMatch.Date)
	if err != nil {
		return model.NFLMatch{}, fmt.Errorf("failed to parse date: %v", err)
	}

	match := model.NFLMatch{
		ID:        rawMatch.ID,
		UID:       rawMatch.UID,
		Name:      rawMatch.Name,
		ShortName: rawMatch.ShortName,
		Date:      date,
		Venue: model.NFLVenue{
			ID:       rawMatch.Venue.ID,
			FullName: rawMatch.Venue.FullName,
			Address: model.NFLVenueAddress{
				City:    rawMatch.Venue.Address.City,
				State:   rawMatch.Venue.Address.State,
				ZipCode: rawMatch.Venue.Address.ZipCode,
			},
			Grass:  rawMatch.Venue.Grass,
			Indoor: rawMatch.Venue.Indoor,
		},
		Status: model.NFLStatus{
			Clock:        rawMatch.Status.Clock,
			DisplayClock: rawMatch.Status.DisplayClock,
			Period:       rawMatch.Status.Period,
			Type: model.NFLStatusType{
				ID:          rawMatch.Status.Type.ID,
				Name:        rawMatch.Status.Type.Name,
				State:       rawMatch.Status.Type.State,
				Completed:   rawMatch.Status.Type.Completed,
				Description: rawMatch.Status.Type.Description,
				Detail:      rawMatch.Status.Type.Detail,
				ShortDetail: rawMatch.Status.Type.ShortDetail,
			},
		},
	}

	for _, comp := range rawMatch.Competitions {
		competition := model.NFLCompetition{
			ID:         comp.ID,
			Attendance: comp.Attendance,
			Type: model.NFLType{
				ID:           comp.Type.ID,
				Text:         comp.Type.Text,
				Abbreviation: comp.Type.Abbreviation,
				Slug:         comp.Type.Slug,
			},
			NeutralSite: comp.NeutralSite,
			Situation: model.NFLSituation{
				Down:         comp.Situation.Down,
				YardLine:     comp.Situation.YardLine,
				Distance:     comp.Situation.Distance,
				IsRedZone:    comp.Situation.IsRedZone,
				HomeTimeouts: comp.Situation.HomeTimeouts,
				AwayTimeouts: comp.Situation.AwayTimeouts,
			},
			HasDefensiveStats: comp.HasDefensiveStats,
		}

		for _, competitor := range comp.Competitors {
			competition.Competitors = append(competition.Competitors, model.NFLCompetitor{
				ID:       competitor.ID,
				HomeAway: competitor.HomeAway,
				Winner:   competitor.Winner,
				Team: model.NFLTeam{
					ID:          competitor.Team.ID,
					DisplayName: competitor.Team.DisplayName,
					Logos:       competitor.Team.Logos,
				},
			})
		}

		match.Competitions = append(match.Competitions, competition)
	}

	for _, img := range rawMatch.Venue.Images {
		match.Venue.Images = append(match.Venue.Images, model.NFLVenueImage{
			Href:   img.Href,
			Width:  img.Width,
			Height: img.Height,
			Alt:    img.Alt,
			Rel:    img.Rel,
		})
	}

	return match, nil
}

func (c *NFLDataClient) FetchNFLTeams() ([]model.NFLTeam, error) {
	url := fmt.Sprintf("%s/nfl-teams", c.baseURL)
	return c.fetchTeams(url)
}

func (c *NFLDataClient) FetchNFLTeamDetail(teamID string) (model.NFLTeam, error) {
	url := fmt.Sprintf("%s/nfl-team?team_id=%s", c.baseURL, teamID)
	return c.fetchTeamDetail(url)
}

func (c *NFLDataClient) FetchNFLLiveScore() ([]model.NFLMatch, error) {
	url := fmt.Sprintf("%s/nfl-livescores", c.baseURL)
	return c.fetchMatches(url)
}

func (c *NFLDataClient) FetchNFLMatches() ([]model.NFLMatch, error) {
	url := fmt.Sprintf("%s/nfl-matches", c.baseURL) // Adjust the endpoint based on the actual API
	return c.fetchMatches(url)
}

func (c *NFLDataClient) fetchTeams(url string) ([]model.NFLTeam, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-rapidapi-key", c.apiKey)
	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		Teams []model.NFLTeam `json:"teams"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Teams, nil
}

func (c *NFLDataClient) fetchTeamDetail(url string) (model.NFLTeam, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return model.NFLTeam{}, err
	}

	req.Header.Set("x-rapidapi-key", c.apiKey)
	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

	resp, err := c.client.Do(req)
	if err != nil {
		return model.NFLTeam{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.NFLTeam{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var team model.NFLTeam
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return model.NFLTeam{}, err
	}

	return team, nil
}

// func (c *NFLDataClient) fetchPlayerDetail(url string) (model.NFLPlayer, error) {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return model.NFLPlayer{}, err
// 	}

// 	req.Header.Set("x-rapidapi-key", c.apiKey)
// 	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return model.NFLPlayer{}, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return model.NFLPlayer{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
// 	}

// 	var player model.NFLPlayer
// 	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
// 		return model.NFLPlayer{}, err
// 	}

// 	return player, nil
// }

// func (c *NFLDataClient) fetchPlayers(url string) ([]model.NFLPlayer, error) {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("x-rapidapi-key", c.apiKey)
// 	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
// 	}

// 	var result struct {
// 		Players []model.NFLPlayer `json:"players"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}

// 	return result.Players, nil
// }

func (c *NFLDataClient) fetchMatches(url string) ([]model.NFLMatch, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-rapidapi-key", c.apiKey)
	req.Header.Set("x-rapidapi-host", "nfl-api-data.p.rapidapi.com")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		Matches []model.NFLMatch `json:"matches"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Matches, nil
}
