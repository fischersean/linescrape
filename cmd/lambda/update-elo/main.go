package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fischersean/linescrape/internal/database"
	"github.com/fischersean/linescrape/pkg/fte"
	"github.com/fischersean/linescrape/pkg/game"

	"log"
	"time"
)

func shouldUpdate(v game.Projection) bool {

	yesterday := func(t time.Time) bool {
		s := time.Since(t)
		return s <= 36*time.Hour && s > 0
	}(v.GameDate)

	withinSixDays := func(t time.Time) bool {
		u := time.Until(t)
		return u <= 6*24*time.Hour && u > 0
	}(v.GameDate)

	return yesterday || withinSixDays
}

func putLeague(leagueProjs map[string][]game.Projection) error {

	for _, v := range leagueProjs {
		for _, v2 := range v {
			if !shouldUpdate(v2) {
				continue
			}
			if err := database.PutGameProjection(v2); err != nil {
				return err
			}
		}
	}

	return nil

}

func Handler() error {

	leagueList := []string{
		"NFL",
		"NBA",
	}

	database.Init()
	for _, v := range leagueList {
		p, err := fte.FetchLatest(v)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if err = putLeague(p); err != nil {
			log.Fatalf(err.Error())
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
