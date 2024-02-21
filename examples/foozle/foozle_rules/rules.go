package foozle_rules

import (
	"time"

	"github.com/hudsn/dac"
	"github.com/hudsn/dac/examples/foozle"
)

var Rules map[string]foozle.RuleInitFunc = map[string]foozle.RuleInitFunc{
	"evil_foozle": EvilFoozleInit,
	"evil_tag":    EvilTagInit,
}

type defaultData struct {
	metadata dac.RuleMeta[foozle.Foozle]
}

func (dd *defaultData) defaultResult() dac.DetectionResult {

	retVal := dac.DetectionResult{
		RuleId:         dd.metadata.InternalId,
		IsTuning:       dd.metadata.IsTuning,
		Timestamp:      time.Now().UTC(),
		Title:          dd.metadata.DisplayName,
		Description:    dd.metadata.DefaultDescription,
		Classification: dd.metadata.DefaultClassification,
		IsMatch:        false,
		Notes:          nil,
	}

	return retVal
}
