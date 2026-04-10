package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/types"
	"github.com/axone-protocol/axoned/v15/x/logic/util"
)

func (k Keeper) validateProgram(ctx context.Context, params types.Params, source string) error {
	ctx = k.enhanceContext(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	consumeIOGas(sdkCtx.GasMeter(), uint64(len(source)), params.GasPolicy.IoCoeff)

	i, _, err := k.newInterpreter(ctx, params)
	if err != nil {
		if limitErr := util.AsLimitExceededError(ctx, err); limitErr != nil {
			return limitErr
		}

		return errorsmod.Wrapf(types.ErrInternal, "error creating interpreter: %v", err.Error())
	}

	if err := i.ExecContext(ctx, source); err != nil {
		if limitErr := util.AsLimitExceededError(ctx, err); limitErr != nil {
			return limitErr
		}

		return errorsmod.Wrapf(types.ErrInvalidArgument, "error validating program: %v", err.Error())
	}

	return nil
}

func (k Keeper) GetStoredProgram(ctx sdk.Context, programID []byte) (types.StoredProgram, bool, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.StoredProgramKey(programID))
	if bz == nil {
		return types.StoredProgram{}, false, nil
	}

	var program types.StoredProgram
	if err := k.cdc.Unmarshal(bz, &program); err != nil {
		return types.StoredProgram{}, false, err
	}

	return program, true, nil
}

func (k Keeper) SetStoredProgram(ctx sdk.Context, programID []byte, program types.StoredProgram) error {
	store := ctx.KVStore(k.storeKey)

	bz, err := k.cdc.Marshal(&program)
	if err != nil {
		return err
	}

	store.Set(types.StoredProgramKey(programID), bz)

	return nil
}

func (k Keeper) GetProgramPublication(
	ctx sdk.Context, publisher, programID []byte,
) (types.ProgramPublication, bool, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProgramPublicationKey(publisher, programID))
	if bz == nil {
		return types.ProgramPublication{}, false, nil
	}

	var publication types.ProgramPublication
	if err := k.cdc.Unmarshal(bz, &publication); err != nil {
		return types.ProgramPublication{}, false, err
	}

	return publication, true, nil
}

func (k Keeper) SetProgramPublication(
	ctx sdk.Context, publisher, programID []byte, publication types.ProgramPublication,
) error {
	store := ctx.KVStore(k.storeKey)

	bz, err := k.cdc.Marshal(&publication)
	if err != nil {
		return err
	}

	store.Set(types.ProgramPublicationKey(publisher, programID), bz)

	return nil
}

func (k Keeper) EnsureProgramPublication(
	ctx sdk.Context, publisher, programID []byte, publishedAt int64,
) error {
	_, found, err := k.GetProgramPublication(ctx, publisher, programID)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	return k.SetProgramPublication(ctx, publisher, programID, types.ProgramPublication{
		PublishedAt: publishedAt,
	})
}

func (k Keeper) IterateStoredPrograms(
	ctx sdk.Context, walkFn func(programID []byte, program types.StoredProgram) (stop bool, err error),
) error {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.StoredProgramKeyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		programID, err := types.ParseStoredProgramKey(iter.Key())
		if err != nil {
			return err
		}

		var program types.StoredProgram
		if err := k.cdc.Unmarshal(iter.Value(), &program); err != nil {
			return err
		}

		stop, err := walkFn(programID, program)
		if err != nil || stop {
			return err
		}
	}

	return nil
}

func (k Keeper) IterateProgramPublications(
	ctx sdk.Context,
	walkFn func(publisher, programID []byte, publication types.ProgramPublication) (stop bool, err error),
) error {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.ProgramPublicationKeyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		publisher, programID, err := types.ParseProgramPublicationKey(iter.Key())
		if err != nil {
			return err
		}

		var publication types.ProgramPublication
		if err := k.cdc.Unmarshal(iter.Value(), &publication); err != nil {
			return err
		}

		stop, err := walkFn(publisher, programID, publication)
		if err != nil || stop {
			return err
		}
	}

	return nil
}
