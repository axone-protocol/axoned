% chain.pl
% Chain-related predicates for querying chain information such as block headers and Comet block info.

% comet_info(?CometInfo) is det.
%
% Unifies CometInfo with the current Comet block info dict exposed by the VFS.
%
% Returned term shape:
% ```prolog
% comet{
%   validators_hash: [Byte],
%   proposer_address: [Byte],
%   evidence: [evidence{
%     type: Type,
%     validator: validator{address:[Byte], power:Power},
%     height: Height,
%     time: Time,
%     total_voting_power: TotalVotingPower
%   }],
%   last_commit: commit_info{
%     round: Round,
%     votes: [vote_info{
%       block_id_flag: BlockIDFlag,
%       validator: validator{address:[Byte], power:Power}
%     }]
%   }
% }.
% ```
%
% where:
% - Byte is an integer in [0,255].
% - Time is a Unix timestamp in seconds (0 when unset).
% - Empty lists are returned when data is unavailable.
comet_info(CometInfo) :-
  setup_call_cleanup(
    open('/v1/sys/comet/@', read, Stream, [type(text)]),
    read_term(Stream, CometInfo, []),
    close(Stream)
  ).

% header_info(?HeaderInfo) is det.
%
% Unifies HeaderInfo with the current SDK header info dict exposed by the VFS.
%
% Returned term shape:
% ```prolog
% header{
%   height: Height,
%   hash: [Byte],
%   time: Time,
%   chain_id: ChainID,
%   app_hash: [Byte]
% }.
% ```
%
% where:
% - Height is the current block height.
% - Time is a Unix timestamp in seconds.
% - ChainID is an atom (quoted if needed).
% - Byte is an integer in [0,255].
header_info(HeaderInfo) :-
  setup_call_cleanup(
    open('/v1/sys/header/@', read, Stream, [type(text)]),
    read_term(Stream, HeaderInfo, []),
    close(Stream)
  ).
