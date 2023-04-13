package util

import "reflect"

// DerefOrDefault returns the value of the pointer if it is not nil, otherwise returns the default value.
func DerefOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// NonZeroOrDefault returns the value of the argument if it is not nil and not zero, otherwise returns the default value.
func NonZeroOrDefault[T any](v, defaultValue T) T {
	v1 := reflect.ValueOf(v)
	if v1.IsValid() && !v1.IsZero() {
		return v
	}
	return defaultValue
}
