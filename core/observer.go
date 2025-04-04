package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Observer struct {
	Repo           *PlayerRepo
	Crawler        Crawler
	GotOnline      chan DbPlayer
	GotOffline     chan DbPlayer
	UpdatedPlayer  chan PlayerId
	AddedPlayer    chan PlayerId
	StatsInterval  time.Duration
	OnlineInterval time.Duration

	wg sync.WaitGroup
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
		not_found := err == ERR_PLAYER_NOT_FOUND
		if err != nil && !not_found {
			pErr(err)
			continue
		}

		if not_found {
			id, err := this.Repo.AddPlayer(player)
			if err != nil {
				pErr(err)
				continue
			}
			this.AddedPlayer <- id

			p, err = this.Repo.GetPlayerByName(player.Name)
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
		p, err := this.Repo.GetPlayerByName(player.Name)
		not_found := err == ERR_PLAYER_NOT_FOUND
		if err != nil && !not_found {
			log.Println(err)
			continue
		}

		if !not_found {
			err := this.Repo.UpdatePlayer(p.ID, player)
			if err != nil {
				log.Println(err)
			}
			this.UpdatedPlayer <- p.ID
			continue
		}

		id, err := this.Repo.AddPlayer(player)
		if err != nil {
			log.Println(err)
		}
		this.AddedPlayer <- id
	}

	return nil
}

func (this *Observer) observePlayers() {
	pageCount := 20
	if os.Getenv("ENV") == "dev" {
		pageCount = 1
	}

	for page := 1; page <= pageCount; page++ {
		err := this.observePlayersPage(page)
		if err != nil {
			log.Println(err)
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
			case <-time.After(this.OnlineInterval):
			}
		}
	}()

	go func() {
		log.Println("observer: started observing stats")
		defer this.wg.Done()
		defer log.Println("observer: stopped observing stats")

		for {
			this.observePlayers()

			select {
			case <-ctx.Done():
				return
			case <-time.After(this.StatsInterval):
			}
		}
	}()

	<-ctx.Done()
	this.stop()
}

func (this *Observer) stop() {
	log.Println("observer: stopping...")

	this.wg.Wait()
	close(this.GotOnline)
	close(this.GotOffline)
	close(this.UpdatedPlayer)
	close(this.AddedPlayer)
}
