package dac_test

import (
	"testing"

	"github.com/hudsn/dac"
	"github.com/hudsn/dac/examples/foozle"
	"github.com/hudsn/dac/examples/foozle/foozle_rules"
)

// ensuring that engine creation and processing works at a basic level
func TestDetection(t *testing.T) {
	testEngine := dac.NewEngine()

	testTelemetry := dac.Telemetry{
		TelemetryType: "foozle",
		Data: foozle.Foozle{
			Name:   "evil foozle",
			IsEvil: true,
			Tags:   []string{"evil"},
		},
	}

	testErrTelemetry := dac.Telemetry{
		TelemetryType: "fake",
		Data: struct {
			Who   string
			Cares string
		}{
			Who:   "asdf",
			Cares: "fdsa",
		},
	}

	foozleEvaluator, err := foozle.NewEvaluatorService(foozle_rules.Rules)
	if err != nil {
		t.Fatal(err)
	}

	testEngine.AddEvaluator("foozle", foozleEvaluator)

	t.Run("check results count", func(t *testing.T) {
		results, _ := testEngine.Detect(testTelemetry)

		wantCount := 2
		if len(results) != wantCount {
			t.Errorf("expected %d detection results, instead got %d", wantCount, len(results))
		}
	})

	t.Run("check index not found", func(t *testing.T) {
		_, ok := testEngine.Detect(testErrTelemetry)
		if ok {
			t.Error("expected second return value of 'Detect()' to be false. got true")
		}
	})
}
