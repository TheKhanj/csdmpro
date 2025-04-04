package core

import (
	"testing"

	"github.com/thekhanj/csdmpro/db"
)

func TestPlayerRepoPlayer(t *testing.T) {
	f := db.FakeDbFactory{}
	db, err := f.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Deinit()

	repo, err := CreatePlayerRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.AddPlayer(Player{
		Name:    "thekhanj",
		Country: "iran",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetPlayerByName("thekhanj")
	if err == ERR_PLAYER_NOT_FOUND {
		t.Fatal("player thekhanj must exist in database")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlayerRepoOnline(t *testing.T) {
	f := db.FakeDbFactory{}
	db, err := f.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Deinit()

	repo, err := CreatePlayerRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.AddPlayer(Player{
		Name:    "thekhanj",
		Country: "iran",
	})
	if err != nil {
		t.Fatal(err)
	}

	isOnline, err := repo.IsOnline("thekhanj")
	if err != nil {
		t.Fatal(err)
	}
	if isOnline {
		t.Fatal("player should not be online before adding to database")
	}

	p, err := repo.GetPlayerByName("thekhanj")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.AddOnlinePlayer(p.ID)
	if err != nil {
		t.Fatal(err)
	}

	isOnline, err = repo.IsOnline("thekhanj")
	if err != nil {
		t.Fatal(err)
	}
	if !isOnline {
		t.Fatal("player should be online after adding to database")
	}

	onlines, err := repo.Onlines()
	if err != nil {
		t.Fatal(err)
	}
	if len(onlines) != 1 {
		t.Fatal("number of online players must be 1")
	}

	err = repo.RemoveOnlinePlayer(p.ID)
	if err != nil {
		t.Fatal(err)
	}

	isOnline, err = repo.IsOnline("thekhanj")
	if err != nil {
		t.Fatal(err)
	}
	if isOnline {
		t.Fatal("player should not be online after removing it from database")
	}
}
