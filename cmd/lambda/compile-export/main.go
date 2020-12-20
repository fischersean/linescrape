package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/fischersean/linescrape/pkg/fte"
	"github.com/fischersean/linescrape/pkg/game"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type Response struct {
	Body       string            `json:"body"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
}

type LineProjectionResult struct {
	Home                   string
	Visiting               string
	HomeMoneyLine          int64
	VisitingMoneyLine      int64
	HomeWinProbability     float64
	VisitingWinProbability float64
	GameTime               string
}

func Handler() (Response, error) {

	lpr, err := generateExport()
	if err != nil {
		return Response{StatusCode: 504}, err
	}

	body, err := json.Marshal(lpr)
	if err != nil {
		return Response{StatusCode: 504}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	return Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: buf.String(),
	}, err
}

func generateExport() (lpr []LineProjectionResult, err error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	// Fetch the latest lines
	odds, err := getLatestNflOdds(svc)

	if err != nil {
		return lpr, err
	}

	log.Println("#v", odds)
	for _, v := range odds {
		proj, err := fetchProjection(svc, v)
		if err != nil && err.Error() != "Could not find item matching query expression" {
			return lpr, err
		}

		lpr = append(lpr, LineProjectionResult{
			Home:                   v.HomeTeam,
			Visiting:               v.VisitingTeam,
			GameTime:               v.GameTime,
			HomeMoneyLine:          v.HomeMoneyLine,
			VisitingMoneyLine:      v.VisitingMoneyLine,
			HomeWinProbability:     proj.HomeWinProbability,
			VisitingWinProbability: proj.VisitingWinProbability,
		})
	}

	return lpr, err
}

func getLatestNflOdds(svc *dynamodb.DynamoDB) (odds []game.Line, err error) {

	tableName := "game-odds"
	league := "nfl"

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(league),
			},
		},
		KeyConditionExpression: aws.String("league = :v1"),
		TableName:              aws.String(tableName),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int64(1),
	}

	result, err := svc.Query(input)

	if err != nil {
		return odds, err
	}

	if *result.Count != 1 {
		return odds, errors.New("Could not find item matching query expression")
	}

	item := []game.GameOddsTableItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)

	return item[0].Odds, err
}

func fetchProjection(svc *dynamodb.DynamoDB, odds game.Line) (p game.Projection, err error) {

	gameDate, err := time.Parse("2006-01-02 15:04:05", odds.GameTime)

	if err != nil {
		return p, err
	}

	gid := fmt.Sprintf("%s%s%s%s",
		"NFL",
		gameDate.Format("20060102"),
		fte.FromCommon[odds.HomeTeam],
		fte.FromCommon[odds.VisitingTeam],
	)

	log.Printf("%s", gid)
	tableName := "win-projections"

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(gid),
			},
			":v2": {
				S: aws.String("FTEQBELO"),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("source"),
		},
		KeyConditionExpression: aws.String("gameId = :v1"),
		FilterExpression:       aws.String("#S = :v2"),
		TableName:              aws.String(tableName),
		Limit:                  aws.Int64(1),
	}

	result, err := svc.Query(input)

	if err != nil {
		return p, err
	}

	if *result.Count != 1 {
		log.Printf("Error finding projections for %s", gid)
		return p, errors.New("Could not find item matching query expression")
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &p)

	return p, err

}

func main() {

	lambda.Start(Handler)

}
