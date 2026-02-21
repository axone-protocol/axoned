package comet

import (
	"context"
	"io/fs"

	"github.com/axone-protocol/prolog/v3/engine"

	corecomet "cosmossdk.io/core/comet"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/virtualfile"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

const (
	atPath              = "@"
	validatorsHashPath  = "validators_hash"
	proposerAddressPath = "proposer_address"
	evidencePath        = "evidence"
	lastCommitPath      = "last_commit"
	lastCommitRoundPath = "last_commit/round"
	lastCommitVotesPath = "last_commit/votes"
)

var (
	atomComet            = engine.NewAtom("comet")
	atomValidatorsHash   = engine.NewAtom("validators_hash")
	atomProposerAddress  = engine.NewAtom("proposer_address")
	atomEvidence         = engine.NewAtom("evidence")
	atomLastCommit       = engine.NewAtom("last_commit")
	atomCommitInfo       = engine.NewAtom("commit_info")
	atomRound            = engine.NewAtom("round")
	atomVotes            = engine.NewAtom("votes")
	atomVoteInfo         = engine.NewAtom("vote_info")
	atomBlockIDFlag      = engine.NewAtom("block_id_flag")
	atomValidator        = engine.NewAtom("validator")
	atomAddress          = engine.NewAtom("address")
	atomPower            = engine.NewAtom("power")
	atomType             = engine.NewAtom("type")
	atomHeight           = engine.NewAtom("height")
	atomTime             = engine.NewAtom("time")
	atomTotalVotingPower = engine.NewAtom("total_voting_power")
)

type vfs struct {
	ctx context.Context
}

type cometTerms struct {
	validatorsHash  engine.Term
	proposerAddress engine.Term
	evidence        engine.Term
	lastCommit      engine.Term
	lastCommitRound engine.Term
	lastCommitVotes engine.Term
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates the /v1/sys/comet snapshot filesystem.
func NewFS(ctx context.Context) fs.ReadFileFS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	data, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return virtualfile.New(name, data, prolog.ResolveHeaderInfo(sdkCtx).Time), nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	return f.readFile("readfile", name)
}

func (f *vfs) readFile(op, name string) ([]byte, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	content, err := renderFile(sdkCtx.CometInfo(), subpath)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	return content, nil
}

func renderFile(info corecomet.BlockInfo, subpath string) ([]byte, error) {
	terms, err := newCometTerms(info)
	if err != nil {
		return nil, err
	}

	switch subpath {
	case atPath:
		dictTerm, err := terms.dict()
		if err != nil {
			return nil, err
		}
		return prologterm.Render(dictTerm, true)
	case validatorsHashPath:
		return prologterm.Render(terms.validatorsHash, true)
	case proposerAddressPath:
		return prologterm.Render(terms.proposerAddress, true)
	case evidencePath:
		return prologterm.Render(terms.evidence, true)
	case lastCommitPath:
		return prologterm.Render(terms.lastCommit, true)
	case lastCommitRoundPath:
		return prologterm.Render(terms.lastCommitRound, true)
	case lastCommitVotesPath:
		return prologterm.Render(terms.lastCommitVotes, true)
	default:
		return nil, fs.ErrNotExist
	}
}

func newCometTerms(info corecomet.BlockInfo) (cometTerms, error) {
	terms := cometTerms{
		validatorsHash:  prolog.BytesToByteListTerm(nil),
		proposerAddress: prolog.BytesToByteListTerm(nil),
		evidence:        engine.List(),
		lastCommitRound: engine.Integer(0),
		lastCommitVotes: engine.List(),
	}

	if info != nil {
		terms.validatorsHash = prolog.BytesToByteListTerm(copyBytes(info.GetValidatorsHash()))
		terms.proposerAddress = prolog.BytesToByteListTerm(copyBytes(info.GetProposerAddress()))

		evidenceTerm, err := evidenceListToTerm(info.GetEvidence())
		if err != nil {
			return cometTerms{}, err
		}
		terms.evidence = evidenceTerm

		lastCommitRound, lastCommitVotes, err := lastCommitTerms(info.GetLastCommit())
		if err != nil {
			return cometTerms{}, err
		}
		terms.lastCommitRound = lastCommitRound
		terms.lastCommitVotes = lastCommitVotes
	}

	lastCommitTerm, err := commitInfoToTerm(terms.lastCommitRound, terms.lastCommitVotes)
	if err != nil {
		return cometTerms{}, err
	}
	terms.lastCommit = lastCommitTerm

	return terms, nil
}

func (t cometTerms) dict() (engine.Term, error) {
	return engine.NewDict([]engine.Term{
		atomComet,
		atomValidatorsHash, t.validatorsHash,
		atomProposerAddress, t.proposerAddress,
		atomEvidence, t.evidence,
		atomLastCommit, t.lastCommit,
	})
}

func commitInfoToTerm(round, votes engine.Term) (engine.Term, error) {
	return engine.NewDict([]engine.Term{
		atomCommitInfo,
		atomRound, round,
		atomVotes, votes,
	})
}

func lastCommitTerms(commitInfo corecomet.CommitInfo) (engine.Term, engine.Term, error) {
	round := engine.Integer(0)
	votes := engine.List()

	if commitInfo == nil {
		return round, votes, nil
	}

	round = engine.Integer(int64(commitInfo.Round()))
	votesTerm, err := voteInfosToTerm(commitInfo.Votes())
	if err != nil {
		return nil, nil, err
	}

	return round, votesTerm, nil
}

func voteInfosToTerm(votes corecomet.VoteInfos) (engine.Term, error) {
	if votes == nil {
		return engine.List(), nil
	}

	terms := make([]engine.Term, 0, votes.Len())
	for i := 0; i < votes.Len(); i++ {
		voteTerm, err := voteInfoToTerm(votes.Get(i))
		if err != nil {
			return nil, err
		}
		terms = append(terms, voteTerm)
	}

	return engine.List(terms...), nil
}

func voteInfoToTerm(voteInfo corecomet.VoteInfo) (engine.Term, error) {
	blockIDFlag := int64(0)
	var validator corecomet.Validator
	if voteInfo != nil {
		blockIDFlag = int64(voteInfo.GetBlockIDFlag())
		validator = voteInfo.Validator()
	}

	validatorTerm, err := validatorToTerm(validator)
	if err != nil {
		return nil, err
	}

	return engine.NewDict([]engine.Term{
		atomVoteInfo,
		atomBlockIDFlag, engine.Integer(blockIDFlag),
		atomValidator, validatorTerm,
	})
}

func evidenceListToTerm(evidenceList corecomet.EvidenceList) (engine.Term, error) {
	if evidenceList == nil {
		return engine.List(), nil
	}

	terms := make([]engine.Term, 0, evidenceList.Len())
	for i := 0; i < evidenceList.Len(); i++ {
		evidenceTerm, err := evidenceToTerm(evidenceList.Get(i))
		if err != nil {
			return nil, err
		}
		terms = append(terms, evidenceTerm)
	}

	return engine.List(terms...), nil
}

func evidenceToTerm(evidence corecomet.Evidence) (engine.Term, error) {
	typeTerm := engine.Integer(0)
	heightTerm := engine.Integer(0)
	timeTerm := engine.Integer(0)
	totalVotingPowerTerm := engine.Integer(0)
	var validator corecomet.Validator
	if evidence != nil {
		typeTerm = engine.Integer(int64(evidence.Type()))
		heightTerm = engine.Integer(evidence.Height())
		totalVotingPowerTerm = engine.Integer(evidence.TotalVotingPower())
		validator = evidence.Validator()

		t := evidence.Time()
		if !t.IsZero() {
			timeTerm = engine.Integer(t.Unix())
		}
	}

	validatorTerm, err := validatorToTerm(validator)
	if err != nil {
		return nil, err
	}

	return engine.NewDict([]engine.Term{
		atomEvidence,
		atomType, typeTerm,
		atomValidator, validatorTerm,
		atomHeight, heightTerm,
		atomTime, timeTerm,
		atomTotalVotingPower, totalVotingPowerTerm,
	})
}

func validatorToTerm(validator corecomet.Validator) (engine.Term, error) {
	address := []byte(nil)
	power := int64(0)
	if validator != nil {
		address = copyBytes(validator.Address())
		power = validator.Power()
	}

	return engine.NewDict([]engine.Term{
		atomValidator,
		atomAddress, prolog.BytesToByteListTerm(address),
		atomPower, engine.Integer(power),
	})
}

func copyBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}

	cpy := make([]byte, len(b))
	copy(cpy, b)

	return cpy
}
