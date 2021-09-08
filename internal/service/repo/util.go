package repo

import "errors"

// stripNilValues removes empty strings and nil values from map https://github.com/Masterminds/squirrel/issues/66
func stripNilValues(in map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	for k, v := range in {
		if v != nil && v != "" {
			out[k] = v
		}
	}

	if len(out) == 0 {
		return nil, errors.New("provide at least one field to update resource partially")
	}

	return out, nil
}
