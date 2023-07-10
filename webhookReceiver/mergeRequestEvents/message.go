package mergeRequestEvents

import (
	"fmt"

	"github.com/go-playground/webhooks/v6/gitlab"
)

func CreateMessage(info gitlab.MergeRequestEventPayload) string {
	objAttr := info.ObjectAttributes
	switch objAttr.Action {
	case "open":
		return fmt.Sprintf("User %s opened a new MR in %s!\n%s into %s\nTitle: %s\nLink: %s",
			info.User.Name, info.Project.Name, objAttr.SourceBranch, objAttr.TargetBranch, objAttr.Title, objAttr.URL)

	case "close":
		return fmt.Sprintf("User %s closed MR in %s!\n%s into %s\nTitle: %s\nLink: %s",
			info.User.Name, info.Project.Name, objAttr.SourceBranch, objAttr.TargetBranch, objAttr.Title, objAttr.URL)

	case "approved":
		return fmt.Sprintf("User %s approved MR in %s!\n%s into %s\nTitle: %s\nLink: %s",
			info.User.Name, info.Project.Name, objAttr.SourceBranch, objAttr.TargetBranch, objAttr.Title, objAttr.URL)

	case "merge":
		return fmt.Sprintf("User %s merged MR in %s!\n%s into %s\nTitle: %s\nLink: %s",
			info.User.Name, info.Project.Name, objAttr.SourceBranch, objAttr.TargetBranch, objAttr.Title, objAttr.URL)

	case "update":
		return fmt.Sprintf("User %s updated MR in %s!\n%s into %s\nTitle: %s\nLink: %s",
			info.User.Name, info.Project.Name, objAttr.SourceBranch, objAttr.TargetBranch, objAttr.Title, objAttr.URL)

	case "unapproved":
		return fmt.Sprintf("User %s unapproved MR in %s!\n%s into %s\nTitle: %s\nLink: %s",
			info.User.Name, info.Project.Name, objAttr.SourceBranch, objAttr.TargetBranch, objAttr.Title, objAttr.URL)
	}

	return "Unknown action lol"
}
