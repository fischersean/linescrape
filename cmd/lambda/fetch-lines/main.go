package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fischersean/linescrape/internal/database"
	"github.com/fischersean/linescrape/pkg/caesars"
	"github.com/fischersean/linescrape/pkg/game"
	"github.com/fischersean/linescrape/pkg/mybookie"

	"errors"
	"fmt"
	"log"
	"time"
)

type Request struct {
	Source string `json:"source"`
	League string `json:"league"`
}

func Handler(request Request) (database.GameOddsItem, error) {
	// TODO: Need to abstract away the calls to colly in this function

	log.Printf("%+v", request)

	sourceMap := map[string]func(string) ([]game.Line, error){
		"mybookie": mybookie.FetchLines,
		"caesars":  caesars.FetchLines,
	}
	if _, ok := sourceMap[request.Source]; !ok {
		return database.GameOddsItem{}, errors.New(fmt.Sprintf("Source not supported: %s", request.Source))
	}

	odds, err := sourceMap[request.Source](request.League)
	if err != nil {
		return database.GameOddsItem{}, err
	}

	res := database.GameOddsItem{
		Odds:      odds,
		TimeStamp: time.Now(),
		League:    request.League,
		Source:    request.Source,
	}

	database.Init()
	err = database.PutGameOddsItem(res)

	return res, err
}

func main() {
	lambda.Start(Handler)
}
