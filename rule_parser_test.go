package dac_test

import (
	"embed"
	"encoding/json"
	"testing"

	"github.com/hudsn/dac"
)

//go:embed rule_parser_test.yml
var testYML embed.FS

func TestReadYML(t *testing.T) {
	type fakeData struct {
		MadeUp   string `yaml:"made_up_field"`
		NumThing int    `yaml:"another_real_field_but_an_int_this_time"`
	}

	ruleMetas, err := dac.ReadYML[fakeData](testYML)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("number of rules", func(t *testing.T) {
		if len(ruleMetas) != 1 {
			t.Errorf("expected one extracted rule, instead got %d", len(ruleMetas))
		}
	})

	meta := ruleMetas[0]

	t.Run("check analysis steps", func(t *testing.T) {
		analysisStepsWant := "IDK just like figure it out silly analyst lol"
		if meta.AnalysisSteps != analysisStepsWant {
			t.Errorf("expected analysis steps to be %q, got %q", analysisStepsWant, meta.AnalysisSteps)
		}
	})

	t.Run("check classification", func(t *testing.T) {
		classificationWant := dac.INFORMATIONAL
		if meta.DefaultClassification != classificationWant {
			t.Errorf("expected classification to be %q, got %q", classificationWant, meta.DefaultClassification)
		}
	})

	t.Run("check active", func(t *testing.T) {
		activeWant := true
		if meta.IsActive != activeWant {
			t.Errorf("expected active to be %v, got %v", activeWant, meta.IsActive)
		}
	})

	t.Run("check tuning", func(t *testing.T) {
		tuningWant := false
		if meta.IsTuning != tuningWant {
			t.Errorf("expected tuning to be %v, got %v", tuningWant, meta.IsTuning)
		}
	})

	t.Run("check false positives", func(t *testing.T) {
		falsePositivesWant := "None this is the best rule ever and will never ever return any false positives lmao"
		if meta.FalsePositives != falsePositivesWant {
			t.Errorf("expected false positive field to be %q, instead got %q", falsePositivesWant, meta.FalsePositives)
		}
	})

	t.Run("check references", func(t *testing.T) {
		referencesLenWant := 1
		if len(meta.References) != referencesLenWant {
			t.Errorf("expected references to be of len 1, instead got %d", len(meta.References))
		}
		references0Want := "example.com"
		if meta.References[0] != references0Want {
			t.Errorf("expected references value to be %q, instead got %q", references0Want, meta.References[0])
		}
	})

	t.Run("check test reference file", func(t *testing.T) {
		wantTestBase := "non-existent-file.json"
		if meta.TestBaseFile != wantTestBase {
			t.Errorf("expected test base file name to be %q. instead got %q", wantTestBase, meta.TestBaseFile)
		}
	})

	for _, test := range meta.Tests {

		wantBytes, err := json.Marshal(test.Want)
		if err != nil {
			t.Fatal(err)
		}

		wantResult := dac.DetectionResult{}
		err = json.Unmarshal(wantBytes, &wantResult)
		if err != nil {
			t.Fatal(err)
		}

		t.Run("check detection notes", func(t *testing.T) {
			if !wantResult.HasNotes() {
				t.Error("expected to parse out a single desired note from the testing section. instead got none")
			}
			if len(wantResult.Notes) != 1 {
				t.Errorf("expected to parse out a single desired note from the testing section. instead got %d", len(wantResult.Notes))
			}

			wantNotes := dac.DetectionNote{
				Title:       "some note I guess",
				Description: "this is super evil you should be scared",
			}

			if wantResult.Notes[0].Title != wantNotes.Title {
				t.Errorf("expected to parse out a note with title %q. instead got %q", wantNotes.Title, wantResult.Notes[0].Title)

			}

			if wantResult.Notes[0].Description != wantNotes.Description {
				t.Errorf("expected to parse out a note with description %q. instead got %q", wantNotes.Description, wantResult.Notes[0].Description)
			}
		})

		t.Run("check is_match", func(t *testing.T) {
			if wantResult.IsMatch != true {
				t.Error("expected 'is_match' requirement to be parsed as true. instead got false.")
			}
		})

		t.Run("check test override values", func(t *testing.T) {
			override := test.Override

			if override.MadeUp != "bananarama" {
				t.Errorf("expected override field to parse out 'bananarama'. instead got %v", override.MadeUp)
			}
			if override.NumThing != 999 {
				t.Errorf("expected override field to parse out 999. instead got %v", override.NumThing)
			}
		})

	}

}
