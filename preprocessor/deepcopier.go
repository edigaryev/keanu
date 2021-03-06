package preprocessor

import (
	"bytes"
	"encoding/gob"
	"gopkg.in/yaml.v2"
)

// This is rather inefficient and error-prone (due to the need to manually register unknown types),
// but nevertheless works flawlessly for yaml.v2 structures, compared to other alternatives.
func deepcopy(dst, src interface{}) error {
	// Register unknown types
	// https://golang.org/pkg/encoding/gob/#Register
	gob.Register(yaml.MapSlice{})
	gob.Register([]interface{}{})

	var tmp bytes.Buffer

	if err := gob.NewEncoder(&tmp).Encode(src); err != nil {
		return err
	}

	return gob.NewDecoder(&tmp).Decode(dst)
}
