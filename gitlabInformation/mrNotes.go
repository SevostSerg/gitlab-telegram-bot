package gitlabInformation

import (
	config "GitlabTgBot/configuration"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Notes []Note

type Note struct {
	ID              int64           `json:"id"`
	Type            interface{}     `json:"type"`
	Body            string          `json:"body"`
	Attachment      interface{}     `json:"attachment"`
	Author          GitlabUser      `json:"author"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	System          bool            `json:"system"`
	NoteableID      int64           `json:"noteable_id"`
	NoteableType    string          `json:"noteable_type"`
	Resolvable      bool            `json:"resolvable"`
	Confidential    bool            `json:"confidential"`
	NoteableIid     int64           `json:"noteable_iid"`
	CommandsChanges CommandsChanges `json:"commands_changes"`
}

type CommandsChanges struct {
}

func GetMRNotesList(accessToken string, projectID int64, mergeRequestID int64) ([]Note, error) {
	GitlabURL := config.GetConfigInstance().GitlabURL
	uri := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests/%s/notes?private_token=%s&per_page=300", GitlabURL, strconv.FormatInt(projectID, 10), strconv.FormatInt(mergeRequestID, 10), accessToken)
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var notes []Note
	json.Unmarshal(contents, &notes)
	return notes, nil
}
