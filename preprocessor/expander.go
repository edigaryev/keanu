package preprocessor

import (
	"gopkg.in/yaml.v2"
)

// Callback function to be called by traverse().
//
// Traversal continues until this function decides that
// no more artifacts can be collected into out slice
// and returns true.
type callback func(root *yaml.MapItem, item *yaml.MapItem, out *[]yaml.MapItem) (bool, error)

// Implements preorder traversal of the YAML parse tree.
func traverse(root *yaml.MapItem, item *yaml.MapItem, f callback, out *[]yaml.MapItem) error {
	done, err := f(root, item, out)
	if done {
		return err
	}

	return process(root, item.Value, f, out)
}

// Recursive descent helper for traverse().
func process(root *yaml.MapItem, something interface{}, f callback, out *[]yaml.MapItem) (err error) {
	switch obj := something.(type) {
	case yaml.MapSlice:
		// YAML mapping node
		for i := range obj {
			if err := traverse(root, &obj[i], f, out); err != nil {
				return err
			}
		}
	case yaml.MapItem:
		// YAML mapping KV pair
		err = traverse(root, &obj, f, out)
	case []interface{}:
		// YAML sequence node
		for _, obj := range obj {
			if err := process(root, obj, f, out); err != nil {
				return err
			}
		}
	}

	return
}

// Expands one MapItem if it holds a "matrix" into multiple MapItem's
// and writes result to out.
func expandOneMatrix(root *yaml.MapItem, item *yaml.MapItem, out *[]yaml.MapItem) (bool, error) {
	// Potential "matrix" modifier can only be found in a map
	obj, ok := item.Value.(yaml.MapSlice)
	if !ok {
		return false, nil
	}

	// Split the map into two slices:
	// * nonMatrixSlice contains non-"matrix" items
	// * matrixSlice contains items with the "matrix" key
	var nonMatrixSlice []yaml.MapItem
	var matrixSlice []yaml.MapItem

	for _, sliceItem := range obj {
		if sliceItem.Key == "matrix" {
			matrixSlice = append(matrixSlice, sliceItem)
		} else {
			nonMatrixSlice = append(nonMatrixSlice, sliceItem)
		}
	}

	// Keep going deeper if no "matrix" modifiers were found at this level
	if len(matrixSlice) == 0 {
		return false, nil
	}

	// Take the first "matrix" modifier found and pour the rest (if any)
	// into a slice of non-"matrix" items for further processing
	// by the forthcoming invocations of singlePass().
	matrix, matrixSlice := matrixSlice[0], matrixSlice[1:]
	nonMatrixSlice = append(nonMatrixSlice, matrixSlice...)

	// Extract parametrizations from "matrix" modifier we've selected
	var parametrizations []yaml.MapSlice

	switch obj := matrix.Value.(type) {
	case yaml.MapSlice:
		for _, sliceItem := range obj {
			// Inherit "matrix" siblings
			var tmp yaml.MapSlice
			tmp = append(tmp, nonMatrixSlice...)

			// Generate a single parametrization
			tmp = append(tmp, sliceItem)
			parametrizations = append(parametrizations, tmp)
		}
	case []interface{}:
		for _, listItem := range obj {
			// Inherit "matrix" siblings
			var tmp yaml.MapSlice
			tmp = append(tmp, nonMatrixSlice...)

			// Ensure that matrix with a list contains only maps as it's items
			//
			// This restriction was made purely for simplicity's sake and can be lifted in the future.
			innerSlice, ok := listItem.(yaml.MapSlice)
			if !ok {
				return true, ErrMatrixNeedsListOfMaps
			}

			// Generate a single parametrization
			tmp = append(tmp, innerSlice...)
			parametrizations = append(parametrizations, tmp)
		}
	default:
		// Semantics is undefined for "matrix" modifiers without a collection inside
		return true, ErrMatrixNeedsCollection
	}

	// The Tricky Part™
	//
	// Produces a new diverged root for each parametrization,
	// with a side-effect that the sub-tree of the original root
	// will be overwritten by our parametrization and thus made dirty.
	//
	// However this is fine, because we never re-use the old root anyways
	// and stop processing straight after the parametrization is complete.
	for _, parametrization := range parametrizations {
		item.Value = parametrization

		var divergedRoot yaml.MapItem
		if err := deepcopy(&divergedRoot, *root); err != nil {
			return true, err
		}

		*out = append(*out, divergedRoot)
	}

	return true, nil
}
