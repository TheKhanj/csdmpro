package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type StartController struct {
	WatchlistController *WatchlistController
}

func (this *StartController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/start").
		AddMethod("", "Index")
}

func (this *StartController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	return this.WatchlistController.Index(ctx)
}

var _ tgool.Controller = (*StartController)(nil)
