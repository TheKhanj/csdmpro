package controllers

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/tg/service"
	"github.com/thekhanj/tgool"
)

type WatchlistController struct {
	PlayerRepo *core.PlayerRepo
	Service    *service.WatchlistService
}

func (this *WatchlistController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/watchlist").
		AddMethod("", "Index").
		AddMethod("add-players/:page", "AddPlayersIndex").
		AddMethod("remove-players", "RemovePlayersIndex").
		AddMethod("a/post/players/:playerId", "AddPlayer").
		AddMethod("a/delete/players/:playerId", "RemovePlayer")
}

func (this *WatchlistController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `🕵️ Watchlist Feature

Keep track of your favorite players! I'll notify you when they join or leave the server.
Just add them to your watchlist, and I'll handle the rest. 🚀`

	tps, err := this.Service.GetTracking(chatId)
	if err != nil {
		return nil, err
	}

	txt += "\n\n"
	if len(tps) == 0 {
		txt += "👀 You’re not tracking anyone yet."
	} else {
		txt += "👀 Currently Tracked Players:\n"
		for _, tp := range tps {
			// todo: show red for offline

			var status string
			if tp.IsOnline {
				status = "🟢"
			} else {
				status = "🔴"
			}
			txt += fmt.Sprintf("%s %s\n", status, tp.Player.Name)
		}
	}

	msg := tgbotapi.NewMessage(chatId, txt)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"➕ Add to Watchlist",
				"/watchlist/add-players/0",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"➖ Remove from Watchlist",
				"/watchlist/remove-players",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔄 Refresh List",
				"/watchlist",
			),
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (this *WatchlistController) AddPlayersIndex(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()
	page, err := strconv.Atoi(ctx.Params().ByName("page"))
	if err != nil {
		return nil, err
	}
	txt := `➕ Add to Watchlist

Select a player to add to your watchlist`

	players, err := this.PlayerRepo.List(page*20, 20)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(chatId, txt)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	getTwoPlayerKeyboard := func(
		p1 *core.Player, p2 *core.Player,
	) ([]tgbotapi.InlineKeyboardButton, error) {
		// TODO: what is this shit? fix it
		p1Id, err := this.PlayerRepo.GetPlayerId(p1.Name)
		if err != nil {
			return nil, err
		}
		p2Id, err := this.PlayerRepo.GetPlayerId(p2.Name)
		if err != nil {
			return nil, err
		}

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("✔️ %s", p1.Name),
				fmt.Sprintf("/watchlist/a/post/players/%d", p1Id),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("✔️ %s", p2.Name),
				fmt.Sprintf("/watchlist/a/post/players/%d", p2Id),
			),
		)

		return row, nil
	}

	for i := 0; i < len(players); i += 2 {
		p1 := players[i]
		p2 := players[i+1]
		row, err := getTwoPlayerKeyboard(&p1, &p2)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	paginationButtons := []tgbotapi.InlineKeyboardButton{}
	if page != 0 {
		paginationButtons = append(
			paginationButtons, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("⬅️ Previous Page"),
				fmt.Sprintf("/watchlist/add-players/%d", page-1),
			),
		)
	}

	// TODO: fix
	// if page!=end

	fmt.Printf("/watchlist/add-players/%d", page+1)
	paginationButtons = append(
		paginationButtons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Next Page ➡️"),
			fmt.Sprintf("/watchlist/add-players/%d", page+1),
		),
	)

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(paginationButtons...),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔙 Back",
				"/watchlist",
			),
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (this *WatchlistController) AddPlayer(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	selectedPlayerId := ctx.Params().ByName("playerId")
	log.Printf("selected player (%s)", selectedPlayerId)

	ctx.Redirect("/watchlist")

	ctx.Bot().Request(
		tgbotapi.NewCallback(
			ctx.Update().CallbackQuery.ID,
			fmt.Sprintf("player %s added to watchlist", selectedPlayerId),
		),
	)

	return this.Index(ctx)
}

func (this *WatchlistController) RemovePlayersIndex(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `➖ Remove from Watchlist

Select a player to remove from your watchlist`

	currentPlayers := []string{"player-1", "player-2"}

	msg := tgbotapi.NewMessage(chatId, txt)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, player := range currentPlayers {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("🚫 %s", player),
					fmt.Sprintf("/watchlist/a/delete/players/%s", player),
				),
			),
		)
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔙 Back",
				"/watchlist",
			),
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (this *WatchlistController) RemovePlayer(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	selectedPlayer := ctx.Params().ByName("player")
	log.Printf("selected player (%s)", selectedPlayer)

	ctx.Redirect("/watchlist")

	ctx.Bot().Request(
		tgbotapi.NewCallback(
			ctx.Update().CallbackQuery.ID,
			fmt.Sprintf("player %s removed from watchlist", selectedPlayer),
		),
	)

	return this.Index(ctx)
}

var _ tgool.Controller = (*WatchlistController)(nil)
