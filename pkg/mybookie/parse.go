package mybookie

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/fischersean/linescrape/pkg/game"
	"github.com/gocolly/colly"
	"strconv"
)

type spread struct {
	Team   string
	VsTeam string
	Points float64
	Line   int64
}

type moneyLine struct {
	Team   string
	VsTeam string
	Line   int64
}

//type total struct {
//Team   string
//VsTeam string
//Points float64
//Line   int64
//}

type tagData struct {
	Team   string
	VsTeam string
	Spread float64
	Odds   int64
}

// Line is the basic data structure used within betting
//type Line struct {
//GameTime             string
//HomeTeam             string
//VisitingTeam         string
//HomeSpreadPoints     float64
//VisitingSpreadPoints float64
//HomeSpreadLine       int64
//VisitingSpreadLine   int64
//HomeMoneyLine        int64
//VisitingMoneyLine    int64
//GameOver             float64
//GameOverLine         int64
//GameUnderLine        int64
//}

func parseTags(dom *goquery.Selection, selector string) (t tagData, err error) {

	sd := dom.Find(selector)

	var gSpreadPoints string
	var gHSpreadLine string

	var exists bool

	if t.Team, exists = sd.Attr("data-team"); !exists {
		return t, errors.New("Could not locate team data")
	}

	if t.VsTeam, exists = sd.Attr("data-team-vs"); !exists {
		return t, errors.New("Could not locate team data")
	}

	if gSpreadPoints, exists = sd.Attr("data-spread"); !exists {
		return t, errors.New("Could not locate spread data")
	}

	if t.Spread, err = strconv.ParseFloat(gSpreadPoints, 64); err != nil && gSpreadPoints != "" {
		return t, err
	}

	if gHSpreadLine, exists = sd.Attr("data-odds"); !exists {
		return t, errors.New("Could not locate spread odds data")
	}

	if t.Odds, err = strconv.ParseInt(gHSpreadLine, 10, 64); err != nil && gHSpreadLine != "" {
		return t, err
	}
	return t, nil
}

func parseSpread(dom *goquery.Selection) (s spread, err error) {

	t, err := parseTags(dom, "[data-wager-type='sp']")

	if err != nil {
		return s, err
	}

	s.Team = t.Team
	s.VsTeam = t.VsTeam
	s.Points = t.Spread
	s.Line = t.Odds

	return s, err
}

func parseGameSpread(homeDOM *goquery.Selection, visitingDOM *goquery.Selection, odds *game.Line) (err error) {

	vSpread, err := parseSpread(visitingDOM)

	if err != nil {
		return err
	}

	hSpread, err := parseSpread(homeDOM)

	if err != nil {
		return err
	}

	odds.HomeSpreadPoints = hSpread.Points
	odds.VisitingSpreadPoints = vSpread.Points
	odds.HomeSpreadLine = hSpread.Line
	odds.VisitingSpreadLine = vSpread.Line

	return err
}

func parseMoneyLine(dom *goquery.Selection) (m moneyLine, err error) {

	t, err := parseTags(dom, "[data-wager-type='ml']")

	if err != nil {
		return m, err
	}

	m.Team = t.Team
	m.VsTeam = t.VsTeam
	m.Line = t.Odds

	return m, err
}

func parseGameMoneyLine(homeDOM *goquery.Selection, visitingDOM *goquery.Selection, odds *game.Line) (err error) {

	vSpread, err := parseMoneyLine(visitingDOM)

	if err != nil {
		return err
	}

	hSpread, err := parseMoneyLine(homeDOM)

	if err != nil {
		return err
	}

	odds.HomeMoneyLine = hSpread.Line
	odds.VisitingMoneyLine = vSpread.Line

	return err
}

func parseGameOver(homeDOM *goquery.Selection, visitingDOM *goquery.Selection, odds *game.Line) (err error) {

	th, err := parseTags(homeDOM, "[data-wager-type='to']")

	if err != nil {
		return err
	}

	tv, err := parseTags(visitingDOM, "[data-wager-type='to']")

	if err != nil {
		return err
	}

	odds.GameOver = th.Spread
	odds.GameOverLine = tv.Odds
	odds.GameUnderLine = th.Odds

	return err
}

// ParseOdds parses a HTML doc pulled from mybookie's nfl page and return a slice of odds
func ParseOdds(e *colly.HTMLElement) (odds game.Line, err error) {

	// Find game date
	var exists bool
	gtime := e.DOM.Find(".game-line__time__date__hour")
	odds.GameTime, exists = gtime.Attr("data-time")

	if !exists {
		return odds, errors.New("Could not find date time string for line")
	}

	// Find Home team odds
	hDom := e.DOM.Find(".game-line__home-line")
	vDom := e.DOM.Find(".game-line__visitor-line")

	odds.HomeTeam, exists = e.DOM.Find(".game-line__home-team__name").Attr("title")

	if !exists {
		return odds, errors.New("Could not parse home name")
	}

	odds.VisitingTeam, exists = e.DOM.Find(".game-line__visitor-team__name").Attr("title")

	if !exists {
		return odds, errors.New("Could not parse visitor name")
	}

	// Parse spread
	err = parseGameSpread(hDom, vDom, &odds)

	if err != nil {
		return odds, err
	}

	// Parse money line
	err = parseGameMoneyLine(hDom, vDom, &odds)

	if err != nil {
		return odds, err
	}

	err = parseGameOver(hDom, vDom, &odds)

	if err != nil {
		return odds, err
	}

	return odds, err
}
