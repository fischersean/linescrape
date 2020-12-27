package fte

var (
	NbaLatestUrl  = "https://projects.fivethirtyeight.com/nba-model/nba_elo_latest.csv"
	NbaHistoryUrl = "https://projects.fivethirtyeight.com/nba-model/nba_elo.csv"
)

type NbaEloProjection struct {
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
	CarmElo1Pre  float64 `csv:"carm-elo1_pre"`
	CarmElo2Pre  float64 `csv:"carm-elo2_pre"`
	CarmEloProb1 float64 `csv:"carm-elo_prob1"`
	CarmEloProb2 float64 `csv:"carm-elo_prob2"`
	CarmElo1Post float64 `csv:"carm-elo1_post"`
	CarmElo2Post float64 `csv:"carm-elo2_post"`
	Raptor1Pre   float64 `csv:"raptor1_pre"`
	Raptor2Pre   float64 `csv:"raptor2_pre"`
	RaptorProb1  float64 `csv:"raptor_prob1"`
	RaptorProb2  float64 `csv:"raptor_prob2"`
	Score1       float64 `csv:"score1"`
	Score2       float64 `csv:"score2"`
}

var ToCommonNbaName = map[string]string{
	"PHI": "Philadelphia 76ers",
	"CLE": "Cleveland Cavaliers",
	"IND": "Indiana Pacers",
	"ORL": "Orlando Magic",
	"BOS": "Boston Celtics",
	"TOR": "Toronto Raptors",
	"MIN": "Minnesota Timberwolves",
	"MEM": "Memphis Grizzlies",
	"CHI": "Chicago Bulls",
	"DEN": "Denver Nuggets",
	"POR": "Portland Trail Blazers",
	"PHO": "Phoenix Suns",
	"MIA": "Miami Heat",
	"MIL": "Milwaukee Bucks",
	"CHO": "Charlotte Hornets",
	"DET": "Detroit Pistons",
	"WAS": "Washington Wizards",
	"NYK": "New York Knicks",
	"SAS": "San Antonio Spurs",
	"UTA": "Utah Jazz",
	"SAC": "Sacramento Kings",
	"LAC": "LA Clippers",
	"NOP": "New Orleans Pelicans",
	"ATL": "Atlanta Hawks",
	"OKC": "Oklahoma City Thunder",
	"DAL": "Dallas Mavericks",
	"HOU": "Houston Rockets",
	"GSW": "Golden State Warriors",
}

var FromCommonNbaName = map[string]string{
	"Philadelphia 76ers":     "PHI",
	"Cleveland Cavaliers":    "CLE",
	"Indiana Pacers":         "IND",
	"Orlando Magic":          "ORL",
	"Boston Celtics":         "BOS",
	"Toronto Raptors":        "TOR",
	"Minnesota Timberwolves": "MIN",
	"Memphis Grizzlies":      "MEM",
	"Chicago Bulls":          "CHI",
	"Denver Nuggets":         "DEN",
	"Portland Trail Blazers": "POR",
	"Phoenix Suns":           "PHO",
	"Miami Heat":             "MIA",
	"Milwaukee Bucks":        "MIL",
	"Charlotte Hornets":      "CHO",
	"Detroit Pistons":        "DET",
	"Washington Wizards":     "WAS",
	"New York Knicks":        "NYK",
	"San Antonio Spurs":      "SAS",
	"Utah Jazz":              "UTA",
	"Sacramento Kings":       "SAC",
	"LA Clippers":            "LAC",
	"New Orleans Pelicans":   "NOP",
	"Atlanta Hawks":          "ATL",
	"Oklahoma City Thunder":  "OKC",
	"Dallas Mavericks":       "DAL",
	"Houston Rockets":        "HOU",
	"Golden State Warriors":  "GSW",
}
