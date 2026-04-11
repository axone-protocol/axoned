% base64.pl
% Base64 helpers implemented in Prolog.

:- consult('/v1/lib/error.pl').

%! base64_encoded(?Plain, ?Encoded, +Options) is det.
%
% Relates a text value to its Base64-encoded representation as specified by
% [RFC 4648](https://rfc-editor.org/rfc/rfc4648.html).
%
% The predicate follows a functional direction:
% - when `Plain` is instantiated, it encodes `Plain` into `Encoded`;
% - otherwise, when `Encoded` is instantiated, it decodes `Encoded` into `Plain`;
% - otherwise, it throws `instantiation_error`.
%
% `Plain` may be an atom, a list of characters, or a list of character codes.
% `Encoded` may be an atom, a list of characters, or a list of character codes.
%
% Supported options are:
% - `charset(+Charset)` where `Charset` is `classic` (default) or `url`;
% - `padding(+Boolean)` where `Boolean` is `true` (default) or `false`;
% - `as(+Type)` where `Type` is `string` (default) or `atom`;
% - `encoding(+Encoding)` to translate between text and bytes, defaulting to `utf8`.
base64_encoded(Plain, Encoded, Options) :-
  with_context(base64_encoded/3, base64_options(Options, Charset, Padding, As, Encoding)),
  ( nonvar(Plain)
  -> with_context(base64_encoded/3, string_bytes(Plain, PlainBytes, Encoding)),
     phrase(base64_bytes(Padding, PlainBytes, Charset), EncodedCodes),
     base64_encoded_output(As, Encoded, EncodedCodes)
  ; nonvar(Encoded)
  -> with_context(base64_encoded/3, string_bytes(Encoded, EncodedCodes, text)),
     ( phrase(base64_bytes(Padding, PlainBytes, Charset), EncodedCodes)
     -> base64_decoded_output(As, Plain, PlainBytes, Encoding)
     ;  throw(error(domain_error(encoding(base64), Encoded), base64_encoded/3))
     )
  ; throw(error(instantiation_error, base64_encoded/3))
  ),
  !.

%! base64url(?Plain, ?Encoded) is det.
%
% Relates a text value to its URL-safe Base64 representation.
%
% The predicate is equivalent to `base64_encoded/3` with options
% `[as(atom), encoding(utf8), charset(url), padding(false)]`.
base64url(Plain, Encoded) :-
  base64_encoded(Plain, Encoded, [as(atom), encoding(utf8), charset(url), padding(false)]).

%! base64(?Plain, ?Encoded) is det.
%
% Relates a text value to its classic padded Base64 representation.
%
% The predicate is equivalent to `base64_encoded/3` with options
% `[as(atom), encoding(utf8), charset(classic), padding(true)]`.
base64(Plain, Encoded) :-
  base64_encoded(Plain, Encoded, [as(atom), encoding(utf8), charset(classic), padding(true)]).

base64_options(Options, Charset, Padding, As, Encoding) :-
  base64_option(charset, Options, classic, Charset),
  must_be(atom, Charset),
  base64_charset(Charset),
  base64_option(padding, Options, true, Padding),
  base64_padding(Padding),
  base64_option(as, Options, string, As),
  must_be(atom, As),
  base64_output_type(As),
  base64_option(encoding, Options, utf8, Encoding),
  must_be(atom, Encoding).

base64_option(Name, Options, Default, Value) :-
  ( base64_option_term(Name, Options, Value)
  -> true
  ;  Value = Default
  ).

base64_option_term(_, [], _) :-
  !,
  fail.
base64_option_term(Name, [Option | Rest], Value) :-
  !,
  ( base64_option_match(Name, Option, Value)
  -> true
  ; base64_option_known(Option)
  -> base64_option_term(Name, Rest, Value)
  ; throw(error(type_error(option, Option), base64_encoded/3))
  ).
base64_option_term(Name, Option, Value) :-
  ( base64_option_match(Name, Option, Value)
  -> true
  ; base64_option_known(Option)
  -> fail
  ; throw(error(type_error(option, Option), base64_encoded/3))
  ).

base64_option_match(Name, Option, Value) :-
  nonvar(Option),
  compound(Option),
  Option =.. [Name, Value].

base64_option_known(Option) :-
  nonvar(Option),
  compound(Option),
  Option =.. [_Name, _Value].

base64_charset(classic).
base64_charset(url).
base64_charset(Charset) :-
  throw(error(domain_error(charset, Charset), base64_encoded/3)).

base64_padding(true).
base64_padding(false).
base64_padding(Padding) :-
  throw(error(domain_error(padding, Padding), base64_encoded/3)).

base64_output_type(string).
base64_output_type(atom).
base64_output_type(As) :-
  throw(error(domain_error(as, As), base64_encoded/3)).

base64_encoded_output(atom, Encoded, Codes) :-
  atom_codes(Encoded, Codes).
base64_encoded_output(string, Encoded, Codes) :-
  base64_codes_chars(Codes, Encoded).

base64_decoded_output(atom, Plain, Bytes, Encoding) :-
  with_context(base64_encoded/3, string_bytes(Chars, Bytes, Encoding)),
  atom_chars(Plain, Chars).
base64_decoded_output(string, Plain, Bytes, Encoding) :-
  with_context(base64_encoded/3, string_bytes(Plain, Bytes, Encoding)).

base64_codes_chars([], []).
base64_codes_chars([Code | Rest], [Char | Chars]) :-
  char_code(Char, Code),
  base64_codes_chars(Rest, Chars).

base64_bytes(Padding, Input, Charset) -->
  { nonvar(Input) },
  !,
  base64_encode_bytes(Padding, Input, Charset).
base64_bytes(Padding, Output, Charset) -->
  base64_decode_bytes(Padding, Output, Charset).

base64_encode_bytes(Padding, [I0, I1, I2 | Rest], Charset) -->
  !,
  [O0, O1, O2, O3],
  {
    A is (I0 << 16) + (I1 << 8) + I2,
    O00 is (A >> 18) /\ 63,
    O01 is (A >> 12) /\ 63,
    O02 is (A >> 6) /\ 63,
    O03 is A /\ 63,
    base64_char(Charset, O00, O0),
    base64_char(Charset, O01, O1),
    base64_char(Charset, O02, O2),
    base64_char(Charset, O03, O3)
  },
  base64_encode_bytes(Padding, Rest, Charset).
base64_encode_bytes(true, [I0, I1], Charset) -->
  !,
  [O0, O1, O2, 61],
  {
    A is (I0 << 16) + (I1 << 8),
    O00 is (A >> 18) /\ 63,
    O01 is (A >> 12) /\ 63,
    O02 is (A >> 6) /\ 63,
    base64_char(Charset, O00, O0),
    base64_char(Charset, O01, O1),
    base64_char(Charset, O02, O2)
  }.
base64_encode_bytes(true, [I0], Charset) -->
  !,
  [O0, O1, 61, 61],
  {
    A is I0 << 16,
    O00 is (A >> 18) /\ 63,
    O01 is (A >> 12) /\ 63,
    base64_char(Charset, O00, O0),
    base64_char(Charset, O01, O1)
  }.
base64_encode_bytes(false, [I0, I1], Charset) -->
  !,
  [O0, O1, O2],
  {
    A is (I0 << 16) + (I1 << 8),
    O00 is (A >> 18) /\ 63,
    O01 is (A >> 12) /\ 63,
    O02 is (A >> 6) /\ 63,
    base64_char(Charset, O00, O0),
    base64_char(Charset, O01, O1),
    base64_char(Charset, O02, O2)
  }.
base64_encode_bytes(false, [I0], Charset) -->
  !,
  [O0, O1],
  {
    A is I0 << 16,
    O00 is (A >> 18) /\ 63,
    O01 is (A >> 12) /\ 63,
    base64_char(Charset, O00, O0),
    base64_char(Charset, O01, O1)
  }.
base64_encode_bytes(_, [], _) -->
  [].

base64_decode_bytes(true, Bytes, Charset) -->
  [C0, C1, C2, C3],
  !,
  {
    base64_char(Charset, B0, C0),
    base64_char(Charset, B1, C1)
  },
  !,
  {
    ( C3 =:= 61
    -> ( C2 =:= 61
       -> A is (B0 << 18) + (B1 << 12),
          I0 is (A >> 16) /\ 255,
          Bytes = [I0 | Rest]
       ;  base64_char(Charset, B2, C2),
          A is (B0 << 18) + (B1 << 12) + (B2 << 6),
          I0 is (A >> 16) /\ 255,
          I1 is (A >> 8) /\ 255,
          Bytes = [I0, I1 | Rest]
       )
    ;  base64_char(Charset, B2, C2),
       base64_char(Charset, B3, C3),
       A is (B0 << 18) + (B1 << 12) + (B2 << 6) + B3,
       I0 is (A >> 16) /\ 255,
       I1 is (A >> 8) /\ 255,
       I2 is A /\ 255,
       Bytes = [I0, I1, I2 | Rest]
    )
  },
  base64_decode_bytes(true, Rest, Charset).
base64_decode_bytes(false, Bytes, Charset) -->
  [C0, C1, C2, C3],
  !,
  {
    base64_char(Charset, B0, C0),
    base64_char(Charset, B1, C1),
    base64_char(Charset, B2, C2),
    base64_char(Charset, B3, C3),
    A is (B0 << 18) + (B1 << 12) + (B2 << 6) + B3,
    I0 is (A >> 16) /\ 255,
    I1 is (A >> 8) /\ 255,
    I2 is A /\ 255,
    Bytes = [I0, I1, I2 | Rest]
  },
  base64_decode_bytes(false, Rest, Charset).
base64_decode_bytes(false, [I0, I1], Charset) -->
  [C0, C1, C2],
  !,
  {
    base64_char(Charset, B0, C0),
    base64_char(Charset, B1, C1),
    base64_char(Charset, B2, C2),
    A is (B0 << 18) + (B1 << 12) + (B2 << 6),
    I0 is (A >> 16) /\ 255,
    I1 is (A >> 8) /\ 255
  }.
base64_decode_bytes(false, [I0], Charset) -->
  [C0, C1],
  !,
  {
    base64_char(Charset, B0, C0),
    base64_char(Charset, B1, C1),
    A is (B0 << 18) + (B1 << 12),
    I0 is (A >> 16) /\ 255
  }.
base64_decode_bytes(_, [], _) -->
  [].

base64_char(classic, Value, Code) :-
  base64_classic_char(Value, Code).
base64_char(url, Value, Code) :-
  ( nonvar(Value)
  -> base64_url_value_char(Value, Code)
  ;  base64_url_char_value(Code, Value)
  ).

base64_classic_char(Value, Code) :-
  ( nonvar(Value)
  -> base64_classic_value_char(Value, Code)
  ;  base64_classic_char_value(Code, Value)
  ).

base64_classic_value_char(Value, Code) :-
  Value >= 0,
  Value =< 25,
  !,
  Code is 65 + Value.
base64_classic_value_char(Value, Code) :-
  Value >= 26,
  Value =< 51,
  !,
  Code is 97 + Value - 26.
base64_classic_value_char(Value, Code) :-
  Value >= 52,
  Value =< 61,
  !,
  Code is 48 + Value - 52.
base64_classic_value_char(62, 43).
base64_classic_value_char(63, 47).

base64_classic_char_value(Code, Value) :-
  Code >= 65,
  Code =< 90,
  !,
  Value is Code - 65.
base64_classic_char_value(Code, Value) :-
  Code >= 97,
  Code =< 122,
  !,
  Value is Code - 97 + 26.
base64_classic_char_value(Code, Value) :-
  Code >= 48,
  Code =< 57,
  !,
  Value is Code - 48 + 52.
base64_classic_char_value(43, 62).
base64_classic_char_value(47, 63).

base64_url_value_char(62, 45) :-
  !.
base64_url_value_char(63, 95) :-
  !.
base64_url_value_char(Value, Code) :-
  base64_classic_value_char(Value, Code).

base64_url_char_value(45, 62) :-
  !.
base64_url_char_value(95, 63) :-
  !.
base64_url_char_value(Code, Value) :-
  base64_classic_char_value(Code, Value).
