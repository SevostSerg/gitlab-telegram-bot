package webhookReceiver

import (
	"GitlabTgBot/gitlabInformation"
	"GitlabTgBot/webhookReceiver/commentEvents"
	"GitlabTgBot/webhookReceiver/mergeRequestEvents"
	"GitlabTgBot/webhookReceiver/mrRoulette"
	"GitlabTgBot/webhookReceiver/pipelineEvents"
	"GitlabTgBot/webhookReceiver/pushEvents"
	"GitlabTgBot/webhookReceiver/tagEvents"

	"log"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func ParseWebhook(payload interface{}) {
	if isItSilentTime() {
		go handleHookSilent(payload)
		return
	}

	switch payload := payload.(type) {
	case gitlab.PushEventPayload:
		pushEvents.Handle(payload, GetProjectUsers(int(payload.Project.ID)))

	case gitlab.CommentEventPayload:
		mrRoulette.ChangeReviewerIfNecessary(&mrRoulette.MRRoulettePayload{
			ProjectID:      payload.ProjectID,
			MergeRequestID: payload.MergeRequest.IID,
			GitlabAction:   payload.ObjectAttributes.Action,
			Description:    payload.ObjectAttributes.Description,
			HookType:       mrRoulette.CommentHook,
		})

		commentEvents.Handle(payload, GetProjectUsers(int(payload.Project.ID)))

	case gitlab.MergeRequestEventPayload:
		mrRoulette.ChangeReviewerIfNecessary(&mrRoulette.MRRoulettePayload{
			ProjectID:      payload.Project.ID,
			MergeRequestID: payload.ObjectAttributes.IID,
			GitlabAction:   payload.ObjectAttributes.Action,
			Description:    "",
			HookType:       mrRoulette.MRHook,
		})
		if mergeRequestEvents.IsMerged(payload) {
			return
		}

		mergeRequestEvents.Handle(payload, GetProjectUsers(int(payload.Project.ID)))

	case gitlab.TagEventPayload:
		tagEvents.Handle(payload, GetProjectUsers(int(payload.Project.ID)))

	case gitlab.PipelineEventPayload:
		pipelineEvents.Handle(payload, GetProjectUsers(int(payload.Project.ID)))

	default:
		log.Print("Unknown payload!")
	}
}

func GetProjectUsers(projectId int) []gitlabInformation.ProjectUser {
	projectUsers, err := gitlabInformation.GetUsersInProject(projectId)
	if err != nil {
		log.Panicf("Unable to get user list from gitlab! %s", err.Error())
	}

	return projectUsers
}
