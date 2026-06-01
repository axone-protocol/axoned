//nolint:gocognit,lll
package predicate

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"
	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/testutil"
)

func TestXVerify(t *testing.T) {
	const (
		badPublicKeyLength     = "e,d,2,5,5,1,9,:, ,b,a,d, ,p,u,b,l,i,c, ,k,e,y, ,l,e,n,g,t,h,:, ,3,3"
		failedToParsePublicKey = "f,a,i,l,e,d, ,t,o, ,p,a,r,s,e, ,c,o,m,p,r,e,s,s,e,d, ,p,u,b,l,i,c, ,k,e,y, ,(,f,i,r,s,t, ,1,0, ,b,y,t,e,s,),:, ,0,2,1,3,c,8,4,2,6,b,e,4,7,1,e,5,5,5,0,6"
	)

	Convey("Given a test cases", t, func() {
		cases := []struct {
			program     string
			query       string
			wantResult  []testutil.TermResults
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
				wantResult:  []testutil.TermResults{{}},
				wantSuccess: true,
			},
			{ // All good with hex encoding
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				eddsa_verify(PubKey, '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Sig, encoding(hex)).`,
				query:       `verify.`,
				wantResult:  []testutil.TermResults{{}},
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
				wantError:   fmt.Errorf("error(syntax_error([%s]),hex_bytes/2)", badPublicKeyLength),
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
			{ // Incorrect algo
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				eddsa_verify(PubKey, Msg, Sig, [encoding(octet), type(foo)]).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(cryptographic_algorithm,foo),hex_bytes/2)"),
			},
			{ // Unsupported algo
				program: `verify :-
				hex_bytes('53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				eddsa_verify(PubKey, Msg, Sig, [encoding(octet), type(secp256k1)]).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(cryptographic_algorithm,secp256k1),hex_bytes/2)"),
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
				wantResult:  []testutil.TermResults{{}},
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
				wantError:   fmt.Errorf("error(syntax_error([%s]),hex_bytes/2)", failedToParsePublicKey),
			},
			{ // Unsupported algo
				program: `verify :-
				hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
				ecdsa_verify(PubKey, Msg, Sig, [encoding(octet), type(foo)]).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(cryptographic_algorithm,foo),hex_bytes/2)"),
			},
			{
				// Wrong msg
				program: `verify :-
				hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
				hex_bytes('e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3bf', Msg),
				hex_bytes('30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8ed7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e47', Sig),
				ecdsa_verify(PubKey, Msg, Sig, encoding(octet)).`,
				query:       `verify.`,
				wantResult:  []testutil.TermResults{{}},
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
				wantResult:  []testutil.TermResults{{}},
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
				wantResult:  []testutil.TermResults{{}},
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
				wantResult:  []testutil.TermResults{{}},
				wantSuccess: false,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := dbm.NewMemDB()
					stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewComprehensiveInterpreterMust(ctx)
						interpreter.Register4(engine.NewAtom("ecdsa_verify"), ECDSAVerify)
						interpreter.Register4(engine.NewAtom("eddsa_verify"), EDDSAVerify)

						err := interpreter.Compile(ctx, ":- consult('/v1/lib/crypto.pl').")
						So(err, ShouldBeNil)
						err = interpreter.Compile(ctx, tc.program)
						So(err, ShouldBeNil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldBeNil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []testutil.TermResults
									for sols.Next() {
										m := testutil.TermResults{}
										err := sols.Scan(m)
										So(err, ShouldBeNil)

										got = append(got, m)
									}
									if tc.wantError != nil {
										So(sols.Err(), ShouldNotBeNil)
										So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
									} else {
										So(sols.Err(), ShouldBeNil)

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
