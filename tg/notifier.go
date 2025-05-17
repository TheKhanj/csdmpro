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
	gotOnline      chan core.PlayerId
	gotOffline     chan core.PlayerId
	usernameChange chan core.UsernameChange

	observer      *core.Observer
	watchlistRepo *repo.WatchlistRepo
	playerRepo    *core.PlayerRepo
	bot           *tgbotapi.BotAPI
	wg            sync.WaitGroup
}

func (this *Notifier) Start(ctx context.Context) {
	log.Println("notifier: started")
	defer log.Println("notifier: stopped")

	this.wg.Add(3)

	go func() {
		defer this.wg.Done()
		this.gotOnline = this.observer.Bus.Sub(core.GotOnlineTopic)

		this.handleEvent(this.gotOnline, true)
	}()
	go func() {
		defer this.wg.Done()
		this.gotOffline = this.observer.Bus.Sub(core.GotOfflineTopic)

		this.handleEvent(this.gotOffline, false)
	}()
	go func() {
		defer this.wg.Done()
		this.usernameChange = this.observer.UsernameChangeBus.Sub(core.UsernameChangeTopic)

		this.handleUsernameChange()
	}()

	<-ctx.Done()
	this.stop()
}

func (this *Notifier) handleUsernameChange() {
	for u := range this.usernameChange {
		log.Printf("notifier: username changed %s -> %s", u.From, u.To)
		msg := fmt.Sprintf("👀 Player %s changed their username to %s", u.From, u.To)
		myChatId := int64(209245565)

		this.bot.Send(tgbotapi.NewMessage(myChatId, msg))
	}
}

func (this *Notifier) stop() {
	log.Println("notifier: stopping...")

	go this.observer.Bus.Unsub(this.gotOnline)
	go this.observer.Bus.Unsub(this.gotOffline)
	go this.observer.UsernameChangeBus.Unsub(this.usernameChange)

	this.wg.Wait()
}

func (this *Notifier) handleEvent(events chan core.PlayerId, gotOnline bool) {
	for playerId := range events {
		player, err := this.playerRepo.GetPlayer(playerId)
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
			msg = fmt.Sprintf("🟢 Player %s got online", player.Player.Name)
			log.Printf("notifier: player %s got online", player.Player.Name)
		} else {
			msg = fmt.Sprintf("🔴 Player %s got offline", player.Player.Name)
			log.Printf("notifier: player %s got offline", player.Player.Name)
		}

		for _, chatId := range chatIds {
			this.bot.Send(tgbotapi.NewMessage(chatId, msg))
		}
	}
}
