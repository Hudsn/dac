package foozle_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hudsn/dac"
	"github.com/hudsn/dac/examples/foozle"
	"github.com/hudsn/dac/examples/foozle/foozle_rules"
)

// run each test case of each yml file

func TestRules(t *testing.T) {
	ruleMetas, err := dac.ReadYML[foozle.Foozle](foozle.FoozleRules)
	if err != nil {
		t.Fatal(err)
	}
	for _, metaEntry := range ruleMetas {
		exampleFoozle := mustGetExampleFoozle(t, metaEntry.TestBaseFile)
		rule := mustGetRule(t, metaEntry)

		for i, test := range metaEntry.Tests {

			t.Run(fmt.Sprintf("check test override. %s yaml test #%d", metaEntry.InternalId, i), func(t *testing.T) {
				overrideBytes, err := json.Marshal(test.Override)
				if err != nil {
					t.Fatal(err)
				}

				// using a weird nullable conversion with pointers so that we can check for absence of types lke bools (ex: is_evil)
				// probably a better way to do this, but this scuffed way seems to work for now.
				tempFoozle := foozle.NullableFoozle{}
				err = json.Unmarshal(overrideBytes, &tempFoozle)
				if err != nil {
					t.Fatal(err)
				}

				newFoozle := tempFoozle.ToFoozle()

				if newFoozle.HasIsEvil {
					exampleFoozle.IsEvil = newFoozle.IsEvil
				}
				if newFoozle.HasName {
					exampleFoozle.Name = newFoozle.Name
				}
				if newFoozle.HasTags {
					exampleFoozle.Tags = newFoozle.Tags
				}
			})

			// BEGIN RESULT CHECKING -- CAN REUSE THIS SEGMENT IN NON-MADE-UP IMPLEMENTATIONS
			t.Run(test.It, func(t *testing.T) {
				gotResult, err := rule.Evaluate(exampleFoozle)
				if err != nil {
					t.Fatal(err)
				}
				wantResult := test.Want
				if gotResult.IsMatch != wantResult.IsMatch {
					t.Errorf("%s : %s -- incorrect match values. expected %v. instead got %v.", metaEntry.DisplayName, test.It, wantResult.IsMatch, gotResult.IsMatch)
				}

				if wantResult.HasTitle() {
					if gotResult.Title != wantResult.Title {
						t.Errorf("%s : %s -- expected detection title %q. instead got %q.", metaEntry.DisplayName, test.It, wantResult.Title, gotResult.Title)
					}
				}
				if wantResult.HasDescription() {
					if gotResult.Description != wantResult.Description {
						t.Errorf("%s : %s -- expected detection description %q. instead got %q.", metaEntry.DisplayName, test.It, wantResult.Description, gotResult.Description)
					}
				}

				if wantResult.Notes == nil && len(gotResult.Notes) > 0 {
					t.Errorf("wanted no notes, got %d", len(gotResult.Notes))
				}
				if wantResult.HasNotes() {
					if len(wantResult.Notes) != len(gotResult.Notes) {
						t.Errorf("%s : %s -- expected notes length of %d, instead got %d.", metaEntry.DisplayName, test.It, len(wantResult.Notes), len(gotResult.Notes))
					}
					for i, note := range wantResult.Notes {
						if gotResult.Notes[i].Title != note.Title ||
							gotResult.Notes[i].Description != note.Description {
							t.Errorf("expected notes to be equal. Instead got: \n%v\nwhen wanted:\n%v\n\n", gotResult.Notes[i], note)
						}
					}
				}

			})

		}
	}
}

func mustGetExampleFoozle(t *testing.T, foozleName string) foozle.Foozle {
	t.Helper()
	exampleData, ok := foozle.ExampleDataMap[foozleName]
	if !ok {
		t.Fatal("could not find example data mapping in foozle package")
	}
	var unwrap struct {
		Data foozle.Foozle `json:"data"`
	}
	err := json.Unmarshal(exampleData, &unwrap)
	if err != nil {
		t.Fatal(err)
	}
	return unwrap.Data
}

func mustGetRule(t *testing.T, ruleMeta dac.RuleMeta[foozle.Foozle]) foozle.FoozleRule {
	t.Helper()
	initFunc, ok := foozle_rules.Rules[ruleMeta.InternalId]
	if !ok {
		t.Fatal("unable to find rule for test")
	}
	rule := initFunc(ruleMeta)
	return rule
}
