package middlewares

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/tg/controllers"
	"github.com/thekhanj/tgool"
)

type BilakhMiddleware struct {
	bilakhChatIds map[int64]bool
}

func (this *BilakhMiddleware) Handle(
	ctx tgool.Context, next func(),
) tgbotapi.Chattable {
	chatId := ctx.GetChatId()

	if !this.bilakhChatIds[chatId] {
		next()
		return nil
	}

	m := tgool.NewControllerMiddleware(&controllers.BilakhController{})

	ret := m.Handle(ctx, func() {})
	if ret != nil {
		return ret
	}

	defaultMiddleWare := tgool.DefaultMiddleWare{}
	return defaultMiddleWare.Handle(ctx, next)
}

func NewBilakhMiddleware(bilakhChatIds []int64) *BilakhMiddleware {
	bilakhMap := make(map[int64]bool, 0)

	for _, chatId := range bilakhChatIds {
		bilakhMap[chatId] = true
	}

	return &BilakhMiddleware{bilakhChatIds: bilakhMap}
}

var _ tgool.Middleware = (*BilakhMiddleware)(nil)
