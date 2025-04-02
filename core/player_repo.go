package core

import (
	"database/sql"
)

type PlayerId int

type PlayerRepo struct {
	Database *sql.DB
}

func (this *PlayerRepo) AddPlayer(player Player) error {
	insertSQL := `INSERT INTO players (name, country) VALUES (?, ?)`
	_, err := this.Database.Exec(insertSQL, player.Name, player.Country)
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
	rows, err := this.Database.Query(`
		SELECT player.name, player.country
		FROM players_online as online
		INNER JOIN players as player on player.id = online.player_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]Player, 0, 0)

	for rows.Next() {
		var name string
		var country string
		err = rows.Scan(&name, &country)
		if err != nil {
			return nil, err
		}

		players = append(players, Player{
			Name:    name,
			Country: country,
		})
	}

	return players, nil
}

func (this *PlayerRepo) GetPlayerId(name string) (int, error) {
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

	var id int
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (this *PlayerRepo) AddOnlinePlayer(playerId int) error {
	insertSQL := `INSERT INTO players_online (player_id) VALUES (?)`
	_, err := this.Database.Exec(insertSQL, playerId)
	return err
}

func (this *PlayerRepo) RemoveOnlinePlayer(playerId int) error {
	insertSQL := `DELETE FROM players_online WHERE player_id = ?`
	_, err := this.Database.Exec(insertSQL, playerId)
	return err
}

func (this *PlayerRepo) List(offset int32, limit int32) ([]Player, error) {
	rows, err := this.Database.Query(`
		SELECT p.name, p.country
		FROM players as p
		ORDER BY p.id
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]Player, 0, 0)

	for rows.Next() {
		var name string
		var country string
		err = rows.Scan(&name, &country)
		if err != nil {
			return nil, err
		}

		players = append(players, Player{
			Name:    name,
			Country: country,
		})
	}

	return players, nil
}

func CreatePlayerRepo(db *sql.DB) (*PlayerRepo, error) {
	createPlayersStatsTable := `CREATE TABLE IF NOT EXISTS players (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		country TEXT
	);`
	_, err := db.Exec(createPlayersStatsTable)
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
