package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/go-gota/gota/dataframe"

	mb "github.com/fischersean/linescrape/pkg/mybookie"

	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type GameOddsItem struct {
	TimeStamp time.Time     `json:"time_stamp"`
	Odds      []mb.GameLine `json:"odds"`
	League    string        `json:"league"`
}

func getLatestNflOdds() (odds []mb.GameLine, err error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

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

	item := []GameOddsItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)

	return item[0].Odds, err
}

func getLatestNflElo() (*os.File, error) {

	bucket := "game-odds"
	item := "nfl-elo-latest.csv"

	file, err := ioutil.TempFile("tmp", "*")
	if err != nil {
		return file, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	if err != nil {
		return file, err
	}

	return file, err

}

func main() {
	// Pull latest elo data

	file, err := getLatestNflElo()
	defer file.Close()

	odds, err := getLatestNflOdds()

	if err != nil {
		log.Fatalf(err.Error())
	}

	//log.Printf("%#v", odds)

	df := dataframe.LoadStructs(odds)
	df2 := dataframe.ReadCSV(file)

	 //Load into a dataframe
	log.Println(df)
	log.Println(df2)

}
