package core

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const CSDMPRO_SITE = "https://www.csdm.pro"

type Player struct {
	Name     string
	Country  string
	Rank     int
	Score    int
	Kills    int
	Deaths   int
	Accuracy int
}

type Crawler interface {
	Stats(int) ([]Player, error)
	Online() ([]Player, error)
}

type HttpCrawler struct{}

func (this *HttpCrawler) Stats(page int) ([]Player, error) {
	resp, err := http.Get(CSDMPRO_SITE + fmt.Sprintf("/stats?p=%d", page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return this.parseBody(resp.Body)
}

func (this *HttpCrawler) Online() ([]Player, error) {
	resp, err := http.Get(CSDMPRO_SITE)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return this.parseBody(resp.Body)
}

func (this *HttpCrawler) parseBody(body io.ReadCloser) ([]Player, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	players := make([]Player, 0, 50)
	getIntCol := func(row *goquery.Selection, col int) (int, error) {
		str := strings.TrimSpace(row.Children().Eq(col).Text())

		return strconv.Atoi(str)
	}

	doc.Find(".stat tbody tr").Each(func(index int, row *goquery.Selection) {
		cols := row.Children().Eq(1)
		img := cols.Children().First()
		imgSrc, _ := img.Attr("src")
		username := strings.TrimSpace(cols.Text())

		rank, err := getIntCol(row, 0)
		score, err := getIntCol(row, 2)
		kills, err := getIntCol(row, 3)
		deaths, err := getIntCol(row, 4)
		accuracyStr := strings.TrimSpace(row.Children().Eq(6).Text())
		re := regexp.MustCompile("([0-9]*)")
		accuracyStrCleanedMatches := re.FindStringSubmatch(accuracyStr)
		if len(accuracyStrCleanedMatches) == 0 {
			log.Println(errors.New("accuracy regex not matched!"))
			return
		}
		accuracyStrCleaned := accuracyStrCleanedMatches[0]
		accuracy, err := strconv.Atoi(accuracyStrCleaned)
		if err != nil {
			log.Println(err)
			return
		}

		player := Player{
			Name:     username,
			Country:  imgSrc,
			Rank:     rank,
			Score:    score,
			Kills:    kills,
			Deaths:   deaths,
			Accuracy: accuracy,
		}
		players = append(players, player)
	})

	return players, nil
}

var _ Crawler = (*HttpCrawler)(nil)

type StubCrawler struct {
	Onlines []Player
	Players []Player
}

func (this *StubCrawler) Stats(page int) ([]Player, error) {
	l := (page - 1) * 50
	r := page * 50

	if r > len(this.Players) {
		return this.Players[:], nil
	}

	ret := make([]Player, 0, r-l)
	for _, p := range this.Players[l:r] {
		player := Player{}
		player = p
		ret = append(ret, player)
	}

	return ret, nil
}

func (this *StubCrawler) Online() ([]Player, error) {
	ret := make([]Player, 0, len(this.Onlines))
	for _, p := range this.Onlines {
		player := Player{}
		player = p
		ret = append(ret, player)
	}

	return ret, nil
}
