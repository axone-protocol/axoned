package credential

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2020"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	"github.com/axone-protocol/axoned/v9/x/logic/util"
)

const (
	flagOverwrite = "overwrite"
	flagDate      = "date"
	flagSchemaMap = "schema-map"
	flagPurpose   = "purpose"
)

const (
	symbolStdIn             = "-"
	symbolKeyValueSeparator = "="
	symbolHome              = "~"
)

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
	cmd.Flags().StringSlice(flagSchemaMap, []string{}, fmt.Sprintf(
		"Map original URIs to alternative URIs for resolving JSON-LD schemas. "+
			"Useful for redirecting network-based URIs to local filesystem paths or "+
			"other URIs. Each mapping should be in the format 'originalURI=alternativeURI'. "+
			"Multiple mappings can be provided by repeating the flag. Example usage: "+
			"--%[1]s originalURI1=alternativeURI1 --%[1]s originalURI2=alternativeURI2",
		flagSchemaMap))
	cmd.Flags().String(flagPurpose, "assertionMethod", "Proof that describes credential purpose, helps prevent it from being "+
		"misused for some other purpose. Example of commonly used proof purpose values:  "+
		"authentication, assertionMethod, keyAgreement, capabilityDelegation, capabilityInvocation.")

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

	filename, err := expandPath(args[0])
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidRequest, "failed to expand filename: %v", err)
	}
	bs, err := readFromFileOrStdin(filename)
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrIO, "failed to read file: %v", err)
	}

	schemaMap, err := parseStringSliceAsMap(cmd, flagSchemaMap)
	if err != nil {
		return err
	}
	documentLoader := newDocumentLoader(schemaMap)
	vc, err := loadVerifiableCredential(documentLoader, bs)
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
	date, err := parseStringAsDate(cmd, flagDate)
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidType, "%s is not a valid date: %v", flagDate, err)
	}
	purpose, err := cmd.Flags().GetString(flagPurpose)
	if err != nil {
		return errorsmod.Wrapf(sdkerr.ErrInvalidType, "%s is not a valid string: %v", flagPurpose, err)
	}
	err = signVerifiableCredential(documentLoader, vc, signer, date, purpose)
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

// mappedDocumentLoader customizes the JSON-LD document loading process
// by mapping URIs to alternative URIs, utilizing a delegate loader for the actual loading.
type mappedDocumentLoader struct {
	schemaMap map[string]string
	delegate  ld.DocumentLoader
}

func (m mappedDocumentLoader) LoadDocument(url string) (*ld.RemoteDocument, error) {
	if mapped, ok := m.schemaMap[url]; ok {
		expanded, err := expandPath(mapped)
		if err != nil {
			return nil, fmt.Errorf("failed to expand mapped URI: %w", err)
		}
		url = expanded
	}

	return m.delegate.LoadDocument(url)
}

// newDocumentLoader returns a JSON-LD document loader that can be used to load schemas, with the provided schema map.
// The loader will first find schemas in the provided map for a given URI, and use it if found. In any case it will
// use the default document loader to load the schema given the URI.
func newDocumentLoader(schemaMap map[string]string) ld.DocumentLoader {
	return &mappedDocumentLoader{
		schemaMap: schemaMap,
		delegate:  ld.NewDefaultDocumentLoader(nil),
	}
}

func mkKeyringSigner(keyring keyring.Keyring, uid string) KeyringSigner {
	return KeyringSigner{keyring: keyring, uid: uid}
}

// readFromFileOrStdin reads content from the given filename or from stdin if "-" is passed as the filename.
func readFromFileOrStdin(filename string) ([]byte, error) {
	if filename == symbolStdIn {
		return io.ReadAll(os.Stdin)
	}

	return os.ReadFile(filename)
}

// expandPath expands the given path, replacing the "~" symbol with the user's home directory.
func expandPath(path string) (string, error) {
	if len(path) == 0 || !strings.HasPrefix(path, symbolHome) {
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

func loadVerifiableCredential(documentLoader ld.DocumentLoader, bs []byte) (*verifiable.Credential, error) {
	return verifiable.ParseCredential(
		bs,
		verifiable.WithDisabledProofCheck(),
		verifiable.WithJSONLDValidation(),
		verifiable.WithJSONLDDocumentLoader(documentLoader))
}

func signVerifiableCredential(
	documentLoader ld.DocumentLoader, vc *verifiable.Credential, signer KeyringSigner, date time.Time, purpose string,
) error {
	didKeyID, err := signer.DIDKeyID()
	if err != nil {
		return err
	}

	pubKey, err := signer.PubKey()
	if err != nil {
		return err
	}

	switch pubKey.(type) {
	case *ed25519.PubKey:
		return vc.AddLinkedDataProof(&verifiable.LinkedDataProofContext{
			Created:                 &date,
			SignatureType:           "Ed25519Signature2020",
			Suite:                   ed25519signature2020.New(suite.WithSigner(signer)),
			SignatureRepresentation: verifiable.SignatureProofValue,
			VerificationMethod:      didKeyID,
			Purpose:                 purpose,
		}, jsonld.WithDocumentLoader(documentLoader))
	case *secp256k1.PubKey:
		return vc.AddLinkedDataProof(&verifiable.LinkedDataProofContext{
			Created:                 &date,
			SignatureType:           "EcdsaSecp256k1Signature2019",
			Suite:                   ecdsasecp256k1signature2019.New(suite.WithSigner(signer)),
			SignatureRepresentation: verifiable.SignatureJWS,
			VerificationMethod:      didKeyID,
			Purpose:                 purpose,
		}, jsonld.WithDocumentLoader(documentLoader))
	default:
		return fmt.Errorf("invalid pubkey type: %s; expected oneof %+q",
			pubKey.Type(), []string{(&ed25519.PubKey{}).Type(), (&secp256k1.PubKey{}).Type()})
	}
}

func parseStringAsDate(cmd *cobra.Command, flag string) (time.Time, error) {
	dateStr, err := cmd.Flags().GetString(flag)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", flag, err)
	}
	if dateStr == "" {
		return time.Now(), nil
	}
	return time.Parse(time.RFC3339, dateStr)
}

func parseStringSliceAsMap(cmd *cobra.Command, flag string) (map[string]string, error) {
	mappings, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", flag, err)
	}

	schemaMap := make(map[string]string, len(mappings))
	for _, mapping := range mappings {
		parts := strings.SplitN(mapping, symbolKeyValueSeparator, 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid mapping: %s", mapping)
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		schemaMap[key] = value
	}

	return schemaMap, nil
}
