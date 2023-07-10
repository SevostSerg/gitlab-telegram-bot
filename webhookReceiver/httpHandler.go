package webhookReceiver

import (
	"log"

	"net/http"

	"github.com/go-playground/webhooks/v6/gitlab"

	config "GitlabTgBot/configuration"
)

const (
	path string = "/"
)

var (
	gitlabWebhook *gitlab.Webhook
)

func StartWebhookReceiving() {
	port := config.GetConfigInstance().Port
	gitlabWebhook = CreateWebhook()
	log.Print("Listen on port: " + port)
	http.HandleFunc(path, GetHook)
	http.ListenAndServe(":"+port, nil)
}

func GetHook(w http.ResponseWriter, r *http.Request) {
	payload, err := gitlabWebhook.Parse(r, gitlab.PushEvents, gitlab.TagEvents, gitlab.CommentEvents,
		gitlab.CommentEvents, gitlab.IssuesEvents, gitlab.ConfidentialIssuesEvents,
		gitlab.MergeRequestEvents, gitlab.JobEvents, gitlab.PipelineEvents,
		gitlab.WikiPageEvents)
	if err != nil {
		if err == gitlab.ErrEventNotFound {
			log.Print("Unknown event in http handle func")
			return
		}
	}

	ParseWebhook(payload)
}

func CreateWebhook() *gitlab.Webhook {
	webhookToken := config.GetConfigInstance().WebhookToken
	var webhook *gitlab.Webhook
	if webhookToken == "" {
		// without secret token
		instance, err := gitlab.New()
		if err != nil {
			log.Panic(err)
		}

		webhook = instance
	} else {
		instance, err := gitlab.New(gitlab.Options.Secret(webhookToken))
		if err != nil {
			log.Panic(err)
		}

		webhook = instance
	}

	return webhook
}
