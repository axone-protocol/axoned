package app

import (
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"

	corestoretypes "cosmossdk.io/core/store"
	circuitante "cosmossdk.io/x/circuit/ante"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	"github.com/cosmos/ibc-go/v8/modules/core/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions

	IBCKeeper             *keeper.Keeper
	WasmConfig            *wasmTypes.WasmConfig
	WasmKeeper            *wasmkeeper.Keeper
	TXCounterStoreService corestoretypes.KVStoreService
	CircuitKeeper         *circuitkeeper.Keeper
}

// IBCDisabledDecorator rejects IBC transactions when IBC is disabled.
type IBCDisabledDecorator struct{}

// NewIBCDisabledDecorator creates a new IBCDisabledDecorator.
func NewIBCDisabledDecorator() IBCDisabledDecorator {
	return IBCDisabledDecorator{}
}

// AnteHandle rejects IBC transactions when IBCEnabled is false.
func (idd IBCDisabledDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx,
	simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	if !IBCEnabled {
		for _, msg := range tx.GetMsgs() {
			if isIBCMessage(msg) {
				return ctx, sdkerrors.ErrUnauthorized.Wrapf("IBC is disabled: %T", msg)
			}
		}
	}
	return next(ctx, tx, simulate)
}

// isIBCMessage checks if a message is an IBC-related message.
func isIBCMessage(msg sdk.Msg) bool {
	switch msg.(type) {
	// IBC Transfer messages
	case *ibctransfertypes.MsgTransfer:
		return true

	// IBC Core messages
	case *ibcchanneltypes.MsgChannelOpenInit,
		*ibcchanneltypes.MsgChannelOpenTry,
		*ibcchanneltypes.MsgChannelOpenAck,
		*ibcchanneltypes.MsgChannelOpenConfirm,
		*ibcchanneltypes.MsgChannelCloseInit,
		*ibcchanneltypes.MsgChannelCloseConfirm,
		*ibcchanneltypes.MsgRecvPacket,
		*ibcchanneltypes.MsgTimeout,
		*ibcchanneltypes.MsgTimeoutOnClose,
		*ibcchanneltypes.MsgAcknowledgement:
		return true

	// IBC Client messages
	case *ibcclienttypes.MsgCreateClient,
		*ibcclienttypes.MsgUpdateClient,
		*ibcclienttypes.MsgUpgradeClient,
		*ibcclienttypes.MsgSubmitMisbehaviour, //nolint:staticcheck
		*ibcclienttypes.MsgRecoverClient,
		*ibcclienttypes.MsgIBCSoftwareUpgrade:
		return true

	// IBC Connection messages
	case *ibcconnectiontypes.MsgConnectionOpenInit,
		*ibcconnectiontypes.MsgConnectionOpenTry,
		*ibcconnectiontypes.MsgConnectionOpenAck,
		*ibcconnectiontypes.MsgConnectionOpenConfirm:
		return true

	// IBC Fee messages
	case *ibcfeetypes.MsgRegisterPayee,
		*ibcfeetypes.MsgRegisterCounterpartyPayee,
		*ibcfeetypes.MsgPayPacketFee,
		*ibcfeetypes.MsgPayPacketFeeAsync:
		return true

	// ICA messages
	case *icacontrollertypes.MsgSendTx:
		return true

	default:
		return false
	}
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errors.New("account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return nil, errors.New("bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return nil, errors.New("sign mode handler is required for ante builder")
	}
	if options.WasmConfig == nil {
		return nil, errors.New("wasm config is required for ante builder")
	}
	if options.TXCounterStoreService == nil {
		return nil, errors.New("wasm store service is required for ante builder")
	}

	sigGasConsumer := options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewIBCDisabledDecorator(),       // reject IBC transactions when IBCEnabled is false
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
