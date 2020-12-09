// snippet-comment:[These are tags for the AWS doc team's sample catalog. Do not remove.]
// snippet-sourceauthor:[Doug-AWS]
// snippet-sourcedescription:[DynamoDBCreateItem.go creates an item in an Amazon DynamoDB table.]
// snippet-keyword:[Amazon DynamoDB]
// snippet-keyword:[PutItem function]
// snippet-keyword:[Go]
// snippet-sourcesyntax:[go]
// snippet-service:[dynamodb]
// snippet-keyword:[Code Sample]
// snippet-sourcetype:[full-example]
// snippet-sourcedate:[2019-03-19]
/*
   Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/
// snippet-start:[dynamodb.go.create_item]
package main

// snippet-start:[dynamodb.go.create_item.imports]
import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fischersean/linescrape/pkg/mybookie"
	"github.com/gocolly/colly"
	"log"
	"os"
	"time"
)

// snippet-end:[dynamodb.go.create_item.imports]

// snippet-start:[dynamodb.go.create_item.struct]
// Create struct to hold info about new item
type Request struct {
	League string `json:"league"`
}

type Resonse struct {
	TimeStamp time.Time           `json:"time_stamp"`
	Odds      []mybookie.GameLine `json:"odds"`
	League    string              `json:"league"`
}

// snippet-end:[dynamodb.go.create_item.struct]

func main() {
	request := Request{
		League: "nfl",
	}
	// snippet-start:[dynamodb.go.create_item.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.create_item.session]

	// snippet-start:[dynamodb.go.create_item.assign_struct]
	siteUrl := fmt.Sprintf("https://mybookie.ag/sportsbook/%s/", request.League)

	var odds []mybookie.GameLine

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
		log.Fatalf(err.Error())
	}

	item := Resonse{
		Odds:      odds,
		TimeStamp: time.Now(),
		League:    request.League,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// snippet-end:[dynamodb.go.create_item.assign_struct]

	// snippet-start:[dynamodb.go.create_item.call]
	// Create item in table Movies
	tableName := "game-odds"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
