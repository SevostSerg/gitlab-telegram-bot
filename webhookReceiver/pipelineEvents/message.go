package pipelineEvents

import (
	"fmt"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreateMessage(info gitlab.PipelineEventPayload) string {
	if info.ObjectAttributes.Tag {
		return fmt.Sprintf("Pipeline status: %s\nLink for details: %s", info.ObjectAttributes.Status, info.Commit.URL)
	}

	return fmt.Sprintf("Your pipeline status: %s\nLink for details: %s", info.ObjectAttributes.Status, info.Commit.URL)
}
