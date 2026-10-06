% staking.pl
% Staking-related predicates for delegations, validator bond status, and eligible balances.

:- consult('/v1/lib/error.pl').
:- consult('/v1/lib/bech32.pl').
:- consult('/v1/lib/lists.pl').

%! staking_delegations(+Delegator, -Delegations) is det.
%
% Unifies Delegations with all staking delegations of Delegator, regardless of validator status.
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

%! staking_validators(+Status, -Validators) is det.
%
% Unifies Validators with validators having bond Status: bonded, unbonding, or unbonded.
% Status must be an atom. Each item is validator{operator: Operator, status: Status}.
% Operator is a validator operator address atom. Bond status is independent of jailed status:
% bonded means membership in the active staking validator set, not an unjailed guarantee.
staking_validators(Status, Validators) :-
  with_context(staking_validators/2, must_be(atom, Status)),
  with_context(staking_validators/2, must_be(oneof([bonded, unbonding, unbonded]), Status)),
  atom_concat('/v1/var/lib/staking/validators/', Status, Prefix),
  atom_concat(Prefix, '/@', Path),
  staking_read_collection(Path, Validators).

%! staking_bonded_validators(-Validators) is det.
%
% Unifies Validators with validators in BOND_STATUS_BONDED, including any jailed validators
% whose bond status is still bonded. This is staking_validators(bonded, Validators).
staking_bonded_validators(Validators) :-
  with_context(staking_bonded_validators/1, staking_validators(bonded, Validators)).

%! staking_delegations_to_bonded_validators(+Delegator, -Delegations) is det.
%
% Unifies Delegations with the delegations of Delegator whose validators have bond status
% bonded. It preserves delegation order and balances. It does not additionally filter jailed
% validators. Unbondings, redelegation entries, and distribution rewards are not included.
% Delegator must be an AXONE account address atom in Bech32 format.
staking_delegations_to_bonded_validators(Delegator, Delegations) :-
  with_context(staking_delegations_to_bonded_validators/2, (
    staking_delegations(Delegator, All),
    staking_bonded_validators(Validators),
    staking_filter_bonded(All, Validators, Delegations)
  )).

%! staking_bonded_delegated_balance(+Delegator, -Balance) is det.
%
% Unifies Balance with coin(BondDenom, Amount), the sum of delegations to bonded validators.
% Empty eligible collections return coin(BondDenom, 0), using the staking bond denomination.
% Amount is an integer when it fits in int64, otherwise a decimal atom. Addition is exact,
% including amounts and totals larger than int64. A mismatching coin denomination raises
% domain_error(staking_bond_denom, Denom); different denominations are never added.
% Unbondings, redelegation entries, rewards, and bank balances are excluded.
staking_bonded_delegated_balance(Delegator, coin(Denom, Amount)) :-
  with_context(staking_bonded_delegated_balance/2, (
    staking_delegations_to_bonded_validators(Delegator, Delegations),
    staking_read_collection('/v1/var/lib/staking/params/@', [staking_params{bond_denom: Denom}]),
    staking_sum_balances(Delegations, Denom, [48], Digits),
    staking_reverse(Digits, [], Codes),
    atom_codes(Decimal, Codes),
    length(Codes, Length),
    ( (Length < 19 ; Length =:= 19, Decimal @=< '9223372036854775807')
    -> number_codes(Amount, Codes)
    ; Amount = Decimal
    )
  )).

staking_filter_bonded([], _, []).
staking_filter_bonded([Delegation | Rest], Validators, Eligible) :-
  Delegation = delegation{validator: Operator, balance: _},
  ( member(validator{operator: Operator, status: bonded}, Validators)
  -> Eligible = [Delegation | Tail]
  ; Eligible = Tail
  ),
  staking_filter_bonded(Rest, Validators, Tail).

staking_sum_balances([], _, Sum, Sum).
staking_sum_balances([delegation{validator: _, balance: coin(CoinDenom, Amount)} | Rest], Denom, Sum0, Sum) :-
  ( CoinDenom == Denom
  -> true
  ; throw(error(domain_error(staking_bond_denom, CoinDenom), staking_bonded_delegated_balance/2))
  ),
  ( integer(Amount) -> number_codes(Amount, Codes) ; atom_codes(Amount, Codes) ),
  staking_reverse(Codes, [], Digits),
  staking_add_digits(Sum0, Digits, 0, Sum1),
  staking_sum_balances(Rest, Denom, Sum1, Sum).

% Decimal addition on least-significant-first digit codes avoids int64 overflow.
staking_add_digits([], [], 0, []) :- !.
staking_add_digits([], [], 1, [49]) :- !.
staking_add_digits(Left, Right, Carry, [Code | Rest]) :-
  staking_digit(Left, A, As),
  staking_digit(Right, B, Bs),
  Total is A + B + Carry,
  Code is 48 + Total mod 10,
  NextCarry is Total div 10,
  staking_add_digits(As, Bs, NextCarry, Rest).

staking_digit([], 0, []).
staking_digit([Code | Rest], Digit, Rest) :- Digit is Code - 48.

staking_reverse([], Acc, Acc).
staking_reverse([Head | Tail], Acc, Reversed) :-
  staking_reverse(Tail, [Head | Acc], Reversed).

staking_collection(Context, Delegator, Collection, Positions) :-
  with_context(Context, must_be(atom, Delegator)),
  with_context(Context, bech32_address(_, Delegator)),
  atom_concat('/v1/var/lib/staking/', Delegator, Prefix),
  atom_concat(Prefix, '/', PathPrefix),
  atom_concat(PathPrefix, Collection, CollectionPath),
  atom_concat(CollectionPath, '/@', Path),
  staking_read_collection(Path, Positions).

staking_read_collection(Path, Terms) :-
  setup_call_cleanup(
    open(Path, read, Stream, [type(text)]),
    read_staking_terms(Stream, Terms),
    close(Stream)
  ).

read_staking_terms(Stream, Terms) :-
  read_term(Stream, Term, []),
  (   Term == end_of_file
  ->  Terms = []
  ;   Terms = [Term | Rest],
      read_staking_terms(Stream, Rest)
  ).
