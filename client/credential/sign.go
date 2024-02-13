package credential

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	"github.com/okp4/okp4d/x/logic/util"
)

var (
	flagOverwrite = "overwrite"
	flagDate      = "date"
)

var valueDash = "-"

func SignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign [file] [flags]",
		Short: "Sign a W3C Verifiable Credential provided as a file or stdin",
		Long: `Sign a W3C Verifiable Credential;

It will read a verifiable credential from a file (or stdin), sign it, and print the JSON-LD signed credential to stdout.
			`,
		Args: cobra.ExactArgs(1),
		RunE: runSignCmd,
	}
	cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	cmd.Flags().Bool(flagOverwrite, false, "Overwrite existing signatures with a new one. If disabled, new signature will be appended")
	cmd.Flags().String(flagDate, "", "Date of the signature provided in RFC3339 format. If not provided, current time will be used")
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	flags.AddKeyringFlags(cmd.Flags())

	return cmd
}

func runSignCmd(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	kuid, err := cmd.Flags().GetString(flags.FlagFrom)
	if err != nil {
		return fmt.Errorf("%s: %w", flags.FlagFrom, err)
	}
	k, err := fetchKey(clientCtx.Keyring, kuid)
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidType, "%s is not a valid name or address: %v", args[0], err)
	}
	signer := mkKeyringSigner(clientCtx.Keyring, k.Name)

	filename, err := expand(args[0])
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidRequest, "failed to expand filename: %v", err)
	}
	bs, err := readFromFileOrStdin(filename)
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrIO, "failed to read file: %v", err)
	}

	vc, err := loadVerifiableCredential(bs)
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidRequest, "failed to load verifiable credential: %v", err)
	}

	overrideProofs, err := cmd.Flags().GetBool(flagOverwrite)
	if err != nil {
		return fmt.Errorf("%s: %w", flagOverwrite, err)
	}
	if overrideProofs {
		vc.Proofs = nil
	}

	date, err := getFlagAsDate(cmd, flagDate)
	if err != nil {
		return err
	}
	err = signVerifiableCredential(vc, signer, date)
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidRequest, "failed to sign: %v", err)
	}

	marshalled, err := vc.MarshalJSON()
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidRequest, "failed to marshal signed credential: %v", err)
	}
	cmd.Println(string(marshalled))

	return nil
}

type KeyringSigner struct {
	keyring keyring.Keyring
	uid     string
}

// Sign will sign document and return signature.
func (m KeyringSigner) Sign(data []byte) ([]byte, error) {
	bs, _, err := m.keyring.Sign(m.uid, data, signing.SignMode_SIGN_MODE_DIRECT)
	return bs, err
}

func (m KeyringSigner) Alg() string {
	return "unknown"
}

func (m KeyringSigner) PubKey() (cryptotypes.PubKey, error) {
	record, err := m.keyring.Key(m.uid)
	if err != nil {
		return nil, err
	}

	return record.GetPubKey()
}

func (m KeyringSigner) DIDKeyID() (string, error) {
	pk, err := m.PubKey()
	if err != nil {
		return "", err
	}

	return util.CreateDIDKeyIDByPubKey(pk)
}

func mkKeyringSigner(keyring keyring.Keyring, uid string) KeyringSigner {
	return KeyringSigner{keyring: keyring, uid: uid}
}

// readFromFileOrStdin reads content from the given filename or from stdin if "-" is passed as the filename.
func readFromFileOrStdin(filename string) ([]byte, error) {
	if filename == valueDash {
		return io.ReadAll(os.Stdin)
	}

	return os.ReadFile(filename)
}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func fetchKey(kb keyring.Keyring, keyref string) (*keyring.Record, error) {
	k, err := kb.Key(keyref)
	if err == nil || !errorsmod.IsOf(err, sdkerr.ErrIO, sdkerr.ErrKeyNotFound) {
		return k, err
	}

	accAddr, err := sdk.AccAddressFromBech32(keyref)
	if err != nil {
		return k, err
	}

	return kb.KeyByAddress(accAddr)
}

func loadVerifiableCredential(bs []byte) (*verifiable.Credential, error) {
	documentLoader := ld.NewDefaultDocumentLoader(nil)
	return verifiable.ParseCredential(
		bs,
		verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDValidation(),
		verifiable.WithJSONLDDocumentLoader(documentLoader))
}

func signVerifiableCredential(vc *verifiable.Credential, signer KeyringSigner, date time.Time) error {
	documentLoader := ld.NewDefaultDocumentLoader(nil)
	didKeyID, err := signer.DIDKeyID()
	if err != nil {
		return err
	}

	return vc.AddLinkedDataProof(&verifiable.LinkedDataProofContext{
		Created:                 &date,
		SignatureType:           "Ed25519Signature2018",
		Suite:                   ed25519signature2018.New(suite.WithSigner(signer)),
		SignatureRepresentation: verifiable.SignatureProofValue,
		VerificationMethod:      didKeyID,
	}, jsonld.WithDocumentLoader(documentLoader))
}

func getFlagAsDate(cmd *cobra.Command, flag string) (time.Time, error) {
	dateStr, err := cmd.Flags().GetString(flag)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", flag, err)
	}
	if dateStr == "" {
		return time.Now(), nil
	}
	return time.Parse(time.RFC3339, dateStr)
}
