package controllers

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/tgool"
)

type BilakhController struct{}

func (this *BilakhController) AddRoutes(b *tgool.RouterBuilder) {
	b.SetPrefixRoute("/start").
		AddMethod("", "Index").
		SetPrefixRoute("/bilakh").
		AddMethod("", "Index").
		AddMethod("/accept-your-fate", "AcceptYourFate").
		AddMethod("/yes-double-it", "YesDoubleIt").
		AddMethod("/who-sent-this", "WhoSentThis").
		AddMethod("/send-me-one-more-then", "SendMeOneMoreThen").
		AddMethod("/return-the-bilakh", "ReturnTheBilakh").
		AddMethod("/show-me-more", "ShowMeMore").
		AddMethod("/bilakh-with-fireworks", "BilakhWithFireworks").
		AddMethod("/bilakh-with-dramatic-music", "BilakhWithDramaticMusic").
		AddMethod("/bilakh-with-screaming-goat", "BialkhWithScreamingGoat").
		AddMethod("/bilakh-in-slow-motion", "BilakhInSlowMotion")
}

func (this *BilakhController) Index(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `Surprise! 🎁

	You’ve received a certified Bilakh™ from the universe. 👍`

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🫡 Accept your fate",
				"/bilakh/accept-your-fate",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"🤨 Who sent this?!",
				"/bilakh/who-sent-this",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔁 Return the Bilakh",
				"/bilakh/return-the-bilakh",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"💥 Show me more!",
				"/bilakh/show-me-more",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) AcceptYourFate(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `Wise choice. Resistance is futile. 👏`

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔁 Yes, double it",
				"/bilakh/yes-double-it",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"🙈 No thanks",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) YesDoubleIt(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := "Here's a double Bilakh for you honey. 👍👍❤️"

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🥰 Appretiated",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) WhoSentThis(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `Sorry, that info is classified. 🕵️`

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"👍 Send me one more then",
				"/bilakh/send-me-one-more-then",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"😭 Don't embarrase me anymore!",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) SendMeOneMoreThen(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `Here's one more Bilakh for you dear. ❤️👍`

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"❤️ Thanks",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) ReturnTheBilakh(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := "Sorry, this bilakh is non-refundable — it's handcrafted just for you. 🎁"

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"❤️🥹 Thanks then",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) ShowMeMore(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `You’re now entering Bilakh Heaven. 🥳`

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔥 Bilakh with fireworks",
				"/bilakh/bilakh-with-fireworks",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🎵 Bilakh with dramatic music",
				"/bilakh/bilakh-with-dramatic-music",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"😱 Bilakh with screaming goat",
				"/bilakh/bilakh-with-screaming-goat",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🐢 Bilakh in slow motion",
				"/bilakh/bilakh-in-slow-motion",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) BilakhWithFireworks(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `Here's a Bilakh with fireworks for you. 👍🔥🧨🎆`

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔥😍 Coooool",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) BilakhWithDramaticMusic(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := `Here's a Bilakh with some music for you. 👍🎵🎹🎷🎸🎺`

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"☺️😍 I love it",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) BialkhWithScreamingGoat(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt := "Here's a Bilakh with an screaming goat for you. 👍🐐😱"

	msg := tgbotapi.NewMessage(chatId, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🐑 Beeeee",
				"/bilakh",
			),
		),
	)

	return msg, nil
}

func (this *BilakhController) BilakhInSlowMotion(
	ctx tgool.Context,
) (tgbotapi.Chattable, error) {
	chatId := ctx.GetChatId()

	txt1 := "Here's a Biiii..."
	msg1 := tgbotapi.NewMessage(chatId, txt1)
	ctx.Bot().Send(msg1)

	time.Sleep(time.Second * 2)
	txt2 := "...llllaaaakhhhhh for you. 👍🐢"
	msg2 := tgbotapi.NewMessage(chatId, txt2)

	msg2.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🐌 Than... ...kksss",
				"/bilakh",
			),
		),
	)

	return msg2, nil
}

var _ tgool.Controller = (*BilakhController)(nil)
