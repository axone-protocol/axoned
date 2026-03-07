% bank.pl
% Bank-related predicates for querying account balances.

:- consult('/v1/lib/error.pl').
:- consult('/v1/lib/bech32.pl').
:- consult('/v1/lib/lists.pl').

%! bank_balances(+Address, -Balances) is det.
%
% Unifies Balances with the list of coin balances for the given account Address.
% The address must be instantiated (non-variable) and in Bech32 format.
%
% Returned term shape:
% ```prolog
% [Denom-Amount, ...]
% ```
%
% where:
% - Denom is an atom representing the coin denomination.
% - Amount is an integer when it fits in int64, otherwise an atom preserving full precision.
% - The list is sorted by denomination.
%
% Throws instantiation_error if Address is a variable.
% Throws domain_error(valid_encoding(bech32), Address) if Address is not a valid Bech32 address.
%
% Examples:
% ```prolog
% ?- bank_balances('axone1...', Balances).
% Balances = [uatom-100, uaxone-200].
% ```
bank_balances(Address, Balances) :-
  bank_must_be(nonvar, Address, bank_balances/2),
  validate_bech32_address(Address, bank_balances/2),
  atom_concat('/v1/state/bank/', Address, Path1),
  atom_concat(Path1, '/balances/@', Path),
  setup_call_cleanup(
    open(Path, read, Stream, [type(text)]),
    read_terms_from_stream(Stream, Balances),
    close(Stream)
  ).

%! bank_spendable_balances(+Address, -Balances) is det.
%
% Unifies Balances with the list of spendable coin balances for the given account Address.
% The address must be instantiated (non-variable) and in Bech32 format.
%
% Returned term shape:
% ```prolog
% [Denom-Amount, ...]
% ```
%
% where:
% - Denom is an atom representing the coin denomination.
% - Amount is an integer when it fits in int64, otherwise an atom preserving full precision.
% - The list is sorted by denomination.
%
% Throws instantiation_error if Address is a variable.
% Throws domain_error(valid_encoding(bech32), Address) if Address is not a valid Bech32 address.
%
% Examples:
% ```prolog
% ?- bank_spendable_balances('axone1...', Balances).
% Balances = [uatom-100, uaxone-200].
% ```
bank_spendable_balances(Address, Balances) :-
  bank_must_be(nonvar, Address, bank_spendable_balances/2),
  validate_bech32_address(Address, bank_spendable_balances/2),
  atom_concat('/v1/state/bank/', Address, Path1),
  atom_concat(Path1, '/spendable/@', Path),
  setup_call_cleanup(
    open(Path, read, Stream, [type(text)]),
    read_terms_from_stream(Stream, Balances),
    close(Stream)
  ).

%! bank_locked_balances(+Address, -Balances) is det.
%
% Unifies Balances with the list of locked coin balances for the given account Address.
% The address must be instantiated (non-variable) and in Bech32 format.
%
% Returned term shape:
% ```prolog
% [Denom-Amount, ...]
% ```
%
% where:
% - Denom is an atom representing the coin denomination.
% - Amount is an integer when it fits in int64, otherwise an atom preserving full precision.
% - The list is sorted by denomination.
%
% Throws instantiation_error if Address is a variable.
% Throws domain_error(valid_encoding(bech32), Address) if Address is not a valid Bech32 address.
%
% Examples:
% ```prolog
% ?- bank_locked_balances('axone1...', Balances).
% Balances = [uatom-100, uaxone-200].
% ```
bank_locked_balances(Address, Balances) :-
  bank_must_be(nonvar, Address, bank_locked_balances/2),
  validate_bech32_address(Address, bank_locked_balances/2),
  atom_concat('/v1/state/bank/', Address, Path1),
  atom_concat(Path1, '/locked/@', Path),
  setup_call_cleanup(
    open(Path, read, Stream, [type(text)]),
    read_terms_from_stream(Stream, Balances),
    close(Stream)
  ).

% validate_bech32_address(+Address, +Context) is det.
%
% Verifies that Address is valid Bech32 and maps only this specific validation error
% to the given predicate context while leaving all other errors untouched.
validate_bech32_address(Address, Context) :-
  catch(
    bech32_address(_, Address),
    Error,
    rethrow_bech32_error(Error, Address, Context)
  ).

bank_must_be(Type, Term, Context) :-
  catch(
    must_be(Type, Term),
    error(Formal, must_be/2),
    throw(error(Formal, Context))
  ).

rethrow_bech32_error(error(Formal, bech32_address/2), Address, Context) :-
  !,
  normalize_bech32_formal(Formal, Address, Normalized),
  throw(error(Normalized, Context)).
rethrow_bech32_error(error(Formal, _, _), Address, Context) :-
  normalize_bech32_formal(Formal, Address, Normalized),
  throw(error(Normalized, Context)).

normalize_bech32_formal(domain_error(encoding(bech32), _), Address, domain_error(valid_encoding(bech32), Address)) :-
  !.
normalize_bech32_formal(domain_error(valid_encoding(bech32), _), Address, domain_error(valid_encoding(bech32), Address)) :-
  !.
normalize_bech32_formal(Formal, _Address, Formal).

% read_terms_from_stream(+Stream, -Terms) is det.
%
% Helper predicate to read all terms from a stream into a list.
read_terms_from_stream(Stream, Terms) :-
  read_term(Stream, Term, []),
  (   Term == end_of_file
  ->  Terms = []
  ;   Terms = [Term | Rest],
      read_terms_from_stream(Stream, Rest)
  ).
