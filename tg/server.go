package tg

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type Server struct {
	bot    *tgbotapi.BotAPI
	router *tgool.Router
}

func (this *Server) Listen() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := this.bot.GetUpdatesChan(updateConfig)

	tgoolEngine := tgool.NewEngine(
		this.router,
		this.bot,
	)

	log.Println("tg: waiting for updates...")
	tgoolEngine.HandleUpdates(updates)
}
