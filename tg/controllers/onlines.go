package controllers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/tgool"
)

type OnlinesController struct {
	PlayerRepo *core.PlayerRepo
}

func (this *OnlinesController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/onlines/").
		AddMethod("", "Index")
}

func (this *OnlinesController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	players, err := this.PlayerRepo.Onlines()
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(
		ctx.GetChatId(),
		"📊 Live player stats from the battlefield — updated in real-time.",
	)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	getTwoPlayerKeyboard := func(
		p1 *core.DbPlayer, p2 *core.DbPlayer,
	) []tgbotapi.InlineKeyboardButton {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("🟢 %s", p1.Player.Name),
				fmt.Sprintf("/players/%d", p1.ID),
			),
		)
		if p2 != nil {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("🟢 %s", p2.Player.Name),
				fmt.Sprintf("/players/%d", p2.ID),
			),
			)
		}

		return row
	}

	for i := 0; i < len(players); i += 2 {
		p1 := &players[i]

		var p2 *core.DbPlayer = nil
		if i+1 < len(players) {
			p2 = &players[i]
		}

		row := getTwoPlayerKeyboard(p1, p2)
		rows = append(rows, row)
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔄 Refresh List",
				"/onlines",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔙 Back",
				"/start",
			),
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

var _ tgool.Controller = (*OnlinesController)(nil)
