package trial

import (
	"errors"
	"testing"
)

func TestVisitorReceivesCourseInformation(t *testing.T) {
	service := NewService(NewFixtureStore(FixturePages()))
	page, err := service.Page("go-from-zero")
	if err != nil {
		t.Fatal(err)
	}
	if page.Title != "Go 从零到上线" || page.VideoNote == "" || page.ButtonLabel != "开始试听" || !page.Available {
		t.Fatalf("unexpected course page: %+v", page)
	}
}

func TestVisitorFollowsQuestionnaireDriveAndCommunityInOrder(t *testing.T) {
	service := NewService(NewFixtureStore(FixturePages()))
	started, err := service.Start("go-from-zero")
	if err != nil {
		t.Fatal(err)
	}
	if started.Total != 3 {
		t.Fatalf("expected three destinations, got %d", started.Total)
	}
	expected := []EntryKind{EntryQuestionnaire, EntryDrive, EntryCommunity}
	for index, kind := range expected {
		step, nextErr := service.Next(started.VisitID)
		if nextErr != nil || step.Done || step.Kind != kind || step.Step != index+1 {
			t.Fatalf("unexpected destination at %d: %+v, %v", index, step, nextErr)
		}
	}
	finished, err := service.Next(started.VisitID)
	if err != nil || !finished.Done || finished.ClosingCopy == "" {
		t.Fatalf("expected completion message, got %+v, %v", finished, err)
	}
}

func TestVisitorSequenceExcludesDisabledCourseEntry(t *testing.T) {
	service := NewService(NewFixtureStore(FixturePages()))
	started, err := service.Start("go-disabled-drive")
	if err != nil {
		t.Fatal(err)
	}
	if started.Total != 2 {
		t.Fatalf("expected two available destinations, got %d", started.Total)
	}
	first, err := service.Next(started.VisitID)
	if err != nil || first.Kind != EntryQuestionnaire {
		t.Fatalf("expected questionnaire first, got %+v, %v", first, err)
	}
	second, err := service.Next(started.VisitID)
	if err != nil || second.Kind != EntryCommunity {
		t.Fatalf("expected community second, got %+v, %v", second, err)
	}
	finished, err := service.Next(started.VisitID)
	if err != nil || !finished.Done {
		t.Fatalf("expected sequence to finish, got %+v, %v", finished, err)
	}
}

func TestVisitorSeesClosingCopyAfterAccessLimit(t *testing.T) {
	service := NewService(NewFixtureStore(FixturePages()))
	if _, err := service.Start("go-limited"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Start("go-limited"); !errors.Is(err, ErrAccessLimit) {
		t.Fatalf("expected access limit, got %v", err)
	}
	page, err := service.Page("go-limited")
	if err != nil || page.Available || page.ClosingCopy != "本期试听名额已满，感谢你的关注。" {
		t.Fatalf("unexpected closed page: %+v, %v", page, err)
	}
}

func TestVisitorCannotOpenInactiveCourse(t *testing.T) {
	service := NewService(NewFixtureStore(FixturePages()))
	if _, err := service.Start("go-archived"); !errors.Is(err, ErrPageInactive) {
		t.Fatalf("expected inactive status, got %v", err)
	}
}

func TestTeacherUpdatesCourseContentForVisitor(t *testing.T) {
	service := NewService(NewFixtureStore(FixturePages()))
	page, err := service.AdminPage("go-from-zero")
	if err != nil {
		t.Fatal(err)
	}
	updated, err := service.UpdatePage(page.Slug, PageUpdate{
		Title: page.Title, Introduction: "新的课程介绍", VideoNote: "新的试听视频说明",
		ButtonLabel: "领取新版资料", ClosingCopy: "本次试听已结束。", AccessLimit: 8, Entries: page.Entries,
	})
	if err != nil || updated.Introduction != "新的课程介绍" || updated.AccessLimit != 8 {
		t.Fatalf("unexpected saved course: %+v, %v", updated, err)
	}
	visitor, err := service.Page(page.Slug)
	if err != nil || visitor.VideoNote != "新的试听视频说明" || visitor.ButtonLabel != "领取新版资料" || visitor.ClosingCopy != "本次试听已结束。" {
		t.Fatalf("unexpected visitor page after save: %+v, %v", visitor, err)
	}
}

func TestTeacherCannotSaveIncompleteCourseContent(t *testing.T) {
	service := NewService(NewFixtureStore(FixturePages()))
	page, err := service.AdminPage("go-from-zero")
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.UpdatePage(page.Slug, PageUpdate{Title: page.Title, Entries: page.Entries})
	if !errors.Is(err, ErrInvalidPage) {
		t.Fatalf("expected invalid course status, got %v", err)
	}
}
