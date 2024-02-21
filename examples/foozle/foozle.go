package foozle

import (
	"embed"
	"log"

	"github.com/hudsn/dac"
)

//go:embed foozle_rules/*.yaml
var FoozleRules embed.FS

//go:embed example.json
var ExampleFoozle []byte

var ExampleDataMap map[string][]byte = map[string][]byte{
	"example_foozle.json": ExampleFoozle,
}

type NullableFoozle struct {
	Name   *string  `json:"name,omitempty" yaml:"name,omitempty"`
	IsEvil *bool    `json:"is_evil,omitempty" yaml:"is_evil,omitempty"`
	Tags   []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

func (nf *NullableFoozle) ToFoozle() Foozle {
	retVal := &Foozle{}
	fieldMetas := &FoozleMeta{}
	if nf.Name != nil {
		retVal.Name = *nf.Name
		fieldMetas.HasName = true
	}
	if nf.IsEvil != nil {
		retVal.IsEvil = *nf.IsEvil
		fieldMetas.HasIsEvil = true
	}
	if nf.Tags != nil {
		retVal.Tags = nf.Tags
		fieldMetas.HasTags = true
	}

	retVal.FoozleMeta = *fieldMetas

	return *retVal
}

type Foozle struct {
	Name       string   `json:"name" yaml:"name"`
	IsEvil     bool     `json:"is_evil" yaml:"is_evil"`
	Tags       []string `json:"tags" yaml:"tags"`
	FoozleMeta `json:"-"`
}

type FoozleMeta struct {
	HasIsEvil bool
	HasName   bool
	HasTags   bool
}

type FoozleRule interface {
	Evaluate(Foozle) (dac.DetectionResult, error)
}

type RuleInitFunc func(dac.RuleMeta[Foozle]) FoozleRule
type FoozleInitMap map[string]RuleInitFunc

type FoozleEvaluator struct {
	Rules []FoozleRule
}

func (fe *FoozleEvaluator) GetMatches(data any) []dac.DetectionResult {
	switch v := data.(type) {
	case Foozle:
		return fe.getMatches(v)
	default:
		return nil
	}
}

func (fe *FoozleEvaluator) getMatches(data Foozle) []dac.DetectionResult {
	results := []dac.DetectionResult{}
	for _, entry := range fe.Rules {
		detectionResult, err := entry.Evaluate(data)
		if err != nil {
			log.Println(err)
			continue
		}
		if detectionResult.IsMatch {
			results = append(results, detectionResult)
		}
	}

	return results
}

func NewEvaluatorService(ruleMap FoozleInitMap) (dac.EvaluatorService, error) {

	ruleList := []FoozleRule{}

	ruleMetas, err := dac.ReadYML[Foozle](FoozleRules)
	if err != nil {
		return nil, err
	}

	for _, ruleMeta := range ruleMetas {
		// don't initialize any inactive rules
		if !ruleMeta.IsActive {
			continue
		}
		if initFuncEntry, ok := ruleMap[ruleMeta.InternalId]; ok {
			rule := initFuncEntry(ruleMeta)
			ruleList = append(ruleList, rule)
		}
	}

	return &FoozleEvaluator{
		Rules: ruleList,
	}, nil
}
