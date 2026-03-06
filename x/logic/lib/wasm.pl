% wasm.pl
% CosmWasm device helpers.

:- consult('/v1/lib/error.pl').
:- consult('/v1/lib/dev.pl').

%! wasm_query(+Address, +RequestBytes, -ResponseBytes) is det.
%
% Executes a CosmWasm smart query against the contract at Address.
%
% - Address must be a valid Bech32 account address.
% - RequestBytes is the exact query payload as bytes (typically UTF-8 JSON).
% - ResponseBytes is unified with the raw response bytes returned by the contract.
%
% Both RequestBytes and ResponseBytes use lists of integers in [0,255].
wasm_query(Address, RequestBytes, ResponseBytes) :-
  must_be(nonvar, Address),
  validate_bech32_address(Address, wasm_query/3),
  must_be(list(byte), RequestBytes),
  wasm_query_path(Address, Path),
  dev_call(Path, binary, wasm_write(RequestBytes), wasm_read(ResponseBytes)).

wasm_query_path(Address, Path) :-
  atom_concat('/v1/dev/wasm/', Address, Prefix),
  atom_concat(Prefix, '/query', Path).

wasm_write(Stream, RequestBytes) :-
  dev_write_bytes(Stream, RequestBytes).

wasm_read(Stream, ResponseBytes) :-
  dev_read_bytes(Stream, ResponseBytes).

validate_bech32_address(Address, Context) :-
  catch(
    bech32_address(_, Address),
    Error,
    normalize_bech32_error(Error, Address, Context)
  ).

normalize_bech32_error(Error, Address, Context) :-
  (   bech32_domain_error(Error)
  ->  throw(error(domain_error(encoding(bech32), Address), Context))
  ;   throw(Error)
  ).

bech32_domain_error(error(domain_error(encoding(bech32), _), _)).
bech32_domain_error(error(domain_error(encoding(bech32), _), _, _)).
