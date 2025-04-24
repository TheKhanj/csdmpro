package db

import (
	"database/sql"
	"log"

	"github.com/google/wire"
	_ "github.com/mattn/go-sqlite3"
)

type Database *sql.DB

func ProvideDb() Database {
	db, err := OpenDb("database.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var DbModule = wire.NewSet(ProvideDb)
