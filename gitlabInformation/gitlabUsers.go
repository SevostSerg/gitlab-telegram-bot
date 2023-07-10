package gitlabInformation

import (
	config "GitlabTgBot/configuration"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GitlabUser struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	State     State   `json:"state"`
	AvatarURL *string `json:"avatar_url"`
	WebURL    string  `json:"web_url"`
}

type State string

// Users limit = 100
const (
	Active               State  = "active"
	usersRequestAddition string = "/api/v4/users?per_page=100&private_token="
)

func GetUserList(accessToken string) ([]GitlabUser, error) {
	URL := config.GetConfigInstance().GitlabURL
	response, err := http.Get(URL + usersRequestAddition + accessToken)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}

	var userList []GitlabUser
	json.Unmarshal(contents, &userList)
	return userList, nil
}
