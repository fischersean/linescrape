package caesars

import (
	"time"
)

type rawResponse struct {
	Data data `json:"data"`
}
type competitionSummary struct {
	CompetitionDrilldownTagID string   `json:"competitionDrilldownTagId"`
	TypeIds                   []string `json:"typeIds"`
	EventCount                int      `json:"eventCount"`
}
type externalIds struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
}
type alternativeNames struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type teams struct {
	Side             string             `json:"side"`
	Name             string             `json:"name"`
	AlternativeNames []alternativeNames `json:"alternativeNames"`
}
type Type struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	DisplayOrder       string `json:"displayOrder"`
	FixedOddsAvailable bool   `json:"fixedOddsAvailable"`
}
type class struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayOrder string `json:"displayOrder"`
}
type category struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayOrder string `json:"displayOrder"`
	Code         string `json:"code"`
}
type prices struct {
	Numerator    int     `json:"numerator"`
	Denominator  int     `json:"denominator"`
	Decimal      float64 `json:"decimal"`
	DisplayOrder int     `json:"displayOrder"`
	PriceType    string  `json:"priceType"`
	HandicapLow  string  `json:"handicapLow"`
	HandicapHigh string  `json:"handicapHigh"`
}
type outcomes struct {
	ID           string   `json:"id"`
	MarketID     string   `json:"marketId"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	SubType      string   `json:"subType"`
	DisplayOrder int      `json:"displayOrder"`
	Active       bool     `json:"active"`
	Displayed    bool     `json:"displayed"`
	Status       string   `json:"status"`
	Channels     []string `json:"channels"`
	Resulted     bool     `json:"resulted"`
	RcpAvailable bool     `json:"rcpAvailable"`
	RunnerNumber int      `json:"runnerNumber"`
	RetailCode   string   `json:"retailCode"`
	Prices       []prices `json:"prices"`
}
type markets struct {
	ID                 string     `json:"id"`
	EventID            string     `json:"eventId"`
	TemplateMarketID   string     `json:"templateMarketId"`
	Type               string     `json:"type"`
	SubType            string     `json:"subType"`
	Name               string     `json:"name"`
	Active             bool       `json:"active"`
	Displayed          bool       `json:"displayed"`
	Status             string     `json:"status"`
	BetInRun           bool       `json:"betInRun"`
	DisplayOrder       int        `json:"displayOrder"`
	Channels           []string   `json:"channels"`
	HandicapValue      float64    `json:"handicapValue"`
	LivePriceAvailable bool       `json:"livePriceAvailable"`
	CashoutAvailable   bool       `json:"cashoutAvailable"`
	MinimumAccumulator int        `json:"minimumAccumulator"`
	MaximumAccumulator int        `json:"maximumAccumulator"`
	FixedOddsAvailable bool       `json:"fixedOddsAvailable"`
	RcpAvailable       bool       `json:"rcpAvailable"`
	Outcomes           []outcomes `json:"outcomes"`
}
type events struct {
	ID                        string        `json:"id"`
	ExternalIds               []externalIds `json:"externalIds"`
	Name                      string        `json:"name"`
	Active                    bool          `json:"active"`
	Displayed                 bool          `json:"displayed"`
	Status                    string        `json:"status"`
	DisplayOrder              int           `json:"displayOrder"`
	Channels                  []string      `json:"channels"`
	SortCode                  string        `json:"sortCode"`
	StartTime                 time.Time     `json:"startTime"`
	Started                   bool          `json:"started"`
	LiveNow                   bool          `json:"liveNow"`
	LiveBettingAvailable      bool          `json:"liveBettingAvailable"`
	Resulted                  bool          `json:"resulted"`
	Settled                   bool          `json:"settled"`
	CashoutAvailable          bool          `json:"cashoutAvailable"`
	SportID                   string        `json:"sportId"`
	CompetitionDrilldownTagID string        `json:"competitionDrilldownTagId"`
	RaceNumber                int           `json:"raceNumber"`
	IsVoid                    bool          `json:"isVoid"`
	FixedOddsAvailable        bool          `json:"fixedOddsAvailable"`
	RcpAvailable              bool          `json:"rcpAvailable"`
	StatisticsAvailable       bool          `json:"statisticsAvailable"`
	Teams                     []teams       `json:"teams"`
	Type                      Type          `json:"type"`
	Class                     class         `json:"class"`
	Category                  category      `json:"category"`
	MarketCount               int           `json:"marketCount"`
	Markets                   []markets     `json:"markets"`
}
type timeBandEvents struct {
	Type               string               `json:"type"`
	Date               string               `json:"date"`
	CompetitionSummary []competitionSummary `json:"competitionSummary"`
	Events             []events             `json:"events"`
}
type data struct {
	TimeBandEvents []timeBandEvents `json:"timeBandEvents"`
}
