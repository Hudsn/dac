# DAC (Detection as Code)
A small detection-as-code implementation written in Go

This package mainly contains some struct and interface definitions that outline the general shape of data that makes up a minimalist detection-as-code pipeline and provides a small but hopefully comprehensive example implementation package to reference (or copy/paste).

## Usage

Note that **this is not a boilerplate or code generator**, and due to the varying nature of different telemetry sources you will have to implement functionality that uses these data types yourself. 

However, this can be used as a starting point "skeleton" for an initial detection as code pipeline. 

See the `examples` folder for a sample implementation for running detections against an imaginary telemetry type called a `foozle`.

You can run the foozle simulator to output an example detection stream to your terminal just to illustrate what it looks like:

```bash
go run examples/cmd/foozle_simulator.go
```
But if you're lazy like me and don't really want to clone the repo or run the command, here's output from our example implementation pipeline to detect evil foozles:

```json
{
	"foozle": {
		"name": "test foozle",
		"is_evil": true,
		"tags": [
			"evil"
		]
	},
	"results": [
		{
			"rule_id": "evil_foozle",
			"is_tuning": false,
			"timestamp": "2024-02-20T15:58:52.66923Z",
			"title": "Evil foozle",
			"description": "This rule is meant to identify evil foozles based on the 'is_evil' field in the foozle data structure.",
			"classification": "malicious",
			"is_match": true,
			"notes": [
				{
					"title": "very evil foozle",
					"description": "this is indicative of an exceptionally evil foozle"
				}
			]
		},
		{
			"rule_id": "evil_tag",
			"is_tuning": false,
			"timestamp": "2024-02-20T15:58:52.66923Z",
			"title": "Evil tag",
			"description": "This rule is meant to identify evil foozles based on the tags containing an 'evil' value in the foozle data structure.",
			"classification": "malicious",
			"is_match": true,
			"notes": null
		}
	]
}
```

---
## For the curious: Detection Entities

***NOTE: YOU DO NOT NEED TO READ ALL THIS EXTRA STUFF TO START USING THIS - IT'S SIMPLY MADE AVAILABLE FOR THOSE WHO ARE CURIOUS AS TO HOW THIS WORKS OR WANT TO REFERENCE IT LATER (ME BECAUSE I'LL COME BACK IN 2 WEEKS AND FORGET EVERYTHING I WROTE)***

The type definitions in the `detection_entities.go` file should serve as a good reference for the "shape" of data in our detection pipeline. There should be comments containing a general description of each data type. However a summary of each significant data type is included here as well...

### Engine

This is top-level struct that we'll be using to evaluate any incoming telemetry source via `Engine.Detect(telemetry)`. This will effectively just run the telemetry data through a series of detection evaluator functions and ones that return `true` will generate detection items that you can then do with as you please. Some ideas might be sending webhooks or emails to generate alerts.

The Engine struct contains one data type `SourceIndex` which holds a map of possible data source index values that point to specific telemetry evaluators. 

For example you might have `event_logs` as a key in `SourceIndex` that points to a value of an [`EvaluatorService`](###EvaluatorService) implementation containing detection rules and helper functions for transforming and evaluating event logs. When telemetry comes in, we can use its [`TelemetryType`](###Telemetry) to route it to the approriate evaluator service.


You'll create the engine with something like:

```go
myEngine := NewEngine()
```
When you've implemented your own service to satisfy the [`EvaluatorService`](###EvaluatorService) interface, you can add it to your engine like so:

```go
// you'll need to implement your own evaluator. See the examples folder for one possible approach.
eventLogEvaluator := myOwnPackage.NewEventLogEvaluator()

//event_logs is the imaginary name of our telemetry source type.
myEngine.AddEvaluator("event_logs", eventLogEvaluator)
```

Finally, when you have your evaluator(s) added to your engine you can pass telemetry to it for evaluation: 

```go
detectionResults := myEngine.Detect(eventLogItem)
```

### Telemetry

This is the source type and actual raw telemetry data to be evaluated by an evaluator service. Note that when you implement the evaluator service you'll probably need to use a type assertion on the incoming telemetry data since we initially declare it as `any` in this struct, and will pass this `Data` field into the evaluator's `GetMatches(any)` function. Reference the `GetMatches` function in `examples/foozle/foozle.go` for an example on how you might do this.

### EvaluatorService

This is responsible for actually processing the `Data` part of a a telemetry struct.

You will need to implement this yourself and satisfy the interface, which is simply creating a type with a `GetMatches()` function. This interface is defined in the `detection_entities.go` file, and an example implementation can be found in `examples/foozle/foozle.go`.

Once you create an evaluator service you can add it to an engine struct using `Engine.AddEvaluator()`.

### DetectionResult 

This struct represents a resulting detection item. 
Any telemetry that evaluates to 

The fields that are contained in this struct should mostly be self explanatory based on name, except for maybe the `Notes` field.  

`Notes` can contain any conditional just-in-time context that should be presented alongside the other detection data. For example, you may have a generic detection for malicious activity, and then in your detection logic you might add a note for a particular subset of that activity that is closely tied to a particular threat actor. Or maybe you simply wish to point analysts consuming this data to a particular field of the telemetry to highlight the "why" of the detection.

### RuleMeta 
This is metadata for any given detection rule. In our example implementation we define these values in the `.yaml` files in the `examples/foozle/foozle_rules` directory. 

A brief summary of each field is included below:

- **Sourcetype**: The name of the index the rule is being evaluated against. This is going to be the same the `SourceIndex` key in our [Engine](###Engine) struct.
- **InternalID**: The internal name of our rule. This can be used as values in other maps containing rule data. For example your InternalID `my_detection_rule` might point to a struct or function (or whatever type you want to implement) that represents the actual detection logic. An example map can be found in `examples/foozle/foozle_rules/rules.go`
- **DisplayName**: The much more readable version of your rule name and should ideally be the name/title that is presented to analysts. Ex: "My Detection Rule"
- **CreatedAt/UpdatedAt**: Time the detection rule was created or updated, respectively.
- **IsActive**: Whether the rule is active or not. Your evaluator service should not process rules that aren't active.
- **IsTuning** : Whether the rule is in tuning or not. Your evaluator service should still generate [DetectionResult](###DetectionResult) items for tuning rules with the `IsTuning` field set to true. The default result example in `examples/foozle/foozle_rules/rules.go` demonstrates this.
- **DefaultDescription**: The description to attach to the `Description` field of `DetectionResult` if it is left unchanged.
- **DefaultClassification** The classification to attach to the `Classification` field of `DetectionResult` if it is left unchanged.
- **AnalysisStels**: Describes to analysts what actions should be taken to make a determination on whether the detected items are true positives.
- **FalsePositives**: Conditions and situations in which known false positives may occur.
- **References**: List of URLs that can be referenced to provide additional context or source material that was used in the development of the detection rule.
- **TestBaseFile**: Local file which contains example telemetry (probably in JSON format) to use when testing against the rule. In our example this is found at `examples/foozle/example.json`.
- **Tests**: A list of test cases to run against the rule logic. This is the `yaml` representation of [`TestCase`](###TestCase). In our example, these tests are parsed and run via `examples/foozle/foozle_test.go`. 

### TestCase
Struct to define the test cases for any given detection rule.

A summary of the fields are below: 
- **It**: A string describing the goals of the test. Think of it like `**It** expects the rule to evaluate to true when x, y, and z happen.`
- **Want**: The `yaml` representation of which [`DetectionResult`](###DetectionResult) fields we specifically expect the rule to return in this test. Your implementation of this should be prepopulated from a default result (like in `examples/foozle/foozle_rules/rules.go`) so you only need to specify fields you expect to be different from the default, resulting in less `yaml` content.
- **Override**: The `yaml` representation of the telemetry that we want the test to use. Your implementation of this should be prepopulated from the default data defined in `TestBaseFile` above, so you only need to specify which fields to change from the default data, resulting in less `yaml` content.

See references to `Override` in `examples/foozle/foozle_test.go` for an example.

---