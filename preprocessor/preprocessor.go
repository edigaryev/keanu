package preprocessor

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// ErrNeedMap is returned when the input document passed to preprocessor
// does not contain a map as it's "outer" layer.
var ErrNeedMap = errors.New("first document should be a map")

// ErrMatrixNeedsCollection is returned when the matrix modifier
// does not contain a collection (either map or slice) inside.
var ErrMatrixNeedsCollection = errors.New("matrix should contain a collection")

// ErrMatrixNeedsListOfMaps is returned when the matrix modifier contains
// something other than maps (e.g. lists or scalars) as it's items.
var ErrMatrixNeedsListOfMaps = errors.New("matrix with a list can only contain maps as it's items")

// A Preprocessor for YAML files with matrix modifiers.
type Preprocessor struct {
	currentTrees yaml.MapSlice
}

// New constructs preprocessor by loading the first YAML document from the supplied byte slice.
func New(in []byte) (*Preprocessor, error) {
	p := &Preprocessor{}

	err := yaml.Unmarshal(in, &p.currentTrees)
	if err != nil {
		return nil, ErrNeedMap
	}

	return p, nil
}

// NewFromFile constructs preprocessor by loading the first YAML document from the file located at path.
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
	var newParsedTree yaml.MapSlice
	var expanded bool

	if len(p.currentTrees) == 0 {
		return false, nil
	}

	for i := range p.currentTrees {
		var treeToExpand yaml.MapItem
		// deepcopy since expandIfMatrix has side effects
		if err := deepcopy(&treeToExpand, p.currentTrees[i]); err != nil {
			return false, err
		}

		var expandedTrees []yaml.MapItem
		expandedTreesCollector := func(item *yaml.MapItem) (bool, error) {
			newTrees, expandErr := expandIfMatrix(&treeToExpand, item)
			// stop once found any expansion
			if len(newTrees) != 0 {
				expandedTrees = newTrees
				return true, nil
			}
			return false, expandErr
		}

		if err := traverse(&treeToExpand, expandedTreesCollector); err != nil {
			return true, err
		}

		if len(expandedTrees) == 0 {
			newParsedTree = append(newParsedTree, treeToExpand)
		} else {
			newParsedTree = append(newParsedTree, expandedTrees...)
			expanded = true
		}
	}

	p.currentTrees = newParsedTree

	return expanded, nil
}

// Run preprocesses matrix modifiers in the loaded YAML document.
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

// Dump marshals the loaded (and possibly preprocessed) YAML document back.
func (p *Preprocessor) Dump() ([]byte, error) {
	return yaml.Marshal(&p.currentTrees)
}
