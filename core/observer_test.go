package core

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/thekhanj/csdmpro/db"
)

type TestingObserverFactory struct {
	Crawler  *StubCrawler
	Observer *Observer

	dbF db.FakeDbFactory
}

func (this *TestingObserverFactory) Init(t *testing.T) {
	db, err := this.dbF.Init()
	if err != nil {
		t.Fatal(err)
	}

	repo, err := CreatePlayerRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	this.Crawler = &StubCrawler{
		Onlines: make([]Player, 0),
		Players: make([]Player, 0),
	}

	this.Observer = NewObserver(repo, this.Crawler, 0, 0)
}

func (this *TestingObserverFactory) Deinit() {
	this.dbF.Deinit()
}

func TestObserverSimply(t *testing.T) {
	tof := TestingObserverFactory{}
	tof.Init(t)
	defer tof.Deinit()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
	defer cancel()

	go tof.Observer.Start(ctx)
	go func() {
		for {
			select {
			case <-tof.Observer.UpdatedPlayer:
			case <-tof.Observer.AddedPlayer:
				// throw away channels
			}
		}
	}()

	c := tof.Crawler

	player := Player{
		Name:    "thekhanj",
		Country: "Iran 😭, fuck IRAN, ISLAMIC REPUBLIC to be more accurate. Iran is lovely❤️",
	}
	c.Onlines = append(c.Onlines, player)

	select {
	case <-ctx.Done():
		t.Fatal("expected online player event to pass in")
	case p := <-tof.Observer.GotOnline:
		if p.Player.Name != player.Name {
			t.Fatal("player name does not match")
		}
	}

	c.Onlines = c.Onlines[:len(c.Onlines)-1]

	select {
	case <-ctx.Done():
		t.Fatal("expected offline player event to pass in")
	case p := <-tof.Observer.GotOffline:
		if p.Player.Name != player.Name {
			t.Fatal("player name does not match")
		}
	}
}

func TestObserverStats(t *testing.T) {
	tof := TestingObserverFactory{}
	tof.Init(t)
	defer tof.Deinit()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	go tof.Observer.Start(ctx)
	go func() {
		for {
			select {
			case <-tof.Observer.UpdatedPlayer:
			case <-tof.Observer.AddedPlayer:
				// throw away channels
			}
		}
	}()

	c := tof.Crawler

	oneRank := 1
	player := Player{
		Name:    "thekhanj",
		Country: "Iran 😭, fuck IRAN, ISLAMIC REPUBLIC to be more accurate. Iran is lovely❤️",
		Rank:    &oneRank,
	}
	c.Players = append(c.Players, player)

	select {
	case <-ctx.Done():
		t.Fatal("expected new player event to pass in")
	case id := <-tof.Observer.AddedPlayer:
		t.Logf("new player with id %d created", id)
	}

	twoRank := 2
	c.Players[0].Rank = &twoRank

	select {
	case <-ctx.Done():
		t.Fatal("expected update player event to pass in")
	case id := <-tof.Observer.UpdatedPlayer:
		t.Logf("player with id %d updated", id)
	}
}

func TestObserverMultipleEvents(t *testing.T) {
	tof := TestingObserverFactory{}
	tof.Init(t)
	defer tof.Deinit()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	go tof.Observer.Start(ctx)
	go func() {
		for {
			select {
			case <-tof.Observer.UpdatedPlayer:
			case <-tof.Observer.AddedPlayer:
				// throw away channels
			}
		}
	}()

	c := tof.Crawler

	players := make([]Player, 0, 100)
	for i := 0; i < 100; i++ {
		players = append(players, Player{
			Name:    fmt.Sprintf("player-%d", i),
			Country: "anything",
		})
	}
	player_index := 0
	for i := 0; i < 10; i++ {
		p := players[player_index]
		player_index++

		c.Onlines = append(c.Onlines, p)

		event := <-tof.Observer.GotOnline
		if event.Player.Name != p.Name {
			t.Fatalf("unexpected player name")
		}
	}

	for i := 0; i < 5; i++ {
		p := players[player_index-1-i]

		c.Onlines = c.Onlines[0 : len(c.Onlines)-1]

		event := <-tof.Observer.GotOffline
		if event.Player.Name != p.Name {
			t.Fatalf("unexpected player name")
		}
	}

	for i := 0; i < 10; i++ {
		p := players[player_index]
		player_index++

		c.Onlines = append(c.Onlines, p)

		event := <-tof.Observer.GotOnline
		if event.Player.Name != p.Name {
			t.Fatalf("unexpected player name")
		}
	}

	player_index = 5
	for i := 0; i < 5; i++ {
		p := players[player_index]
		player_index++

		c.Onlines = append(c.Onlines, p)

		event := <-tof.Observer.GotOnline
		if event.Player.Name != p.Name {
			t.Fatalf("unexpected player name")
		}
	}
}
