package core

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type SnapshotPlayer struct {
	ID     PlayerId
	Player Player
	Time   int64
}

type SnapshotRepo struct {
	Database *sql.DB
}

func (this *SnapshotRepo) Get(name string) (SnapshotPlayer, error) {
	rows, err := this.Database.Query(fmt.Sprintf(`
		SELECT %s
		FROM snapshot as s
		WHERE name = ?
		ORDER BY s.rank ASC
		LIMIT 1
	`, this.getPlayerFields("s.")), name)
	if err != nil {
		return SnapshotPlayer{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return SnapshotPlayer{}, ERR_PLAYER_NOT_FOUND
	}

	return this.scanPlayer(rows)
}

func (this *SnapshotRepo) getPlayerFields(prefix string) string {
	ret := ""
	ret += fmt.Sprintf("%s%s, ", prefix, "id")
	ret += fmt.Sprintf("%s%s, ", prefix, "time")
	ret += fmt.Sprintf("%s%s, ", prefix, "name")
	ret += fmt.Sprintf("%s%s, ", prefix, "country")
	ret += fmt.Sprintf("%s%s, ", prefix, "rank")
	ret += fmt.Sprintf("%s%s, ", prefix, "score")
	ret += fmt.Sprintf("%s%s, ", prefix, "kills")
	ret += fmt.Sprintf("%s%s, ", prefix, "deaths")
	ret += fmt.Sprintf("%s%s ", prefix, "accuracy")
	return ret
}

func (this *SnapshotRepo) scanPlayer(rows *sql.Rows) (SnapshotPlayer, error) {
	var p SnapshotPlayer
	var rank sql.NullInt32

	err := rows.Scan(
		&p.ID, &p.Time,
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

func (this *SnapshotRepo) FindPossibleCandidates(
	player Player, timeDiff time.Duration,
) ([]SnapshotPlayer, error) {
	min := int(timeDiff / time.Minute)
	killMin, killMax := player.Kills, player.Kills+MAX_KILL_SPEED*min
	deathMin, deathMax := player.Deaths, player.Deaths+MAX_DEATH_SPEED*min
	rankMin := *player.Rank + MIN_RANK_SPEED*min
	rankMax := *player.Rank + MAX_RANK_SPEED*min
	accuracyMin := player.Accuracy + MIN_ACCURACY_SPEED*min
	accuracyMax := player.Accuracy + MAX_ACCURACY_SPEED*min

	q := fmt.Sprintf(`
		SELECT %s
		FROM snapshot as s
		WHERE
			? <= kills AND kills <= ? AND
			? <= deaths AND deaths <= ? AND
			? <= rank AND rank <= ? AND
			? <= accuracy AND accuracy <= ?
	`, this.getPlayerFields("s."))

	log.Printf(
		"searching for candicate %s: "+
			"kills: (%d)[%d, %d], "+
			"deaths: (%d)[%d, %d], "+
			"rank: (%d)[%d, %d], "+
			"accuracy: (%d)[%d, %d], ",
		player.Name,
		player.Kills, killMin, killMax,
		player.Deaths, deathMin, deathMax,
		*player.Rank, rankMin, rankMax,
		player.Accuracy, accuracyMin, accuracyMax,
	)
	rows, err := this.Database.Query(
		q,
		killMin, killMax,
		deathMin, deathMax,
		rankMin, rankMax,
		accuracyMin, accuracyMax,
	)
	if err != nil {
		return nil, err
	}

	players := make([]SnapshotPlayer, 0, 0)

	for rows.Next() {
		p, err := this.scanPlayer(rows)
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	if len(players) == 0 {
		return nil, ERR_NO_CANDIDATE_FOUND
	}

	return players, nil
}

func (this *SnapshotRepo) Update(players []Player) error {
	log.Println("snapshot repo: truncating database...")
	err := this.truncate()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, p := range players {
		err = this.addPlayer(p, now)
		if err != nil {
			log.Printf("snapshot: update: %s", err)
		}
	}

	return nil
}

func (this *SnapshotRepo) addPlayer(player Player, now int64) error {
	q := `
		INSERT INTO snapshot (time, name, country, rank, score, kills, deaths, accuracy)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := this.Database.Exec(q, now,
		player.Name, player.Country, *player.Rank,
		player.Score, player.Kills, player.Deaths, player.Accuracy)

	return err
}

func (this *SnapshotRepo) truncate() error {
	_, err := this.Database.Exec("DELETE FROM snapshot;")
	return err
}

func CreateSnapshotRepo(db *sql.DB) (*SnapshotRepo, error) {
	createSnapshotTable := `CREATE TABLE IF NOT EXISTS snapshot (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time INTEGER,
		name TEXT,
		country TEXT,
		rank INTEGER,
		score INTEGER,
		kills INTEGER,
		deaths INTEGER,
		accuracy INTEGER
	);`
	_, err := db.Exec(createSnapshotTable)
	if err != nil {
		return nil, err
	}

	createIndexes := `
		CREATE INDEX IF NOT EXISTS idx_snapshot_name ON snapshot(name);
		CREATE INDEX IF NOT EXISTS idx_snapshot_kills ON snapshot(kills);
	`
	_, err = db.Exec(createIndexes)
	if err != nil {
		return nil, err
	}

	return &SnapshotRepo{Database: db}, nil
}
