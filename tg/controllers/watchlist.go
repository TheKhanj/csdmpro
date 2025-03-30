package controllers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type WatchlistController struct{}

func (this *WatchlistController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/watchlist").
		AddMethod("", "Index").
		AddMethod("/users-add", "UsersAddIndex").
		AddMethod("/users-remove", "UsersRemoveIndex")
}

func (this *WatchlistController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `🕵️ Watchlist Feature

Keep track of your favorite players! I'll notify you when they join or leave the server.
Just add them to your watchlist, and I'll handle the rest. 🚀`

	currentUsers := []string{"user-1", "user-2"}
	txt += "\n\n"
	if len(currentUsers) == 0 {
		txt += "👀 You’re not tracking anyone yet."
	} else {
		txt += "👀 Currently Tracked Players:\n\n"
		for _, user := range currentUsers {
			// todo: show red for offline
			txt += fmt.Sprintf("🟢 %s\n", user)
		}
	}

	msg := tgbotapi.NewMessage(chatId, txt)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"➕ Add to Watchlist",
				"/watchlist/users-add",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"➖ Remove from Watchlist",
				"/watchlist/users-remove",
			),
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (this *WatchlistController) UsersAddIndex(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `➕ Add to Watchlist

Select a player to add to your watchlist`

	allUsers := []string{"user-1", "user-2"}

	msg := tgbotapi.NewMessage(chatId, txt)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	// TODO: add pagination
	for _, user := range allUsers {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("✔️ %s", user),
					"/watchlist/users-add/%s",
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

func (this *WatchlistController) UsersRemoveIndex(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `➖ Remove from Watchlist

Select a player to remove from your watchlist`

	currentUsers := []string{"user-1", "user-2"}

	msg := tgbotapi.NewMessage(chatId, txt)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, user := range currentUsers {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("🚫 %s", user),
					"/watchlist/users-remove/%s",
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

var _ tgool.Controller = (*WatchlistController)(nil)
