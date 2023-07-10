package tagEvents

import (
	botTypes "GitlabTgBot/botTypes"
	gitlabInfo "GitlabTgBot/gitlabInformation"
	messageReceiver "GitlabTgBot/messageReceiver"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreateTagHook(payload gitlab.TagEventPayload) *botTypes.WebhookInformation {
	return &botTypes.WebhookInformation{
		Type:      botTypes.Tag,
		Message:   CreateMessage(payload),
		ProjectID: int(payload.ProjectID),
		Author:    payload.UserUsername,
	}
}

func Handle(payload gitlab.TagEventPayload, usersRelatedToThisHook []gitlabInfo.ProjectUser) {
	for i := range usersRelatedToThisHook {
		user := messageReceiver.GetUserWithGitlabUsername(usersRelatedToThisHook[i].Username)
		if user != nil {
			messageReceiver.CreateWebhookMessage(CreateTagHook(payload), user.GitlabUsername)
		}
	}
}
