package botTypes

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserMessage struct {
	User                  *User
	Update                tgbotapi.Update
	Type                  TypeOfMessage
	LoopStatus            LoopStatus
	ChatID                int64
	RegistrationCompleted bool
}

type TypeOfMessage int

const (
	Unknown TypeOfMessage = iota
	Help
	InformationAboutUser
	WebhookStatusChange
	UserListRequest
	ProjectListRequest
	WebhookURLInformation
	NoteTimeLimitChange
	TokenChange
	KeyboardRequest
	GeraldMode
	StatsRequest
)

type LoopStatus int

const (
	NotALoopMessage LoopStatus = iota
	InProgress
	Completed
)

func AnalyzeMessage(update tgbotapi.Update, userList map[int64]*User, usersInLoop map[int64]*UserMessage) *UserMessage {
	loopStatus, MessageType := IdentifyTypeAndLoopStatus(update, usersInLoop)
	user, userRegistred := userList[update.Message.Chat.ID]
	return &UserMessage{
		User:                  user,
		Update:                update,
		Type:                  MessageType,
		LoopStatus:            loopStatus,
		ChatID:                update.Message.Chat.ID,
		RegistrationCompleted: userRegistred,
	}
}

func IdentifyMessageType(update tgbotapi.Update) TypeOfMessage {
	switch update.Message.Text {
	case "Help", "/help", "How to add notifications", "/start":
		return Help

	case "My Information":
		return InformationAboutUser

	case "Set time limit for notifications":
		return NoteTimeLimitChange

	case "Disable push",
		"Enable push",
		"Disable MR",
		"Enable MR",
		"Disable notes",
		"Enable notes",
		"Disable pipelines",
		"Enable pipelines",
		"Disable tags",
		"Enable tags":
		return WebhookStatusChange

	case "Users":
		return UserListRequest

	case "Projects":
		return ProjectListRequest

	case "Webhook Url":
		return WebhookURLInformation

	case "GeraldMode":
		return GeraldMode

	case "Change access token":
		return TokenChange

	case "Stats",
		"/stats":
		return StatsRequest

	case "/showKeyboard":
		return KeyboardRequest

	default:
		return Unknown
	}
}

func IdentifyTypeAndLoopStatus(update tgbotapi.Update, usersInLoop map[int64]*UserMessage) (LoopStatus, TypeOfMessage) {
	if message, contains := usersInLoop[update.Message.Chat.ID]; contains {
		return InProgress, message.Type
	} else {
		typeOfMessage := IdentifyMessageType(update)
		switch typeOfMessage {
		case NoteTimeLimitChange, TokenChange, GeraldMode:
			return InProgress, typeOfMessage

		default:
			return NotALoopMessage, typeOfMessage
		}
	}
}
