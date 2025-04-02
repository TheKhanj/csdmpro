package core

import (
	"testing"

	"github.com/thekhanj/csdmpro/db"
)

func TestPlayerRepoPlayer(t *testing.T) {
	f := db.FakeDbFactory{}
	db, err := f.Init()
	if err != nil {
		t.Error(err)
	}
	defer f.Deinit()

	repo, err := CreatePlayerRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.AddPlayer(Player{
		Name:    "thekhanj",
		Country: "iran",
	})
	if err != nil {
		t.Error(err)
	}

	_, err = repo.GetPlayerId("thekhanj")
	if err != nil {
		t.Error(err)
	}

	exists, err := repo.PlayerExists("thekhanj")
	if err != nil {
		t.Error(err)
	}
	if exists == false {
		t.Error("player thekhanj must exist in database")
	}
}

func TestPlayerRepoOnline(t *testing.T) {
	f := db.FakeDbFactory{}
	db, err := f.Init()
	if err != nil {
		t.Error(err)
	}
	defer f.Deinit()

	repo, err := CreatePlayerRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.AddPlayer(Player{
		Name:    "thekhanj",
		Country: "iran",
	})
	if err != nil {
		t.Error(err)
	}

	isOnline, err := repo.IsOnline("thekhanj")
	if err != nil {
		t.Error(err)
	}
	if isOnline {
		t.Error("player should not be online before adding to database")
	}

	id, err := repo.GetPlayerId("thekhanj")
	if err != nil {
		t.Error(err)
	}

	err = repo.AddOnlinePlayer(id)
	if err != nil {
		t.Error(err)
	}

	isOnline, err = repo.IsOnline("thekhanj")
	if err != nil {
		t.Error(err)
	}
	if !isOnline {
		t.Error("player should be online after adding to database")
	}

	onlines, err := repo.Onlines()
	if err != nil {
		t.Error(err)
	}
	if len(onlines) != 1 {
		t.Error("number of online players must be 1")
	}

	err = repo.RemoveOnlinePlayer(id)
	if err != nil {
		t.Error(err)
	}

	isOnline, err = repo.IsOnline("thekhanj")
	if err != nil {
		t.Error(err)
	}
	if isOnline {
		t.Error("player should not be online after removing it from database")
	}
}
