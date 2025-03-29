package main

import (
	"fmt"
	"testing"
)

func TestHttpCrawlerStats(t *testing.T) {
	c := HttpCrawler{}
	players, err := c.Stats(1)
	if err != nil {
		t.Error(err)
	}

	if len(players) != 50 {
		t.Error("number of players on first page is expected to be 50!")
	}

	khanjFound := false
	for _, player := range players {
		if player.Name == "thekhanj" {
			khanjFound = true
		}
	}

	if khanjFound == false {
		// yeah i expect myself to always be on first page :)
		t.Error("thekhanj was expected to be on first page!")
	}
}

func TestHttpCrawlerOnline(t *testing.T) {
	t.SkipNow()

	c := HttpCrawler{}
	players, err := c.Online()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%d players are online\n", len(players))
}
