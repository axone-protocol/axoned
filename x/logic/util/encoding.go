package util

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"golang.org/x/net/html/charset"
)

// Decode converts a byte slice from a specified encoding.
// Decode function is the reverse of encode function.
func Decode(bs []byte, label string) ([]byte, error) {
	switch label {
	case "", "text":
		return bs, nil
	case "octet":
		var buffer bytes.Buffer
		for _, b := range bs {
			buffer.WriteRune(rune(b))
		}
		return buffer.Bytes(), nil
	default:
		encoding, _ := charset.Lookup(label)
		if encoding == nil {
			return nil, fmt.Errorf("invalid encoding: %s", label)
		}
		return encoding.NewDecoder().Bytes(bs)
	}
}

// Encode converts a byte slice to a specified encoding.
//
// In case of:
//   - empty encoding label or 'text', return the original bytes without modification.
//   - 'octet', decode the bytes as unicode code points and return the resulting bytes. If a code point is greater than
//     0xff, an error is returned.
//   - any other encoding label, convert the bytes to the specified encoding.
func Encode(bs []byte, label string) ([]byte, error) {
	switch label {
	case "", "text":
		return bs, nil
	case "octet":
		result := make([]byte, 0, len(bs))
		for i := 0; i < len(bs); {
			runeValue, width := utf8.DecodeRune(bs[i:])

			if runeValue > 0xff {
				return nil, fmt.Errorf("cannot convert character '%c' to %s", runeValue, label)
			}
			result = append(result, byte(runeValue))
			i += width
		}
		return result, nil
	default:
		encoding, _ := charset.Lookup(label)
		if encoding == nil {
			return nil, fmt.Errorf("invalid encoding: %s", label)
		}
		return encoding.NewEncoder().Bytes(bs)
	}
}
