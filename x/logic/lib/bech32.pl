% bech32.pl
% Bech32 helpers backed by the codec VFS device.

:- consult('/v1/lib/error.pl').
:- consult('/v1/lib/dev.pl').

%! bech32_address(?Address, ?Bech32) is det.
%
% Converts between a Bech32 atom and its Address pair representation.
%
% The predicate follows a functional direction:
% - when Address is ground, it encodes Address into Bech32;
% - otherwise, when Bech32 is ground, it decodes Bech32 into Address;
% - otherwise, it throws instantiation_error.
%
% Address is represented as Hrp-Bytes where:
% - Hrp is an atom
% - Bytes is a proper list of byte integers in [0,255]
bech32_address(Address, Bech32) :-
  ground(Address),
  bech32_encode_response(Address, ok(Encoded)),
  !,
  Encoded = Bech32.
bech32_address(Address, _) :-
  ground(Address),
  bech32_encode_response(Address, domain_error),
  !,
  throw(error(domain_error(valid_encoding(bech32), Address), bech32_address/2)).
bech32_address(Address, _) :-
  ground(Address),
  bech32_encode_response(Address, system_error),
  !,
  throw(error(system_error, bech32_address/2)).
bech32_address(Address, Bech32) :-
  ground(Bech32),
  bech32_decode_response(Bech32, ok(Decoded)),
  !,
  Address = Decoded.
bech32_address(_, Bech32) :-
  ground(Bech32),
  bech32_decode_response(Bech32, domain_error),
  !,
  throw(error(domain_error(valid_encoding(bech32), Bech32), bech32_address/2)).
bech32_address(_, Bech32) :-
  ground(Bech32),
  bech32_decode_response(Bech32, system_error),
  !,
  throw(error(system_error, bech32_address/2)).
bech32_address(_, _) :-
  throw(error(instantiation_error, bech32_address/2)).

bech32_encode_response(Address, Outcome) :-
  bech32_encode_request(Address, HrpChars, HexChars),
  dev_call(
    '/v1/dev/codec/bech32',
    text,
    bech32_write_encode_request(HrpChars, HexChars),
    bech32_read_response(Response)
  ),
  bech32_encode_outcome(Response, Outcome).

bech32_decode_response(Bech32, Outcome) :-
  bech32_decode_request(Bech32, Bech32Chars),
  dev_call(
    '/v1/dev/codec/bech32',
    text,
    bech32_write_decode_request(Bech32Chars),
    bech32_read_response(Response)
  ),
  bech32_decode_outcome(Response, Outcome).

bech32_encode_request(Address, HrpChars, HexChars) :-
  bech32_must_be(pair, Address),
  Address = Hrp-Bytes,
  bech32_must_be(atom, Hrp),
  bech32_request_token_chars(Hrp, Address, HrpChars),
  bech32_must_be(list(byte), Bytes),
  bech32_bytes_hex_chars(Bytes, HexChars).

bech32_decode_request(Bech32, Bech32Chars) :-
  bech32_must_be(atom, Bech32),
  bech32_request_token_chars(Bech32, Bech32, Bech32Chars).

bech32_must_be(Type, Term) :-
  catch(
    must_be(Type, Term),
    error(Formal, must_be/2),
    throw(error(Formal, bech32_address/2))
  ).

bech32_request_token_chars(Atom, InvalidValue, Chars) :-
  atom_chars(Atom, Chars),
  ( Chars \= [],
    bech32_token_chars(Chars)
  -> true
  ;  throw(error(domain_error(valid_encoding(bech32), InvalidValue), bech32_address/2))
  ).

bech32_token_chars([]).
bech32_token_chars([Char | Rest]) :-
  char_code(Char, Code),
  Code > 32,
  Code =\= 127,
  bech32_token_chars(Rest).

bech32_bytes_hex_chars([], []).
bech32_bytes_hex_chars([Byte | Rest], [HiChar, LoChar | HexChars]) :-
  LoNibble is Byte mod 16,
  bech32_high_nibble(Byte, HiNibble),
  bech32_hex_char(HiNibble, HiChar),
  bech32_hex_char(LoNibble, LoChar),
  bech32_bytes_hex_chars(Rest, HexChars).

bech32_high_nibble(Byte, 0) :-
  Byte < 16,
  !.
bech32_high_nibble(Byte, HiNibble) :-
  Byte1 is Byte - 16,
  bech32_high_nibble(Byte1, PrevNibble),
  HiNibble is PrevNibble + 1.

bech32_hex_char(0, '0').
bech32_hex_char(1, '1').
bech32_hex_char(2, '2').
bech32_hex_char(3, '3').
bech32_hex_char(4, '4').
bech32_hex_char(5, '5').
bech32_hex_char(6, '6').
bech32_hex_char(7, '7').
bech32_hex_char(8, '8').
bech32_hex_char(9, '9').
bech32_hex_char(10, 'a').
bech32_hex_char(11, 'b').
bech32_hex_char(12, 'c').
bech32_hex_char(13, 'd').
bech32_hex_char(14, 'e').
bech32_hex_char(15, 'f').

bech32_write_encode_request(Stream, HrpChars, HexChars) :-
  bech32_put_chars(Stream, ['e', 'n', 'c', 'o', 'd', 'e', ' ']),
  bech32_put_chars(Stream, HrpChars),
  put_char(Stream, ' '),
  bech32_put_chars(Stream, HexChars),
  put_char(Stream, '\n').

bech32_write_decode_request(Stream, Bech32Chars) :-
  bech32_put_chars(Stream, ['d', 'e', 'c', 'o', 'd', 'e', ' ']),
  bech32_put_chars(Stream, Bech32Chars),
  put_char(Stream, '\n').

bech32_put_chars(_, []).
bech32_put_chars(Stream, [Char | Rest]) :-
  put_char(Stream, Char),
  bech32_put_chars(Stream, Rest).

bech32_read_response(Stream, Response) :-
  catch(
    read_term(Stream, Response0, []),
    _,
    throw(error(system_error, bech32_address/2))
  ),
  ( Response0 == end_of_file
  -> throw(error(system_error, bech32_address/2))
  ;  Response = Response0
  ).

bech32_encode_outcome(ok(Bech32), ok(Bech32)) :-
  atom(Bech32),
  !.
bech32_encode_outcome(error(Code), domain_error) :-
  bech32_domain_error_code(Code),
  !.
bech32_encode_outcome(_, system_error).

bech32_decode_outcome(ok(Address), ok(Address)) :-
  Address = _-_,
  !.
bech32_decode_outcome(error(Code), domain_error) :-
  bech32_domain_error_code(Code),
  !.
bech32_decode_outcome(_, system_error).

bech32_domain_error_code(invalid_bech32).
bech32_domain_error_code(invalid_hrp).
bech32_domain_error_code(invalid_bytes).
