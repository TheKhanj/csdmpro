package core

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type TestingRepoFactory struct {
	Repo *PlayerRepo

	tmpFilePath string
	db          *sql.DB
}

func (this *TestingRepoFactory) Deinit() error {
	if this.tmpFilePath != "" {
		err := os.Remove(this.tmpFilePath)
		if err != nil {
			return err
		}
	}
	if this.db != nil {
		err := this.db.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *TestingRepoFactory) Init() error {
	tempFile, err := os.CreateTemp("", "csdmpro-test-*.db")
	if err != nil {
		return err
	}
	tempFile.Close()

	this.tmpFilePath = tempFile.Name()
	this.db, err = sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return err
	}

	f := &PlayerRepoFactory{
		Database: this.db,
	}

	err = f.assertTables()
	if err != nil {
		return err
	}

	repo, err := f.Create()
	if err != nil {
		return err
	}

	this.Repo = repo
	return nil
}

func TestPlayerRepoPlayer(t *testing.T) {
	trf := TestingRepoFactory{}
	err := trf.Init()
	if err != nil {
		t.Error(err)
	}
	defer trf.Deinit()

	repo := trf.Repo

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
	trf := TestingRepoFactory{}
	err := trf.Init()
	if err != nil {
		t.Error(err)
	}
	defer trf.Deinit()

	repo := trf.Repo

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
