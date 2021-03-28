package mybookie

import (
	"github.com/gocolly/colly"

	"github.com/fischersean/linescrape/pkg/game"

	"errors"
	"fmt"
	"log"
)

var (
	SupportedLeagues = []string{
		"NFL",
		"college-football",
		"NBA",
		"college-basketball",
		"MLB",
	}

	SupportedSports = map[string]string{
		"Football":   "",
		"Basketball": "",
		"Baseball":   "",
	}
)

func FetchLines(league string) ([]game.Line, error) {
	requestNames := map[string]string{
		"NFL":                "nfl",
		"NBA":                "nba",
		"MLB":                "",
		"college-football":   "college-football",
		"college-basketball": "ncaa-basketball",
	}

	cleanName, ok := requestNames[league]
	if !ok {
		return []game.Line{}, errors.New(fmt.Sprintf("%s: %s", "Invalid league type", league))
	}

	siteUrl := fmt.Sprintf("https://mybookie.ag/sportsbook/%s/", cleanName)

	var odds []game.Line

	c := colly.NewCollector()

	c.OnHTML(".game-line", func(e *colly.HTMLElement) {
		g, err := ParseOdds(e)
		if err != nil {
			log.Printf("Could not parse element: %s", err.Error())
			return
		}
		odds = append(odds, g)
	})

	err := c.Visit(siteUrl)

	if err != nil {
		return []game.Line{}, errors.New("Could not reach mybookie.ag")
	}

	return odds, nil
}
