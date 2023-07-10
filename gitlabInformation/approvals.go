package gitlabInformation

import (
	config "GitlabTgBot/configuration"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Approvals struct {
	UserHasApproved bool         `json:"user_has_aproved"`
	UserCanApprove  bool         `json:"user_can_approve"`
	Approved        bool         `json:"approved"`
	ApprovedBy      []ApprovedBy `json:"approved_by"`
}

type ApprovedBy struct {
	User GitlabUser `json:"user"`
}

func GetApprovalsList(accessToken string, projectID int64, mergeRequestID int64) (map[string]GitlabUser, error) {
	GitlabURL := config.GetConfigInstance().GitlabURL
	response, err := http.Get(fmt.Sprintf("%s/api/v4/projects/%s/merge_requests/%s/approvals?private_token=%s", GitlabURL, strconv.FormatInt(projectID, 10), strconv.FormatInt(mergeRequestID, 10), accessToken))
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var approvalsList Approvals
	json.Unmarshal(contents, &approvalsList)
	approvalsMap := make(map[string]GitlabUser)
	for i := range approvalsList.ApprovedBy {
		approvalsMap[approvalsList.ApprovedBy[i].User.Username] = approvalsList.ApprovedBy[i].User
	}

	return approvalsMap, nil
}
