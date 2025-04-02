package tg

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/tg/repo"
)

type Notifier struct {
	gotOnline     chan core.Player
	gotOffline    chan core.Player
	watchlistRepo *repo.WatchlistRepo
	playerRepo    *core.PlayerRepo
	bot           *tgbotapi.BotAPI
}

func (this *Notifier) Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		this.handleEvent(this.gotOnline, true)
	}()

	go func() {
		defer wg.Done()
		go this.handleEvent(this.gotOffline, false)
	}()

	wg.Wait()
}

func (this *Notifier) handleEvent(events chan core.Player, gotOnline bool) {
	for p := range events {
		playerId, err := this.playerRepo.GetPlayerId(p.Name)
		if err != nil {
			log.Printf("notifier: %s", err.Error())
			continue
		}

		chatIds, err := this.watchlistRepo.GetInterested(playerId)
		if err != nil {
			log.Printf("notifier: %s", err.Error())
			continue
		}

		var msg string
		if gotOnline {
			msg = fmt.Sprintf("🟢 Player *%s* got online", p.Name)
			log.Printf("notifier: player %s got online", p.Name)
		} else {
			msg = fmt.Sprintf("🔴 Player *%s* got offline", p.Name)
			log.Printf("notifier: player %s got offline", p.Name)
		}

		for _, chatId := range chatIds {
			this.bot.Send(tgbotapi.NewMessage(chatId, msg))
		}
	}
}
