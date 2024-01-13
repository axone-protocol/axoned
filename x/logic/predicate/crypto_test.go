//nolint:gocognit,lll
package predicate

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
)

func TestCryptoOperations(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program     string
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				program: `test(Hex) :- crypto_data_hash('hello world', Hash, []), hex_bytes(Hex, Hash).`,
				query:   `test(Hex).`,
				wantResult: []types.TermResults{{
					"Hex": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				}},
				wantSuccess: true,
			},
			{
				program:     `test :- crypto_data_hash('hello world', Hash, []), hex_bytes('b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9', Hash).`,
				query:       `test.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				program: `test(Hex) :- crypto_data_hash('hello world', Hash, [algorithm(sha256)]), hex_bytes(Hex, Hash).`,
				query:   `test(Hex).`,
				wantResult: []types.TermResults{{
					"Hex": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				}},
				wantSuccess: true,
			},
			{
				program: `test(Hex) :- crypto_data_hash('hello world', Hash, [encoding(utf8)]), hex_bytes(Hex, Hash).`,
				query:   `test(Hex).`,
				wantResult: []types.TermResults{{
					"Hex": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				}},
				wantSuccess: true,
			},
			{
				program: `test(Hex) :- crypto_data_hash('68656c6c6f20776f726c64', Hash, [encoding(hex)]), hex_bytes(Hex, Hash).`,
				query:   `test(Hex).`,
				wantResult: []types.TermResults{{
					"Hex": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				}},
				wantSuccess: true,
			},
			{
				program: `test(Hex) :- crypto_data_hash([104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100], Hash, [algorithm(sha256),encoding(octet)]), hex_bytes(Hex, Hash).`,
				query:   `test(Hex).`,
				wantResult: []types.TermResults{{
					"Hex": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				}},
				wantSuccess: true,
			},
			{
				query:       ` crypto_data_hash('hello world', Hash, [algorithm(cheh)]).`,
				wantError:   fmt.Errorf("error(type_error(hash_algorithm,cheh),crypto_data_hash/3)"),
				wantSuccess: false,
			},
			{
				program: `test(Hex) :- crypto_data_hash('hello world', Hash, [algorithm(sha512)]), hex_bytes(Hex, Hash).`,
				query:   `test(Hex).`,
				wantResult: []types.TermResults{{
					"Hex": "'309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f'",
				}},
				wantSuccess: true,
			},
			{
				program: `test(Hex) :- crypto_data_hash('hello world', Hash, [algorithm(md5)]), hex_bytes(Hex, Hash).`,
				query:   `test(Hex).`,
				wantResult: []types.TermResults{{
					"Hex": "'5eb63bbbe01eeed093cb22bb8f5acdc3'",
				}},
				wantSuccess: true,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewLightInterpreterMust(ctx)
						interpreter.Register3(engine.NewAtom("crypto_data_hash"), CryptoDataHash)
						interpreter.Register2(engine.NewAtom("hex_bytes"), HexBytes)

						err := interpreter.Compile(ctx, tc.program)
						So(err, ShouldEqual, nil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldEqual, nil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []types.TermResults
									for sols.Next() {
										m := types.TermResults{}
										err := sols.Scan(m)
										So(err, ShouldEqual, nil)

										got = append(got, m)
									}
									if tc.wantError != nil {
										So(sols.Err(), ShouldNotEqual, nil)
										So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
									} else {
										So(sols.Err(), ShouldEqual, nil)

										if tc.wantSuccess {
											So(len(got), ShouldBeGreaterThan, 0)
											So(len(got), ShouldEqual, len(tc.wantResult))
											for iGot, resultGot := range got {
												for varGot, termGot := range resultGot {
													So(testutil.ReindexUnknownVariables(termGot), ShouldEqual, tc.wantResult[iGot][varGot])
												}
											}
										} else {
											So(len(got), ShouldEqual, 0)
										}
									}
								})
							})
						})
					})
				})
			})
		}
	})
}

func TestXVerify(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program     string
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			// ed25519
			{ // All good
				program: `verify :-
			hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
			hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
			hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
			eddsa_verify(PubKey, Msg, Sig, [encoding(octet), type(ed25519)]).`,
				query:       `verify.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{ // All good with hex encoding
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				eddsa_verify(PubKey, '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Sig, encoding(hex)).`,
				query:       `verify.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{ // Wrong Msg
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9e', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				eddsa_verify(PubKey, Msg, Sig, encoding(octet)).`,
				query:       `verify.`,
				wantSuccess: false,
			},
			{ // Wrong public key
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5b5b', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				eddsa_verify(PubKey, Msg, Sig, encoding(octet)).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError: fmt.Errorf("error(syntax_error([%s]),unknown)",
					strings.Join(strings.Split("ed25519: bad public key length: 33", ""), ",")),
			},
			{ // Wrong signature
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff', Sig),
				eddsa_verify(PubKey, Msg, Sig, encoding(octet)).`,
				query:       `verify.`,
				wantSuccess: false,
			},
			{ // Unsupported algo
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				eddsa_verify(PubKey, Msg, Sig, [encoding(octet), type(foo)]).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(cryptographic_algorithm,foo),eddsa_verify/4)"),
			},
			// ECDSA - secp256r1
			{
				// All good
				program: `verify :-
				hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
				hex_bytes('e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3be', Msg),
				hex_bytes('30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8ed7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e47', Sig),
				ecdsa_verify(PubKey, Msg, Sig, [encoding(octet), type(secp256r1)]).`,
				query:       `verify.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{ // Invalid secp signature
				program: `verify :-
				hex_bytes('0213c8426be471e55506f7ce4f7df557', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				ecdsa_verify(PubKey, Msg, Sig, encoding(octet)).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError: fmt.Errorf("error(syntax_error([%s]),unknown)",
					strings.Join(strings.Split("failed to parse compressed public key (first 10 bytes): 0213c8426be471e55506", ""), ",")),
			},
			{ // Unsupported algo
				program: `verify :-
				hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				ecdsa_verify(PubKey, Msg, Sig, [encoding(octet), type(foo)]).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(cryptographic_algorithm,foo),ecdsa_verify/4)"),
			},
			{
				// Wrong msg
				program: `verify :-
				hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
				hex_bytes('e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3bf', Msg),
				hex_bytes('30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8ed7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e47', Sig),
				ecdsa_verify(PubKey, Msg, Sig, encoding(octet)).`,
				query:       `verify.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: false,
			},
			{
				// Wrong signature
				program: `verify :-
				hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
				hex_bytes('e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3be', Msg),
				hex_bytes('30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8ed7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e48', Sig),
				ecdsa_verify(PubKey, Msg, Sig, encoding(octet)).`,
				query:       `verify.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: false,
			},
			// ECDSA - secp256k1
			{
				// All good
				program: `verify :-
				hex_bytes('026b5450187ee9c63ba9e42cb6018d8469c903aca116178e223de76e49fe63b71c', PubKey),
				hex_bytes('dece063885d3648078f903b6a3e8989f649dc3368cd9c8d69755ed9dcb6a0995', Msg),
				hex_bytes('304402201448201bb4408549b0997f4b9ad9ed36f3cf8bb9c433fc7f3ba48c6b6e39476e022053f7d056f7ffeab9a79f3a36bc2ba969ddd530a3a1495d1ed7bba00039820223', Sig),
				ecdsa_verify(PubKey, Msg, Sig, [encoding(octet), type(secp256k1)]).`,
				query:       `verify.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				// Wrong signature
				program: `verify :-
				hex_bytes('026b5450187ee9c63ba9e42cb6018d8469c903aca116178e223de76e49fe63b71c', PubKey),
				hex_bytes('dece063885d3648078f903b6a3e8989f649dc3368cd9c8d69755ed9dcb6a0996', Msg),
				hex_bytes('304402201448201bb4408549b0997f4b9ad9ed36f3cf8bb9c433fc7f3ba48c6b6e39476e022053f7d056f7ffeab9a79f3a36bc2ba969ddd530a3a1495d1ed7bba00039820223', Sig),
				ecdsa_verify(PubKey, Msg, Sig, [encoding(octet), type(secp256k1)]).`,
				query:       `verify.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: false,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewLightInterpreterMust(ctx)
						interpreter.Register2(engine.NewAtom("hex_bytes"), HexBytes)
						interpreter.Register4(engine.NewAtom("eddsa_verify"), EDDSAVerify)
						interpreter.Register4(engine.NewAtom("ecdsa_verify"), ECDSAVerify)

						err := interpreter.Compile(ctx, tc.program)
						So(err, ShouldBeNil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldEqual, nil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []types.TermResults
									for sols.Next() {
										m := types.TermResults{}
										err := sols.Scan(m)
										So(err, ShouldEqual, nil)

										got = append(got, m)
									}
									if tc.wantError != nil {
										So(sols.Err(), ShouldNotEqual, nil)
										So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
									} else {
										So(sols.Err(), ShouldEqual, nil)

										if tc.wantSuccess {
											So(len(got), ShouldBeGreaterThan, 0)
											So(len(got), ShouldEqual, len(tc.wantResult))
											for iGot, resultGot := range got {
												for varGot, termGot := range resultGot {
													So(testutil.ReindexUnknownVariables(termGot), ShouldEqual, tc.wantResult[iGot][varGot])
												}
											}
										} else {
											So(len(got), ShouldEqual, 0)
										}
									}
								})
							})
						})
					})
				})
			})
		}
	})
}
