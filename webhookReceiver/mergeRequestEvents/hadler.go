package mergeRequestEvents

import (
	botTypes "GitlabTgBot/botTypes"
	gitlabInfo "GitlabTgBot/gitlabInformation"
	messageReceiver "GitlabTgBot/messageReceiver"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreateMRHook(payload gitlab.MergeRequestEventPayload) *botTypes.WebhookInformation {
	return &botTypes.WebhookInformation{
		Type:      botTypes.MR,
		Message:   CreateMessage(payload),
		ProjectID: int(payload.Project.ID),
		Author:    payload.User.UserName,
	}
}

func Handle(payload gitlab.MergeRequestEventPayload, usersRelatedToThisHook []gitlabInfo.ProjectUser) {
	for i := range usersRelatedToThisHook {
		user := messageReceiver.GetUserWithGitlabUsername(usersRelatedToThisHook[i].Username)
		if user != nil {
			messageReceiver.CreateWebhookMessage(CreateMRHook(payload), user.GitlabUsername)
		}
	}
}

func IsMerged(payload gitlab.MergeRequestEventPayload) bool {
	return payload.ObjectAttributes.State == "merged"
}
