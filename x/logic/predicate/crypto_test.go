//nolint:gocognit,lll
package predicate

import (
	"fmt"
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

func TestCryptoHash(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program     string
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				query: `sha_hash('foo', Hash).`,
				wantResult: []types.TermResults{{
					"Hash": "[44,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]",
				}},
				wantSuccess: true,
			},
			{
				query:       `sha_hash(Foo, Hash).`,
				wantResult:  []types.TermResults{},
				wantError:   fmt.Errorf("sha_hash/2: invalid data type: engine.Variable, should be Atom"),
				wantSuccess: false,
			},
			{
				query: `sha_hash('bar',
[44,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]).`,
				wantSuccess: false,
			},
			{
				query: `sha_hash('bar',
[252,222,43,46,219,165,107,244,8,96,31,183,33,254,155,92,51,141,16,238,66,158,160,79,174,85,17,182,143,191,143,185]).`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				query: `sha_hash('bar',
[345,222,43,46,219,165,107,244,8,96,31,183,33,254,155,92,51,141,16,238,66,158,160,79,174,85,17,182,143,191,143,185]).`,
				wantSuccess: false,
			},
			{
				program: `test :- sha_hash('bar', H),
H == [252,222,43,46,219,165,107,244,8,96,31,183,33,254,155,92,51,141,16,238,66,158,160,79,174,85,17,182,143,191,143,185].`,
				query:       `test.`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				program: `test :- sha_hash('bar', H),
H == [2252,222,43,46,219,165,107,244,8,96,31,183,33,254,155,92,51,141,16,238,66,158,160,79,174,85,17,182,143,191,143,185].`,
				query:       `test.`,
				wantSuccess: false,
			},
			{
				query: `hex_bytes(Hex,
[44,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]).`,
				wantResult:  []types.TermResults{{"Hex": "'2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae'"}},
				wantSuccess: true,
			},
			{
				query: `hex_bytes('2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae', Bytes).`,
				wantResult: []types.TermResults{{
					"Bytes": "[44,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]",
				}},
				wantSuccess: true,
			},
			{
				query: `hex_bytes('2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae',
[44,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]).`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				query: `hex_bytes('3c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae',
[44,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]).`,
				wantSuccess: false,
			},
			{
				query: `hex_bytes('fail',
[44,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]).`,
				wantError:   fmt.Errorf("hex_bytes/2: failed decode hexadecimal encoding/hex: invalid byte: U+0069 'i'"),
				wantSuccess: false,
			},
			{
				query: `hex_bytes('2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae',
[45,38,180,107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]).`,
				wantSuccess: false,
			},
			{
				query: `hex_bytes('2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae',
[345,38,'hey',107,104,255,198,143,249,155,69,60,29,48,65,52,19,66,45,112,100,131,191,160,249,138,94,136,98,102,231,174]).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("hex_bytes/2: failed convert list into bytes: invalid integer value in list at position 1: 345 is out of byte range (0-255)"),
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
						interpreter.Register2(engine.NewAtom("sha_hash"), SHAHash)
						interpreter.Register2(engine.NewAtom("hex_bytes"), HexBytes)

						err := interpreter.Compile(ctx, tc.program)
						So(err, ShouldBeNil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldBeNil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []types.TermResults
									for sols.Next() {
										m := types.TermResults{}
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
				wantError:   fmt.Errorf("eddsa_verify/4: failed to verify signature: ed25519: bad public key length: 33"),
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
				wantError:   fmt.Errorf("eddsa_verify/4: invalid type: foo. Possible values: ed25519"),
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
				wantError:   fmt.Errorf("ecdsa_verify/4: failed to verify signature: failed to parse compressed public key"),
			},
			{ // Unsupported algo
				program: `verify :-
			hex_bytes('0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b', PubKey),
			hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
			hex_bytes('889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee00c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b', Sig),
			ecdsa_verify(PubKey, Msg, Sig, [encoding(octet), type(foo)]).`,
				query:       `verify.`,
				wantSuccess: false,
				wantError:   fmt.Errorf("ecdsa_verify/4: invalid type: foo. Possible values: secp256r1, secp256k1"),
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
			/*
				{
					// All good
					program: `verify :-
				hex_bytes('2d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0d0a4d466b77457759484b6f5a497a6a3043415159494b6f5a497a6a3044415163445167414545386843612b527835565547393835506666565870433478446643660d0a6b75747a4c4b4d49586e6c383735734543525036654b4b79704c70514564564752526b3551396f687a64766b49392b583850756d666766356d673d3d0d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d', PubKey),
				hex_bytes('9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', Msg),
				hex_bytes('79ca2e9d19418ad6ed1dd1bc2fd0859458da001b7dfbe1ab4a578b17aeb4bd77ab14228070a57f11a98230c3a49890148c9f0f1c0e88de6d83227e99ab62b0dc', Sig),
				ecdsa_verify(PubKey, Msg, Sig, [encoding(octet), type(secp256k1)]).`,
					query:       `verify.`,
					wantResult:  []types.TermResults{{}},
					wantSuccess: true,
				},*/
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
						interpreter.Register4(engine.NewAtom("eddsa_verify"), EdDSAVerify)
						interpreter.Register4(engine.NewAtom("ecdsa_verify"), ECDSAVerify)

						err := interpreter.Compile(ctx, tc.program)
						So(err, ShouldBeNil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldBeNil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []types.TermResults
									for sols.Next() {
										m := types.TermResults{}
										err := sols.Scan(m)
										So(err, ShouldBeNil)

										got = append(got, m)
									}
									if tc.wantError != nil {
										So(sols.Err(), ShouldBeError, tc.wantError.Error())
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
