package dac

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadYML[overrideType any](myFS fs.FS) ([]RuleMeta[overrideType], error) {
	var retErr error
	rules := []RuleMeta[overrideType]{}
	fs.WalkDir(myFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {

			b, err := fs.ReadFile(myFS, path)
			if err != nil {
				retErr = err
				return fs.SkipAll
			}

			newRule, err := parseYML[overrideType](b)
			if err != nil {
				retErr = err
				return fs.SkipAll
			}

			matchCallCount := strings.Count(string(b), "is_match:")
			if matchCallCount < len(newRule.Tests) {
				retErr = fmt.Errorf("%s: require is_match to be called for each test case. currently missing %d", path, len(newRule.Tests)-matchCallCount)
				return fs.SkipAll
			}

			rules = append(rules, newRule)
		}

		return nil
	})

	return rules, retErr
}

func parseYML[overrideType any](ymlBytes []byte) (RuleMeta[overrideType], error) {
	ruleMeta := &RuleMeta[overrideType]{}
	err := yaml.Unmarshal(ymlBytes, ruleMeta)
	if err != nil {
		return *ruleMeta, err
	}

	return *ruleMeta, nil
}
