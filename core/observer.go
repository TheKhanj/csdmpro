package core

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cskr/pubsub/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Topic int

const (
	GotOnlineTopic Topic = iota
	GotOfflineTopic
	AddedPlayerTopic
	UpdatedPlayerTopic
)

type Bus = *pubsub.PubSub[Topic, PlayerId]

type Observer struct {
	Bus Bus

	repo           *PlayerRepo
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
	err := this.repo.Unrank(*player.Rank)
	if err != nil {
		return err
	}

	p, err := this.repo.GetPlayerByName(player.Name)
	not_found := err == ERR_PLAYER_NOT_FOUND
	if not_found {
		pId, err := this.repo.AddPlayer(player)
		if err != nil {
			return err
		}
		p, err = this.repo.GetPlayer(pId)
		if err != nil {
			return err
		}

		this.Bus.Pub(p.ID, AddedPlayerTopic)
	} else {
		err := this.repo.UpdatePlayer(p.ID, player)
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
			err := this.repo.MarkOnline(id)
			if err != nil {
				log.Printf("observer: %s", err)
			}
			this.Bus.Pub(id, GotOnlineTopic)
		}
	}

	for id, wasOnline := range wasOnline {
		if wasOnline && !isOnline[id] {
			err := this.repo.MarkOffline(id)
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
		dbp, err := this.repo.GetPlayerByName(p.Name)
		if err != nil {
			return nil, err
		}

		ret[dbp.ID] = true
	}

	return ret, nil
}

func (this *Observer) getWasOnline() (OnlineMap, error) {
	wereOnlines, err := this.repo.Onlines()
	if err != nil {
		return nil, err
	}

	ret := make(OnlineMap, len(wereOnlines))

	for _, p := range wereOnlines {
		ret[p.ID] = true
	}

	return ret, nil
}

func (this *Observer) observeStatsPage(page int) error {
	players, err := this.crawler.Stats(page)
	if err != nil {
		return err
	}

	for _, player := range players {
		err := this.handlePlayer(player)
		if err != nil {
			log.Printf("observer: stats: page %d: %s", page, err)
			continue
		}
	}

	return nil
}

func (this *Observer) observeStats(ctx context.Context) {
	pageCount := 20
	if os.Getenv("ENV") == "dev" {
		pageCount = 1
	}

	for page := 1; page <= pageCount; page++ {
		select {
		case <-ctx.Done():
			return
		default:
			err := this.observeStatsPage(page)
			if err != nil {
				log.Printf("observer: stats: page %d: %s", page, err)
			}
		}
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
	repo *PlayerRepo, crawler Crawler,
	statsInterval time.Duration, onlineInterval time.Duration,
) *Observer {
	return &Observer{
		Bus: pubsub.New[Topic, PlayerId](0),

		repo:           repo,
		crawler:        crawler,
		statsInterval:  statsInterval,
		onlineInterval: onlineInterval,
	}
}
