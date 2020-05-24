package preprocessor

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var ErrNeedMap = errors.New("first document should be a map")
var ErrMatrixNeedsCollection = errors.New("matrix should contain a collection")
var ErrMatrixNeedsListOfMaps = errors.New("matrix with a list can only contain maps as it's items")

type Preprocessor struct {
	firstDocument yaml.MapSlice
}

// Constructs a new preprocessor by loading the first YAML document from the supplied byte slice.
func New(in []byte) (*Preprocessor, error) {
	p := &Preprocessor{}

	err := yaml.Unmarshal(in, &p.firstDocument)
	if err != nil {
		return nil, ErrNeedMap
	}

	return p, nil
}

// Constructs a new preprocessor by loading the first YAML document from the file located at path.
func NewFromFile(path string) (*Preprocessor, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	in, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return New(in)
}

// Recursively processes each "outer" map key of the loaded YAML document
// in an attempt to produce multiple keys as a result of matrix expansion.
func (p *Preprocessor) singlePass() (bool, error) {
	var newParseTree yaml.MapSlice
	var expanded bool

	if len(p.firstDocument) == 0 {
		return false, nil
	}

	for _, mapItem := range p.firstDocument {
		var out []yaml.MapItem
		if err := traverse(&mapItem, &mapItem, expandOneMatrix, &out); err != nil {
			return true, err
		}

		if len(out) == 0 {
			newParseTree = append(newParseTree, mapItem)
		} else {
			newParseTree = append(newParseTree, out...)
			expanded = true
		}
	}

	p.firstDocument = newParseTree

	return expanded, nil
}

// Preprocesses matrix modifiers in the loaded YAML document.
func (p *Preprocessor) Run() error {
	for {
		expanded, err := p.singlePass()
		if err != nil {
			return err
		}

		// Consider the preprocessing done once singlePass() stops expanding the document
		// (which means no "matrix" modifier were to be found)
		if !expanded {
			return nil
		}
	}
}

// Marshals the loaded (and possibly preprocessed) YAML document back.
func (p *Preprocessor) Dump() ([]byte, error) {
	return yaml.Marshal(&p.firstDocument)
}
