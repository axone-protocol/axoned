package util

import (
	"errors"
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"
	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func TestAsLimitExceededError(t *testing.T) {
	Convey("Given a context with a finite gas meter", t, func() {
		gasMeter := storetypes.NewGasMeter(100)
		gasMeter.ConsumeGas(17, "seed")
		ctx := newSDKContextWithGasMeter(gasMeter)

		testCases := []struct {
			name          string
			err           error
			expectedError string
			expectSameErr bool
		}{
			{
				name: "nil error",
				err:  nil,
			},
			{
				name:          "already normalized error",
				err:           errorsmod.Wrapf(types.ErrLimitExceeded, "already normalized"),
				expectedError: "already normalized: limit exceeded",
				expectSameErr: true,
			},
			{
				name: "non prolog error",
				err:  errors.New("boom"),
			},
			{
				name: "meter resource error",
				err: engine.NewException(
					engine.Atom("error").Apply(
						engine.Atom("resource_error").Apply(engine.Atom("instruction")),
						engine.Atom("root"),
					),
					nil,
				),
				expectedError: "out of gas: logic <Instruction> (17/100): limit exceeded",
			},
			{
				name: "panic error",
				err: engine.NewException(
					engine.Atom("error").Apply(
						engine.AtomPanicError.Apply(engine.Atom("maximum number of variables reached")),
					),
					nil,
				),
				expectedError: "maximum number of variables reached: limit exceeded",
			},
			{
				name: "unknown resource error with context",
				err: engine.NewException(
					engine.Atom("error").Apply(
						engine.Atom("resource_error").Apply(engine.Atom("foo")),
						engine.Atom("root"),
					),
					nil,
				),
			},
			{
				name: "unrecognized prolog exception",
				err: engine.NewException(
					engine.Atom("error").Apply(engine.Atom("unexpected")),
					nil,
				),
			},
		}

		for _, tc := range testCases {
			Convey(fmt.Sprintf("When normalizing %s", tc.name), func() {
				got := AsLimitExceededError(ctx, tc.err)

				if tc.err == nil || (tc.expectedError == "" && !tc.expectSameErr) {
					Convey("Then no limit exceeded error should be returned", func() {
						So(got, ShouldBeNil)
					})

					return
				}

				Convey("Then a consistent result should be returned", func() {
					So(got, ShouldNotBeNil)
					So(errors.Is(got, types.ErrLimitExceeded), ShouldBeTrue)
					So(got.Error(), ShouldEqual, tc.expectedError)
				})

				if tc.expectSameErr {
					Convey("And the original error should be preserved", func() {
						So(got, ShouldEqual, tc.err)
					})
				}
			})
		}
	})
}

func TestMeterDescriptor(t *testing.T) {
	Convey("Given VM meter resource descriptors", t, func() {
		testCases := []struct {
			resource string
			expected string
			ok       bool
		}{
			{resource: "instruction", expected: "Instruction", ok: true},
			{resource: "arith_node", expected: "ArithNode", ok: true},
			{resource: "compare_step", expected: "CompareStep", ok: true},
			{resource: "copy_node", expected: "CopyNode", ok: true},
			{resource: "list_cell", expected: "ListCell", ok: true},
			{resource: "unify_step", expected: "UnifyStep", ok: true},
			{resource: "custom", expected: "", ok: false},
		}

		for _, tc := range testCases {
			Convey(fmt.Sprintf("When converting %s", tc.resource), func() {
				Convey("Then the descriptor should match expectations", func() {
					actual, ok := meterDescriptor(tc.resource)

					So(ok, ShouldEqual, tc.ok)
					So(actual, ShouldEqual, tc.expected)
				})
			})
		}
	})
}

func newSDKContextWithGasMeter(gasMeter storetypes.GasMeter) sdk.Context {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	return sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger()).WithGasMeter(gasMeter)
}
