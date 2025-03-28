package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const CSDMPRO_SITE = "https://www.csdm.pro"

type Player struct {
	Name    string
	Country string
}

type Crawler struct{}

func (this *Crawler) Stats(page int) ([]Player, error) {
	resp, err := http.Get(CSDMPRO_SITE + fmt.Sprintf("/stats?p=%d", page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return this.parseBody(resp.Body)
}

func (this *Crawler) Online() ([]Player, error) {
	resp, err := http.Get(CSDMPRO_SITE)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return this.parseBody(resp.Body)
}

func (this *Crawler) parseBody(body io.ReadCloser) ([]Player, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	players := make([]Player, 0, 50)
	doc.Find(".stat tbody tr").Each(func(index int, row *goquery.Selection) {
		cols := row.Children().Eq(1)
		img := cols.Children().First()
		imgSrc, _ := img.Attr("src")
		username := strings.TrimSpace(cols.Text())

		player := Player{
			Name:    username,
			Country: imgSrc,
		}
		players = append(players, player)
	})

	return players, nil
}
