package main

import (
	"context"
	"testing"
	"time"
)

type TestingObserverFactory struct {
	Crawler  *StubCrawler
	Observer *Observer

	trf TestingRepoFactory
}

func (this *TestingObserverFactory) Init(t *testing.T) {
	err := this.trf.Init()
	if err != nil {
		t.Error(err)
	}

	repo := this.trf.Repo

	this.Crawler = &StubCrawler{
		Onlines: make([]Player, 0),
		Players: make([]Player, 0),
	}

	this.Observer = &Observer{
		Repo:           repo,
		Crawler:        this.Crawler,
		GotOnline:      make(chan Player),
		GotOffline:     make(chan Player),
		StatsInterval:  0,
		OnlineInterval: 0,
	}
}

func (this *TestingObserverFactory) Deinit() {
	defer this.trf.Deinit()
}

func TestObserverSimply(t *testing.T) {
	tof := TestingObserverFactory{}
	tof.Init(t)
	defer tof.Deinit()

	go tof.Observer.Start()
	defer tof.Observer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c := tof.Crawler

	player := Player{
		Name:    "thekhanj",
		Country: "Iran 😭, fuck IRAN, ISLAMIC REPUBLIC to be more accurate. Iran is lovely❤️",
	}
	c.Onlines = append(c.Onlines, player)

	select {
	case <-ctx.Done():
		t.Error("expected online player event to pass in")
	case p := <-tof.Observer.GotOnline:
		if p.Name != player.Name {
			t.Error("player name does not match")
		}
	}

	c.Onlines = c.Onlines[:len(c.Onlines)-1]

	select {
	case <-ctx.Done():
		t.Error("expected offline player event to pass in")
	case p := <-tof.Observer.GotOffline:
		if p.Name != player.Name {
			t.Error("player name does not match")
		}
	}
}
