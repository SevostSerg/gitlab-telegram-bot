package pushEvents

import (
	"fmt"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreateMessage(info gitlab.PushEventPayload) string {
	commit := info.Commits[len(info.Commits)-1]
	return fmt.Sprintf("User %s pushed to %s!\nCommit: %s\nLink: %s",
		info.UserName, info.Project.Name, commit.Title, commit.URL)
}
