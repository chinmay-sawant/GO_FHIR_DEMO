package utils //nolint:revive

import "time"

// CreateStringPtr returns a pointer to the given string value.
func CreateStringPtr(s string) *string {
	return &s
}

// CreateTimePtr returns a pointer to a time value parsed from string.
func CreateTimePtr(t string) *time.Time {
	var ft time.Time
	if err := ft.UnmarshalJSON([]byte(`"` + t + `"`)); err != nil {
		return nil
	}
	return &ft
}
