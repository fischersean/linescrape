package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/fischersean/linescrape/pkg/fte"
	"github.com/fischersean/linescrape/pkg/game"

	"errors"
	"fmt"
	//"log"
	"strings"
	"time"
)

var (
	// Session is the shared aws session
	Session *session.Session

	// Service is the shared dynamo db connection
	Service *dynamodb.DynamoDB
)

// GameOddsItem is the structure expected for a single item from 'game-odds'
type GameOddsItem struct {
	TimeStamp time.Time   `json:"time_stamp"`
	League    string      `json:"league"`
	Odds      []game.Line `json:"odds"`
	Source    string      `json:"source"`
}

// Init connects the database to the globaly configured dynamo db
func Init() {

	Session = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	Service = dynamodb.New(Session)

}

// FetchLatestNflOdds returns the most recently added set of game odds
func FetchLatestOdds(dataSource string, league string) ([]game.Line, error) {

	var odds []game.Line

	tableName := "game-odds"

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(league),
			},
			":v2": {
				S: aws.String(dataSource),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#L": aws.String("league"),
			"#S": aws.String("source"),
		},
		KeyConditionExpression: aws.String("#L = :v1"),
		FilterExpression:       aws.String("#S = :v2"),
		TableName:              aws.String(tableName),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int64(2), // Will alwways be within n * # of sources of top
	}

	result, err := Service.Query(input)

	if err != nil {
		return odds, err
	}

	if *result.Count != 1 {
		return odds, errors.New("Could not find item matching query expression")
	}

	item := []GameOddsItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)

	return item[0].Odds, err

}

// FetchLatestNflOdds returns the most recently added set of game odds
func FetchLatestNflOdds(dataSource string) ([]game.Line, error) {

	var odds []game.Line

	tableName := "game-odds"
	league := "nfl"

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(league),
			},
			":v2": {
				S: aws.String(dataSource),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("source"),
		},
		KeyConditionExpression: aws.String("league = :v1"),
		FilterExpression:       aws.String("#S = :v2"),
		TableName:              aws.String(tableName),
		ScanIndexForward:       aws.Bool(false),
		Limit:                  aws.Int64(1),
	}

	result, err := Service.Query(input)

	if err != nil {
		return odds, err
	}

	if *result.Count != 1 {
		return odds, errors.New("Could not find item matching query expression")
	}

	item := []GameOddsItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)

	return item[0].Odds, err

}

// FetchProjection returns the projection matching the given game info and source
func FetchProjectionV2(gameDate time.Time, source string) ([]game.Projection, error) {

	var p []game.Projection

	//log.Printf("%s", gid)
	tableName := "game-projections"
	indexName := "ByDate"

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(source),
			},
			":v2": {
				S: aws.String(gameDate.Format(time.RFC3339)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#S":  aws.String("source"),
			"#GD": aws.String("gameDate"),
		},
		KeyConditionExpression: aws.String("#S = :v1 and #GD = :v2"),
		TableName:              aws.String(tableName),
		IndexName:              aws.String(indexName),
	}

	result, err := Service.Query(input)

	if err != nil {
		return p, err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &p)

	// Rename to common
	var nameMap map[string]string
	if source == "FTEQBELO" {
		nameMap = fte.ToCommon
	} else {
		// we will just return early and not rename the items
		return p, nil
	}

	for i, v := range p {
		fmt.Println(p[i])
		p[i].Home = nameMap[v.Home]
		p[i].Visiting = nameMap[v.Visiting]
	}

	return p, err

}

// FetchProjection returns the projection matching the given game info and source
func FetchProjection(odds game.Line, projectionSource string) (game.Projection, error) {

	var p game.Projection

	gameDate, err := time.Parse("2006-01-02 15:04:05", odds.GameTime)
	if err != nil {
		// Try the caesars format
		parts := strings.Split(odds.GameTime, " ")
		gameDate, err = time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", parts[0], parts[1]))

		if err != nil {
			return p, err

		}
	}

	var nameMap map[string]string

	if projectionSource == "FTEQBELO" {
		nameMap = fte.FromCommon
	} else {
		return p, errors.New(fmt.Sprintf("Projection source not supported: %s", projectionSource))
	}

	gid := fmt.Sprintf("%s%s%s%s",
		"NFL",
		gameDate.Format("20060102"),
		nameMap[odds.HomeTeam],
		nameMap[odds.VisitingTeam],
	)

	//log.Printf("%s", gid)
	tableName := "game-projections"

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(gid),
			},
			":v2": {
				S: aws.String(projectionSource),
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

	result, err := Service.Query(input)

	if err != nil {
		return p, err
	}

	if *result.Count != 1 {
		//log.Printf("Error finding projections for %s", gid)
		return p, errors.New("Could not find item matching query expression")
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &p)

	return p, err

}

func PutGameProjection(p game.Projection) error {

	tableName := "game-projections"

	av, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = Service.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// PutGameOddsItem puts a single game odds item into the game-odds table
func PutGameOddsItem(odds GameOddsItem) error {

	tableName := "game-odds"

	av, err := dynamodbattribute.MarshalMap(odds)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	//log.Printf("%#v", odds)
	_, err = Service.PutItem(input)

	return err

}
