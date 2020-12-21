package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/gocolly/colly"

	"github.com/fischersean/linescrape/internal/database"
	"github.com/fischersean/linescrape/pkg/game"
	"github.com/fischersean/linescrape/pkg/mybookie"

	"errors"
	"fmt"
	"log"
	"time"
)

type Request struct {
	League string `json:"league"`
}

func Handler(request Request) (database.GameOddsItem, error) {
	// TODO: Need to abstract away the calls to colly in this function

	log.Printf("%+v", request)

	if request.League != "nfl" && request.League != "college-football" {
		return database.GameOddsItem{}, errors.New(fmt.Sprintf("%s: %s", "Invalid league type", request.League))
	}

	siteUrl := fmt.Sprintf("https://mybookie.ag/sportsbook/%s/", request.League)

	var odds []game.Line

	c := colly.NewCollector()

	c.OnHTML(".game-line", func(e *colly.HTMLElement) {
		g, err := mybookie.ParseOdds(e)
		if err != nil {
			log.Fatalf(err.Error())
		}
		odds = append(odds, g)
	})

	err := c.Visit(siteUrl)

	if err != nil {
		return database.GameOddsItem{}, errors.New("Could not reach mybookie.ag")
	}

	res := database.GameOddsItem{
		Odds:      odds,
		TimeStamp: time.Now(),
		League:    request.League,
		Source:    "mybookie",
	}

	database.Init()
	err = database.PutGameOddsItem(res)

	return res, err
}

func main() {
	lambda.Start(Handler)
}
