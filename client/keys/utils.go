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

type bechKeyOutFn func(k *cryptokeyring.Record) (keys.KeyOutput, error)

func printKeyringRecord(w io.Writer, k *cryptokeyring.Record, bechKeyOut bechKeyOutFn, output string) error {
	ko, err := mkKeyOutput(k, bechKeyOut)
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

	for i, r := range records {
		kko, err := mkKeyOutput(r, keys.MkAccKeyOutput)
		if err != nil {
			return nil, err
		}

		kos[i] = kko
	}

	return kos, nil
}

func mkKeyOutput(record *cryptokeyring.Record, bechKeyOut bechKeyOutFn) (KeyOutput, error) {
	kko, err := bechKeyOut(record)
	if err != nil {
		return KeyOutput{}, err
	}
	pk, err := record.GetPubKey()
	if err != nil {
		return KeyOutput{}, err
	}
	did, _ := util.CreateDIDKeyByPubKey(pk)

	return KeyOutput{
		KeyOutput: kko,
		DID:       did,
	}, nil
}
