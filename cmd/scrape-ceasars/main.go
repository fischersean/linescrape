package main

import (
	caesars "github.com/fischersean/linescrape/pkg/ceasars"

	"github.com/fischersean/linescrape/internal/database"

	"log"
	"time"
)

func main() {
	o, err := caesars.FetchLines("NBA")

	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, v := range o {
		log.Printf("%s v %s: %d, %d", v.HomeTeam, v.VisitingTeam, v.HomeMoneyLine, v.VisitingMoneyLine)
	}

	_ = database.GameOddsItem{
		//pItem := database.GameOddsItem{
		TimeStamp: time.Now(),
		League:    "nba",
		Source:    "caesars",
		Odds:      o,
	}

	//database.Init()

	//err = database.PutGameOddsItem(pItem)

	if err != nil {
		log.Fatalf(err.Error())
	}

}
