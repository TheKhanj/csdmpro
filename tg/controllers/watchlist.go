package controllers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/tgool"
)

type WatchlistController struct {
	Repo *core.PlayerRepo
}

func (this *WatchlistController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/watchlist").
		AddMethod("", "Index").
		AddMethod("add-users", "AddUsersIndex").
		AddMethod("remove-users", "RemoveUsersIndex").
		AddMethod("a/post/users/:user", "AddUser").
		AddMethod("a/delete/users/:user", "RemoveUser")
}

func (this *WatchlistController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `🕵️ Watchlist Feature

Keep track of your favorite players! I'll notify you when they join or leave the server.
Just add them to your watchlist, and I'll handle the rest. 🚀`

	players, err := this.Repo.List(0, 10)
	if err != nil {
		return nil, err
	}
	currentUsers := []string{}
	for _, p := range players {
		currentUsers = append(currentUsers, p.Name)
	}

	txt += "\n\n"
	if len(currentUsers) == 0 {
		txt += "👀 You’re not tracking anyone yet."
	} else {
		txt += "👀 Currently Tracked Players:\n"
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
				"/watchlist/add-users",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"➖ Remove from Watchlist",
				"/watchlist/remove-users",
			),
		),
	)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func (this *WatchlistController) AddUsersIndex(
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
					fmt.Sprintf("/watchlist/a/post/users/%s", user),
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

func (this *WatchlistController) AddUser(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	selectedUser := ctx.Params().ByName("user")
	log.Printf("selected user (%s)", selectedUser)

	ctx.Redirect("/watchlist")

	ctx.Bot().Request(
		tgbotapi.NewCallback(
			ctx.Update().CallbackQuery.ID,
			fmt.Sprintf("user %s added to watchlist", selectedUser),
		),
	)

	return this.Index(ctx)
}

func (this *WatchlistController) RemoveUsersIndex(
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
					fmt.Sprintf("/watchlist/a/delete/users/%s", user),
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

func (this *WatchlistController) RemoveUser(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	selectedUser := ctx.Params().ByName("user")
	log.Printf("selected user (%s)", selectedUser)

	ctx.Redirect("/watchlist")

	ctx.Bot().Request(
		tgbotapi.NewCallback(
			ctx.Update().CallbackQuery.ID,
			fmt.Sprintf("user %s removed from watchlist", selectedUser),
		),
	)

	return this.Index(ctx)
}

var _ tgool.Controller = (*WatchlistController)(nil)
