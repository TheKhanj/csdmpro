package core

import (
	"context"
	"log"
	"time"

	"github.com/google/wire"
	"github.com/thekhanj/csdmpro/db"
)

func ProvidePlayerRepo(db db.Database) *PlayerRepo {
	repo, err := CreatePlayerRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func ProvideObserver(repo *PlayerRepo) *Observer {
	return &Observer{
		Repo:           repo,
		Crawler:        &HttpCrawler{},
		GotOnline:      make(chan Player, 0),
		GotOffline:     make(chan Player, 0),
		StatsInterval:  time.Minute * 20,
		OnlineInterval: time.Minute,
		Ctx:            context.Background(),
	}
}

var CoreModule = wire.NewSet(
	db.DbModule,
	ProvideObserver, ProvidePlayerRepo,
)
