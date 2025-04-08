package core

import (
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
	return NewObserver(
		repo, &HttpCrawler{}, time.Minute*20, time.Minute,
	)
}

var CoreModule = wire.NewSet(
	db.DbModule,
	ProvideObserver, ProvidePlayerRepo,
)
