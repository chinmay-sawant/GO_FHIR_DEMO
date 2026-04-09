// Package queryparams provides helpers for flattening query values.
package queryparams

import "net/url"

// FromValues returns a flattened string map using the first value for each key.
func FromValues(values url.Values) map[string]string {
	params := make(map[string]string, len(values))
	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		params[key] = vals[0]
	}
	return params
}
