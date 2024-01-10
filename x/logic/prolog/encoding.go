package prolog

import (
	"bytes"
	"unicode/utf8"

	"github.com/ichiban/prolog/engine"
	"golang.org/x/net/html/charset"
)

// Decode converts a byte slice from a specified encoding.
// Decode function is the reverse of encode function.
func Decode(bs []byte, label engine.Atom, env *engine.Env) ([]byte, error) {
	switch label {
	case AtomEmpty, AtomText:
		return bs, nil
	case AtomOctet:
		var buffer bytes.Buffer
		for _, b := range bs {
			buffer.WriteRune(rune(b))
		}
		return buffer.Bytes(), nil
	default:
		encoding, _ := charset.Lookup(label.String())
		if encoding == nil {
			return nil, engine.DomainError(ValidCharset(), label, env)
		}
		result, err := encoding.NewDecoder().Bytes(bs)
		if err != nil {
			culprit := BytesToCodepointListTermWithDefault(bs, env)
			return nil, engine.DomainError(ValidEncoding(label.String(), err), culprit, env)
		}
		return result, nil
	}
}

// Encode converts a byte slice to a specified encoding.
//
// In case of:
//   - empty encoding label or 'text', return the original bytes without modification.
//   - 'octet', decode the bytes as unicode code points and return the resulting bytes. If a code point is greater than
//     0xff, an error is returned.
//   - any other encoding label, convert the bytes to the specified encoding.
func Encode(bs []byte, label engine.Atom, env *engine.Env) ([]byte, error) {
	switch label {
	case AtomEmpty, AtomText:
		return bs, nil
	case AtomOctet:
		result := make([]byte, 0, len(bs))
		for i := 0; i < len(bs); {
			runeValue, width := utf8.DecodeRune(bs[i:])

			if runeValue > 0xff {
				culprit := BytesToCodepointListTermWithDefault(bs, env)
				return nil, engine.DomainError(ValidByte(int64(runeValue)), culprit, env)
			}
			result = append(result, byte(runeValue))
			i += width
		}
		return result, nil
	default:
		encoding, _ := charset.Lookup(label.String())
		if encoding == nil {
			return nil, engine.DomainError(ValidCharset(), label, env)
		}
		result, err := encoding.NewEncoder().Bytes(bs)
		if err != nil {
			culprit := BytesToCodepointListTermWithDefault(bs, env)
			return nil, engine.DomainError(ValidEncoding(label.String(), err), culprit, env)
		}
		return result, nil
	}
}
