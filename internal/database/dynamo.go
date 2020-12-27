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
	"log"
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

// FetchProjection returns the projection matching the given game info and source
func FetchProjectionV2(gameDate time.Time, source string) ([]game.Projection, error) {

	var p []game.Projection

	tableName := "win-projections"
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
	log.Printf("%#v", p)

	// Rename to common
	var nameMap map[string]string
	if source == "FTEQBELO" {
		nameMap = fte.ToCommonNflName
	} else if source == "FTERAPTORELO"{
		nameMap = fte.ToCommonNbaName
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

func PutGameProjection(p game.Projection) error {

	tableName := "win-projections"

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
