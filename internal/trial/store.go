package trial

import (
	"fmt"
	"sync"
)

type FixtureStore struct {
	mu     sync.Mutex
	pages  map[string]Page
	visits map[string]Visit
	counts map[string]int
	nextID int
}

func NewFixtureStore(pages []Page) *FixtureStore {
	pageMap := make(map[string]Page, len(pages))
	counts := make(map[string]int, len(pages))
	for _, page := range pages {
		page.Entries = append([]Entry(nil), page.Entries...)
		pageMap[page.Slug] = page
	}
	return &FixtureStore{pages: pageMap, visits: make(map[string]Visit), counts: counts}
}

func (s *FixtureStore) Page(slug string) (Page, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	page, ok := s.pages[slug]
	if !ok {
		return Page{}, false
	}
	page.Entries = append([]Entry(nil), page.Entries...)
	return page, true
}

func (s *FixtureStore) Count(slug string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counts[slug]
}

func (s *FixtureStore) SavePage(page Page) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pages[page.Slug]; !ok {
		return false
	}
	page.Entries = append([]Entry(nil), page.Entries...)
	s.pages[page.Slug] = page
	return true
}

func (s *FixtureStore) BeginVisit(slug string) (Visit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	page, ok := s.pages[slug]
	if !ok {
		return Visit{}, ErrPageNotFound
	}
	if !page.Active {
		return Visit{}, ErrPageInactive
	}
	if page.AccessLimit > 0 && s.counts[page.Slug] >= page.AccessLimit {
		return Visit{}, ErrAccessLimit
	}
	s.counts[page.Slug]++
	s.nextID++
	entries := activeEntries(append([]Entry(nil), page.Entries...))
	visit := Visit{ID: "visit-" + formatID(s.nextID), PageSlug: page.Slug, Entries: append([]Entry(nil), entries...)}
	s.visits[visit.ID] = visit
	return visit, nil
}

func (s *FixtureStore) AdvanceVisit(id string) (Entry, int, int, string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	visit, ok := s.visits[id]
	if !ok {
		return Entry{}, 0, 0, "", false, ErrVisitNotFound
	}
	page, ok := s.pages[visit.PageSlug]
	if !ok {
		return Entry{}, 0, 0, "", false, ErrPageNotFound
	}
	if visit.Cursor >= len(visit.Entries) {
		return Entry{}, visit.Cursor, len(visit.Entries), page.ClosingCopy, true, nil
	}
	entry := visit.Entries[visit.Cursor]
	visit.Cursor++
	s.visits[visit.ID] = visit
	return entry, visit.Cursor, len(visit.Entries), "", false, nil
}

func formatID(value int) string {
	return fmt.Sprintf("%03d", value)
}
