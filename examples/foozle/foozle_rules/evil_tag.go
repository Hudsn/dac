package foozle_rules

import (
	"github.com/hudsn/dac"
	"github.com/hudsn/dac/examples/foozle"
)

type evilTag struct {
	defaultData
}

func (r evilTag) Evaluate(f foozle.Foozle) (dac.DetectionResult, error) {
	result := r.defaultResult()

	for _, tag := range f.Tags {
		if tag == "evil" {
			result.IsMatch = true
		}
	}

	return result, nil
}

func EvilTagInit(meta dac.RuleMeta[foozle.Foozle]) foozle.FoozleRule {
	defaults := defaultData{
		metadata: meta,
	}

	return evilTag{
		defaultData: defaults,
	}
}
