package core

import (
	"database/sql"
	"errors"
	"fmt"
)

type PlayerId int

type PlayerRepo struct {
	Database *sql.DB
}

func (this *PlayerRepo) AddPlayer(player Player) error {
	insertSQL := `
	INSERT INTO players (name, country, rank, score, kills, deaths, accuracy)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := this.Database.Exec(
		insertSQL,
		player.Name, player.Country,
		player.Rank, player.Score, player.Kills,
		player.Deaths, player.Accuracy,
	)
	return err
}

func (this *PlayerRepo) PlayerExists(name string) (bool, error) {
	rows, err := this.Database.Query(`
		SELECT name
		FROM players
		WHERE name = ?
		LIMIT 1
	`, name)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	cnt := 0
	for rows.Next() {
		cnt++
	}

	return cnt == 1, nil
}

func (this *PlayerRepo) IsOnline(name string) (bool, error) {
	onlines, err := this.Onlines()
	if err != nil {
		return false, err
	}

	for _, online := range onlines {
		if online.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func (this *PlayerRepo) Onlines() ([]Player, error) {
	rows, err := this.Database.Query(
		fmt.Sprintf(
			`SELECT %s
		FROM players_online as online
		INNER JOIN players as player on player.id = online.player_id
		ORDER BY rank
	`, this.getPlayerFields("player.")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]Player, 0, 0)

	for rows.Next() {
		p, err := this.scanPlayer(rows)
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	return players, nil
}

func (this *PlayerRepo) scanPlayer(rows *sql.Rows) (Player, error) {
	var p Player
	var id int
	p.ID = &id

	err := rows.Scan(
		&id,
		&p.Name, &p.Country,
		&p.Rank, &p.Score, &p.Kills,
		&p.Deaths, &p.Accuracy,
	)

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

// Deprecated: don't use this
func (this *PlayerRepo) GetPlayerId(name string) (PlayerId, error) {
	rows, err := this.Database.Query(`
		SELECT id
		FROM players
		WHERE name = ?
		LIMIT 1
	`, name)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rows.Next()

	var id PlayerId
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (this *PlayerRepo) GetPlayer(id PlayerId) (Player, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM players as p
		WHERE id = ?
		LIMIT 1
	`, this.getPlayerFields("p.")), id)
	if err != nil {
		return Player{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return Player{}, errors.New("player not found")
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

func (this *PlayerRepo) List(offset int, limit int) ([]Player, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM players as p
		ORDER BY p.id
		LIMIT ? OFFSET ?
	`, this.getPlayerFields("p.")), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]Player, 0, 0)

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

	createPlayersRankIndex := `CREATE INDEX IF NOT EXISTS idx_rank ON players(rank)`

	_, err = db.Exec(createPlayersRankIndex)
	if err != nil {
		return nil, err
	}

	createPlayersOnlineTable := `CREATE TABLE IF NOT EXISTS players_online (
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
