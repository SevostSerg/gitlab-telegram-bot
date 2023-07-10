package userProperties

import (
	botTypes "GitlabTgBot/botTypes"
	botDB "GitlabTgBot/db"
	extensions "GitlabTgBot/extensions"
	"fmt"
	"strconv"
	"strings"
)

const (
	dbTrue         int    = 1
	dbFalse        int    = 0
	disableMessage string = "Option %s successfully disabled"
	enableMessage  string = "Option %s successfully enabled"
)

// Here we can enable/disable every webhook notification status. If user disable or enable anything,
// method returns user with changed webhook parmeter and starts new goroutine,
// which updates this info in DB.
func ChangeWebhookNotifStatus(message *botTypes.UserMessage) string {
	// Commnd must be same as "disable/enable webhook"
	// For example: "disable push" or "enable mr"
	splittedCommand := strings.Split(message.Update.Message.Text, " ")
	command, option := splittedCommand[0], splittedCommand[1]
	return CheckOptionStatus(command, option, message.User)
}

// Check for is this option already has same status as user wants to change
func IsOptionAlreadyEqual(command string, isNotifEnabled bool) (string, bool) {
	if command == "Disable" && !isNotifEnabled {
		return "This option is already disabled", true
	}

	if command == "Enable" && isNotifEnabled {
		return "This option is already enabled", true
	}

	return "", false
}

func CheckOptionStatus(command string, option string, user *botTypes.User) string {
	webhookType, currentStatus := GetWebhookTypeAndStatus(user, option)
	alreadyEqualMessage, isEqual := IsOptionAlreadyEqual(command, currentStatus)
	if isEqual {
		return alreadyEqualMessage
	}

	var userOption *int
	var changingOptionName string
	switch webhookType {
	case botTypes.Push:
		userOption = &user.PushOption
		changingOptionName = botDB.PushOptionColumn

	case botTypes.MR:
		userOption = &user.MergeRequestOption
		changingOptionName = botDB.MergeRequestOptionColumn

	case botTypes.Comment:
		userOption = &user.CommentsOption
		changingOptionName = botDB.CommentsOptionColumn

	case botTypes.Pipeline:
		userOption = &user.PipelineOption
		changingOptionName = botDB.PipelineOptionColumn

	case botTypes.Tag:
		userOption = &user.TagOption
		changingOptionName = botDB.TagOptionColumn

	default:
		userOption = nil
	}

	if command == "Disable" {
		*userOption = dbFalse
		botDB.UpdateValInDB(user.ChatID, changingOptionName, strconv.Itoa(dbFalse), botDB.UsersTableName)
		return fmt.Sprintf(disableMessage, option)
	}

	if command == "Enable" {
		*userOption = dbTrue
		botDB.UpdateValInDB(user.ChatID, changingOptionName, strconv.Itoa(dbTrue), botDB.UsersTableName)
		return fmt.Sprintf(enableMessage, option)
	}

	return fmt.Sprintf("Unable command %s!", option)
}

func IsThisNotifAllowed(user *botTypes.User, webhook *botTypes.WebhookInformation) bool {
	switch webhook.Type {
	case botTypes.Push:
		return user.PushOption == dbTrue

	case botTypes.MR:
		return user.MergeRequestOption == dbTrue

	case botTypes.Comment:
		return user.CommentsOption == dbTrue

	case botTypes.Tag:
		return user.TagOption == dbTrue

	case botTypes.Pipeline:
		return user.PipelineOption == dbTrue

	default:
		return false
	}
}

func IsUserAuthorOfThisActivity(user *botTypes.User, webhook *botTypes.WebhookInformation) bool {
	return user.GitlabUsername == webhook.Author
}

func DoesUserHaveAccess(token string) bool {
	return extensions.IsTokenCorrect(token)
}

func GetWebhookTypeAndStatus(user *botTypes.User, option string) (botTypes.WebhookType, bool) {
	switch option {
	case "push":
		return botTypes.Push, user.PushOption == dbTrue

	case "MR":
		return botTypes.MR, user.MergeRequestOption == dbTrue

	case "notes":
		return botTypes.Comment, user.CommentsOption == dbTrue

	case "pipelines":
		return botTypes.Pipeline, user.PipelineOption == dbTrue

	case "tags":
		return botTypes.Tag, user.TagOption == dbTrue
	}

	return botTypes.UnknownHook, false
}
