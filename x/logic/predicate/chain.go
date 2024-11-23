package predicate

import (
	"context"

	"github.com/axone-protocol/prolog/engine"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/proto/tendermint/version"

	"github.com/axone-protocol/axoned/v10/x/logic/prolog"
)

// BlockHeader is a predicate which unifies the given term with the current block header.
//
// # Signature
//
//	block_header(?Header) is det
//
// where:
//
//   - Header is a Dict representing the current chain header at the time of the query.
func BlockHeader(vm *engine.VM, header engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := prolog.UnwrapSDKContext(ctx, env)
		if err != nil {
			return engine.Error(err)
		}

		headerDict, err := blockHeaderToTerm(sdkContext.BlockHeader())
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, headerDict, header, cont, env)
	})
}

func blockHeaderToTerm(header cmtproto.Header) (engine.Term, error) {
	version, err := consensusToTerm(header.Version)
	if err != nil {
		return nil, err
	}

	lastBlockID, err := blockIDToTerm(header.LastBlockId)
	if err != nil {
		return nil, err
	}

	return engine.NewDict(
		[]engine.Term{
			engine.NewAtom("header"),
			engine.NewAtom("version"), version,
			engine.NewAtom("chain_id"), engine.NewAtom(header.ChainID),
			engine.NewAtom("height"), engine.Integer(header.Height),
			engine.NewAtom("time"), engine.Integer(header.Time.Unix()),
			engine.NewAtom("last_block_id"), lastBlockID,
			engine.NewAtom("last_commit_hash"), prolog.BytesToByteListTerm(header.LastCommitHash),
			engine.NewAtom("data_hash"), prolog.BytesToByteListTerm(header.DataHash),
			engine.NewAtom("validators_hash"), prolog.BytesToByteListTerm(header.ValidatorsHash),
			engine.NewAtom("next_validators_hash"), prolog.BytesToByteListTerm(header.NextValidatorsHash),
			engine.NewAtom("consensus_hash"), prolog.BytesToByteListTerm(header.ConsensusHash),
			engine.NewAtom("app_hash"), prolog.BytesToByteListTerm(header.AppHash),
			engine.NewAtom("last_results_hash"), prolog.BytesToByteListTerm(header.LastResultsHash),
			engine.NewAtom("evidence_hash"), prolog.BytesToByteListTerm(header.EvidenceHash),
			engine.NewAtom("proposer_address"), prolog.BytesToByteListTerm(header.ProposerAddress),
		},
	)
}

func partSetHeaderToTerm(partSetHeader cmtproto.PartSetHeader) (engine.Dict, error) {
	return engine.NewDict(
		[]engine.Term{
			engine.NewAtom("part_set_header"),
			engine.NewAtom("total"), engine.Integer(partSetHeader.Total),
			engine.NewAtom("hash"), prolog.BytesToByteListTerm(partSetHeader.Hash),
		})
}

func consensusToTerm(consensus version.Consensus) (engine.Dict, error) {
	return engine.NewDict([]engine.Term{
		engine.NewAtom("consensus"),
		engine.NewAtom("block"), engine.Integer(consensus.Block), //nolint:gosec // disable G115 as it's unlikely to be a problem
		engine.NewAtom("app"), engine.Integer(consensus.App), //nolint:gosec // disable G115 as it's unlikely to be a problem
	})
}

func blockIDToTerm(blockID cmtproto.BlockID) (engine.Dict, error) {
	partSetHeader, err := partSetHeaderToTerm(blockID.PartSetHeader)
	if err != nil {
		return nil, err
	}

	return engine.NewDict(
		[]engine.Term{
			engine.NewAtom("block_id"),
			engine.NewAtom("hash"), prolog.BytesToByteListTerm(blockID.Hash),
			engine.NewAtom("part_set_header"), partSetHeader,
		})
}

// ChainID is a predicate which unifies the given term with the current chain ID. The signature is:
//
// The signature is as follows:
//
//	chain_id(?ID)
//
// where:
//   - ID represents the current chain ID at the time of the query.
func ChainID(vm *engine.VM, chainID engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := prolog.UnwrapSDKContext(ctx, env)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, chainID, engine.NewAtom(sdkContext.ChainID()), cont, env)
	})
}
