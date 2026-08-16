package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"course-trial/internal/trial"
)

func TestCoursePageEndpointReturnsVisitorContent(t *testing.T) {
	server := httptest.NewServer(New(trial.NewService(trial.NewFixtureStore(trial.FixturePages()))))
	defer server.Close()
	response, err := http.Get(server.URL + "/api/pages/go-from-zero")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected page response, got %d", response.StatusCode)
	}
	var page trial.PageView
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Title == "" || page.VideoNote == "" || page.ButtonLabel == "" {
		t.Fatalf("missing visitor content: %+v", page)
	}
}

func TestCourseVisitEndpointMovesThroughDestinations(t *testing.T) {
	server := httptest.NewServer(New(trial.NewService(trial.NewFixtureStore(trial.FixturePages()))))
	defer server.Close()
	response, err := http.Post(server.URL+"/api/pages/go-from-zero/visits", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var started trial.VisitStarted
	if response.StatusCode != http.StatusCreated || json.NewDecoder(response.Body).Decode(&started) != nil {
		t.Fatalf("unexpected visit response: %d", response.StatusCode)
	}
	next, err := http.Post(server.URL+"/api/visits/"+started.VisitID+"/next", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer next.Body.Close()
	var step trial.Step
	if next.StatusCode != http.StatusOK || json.NewDecoder(next.Body).Decode(&step) != nil || step.Kind != trial.EntryQuestionnaire {
		t.Fatalf("unexpected first destination: %d %+v", next.StatusCode, step)
	}
}

func TestCourseVisitEndpointReportsClosingCopyWhenLimitIsReached(t *testing.T) {
	server := httptest.NewServer(New(trial.NewService(trial.NewFixtureStore(trial.FixturePages()))))
	defer server.Close()
	first, err := http.Post(server.URL+"/api/pages/go-limited/visits", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	first.Body.Close()
	second, err := http.Post(server.URL+"/api/pages/go-limited/visits", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Body.Close()
	var result map[string]string
	if second.StatusCode != http.StatusConflict || json.NewDecoder(second.Body).Decode(&result) != nil || result["closing_copy"] != "本期试听名额已满，感谢你的关注。" {
		t.Fatalf("unexpected closed visit response: %d %+v", second.StatusCode, result)
	}
}

func TestAdminEndpointUpdatesCourseContent(t *testing.T) {
	server := httptest.NewServer(New(trial.NewService(trial.NewFixtureStore(trial.FixturePages()))))
	defer server.Close()
	pageResponse, err := http.Get(server.URL + "/api/admin/pages/go-from-zero")
	if err != nil {
		t.Fatal(err)
	}
	var page trial.AdminPage
	if pageResponse.StatusCode != http.StatusOK || json.NewDecoder(pageResponse.Body).Decode(&page) != nil {
		pageResponse.Body.Close()
		t.Fatalf("unexpected admin page response: %d", pageResponse.StatusCode)
	}
	pageResponse.Body.Close()
	update := trial.PageUpdate{Title: page.Title, Introduction: "老师保存的新介绍", VideoNote: page.VideoNote, ButtonLabel: "打开新资料", ClosingCopy: page.ClosingCopy, AccessLimit: page.AccessLimit, Entries: page.Entries}
	body, err := json.Marshal(update)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodPut, server.URL+"/api/admin/pages/go-from-zero", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var saved trial.AdminPage
	if response.StatusCode != http.StatusOK || json.NewDecoder(response.Body).Decode(&saved) != nil || saved.Introduction != "老师保存的新介绍" || saved.ButtonLabel != "打开新资料" {
		t.Fatalf("unexpected saved admin page: %d %+v", response.StatusCode, saved)
	}
}
