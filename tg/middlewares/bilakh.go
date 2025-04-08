package middlewares

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/tg/controllers"
	"github.com/thekhanj/csdmpro/tg/repo"
	"github.com/thekhanj/tgool"
)

type BilakhMiddleware struct {
	repo *repo.BilakhRepo
}

func (this *BilakhMiddleware) Handle(
	ctx tgool.Context, next func(),
) tgbotapi.Chattable {
	chatId := ctx.GetChatId()

	isBilakhed, err := this.repo.IsBilakhed(chatId)
	if err != nil {
		log.Println(err)
		return nil
	}
	if !isBilakhed {
		next()
		return nil
	}

	m := tgool.NewControllerMiddleware(&controllers.BilakhController{})

	ret := m.Handle(ctx, func() {})
	return ret
}

func NewBilakhMiddleware(bilakhRepo *repo.BilakhRepo) *BilakhMiddleware {
	return &BilakhMiddleware{bilakhRepo}
}

var _ tgool.Middleware = (*BilakhMiddleware)(nil)
