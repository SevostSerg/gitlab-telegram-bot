package gitlabInformation

import (
	config "GitlabTgBot/configuration"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type ProjectUser struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Username    string      `json:"username"`
	State       string      `json:"state"`
	AvatarURL   interface{} `json:"avatar_url"`
	WebURL      string      `json:"web_url"`
	AccessLevel int64       `json:"access_level"`
	CreatedAt   string      `json:"created_at"`
	ExpiresAt   interface{} `json:"expires_at"`
}

const (
	tokenAddition    = "?private_token="
	membersLinkPart0 = "/api/v4/projects/"
	membersLinkPart1 = "/members/all"
)

// This method is necessary for webhooks logic. Using project ID it gets all users in project.
func GetUsersInProject(id int) ([]ProjectUser, error) {
	token := config.GetConfigInstance().GitlabAccessToken

	URL := config.GetConfigInstance().GitlabURL
	response, err := http.Get(URL + membersLinkPart0 + strconv.Itoa(id) + membersLinkPart1 + tokenAddition + token)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var projectUsers []ProjectUser
	err = json.Unmarshal(contents, &projectUsers)
	if err != nil {
		return nil, err
	}

	return projectUsers, nil
}
