package foozle_rules

import (
	"github.com/hudsn/dac"
	"github.com/hudsn/dac/examples/foozle"
)

type evilFoozle struct {
	defaultData
}

func (r evilFoozle) Evaluate(f foozle.Foozle) (dac.DetectionResult, error) {
	result := r.defaultResult()
	if f.IsEvil {
		note := dac.DetectionNote{
			Title:       "very evil foozle",
			Description: "this is indicative of an exceptionally evil foozle",
		}
		result.AddNote(note)
		result.SetMatch(true)
	}

	return result, nil
}

func EvilFoozleInit(meta dac.RuleMeta[foozle.Foozle]) foozle.FoozleRule {
	defaults := defaultData{
		metadata: meta,
	}

	return evilFoozle{
		defaultData: defaults,
	}
}
