package messageReceiver

import (
	"fmt"
	"log"
	"strconv"

	botTypes "GitlabTgBot/botTypes"
	extensions "GitlabTgBot/extensions"
	registration "GitlabTgBot/registration"
	userProperties "GitlabTgBot/userProperties"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var tgBot *tgbotapi.BotAPI
var updatesChan tgbotapi.UpdatesChannel

// Key: tg bot chat ID
var usersChatIDKey map[int64]*botTypes.User

// Key: gitlab username
var usersGitlabUsernameKey map[string]*botTypes.User

var usersInLoop map[int64]*botTypes.UserMessage = make(map[int64]*botTypes.UserMessage)

func Start(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	tgBot = bot
	updatesChan = updates
	usersChatIDKey, usersGitlabUsernameKey = botTypes.GetUserListFromDB()
	StartReceivingMessages()
}

func StartReceivingMessages() {
	for update := range updatesChan {
		if update.Message == nil { // ignore any non-Message Updates and /start
			continue
		}

		if extensions.IsUserBannedForSpam(update.Message.Chat.ID, tgBot) {
			continue
		}

		RecognizeUserMessage(botTypes.AnalyzeMessage(update, usersChatIDKey, usersInLoop))
	}
}

func CreateWebhookMessage(webhook *botTypes.WebhookInformation, gitlabUserName string) {
	if user, isUserRegistred := usersChatIDKey[GetUserWithGitlabUsername(gitlabUserName).ChatID]; isUserRegistred && CanUserReceiveMessage(user, webhook) {
		tgBot.Send(tgbotapi.NewMessage(user.ChatID, webhook.Message))
		log.Printf("Webhook event(%s): Message successfully sended to %s!", strconv.Itoa(int(webhook.Type)), gitlabUserName)
	} else {
		var cond string
		if usersGitlabUsernameKey[gitlabUserName].GitlabUsername == webhook.Author {
			cond = "User is author."
		} else {
			cond = "User disabled this option."
		}

		log.Printf("Webhook event: Cannot send this message to user %s! %s", gitlabUserName, cond)
	}
}

func CanUserReceiveMessage(user *botTypes.User, webhook *botTypes.WebhookInformation) bool {
	if webhook.Type == botTypes.Pipeline {
		return userProperties.IsThisNotifAllowed(user, webhook) && userProperties.DoesUserHaveAccess(user.Token)
	}

	return userProperties.IsThisNotifAllowed(user, webhook) && !userProperties.IsUserAuthorOfThisActivity(user, webhook) &&
		userProperties.DoesUserHaveAccess(user.Token)
}

func RecognizeUserMessage(message *botTypes.UserMessage) {
	response, err := ReadUserCommand(message)
	if err != nil {
		log.Print(err)
	}

	if message.RegistrationCompleted {
		for i := range response {
			tgBot.Send(tgbotapi.NewMessage(message.Update.Message.Chat.ID, response[i]))
			log.Printf("%s: %s", message.User.GitlabUsername, message.Update.Message.Text)
		}
	}
}

func RegisterUser(update tgbotapi.Update) string {
	newUser, responce := registration.CheckInfo(update.Message)
	tgBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, responce))
	if newUser != nil {
		usersChatIDKey[newUser.ChatID] = newUser
		usersGitlabUsernameKey[newUser.GitlabUsername] = newUser
		log.Printf("Registration: User %s was registered", update.Message.Chat.UserName)
	}

	return responce
}

func GetUserWithGitlabUsername(gitlabUsername string) *botTypes.User {
	res, contains := usersGitlabUsernameKey[gitlabUsername]
	if !contains {
		return nil
	}

	user, _ := usersChatIDKey[res.ChatID]
	return user
}

func SendToEveryone(payload tgbotapi.Update) {
	for chatID := range usersChatIDKey {
		if payload.Message.Chat.ID == chatID {
			continue
		}

		tgBot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("%s:\n%s", payload.Message.From.FirstName, payload.Message.Text)))
		if payload.Message.Photo != nil {
			photo := *payload.Message.Photo
			tgBot.Send(tgbotapi.NewPhotoShare(chatID, photo[len(photo)-1].FileID))
		}
	}
}
