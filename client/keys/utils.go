package keys

import (
	"encoding/json"
	"fmt"
	"io"

	"sigs.k8s.io/yaml"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/okp4/okp4d/x/logic/util"
)

// KeyOutput is the output format for keys when listing them.
// It is an improved copy of the KeyOutput from the keys module (github.com/cosmos/cosmos-sdk/client/keys/types.go).
type KeyOutput struct {
	keys.KeyOutput
	DID string `json:"did,omitempty" yaml:"did"`
}

// bechKeyOutFn is a function that converts a record into a KeyOutput and returns an error if it fails.
type bechKeyOutFn func(k *cryptokeyring.Record) (KeyOutput, error)

func printKeyringRecord(w io.Writer, k *cryptokeyring.Record, bechKeyOut bechKeyOutFn, output string) error {
	ko, err := bechKeyOut(k)
	if err != nil {
		return err
	}

	switch output {
	case flags.OutputFormatText:
		if err := printTextRecords(w, []KeyOutput{ko}); err != nil {
			return err
		}

	case flags.OutputFormatJSON:
		out, err := json.Marshal(ko)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, string(out)); err != nil {
			return err
		}
	}

	return nil
}

func printKeyringRecords(w io.Writer, records []*cryptokeyring.Record, output string) error {
	kos, err := mkKeysOutput(records)
	if err != nil {
		return err
	}

	switch output {
	case flags.OutputFormatText:
		if err := printTextRecords(w, kos); err != nil {
			return err
		}

	case flags.OutputFormatJSON:
		out, err := json.Marshal(kos)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, "%s", out); err != nil {
			return err
		}
	}

	return nil
}

func printTextRecords(w io.Writer, kos []KeyOutput) error {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, string(out)); err != nil {
		return err
	}

	return nil
}

func mkKeysOutput(records []*cryptokeyring.Record) ([]KeyOutput, error) {
	kos := make([]KeyOutput, len(records))
	bechKeyOut := toBechKeyOutFn(keys.MkAccKeyOutput)
	for i, r := range records {
		kko, err := bechKeyOut(r)
		if err != nil {
			return nil, err
		}

		kos[i] = kko
	}

	return kos, nil
}

// toBechKeyOutFn converts a function that returns a KeyOutput and an error into a function that returns
// an extended KeyOutput and an error.
func toBechKeyOutFn(in func(k *cryptokeyring.Record) (keys.KeyOutput, error)) bechKeyOutFn {
	return func(k *cryptokeyring.Record) (KeyOutput, error) {
		ko, err := in(k)
		if err != nil {
			return KeyOutput{}, err
		}

		pk, err := k.GetPubKey()
		if err != nil {
			return KeyOutput{}, err
		}
		did, _ := util.CreateDIDKeyByPubKey(pk)

		return KeyOutput{
			KeyOutput: ko,
			DID:       did,
		}, nil
	}
}
