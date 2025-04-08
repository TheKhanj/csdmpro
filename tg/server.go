package tg

import (
	"context"
	"log"
	"os"

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

	tgDone := make(chan struct{})
	go func() {
		defer close(tgDone)

		tgoolEngine.HandleUpdates(updates)
	}()

	<-ctx.Done()
	this.bot.StopReceivingUpdates()

	if os.Getenv("ENV") == "dev" {
		log.Println("tg: server forecfully stopped")
		return
	}

	<-tgDone
}
