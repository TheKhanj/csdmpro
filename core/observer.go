package core

import (
	"context"
	"log"
	"math"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/cskr/pubsub/v2"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MAX_KILL_SPEED     = 20
	MAX_DEATH_SPEED    = 10
	MIN_RANK_SPEED     = -10
	MAX_RANK_SPEED     = 100
	MIN_ACCURACY_SPEED = -1
	MAX_ACCURACY_SPEED = 1
)

type Topic int

const (
	GotOnlineTopic Topic = iota
	GotOfflineTopic
	AddedPlayerTopic
	UpdatedPlayerTopic
)

type Bus = *pubsub.PubSub[Topic, PlayerId]

const UsernameChangeTopic = iota

type UsernameChange struct {
	From string
	To   string
}

type UsernameChangeBus = *pubsub.PubSub[int, UsernameChange]

type Observer struct {
	Bus               Bus
	UsernameChangeBus UsernameChangeBus

	playerRepo     *PlayerRepo
	snapshotRepo   *SnapshotRepo
	crawler        Crawler
	statsInterval  time.Duration
	onlineInterval time.Duration

	wg sync.WaitGroup
}

func (this *Observer) observeOnlinePlayers() error {
	players, err := this.crawler.Online()
	if err != nil {
		return err
	}

	for _, player := range players {
		err := this.handlePlayer(player)
		if err != nil {
			log.Printf("observer: onlines: %s", err)
			continue
		}
	}

	return this.handleOnlinePlayers(players)
}

func (this *Observer) handlePlayer(player Player) error {
	err := this.playerRepo.Unrank(*player.Rank)
	if err != nil {
		return err
	}

	p, err := this.playerRepo.GetPlayerByName(player.Name)
	not_found := err == ERR_PLAYER_NOT_FOUND
	if not_found {
		pId, err := this.playerRepo.AddPlayer(player)
		if err != nil {
			return err
		}
		p, err = this.playerRepo.GetPlayer(pId)
		if err != nil {
			return err
		}

		this.Bus.Pub(p.ID, AddedPlayerTopic)
	} else {
		err := this.playerRepo.UpdatePlayer(p.ID, player)
		if err != nil {
			return err
		}

		this.Bus.Pub(p.ID, UpdatedPlayerTopic)
	}

	return nil
}

func (this *Observer) handleOnlinePlayers(areOnlines []Player) error {
	wasOnline, err := this.getWasOnline()
	if err != nil {
		return err
	}
	isOnline, err := this.getIsOnline(areOnlines)
	if err != nil {
		return err
	}

	for id, isOnline := range isOnline {
		if isOnline && !wasOnline[id] {
			err := this.playerRepo.MarkOnline(id)
			if err != nil {
				log.Printf("observer: %s", err)
			}
			this.Bus.Pub(id, GotOnlineTopic)
		}
	}

	for id, wasOnline := range wasOnline {
		if wasOnline && !isOnline[id] {
			err := this.playerRepo.MarkOffline(id)
			if err != nil {
				log.Printf("observer: %s", err)
			}
			this.Bus.Pub(id, GotOfflineTopic)
		}
	}

	return nil
}

type OnlineMap = map[PlayerId]bool

func (this *Observer) getIsOnline(players []Player) (OnlineMap, error) {
	ret := make(OnlineMap, 0)

	for _, p := range players {
		dbp, err := this.playerRepo.GetPlayerByName(p.Name)
		if err != nil {
			return nil, err
		}

		ret[dbp.ID] = true
	}

	return ret, nil
}

func (this *Observer) getWasOnline() (OnlineMap, error) {
	wereOnlines, err := this.playerRepo.Onlines()
	if err != nil {
		return nil, err
	}

	ret := make(OnlineMap, len(wereOnlines))

	for _, p := range wereOnlines {
		ret[p.ID] = true
	}

	return ret, nil
}

func (this *Observer) observeStatsPage(page int) ([]Player, error) {
	players, err := this.crawler.Stats(page)
	if err != nil {
		return nil, err
	}

	snapshotTime := time.Now().Unix()
	for _, player := range players {
		if *player.Rank <= 200 {
			err := this.handlePossibleUsernameChange(player, snapshotTime)
			if err != nil {
				log.Printf("observer: stats: page %d: username change: %s", page, err)
				continue
			}
		}
		err = this.handlePlayer(player)
		if err != nil {
			log.Printf("observer: stats: page %d: %s", page, err)
			continue
		}
	}

	return players, nil
}

func (this *Observer) handlePossibleUsernameChange(
	player Player, snapshotTime int64,
) error {
	prevSnapshot, err := this.snapshotRepo.Get(player.Name)
	if err != nil {
		return err
	}

	timeDiffUnix := snapshotTime - prevSnapshot.Time
	timeDiff := time.Duration(timeDiffUnix) * time.Second
	if this.isSimilarPlayer(prevSnapshot.Player, player, timeDiff) {
		return nil
	}
	log.Printf("user %s not found in snapshots", player.Name)

	from, err := this.findPreviousUsername(player, timeDiff)
	if err != nil {
		return err
	}

	ev := UsernameChange{
		From: from,
		To:   player.Name,
	}

	go this.UsernameChangeBus.Pub(ev, UsernameChangeTopic)
	return nil
}

func (this *Observer) findPreviousUsername(
	player Player, timeDiff time.Duration,
) (string, error) {
	prevs, err := this.snapshotRepo.FindPossibleCandidates(player, timeDiff)
	if err != nil {
		return "", err
	}

	if len(prevs) == 1 {
		return prevs[0].Player.Name, nil
	}

	type Candidate struct {
		variance int64
		player   Player
	}
	c := make([]Candidate, len(prevs))
	for _, prev := range prevs {
		c = append(c,
			Candidate{
				variance: this.calculateVariance(prev.Player, player, timeDiff),
				player:   player,
			},
		)
	}

	slices.SortFunc(c, func(a, b Candidate) int {
		if a.variance < b.variance {
			return 1
		}
		return 0
	})

	return c[0].player.Name, nil
}

func (this *Observer) calculateVariance(
	prev, curr Player, timeDiff time.Duration,
) int64 {
	if !this.isSimilarPlayer(prev, curr, timeDiff) {
		return math.MaxInt64
	}

	min := timeDiff / time.Minute
	killSpeed := float64(curr.Kills-prev.Kills) / float64(min)
	deathSpeed := float64(curr.Deaths-prev.Deaths) / float64(min)
	rankSpeed := float64(*curr.Rank-*prev.Rank) / float64(min)
	accuracySpeed := float64(curr.Accuracy-prev.Accuracy) / float64(min)

	ret := killSpeed*killSpeed +
		deathSpeed*deathSpeed +
		rankSpeed*rankSpeed +
		accuracySpeed*accuracySpeed*400

	return int64(ret)
}

func (this *Observer) isSimilarPlayer(
	prev, curr Player, timeDiff time.Duration,
) bool {
	if curr.Kills < prev.Kills || curr.Deaths < prev.Deaths {
		return false
	}

	min := timeDiff / time.Minute
	killSpeed := float64(curr.Kills-prev.Kills) / float64(min)
	deathSpeed := float64(curr.Deaths-prev.Deaths) / float64(min)
	rankSpeed := float64(*curr.Rank-*prev.Rank) / float64(min)
	accuracySpeed := float64(curr.Accuracy-prev.Accuracy) / float64(min)

	return killSpeed <= MAX_KILL_SPEED &&
		deathSpeed <= MAX_DEATH_SPEED &&
		MIN_RANK_SPEED <= rankSpeed && rankSpeed <= MAX_RANK_SPEED &&
		MIN_ACCURACY_SPEED <= accuracySpeed && accuracySpeed <= MAX_ACCURACY_SPEED
}

func (this *Observer) observeStats(ctx context.Context) {
	pageCount := 20
	if os.Getenv("ENV") == "dev" {
		pageCount = 1
	}

	snapshot := make([]Player, 0)

	for page := 1; page <= pageCount; page++ {
		select {
		case <-ctx.Done():
			return
		default:
			players, err := this.observeStatsPage(page)
			if err != nil {
				log.Printf("observer: stats: page %d: %s", page, err)
			} else {
				snapshot = append(snapshot, players...)
			}
		}
	}

	err := this.snapshotRepo.Update(snapshot)
	if err != nil {
		log.Printf("observer: stats: updating snapshot: %s", err)
	}
}

func (this *Observer) Start(ctx context.Context) {
	log.Println("observer: started")
	defer log.Println("observer: stopped")

	this.wg.Add(2)

	go func() {
		log.Println("observer: started observing onlines")
		defer this.wg.Done()
		defer log.Println("observer: stopped observing onlines")

		for {
			err := this.observeOnlinePlayers()
			if err != nil {
				log.Println(err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(this.onlineInterval):
			}
		}
	}()

	go func() {
		log.Println("observer: started observing stats")
		defer this.wg.Done()
		defer log.Println("observer: stopped observing stats")

		for {
			this.observeStats(ctx)

			select {
			case <-ctx.Done():
				return
			case <-time.After(this.statsInterval):
			}
		}
	}()

	<-ctx.Done()
	this.stop()
}

func (this *Observer) stop() {
	log.Println("observer: stopping...")

	this.wg.Wait()
}

func NewObserver(
	repo *PlayerRepo, snapshotRepo *SnapshotRepo,
	crawler Crawler,
	statsInterval time.Duration, onlineInterval time.Duration,
) *Observer {
	return &Observer{
		Bus:               pubsub.New[Topic, PlayerId](0),
		UsernameChangeBus: pubsub.New[int, UsernameChange](0),

		playerRepo:     repo,
		snapshotRepo:   snapshotRepo,
		crawler:        crawler,
		statsInterval:  statsInterval,
		onlineInterval: onlineInterval,
	}
}
