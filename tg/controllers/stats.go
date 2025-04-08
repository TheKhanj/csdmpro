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
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("(%d) %s", p1.Player.Rank, p1.Player.Name),
				fmt.Sprintf("/players/%d", p1.ID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("(%d) %s", p2.Player.Rank, p2.Player.Name),
				fmt.Sprintf("/players/%d", p2.ID),
			),
		)

		return row
	}

	for i := 0; i < len(players)-1; i += 2 {
		row := getTwoPlayerKeyboard(&players[i], &players[i+1])
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
			p.Player.Rank,
			p.Player.Score,
			p.Player.Kills,
			p.Player.Deaths,
			p.Player.Accuracy,
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
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
