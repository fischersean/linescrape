package fte

import (
	"github.com/gocarina/gocsv"

	"github.com/fischersean/linescrape/pkg/game"

	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func getGameId(league, date, t1, t2 string) string {

	parts := strings.Split(date, "-")
	return fmt.Sprintf("%s%s%s%s%s%s",
		league,
		parts[0], parts[1], parts[2],
		t1, t2,
	)

}

func gameDate(date string) time.Time {

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}

	return t

}

func fetchNba() (map[string][]game.Projection, error) {

	res, err := http.Get(NbaLatestUrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	projs := map[string][]game.Projection{}
	var fteProjs []NbaEloProjection

	cr := csv.NewReader(res.Body)
	err = gocsv.UnmarshalCSV(cr, &fteProjs)

	if err != nil {
		return nil, err
	}

	for _, v := range fteProjs {
		p0 := game.Projection{
			GameId:                 getGameId("NBA", v.Date, v.Team1, v.Team2),
			League:                 "NBA",
			Source:                 "FTERAPTORELO",
			Home:                   v.Team1,
			Visiting:               v.Team2,
			HomeWinProbability:     v.RaptorProb1,
			VisitingWinProbability: v.RaptorProb2,
			GameDate:               gameDate(v.Date),
			Season:                 v.Season,
			Playoff:                v.Playoff,
		}
		projs["FTERAPTORELO"] = append(projs["FTERAPTORELO"], p0)

		p1 := game.Projection{
			GameId:                 getGameId("NBA", v.Date, v.Team1, v.Team2),
			League:                 "NBA",
			Source:                 "FTEELO",
			Home:                   v.Team1,
			Visiting:               v.Team2,
			HomeWinProbability:     v.EloProb1,
			VisitingWinProbability: v.EloProb2,
			GameDate:               gameDate(v.Date),
			Season:                 v.Season,
			Playoff:                v.Playoff,
		}
		projs["FTEELO"] = append(projs["FTEELO"], p1)
	}

	return projs, nil

}

func fetchNfl() (map[string][]game.Projection, error) {

	res, err := http.Get(NflLatestUrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	projs := map[string][]game.Projection{}
	var fteProjs []NflEloProjection

	cr := csv.NewReader(res.Body)
	err = gocsv.UnmarshalCSV(cr, &fteProjs)

	if err != nil {
		return nil, err
	}

	for _, v := range fteProjs {
		p0 := game.Projection{
			GameId:                 getGameId("NFL", v.Date, v.Team1, v.Team2),
			League:                 "NFL",
			Source:                 "FTEQBELO",
			Home:                   v.Team1,
			Visiting:               v.Team2,
			HomeWinProbability:     v.QBEloProb1,
			VisitingWinProbability: v.QBEloProb2,
			GameDate:               gameDate(v.Date),
			Season:                 v.Season,
			Playoff:                v.Playoff,
		}
		projs["FTEQBELO"] = append(projs["FTEQBELO"], p0)

		p1 := game.Projection{
			GameId:                 getGameId("NFL", v.Date, v.Team1, v.Team2),
			League:                 "NFL",
			Source:                 "FTEELO",
			Home:                   v.Team1,
			Visiting:               v.Team2,
			HomeWinProbability:     v.EloProb1,
			VisitingWinProbability: v.EloProb2,
			GameDate:               gameDate(v.Date),
			Season:                 v.Season,
			Playoff:                v.Playoff,
		}
		projs["FTEELO"] = append(projs["FTEELO"], p1)
	}

	return projs, nil

}

// FetchLatest retrieves the latest projections from fivethirtyeight.com
func FetchLatest(league string) (map[string][]game.Projection, error) {
	leagueMap := map[string]func() (map[string][]game.Projection, error){
		"NFL": fetchNfl,
		"NBA": fetchNba,
	}

	if v, ok := leagueMap[league]; ok {
		return v()
	}

	return nil, errors.New(fmt.Sprintf("League not supported: %s", league))
}
