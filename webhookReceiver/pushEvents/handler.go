package pushEvents

import (
	botTypes "GitlabTgBot/botTypes"
	gitlabInfo "GitlabTgBot/gitlabInformation"
	messageReceiver "GitlabTgBot/messageReceiver"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreatePushHook(payload gitlab.PushEventPayload) *botTypes.WebhookInformation {
	return &botTypes.WebhookInformation{
		Type:      botTypes.Push,
		Message:   CreateMessage(payload),
		ProjectID: int(payload.ProjectID),
		Author:    payload.UserUsername,
	}
}

func Handle(payload gitlab.PushEventPayload, usersRelatedToThisHook []gitlabInfo.ProjectUser) {
	for i := range usersRelatedToThisHook {
		user := messageReceiver.GetUserWithGitlabUsername(usersRelatedToThisHook[i].Username)
		if user != nil {
			messageReceiver.CreateWebhookMessage(CreatePushHook(payload), user.GitlabUsername)
		}
	}
}
