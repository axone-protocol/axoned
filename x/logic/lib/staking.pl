% staking.pl
% Staking-related predicates for querying delegator positions.

:- consult('/v1/lib/error.pl').
:- consult('/v1/lib/bech32.pl').

%! staking_delegations(+Delegator, -Delegations) is det.
%
% Unifies Delegations with the active staking delegations of Delegator.
% Delegator must be an AXONE account address atom in Bech32 format.
%
% Each item has the shape:
% ```prolog
% delegation{validator: Validator, balance: coin(Denom, Amount)}
% ```
%
% Validator and Denom are atoms. Amount is an integer when it fits in int64,
% otherwise an atom preserving its full decimal precision.
staking_delegations(Delegator, Delegations) :-
  staking_collection(staking_delegations/2, Delegator, delegations, Delegations).

%! staking_unbonding_delegations(+Delegator, -Unbondings) is det.
%
% Unifies Unbondings with the unbonding staking delegations of Delegator.
% Delegator must be an AXONE account address atom in Bech32 format.
%
% Each item has the shape:
% ```prolog
% unbonding_delegation{
%   validator: Validator,
%   entries: [unbonding_entry{balance: coin(Denom, Amount), completion_time: Time}, ...]
% }
% ```
%
% Time is a Unix timestamp in seconds. Balance is the amount currently expected
% at completion, after any applicable slashing.
staking_unbonding_delegations(Delegator, Unbondings) :-
  staking_collection(staking_unbonding_delegations/2, Delegator, unbonding_delegations, Unbondings).

%! staking_redelegations(+Delegator, -Redelegations) is det.
%
% Unifies Redelegations with the in-flight redelegations of Delegator.
% Delegator must be an AXONE account address atom in Bech32 format.
%
% Each item has the shape:
% ```prolog
% redelegation{
%   source_validator: SourceValidator,
%   destination_validator: DestinationValidator,
%   entries: [redelegation_entry{balance: coin(Denom, Amount), completion_time: Time}, ...]
% }
% ```
%
% Time is a Unix timestamp in seconds.
staking_redelegations(Delegator, Redelegations) :-
  staking_collection(staking_redelegations/2, Delegator, redelegations, Redelegations).

%! staking_positions(+Delegator, -Positions) is det.
%
% Unifies Positions with all staking-module positions of Delegator.
% Delegator must be an AXONE account address atom in Bech32 format.
%
% Positions has the shape:
% ```prolog
% staking{
%   delegations: [delegation{...}, ...],
%   unbonding_delegations: [unbonding_delegation{...}, ...],
%   redelegations: [redelegation{...}, ...]
% }
% ```
%
% This predicate covers staking-module state only. Delegation rewards belong to
% the distribution module and are not included.
staking_positions(Delegator, Positions) :-
  staking_delegations(Delegator, Delegations),
  staking_unbonding_delegations(Delegator, Unbondings),
  staking_redelegations(Delegator, Redelegations),
  Positions = staking{
    delegations: Delegations,
    unbonding_delegations: Unbondings,
    redelegations: Redelegations
  }.

staking_collection(Context, Delegator, Collection, Positions) :-
  with_context(Context, must_be(atom, Delegator)),
  with_context(Context, bech32_address(_, Delegator)),
  atom_concat('/v1/var/lib/staking/', Delegator, Prefix),
  atom_concat(Prefix, '/', PathPrefix),
  atom_concat(PathPrefix, Collection, CollectionPath),
  atom_concat(CollectionPath, '/@', Path),
  setup_call_cleanup(
    open(Path, read, Stream, [type(text)]),
    read_staking_terms(Stream, Positions),
    close(Stream)
  ).

read_staking_terms(Stream, Terms) :-
  read_term(Stream, Term, []),
  (   Term == end_of_file
  ->  Terms = []
  ;   Terms = [Term | Rest],
      read_staking_terms(Stream, Rest)
  ).
