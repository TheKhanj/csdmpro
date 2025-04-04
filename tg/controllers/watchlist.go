package controllers

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/tg/repo"
	"github.com/thekhanj/csdmpro/tg/service"
	"github.com/thekhanj/tgool"
)

type WatchlistController struct {
	PlayerRepo    *core.PlayerRepo
	WatchlistRepo *repo.WatchlistRepo
	Service       *service.WatchlistService
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
			txt += fmt.Sprintf("%s %s\n", status, tp.DbPlayer.Player.Name)
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
		p1 *core.DbPlayer, p2 *core.DbPlayer,
	) ([]tgbotapi.InlineKeyboardButton, error) {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("✔️ %s", p1.Player.Name),
				fmt.Sprintf("/watchlist/a/post/players/%d", p1.ID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("✔️ %s", p2.Player.Name),
				fmt.Sprintf("/watchlist/a/post/players/%d", p2.ID),
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
	chatId := ctx.GetChatId()
	playerId, err := strconv.Atoi(ctx.Params().ByName("playerId"))
	if err != nil {
		return nil, err
	}

	id := core.PlayerId(playerId)
	player, err := this.PlayerRepo.GetPlayer(id)
	if err != nil {
		return nil, err
	}

	isWatched, err := this.WatchlistRepo.IsInWatchlist(chatId, id)
	if err != nil {
		return nil, err
	}

	if isWatched {
		ctx.Bot().Request(
			tgbotapi.NewCallback(
				ctx.Update().CallbackQuery.ID,
				fmt.Sprintf("player %s is already in the watchlist", player.Player.Name),
			),
		)
	} else {
		err = this.WatchlistRepo.Add(chatId, id)
		if err != nil {
			return nil, err
		}

		ctx.Bot().Request(
			tgbotapi.NewCallback(
				ctx.Update().CallbackQuery.ID,
				fmt.Sprintf("player %s added into the watchlist", player.Player.Name),
			),
		)
	}

	ctx.Redirect("/watchlist")

	return this.Index(ctx)
}

func (this *WatchlistController) RemovePlayersIndex(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `➖ Remove from Watchlist

Select a player to remove from your watchlist`

	tps, err := this.Service.GetTracking(chatId)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(chatId, txt)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, tp := range tps {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("🚫 %s", tp.DbPlayer.Player.Name),
					fmt.Sprintf("/watchlist/a/delete/players/%d", tp.DbPlayer.ID),
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
	chatId := ctx.GetChatId()
	playerId, err := strconv.Atoi(ctx.Params().ByName("playerId"))
	if err != nil {
		return nil, err
	}

	id := core.PlayerId(playerId)
	player, err := this.PlayerRepo.GetPlayer(id)
	if err != nil {
		return nil, err
	}

	err = this.WatchlistRepo.Remove(chatId, id)
	if err != nil {
		return nil, err
	}

	ctx.Bot().Request(
		tgbotapi.NewCallback(
			ctx.Update().CallbackQuery.ID,
			fmt.Sprintf("player %s removed from the watchlist", player.Player.Name),
		),
	)

	ctx.Redirect("/watchlist")

	return this.Index(ctx)
}

var _ tgool.Controller = (*WatchlistController)(nil)
