package provider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/3rein/terraform-provider-data-utils/internal/config"

	"gopkg.in/yaml.v3"

	"io/fs"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type TestCase struct {
	Name   string `json:"-" yaml:"-"`
	Format string `json:"-" yaml:"-"`
	Config struct {
		Format                      string `json:"format" yaml:"format"`
		WithOverride                bool   `json:"with_override" yaml:"with_override"`
		WithAppendSlice             bool   `json:"with_append_slice" yaml:"with_append_slice"`
		WithOverwriteWithEmptyValue bool   `json:"with_overwrite_with_empty_value" yaml:"with_overwrite_with_empty_value"`
		WithSliceDeepCopy           bool   `json:"with_slice_deep_copy" yaml:"with_slice_deep_copy"`
	} `json:"config" yaml:"config"`
	Inputs []map[string]interface{} `json:"inputs" yaml:"inputs"`
	Output map[string]interface{}   `json:"output" yaml:"output"`
}

func readTestCases(t *testing.T, testCaseFolder string) []*TestCase {
	results := make([]*TestCase, 0)
	err := filepath.Walk(testCaseFolder, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			t.Logf("skipping unsupported file type: %s", path)
			return nil
		}

		t.Logf("Parsing test case file file: %s", path)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		testCase := TestCase{
			Name: info.Name(),
		}
		switch strings.ToLower(ext) {
		case ".json":
			testCase.Format = "JSON"
			err = json.Unmarshal(b, &testCase)
		case ".yaml":
			testCase.Format = "YAML"
			err = yaml.Unmarshal(b, &testCase)
		case ".yml":
			testCase.Format = "YAML"
			err = yaml.Unmarshal(b, &testCase)

		}
		if err != nil {
			return fmt.Errorf(" %w: failed to unmarshal input file: %s", err, path)
		}
		results = append(results, &testCase)
		return nil
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return results

}

func TestHappyPathDeepMerges(t *testing.T) {
	testCases := readTestCases(t, "../../resources/success")

	for _, tc := range testCases {
		t.Run(tc.Name, func(it *testing.T) {
			runTestCase(it, tc)

		})
	}
}

func runTestCase(t *testing.T, testCase *TestCase) {
	rawInputs := make([]interface{}, 0)
	for _, input := range testCase.Inputs {
		switch testCase.Format {
		case "JSON":
			rawInput, err := json.Marshal(input)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			rawInputs = append(rawInputs, string(rawInput))
		case "YAML":
			rawInput, err := yaml.Marshal(input)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			rawInputs = append(rawInputs, string(rawInput))
		}
	}
	t.Logf("Format: %s", testCase.Format)
	t.Logf("WithOverride: %t", testCase.Config.WithOverride)
	t.Logf("WithOverwriteWithEmptyValue: %t", testCase.Config.WithOverwriteWithEmptyValue)
	t.Logf("WithAppendSlice: %t", testCase.Config.WithAppendSlice)
	t.Logf("WithSliceDeepCopy: %t", testCase.Config.WithSliceDeepCopy)

	merged, diagnostics := merge(rawInputs, &config.Config{
		Format:                      testCase.Format,
		WithSliceDeepCopy:           testCase.Config.WithSliceDeepCopy,
		WithAppendSlice:             testCase.Config.WithAppendSlice,
		WithOverwriteWithEmptyValue: testCase.Config.WithOverwriteWithEmptyValue,
		WithOverride:                testCase.Config.WithOverride,
	})

	if diagnostics.HasError() {
		for _, d := range diagnostics {
			t.Error(fmt.Errorf(d.Detail))
		}
		t.FailNow()
	}

	rawOutput := make(map[string]interface{})
	switch testCase.Format {
	case "JSON":
		err := json.Unmarshal([]byte(merged), &rawOutput)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

	case "YAML":
		err := yaml.Unmarshal([]byte(merged), &rawOutput)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

	}
	fmt.Printf("Expected: %v\n", testCase.Output)
	fmt.Printf("Actual:   %v\n", rawOutput)

	if !reflect.DeepEqual(testCase.Output, rawOutput) {
		t.Error(fmt.Errorf("output does not match expected value"))
		t.FailNow()
	}
}
