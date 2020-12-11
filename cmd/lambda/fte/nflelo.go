package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/gocarina/gocsv"

	"github.com/fischersean/linescrape/pkg/fte"
	"github.com/fischersean/linescrape/pkg/game"

	"encoding/csv"
	"fmt"
	"io"
	//"log"
	"net/http"
	"strings"
	"time"
)

func marshalProjections(r io.Reader) (p []fte.EloProjection, err error) {

	cr := csv.NewReader(r)
	err = gocsv.UnmarshalCSV(cr, &p)

	return p, err
}

func refreshProjDynamo(sess *session.Session, f io.Reader) (err error) {

	getGameId := func(p fte.EloProjection) string {
		parts := strings.Split(p.Date, "-")
		return fmt.Sprintf("%s%s%s%s%s%s",
			"NFL",
			parts[0], parts[1], parts[2],
			p.Team1, p.Team2,
		)
	}

	gameDate := func(p fte.EloProjection) time.Time {

		t, err := time.Parse("2006-01-02", p.Date)

		if err != nil {
			return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
		}

		return t
	}

	shouldUpdate := func(v game.Projection) bool {

		yesterday := func(t time.Time) bool {
			s := time.Since(t)
			return s <= 36*time.Hour && s > 0
		}(v.GameDate)

		withinSixDays := func(t time.Time) bool {
			u := time.Until(t)
			return u <= 6*24*time.Hour && u > 0
		}(v.GameDate)

		//if yesterday || withinSixDays {
		//log.Printf("Updating game on %s", v.GameDate)
		//if yesterday {
		//log.Printf("Is yesterday")
		//}
		//if withinSixDays {
		//log.Println("Is within +6 days")
		//}
		//}

		return yesterday || withinSixDays
	}

	var projs []game.Projection

	proj, err := marshalProjections(f)

	if err != nil {
		return err
	}

	for _, v := range proj {
		p := game.Projection{
			GameId:                 getGameId(v),
			League:                 "NFL",
			Source:                 "FTEQBELO",
			Home:                   v.Team1,
			Visiting:               v.Team2,
			HomeWinProbability:     v.QBEloProb1,
			VisitingWinProbability: v.QBEloProb2,
			GameDate:               gameDate(v),
			Season:                 v.Season,
			Playoff:                v.Playoff,
		}

		projs = append(projs, p)
	}

	//projs = projs[0:20]

	svc := dynamodb.New(sess)
	for _, v := range projs {
		if !shouldUpdate(v) {
			continue
		}

		av, err := dynamodbattribute.MarshalMap(v)
		if err != nil {
			return err
		}

		input := &dynamodb.PutItemInput{
			//ExpressionAttributeNames: map[string]*string{
			//"#H":   aws.String("home"),
			//"#V":   aws.String("visiting"),
			//"#HWP": aws.String("homeWinProb"),
			//"#VWP": aws.String("visitingWinProb"),
			//"#SC":  aws.String("source"),
			//"#GD":  aws.String("gameDate"),
			//"#LG":  aws.String("league"),
			//"#S":   aws.String("season"),
			//"#PO":  aws.String("playoff"),
			//},
			//ExpressionAttributeValues: av,
			//Key: map[string]*dynamodb.AttributeValue{
			//"gameId": {
			//S: aws.String(v.GameId),
			//},
			//},
			Item:      av,
			TableName: aws.String("win-projections"),
			//UpdateExpression: aws.String(`SET #H = :home,
			//#V = :visiting,
			//#HWP = :homeWinProb,
			//#VWP = :visitingWinProb,
			//#SC = :source,
			//#LG = :league,
			//#GD = :gameDate,
			//#S = :season,
			//#PO = :playoff`),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			return err
		}
	}

	return err
}

func refreshEloS3(sess *session.Session, f io.Reader) (err error) {
	bucket := "game-odds"
	filename := "nfl-elo-latest.csv"

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   f,
	})

	return err

}

func Handler() error {

	url := "https://projects.fivethirtyeight.com/nfl-api/nfl_elo_latest.csv"

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	if err != nil {
		return err
	}

	err = refreshProjDynamo(sess, res.Body)
	if err != nil {
		return err
	}

	err = refreshEloS3(sess, res.Body)

	return err

}

func main() {
	//err := Handler()
	//log.Printf("%#v", err)
	lambda.Start(Handler)
}
