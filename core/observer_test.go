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
	Repo     *PlayerRepo

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

	this.Repo = repo
	this.Crawler = NewStubCrawler()
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

	name := tof.Crawler.AddPlayer()
	tof.Crawler.MakeOnline(name)

	gotOnline := tof.Observer.Bus.Sub(GotOnlineTopic)
	gotOffline := tof.Observer.Bus.Sub(GotOfflineTopic)
	defer func() {
		go tof.Observer.Bus.Unsub(gotOnline)
		go tof.Observer.Bus.Unsub(gotOffline)
		for {
			select {
			case <-gotOnline:
			case <-gotOffline:
			default:
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		t.Fatal("expected online player event to pass in")
	case pId := <-gotOnline:
		if pId != 1 {
			t.Fatal("player id is not 1")
		}
		p, _ := tof.Repo.GetPlayer(pId)
		if *p.Player.Rank != 1 {
			t.Fatal(
				fmt.Sprintf(
					"expected first player to have rank 1 got %d", *p.Player.Rank,
				),
			)
		}
	}

	tof.Crawler.MakeOffline(name)

	select {
	case <-ctx.Done():
		t.Fatal("expected offline player event to pass in")
	case pId := <-gotOffline:
		if pId != 1 {
			t.Fatal("player id is not 1")
		}
	}
}

func TestObserverMultipleEvents(t *testing.T) {
	tof := TestingObserverFactory{}
	tof.Init(t)
	defer tof.Deinit()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	go tof.Observer.Start(ctx)

	gotOnline := tof.Observer.Bus.Sub(GotOnlineTopic)
	gotOffline := tof.Observer.Bus.Sub(GotOfflineTopic)
	defer func() {
		go tof.Observer.Bus.Unsub(gotOnline)
		go tof.Observer.Bus.Unsub(gotOffline)
		for {
			select {
			case <-gotOnline:
			case <-gotOffline:
			default:
				return
			}
		}
	}()

	for i := 0; i < 50; i++ {
		tof.Crawler.AddPlayer()
	}
	players, _ := tof.Crawler.Stats(1)

	for i := 0; i < 10; i++ {
		p := players[i]

		tof.Crawler.MakeOnline(p.Name)

		eventPlayerId := <-gotOnline
		eventPlayer, err := tof.Repo.GetPlayer(eventPlayerId)
		if err != nil {
			t.Fatal(err)
		}
		if eventPlayer.Player.Name != p.Name {
			t.Fatalf("unexpected player name (%s, %s)", eventPlayer.Player.Name, p.Name)
		}
	}

	onlines, _ := tof.Crawler.Online()
	for i := 0; i < 5; i++ {
		p := onlines[i]

		tof.Crawler.MakeOffline(p.Name)

		eventPlayerId := <-gotOffline
		eventPlayer, err := tof.Repo.GetPlayer(eventPlayerId)
		if err != nil {
			t.Fatal(err)
		}
		if eventPlayer.Player.Name != p.Name {
			t.Fatalf("unexpected player name")
		}
	}
}
