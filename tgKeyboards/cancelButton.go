package tgKeyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func CreateCancelButtom() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cancel"),
		),
	)
}

func ShowCancelButton(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Operation can be canceled with \"Cancel\" button")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
