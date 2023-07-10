package messageReceiver

import (
	"fmt"
	"log"
	"strconv"

	botTypes "GitlabTgBot/botTypes"
	config "GitlabTgBot/configuration"
	extensions "GitlabTgBot/extensions"
	gitlabInfo "GitlabTgBot/gitlabInformation"
	keyboards "GitlabTgBot/tgKeyboards"
	userProperties "GitlabTgBot/userProperties"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	helloMessage          string = "Hello %s %s! Use button 'Help' to get instructions. At first enter your private token."
	tokenHelp             string = "To get your private token, go to the profile settings in Gitlab, section \"Access Tokens\", check ReadAPI and copy it.\n"
	notifHelpMessage      string = "If the webhook for your project has not been created yet, create it in the project settings, allowing all events. If you don't have access to the settings, go to the maintainer of the project\n"
	helpMessage           string = tokenHelp + notifHelpMessage + showKeyboardMessage
	startMessage          string = "The first step is registration. To register, enter your private token from Gitlab. For more information press \"Help\" or enter /help"
	tokenExpiredMessage   string = "❌Wrong token. If your token is revoked or expired, please update your token!"
	unknownCommandMessage string = "Unknown command. Enter /help\n" + showKeyboardMessage
	showKeyboardMessage   string = "If you are registred and the keybord is hidden, enter /showKeyboard\n"
	geraldKeyword         string = "/gerald:"
	itemsPerMessage       int    = 20
)

// TODO: split it inside into different methods
// Processing of commands of registered users
// Terrible shit
func ReadUserCommand(message *botTypes.UserMessage) ([]string, error) {
	var response []string
	if message.Type == botTypes.Help {
		response[0], _ = HelpUser(message)
		return response, nil
	}

	if !message.RegistrationCompleted {
		response = append(response, RegisterUser(message.Update))
		return response, nil
	}

	if message.Type == botTypes.TokenChange {
		response = append(response, ChangeUserToken(message))
		return response, nil
	}

	if !extensions.IsTokenCorrect(message.User.Token) {
		response = append(response, "Your token is revoked or expired")
		return response, nil
	}

	// features for valid users
	switch message.Type {
	case botTypes.UserListRequest:
		response = GetUserListMessage(message.User.Token)

	case botTypes.ProjectListRequest:
		response = GetProjectsListMessage(message.User.Token)

	case botTypes.InformationAboutUser:
		response = append(response, GetUserInfo(message.User))

	case botTypes.KeyboardRequest:
		user, _ := botTypes.GetUserWithChatID(message.ChatID)
		keyboards.ShowMainKeyboard(tgBot, message.Update, keyboards.CreateMainKeyboard(user))
		response = append(response, "")

	case botTypes.WebhookURLInformation:
		response = append(response, SendWebhookURL())

	case botTypes.WebhookStatusChange:
		response = append(response, userProperties.ChangeWebhookNotifStatus(message))

	case botTypes.NoteTimeLimitChange:
		response = append(response, SetTimeLimit(message))

	case botTypes.GeraldMode:
		response = append(response, CreateGeraldMessage(message))

	default:
		response = append(response, unknownCommandMessage)
	}

	return response, nil
}

// Help for all users(not registred too)
func HelpUser(message *botTypes.UserMessage) (string, error) { // redo
	responce := unknownCommandMessage
	switch message.Update.Message.Text {
	case "/start":
		user, _ := botTypes.GetUserWithChatID(message.ChatID)
		keyboards.ShowMainKeyboard(tgBot, message.Update, keyboards.CreateMainKeyboard(user))
		responce = startMessage

	case "Help", "/help":
		responce = helpMessage

	case "How to add notifications":
		responce = notifHelpMessage
	}

	tgBot.Send(tgbotapi.NewMessage(message.ChatID, responce))
	return "", nil
}

func SendWebhookURL() string {
	return "Link to add webhooks for your projects:\n" + config.GetConfigInstance().WebhookURL +
		"\nAdd the webhook at project settings!\nTo display a list of commands, enter '/help'"
}

func CreateGeraldMessage(message *botTypes.UserMessage) string {
	if message.Update.Message.Text == "GeraldMode" {
		AddToLoop(message)
		return "Enter what you want to say"
	}

	if message.Update.Message.Text == "Cancel" {
		RemoveFromLoop(message)
		return "Gerald mode is disabled"
	}

	SendToEveryone(message.Update)
	return "sended"
}

func CreateRankRequestsKeyboard(message *botTypes.UserMessage) string {
	if message.Update.Message.Text == "Rank requests" {
		AddToLoop(message)
		return "Your requests:"
	}

	if message.Update.Message.Text == "Cancel" {
		RemoveFromLoop(message)
		return "Canceled"
	}

	return "FEATURE TEST"
}

func SetTimeLimit(message *botTypes.UserMessage) string {
	if message.Update.Message.Text == "Set time limit for notifications" {
		AddToLoop(message)
		return "Enter new time limit from 1 to 15 min"
	}

	if message.Update.Message.Text == "Cancel" {
		RemoveFromLoop(message)
		return fmt.Sprintf("Operation canceled. Current time limit(min): %s", strconv.Itoa(message.User.ActivityTimeLimit))
	}

	responce, changingCompleted := userProperties.ChangeTimeLimit(message.Update.Message.Text, message.User)
	if changingCompleted {
		user, _ := botTypes.GetUserWithChatID(message.ChatID)
		keyboards.ShowMainKeyboard(tgBot, message.Update, keyboards.CreateMainKeyboard(user))
		delete(usersInLoop, message.User.ChatID)
	}

	return responce
}

func ChangeUserToken(message *botTypes.UserMessage) string {
	if message.Update.Message.Text == "Change access token" {
		AddToLoop(message)
		return "Enter new token!"
	}

	if message.Update.Message.Text == "Cancel" {
		RemoveFromLoop(message)
		return fmt.Sprintf("Operation canceled. Your token is correct: %s", strconv.FormatBool(extensions.IsTokenCorrect(message.User.Token)))
	}

	response, newToken, operationCompleted := userProperties.UpdateToken(message.Update.Message.Text, message.User.ChatID)
	if operationCompleted {
		RemoveFromLoop(message)
		message.User.Token = newToken
		return "✅Token successfully updated!"
	}

	return response
}

func GetUserListMessage(accessToken string) []string {
	var response []string
	message := "Users:\n"
	gitlabUsersList, err := gitlabInfo.GetUserList(accessToken)
	if err != nil {
		log.Print(err)
		return nil
	}

	for i := range gitlabUsersList {
		if i%itemsPerMessage == 0 {
			response = append(response, message)
			message = ""
		}

		message += fmt.Sprintf("%d) %s\n%s\n", i+1, gitlabUsersList[i].Name, gitlabUsersList[i].WebURL)
	}

	response = append(response, message)
	return response
}

func GetProjectsListMessage(accessToken string) []string {
	var response []string
	message := "Projects:\n"
	gitlabProjectsList, err := gitlabInfo.GetProjectsList(accessToken)
	if err != nil {
		log.Print(err)
		return nil
	}

	for i := range gitlabProjectsList {
		if i%itemsPerMessage == 0 {
			response = append(response, message)
			message = ""
		}

		message += fmt.Sprintf("%d) %s\n%s\n", i+1, gitlabProjectsList[i].Name, gitlabProjectsList[i].WebURL)
	}

	response = append(response, message)
	return response
}

func GetUserInfo(user *botTypes.User) string {
	return fmt.Sprintf("Gitlab username: %s\nBot user role: %s\n~~~"+
		"\nPush hooks: %s\nMR hooks: %s\nComment hooks: %s\nPipeline hooks: %s\nTag hooks: %s\n~~~"+
		"\nTime limit for comments: %s minute(s)\n~~~",
		user.GitlabUsername,
		user.UserRole,
		extensions.GetOptionStatusString(user.PushOption),
		extensions.GetOptionStatusString(user.MergeRequestOption),
		extensions.GetOptionStatusString(user.CommentsOption),
		extensions.GetOptionStatusString(user.PipelineOption),
		extensions.GetOptionStatusString(user.TagOption),
		strconv.Itoa(user.ActivityTimeLimit))
}

func AddToLoop(message *botTypes.UserMessage) {
	usersInLoop[message.User.ChatID] = message
	keyboards.ShowCancelButton(tgBot, message.Update, keyboards.CreateCancelButtom())
}

func RemoveFromLoop(message *botTypes.UserMessage) {
	delete(usersInLoop, message.User.ChatID)
	user, _ := botTypes.GetUserWithChatID(message.ChatID)
	keyboards.ShowMainKeyboard(tgBot, message.Update, keyboards.CreateMainKeyboard(user))
}
