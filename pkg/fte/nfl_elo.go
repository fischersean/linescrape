package fte

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
