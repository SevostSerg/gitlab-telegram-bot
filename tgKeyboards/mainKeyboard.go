package tgKeyboards

import (
	botTypes "GitlabTgBot/botTypes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CreateMainKeyboard(user *botTypes.User) tgbotapi.ReplyKeyboardMarkup {
	if user != nil && user.UserRole == "Admin" {
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Help"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("My Information"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Set time limit for notifications"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("GeraldMode"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Change access token"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Webhook Url"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Disable push"), tgbotapi.NewKeyboardButton("Enable push"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Disable MR"), tgbotapi.NewKeyboardButton("Enable MR"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Disable notes"), tgbotapi.NewKeyboardButton("Enable notes"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Disable pipelines"), tgbotapi.NewKeyboardButton("Enable pipelines"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Disable tags"), tgbotapi.NewKeyboardButton("Enable tags"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Users"), tgbotapi.NewKeyboardButton("Projects"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Stats"),
			),
		)
	}

	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Help"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("My Information"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Set time limit for notifications"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("GeraldMode"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Change access token"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Webhook Url"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Disable push"), tgbotapi.NewKeyboardButton("Enable push"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Disable MR"), tgbotapi.NewKeyboardButton("Enable MR"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Disable notes"), tgbotapi.NewKeyboardButton("Enable notes"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Users"), tgbotapi.NewKeyboardButton("Projects"),
		),
	)
}

func ShowMainKeyboard(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Main keyboard")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
