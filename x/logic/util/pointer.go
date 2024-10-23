package util

import (
	"reflect"
)

// NonZeroOrDefault returns the value of the argument if it is not nil and not zero, otherwise returns the default value.
func NonZeroOrDefault[T any](v, defaultValue T) T {
	if IsZero(v) {
		return defaultValue
	}
	return v
}

// IsZero returns true if the argument is nil or zero.
func IsZero[T any](v T) bool {
	v1 := reflect.ValueOf(v)
	return !v1.IsValid() || v1.IsZero()
}
