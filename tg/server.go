package tg

import (
	"context"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type Server struct {
	bot    *tgbotapi.BotAPI
	router *tgool.Router
}

func (this *Server) Listen(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := this.bot.GetUpdatesChan(updateConfig)

	tgoolEngine := tgool.NewEngine(
		this.router,
		this.bot,
	)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		tgoolEngine.HandleUpdates(updates)
	}()

	<-ctx.Done()
	this.bot.StopReceivingUpdates()

	wg.Wait()
	log.Println("ends too!")
}
