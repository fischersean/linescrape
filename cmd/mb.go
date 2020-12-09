package main

import (
	"encoding/json"
	"github.com/fischersean/linescrape/pkg/mybookie"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func main() {
	siteUrl := "https://mybookie.ag/sportsbook/college-football/"

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

	b, err := json.Marshal(odds)

	if err != nil {
		log.Fatalf(err.Error())
	}

	os.Stdout.Write(b)
}
