package registration

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	gitlabInformation "GitlabTgBot/gitlabInformation"

	extensions "GitlabTgBot/extensions"

	botTypes "GitlabTgBot/botTypes"

	botDB "GitlabTgBot/db"
)

const (
	// Check number of %s and args in Exec! SQLite doesn't support boolean types, so it represents like integer, 0 = false, 1 = true
	// Table: |ChatID|UserName|RegistrationDate|Token|UserRole|GitlabUsername|Push|TagPush|Comments|ConfComments|Issues|ConfIssues|MR|Job|Pipeline|Wiki|Deployment|FeatureFlag|Release|
	insertUserIntoDBCommand    string = "INSERT INTO $0 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	successRegistrationMessage string = "Welcome %s %s!\nEnter /help for more information."
)

//	This map is important for registration logic. If user entered right token,
//
// his chatID adds here, and then user skips first registration stepand gets to
// second step.
var usersWithCorrectTokens map[int64]string = make(map[int64]string)

func RegisterNewUser(message *tgbotapi.Message, gitlabToken string) *botTypes.User {
	token, err := extensions.Encrypt([]byte(gitlabToken))
	if err != nil {
		log.Panic(err)
	}

	userName := message.Chat.FirstName + message.Chat.LastName
	_, err = botDB.GetBotDB().Exec(
		insertUserIntoDBCommand,
		botDB.UsersTableName,
		strconv.FormatInt(message.Chat.ID, 10),
		userName,
		string(token),
		botTypes.DefaultUserRole,
		message.Text,
		true,
		true,
		true,
		true,
		true,
		botTypes.DefaultActivityTimeLimit,
	)
	if err != nil {
		log.Panic(err)
	}

	return &botTypes.User{
		ChatID:             message.Chat.ID,
		UserName:           userName,
		Token:              gitlabToken,
		UserRole:           botTypes.DefaultUserRole,
		GitlabUsername:     message.Text,
		PushOption:         1,
		MergeRequestOption: 1,
		CommentsOption:     1,
		PipelineOption:     1,
		TagOption:          1,
		ActivityTimeLimit:  botTypes.DefaultActivityTimeLimit,
	}
}

// Here the username enered by bot user is searched in the list of gitlab users
func IsUsernameCorrect(username string, token string) bool {
	// get actual userlist
	gitlabUsers, err := gitlabInformation.GetUserList(token)
	if err != nil {
		log.Print(err)
		return false
	}

	for i := range gitlabUsers {
		if gitlabUsers[i].Username == username {
			return true
		}
	}

	return false
}

// Two steps registration. At first step new user enters private token from gitlab. If token is correct,
// user gets to 2nd step, where he must enter right gitlab username. If it's ok, user saves in
// DB and method returns true for updating bot users list in message handler
func CheckInfo(message *tgbotapi.Message) (*botTypes.User, string) {
	token, firstStepComplited := usersWithCorrectTokens[message.Chat.ID]
	if firstStepComplited {
		if IsUsernameCorrect(message.Text, token) {
			newUser := RegisterNewUser(message, token)
			return newUser, fmt.Sprintf(successRegistrationMessage, message.Chat.FirstName, message.Chat.LastName)
		} else {
			return nil, "Invalid username"
		}
	}

	if extensions.IsTokenCorrect(message.Text) {
		usersWithCorrectTokens[message.Chat.ID] = message.Text
		return nil, "Enter your gitlab username"
	}

	return nil, "Invalid token"
}
