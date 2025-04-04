package core

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Observer struct {
	Repo           *PlayerRepo
	Crawler        Crawler
	GotOnline      chan DbPlayer
	GotOffline     chan DbPlayer
	StatsInterval  time.Duration
	OnlineInterval time.Duration

	shuttingDown bool
	wg           sync.WaitGroup
}

func (this *Observer) observeOnlinePlayers() error {
	players, err := this.Crawler.Online()
	if err != nil {
		return err
	}

	pErr := func(err error) {
		log.Println(fmt.Errorf("observer: online: %s", err.Error()))
	}

	for _, player := range players {
		p, err := this.Repo.GetPlayerByName(player.Name)
		if err != nil && err != ERR_PLAYER_NOT_FOUND {
			pErr(err)
			continue
		}
		exists := err != ERR_PLAYER_NOT_FOUND
		if !exists {
			err = this.Repo.AddPlayer(player)
			if err != nil {
				pErr(err)
				continue
			}
		}

		isOnlineAlready, err := this.Repo.IsOnline(player.Name)
		if err != nil {
			pErr(err)
			continue
		}
		if isOnlineAlready {
			continue
		}

		err = this.Repo.AddOnlinePlayer(PlayerId(p.ID))
		if err != nil {
			pErr(err)
			continue
		}

		this.GotOnline <- p
	}

	prevOnlines, err := this.Repo.Onlines()
	if err != nil {
		return err
	}

	for _, prevOnline := range prevOnlines {
		isStillOnline := false
		for _, player := range players {
			if prevOnline.Player.Name == player.Name {
				isStillOnline = true
				break
			}
		}

		if isStillOnline {
			continue
		}

		err = this.Repo.RemoveOnlinePlayer(prevOnline.ID)
		if err != nil {
			pErr(err)
			continue
		}

		this.GotOffline <- prevOnline
	}

	return nil
}

func (this *Observer) observePlayersPage(page int) error {
	players, err := this.Crawler.Stats(page)
	if err != nil {
		return err
	}

	for _, player := range players {
		_, err := this.Repo.GetPlayerByName(player.Name)
		if err != nil && err != ERR_PLAYER_NOT_FOUND {
			log.Println(err)
			continue
		}
		exists := err != ERR_PLAYER_NOT_FOUND

		if exists {
			continue
		}

		err = this.Repo.AddPlayer(player)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (this *Observer) observePlayers() {
	for page := 1; page <= 5; page++ {
		err := this.observePlayersPage(page)
		if err != nil {
			log.Println(err)
		}
	}
}

func (this *Observer) Start(ctx context.Context) {
	this.wg.Add(2)

	go func() {
		defer this.wg.Done()

		for {
			if this.shuttingDown {
				return
			}

			err := this.observeOnlinePlayers()
			if err != nil {
				log.Println(err)
			}
			if this.shuttingDown {
				return
			}

			time.Sleep(this.OnlineInterval)
		}
	}()

	go func() {
		defer this.wg.Done()

		for {
			if this.shuttingDown {
				return
			}
			this.observePlayers()
			if this.shuttingDown {
				return
			}

			time.Sleep(this.StatsInterval)
		}
	}()

	<-ctx.Done()
	this.stop()
}

func (this *Observer) stop() {
	this.shuttingDown = true

	this.wg.Wait()
}
