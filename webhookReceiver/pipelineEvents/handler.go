package pipelineEvents

import (
	botTypes "GitlabTgBot/botTypes"
	gitlabInfo "GitlabTgBot/gitlabInformation"
	messageReceiver "GitlabTgBot/messageReceiver"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreatePipelineHook(payload gitlab.PipelineEventPayload) *botTypes.WebhookInformation {
	return &botTypes.WebhookInformation{
		Type:      botTypes.Pipeline,
		Message:   CreateMessage(payload),
		ProjectID: int(payload.ObjectAttributes.ID),
		Author:    payload.User.UserName,
	}
}

func Handle(payload gitlab.PipelineEventPayload, usersRelatedToThisHook []gitlabInfo.ProjectUser) {
	if payload.ObjectAttributes.Status == "pending" || payload.ObjectAttributes.Status == "running" {
		return
	}

	if payload.ObjectAttributes.Tag {
		for i := range usersRelatedToThisHook {
			user := messageReceiver.GetUserWithGitlabUsername(usersRelatedToThisHook[i].Username)
			if user != nil {
				messageReceiver.CreateWebhookMessage(CreatePipelineHook(payload), user.GitlabUsername)
			}
		}
	}

	messageReceiver.CreateWebhookMessage(CreatePipelineHook(payload), payload.User.UserName)

}
