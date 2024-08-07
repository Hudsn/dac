package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hudsn/dac"
	"github.com/hudsn/dac/examples/foozle"
	"github.com/hudsn/dac/examples/foozle/foozle_rules"
)

func main() {

	engine := dac.NewEngine()

	foozleEval, err := foozle.NewEvaluatorService(foozle_rules.Rules)
	if err != nil {
		log.Fatal(err)
	}

	engine.AddEvaluator("foozle", foozleEval)

	for {
		telem := makeRandomTelem()
		fooz := telem.Data.(foozle.Foozle)

		results, _ := engine.Detect(telem)

		data := struct {
			Foozle  foozle.Foozle         `json:"foozle"`
			Results []dac.DetectionResult `json:"results"`
		}{
			Foozle:  fooz,
			Results: results,
		}
		b, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n\n", string(b))
		time.Sleep(500 * time.Millisecond)
	}
}

func randomBool() bool {
	num := rand.Int() % 2
	return num == 1
}

func makeRandomTelem() dac.Telemetry {
	tags := []string{"evil", "good"}
	idx := rand.Int() % len(tags)
	tags = append(tags[:idx], tags[idx+1:]...)

	isEvil := randomBool()

	f := foozle.Foozle{
		Name:   "test foozle",
		IsEvil: isEvil,
		Tags:   tags,
	}

	return dac.Telemetry{
		TelemetryType: "foozle",
		Data:          f,
	}

}
