package mrRoulette

import (
	config "GitlabTgBot/configuration"
	gitlabInfo "GitlabTgBot/gitlabInformation"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RouletteRequestType int

const (
	RouletteWithRandomOrder RouletteRequestType = iota
	RouletteWithStrictOrder
	ReviewerChangeCommand
)

const (
	RequestWithStrictOrder string = "/mr-roulette <@"
	RequestWithRandOrder   string = "/mr-roulette rand <@"
	ReviewerChangeRequest  string = "/mr-roulette ask <@"
	HttpRequest            string = "%s/api/v4/projects/%s/merge_requests/%s?reviewer_ids=%s&private_token=%s"
)

func ChangeReviewer(payload *MRRoulettePayload, newReviewerID int) {
	client := &http.Client{}
	configInfo := config.GetConfigInstance()
	a := fmt.Sprintf(HttpRequest, configInfo.GitlabURL, strconv.FormatInt(payload.ProjectID, 10), strconv.FormatInt(payload.MergeRequestID, 10), strconv.Itoa(newReviewerID), configInfo.GitlabAccessToken)
	request, err := http.NewRequest(http.MethodPut, a, nil)
	if err != nil {
		log.Panic(err)
	}

	_, err = client.Do(request)
	if err != nil {
		log.Panic(err)
	}
}

func ChangeReviewerIfNecessary(payload *MRRoulettePayload) {
	if !IsRouletteRequest(payload) {
		return
	}

	var newReviewer *gitlabInfo.ProjectUser

	approvedUsers, err := gitlabInfo.GetApprovalsList(config.GetConfigInstance().GitlabAccessToken, payload.ProjectID, payload.MergeRequestID)
	if err != nil {
		log.Panic(err)
	}

	projectUsers, err := gitlabInfo.GetUsersInProject(int(payload.ProjectID))
	if err != nil {
		log.Panic(err)
	}

	switch payload.HookType {
	case CommentHook:
		newReviewer = ChangeOnCommentEvent(payload, approvedUsers, projectUsers)

	case MRHook:
		newReviewer = ChangeOnMrEvent(payload, approvedUsers, projectUsers)
	}

	if newReviewer != nil {
		ChangeReviewer(payload, int(newReviewer.ID))
	}
}

func ChangeOnCommentEvent(payload *MRRoulettePayload, approvedUsers map[string]gitlabInfo.GitlabUser, projectUsers []gitlabInfo.ProjectUser) *gitlabInfo.ProjectUser {
	user, err := GetUsersInRequest(payload.Description, ReviewerChangeCommand)
	if err != nil {
		log.Panic(err)
	}

	newReviewer, err := GetNewReviewer(user, ReviewerChangeCommand, approvedUsers, projectUsers)
	if err != nil {
		log.Panic(err)
	}

	return newReviewer
}

func ChangeOnMrEvent(payload *MRRoulettePayload, approvedUsers map[string]gitlabInfo.GitlabUser, projectUsers []gitlabInfo.ProjectUser) *gitlabInfo.ProjectUser {
	notes, err := gitlabInfo.GetMRNotesList(config.GetConfigInstance().GitlabAccessToken, payload.ProjectID, payload.MergeRequestID)
	if err != nil {
		log.Panic(err)
	}

	var requests []gitlabInfo.Note
	for i := range notes {
		if strings.HasPrefix(notes[i].Body, RequestWithRandOrder) || strings.HasPrefix(notes[i].Body, RequestWithStrictOrder) ||
			strings.HasPrefix(notes[i].Body, ReviewerChangeRequest) {
			requests = append(requests, notes[i])
		}
	}

	if len(requests) == 0 {
		return nil
	}

	actualNote := GetLastMRRouletteCommand(requests)
	requestType, err := IdentifyRequestType(actualNote.Body)
	if err != nil {
		log.Panic(err)
	}

	usersString, err := GetUsersInRequest(actualNote.Body, requestType)
	if err != nil {
		log.Panic(err)
	}

	newReviewer, err := GetNewReviewer(usersString, requestType, approvedUsers, projectUsers)
	if err != nil {
		log.Panic(err)
	}

	return newReviewer
}

func IsRouletteRequest(payload *MRRoulettePayload) bool {
	return payload.HookType == MRHook && (payload.GitlabAction == "approved" || payload.GitlabAction == "unapproved") ||
		payload.HookType == CommentHook && strings.HasPrefix(payload.Description, ReviewerChangeRequest)
}

func GetLastMRRouletteCommand(requests []gitlabInfo.Note) gitlabInfo.Note {
	if len(requests) == 0 {
		log.Panic("Requests cannot be empty! Check its len.")
	}

	lastRequest := requests[0]
	for i := 0; i < len(requests)-1; i++ {
		if requests[i+1].ID > requests[i].ID {
			lastRequest = requests[i+1]
		}
	}

	return lastRequest
}

func IdentifyRequestType(note string) (RouletteRequestType, error) {
	if strings.HasPrefix(note, RequestWithRandOrder) {
		return RouletteWithRandomOrder, nil
	}

	if strings.HasPrefix(note, RequestWithStrictOrder) {
		return RouletteWithStrictOrder, nil
	}

	if strings.HasPrefix(note, ReviewerChangeRequest) {
		return ReviewerChangeCommand, nil
	}

	return 0, fmt.Errorf("Unknown roulette request type! Note body: %s", note)
}

func GetUsersInRequest(note string, requestType RouletteRequestType) ([]string, error) {
	var m *regexp.Regexp
	switch requestType {
	case RouletteWithRandomOrder:
		m = regexp.MustCompile(RequestWithRandOrder)

	case RouletteWithStrictOrder:
		m = regexp.MustCompile(RequestWithStrictOrder)

	case ReviewerChangeCommand:
		m = regexp.MustCompile(ReviewerChangeRequest)
	}

	withoutCommand := m.ReplaceAllString(note, "")
	cuttedLeftSide := strings.Split(withoutCommand, ">")[0]
	var usersString []string
	if requestType != ReviewerChangeCommand {
		usersString = strings.Split(cuttedLeftSide, ", @")
	} else {
		usersString = append(usersString, cuttedLeftSide)
	}

	if len(usersString) == 0 {
		return nil, fmt.Errorf("Smth wrong in note body: %s", note)
	}

	return usersString, nil
}

func GetNewReviewer(users []string, requestType RouletteRequestType, approvedUsers map[string]gitlabInfo.GitlabUser, projectUsers []gitlabInfo.ProjectUser) (*gitlabInfo.ProjectUser, error) {
	var newReviewerIndex int
	switch requestType {
	case RouletteWithRandomOrder:
		var notApprovedUsers []string
		for i := range users {
			if _, contains := approvedUsers[users[i]]; !contains {
				notApprovedUsers = append(notApprovedUsers, users[i])
			}
		}

		if len(notApprovedUsers) != 0 {
			rand.Seed(time.Now().UnixNano())
			maxIndex := len(notApprovedUsers) - 1
			randomIndex := rand.Intn(maxIndex + 1)
			newReviewerIndex = randomIndex
		}

	case RouletteWithStrictOrder:
		for i := range users {
			if _, contains := approvedUsers[users[i]]; !contains {
				newReviewerIndex = i
				break
			}
		}

	case ReviewerChangeCommand:
		newReviewerIndex = 0
	}

	for i := range projectUsers {
		if projectUsers[i].Username == users[newReviewerIndex] {
			return &projectUsers[i], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("MrRoulette: cant' find required user %s", users[newReviewerIndex]))
}
