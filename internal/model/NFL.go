package model

import "time"

// NFLTeam represents an NFL team with its basic details
type NFLTeam struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Logos       string `json:"logos"`
}

// NFLMatch represents an NFL match or game
type NFLMatch struct {
	ID           string           `json:"id"`
	UID          string           `json:"uid"`
	Name         string           `json:"name"`
	ShortName    string           `json:"shortName"`
	Date         time.Time        `json:"date"`
	Competitions []NFLCompetition `json:"competitions"`
	Venue        NFLVenue         `json:"venue"`
	Status       NFLStatus        `json:"status"`
}

// NFLCompetition holds details about the competition including teams, venue, and situation
type NFLCompetition struct {
	ID                string          `json:"id"`
	Attendance        int             `json:"attendance"`
	Type              NFLType         `json:"type"`
	NeutralSite       bool            `json:"neutralSite"`
	Competitors       []NFLCompetitor `json:"competitors"`
	Situation         NFLSituation    `json:"situation"`
	HasDefensiveStats bool            `json:"hasDefensiveStats"`
}

// NFLType represents the type of competition (e.g., standard)
type NFLType struct {
	ID           string `json:"id"`
	Text         string `json:"text"`
	Abbreviation string `json:"abbreviation"`
	Slug         string `json:"slug"`
}

// NFLCompetitor represents the teams playing in the competition (home/away, winner)
type NFLCompetitor struct {
	ID       string  `json:"id"`
	HomeAway string  `json:"homeAway"`
	Winner   bool    `json:"winner"`
	Team     NFLTeam `json:"team"`
}

// NFLVenue represents the venue where the match takes place
type NFLVenue struct {
	ID       string          `json:"id"`
	FullName string          `json:"fullName"`
	Address  NFLVenueAddress `json:"address"`
	Grass    bool            `json:"grass"`
	Indoor   bool            `json:"indoor"`
	Images   []NFLVenueImage `json:"images"`
}

// NFLVenueAddress holds the address details of the venue
type NFLVenueAddress struct {
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
}

// NFLVenueImage represents images of the venue
type NFLVenueImage struct {
	Href   string   `json:"href"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Alt    string   `json:"alt"`
	Rel    []string `json:"rel"`
}

// NFLSituation holds the current situation in the match (e.g., down, yard line, etc.)
type NFLSituation struct {
	Down         int  `json:"down"`
	YardLine     int  `json:"yardLine"`
	Distance     int  `json:"distance"`
	IsRedZone    bool `json:"isRedZone"`
	HomeTimeouts int  `json:"homeTimeouts"`
	AwayTimeouts int  `json:"awayTimeouts"`
}

// NFLStatus represents the current status of the match
type NFLStatus struct {
	Clock        int           `json:"clock"`
	DisplayClock string        `json:"displayClock"`
	Period       int           `json:"period"`
	Type         NFLStatusType `json:"type"`
}

// NFLStatusType holds the type of status (e.g., final, post)
type NFLStatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
	ShortDetail string `json:"shortDetail"`
}
