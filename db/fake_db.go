package db

import (
	"database/sql"
	"os"
)

type FakeDbFactory struct {
	tmpFilePath string
	db          *sql.DB
}

func (this *FakeDbFactory) Init() (*sql.DB, error) {
	tempFile, err := os.CreateTemp("", "csdmpro-test-*.db")
	if err != nil {
		return nil, err
	}
	tempFile.Close()

	this.tmpFilePath = tempFile.Name()
	this.db, err = sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return nil, err
	}

	return this.db, nil
}

func (this *FakeDbFactory) Deinit() error {
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
