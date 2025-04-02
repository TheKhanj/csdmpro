package tg

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/db"
)

func addCoupleOfPlayers(db *sql.DB) error {
	repo, err := core.CreatePlayerRepo(db)
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		err := repo.AddPlayer(core.Player{
			Name:    fmt.Sprintf("player-%d", i),
			Country: "does not matter",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestWatchlistRepoSimply(t *testing.T) {
	f := db.FakeDbFactory{}
	db, err := f.Init()
	if err != nil {
		t.Error(err)
	}
	defer f.Deinit()

	err = addCoupleOfPlayers(db)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	repo, err := CreateWatchlistRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	var chatId int64 = 1
	ids, err := repo.Watchlist(chatId)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Fatal("expected watchlist to initially be empty")
	}

	err = repo.AddToWatchlist(chatId, 1)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.AddToWatchlist(chatId, 2)
	if err != nil {
		t.Fatal(err)
	}

	ids, err = repo.Watchlist(chatId)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 {
		t.Fatal("expected to have 2 people on watchlist")
	}

	isWatched, err := repo.IsInWatchlist(chatId, 2)
	if err != nil {
		t.Fatal(err)
	}
	if isWatched == false {
		t.Fatal("expected player-2 to be watched")
	}

	isWatched, err = repo.IsInWatchlist(chatId, 3)
	if err != nil {
		t.Fatal(err)
	}
	if isWatched == true {
		t.Fatal("expected player-3 not to be watched")
	}

	err = repo.RemoveFromWatchlist(chatId, 2)
	if err != nil {
		t.Fatal(err)
	}
}
