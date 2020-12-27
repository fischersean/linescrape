package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fischersean/linescrape/internal/database"
	"github.com/fischersean/linescrape/pkg/game"

	"bytes"
	"encoding/json"
	//"errors"
	//"fmt"
	"log"
	"time"
)

type Request struct {
	Source    string      `json:"source"`
	GameDates []time.Time `json:"gameDate"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// we are relying on API gateway to guarantee the correct params
	r := Request{}
	r.Source = request.QueryStringParameters["source"]

	var gDatesRaw []string
	if v, ok := request.MultiValueQueryStringParameters["gameDate"]; ok {
		gDatesRaw = v
	} else {
		gDatesRaw = []string{
			request.QueryStringParameters["gameDate"],
		}
	}

	for _, v := range gDatesRaw {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, err
		}
		r.GameDates = append(r.GameDates, t)

	}

	lpr, err := generateExport(r.Source, r.GameDates)
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

func generateExport(source string, gameDates []time.Time) (projs []game.Projection, err error) {

	database.Init()

	for _, v := range gameDates {
		proj, err := database.FetchProjectionV2(v, source)
		if err != nil && err.Error() != "Could not find item matching query expression" {
			return projs, err
		}
		log.Printf("%#v", proj)

		projs = append(projs, proj...)
	}

	return projs, nil

}

func main() {

	lambda.Start(Handler)

}
