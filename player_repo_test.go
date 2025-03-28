package main

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func getRepoFactory(ctx context.Context) (*PlayerRepoFactory, error) {
	tempFile, err := os.CreateTemp("", "csdmpro-test-*.db")
	if err != nil {
		return nil, err
	}
	tempFile.Close()

	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			db.Close()
		}
	}()

	repoFactory := &PlayerRepoFactory{
		Database: db,
	}

	return repoFactory, nil
}

func getRepo(ctx context.Context) (*PlayerRepo, error) {
	f, err := getRepoFactory(ctx)
	if err != nil {
		return nil, err
	}

	err = f.assertTables()
	if err != nil {
		return nil, err
	}

	return f.Create()
}

func TestPlayerRepoPlayer(t *testing.T) {
	repo, err := getRepo(t.Context())
	if err != nil {
		t.Error(err)
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
	repo, err := getRepo(t.Context())
	if err != nil {
		t.Error(err)
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
