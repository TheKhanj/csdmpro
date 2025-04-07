package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type BilakhController struct{}

func (this *BilakhController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/bilakh").
		AddMethod("", "Index").
		SetPrefixRoute("/start").
		AddMethod("", "Index")
}

func (this *BilakhController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := "Surprise! 🎁 It’s a bilakh. Straight from the heart❤️. Or somewhere not so close to it.🍆"

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"More please!",
				"/bilakh/more-please",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"I feel honored",
				"/bilakh/feel-honored",
			),
		),
	)

	return msg, nil
}

var _ tgool.Controller = (*BilakhController)(nil)
