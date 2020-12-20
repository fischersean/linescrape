package fte

import (
	"fmt"
	"strings"
)

type EloProjection struct {
	Date         string  `csv:"date"`
	Season       int     `csv:"season"`
	Neutral      int     `csv:"neutral"`
	Playoff      string  `csv:"playoff"`
	Team1        string  `csv:"team1"`
	Team2        string  `csv:"team2"`
	Elo1Pre      float64 `csv:"elo1_pre"`
	Elo2Pre      float64 `csv:"elo2_pre"`
	EloProb1     float64 `csv:"elo_prob1"`
	EloProb2     float64 `csv:"elo_prob2"`
	EloPost1     float64 `csv:"elo1_post"`
	EloPost2     float64 `csv:"elo2_post"`
	QbEloPre1    float64 `csv:"qbelo1_pre"`
	QbEloPre2    float64 `csv:"qbelo2_pre"`
	QB1          string  `csv:"qb1"`
	QB2          string  `csv:"qb2"`
	QB1ValuePre  float64 `csv:"qb1_value_pre"`
	QB2ValuePre  float64 `csv:"qb2_value_pre"`
	QB1Adj       float64 `csv:"qb1_adj"`
	QB2Adj       float64 `csv:"qb2_adj"`
	QBEloProb1   float64 `csv:"qbelo_prob1"`
	QBEloProb2   float64 `csv:"qbelo_prob2"`
	QB1GameValue float64 `csv:"qb1_game_value"`
	QB2GameValue float64 `csv:"qb2_game_value"`
	QB1ValuePost float64 `csv:"qb1_value_post"`
	QB2ValuePost float64 `csv:"qb2_value_post"`
	QBEloPost1   float64 `csv:"qbelo1_post"`
	QBEloPost2   float64 `csv:"qbelo2_post"`
	Score1       float64 `csv:"score1"`
	Score2       float64 `csv:"score2"`
}

func GenGameID(p EloProjection) string {
	parts := strings.Split(p.Date, "-")
	return fmt.Sprintf("%s%s%s%s%s%s",
		"NFL",
		parts[0], parts[1], parts[2],
		p.Team1, p.Team2,
	)
}

var ToCommon = map[string]string{
	"ARI": "Arizona Cardinals",
	"ATL": "Atlanta Falcons",
	"BAL": "Baltimore Ravens",
	"BUF": "Buffalo Bills",
	"CAR": "Carolina Panthers",
	"CHI": "Chicago Bears",
	"CIN": "Cincinnati Bengals",
	"CLE": "Cleveland Browns",
	"DAL": "Dallas Cowboys",
	"DEN": "Denver Broncos",
	"DET": "Detroit Lions",
	"GB":  "Green Bay Packers",
	"HOU": "Houston Texans",
	"IND": "Indianapolis Colts",
	"JAX": "Jacksonville Jaguars",
	"KC":  "Kansas City Chiefs",
	"LAC": "Los Angeles Chargers",
	"LAR": "Los Angeles Rams",
	"MIA": "Miami Dolphins",
	"MIN": "Minnesota Vikings",
	"NE":  "New England Patriots",
	"NO":  "New Orleans Saints",
	"NYG": "New York Giants",
	"NYJ": "New York Jets",
	"OAK": "Las Vegas Raiders",
	"PHI": "Philadelphia Eagles",
	"PIT": "Pittsburgh Steelers",
	"SEA": "Seattle Seahawks",
	"SF":  "San Francisco 49ers",
	"TB":  "Tampa Bay Buccaneers",
	"TEN": "Tennessee Titans",
	"WSH": "Washington Football Team",
}

var FromCommon = map[string]string{
	"Arizona Cardinals":        "ARI",
	"Atlanta Falcons":          "ATL",
	"Baltimore Ravens":         "BAL",
	"Buffalo Bills":            "BUF",
	"Carolina Panthers":        "CAR",
	"Chicago Bears":            "CHI",
	"Cincinnati Bengals":       "CIN",
	"Cleveland Browns":         "CLE",
	"Dallas Cowboys":           "DAL",
	"Denver Broncos":           "DEN",
	"Detroit Lions":            "DET",
	"Green Bay Packers":        "GB",
	"Houston Texans":           "HOU",
	"Indianapolis Colts":       "IND",
	"Jacksonville Jaguars":     "JAX",
	"Kansas City Chiefs":       "KC",
	"Los Angeles Chargers":     "LAC",
	"Los Angeles Rams":         "LAR",
	"Miami Dolphins":           "MIA",
	"Minnesota Vikings":        "MIN",
	"New England Patriots":     "NE",
	"New Orleans Saints":       "NO",
	"New York Giants":          "NYG",
	"New York Jets":            "NYJ",
	"Las Vegas Raiders":        "OAK",
	"Philadelphia Eagles":      "PHI",
	"Pittsburgh Steelers":      "PIT",
	"Seattle Seahawks":         "SEA",
	"San Francisco 49ers":      "SF",
	"Tampa Bay Buccaneers":     "TB",
	"Tennessee Titans":         "TEN",
	"Washington Football Team": "WSH",
}
