% crypto.pl
% Small crypto-adjacent helpers.

:- consult('/v1/lib/error.pl').

%! crypto_data_hash(+Data, ?Hash, +Options) is det.
%
% Computes the Hash of Data using a configured hashing algorithm.
%
% Options may be a single option term or a list of option terms:
% - `algorithm(+Algorithm)` selects `sha256` (default), `sha512`, or `md5`;
% - `encoding(+Encoding)` selects `utf8` (default), `text`, `hex`, or `octet`.
%
% Data is interpreted according to Encoding and Hash is unified with the
% resulting digest as a list of bytes.
crypto_data_hash(Data, Hash, Options) :-
  crypto_hash_option(algorithm, Options, sha256, Algorithm),
  crypto_hash_algorithm(Algorithm),
  crypto_hash_option(encoding, Options, utf8, Encoding),
  with_context(crypto_data_hash/3, crypto_hash_data_bytes(Encoding, Data, DataBytes)),
  crypto_hash_path(Algorithm, Path),
  crypto_hash_dev_call(Path, DataBytes, HashBytes),
  Hash = HashBytes.

crypto_hash_option(Name, Options, Default, Value) :-
  ( crypto_hash_options_list(Options)
  -> crypto_hash_option_list(Name, Options, Found),
     crypto_hash_option_value(Found, Default, Value)
  ; crypto_hash_valid_option_term(Options)
  -> ( crypto_hash_option_term(Name, Options, Found)
     -> Value = Found
     ;  Value = Default
     )
  ;  throw(error(type_error(option, Options), crypto_data_hash/3))
  ).

crypto_hash_options_list(Options) :-
  nonvar(Options),
  ( Options = []
  ; Options = [_|_]
  ).

crypto_hash_option_list(_, [], missing).
crypto_hash_option_list(Name, [Option | Rest], Found) :-
  crypto_hash_valid_option_term(Option),
  ( crypto_hash_option_term(Name, Option, Value)
  -> Found = found(Value)
  ;  crypto_hash_option_list(Name, Rest, Found)
  ).

crypto_hash_valid_option_term(Option) :-
  ( var(Option)
  -> throw(error(instantiation_error, crypto_data_hash/3))
  ; compound(Option)
  -> true
  ;  throw(error(type_error(option, Option), crypto_data_hash/3))
  ).

crypto_hash_option_term(Name, Option, Value) :-
  compound(Option),
  functor(Option, Name, 1),
  arg(1, Option, Value).

crypto_hash_option_value(found(Value), _, Value).
crypto_hash_option_value(missing, Default, Default).

crypto_hash_algorithm(Algorithm) :-
  ( atom(Algorithm)
  -> crypto_hash_algorithm_atom(Algorithm)
  ;  throw(error(type_error(atom, Algorithm), crypto_data_hash/3))
  ).

crypto_hash_algorithm_atom(sha256) :- !.
crypto_hash_algorithm_atom(sha512) :- !.
crypto_hash_algorithm_atom(md5) :- !.
crypto_hash_algorithm_atom(Algorithm) :-
  throw(error(type_error(hash_algorithm, Algorithm), crypto_data_hash/3)).

crypto_hash_data_bytes(hex, Data, Bytes) :-
  !,
  hex_bytes(Data, Bytes).
crypto_hash_data_bytes(octet, Data, Bytes) :-
  !,
  must_be(list(byte), Data),
  Bytes = Data.
crypto_hash_data_bytes(utf8, Data, Bytes) :-
  !,
  string_bytes(Data, Bytes, utf8).
crypto_hash_data_bytes(text, Data, Bytes) :-
  !,
  string_bytes(Data, Bytes, text).
crypto_hash_data_bytes(Encoding, _, _) :-
  ( atom(Encoding)
  -> throw(error(domain_error(valid_encoding(Encoding), Encoding), crypto_data_hash/3))
  ;  throw(error(type_error(atom, Encoding), crypto_data_hash/3))
  ).

crypto_hash_path(Algorithm, Path) :-
  atom_concat('/v1/dev/crypto/', Algorithm, Path).

crypto_hash_dev_call(Path, DataBytes, HashBytes) :-
  ( current_predicate(dev_call/4)
  -> true
  ;  consult('/v1/lib/dev.pl')
  ),
  dev_call(Path, binary, dev_write_bytes(DataBytes), dev_read_bytes(HashBytes)).

%! eddsa_verify(+PubKey, +Data, +Signature, +Options) is semidet.
%
% Succeeds when Signature is a valid EdDSA signature for Data and PubKey.
%
% PubKey and Signature are lists of bytes.
%
% Options may be a single option term or a list of option terms:
% - `type(+Algorithm)` selects `ed25519` (default);
% - `encoding(+Encoding)` selects how Data is interpreted, defaulting to `hex`.
%
% Supported encodings are `hex`, `octet`, `utf8`, and `text`.
eddsa_verify(PubKey, Data, Signature, Options) :-
  crypto_verify(
    eddsa_verify/4,
    [ed25519],
    ed25519,
    PubKey,
    Data,
    Signature,
    Options
  ).

%! ecdsa_verify(+PubKey, +Data, +Signature, +Options) is semidet.
%
% Succeeds when Signature is a valid ECDSA signature for Data and PubKey.
%
% PubKey is the compressed public key as a list of bytes. Signature is the ASN.1
% encoded signature as a list of bytes.
%
% Options may be a single option term or a list of option terms:
% - `type(+Algorithm)` selects `secp256r1` (default) or `secp256k1`;
% - `encoding(+Encoding)` selects how Data is interpreted, defaulting to `hex`.
%
% Supported encodings are `hex`, `octet`, `utf8`, and `text`.
ecdsa_verify(PubKey, Data, Signature, Options) :-
  crypto_verify(
    ecdsa_verify/4,
    [secp256r1, secp256k1],
    secp256r1,
    PubKey,
    Data,
    Signature,
    Options
  ).

crypto_verify(Context, Algorithms, DefaultAlgorithm, PubKey, Data, Signature, Options) :-
  crypto_verify_option(Context, type, Options, DefaultAlgorithm, Algorithm),
  crypto_verify_algorithm(Context, Algorithms, Algorithm),
  crypto_verify_option(Context, encoding, Options, hex, Encoding),
  with_context(Context, must_be(list(byte), PubKey)),
  with_context(Context, must_be(list(byte), Signature)),
  crypto_verify_data_bytes(Context, Encoding, Data, DataBytes),
  crypto_verify_dev_call(Algorithm, PubKey, DataBytes, Signature, Response),
  crypto_verify_response(Context, Response).

crypto_verify_option(Context, Name, Options, Default, Value) :-
  ( crypto_verify_options_list(Options)
  -> crypto_verify_option_list(Context, Name, Options, Found),
     crypto_verify_option_value(Found, Default, Value)
  ; crypto_verify_valid_option_term(Context, Options)
  -> ( crypto_verify_option_term(Name, Options, Found)
     -> Value = Found
     ;  Value = Default
     )
  ;  throw(error(type_error(option, Options), Context))
  ).

crypto_verify_options_list(Options) :-
  nonvar(Options),
  ( Options = []
  ; Options = [_|_]
  ).

crypto_verify_option_list(_, _, [], missing).
crypto_verify_option_list(Context, Name, [Option | Rest], Found) :-
  crypto_verify_valid_option_term(Context, Option),
  ( crypto_verify_option_term(Name, Option, Value)
  -> Found = found(Value)
  ;  crypto_verify_option_list(Context, Name, Rest, Found)
  ).

crypto_verify_valid_option_term(Context, Option) :-
  ( var(Option)
  -> throw(error(instantiation_error, Context))
  ; ( Option = type(_)
    ; Option = encoding(_)
    )
  -> true
  ;  throw(error(type_error(option, Option), Context))
  ).

crypto_verify_option_term(Name, Option, Value) :-
  compound(Option),
  functor(Option, Name, 1),
  arg(1, Option, Value).

crypto_verify_option_value(found(Value), _, Value).
crypto_verify_option_value(missing, Default, Default).

crypto_verify_algorithm(Context, Algorithms, Algorithm) :-
  ( atom(Algorithm)
  -> ( crypto_member(Algorithm, Algorithms)
     -> true
     ;  throw(error(type_error(cryptographic_algorithm, Algorithm), Context))
     )
  ;  throw(error(type_error(atom, Algorithm), Context))
  ).

crypto_member(X, [X | _]) :-
  !.
crypto_member(X, [_ | Rest]) :-
  crypto_member(X, Rest).

crypto_verify_data_bytes(Context, hex, Data, Bytes) :-
  !,
  with_context(Context, hex_bytes(Data, Bytes)).
crypto_verify_data_bytes(Context, octet, Data, Bytes) :-
  !,
  with_context(Context, must_be(list(byte), Data)),
  Bytes = Data.
crypto_verify_data_bytes(Context, utf8, Data, Bytes) :-
  !,
  with_context(Context, string_bytes(Data, Bytes, utf8)).
crypto_verify_data_bytes(Context, text, Data, Bytes) :-
  !,
  with_context(Context, string_bytes(Data, Bytes, text)).
crypto_verify_data_bytes(Context, Encoding, _, _) :-
  ( atom(Encoding)
  -> throw(error(domain_error(valid_encoding(Encoding), Encoding), Context))
  ;  throw(error(type_error(atom, Encoding), Context))
  ).

crypto_verify_dev_call(Algorithm, PubKey, Data, Signature, Response) :-
  ( current_predicate(dev_call/4)
  -> true
  ;  consult('/v1/lib/dev.pl')
  ),
  atom_concat('/v1/dev/crypto/', Algorithm, Path),
  hex_bytes(PubKeyHex, PubKey),
  hex_bytes(DataHex, Data),
  hex_bytes(SignatureHex, Signature),
  dev_call(
    Path,
    text,
    crypto_verify_write_request(PubKeyHex, DataHex, SignatureHex),
    crypto_verify_read_response(Response)
  ).

crypto_verify_write_request(Stream, PubKeyHex, DataHex, SignatureHex) :-
  crypto_put_chars(Stream, ['v', 'e', 'r', 'i', 'f', 'y', ' ']),
  atom_chars(PubKeyHex, PubKeyChars),
  crypto_put_chars(Stream, PubKeyChars),
  put_char(Stream, ' '),
  atom_chars(DataHex, DataChars),
  crypto_put_chars(Stream, DataChars),
  put_char(Stream, ' '),
  atom_chars(SignatureHex, SignatureChars),
  crypto_put_chars(Stream, SignatureChars),
  put_char(Stream, '\n').

crypto_put_chars(_, []).
crypto_put_chars(Stream, [Char | Rest]) :-
  put_char(Stream, Char),
  crypto_put_chars(Stream, Rest).

crypto_verify_read_response(Stream, Response) :-
  catch(
    read_term(Stream, Response0, []),
    _,
    Response0 = error(system_error)
  ),
  ( Response0 == end_of_file
  -> Response = error(system_error)
  ;  Response = Response0
  ).

crypto_verify_response(_, ok(true)) :-
  !.
crypto_verify_response(_, ok(false)) :-
  !,
  fail.
crypto_verify_response(Context, error(invalid_key)) :-
  !,
  throw(error(syntax_error(invalid_key), Context)).
crypto_verify_response(Context, error(invalid_request)) :-
  !,
  throw(error(system_error, Context)).
crypto_verify_response(Context, error(unsupported_operation)) :-
  !,
  throw(error(system_error, Context)).
crypto_verify_response(Context, error(system_error)) :-
  !,
  throw(error(system_error, Context)).
crypto_verify_response(Context, _) :-
  throw(error(system_error, Context)).

%! hex_bytes(?Hex, ?Bytes) is det.
%
% Relates a hexadecimal text representation to a list of bytes.
%
% - Hex may be an atom, a list of characters, or a list of character codes.
% - Bytes is a proper list of integers in [0,255].
% - At least one argument must be instantiated.
% - When converting Bytes to Hex, Hex is returned as a lowercase atom.
hex_bytes(Hex, Bytes) :-
  ( nonvar(Hex)
  -> with_context(hex_bytes/2, must_be(text, Hex)),
     hex_text_chars(Hex, Chars),
     hex_chars_bytes(Chars, Hex, DecodedBytes),
     ( nonvar(Bytes)
     -> with_context(hex_bytes/2, must_be(list(byte), Bytes)),
        Bytes = DecodedBytes
     ;  Bytes = DecodedBytes
     )
  ; nonvar(Bytes)
  -> with_context(hex_bytes/2, must_be(list(byte), Bytes)),
     hex_bytes_chars(Bytes, Chars),
     atom_chars(HexAtom, Chars),
     Hex = HexAtom
  ; throw(error(instantiation_error, hex_bytes/2))
  ).

hex_text_chars(Hex, Chars) :-
  ( atom(Hex)
  -> atom_chars(Hex, Chars)
  ; has_type(chars, Hex)
  -> Chars = Hex
  ; hex_codes_chars(Hex, Chars)
  ).

hex_codes_chars([], []).
hex_codes_chars([Code | Rest], [Char | Chars]) :-
  char_code(Char, Code),
  hex_codes_chars(Rest, Chars).

hex_chars_bytes([], _, []).
hex_chars_bytes([_], Hex, _) :-
  throw(error(domain_error(valid_encoding(hex), Hex), hex_bytes/2)).
hex_chars_bytes([HiChar, LoChar | Rest], Hex, [Byte | Bytes]) :-
  hex_char_nibble(HiChar, Hex, HiNibble),
  hex_char_nibble(LoChar, Hex, LoNibble),
  Byte is HiNibble * 16 + LoNibble,
  hex_chars_bytes(Rest, Hex, Bytes).

hex_char_nibble(Char, _, Nibble) :-
  char_code(Char, Code),
  Code >= 0'0,
  Code =< 0'9,
  !,
  Nibble is Code - 0'0.
hex_char_nibble(Char, _, Nibble) :-
  char_code(Char, Code),
  Code >= 0'a,
  Code =< 0'f,
  !,
  Nibble is Code - 0'a + 10.
hex_char_nibble(Char, _, Nibble) :-
  char_code(Char, Code),
  Code >= 0'A,
  Code =< 0'F,
  !,
  Nibble is Code - 0'A + 10.
hex_char_nibble(_, Hex, _) :-
  throw(error(domain_error(valid_encoding(hex), Hex), hex_bytes/2)).

hex_bytes_chars([], []).
hex_bytes_chars([Byte | Rest], [HiChar, LoChar | Chars]) :-
  hex_low_nibble(Byte, LoNibble),
  hex_high_nibble(Byte, HiNibble),
  hex_nibble_char(HiNibble, HiChar),
  hex_nibble_char(LoNibble, LoChar),
  hex_bytes_chars(Rest, Chars).

hex_high_nibble(Byte, 0) :-
  Byte < 16,
  !.
hex_high_nibble(Byte, HiNibble) :-
  Byte1 is Byte - 16,
  hex_high_nibble(Byte1, PrevNibble),
  HiNibble is PrevNibble + 1.

hex_low_nibble(Byte, LoNibble) :-
  ( Byte < 16
  -> LoNibble = Byte
  ;  Byte1 is Byte - 16,
     hex_low_nibble(Byte1, LoNibble)
  ).

hex_nibble_char(0, '0').
hex_nibble_char(1, '1').
hex_nibble_char(2, '2').
hex_nibble_char(3, '3').
hex_nibble_char(4, '4').
hex_nibble_char(5, '5').
hex_nibble_char(6, '6').
hex_nibble_char(7, '7').
hex_nibble_char(8, '8').
hex_nibble_char(9, '9').
hex_nibble_char(10, 'a').
hex_nibble_char(11, 'b').
hex_nibble_char(12, 'c').
hex_nibble_char(13, 'd').
hex_nibble_char(14, 'e').
hex_nibble_char(15, 'f').
