package util

import (
	"errors"

	"golang.org/x/net/html/charset"
)

var ErrInvalidCharset = errors.New("invalid charset")

// Decode converts a byte slice from a specified encoding to a string.
// Decode function is the reverse of encode function.
func Decode(bs []byte, label string) (string, error) {
	encoding, _ := charset.Lookup(label)
	if encoding == nil {
		return "", ErrInvalidCharset
	}
	result, err := encoding.NewDecoder().Bytes(bs)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// Encode converts a string to a slice of bytes in a specified encoding.
func Encode(str string, label string) ([]byte, error) {
	encoding, _ := charset.Lookup(label)
	if encoding == nil {
		return nil, ErrInvalidCharset
	}
	return encoding.NewEncoder().Bytes([]byte(str))
}
