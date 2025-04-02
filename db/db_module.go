package db

import (
	"database/sql"
	"log"

	"github.com/google/wire"
)

type Database *sql.DB

func ProvideDb() Database {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var DbModule = wire.NewSet(ProvideDb)
