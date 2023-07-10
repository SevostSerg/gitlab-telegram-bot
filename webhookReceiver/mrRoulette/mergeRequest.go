package mrRoulette

type MRRoulettePayload struct {
	ProjectID      int64
	MergeRequestID int64
	GitlabAction   string
	Description    string
	HookType       HookType
}

type HookType int

const (
	CommentHook HookType = iota
	MRHook
)
