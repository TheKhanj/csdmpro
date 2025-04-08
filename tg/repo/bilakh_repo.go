package repo

import "database/sql"

type BilakhRepo struct {
	db *sql.DB
}

func (this *BilakhRepo) IsBilakhed(chatId int64) (bool, error) {
	sql := `SELECT chat_id FROM bilakhs WHERE chat_id = ?`

	rows, err := this.db.Query(sql, chatId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func CreateBilakhRepo(db *sql.DB) (*BilakhRepo, error) {
	sql := `CREATE TABLE IF NOT EXISTS bilakhs (
		chat_id INTEGER PRIMARY KEY
	);`

	_, err := db.Exec(sql)
	if err != nil {
		return nil, err
	}

	return &BilakhRepo{db}, nil
}
