package game

import "time"

type Projection struct {
	GameId                 string    `json:"gameId"`
	League                 string    `json:"league"`
	Source                 string    `json:"source"`
	HomeWinProbability     float64   `json:"homeWinProb"`
	VisitingWinProbability float64   `json:"visitingWinProb"`
	Home                   string    `json:"home"`
	Visiting               string    `json:"visiting"`
	GameDate               time.Time `json:"gameDate"`
	Season                 int       `json:"season"`
	Playoff                string    `json:"playoff"`
}

type Line struct {
	GameTime             string
	HomeTeam             string
	VisitingTeam         string
	HomeSpreadPoints     float64
	VisitingSpreadPoints float64
	HomeSpreadLine       int64
	VisitingSpreadLine   int64
	HomeMoneyLine        int64
	VisitingMoneyLine    int64
	GameOver             float64
	GameOverLine         int64
	GameUnderLine        int64
}
