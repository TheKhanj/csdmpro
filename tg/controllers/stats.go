package controllers

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/tgool"
)

type StatsController struct {
	PlayerRepo *core.PlayerRepo
}

func (this *StatsController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/stats/:page").
		AddMethod("", "Index").
		SetPrefixRoute("/players").
		AddMethod("/:playerId", "PlayerIndex")
}

func (this *StatsController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	page, err := strconv.Atoi(ctx.Params().ByName("page"))
	if err != nil {
		return nil, err
	}

	players, err := this.PlayerRepo.List(page*20, 20+1)
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
		var r1 int = -1
		if p1.Player.Rank != nil {
			r1 = *p1.Player.Rank
		}
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("(%d) %s", r1, p1.Player.Name),
				fmt.Sprintf("/players/%d", p1.ID),
			),
		)

		if p2 != nil {
			var r2 int = -1
			if p2.Player.Rank != nil {
				r2 = *p2.Player.Rank
			}

			row = append(row,
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("(%d) %s", r2, p2.Player.Name),
					fmt.Sprintf("/players/%d", p2.ID),
				),
			)
		}

		return row
	}

	for i := 0; i < len(players)-1; i += 2 {
		p1 := &players[i]
		var p2 *core.DbPlayer = nil
		if i+1 < len(players) {
			p2 = &players[i+1]
		}

		row := getTwoPlayerKeyboard(p1, p2)
		rows = append(rows, row)
	}

	paginationButtons := []tgbotapi.InlineKeyboardButton{}
	if page != 0 {
		paginationButtons = append(
			paginationButtons, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("⬅️ Previous Page"),
				fmt.Sprintf("/stats/%d", page-1),
			),
		)
	}

	if len(players) == 20+1 {
		fmt.Printf("/watchlist/add-players/%d", page+1)
		paginationButtons = append(
			paginationButtons, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Next Page ➡️"),
				fmt.Sprintf("/stats/%d", page+1),
			),
		)
	}

	if len(paginationButtons) != 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(paginationButtons...))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔄 Refresh List",
				"/stats/0",
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

func (this *StatsController) PlayerIndex(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	playerId, err := strconv.Atoi(ctx.Params().ByName("playerId"))
	if err != nil {
		return nil, err
	}

	p, err := this.PlayerRepo.GetPlayer(core.PlayerId(playerId))
	if err != nil {
		return nil, err
	}

	rank := -1
	if p.Player.Rank != nil {
		rank = *p.Player.Rank
	}

	msg := tgbotapi.NewMessage(ctx.GetChatId(),
		fmt.Sprintf(
			`🎮 Player %s Stats

🌍 Country: %s
🏅 Rank: #%d
📈 Score: %d
🔫 Kills: %d
💀 Deaths: %d
🎯 Accuracy: %d%%`,
			p.Player.Name,
			p.Player.Country,
			rank,
			p.Player.Score,
			p.Player.Kills,
			p.Player.Deaths,
			p.Player.Accuracy,
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔄 Refresh",
				fmt.Sprintf("/players/%d", playerId),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔙 Back",
				"/stats/0",
			),
		),
	)

	return msg, nil
}

var _ tgool.Controller = (*StatsController)(nil)
