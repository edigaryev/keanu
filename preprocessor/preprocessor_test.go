package preprocessor_test

import (
	"github.com/edigaryev/keanu/preprocessor"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"testing"
)

// Retrieves the specified document (where the first document index is 1) from YAML file located at path.
func getDocument(t *testing.T, path string, index int) string {
	newPath := filepath.Join("..", "test", path)

	file, err := os.Open(newPath)
	if err != nil {
		t.Fatalf("%s: %s", newPath, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var document yaml.MapSlice

	for i := 0; i < index; i++ {
		if err := decoder.Decode(&document); err != nil {
			t.Fatalf("%s: %s", newPath, err)
		}
	}

	bytes, err := yaml.Marshal(document)
	if err != nil {
		t.Fatalf("%s: %s", newPath, err)
	}

	return string(bytes)
}

// Unmarshals YAML specified by yamlText to a yaml.MapSlice to simplify comparison.
func yamlAsStruct(t *testing.T, yamlText string) (result yaml.MapSlice) {
	if err := yaml.Unmarshal([]byte(yamlText), &result); err != nil {
		t.Fatal(err)
	}

	return
}

var goodCases = []string{
	// Just a document with an empty map
	"empty.yaml",
	// Examples from Google Docs document
	"gdoc-example1.yaml",
	"gdoc-example2.yaml",
	// Real examples from https://cirrus-ci.org/guide/writing-tasks/
	"real1.yaml",
	"real2.yaml",
	"real3.yaml",
	// Real examples from https://cirrus-ci.org/examples/
	"real4.yaml",
	"real5.yaml",
	// Encountered regressions
	"simple-slice.yaml",
	"simple-list.yaml",
	"doubly-nested-balanced.yaml",
	"doubly-nested-unbalanced.yaml",
	"matrix-inside-of-a-list-of-lists.yaml",
	"matrix-siblings.yaml",
	"multiple-matrices-on-the-same-level.yaml",
}

var badCases = []struct {
	File  string
	Error error
}{
	{"bad-not-a-map.yaml", preprocessor.ErrNeedMap},
	{"bad-matrix-without-collection.yaml", preprocessor.ErrMatrixNeedsCollection},
	{"bad-matrix-with-list-of-scalars.yaml", preprocessor.ErrMatrixNeedsListOfMaps},
}

// Helper for instantiating a Preprocessor, running it and getting result.
//
// Takes run parameter which can be set to false to ignore
// the preprocessing (aka matrix expansion) step.
func runPreprocessor(input string, run bool) (string, error) {
	pp, err := preprocessor.New([]byte(input))
	if err != nil {
		return "", err
	}

	if run {
		err = pp.Run()
		if err != nil {
			return "", err
		}
	}

	outputBytes, err := pp.Dump()
	if err != nil {
		return "", err
	}

	return string(outputBytes), nil
}

// Ensures that preprocessing works as expected.
func TestGoodCases(t *testing.T) {
	for _, goodFile := range goodCases {
		input := getDocument(t, goodFile, 1)
		expectedOutput := getDocument(t, goodFile, 2)

		output, err := runPreprocessor(input, true)
		if err != nil {
			t.Error(err)
			continue
		}

		t.Run(goodFile, func(t *testing.T) {
			diff := deep.Equal(yamlAsStruct(t, expectedOutput), yamlAsStruct(t, output))
			if diff != nil {
				t.Error("found difference")
				for _, d := range diff {
					t.Log(d)
				}
			}
		})
	}
}

// Ensures that we return correct errors for expected edge-cases.
func TestBadCases(t *testing.T) {
	for _, badCase := range badCases {
		newPath := filepath.Join("..", "test", badCase.File)
		p, err := preprocessor.NewFromFile(newPath)
		if err != nil {
			assert.Equal(t, badCase.Error, err)
			continue
		}

		if err = p.Run(); err != nil {
			assert.Equal(t, badCase.Error, err)
			continue
		}

		_, err = p.Dump()
		if err != nil {
			assert.Equal(t, badCase.Error, err)
		}
	}
}

// Ensures that no change is made to the original YAML
// with the absence of Run() invocation.
func TestNoop(t *testing.T) {
	for _, testFile := range goodCases {
		input := getDocument(t, testFile, 1)
		output, err := runPreprocessor(input, false)
		if err != nil {
			t.Fatal(err)
		}

		// Preprocessed YAML document should be identical to the input YAML document
		assert.Equal(t, yamlAsStruct(t, input), yamlAsStruct(t, output))
	}
}

// Ensures that good cases are parsed without errors
// when using NewFromFile constructor (the actual results
// of preprocessing are examined in TestGoodCases).
func TestFromFile(t *testing.T) {
	for _, testFile := range goodCases {
		newPath := filepath.Join("..", "test", testFile)
		p, err := preprocessor.NewFromFile(newPath)
		if err != nil {
			t.Error(err)
			continue
		}
		if err := p.Run(); err != nil {
			t.Error(err)
			continue
		}
		_, err = p.Dump()
		if err != nil {
			t.Error(err)
			continue
		}
	}
}
