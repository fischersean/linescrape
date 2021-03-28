package caesars

import (
	"github.com/fischersean/linescrape/pkg/game"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	leagueTagIds = map[string]string{
		"NFL":                "54",
		"college-football":   "55",
		"ALL-FOOTBALL":       "3",
		"NBA":                "77",
		"college-basketball": "81",
		"ALL-BASKETBALL":     "19",
		"MLB":                "62",
	}

	SupportedLeagues = []string{
		"NFL",
		"college-football",
		"NBA",
		"college-basketball",
		"MLB",
	}

	SupportedSports = map[string]string{
		"Football":   "",
		"Basketball": "",
		"Baseball":   "",
	}
)

func constructUrl(league string, dates []time.Time) string {

	u := &url.URL{
		Scheme: "https",
		Host:   "sb-content.caesarscasino.com",
		Path:   "content-service/api/v1/q/time-band-event-list",
	}

	v := url.Values{}

	v.Add("active", "true")
	v.Add("maxMarkets", "10")
	v.Add("marketSortsIncluded", "HH,HL,MR,WH")
	//v.Add("marketSortsIncluded", "HL")
	//v.Add("marketSortsIncluded", "MR")
	//v.Add("marketSortsIncluded", "WH")
	v.Add("allowedEventSorts", "MTCH")
	v.Add("includeChildMarkets", "true")
	v.Add("prioritisePrimaryMarkets", "true")
	v.Add("includeMedia", "false")
	if tid, ok := leagueTagIds[league]; ok {
		// Although the api appears to support multiple tags per request,
		// I am disallowing it to keep things simpler on my end
		v.Add("drilldownTagIds", tid)
	}
	v.Add("maxTotalItems", "100")
	v.Add("maxEventsPerCompetition", "50")
	v.Add("maxCompetitionsPerSportPerBand", "2")
	v.Add("maxEventsForNextToGo", "0")
	v.Add("startTimeOffsetForNextToGo", "600")

	// Dates to pull
	dList := ""
	for _, d := range dates {
		dList = fmt.Sprintf("%s,%s", dList, d.Format(time.RFC3339))
	}
	v.Add("dates", dList[1:])

	u.RawQuery = v.Encode()

	return u.String()

}

func lineFromEvent(ev events) (odd game.Line, err error) {

	// Set home and away team
	odd.HomeTeam, odd.VisitingTeam = func(ts []teams) (string, string) {

		var ht string
		var vt string

		for _, v := range ts {
			if v.Side == "HOME" {
				ht = v.Name
			}
			if v.Side == "AWAY" {
				vt = v.Name
			}
		}
		return ht, vt

	}(ev.Teams)

	// Set start time
	odd.GameTime = ev.StartTime.String()

	assignOutcomes := func(m markets) (outcomes, outcomes) {
		return m.Outcomes[0], m.Outcomes[1]
	}

	for _, market := range ev.Markets {
		//log.Printf(market.Name)
		if market.Name == "Money Line" {
			func(m markets) {

				hotc, votc := assignOutcomes(m)

				odd.HomeMoneyLine = game.DecimalToAmerican(hotc.Prices[0].Decimal)
				odd.VisitingMoneyLine = game.DecimalToAmerican(votc.Prices[0].Decimal)

			}(market)
		}
		if market.Name == "Spread" || market.Name == "Point Spread" {
			func(m markets) {

				hotc, votc := assignOutcomes(m)

				odd.HomeSpreadLine = game.DecimalToAmerican(hotc.Prices[0].Decimal)
				odd.VisitingSpreadLine = game.DecimalToAmerican(votc.Prices[0].Decimal)

				odd.HomeSpreadPoints, err = strconv.ParseFloat(hotc.Prices[0].HandicapLow, 64)
				odd.VisitingSpreadPoints, err = strconv.ParseFloat(votc.Prices[0].HandicapLow, 64)
				// Ignore the error for now. Worst case scenario is there is no data. not fatal

			}(market)
		}
		if market.Name == "Total Points Over/Under" {
			func(m markets) {

				over, under := func(m markets) (ho outcomes, vo outcomes) {

					for _, otc := range m.Outcomes {
						if otc.Name == "Over" {
							ho = otc
						}
						if otc.Name == "Under" {
							vo = otc
						}
					}
					return ho, vo

				}(market)

				odd.GameOverLine = game.DecimalToAmerican(over.Prices[0].Decimal)
				odd.GameUnderLine = game.DecimalToAmerican(under.Prices[0].Decimal)

				odd.GameOver, err = strconv.ParseFloat(over.Prices[0].HandicapLow, 64)
				// Ignore the error for now. Worst case scenario is there is no data. not fatal

			}(market)
		}
	}

	return odd, err

}

func linesFromTimeBand(tbe timeBandEvents) ([]game.Line, error) {

	var odds []game.Line
	for _, v := range tbe.Events {
		// If not a desired sport, continue
		if _, ok := SupportedSports[v.Category.Name]; !ok {
			continue
		}
		o, err := lineFromEvent(v)
		if err != nil {
			return odds, err
		}
		odds = append(odds, o)
	}
	return odds, nil

}

func FetchLines(league string) ([]game.Line, error) {

	var includeDates []time.Time

	for i := 0; i < 7; i += 1 {
		includeDates = append(includeDates, time.Now().AddDate(0, 0, i))
	}

	if _, ok := leagueTagIds[league]; !ok {
		return nil, errors.New(fmt.Sprintf("League not supported: %s", league))
	}

	//log.Println(constructUrl(league, includeDates))
	res, err := http.Get(constructUrl(league, includeDates))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%s", b)

	var raw rawResponse

	err = json.Unmarshal(b, &raw)

	if err != nil {
		return nil, err
	}

	var odds []game.Line

	for _, v := range raw.Data.TimeBandEvents {
		if v.Type == "LIVE" {
			continue
		}
		o, err := linesFromTimeBand(v)
		if err != nil {
			return odds, err
		}

		odds = append(odds, o...)
	}

	return odds, nil

}
