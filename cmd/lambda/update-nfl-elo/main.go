package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/gocarina/gocsv"

	"github.com/fischersean/linescrape/internal/database"
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

func refreshProjDynamo(f io.Reader) (err error) {

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

	for _, v := range projs {

		if !shouldUpdate(v) {
			continue
		}

		database.Init()
		err = database.PutGameProjection(v)
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

	err = refreshProjDynamo(res.Body)
	if err != nil {
		return err
	}

	err = refreshEloS3(sess, res.Body)

	return err

}

func main() {
	lambda.Start(Handler)
}
