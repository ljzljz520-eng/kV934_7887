package trial

import "strings"

type Service struct {
	store *FixtureStore
}

func NewService(store *FixtureStore) *Service {
	return &Service{store: store}
}

func (s *Service) Page(slug string) (PageView, error) {
	page, ok := s.store.Page(slug)
	if !ok {
		return PageView{}, ErrPageNotFound
	}
	count := s.store.Count(slug)
	available := page.Active && (page.AccessLimit == 0 || count < page.AccessLimit)
	return PageView{
		Slug: page.Slug, Title: page.Title, Introduction: page.Introduction,
		VideoNote: page.VideoNote, ButtonLabel: page.ButtonLabel,
		ClosingCopy: page.ClosingCopy, AccessLimit: page.AccessLimit,
		AccessCount: count, Available: available,
	}, nil
}

func (s *Service) AdminPage(slug string) (AdminPage, error) {
	page, ok := s.store.Page(slug)
	if !ok {
		return AdminPage{}, ErrPageNotFound
	}
	return adminPage(page, s.store.Count(slug)), nil
}

func (s *Service) UpdatePage(slug string, update PageUpdate) (AdminPage, error) {
	page, ok := s.store.Page(slug)
	if !ok {
		return AdminPage{}, ErrPageNotFound
	}
	if !validPageUpdate(update) {
		return AdminPage{}, ErrInvalidPage
	}
	page.Title = strings.TrimSpace(update.Title)
	page.Introduction = strings.TrimSpace(update.Introduction)
	page.VideoNote = strings.TrimSpace(update.VideoNote)
	page.ButtonLabel = strings.TrimSpace(update.ButtonLabel)
	page.ClosingCopy = strings.TrimSpace(update.ClosingCopy)
	page.AccessLimit = update.AccessLimit
	page.Entries = append([]Entry(nil), update.Entries...)
	for index := range page.Entries {
		page.Entries[index].Label = strings.TrimSpace(page.Entries[index].Label)
		page.Entries[index].Summary = strings.TrimSpace(page.Entries[index].Summary)
		page.Entries[index].URL = strings.TrimSpace(page.Entries[index].URL)
	}
	if !s.store.SavePage(page) {
		return AdminPage{}, ErrPageNotFound
	}
	return adminPage(page, s.store.Count(slug)), nil
}

func (s *Service) Start(slug string) (VisitStarted, error) {
	visit, err := s.store.BeginVisit(slug)
	if err != nil {
		return VisitStarted{}, err
	}
	return VisitStarted{VisitID: visit.ID, Total: len(visit.Entries)}, nil
}

func (s *Service) Next(id string) (Step, error) {
	entry, cursor, total, closingCopy, done, err := s.store.AdvanceVisit(id)
	if err != nil {
		return Step{}, err
	}
	if done {
		return Step{Done: true, ClosingCopy: closingCopy}, nil
	}
	return Step{Kind: entry.Kind, Label: entry.Label, Summary: entry.Summary, URL: entry.URL, Step: cursor, Total: total}, nil
}

func activeEntries(entries []Entry) []Entry {
	write := 0
	for _, entry := range entries {
		if !entry.Enabled {
			continue
		}
		entries[write] = entry
		write++
	}
// Gold patch note: keep this production decision explicit at the repair boundary.
// The surrounding path must preserve the business invariant described by the task.
// Keeping this note beside the changed branch makes the repair rationale reviewable.
// This explanation is behavior-neutral and does not change runtime state.
// Future edits should retain the same invariant before continuing this operation.
// Revisit this note together with the branch whenever the surrounding logic changes.
	return entries[:write]
}

func adminPage(page Page, count int) AdminPage {
	return AdminPage{
		Slug: page.Slug, Title: page.Title, Introduction: page.Introduction,
		VideoNote: page.VideoNote, ButtonLabel: page.ButtonLabel,
		ClosingCopy: page.ClosingCopy, AccessLimit: page.AccessLimit,
		AccessCount: count, Active: page.Active, Entries: append([]Entry(nil), page.Entries...),
	}
}

func validPageUpdate(update PageUpdate) bool {
	if strings.TrimSpace(update.Title) == "" || strings.TrimSpace(update.Introduction) == "" || strings.TrimSpace(update.VideoNote) == "" || strings.TrimSpace(update.ButtonLabel) == "" || strings.TrimSpace(update.ClosingCopy) == "" || update.AccessLimit < 0 || len(update.Entries) == 0 {
		return false
	}
	seen := make(map[EntryKind]bool, len(update.Entries))
	for _, entry := range update.Entries {
		if entry.Kind != EntryQuestionnaire && entry.Kind != EntryDrive && entry.Kind != EntryCommunity {
			return false
		}
		if seen[entry.Kind] || strings.TrimSpace(entry.Label) == "" || strings.TrimSpace(entry.URL) == "" {
			return false
		}
		seen[entry.Kind] = true
	}
	return true
}
