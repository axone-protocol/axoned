package codec

import (
	"bytes"
	"encoding/hex"
	"unicode/utf8"

	"github.com/axone-protocol/prolog/v3/engine"

	sdkbech32 "github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
)

var (
	atomOK = engine.NewAtom("ok")

	errInvalidBech32 = prolog.AtomError.Apply(engine.NewAtom("invalid_bech32"))
	errInvalidBytes  = prolog.AtomError.Apply(engine.NewAtom("invalid_bytes"))
)

// bech32Codec implements the Codec interface for Bech32 encoding/decoding.
type bech32Codec struct{}

func init() {
	Register(&bech32Codec{})
}

// Name returns the codec identifier.
func (c *bech32Codec) Name() string {
	return codecNameBech32
}

// Decode processes a Bech32 decode request.
//
// Request format: decode <bech32_string>
// Response: ok(HRP-Bytes) or error(invalid_bech32)
//
// Where:
//   - HRP is the human-readable part as an atom
//   - Bytes is a list of byte integers [0..255]
func (c *bech32Codec) Decode(input []byte) engine.Term {
	tokens, ok := splitBech32Payload(input)
	if !ok || len(tokens) != 1 {
		return errInvalidRequest
	}

	bech32Text := tokens[0]
	hrp, data, err := sdkbech32.DecodeAndConvert(string(bech32Text))
	if err != nil {
		return errInvalidBech32
	}

	return atomOK.Apply(prolog.AtomPair.Apply(
		engine.NewAtom(hrp),
		prolog.BytesToByteListTerm(data),
	))
}

// Encode processes a Bech32 encode request.
//
// Request format: encode <hrp> <hex_bytes>
// Response: ok(Bech32String) or error(invalid_bech32|invalid_bytes)
//
// Where:
//   - hrp is the human-readable part
//   - hex_bytes is the data in hexadecimal format
//
// The codec validates the hex payload locally and delegates Bech32 formatting
// to the SDK.
func (c *bech32Codec) Encode(input []byte) engine.Term {
	tokens, ok := splitBech32Payload(input)
	if !ok || len(tokens) != 2 {
		return errInvalidRequest
	}

	hrpText := tokens[0]
	hexText := tokens[1]

	data, err := decodeHex(hexText)
	if err != nil {
		return errInvalidBytes
	}

	// With fixed ConvertBits parameters (8 -> 5 with padding) and
	// validated base16 input, ConvertAndEncode cannot fail for this call site.
	bech32Address, _ := sdkbech32.ConvertAndEncode(string(hrpText), data)

	return atomOK.Apply(engine.NewAtom(bech32Address))
}

func splitBech32Payload(input []byte) ([][]byte, bool) {
	input = bytes.TrimRight(input, "\r\n")
	input = bytes.Trim(input, " ")
	if len(input) == 0 || !utf8.Valid(input) {
		return nil, false
	}

	for _, b := range input {
		switch {
		case b == ' ':
			continue
		case b == '\t', b == '\n', b == '\r':
			return nil, false
		case b < 0x20 || b == 0x7f:
			return nil, false
		}
	}

	return splitSpaceTokens(input), true
}

func splitSpaceTokens(line []byte) [][]byte {
	tokens := make([][]byte, 0, 3)

	for start := 0; start < len(line); {
		for start < len(line) && line[start] == ' ' {
			start++
		}
		if start >= len(line) {
			break
		}

		end := start
		for end < len(line) && line[end] != ' ' {
			end++
		}
		tokens = append(tokens, line[start:end])
		start = end
	}

	return tokens
}

// decodeHex is a helper to decode hexadecimal text to bytes.
func decodeHex(hexText []byte) ([]byte, error) {
	if len(hexText)%2 != 0 {
		return nil, hex.ErrLength
	}

	data := make([]byte, hex.DecodedLen(len(hexText)))
	if _, err := hex.Decode(data, hexText); err != nil {
		return nil, err
	}

	return data, nil
}
