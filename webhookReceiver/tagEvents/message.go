package tagEvents

import (
	"fmt"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreateMessage(info gitlab.TagEventPayload) string {
	return fmt.Sprintf("User %s created tag in %s!\nWaiting for pipeline result...\nProject URL: %s",
		info.UserName, info.Project.Name, info.Project.WebURL)
}
