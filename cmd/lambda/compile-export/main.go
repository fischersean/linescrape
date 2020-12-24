package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fischersean/linescrape/internal/database"
	"github.com/fischersean/linescrape/pkg/game"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	//"log"
)

type Request struct {
	League     string `json:"league"`
	Source     string `json:"source"`
	Projection string `json:"projection"`
}

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

var (
	enabledLeagues = map[string]bool{
		"NFL":                true,
		"NBA":                true,
		"college-football":   true,
		"college-basketball": true,
	}

	enabledSources = map[string]bool{
		"mybookie": true,
		"caesars":  true,
	}

	enabledProjections = map[string]bool{
		"FTEQBELO": true,
		"none":     true,
	}
)

func validateRequest(r Request) error {

	// this func would be a great candiate for the go 1.16 file embed feature
	if v, ok := enabledLeagues[r.League]; !v || !ok {
		return errors.New(fmt.Sprintf("League not supported: %s", r.League))
	}
	if v, ok := enabledSources[r.Source]; !v || !ok {
		return errors.New(fmt.Sprintf("Line source not supported: %s", r.Source))
	}
	if v, ok := enabledProjections[r.Projection]; !v || !ok {
		return errors.New(fmt.Sprintf("Projection source not supported: %s", r.Projection))
	}

	return nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// we are relying on API gateway to guarantee these are present
	r := Request{}
	r.League = request.QueryStringParameters["league"]
	r.Source = request.QueryStringParameters["source"]
	r.Projection = request.QueryStringParameters["projection"]
	// If no param provided, default to none
	if r.Projection == "" {
		r.Projection = "none"
	}

	if err := validateRequest(r); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	lpr, err := generateExport(r.League, r.Source, r.Projection)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	body, err := json.Marshal(lpr)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: buf.String(),
	}, err

}

func generateExport(league string, source string, projection string) (lpr []LineProjectionResult, err error) {

	database.Init()
	//odds, err := database.FetchLatestNflOdds("mybookie")
	odds, err := database.FetchLatestOdds(source, league)

	if err != nil {
		return lpr, err
	}

	//log.Printf("%#v", odds)
	for _, v := range odds {
		proj := game.Projection{}
		if projection != "none" {
			proj, err = database.FetchProjection(v, projection)
			if err != nil && err.Error() != "Could not find item matching query expression" {
				return lpr, err
			}
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

	return lpr, nil

}

func main() {

	lambda.Start(Handler)

}
