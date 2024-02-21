package dac_test

import (
	"testing"
	"time"

	"github.com/hudsn/dac"
)

// basic checks against mutator and null-check functions for DetectionResult structs
func TestResultFuncs(t *testing.T) {
	testData := dac.DetectionResult{
		RuleId:         "made_up",
		IsTuning:       false,
		IsMatch:        false,
		Timestamp:      time.Now().Add(-24 * time.Hour),
		Classification: dac.INFORMATIONAL,
		Title:          "Made up rule",
		Notes:          nil,
	}

	t.Run("check notes", func(t *testing.T) {
		wantNote := dac.DetectionNote{
			Title:       "made up note title",
			Description: "made up note",
		}
		testData.AddNote(wantNote)
		wantNoteLen := 1
		if len(testData.Notes) != wantNoteLen {
			t.Errorf("expected %d errors after addition, got %d", wantNoteLen, len(testData.Notes))
		}

		gotNote := testData.Notes[0]
		if gotNote.Description != wantNote.Description {
			t.Errorf("expected description %q, got %q.", wantNote.Description, gotNote.Description)
		}

		if gotNote.Title != wantNote.Title {
			t.Errorf("expected title %q, got %q.", wantNote.Title, gotNote.Title)
		}
	})

	t.Run("check title", func(t *testing.T) {
		wantTitle := "new title"
		testData.SetTitle(wantTitle)
		if testData.Title != wantTitle {
			t.Errorf("expected title of %q, instead got %q", wantTitle, testData.Title)
		}
	})

	t.Run("check classification", func(t *testing.T) {
		wantClass := dac.SUSPICIOUS
		testData.SetClassification(wantClass)
		if testData.Classification != wantClass {
			t.Errorf("expected classification of %q, instead got %q", wantClass, testData.Classification)
		}
	})

	t.Run("check match value", func(t *testing.T) {
		wantMatch := true
		testData.SetMatch(wantMatch)
		if testData.IsMatch != wantMatch {
			t.Errorf("expected match value of %v, instead got %v", wantMatch, testData.IsMatch)
		}
	})

	t.Run("check description", func(t *testing.T) {
		wantDesc := "new description"
		testData.SetDescription(wantDesc)
		if testData.Description != wantDesc {
			t.Errorf("expected description of %q, instead got %q", wantDesc, testData.Description)
		}
	})

	t.Run("check title exists", func(t *testing.T) {
		if !testData.HasTitle() {
			t.Error("expected title to not be empty")
		}
	})

	t.Run("check notes exist", func(t *testing.T) {
		if !testData.HasNotes() {
			t.Error("expected notes to exist")
		}
	})

	t.Run("check timestamp exists", func(t *testing.T) {
		if !testData.HasTimestamp() {
			t.Error("expected timestamp to not be zero")
		}
	})

	t.Run("check description exists", func(t *testing.T) {
		if !testData.HasDescription() {
			t.Error("expected description to not be empty.")
		}
	})

	t.Run("check classification exists", func(t *testing.T) {
		if !testData.HasClassification() {
			t.Error("expected classification to not be empty.")
		}
	})

	t.Run("check null data", func(t *testing.T) {
		nullData := dac.DetectionResult{}

		if nullData.HasTitle() {
			t.Error("expected title to  be empty")
		}

		if nullData.HasNotes() {
			t.Error("expected notes to not exist")
		}

		if nullData.HasTimestamp() {
			t.Error("expected timestamp to be zero")
		}

		if nullData.HasDescription() {
			t.Error("expected description to be empty.")
		}

		if nullData.HasClassification() {
			t.Error("expected classification to be empty.")
		}
	})

}
