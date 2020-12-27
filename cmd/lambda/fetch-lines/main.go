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
	League string `json:"league"`
	Source string `json:"source"`
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
)

func validateRequest(r Request) error {

	// this func would be a great candiate for the go 1.16 file embed feature
	if v, ok := enabledLeagues[r.League]; !v || !ok {
		return errors.New(fmt.Sprintf("League not supported: %s", r.League))
	}
	if v, ok := enabledSources[r.Source]; !v || !ok {
		return errors.New(fmt.Sprintf("Line source not supported: %s", r.Source))
	}
	return nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// we are relying on API gateway to guarantee these are present
	r := Request{}
	r.League = request.QueryStringParameters["league"]
	r.Source = request.QueryStringParameters["source"]

	if err := validateRequest(r); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	lpr, err := generateExport(r.League, r.Source)
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

func generateExport(league string, source string) (odds []game.Line, err error) {

	database.Init()
	odds, err = database.FetchLatestOdds(source, league)

	if err != nil {
		return odds, err
	}

	return odds, nil

}

func main() {

	lambda.Start(Handler)

}
