% uri.pl
% URI encoding helpers.

:- consult('/v1/lib/error.pl').

%! uri_encoded(+Component, +Value, -Encoded) is det.
%! uri_encoded(+Component, -Value, +Encoded) is det.
%
% Encoded is the URI encoding for Value.
%
% Component specifies the URI component where the value is used. It is one of
% `query_value`, `fragment`, `path` or `segment`.
%
% Value and Encoded may be atoms, lists of characters, or lists of character
% codes. Generated values are returned as atoms.
uri_encoded(Component, Value, Encoded) :-
  uri_component(Component, URIComponent),
  ( nonvar(Value)
  -> uri_encode(URIComponent, Value, Encoded)
  ; nonvar(Encoded)
  -> uri_decode(Value, Encoded)
  ; throw(error(instantiation_error, uri_encoded/3))
  ).

uri_component(Component, query_value) :-
  nonvar(Component),
  Component == query_value,
  !.
uri_component(Component, fragment) :-
  nonvar(Component),
  Component == fragment,
  !.
uri_component(Component, path) :-
  nonvar(Component),
  Component == path,
  !.
uri_component(Component, segment) :-
  nonvar(Component),
  Component == segment,
  !.
uri_component(Component, _) :-
  var(Component),
  !,
  throw(error(instantiation_error, uri_encoded/3)).
uri_component(Component, _) :-
  throw(error(type_error(uri_component, Component), uri_encoded/3)).

uri_encode(Component, Value, Encoded) :-
  with_context(uri_encoded/3, must_be(text, Value)),
  with_context(uri_encoded/3, string_bytes(Value, ValueBytes, text)),
  uri_escape_bytes(Component, ValueBytes, EncodedBytes),
  atom_codes(EncodedAtom, EncodedBytes),
  Encoded = EncodedAtom.

uri_decode(Value, Encoded) :-
  with_context(uri_encoded/3, must_be(text, Encoded)),
  with_context(uri_encoded/3, string_bytes(Encoded, EncodedBytes, text)),
  uri_unescape_bytes(Encoded, EncodedBytes, DecodedBytes),
  with_context(uri_encoded/3, string_bytes(DecodedChars, DecodedBytes, text)),
  atom_chars(DecodedAtom, DecodedChars),
  Value = DecodedAtom.

uri_escape_bytes(_, [], []).
uri_escape_bytes(Component, [Byte | Rest], Codes) :-
  ( uri_should_escape(Component, Byte)
  -> uri_percent_byte(Byte, [Hi, Lo]),
     Codes = [37, Hi, Lo | Tail]
  ;  Codes = [Byte | Tail]
  ),
  uri_escape_bytes(Component, Rest, Tail).

uri_unescape_bytes(_, [], []).
uri_unescape_bytes(Source, [37 | Rest], [Byte | DecodedRest]) :-
  uri_percent_sequence(Rest, Byte, Next),
  !,
  uri_unescape_bytes(Source, Next, DecodedRest).
uri_unescape_bytes(Source, [37 | Rest], _) :-
  !,
  uri_throw_invalid_escape(Source, [37 | Rest]).
uri_unescape_bytes(Source, [Byte | Rest], [Byte | DecodedRest]) :-
  uri_unescape_bytes(Source, Rest, DecodedRest).

uri_percent_sequence([Hi, Lo | Rest], Byte, Rest) :-
  uri_hex_byte(Hi, HiValue),
  uri_hex_byte(Lo, LoValue),
  Byte is HiValue * 16 + LoValue.

uri_throw_invalid_escape(Source, Bytes) :-
  uri_invalid_escape_fragment(Bytes, Fragment),
  atom_codes('invalid URL escape "', Prefix),
  append(Prefix, Fragment, Partial),
  append(Partial, [34], MessageCodes),
  atom_codes(MessageAtom, MessageCodes),
  atom_chars(MessageAtom, MessageChars),
  throw(error(domain_error(encoding(uri), Source), uri_encoded/3, MessageChars)).

uri_invalid_escape_fragment([37, A, B | _], [37, A, B]) :-
  !.
uri_invalid_escape_fragment([37, A | _], [37, A]) :-
  !.
uri_invalid_escape_fragment([37], [37]).

uri_should_escape(_, Byte) :-
  uri_unreserved_byte(Byte),
  !,
  fail.
uri_should_escape(Component, Byte) :-
  uri_reserved_byte(Byte),
  !,
  uri_reserved_should_escape(Component, Byte).
uri_should_escape(_, _).

uri_reserved_should_escape(path, 63).
uri_reserved_should_escape(path, 58).
uri_reserved_should_escape(segment, 47).
uri_reserved_should_escape(segment, 63).
uri_reserved_should_escape(segment, 58).
uri_reserved_should_escape(query_value, 38).
uri_reserved_should_escape(query_value, 43).
uri_reserved_should_escape(query_value, 58).
uri_reserved_should_escape(query_value, 59).
uri_reserved_should_escape(query_value, 61).

uri_unreserved_byte(Byte) :-
  uri_alpha_byte(Byte).
uri_unreserved_byte(Byte) :-
  uri_digit_byte(Byte).
uri_unreserved_byte(45).
uri_unreserved_byte(46).
uri_unreserved_byte(95).
uri_unreserved_byte(126).

uri_reserved_byte(33).
uri_reserved_byte(36).
uri_reserved_byte(38).
uri_reserved_byte(39).
uri_reserved_byte(40).
uri_reserved_byte(41).
uri_reserved_byte(42).
uri_reserved_byte(43).
uri_reserved_byte(44).
uri_reserved_byte(47).
uri_reserved_byte(58).
uri_reserved_byte(59).
uri_reserved_byte(61).
uri_reserved_byte(63).
uri_reserved_byte(64).

uri_alpha_byte(Byte) :-
  Byte >= 65,
  Byte =< 90.
uri_alpha_byte(Byte) :-
  Byte >= 97,
  Byte =< 122.

uri_digit_byte(Byte) :-
  Byte >= 48,
  Byte =< 57.

uri_percent_byte(Byte, [Hi, Lo]) :-
  HiValue is Byte div 16,
  LoValue is Byte mod 16,
  uri_hex_value_code(HiValue, Hi),
  uri_hex_value_code(LoValue, Lo).

uri_hex_value_code(Value, Code) :-
  Value >= 0,
  Value =< 9,
  Code is 48 + Value.
uri_hex_value_code(Value, Code) :-
  Value >= 10,
  Value =< 15,
  Code is 55 + Value.

uri_hex_byte(Code, Value) :-
  Code >= 48,
  Code =< 57,
  Value is Code - 48.
uri_hex_byte(Code, Value) :-
  Code >= 65,
  Code =< 70,
  Value is Code - 55.
uri_hex_byte(Code, Value) :-
  Code >= 97,
  Code =< 102,
  Value is Code - 87.

uri_path_chars_raw([]).
uri_path_chars_raw(Chars) :-
  uri_path_unit(Chars, Rest),
  uri_path_chars_raw(Rest).

uri_path_unit(['%', Hi, Lo | Rest], Rest) :-
  uri_hex_char(Hi),
  uri_hex_char(Lo).
uri_path_unit(['/' | Rest], Rest).
uri_path_unit([Char | Rest], Rest) :-
  uri_pchar_plain_char(Char).

uri_query_or_fragment_chars_raw([]).
uri_query_or_fragment_chars_raw(Chars) :-
  uri_query_or_fragment_unit(Chars, Rest),
  uri_query_or_fragment_chars_raw(Rest).

uri_query_or_fragment_unit(['%', Hi, Lo | Rest], Rest) :-
  uri_hex_char(Hi),
  uri_hex_char(Lo).
uri_query_or_fragment_unit(['/' | Rest], Rest).
uri_query_or_fragment_unit([QuestionMark | Rest], Rest) :-
  uri_question_mark_char(QuestionMark).
uri_query_or_fragment_unit([Char | Rest], Rest) :-
  uri_pchar_plain_char(Char).

uri_pchar_plain_char(Char) :-
  uri_unreserved_char(Char).
uri_pchar_plain_char(Char) :-
  uri_sub_delim_char(Char).
uri_pchar_plain_char(':').
uri_pchar_plain_char('@').

uri_unreserved_char(Char) :-
  uri_alpha_char(Char).
uri_unreserved_char(Char) :-
  uri_digit_char(Char).
uri_unreserved_char('-').
uri_unreserved_char('.').
uri_unreserved_char('_').
uri_unreserved_char('~').

uri_sub_delim_char('!').
uri_sub_delim_char('$').
uri_sub_delim_char('&').
uri_sub_delim_char('''').
uri_sub_delim_char('(').
uri_sub_delim_char(')').
uri_sub_delim_char('*').
uri_sub_delim_char('+').
uri_sub_delim_char(',').
uri_sub_delim_char(';').
uri_sub_delim_char('=').

uri_alpha_char(Char) :-
  uri_lower_alpha_char(Char).
uri_alpha_char(Char) :-
  uri_upper_alpha_char(Char).

uri_lower_alpha_char(Char) :-
  char_code(Char, Code),
  Code >= 0'a,
  Code =< 0'z.

uri_upper_alpha_char(Char) :-
  char_code(Char, Code),
  Code >= 0'A,
  Code =< 0'Z.

uri_digit_char(Char) :-
  char_code(Char, Code),
  Code >= 0'0,
  Code =< 0'9.

uri_hex_char(Char) :-
  uri_digit_char(Char).
uri_hex_char(Char) :-
  char_code(Char, Code),
  Code >= 0'a,
  Code =< 0'f.
uri_hex_char(Char) :-
  char_code(Char, Code),
  Code >= 0'A,
  Code =< 0'F.

uri_question_mark_char(Char) :-
  char_code(Char, 63).
