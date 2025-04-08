package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type StartController struct{}

func (this *StartController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/start").
		AddMethod("", "Index")
}

func (this *StartController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(
		ctx.GetChatId(),
		`👋 Welcome to the csdmpro Bot!

Track your performance, keep an eye on your watchlist, or check who's online — all in one place.

Choose a page to get started 👇`,
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"📊 Stats",
				"/stats/0",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"🟢 Online Players",
				"/onlines",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"👁️ Watchlist",
				"/watchlist",
			),
		),
	)

	return msg, nil
}

var _ tgool.Controller = (*StartController)(nil)
