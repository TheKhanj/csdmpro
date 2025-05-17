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

func ProvideSnapshotRepo(db db.Database) *SnapshotRepo {
	repo, err := CreateSnapshotRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func ProvideObserver(
	playerRepo *PlayerRepo,
	snapshotRepo *SnapshotRepo,
) *Observer {
	return NewObserver(
		playerRepo, snapshotRepo,
		&HttpCrawler{}, time.Second*30, time.Minute,
	)
}

var CoreModule = wire.NewSet(
	db.DbModule,
	ProvideObserver, ProvidePlayerRepo, ProvideSnapshotRepo,
)
