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
	GotOnline      chan Player
	GotOffline     chan Player
	StatsInterval  time.Duration
	OnlineInterval time.Duration
	Ctx            context.Context

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
		exists, err := this.Repo.PlayerExists(player.Name)
		if err != nil {
			pErr(err)
			continue
		}
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

		playerId, err := this.Repo.GetPlayerId(player.Name)
		if err != nil {
			pErr(err)
			continue
		}

		err = this.Repo.AddOnlinePlayer(playerId)
		if err != nil {
			pErr(err)
			continue
		}

		this.GotOnline <- player
	}

	prevOnlines, err := this.Repo.Onlines()
	if err != nil {
		return err
	}

	for _, prevOnline := range prevOnlines {
		isStillOnline := false
		for _, player := range players {
			if prevOnline.Name == player.Name {
				isStillOnline = true
				break
			}
		}

		if isStillOnline {
			continue
		}

		playerId, err := this.Repo.GetPlayerId(prevOnline.Name)
		if err != nil {
			pErr(err)
			continue
		}
		err = this.Repo.RemoveOnlinePlayer(playerId)
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
		exists, err := this.Repo.PlayerExists(player.Name)
		if err != nil {
			log.Println(err)
			continue
		}

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

func (this *Observer) Start() {
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

			ctx, cancel := context.WithTimeout(this.Ctx, this.OnlineInterval)
			defer cancel()
			<-ctx.Done()
		}
	}()

	go func() {
		defer this.wg.Done()

		for {
			if this.shuttingDown {
				return
			}

			this.observePlayers()

			ctx, cancel := context.WithTimeout(this.Ctx, this.StatsInterval)
			defer cancel()
			<-ctx.Done()
		}
	}()

	this.wg.Wait()
}

func (this *Observer) Stop() {
	this.shuttingDown = true

	this.wg.Wait()
}
