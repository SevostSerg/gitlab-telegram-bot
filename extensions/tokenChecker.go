package extensions

import (
	config "GitlabTgBot/configuration"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const checkRequest string = "/api/v4/users?private_token="

func IsTokenCorrect(token string) bool {
	if strings.Contains(token, " ") {
		return false
	}

	URL := config.GetConfigInstance().GitlabURL
	resp, err := http.Get(URL + checkRequest + token)
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var error401 = regexp.MustCompile(`"message":"401 Unauthorized"`)
	var invalidTokenError = regexp.MustCompile(`"error":"invalid_token"`)
	return !error401.MatchString(string(content)) && !invalidTokenError.MatchString(string(content))
}
