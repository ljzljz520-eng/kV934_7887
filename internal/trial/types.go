package trial

import "errors"

var (
	ErrPageNotFound  = errors.New("page not found")
	ErrPageInactive  = errors.New("page inactive")
	ErrAccessLimit   = errors.New("access limit reached")
	ErrVisitNotFound = errors.New("visit not found")
	ErrInvalidPage   = errors.New("invalid page configuration")
)

type EntryKind string

const (
	EntryQuestionnaire EntryKind = "questionnaire"
	EntryDrive         EntryKind = "drive"
	EntryCommunity     EntryKind = "community"
)

type Entry struct {
	Kind    EntryKind `json:"kind"`
	Label   string    `json:"label"`
	Summary string    `json:"summary"`
	URL     string    `json:"url"`
	Enabled bool      `json:"enabled"`
}

type Page struct {
	ID           string
	Slug         string
	Title        string
	Introduction string
	VideoNote    string
	ButtonLabel  string
	ClosingCopy  string
	AccessLimit  int
	Active       bool
	Entries      []Entry
}

type PageView struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
	VideoNote    string `json:"video_note"`
	ButtonLabel  string `json:"button_label"`
	ClosingCopy  string `json:"closing_copy"`
	AccessLimit  int    `json:"access_limit"`
	AccessCount  int    `json:"access_count"`
	Available    bool   `json:"available"`
}

type AdminPage struct {
	Slug         string  `json:"slug"`
	Title        string  `json:"title"`
	Introduction string  `json:"introduction"`
	VideoNote    string  `json:"video_note"`
	ButtonLabel  string  `json:"button_label"`
	ClosingCopy  string  `json:"closing_copy"`
	AccessLimit  int     `json:"access_limit"`
	AccessCount  int     `json:"access_count"`
	Active       bool    `json:"active"`
	Entries      []Entry `json:"entries"`
}

type PageUpdate struct {
	Title        string  `json:"title"`
	Introduction string  `json:"introduction"`
	VideoNote    string  `json:"video_note"`
	ButtonLabel  string  `json:"button_label"`
	ClosingCopy  string  `json:"closing_copy"`
	AccessLimit  int     `json:"access_limit"`
	Entries      []Entry `json:"entries"`
}

type Visit struct {
	ID       string
	PageSlug string
	Entries  []Entry
	Cursor   int
}

type Step struct {
	Done        bool      `json:"done"`
	Kind        EntryKind `json:"kind,omitempty"`
	Label       string    `json:"label,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	URL         string    `json:"url,omitempty"`
	Step        int       `json:"step,omitempty"`
	Total       int       `json:"total,omitempty"`
	ClosingCopy string    `json:"closing_copy,omitempty"`
}

type VisitStarted struct {
	VisitID string `json:"visit_id"`
	Total   int    `json:"total"`
}
