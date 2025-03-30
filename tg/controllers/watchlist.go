package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type WatchlistController struct{}

func (this *WatchlistController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/watchlist").
		AddMethod("", "Index").WithTitle("Watch list")
}

func (this *WatchlistController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	msg := tgbotapi.NewMessage(
		chatId,
		"First controller",
	)

	return msg, nil
}

var _ tgool.Controller = (*WatchlistController)(nil)
