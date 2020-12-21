package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fischersean/linescrape/internal/database"

	"bytes"
	"encoding/json"
	"log"
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

	database.Init()
	odds, err := database.FetchLatestNflOdds("mybookie")

	if err != nil {
		return lpr, err
	}

	log.Println("#v", odds)
	for _, v := range odds {
		proj, err := database.FetchProjection(v, "FTEQBELO")
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

func main() {

	lambda.Start(Handler)

}
