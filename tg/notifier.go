package tg

import (
	"context"
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/tg/repo"
)

type Notifier struct {
	observer      *core.Observer
	watchlistRepo *repo.WatchlistRepo
	playerRepo    *core.PlayerRepo
	bot           *tgbotapi.BotAPI
	wg            sync.WaitGroup
}

func (this *Notifier) Start(ctx context.Context) {
	log.Println("notifier: started")
	defer log.Println("notifier: stopped")

	this.wg.Add(4)

	throwAway := func(ch chan core.PlayerId, template string) {
		defer this.wg.Done()

		for id := range ch {
			log.Printf(template, id)
		}
	}

	go throwAway(
		this.observer.AddedPlayer,
		"notifier: new player added (id: %d)",
	)
	go throwAway(
		this.observer.UpdatedPlayer,
		"notifier: player updated (id: %d)",
	)
	go func() {
		defer this.wg.Done()

		this.handleEvent(this.observer.GotOnline, true)
	}()
	go func() {
		defer this.wg.Done()
		go this.handleEvent(this.observer.GotOffline, false)
	}()

	<-ctx.Done()
	this.stop()
}

func (this *Notifier) stop() {
	log.Println("notifier: stopping...")
	this.wg.Wait()
}

func (this *Notifier) handleEvent(events chan core.DbPlayer, gotOnline bool) {
	for p := range events {
		player, err := this.playerRepo.GetPlayerByName(p.Player.Name)
		if err != nil {
			log.Printf("notifier: %s", err.Error())
			continue
		}

		chatIds, err := this.watchlistRepo.GetInterested(player.ID)
		if err != nil {
			log.Printf("notifier: %s", err.Error())
			continue
		}

		var msg string
		if gotOnline {
			msg = fmt.Sprintf("🟢 Player %s got online", p.Player.Name)
			log.Printf("notifier: player %s got online", p.Player.Name)
		} else {
			msg = fmt.Sprintf("🔴 Player %s got offline", p.Player.Name)
			log.Printf("notifier: player %s got offline", p.Player.Name)
		}

		momenChatIds := map[int64]bool{
			int64(5701113252): true,
			int64(111398839):  true,
		}
		for _, chatId := range chatIds {
			if momenChatIds[chatId] {
				continue
			}
			this.bot.Send(tgbotapi.NewMessage(chatId, msg))
		}
	}
}
