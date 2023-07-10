package botTypes

type WebhookInformation struct {
	Type      WebhookType
	Message   string
	ProjectID int
	Author    string
}

type WebhookType int

const (
	Push WebhookType = iota
	MR
	Comment
	Build
	Issue
	Pipeline
	Tag
	UnknownHook
)
