package core

import (
	"database/sql"
	"errors"
	"fmt"
)

type PlayerId int

type DbPlayer struct {
	ID     PlayerId
	Player Player
}

type PlayerRepo struct {
	Database *sql.DB
}

var ERR_PLAYER_NOT_FOUND error = errors.New("player not found")

func (this *PlayerRepo) AddPlayer(player Player) (PlayerId, error) {
	insertSQL := `
	INSERT INTO players (name, country, rank, score, kills, deaths, accuracy)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	row, err := this.Database.Exec(
		insertSQL,
		player.Name, player.Country,
		player.Rank, player.Score, player.Kills,
		player.Deaths, player.Accuracy,
	)
	id, err := row.LastInsertId()

	return PlayerId(id), err
}

func (this *PlayerRepo) UpdatePlayer(id PlayerId, player Player) error {
	insertSQL := `
	UPDATE players
	SET name = ?, country = ?, rank = ?, score = ?,
		kills = ?, deaths = ?, accuracy = ?
	WHERE id = ?`

	_, err := this.Database.Exec(
		insertSQL,
		player.Name, player.Country,
		player.Rank, player.Score, player.Kills,
		player.Deaths, player.Accuracy,
		id,
	)
	return err
}

func (this *PlayerRepo) IsOnline(name string) (bool, error) {
	onlines, err := this.Onlines()
	if err != nil {
		return false, err
	}

	for _, online := range onlines {
		if online.Player.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func (this *PlayerRepo) Onlines() ([]DbPlayer, error) {
	rows, err := this.Database.Query(
		fmt.Sprintf(
			`SELECT %s
				FROM players_online as online
				INNER JOIN players as player on player.id = online.player_id
				ORDER BY player.rank ASC
			`,
			this.getPlayerFields("player."),
		))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]DbPlayer, 0, 0)

	for rows.Next() {
		p, err := this.scanPlayer(rows)
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	return players, nil
}

func (this *PlayerRepo) scanPlayer(rows *sql.Rows) (DbPlayer, error) {
	var p DbPlayer
	var id PlayerId

	err := rows.Scan(
		&id,
		&p.Player.Name, &p.Player.Country,
		&p.Player.Rank, &p.Player.Score, &p.Player.Kills,
		&p.Player.Deaths, &p.Player.Accuracy,
	)
	p.ID = id

	return p, err
}

func (this *PlayerRepo) getPlayerFields(prefix string) string {
	ret := ""
	ret += fmt.Sprintf("%s%s, ", prefix, "id")
	ret += fmt.Sprintf("%s%s, ", prefix, "name")
	ret += fmt.Sprintf("%s%s, ", prefix, "country")
	ret += fmt.Sprintf("%s%s, ", prefix, "rank")
	ret += fmt.Sprintf("%s%s, ", prefix, "score")
	ret += fmt.Sprintf("%s%s, ", prefix, "kills")
	ret += fmt.Sprintf("%s%s, ", prefix, "deaths")
	ret += fmt.Sprintf("%s%s ", prefix, "accuracy")
	return ret
}

func (this *PlayerRepo) GetPlayerByName(name string) (DbPlayer, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM players as p
		WHERE name = ?
		ORDER BY p.rank ASC
		LIMIT 1
	`, this.getPlayerFields("p.")), name)
	if err != nil {
		return DbPlayer{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return DbPlayer{}, ERR_PLAYER_NOT_FOUND
	}

	return this.scanPlayer(rows)
}

func (this *PlayerRepo) GetPlayer(id PlayerId) (DbPlayer, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM players as p
		WHERE id = ?
		LIMIT 1
	`, this.getPlayerFields("p.")), id)
	if err != nil {
		return DbPlayer{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return DbPlayer{}, ERR_PLAYER_NOT_FOUND
	}

	return this.scanPlayer(rows)
}

func (this *PlayerRepo) AddOnlinePlayer(playerId PlayerId) error {
	insertSQL := `INSERT INTO players_online (player_id) VALUES (?)`
	_, err := this.Database.Exec(insertSQL, playerId)
	return err
}

func (this *PlayerRepo) RemoveOnlinePlayer(playerId PlayerId) error {
	insertSQL := `DELETE FROM players_online WHERE player_id = ?`
	_, err := this.Database.Exec(insertSQL, playerId)
	return err
}

func (this *PlayerRepo) List(offset int, limit int) ([]DbPlayer, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM players as p
		ORDER BY p.rank ASC
		LIMIT ? OFFSET ?
	`, this.getPlayerFields("p.")), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]DbPlayer, 0, 0)

	for rows.Next() {
		p, err := this.scanPlayer(rows)
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	return players, nil
}

func CreatePlayerRepo(db *sql.DB) (*PlayerRepo, error) {
	createPlayersStatsTable := `CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		country TEXT,
		rank INTEGER,
		score INTEGER,
		kills INTEGER,
		deaths INTEGER,
		accuracy INTEGER
	);`
	_, err := db.Exec(createPlayersStatsTable)
	if err != nil {
		return nil, err
	}

	createPlayersRankIndex := `
		CREATE INDEX IF NOT EXISTS idx_rank ON
		players(rank)`

	_, err = db.Exec(createPlayersRankIndex)
	if err != nil {
		return nil, err
	}

	createPlayersNameRankIndex := `
		CREATE INDEX IF NOT EXISTS idx_name_rank
		ON players(name, rank)`

	_, err = db.Exec(createPlayersNameRankIndex)
	if err != nil {
		return nil, err
	}

	createPlayersOnlineTable := `
		CREATE TABLE IF NOT EXISTS players_online (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			player_id INTEGER UNIQUE,
			FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE RESTRICT
		);`
	_, err = db.Exec(createPlayersOnlineTable)
	if err != nil {
		return nil, err
	}

	return &PlayerRepo{Database: db}, nil
}
