package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/gocolly/colly"

	"github.com/fischersean/linescrape/pkg/mybookie"
	"github.com/fischersean/linescrape/pkg/game"
)

type Request struct {
	League string `json:"league"`
}

type Resonse struct {
	TimeStamp time.Time           `json:"time_stamp"`
	Odds      []game.Line `json:"odds"`
	League    string              `json:"league"`
}

func putResponsDB(res Resonse) (err error) {

	tableName := "game-odds"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(res)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)

	return err
}

func Handler(request Request) (Resonse, error) {
	// TODO: Need to abstract away the calls to colly in this function

	log.Printf("%+v", request)

	if request.League != "nfl" && request.League != "college-football" {
		return Resonse{}, errors.New(fmt.Sprintf("%s: %s", "Invalid league type", request.League))
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
		return Resonse{}, errors.New("Could not reach mybookie.ag")
	}

	res := Resonse{
		Odds:      odds,
		TimeStamp: time.Now(),
		League:    request.League,
	}

	if err = putResponsDB(res); err != nil {
		return Resonse{}, err
	}

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
