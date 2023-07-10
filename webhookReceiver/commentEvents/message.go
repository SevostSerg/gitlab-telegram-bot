package commentEvents

import (
	botTypes "GitlabTgBot/botTypes"
	"fmt"
	"strconv"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreateMessage(info gitlab.CommentEventPayload) string {
	return fmt.Sprintf("Text: %s\nLink: %s", info.ObjectAttributes.Description, info.ObjectAttributes.URL)
}

func ModifyCommentMessage(webhook *botTypes.WebhookInformation, key Key) *botTypes.WebhookInformation {
	var commentWord string
	if len(comments[key]) == 1 {
		commentWord = "comment"
	} else {
		commentWord = "comments"
	}

	addition := fmt.Sprintf("User %s added %s %s to %s! Last:\n", key.Author, strconv.Itoa(len(comments[key])), commentWord, key.Project)
	webhook.Message = addition + webhook.Message
	return webhook
}
