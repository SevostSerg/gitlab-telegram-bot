package commentEvents

import (
	botTypes "GitlabTgBot/botTypes"
	gitlabInfo "GitlabTgBot/gitlabInformation"
	messageReceiver "GitlabTgBot/messageReceiver"
	"time"

	"github.com/go-playground/webhooks/v6/gitlab"
)

type Key struct {
	Project string
	Author  string
	TGName  string
}

var comments map[Key][]botTypes.WebhookInformation = make(map[Key][]botTypes.WebhookInformation)

func CreateCommentHook(payload gitlab.CommentEventPayload) botTypes.WebhookInformation {
	var hook botTypes.WebhookInformation
	hook.ProjectID = int(payload.ProjectID)
	hook.Message = CreateMessage(payload)
	hook.Type = botTypes.Comment
	hook.Author = payload.User.UserName
	return hook
}

func Handle(payload gitlab.CommentEventPayload, usersRelatedToThisHook []gitlabInfo.ProjectUser) {
	for i := range usersRelatedToThisHook {
		userRelatedToThisHook := messageReceiver.GetUserWithGitlabUsername(usersRelatedToThisHook[i].Username)
		if userRelatedToThisHook != nil {
			key := CreateKey(payload.User.Name, userRelatedToThisHook.GitlabUsername, payload.Project.Name)
			webhook := CreateCommentHook(payload)
			// if it's the 1st comment for
			if _, contains := comments[key]; !contains {
				arr := make([]botTypes.WebhookInformation, 0)
				arr = append(arr, webhook)
				comments[key] = arr
				go ListenReceiver(key, *userRelatedToThisHook, webhook)
			} else {
				comments[key] = append(comments[key], webhook)
			}
		}
	}
}

func CreateKey(author string, toUser string, project string) Key {
	return Key{
		Project: project,
		Author:  author,
		TGName:  toUser,
	}
}

func ListenReceiver(key Key, user botTypes.User, webhook botTypes.WebhookInformation) {
	timer := *time.NewTimer(time.Minute * time.Duration(user.ActivityTimeLimit))
	<-timer.C
	hooksForThisUser := comments[key]
	messageReceiver.CreateWebhookMessage(ModifyCommentMessage(&hooksForThisUser[len(hooksForThisUser)-1], key), user.GitlabUsername)
	delete(comments, key)
}
