% crypto.pl
% Small crypto-adjacent helpers.

:- consult('/v1/lib/error.pl').

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
