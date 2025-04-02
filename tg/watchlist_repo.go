package tg

import (
	"database/sql"

	"github.com/thekhanj/csdmpro/core"
)

type WatchlistRepo struct {
	db *sql.DB
}

func (this *WatchlistRepo) Watchlist(chatId int64) ([]core.PlayerId, error) {
	rows, err := this.db.Query(`
	SELECT w.player_id
	FROM watchlist as w
	WHERE w.chat_id = ?
	`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playersIds := make([]core.PlayerId, 0, 0)

	for rows.Next() {
		var playerId core.PlayerId
		err = rows.Scan(&playerId)
		if err != nil {
			return nil, err
		}

		playersIds = append(playersIds, playerId)
	}

	return playersIds, nil
}

func (this *WatchlistRepo) IsInWatchlist(
	chatId int64, playerId core.PlayerId,
) (bool, error) {
	ids, err := this.Watchlist(chatId)
	if err != nil {
		return false, err
	}

	for _, id := range ids {
		if id == playerId {
			return true, nil
		}
	}

	return false, nil
}

func (this *WatchlistRepo) AddToWatchlist(
	chatId int64, playerId core.PlayerId,
) error {
	insertSQL := `INSERT INTO watchlist (chat_id, player_id) VALUES (?, ?)`
	_, err := this.db.Exec(insertSQL, chatId, playerId)
	return err
}

func (this *WatchlistRepo) RemoveFromWatchlist(
	chatId int64, playerId core.PlayerId,
) error {
	insertSQL := `DELETE FROM watchlist WHERE chat_id = ? AND player_id = ?`
	_, err := this.db.Exec(insertSQL, chatId, playerId)
	return err
}

func CreateWatchlistRepo(db *sql.DB) (*WatchlistRepo, error) {
	createPlayersOnlineTable := `CREATE TABLE IF NOT EXISTS watchlist (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER,
		player_id INTEGER,
		FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE RESTRICT,
		UNIQUE(chat_id, player_id)
	);`
	_, err := db.Exec(createPlayersOnlineTable)
	if err != nil {
		return nil, err
	}

	return &WatchlistRepo{db}, nil
}
