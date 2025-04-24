package core

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
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

	var rank any = player.Rank
	if player.Rank != nil {
		rank = *player.Rank
	}

	row, err := this.Database.Exec(
		insertSQL,
		player.Name, player.Country,
		rank, player.Score, player.Kills,
		player.Deaths, player.Accuracy,
	)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()

	return PlayerId(id), err
}

func (this *PlayerRepo) GetByRank(rank int) ([]DbPlayer, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM players as p
		WHERE rank = ?
		ORDER BY p.rank ASC
	`, this.getPlayerFields("p.")), rank)
	if err != nil {
		return []DbPlayer{}, err
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

func (this *PlayerRepo) UpdatePlayer(id PlayerId, player Player) error {
	insertSQL := `
		UPDATE players
		SET name = ?, country = ?, rank = ?, score = ?,
			kills = ?, deaths = ?, accuracy = ?
		WHERE id = ?
	`

	_, err := this.Database.Exec(
		insertSQL,
		player.Name, player.Country,
		player.Rank, player.Score, player.Kills,
		player.Deaths, player.Accuracy,
		id,
	)
	return err
}

func (this *PlayerRepo) Unrank(rank int) error {
	insertSQL := `
	UPDATE players
	SET rank = NULL
	WHERE rank = ?`

	_, err := this.Database.Exec(insertSQL, rank)

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
				FROM onlines AS o
				INNER JOIN players as player on player.id = o.player_id
				WHERE o.end_time IS NULL AND player.rank IS NOT NULL
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
	var rank sql.NullInt32

	err := rows.Scan(
		&p.ID,
		&p.Player.Name, &p.Player.Country,
		&rank, &p.Player.Score, &p.Player.Kills,
		&p.Player.Deaths, &p.Player.Accuracy,
	)

	if rank.Valid {
		var r int
		r = int(rank.Int32)
		p.Player.Rank = &r
	} else {
		p.Player.Rank = nil
	}

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

func (this *PlayerRepo) MarkOnline(playerId PlayerId) error {
	now := time.Now().Unix()
	insertSQL := `
		INSERT INTO onlines (player_id, start_time)
		VALUES (?, ?)
	`
	_, err := this.Database.Exec(insertSQL, playerId, now)
	return err
}

func (this *PlayerRepo) MarkOffline(playerId PlayerId) error {
	now := time.Now().Unix()
	updateSQL := `
		UPDATE onlines
		SET end_time = ?
		WHERE player_id = ? AND end_time IS NULL
	`
	_, err := this.Database.Exec(updateSQL, now, playerId, now)
	return err
}

func (this *PlayerRepo) List(offset int, limit int) ([]DbPlayer, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM players as p
		WHERE p.rank IS NOT NULL
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
		CREATE INDEX IF NOT EXISTS idx_players_rank ON
		players(rank)`
	_, err = db.Exec(createPlayersRankIndex)
	if err != nil {
		return nil, err
	}

	createPlayersNameRankIndex := `
		CREATE INDEX IF NOT EXISTS idx_players_name_rank
		ON players(name, rank)`
	_, err = db.Exec(createPlayersNameRankIndex)
	if err != nil {
		return nil, err
	}

	createOnlinesTable := `
		CREATE TABLE IF NOT EXISTS onlines (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			player_id INTEGER UNIQUE,
			start_time INTEGER NOT NULL,
			end_time INTEGER,
			FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE RESTRICT
		);`
	_, err = db.Exec(createOnlinesTable)
	if err != nil {
		return nil, err
	}

	createOnlinesStartTimeIndex := `
		CREATE INDEX IF NOT EXISTS idx_onlines_start_time
		ON onlines(player_id, start_time, end_time)
	`
	_, err = db.Exec(createOnlinesStartTimeIndex)
	if err != nil {
		return nil, err
	}

	createOnlinesEndTimeIndex := `
		CREATE INDEX IF NOT EXISTS idx_onlines_end_time
		ON onlines(player_id, end_time)
	`
	_, err = db.Exec(createOnlinesEndTimeIndex)
	if err != nil {
		return nil, err
	}

	return &PlayerRepo{Database: db}, nil
}
