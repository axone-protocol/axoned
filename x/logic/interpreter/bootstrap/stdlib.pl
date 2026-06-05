% stdlib.pl
% Common utility predicates auto-loaded with the interpreter bootstrap.

% setup_call_cleanup(:Setup, :Goal, :Cleanup) is det.
%
% Runs Setup once, then Goal, and always executes Cleanup exactly once for
% this deterministic execution path:
% - on success of Goal;
% - on failure of Goal;
% - on exception raised by Goal (then rethrows).
%
% This implementation is intended for deterministic goals in this runtime.
setup_call_cleanup(Setup, Goal, Cleanup) :-
  call(Setup),
  catch(
    (
      call(Goal)
    ;
      call(Cleanup),
      fail
    ),
    Error,
    (
      call(Cleanup),
      throw(Error)
    )
  ),
  call(Cleanup).

% open(+SourceSink, +Mode, -Stream) is det.
%
% Opens SourceSink in Mode and unifies Stream with the opened stream.
% Equivalent to open(SourceSink, Mode, Stream, []).
open(SourceSink, Mode, Stream) :-
  open(SourceSink, Mode, Stream, []).

% retractall(+Head) is det.
%
% Retracts all clauses whose head unifies with Head.
retractall(Head) :-
  retract((Head :- _)),
  fail.
retractall(_).

% close(+Stream) is det.
%
% Closes Stream using default close options.
close(Stream) :-
  close(Stream, []).

% flush_output is det.
%
% Flushes the current output stream.
flush_output :-
  current_output(S),
  flush_output(S).

% at_end_of_stream is semidet.
%
% Succeeds if the current input stream is at or past end of stream.
at_end_of_stream :-
  current_input(S),
  at_end_of_stream(S).

% at_end_of_stream(+Stream) is semidet.
%
% Succeeds if Stream is at or past end of stream.
at_end_of_stream(Stream) :-
  stream_property(Stream, end_of_stream(E)), !,
  memberchk(E, [at, past]).

% get_char(?Char) is det.
%
% Reads the next character from the current input stream.
get_char(Char) :-
  current_input(S),
  get_char(S, Char).

% get_code(?Code) is det.
%
% Reads the next character code from the current input stream.
get_code(Code) :-
  current_input(S),
  get_code(S, Code).

% get_code(+Stream, ?Code) is det.
%
% Reads the next character code from Stream.
% Returns -1 at end of file.
get_code(Stream, Code) :-
  get_char(Stream, Char),
  (Char = end_of_file -> Code = -1; char_code(Char, Code)).

% peek_char(?Char) is det.
%
% Peeks the next character from the current input stream without consuming it.
peek_char(Char) :-
  current_input(S),
  peek_char(S, Char).

% peek_code(?Code) is det.
%
% Peeks the next character code from the current input stream without consuming it.
peek_code(Code) :-
  current_input(S),
  peek_code(S, Code).

% peek_code(+Stream, ?Code) is det.
%
% Peeks the next character code from Stream without consuming it.
% Returns -1 at end of file.
peek_code(Stream, Code) :-
  peek_char(Stream, Char),
  (Char = end_of_file -> Code = -1; char_code(Char, Code)).

%! put_char(+Char) is det.
%
% Writes Char to the current output stream.
put_char(Char) :-
  current_output(S),
  put_char(S, Char).

%! put_code(+Code) is det.
%
% Writes the character represented by Code to the current output stream.
put_code(Code) :-
  current_output(S),
  put_code(S, Code).

% put_code(+Stream, +Code) is det.
%
% Writes the character represented by Code to Stream.
put_code(S, Code) :-
  char_code(Char, Code),
  put_char(S, Char).

%! nl is det.
%
% Writes a newline to the current output stream.
nl :-
  current_output(S),
  nl(S).

% nl(+Stream) is det.
%
% Writes a newline to Stream.
nl(S) :-
  put_char(S, '\n').

% get_byte(?Byte) is det.
%
% Reads the next byte from the current input stream.
get_byte(Byte) :-
  current_input(S),
  get_byte(S, Byte).

% peek_byte(?Byte) is det.
%
% Peeks the next byte from the current input stream without consuming it.
peek_byte(Byte) :-
  current_input(S),
  peek_byte(S, Byte).

% put_byte(+Byte) is det.
%
% Writes Byte to the current output stream.
put_byte(Byte) :-
  current_output(S),
  put_byte(S, Byte).

% read_term(?Term, +Options) is det.
%
% Reads a term from the current input stream with Options.
read_term(Term, Options) :-
  current_input(S),
  read_term(S, Term, Options).

% read(?Term) is det.
%
% Reads a term from the current input stream.
read(Term) :-
  current_input(S),
  read(S, Term).

% read(+Stream, ?Term) is det.
%
% Reads a term from Stream using default read options.
read(Stream, Term) :-
  read_term(Stream, Term, []).

%! read_string(+Stream, ?Length, -String) is det.
%
% Reads characters from Stream and unifies String with an atom containing the
% text read. Length is unified with the number of UTF-8 bytes read. When Length
% is instantiated to a positive integer, reading stops once at least that many
% bytes have been read.
read_string(Stream, Length, String) :-
  catch(
    read_string_(Stream, Length, String),
    error(Formal, _),
    throw(error(Formal, read_string/3))
  ).

read_string_(Stream, Length, String) :-
  read_string_limit(Length, Limit),
  read_string_chars(Stream, Limit, 0, Chars, ByteLength),
  atom_chars(Atom, Chars),
  Length = ByteLength,
  String = Atom.

read_string_limit(Length, Length) :-
  integer(Length),
  Length > 0,
  !.
read_string_limit(_, unbounded).

read_string_chars(_, Limit, ByteLength, [], ByteLength) :-
  read_string_limit_reached(Limit, ByteLength),
  !.
read_string_chars(Stream, Limit, CurrentLength, Chars, ByteLength) :-
  get_char(Stream, Char),
  ( Char = end_of_file
  -> Chars = [],
     ByteLength = CurrentLength
  ;  string_bytes(Char, Bytes, text),
     length(Bytes, CharLength),
     NextLength is CurrentLength + CharLength,
     Chars = [Char | Rest],
     read_string_chars(Stream, Limit, NextLength, Rest, ByteLength)
  ).

read_string_limit_reached(Limit, ByteLength) :-
  integer(Limit),
  ByteLength >= Limit.

%! string_bytes(?String, ?Bytes, +Encoding) is det.
%
% Relates text and its byte representation according to Encoding.
string_bytes(String, Bytes, Encoding) :-
  ( nonvar(String),
    ( Encoding == text
    ; Encoding == utf8
    )
  -> string_text_bytes(String, Bytes)
  ; nonvar(String),
    Encoding == octet
  -> string_octet_bytes(String, Bytes)
  ; nonvar(Bytes),
    Encoding == octet
  -> bytes_octet_string(Bytes, String)
  ; nonvar(Bytes),
    Encoding == text
  -> bytes_text_string(Bytes, String)
  ; nonvar(String)
  -> text_codec_call(encode, text(Encoding, String), Response),
     text_codec_response(string_bytes/3, Response, Bytes)
  ; nonvar(Bytes)
  -> text_codec_call(decode, bytes(Encoding, Bytes), Response),
     text_codec_response(string_bytes/3, Response, String)
  ;  throw(error(instantiation_error, string_bytes/3))
  ).

string_text_bytes(Text, Bytes) :-
  text_chars(Text, Chars),
  chars_utf8_bytes(Chars, Bytes).

string_octet_bytes(Text, Bytes) :-
  text_chars(Text, Chars),
  chars_octet_bytes(Chars, Bytes).

bytes_octet_string(Bytes, String) :-
  bytes_chars(Bytes, Chars),
  String = Chars.

bytes_text_string(Bytes, String) :-
  utf8_bytes_chars(Bytes, Chars),
  String = Chars.

text_chars(Text, Chars) :-
  atom(Text),
  !,
  atom_chars(Text, Chars).
text_chars([], []) :-
  !.
text_chars([Head | Tail], Chars) :-
  atom(Head),
  !,
  char_list_chars([Head | Tail], Chars).
text_chars([Head | Tail], Chars) :-
  integer(Head),
  !,
  code_list_chars([Head | Tail], Chars).
text_chars(Text, _) :-
  throw(error(type_error(text, Text), string_bytes/3)).

char_list_chars([], []).
char_list_chars([Char | Rest], [Char | Chars]) :-
  atom(Char),
  atom_length(Char, 1),
  !,
  char_list_chars(Rest, Chars).
char_list_chars([Char | _], _) :-
  throw(error(type_error(character, Char), string_bytes/3)).

code_list_chars([], []).
code_list_chars([Code | Rest], [Char | Chars]) :-
  integer(Code),
  char_code(Char, Code),
  !,
  code_list_chars(Rest, Chars).
code_list_chars([Code | _], _) :-
  throw(error(type_error(character_code, Code), string_bytes/3)).

chars_utf8_bytes(Chars, Bytes) :-
  chars_utf8_bytes_(Chars, Bytes, []).

chars_utf8_bytes_([], Tail, Tail).
chars_utf8_bytes_([Char | Rest], Bytes, Tail) :-
  char_code(Char, Code),
  utf8_code_bytes(Code, Bytes, RestBytes),
  chars_utf8_bytes_(Rest, RestBytes, Tail).

utf8_code_bytes(Code, [Code | Tail], Tail) :-
  Code =< 127,
  !.
utf8_code_bytes(Code, [B1, B2 | Tail], Tail) :-
  Code =< 2047,
  !,
  B1 is 192 + (Code >> 6),
  B2 is 128 + (Code /\ 63).
utf8_code_bytes(Code, [B1, B2, B3 | Tail], Tail) :-
  Code =< 65535,
  !,
  B1 is 224 + (Code >> 12),
  B2 is 128 + ((Code >> 6) /\ 63),
  B3 is 128 + (Code /\ 63).
utf8_code_bytes(Code, [B1, B2, B3, B4 | Tail], Tail) :-
  B1 is 240 + (Code >> 18),
  B2 is 128 + ((Code >> 12) /\ 63),
  B3 is 128 + ((Code >> 6) /\ 63),
  B4 is 128 + (Code /\ 63).

chars_octet_bytes([], []).
chars_octet_bytes([Char | Rest], [Byte | Bytes]) :-
  char_code(Char, Byte),
  Byte =< 255,
  !,
  chars_octet_bytes(Rest, Bytes).
chars_octet_bytes([Char | _], _) :-
  char_code(Char, Code),
  throw(error(type_error(byte, Code), string_bytes/3)).

bytes_chars([], []).
bytes_chars([Byte | Rest], [Char | Chars]) :-
  integer(Byte),
  Byte >= 0,
  Byte =< 255,
  char_code(Char, Byte),
  !,
  bytes_chars(Rest, Chars).
bytes_chars([Byte | _], _) :-
  throw(error(type_error(byte, Byte), string_bytes/3)).

utf8_bytes_chars([], []).
utf8_bytes_chars([Byte | Rest], [Char | Chars]) :-
  byte_value(Byte),
  utf8_byte_char(Byte, Rest, Code, Next),
  char_code(Char, Code),
  utf8_bytes_chars(Next, Chars).

utf8_byte_char(Byte, Rest, Byte, Rest) :-
  Byte =< 127,
  !.
utf8_byte_char(Byte, [B2 | Rest], Code, Rest) :-
  Byte >= 194,
  Byte =< 223,
  utf8_continuation_byte(B2),
  !,
  Code is ((Byte /\ 31) << 6) + (B2 /\ 63).
utf8_byte_char(Byte, [B2, B3 | Rest], Code, Rest) :-
  Byte >= 224,
  Byte =< 239,
  utf8_continuation_byte(B2),
  utf8_continuation_byte(B3),
  Code is ((Byte /\ 15) << 12) + ((B2 /\ 63) << 6) + (B3 /\ 63),
  Code >= 2048,
  ( Code < 55296
  ; Code > 57343
  ),
  !.
utf8_byte_char(Byte, [B2, B3, B4 | Rest], Code, Rest) :-
  Byte >= 240,
  Byte =< 244,
  utf8_continuation_byte(B2),
  utf8_continuation_byte(B3),
  utf8_continuation_byte(B4),
  Code is ((Byte /\ 7) << 18) + ((B2 /\ 63) << 12) + ((B3 /\ 63) << 6) + (B4 /\ 63),
  Code >= 65536,
  Code =< 1114111,
  !.
utf8_byte_char(_, Rest, 65533, Rest).

utf8_continuation_byte(Byte) :-
  byte_value(Byte),
  Byte >= 128,
  Byte =< 191.

byte_value(Byte) :-
  integer(Byte),
  Byte >= 0,
  Byte =< 255,
  !.
byte_value(Byte) :-
  throw(error(type_error(byte, Byte), string_bytes/3)).

text_codec_call(Command, Payload, Response) :-
  setup_call_cleanup(
    open('/v1/dev/codec/text', read_write, Stream, [type(text)]),
    (text_codec_write_request(Stream, Command, Payload),
     read_term(Stream, Response, [])),
    close(Stream)),
  !.

text_codec_write_request(Stream, Command, Payload) :-
  atom_chars(Command, CommandChars),
  text_codec_put_chars(Stream, CommandChars),
  put_char(Stream, '\n'),
  write_term(Stream, Payload, [quoted(true)]),
  put_char(Stream, '.').

text_codec_response(_, ok(Value), Target) :-
  !,
  copy_term(Value, Copy),
  Target = Copy.
text_codec_response(Context, error(Formal), _) :-
  !,
  throw(error(Formal, Context)).
text_codec_response(Context, _, _) :-
  throw(error(system_error, Context)).

text_codec_put_chars(_, []).
text_codec_put_chars(Stream, [Char | Rest]) :-
  put_char(Stream, Char),
  text_codec_put_chars(Stream, Rest).

% write_term(+Term, +Options) is det.
%
% Writes Term to the current output stream with Options.
write_term(Term, Options) :-
  current_output(S),
  write_term(S, Term, Options).

%! write(+Term) is det.
%
% Writes Term to the current output stream.
write(Term) :-
  current_output(S),
  write(S, Term).

% write(+Stream, +Term) is det.
%
% Writes Term to Stream with default write options.
write(Stream, Term) :-
  write_term(Stream, Term, [numbervars(true)]).

%! writeq(+Term) is det.
%
% Writes Term to the current output stream using quoted syntax.
writeq(Term) :-
  current_output(S),
  writeq(S, Term).

% writeq(+Stream, +Term) is det.
%
% Writes Term to Stream using quoted syntax.
writeq(Stream, Term) :-
  write_term(Stream, Term, [quoted(true), numbervars(true)]).

%! write_canonical(+Term) is det.
%
% Writes Term to the current output stream in canonical form.
write_canonical(Term) :-
  current_output(S),
  write_canonical(S, Term).

% write_canonical(+Stream, +Term) is det.
%
% Writes Term to Stream in canonical form.
write_canonical(Stream, Term) :-
  write_term(Stream, Term, [quoted(true), ignore_ops(true)]).

%! term_to_atom(?Term, ?Atom) is det.
%
% Relates a ground Term with its textual Atom representation.
%
% where:
%
% - Term is a ground term that unifies with the parsed representation of Atom;
% - Atom is an atom containing a canonical textual representation of Term.
%
% When Term is ground, Atom is unified with a canonical textual representation
% that can be parsed back by this predicate. When Atom is instantiated, it is
% parsed back into Term.
%
% The supported syntax matches the canonical text produced here: atoms, quoted
% atoms, numbers, double-quoted strings (lists of one-character atoms), lists and compounds.
%
% Throws:
%
% - error(instantiation_error, term_to_atom/2) when both arguments are variables;
% - error(type_error(atom, Atom), term_to_atom/2) when Atom is instantiated but is not an atom;
% - error(syntax_error(term), term_to_atom/2) when Atom is an atom that does not contain a valid canonical term.
term_to_atom(Term, Atom) :-
  ( nonvar(Atom),
    \+atom(Atom)
  -> throw(error(type_error(atom, Atom), term_to_atom/2))
  ; ground(Term)
  -> term_to_atom_chars(Term, Chars),
     atom_chars(Atom, Chars)
  ; atom(Atom)
  -> atom_chars(Atom, Chars),
     ( Chars = []
     -> Term = ''
     ; phrase(term_to_atom_term(Term), Chars)
     -> true
     ; throw(error(syntax_error(term), term_to_atom/2))
     )
  ; throw(error(instantiation_error, term_to_atom/2))
  ).

term_to_atom_chars([], [Open, Close]) :-
  term_to_atom_lbracket_char(Open),
  term_to_atom_rbracket_char(Close),
  !.
term_to_atom_chars(Term, Chars) :-
  atom(Term),
  !,
  term_to_atom_atom_chars(Term, Chars).
term_to_atom_chars(Term, Chars) :-
  number(Term),
  !,
  number_chars(Term, Chars).
term_to_atom_chars(Term, Chars) :-
  term_to_atom_is_char_list(Term),
  !,
  term_to_atom_double_quoted_string_chars(Term, Chars).
term_to_atom_chars([Head|Tail], [Open|Chars]) :-
  term_to_atom_lbracket_char(Open),
  !,
  term_to_atom_list_chars([Head|Tail], Chars).
term_to_atom_chars(Term, Chars) :-
  compound(Term),
  functor(Term, Functor, Arity),
  term_to_atom_atom_chars(Functor, FunctorChars),
  term_to_atom_lparen_char(Open),
  append(FunctorChars, [Open|ArgsChars], Chars),
  term_to_atom_args_chars(1, Arity, Term, ArgsChars).

term_to_atom_atom_chars(Atom, Chars) :-
  atom_chars(Atom, RawChars),
  ( term_to_atom_is_bare_atom_chars(RawChars)
  -> Chars = RawChars
  ; term_to_atom_is_symbol_atom_chars(RawChars)
  -> Chars = RawChars
  ; term_to_atom_quoted_atom_chars(RawChars, Chars)
  ).

term_to_atom_is_bare_atom_chars([Head|Tail]) :-
  term_to_atom_lowercase_char(Head),
  term_to_atom_is_bare_atom_tail_chars(Tail).

term_to_atom_is_bare_atom_tail_chars([]).
term_to_atom_is_bare_atom_tail_chars([Head|Tail]) :-
  term_to_atom_identifier_char(Head),
  term_to_atom_is_bare_atom_tail_chars(Tail).

term_to_atom_is_symbol_atom_chars([Head|Tail]) :-
  term_to_atom_symbol_char(Head),
  term_to_atom_is_symbol_atom_tail_chars(Tail).

term_to_atom_is_symbol_atom_tail_chars([]).
term_to_atom_is_symbol_atom_tail_chars([Head|Tail]) :-
  term_to_atom_token_char(Head),
  term_to_atom_is_symbol_atom_tail_chars(Tail).

term_to_atom_quoted_atom_chars(RawChars, [Quote|Chars]) :-
  term_to_atom_quote_char(Quote),
  term_to_atom_quoted_atom_body_chars(RawChars, BodyChars),
  append(BodyChars, [Quote], Chars).

term_to_atom_quoted_atom_body_chars([], []).
term_to_atom_quoted_atom_body_chars([Head|Tail], [Escape, Quote|Chars]) :-
  term_to_atom_quote_char(Head),
  !,
  term_to_atom_backslash_char(Escape),
  term_to_atom_quote_char(Quote),
  term_to_atom_quoted_atom_body_chars(Tail, Chars).
term_to_atom_quoted_atom_body_chars([Head|Tail], [Escape, Backslash|Chars]) :-
  term_to_atom_backslash_char(Head),
  !,
  term_to_atom_backslash_char(Escape),
  term_to_atom_backslash_char(Backslash),
  term_to_atom_quoted_atom_body_chars(Tail, Chars).
term_to_atom_quoted_atom_body_chars([Head|Tail], [Head|Chars]) :-
  term_to_atom_quoted_atom_body_chars(Tail, Chars).

% A "string" is represented as a proper list of one-character atoms.
term_to_atom_is_char_list([]).
term_to_atom_is_char_list([H|T]) :-
  atom(H),
  atom_length(H, 1),
  term_to_atom_is_char_list(T).

term_to_atom_double_quoted_string_chars(String, [Quote|Chars]) :-
  term_to_atom_double_quote_char(Quote),
  term_to_atom_double_quoted_body_chars(String, BodyChars),
  append(BodyChars, [Quote], Chars).

term_to_atom_double_quoted_body_chars([], []).
term_to_atom_double_quoted_body_chars([Head|Tail], [Escape, Quote|Chars]) :-
  term_to_atom_double_quote_char(Head),
  !,
  term_to_atom_backslash_char(Escape),
  term_to_atom_double_quote_char(Quote),
  term_to_atom_double_quoted_body_chars(Tail, Chars).
term_to_atom_double_quoted_body_chars([Head|Tail], [Escape, Backslash|Chars]) :-
  term_to_atom_backslash_char(Head),
  !,
  term_to_atom_backslash_char(Escape),
  term_to_atom_backslash_char(Backslash),
  term_to_atom_double_quoted_body_chars(Tail, Chars).
term_to_atom_double_quoted_body_chars([Head|Tail], [Head|Chars]) :-
  term_to_atom_double_quoted_body_chars(Tail, Chars).

term_to_atom_list_chars([], [Close]) :-
  term_to_atom_rbracket_char(Close).
term_to_atom_list_chars([Head|Tail], Chars) :-
  term_to_atom_chars(Head, HeadChars),
  append(HeadChars, TailChars, Chars),
  term_to_atom_list_tail_chars(Tail, TailChars).

term_to_atom_list_tail_chars([], [Close]) :-
  term_to_atom_rbracket_char(Close),
  !.
term_to_atom_list_tail_chars([Head|Tail], [Comma|Chars]) :-
  term_to_atom_comma_char(Comma),
  !,
  term_to_atom_chars(Head, HeadChars),
  append(HeadChars, TailChars, Chars),
  term_to_atom_list_tail_chars(Tail, TailChars).
term_to_atom_list_tail_chars(Tail, [Pipe|Chars]) :-
  term_to_atom_pipe_char(Pipe),
  term_to_atom_chars(Tail, TailChars),
  term_to_atom_rbracket_char(Close),
  append(TailChars, [Close], Chars).

term_to_atom_args_chars(Index, Arity, Term, Chars) :-
  arg(Index, Term, Arg),
  term_to_atom_chars(Arg, ArgChars),
  ( Index =:= Arity
  -> term_to_atom_rparen_char(Close),
     append(ArgChars, [Close], Chars)
  ; term_to_atom_comma_char(Comma),
    append(ArgChars, [Comma|TailChars], Chars),
    Next is Index + 1,
    term_to_atom_args_chars(Next, Arity, Term, TailChars)
  ).

term_to_atom_term(Term) -->
  term_to_atom_blanks,
  term_to_atom_value(Term),
  term_to_atom_blanks.

term_to_atom_value([]) -->
  [Open],
  {term_to_atom_lbracket_char(Open)},
  term_to_atom_blanks,
  [Close],
  {term_to_atom_rbracket_char(Close)},
  !.
term_to_atom_value(List) -->
  [Open],
  {term_to_atom_lbracket_char(Open)},
  term_to_atom_blanks,
  term_to_atom_list_value(List),
  term_to_atom_blanks,
  [Close],
  {term_to_atom_rbracket_char(Close)},
  !.
term_to_atom_value(String) -->
  term_to_atom_double_quoted_string(String),
  !.
term_to_atom_value(Number) -->
  term_to_atom_number(Number),
  !.
term_to_atom_value(Term) -->
  term_to_atom_atom_or_compound(Term).

term_to_atom_list_value([Head|Tail]) -->
  term_to_atom_value(Head),
  term_to_atom_blanks,
  ( [Comma],
    {term_to_atom_comma_char(Comma)},
    term_to_atom_blanks,
    term_to_atom_list_value(Tail)
  ; [Pipe],
    {term_to_atom_pipe_char(Pipe)},
    term_to_atom_blanks,
    term_to_atom_value(Tail)
  ; {Tail = []}
  ).

term_to_atom_double_quoted_string(String) -->
  [Quote],
  {term_to_atom_double_quote_char(Quote)},
  term_to_atom_double_quoted_chars(Chars),
  [Quote],
  {String = Chars}.

term_to_atom_double_quoted_chars([Char|Chars]) -->
  term_to_atom_double_quoted_char(Char),
  !,
  term_to_atom_double_quoted_chars(Chars).
term_to_atom_double_quoted_chars([]) -->
  [].

term_to_atom_double_quoted_char('"') -->
  [Escape, Quote],
  {term_to_atom_backslash_char(Escape), term_to_atom_double_quote_char(Quote)}.
term_to_atom_double_quoted_char(Char) -->
  [Escape, Char],
  {term_to_atom_backslash_char(Escape), term_to_atom_backslash_char(Char)}.
term_to_atom_double_quoted_char(Char) -->
  [Char],
  {\+ term_to_atom_double_quote_char(Char), \+ term_to_atom_backslash_char(Char)}.

term_to_atom_number(Number) -->
  term_to_atom_token(Token),
  {
    Token \= [],
    catch(number_chars(Number, Token), _Error, fail)
  }.

term_to_atom_atom_or_compound(Term) -->
  term_to_atom_functor(Functor),
  term_to_atom_blanks,
  ( [Open],
    {term_to_atom_lparen_char(Open)},
    term_to_atom_blanks,
    term_to_atom_arguments(Args),
    term_to_atom_blanks,
    [Close],
    {term_to_atom_rparen_char(Close)},
    {Term =.. [Functor|Args]}
  ; {Term = Functor}
  ).

term_to_atom_arguments([Arg|Args]) -->
  term_to_atom_value(Arg),
  term_to_atom_blanks,
  ( [Comma],
    {term_to_atom_comma_char(Comma)},
    term_to_atom_blanks,
    term_to_atom_arguments(Args)
  ; {Args = []}
  ).

term_to_atom_functor(Functor) -->
  term_to_atom_quoted_atom(Functor),
  !.
term_to_atom_functor(Functor) -->
  term_to_atom_bare_atom(Functor),
  !.
term_to_atom_functor(Functor) -->
  term_to_atom_symbol_atom(Functor).

term_to_atom_quoted_atom(Atom) -->
  [Quote],
  {term_to_atom_quote_char(Quote)},
  term_to_atom_single_quoted_chars(Chars),
  [Quote],
  {atom_chars(Atom, Chars)}.

term_to_atom_single_quoted_chars([Char|Chars]) -->
  term_to_atom_single_quoted_char(Char),
  !,
  term_to_atom_single_quoted_chars(Chars).
term_to_atom_single_quoted_chars([]) -->
  [].

term_to_atom_single_quoted_char(Char) -->
  [Escape, Char],
  {term_to_atom_backslash_char(Escape), term_to_atom_quote_char(Char)}.
term_to_atom_single_quoted_char(Char) -->
  [Escape, Char],
  {term_to_atom_backslash_char(Escape), term_to_atom_backslash_char(Char)}.
term_to_atom_single_quoted_char(Char) -->
  [Char],
  {\+ term_to_atom_quote_char(Char), \+ term_to_atom_backslash_char(Char)}.

term_to_atom_bare_atom(Atom) -->
  term_to_atom_bare_atom_chars(Chars),
  {atom_chars(Atom, Chars)}.

term_to_atom_symbol_atom(Atom) -->
  term_to_atom_symbol_atom_chars(Chars),
  {atom_chars(Atom, Chars)}.

term_to_atom_bare_atom_chars([Head|Tail]) -->
  [Head],
  {term_to_atom_lowercase_char(Head)},
  term_to_atom_bare_atom_tail_chars(Tail).

term_to_atom_bare_atom_tail_chars([Head|Tail]) -->
  [Head],
  {term_to_atom_identifier_char(Head)},
  !,
  term_to_atom_bare_atom_tail_chars(Tail).
term_to_atom_bare_atom_tail_chars([]) -->
  [].

term_to_atom_symbol_atom_chars([Head|Tail]) -->
  [Head],
  {term_to_atom_symbol_char(Head)},
  term_to_atom_symbol_atom_tail_chars(Tail).

term_to_atom_symbol_atom_tail_chars([Head|Tail]) -->
  [Head],
  {term_to_atom_token_char(Head)},
  !,
  term_to_atom_symbol_atom_tail_chars(Tail).
term_to_atom_symbol_atom_tail_chars([]) -->
  [].

term_to_atom_token([Head|Tail]) -->
  [Head],
  {term_to_atom_token_char(Head)},
  term_to_atom_token_tail(Tail).

term_to_atom_token_tail([Head|Tail]) -->
  [Head],
  {term_to_atom_token_char(Head)},
  !,
  term_to_atom_token_tail(Tail).
term_to_atom_token_tail([]) -->
  [].

term_to_atom_blanks -->
  [Char],
  {term_to_atom_blank(Char)},
  !,
  term_to_atom_blanks.
term_to_atom_blanks -->
  [].

term_to_atom_blank(' ').
term_to_atom_blank('\n').
term_to_atom_blank('\r').
term_to_atom_blank('\t').

term_to_atom_token_char(Char) :-
  \+term_to_atom_blank(Char),
  \+term_to_atom_lparen_char(Char),
  \+term_to_atom_rparen_char(Char),
  \+term_to_atom_lbracket_char(Char),
  \+term_to_atom_rbracket_char(Char),
  \+term_to_atom_comma_char(Char),
  \+term_to_atom_pipe_char(Char).

term_to_atom_identifier_char(Char) :-
  term_to_atom_lowercase_char(Char).
term_to_atom_identifier_char(Char) :-
  term_to_atom_uppercase_char(Char).
term_to_atom_identifier_char(Char) :-
  term_to_atom_digit_char(Char).
term_to_atom_identifier_char('_').

term_to_atom_symbol_char(Char) :-
  term_to_atom_token_char(Char),
  \+term_to_atom_identifier_char(Char).

term_to_atom_lowercase_char(Char) :-
  char_code(Char, Code),
  Code >= 97,
  Code =< 122.

term_to_atom_uppercase_char(Char) :-
  char_code(Char, Code),
  Code >= 65,
  Code =< 90.

term_to_atom_digit_char(Char) :-
  char_code(Char, Code),
  Code >= 48,
  Code =< 57.

term_to_atom_quote_char(Char) :-
  char_code(Char, 39).

term_to_atom_backslash_char(Char) :-
  char_code(Char, 92).

term_to_atom_double_quote_char(Char) :-
  char_code(Char, 34).

term_to_atom_lparen_char(Char) :-
  char_code(Char, 40).

term_to_atom_rparen_char(Char) :-
  char_code(Char, 41).

term_to_atom_comma_char(Char) :-
  char_code(Char, 44).

term_to_atom_lbracket_char(Char) :-
  char_code(Char, 91).

term_to_atom_rbracket_char(Char) :-
  char_code(Char, 93).

term_to_atom_pipe_char(Char) :-
  char_code(Char, 124).

%! source_file(?File) is nondet.
%
% True when File is one of the Prolog source files loaded in the current
% interpreter.
source_file(File) :-
  ( var(File)
  -> source_file_sources(Files),
     source_file_member(File, Files)
  ; atom(File)
  -> source_file_sources(Files),
     source_file_member(File, Files)
  ;  throw(error(type_error(atom, File), source_file/1))
  ).

source_file_sources(Files) :-
  setup_call_cleanup(
    open('/v1/run/source/files', read, Stream, [type(text)]),
    read_term(Stream, Files, []),
    close(Stream)
  ).

source_file_member(File, [File | _]).
source_file_member(File, [_ | Rest]) :-
  source_file_member(File, Rest).

%! atomic_list_concat(+List, ?Atom) is det.
%
% Unifies Atom with the concatenation of the atomic textual representation of
% each element in List.
%
% where:
%
% - List is a proper list of ground terms. Each element is converted using term_to_atom/2, so atoms, numbers,
%   double-quoted strings, lists and compounds are supported;
% - Atom is an atom representing the concatenation of the textual representation of each element in List.
%
% Throws:
%
% - error(instantiation_error, atomic_list_concat/2) when List is insufficiently instantiated;
% - error(type_error(list, List), atomic_list_concat/2) when List is not a proper list.
atomic_list_concat(List, Atom) :-
  atomic_list_concat_must_be_list(List, atomic_list_concat/2),
  atomic_list_concat_parts(List, '', Atom).

%! atomic_list_concat(+List, +Separator, ?Atom) is det.
%
% Unifies Atom with the concatenation of the atomic textual representation of
% each element in List, inserting Separator between adjacent elements.
%
% where:
%
% - List is a proper list of ground terms. Each element is converted using term_to_atom/2, so atoms, numbers,
%   double-quoted strings, lists and compounds are supported;
% - Separator is an atom inserted between adjacent elements;
% - Atom is an atom representing the concatenation of the textual representation of each element in List.
%
% Throws:
%
% - error(instantiation_error, atomic_list_concat/3) when List or Separator is insufficiently instantiated;
% - error(type_error(list, List), atomic_list_concat/3) when List is not a proper list;
% - error(type_error(atom, Separator), atomic_list_concat/3) when Separator is instantiated but is not an atom.
atomic_list_concat(List, Separator, Atom) :-
  atomic_list_concat_must_be_list(List, atomic_list_concat/3),
  ( var(Separator)
  -> throw(error(instantiation_error, atomic_list_concat/3))
  ; atom(Separator)
  -> atomic_list_concat_parts(List, Separator, Atom)
  ; throw(error(type_error(atom, Separator), atomic_list_concat/3))
  ).

atomic_list_concat_parts([], _, '').
atomic_list_concat_parts([Head|Tail], Separator, Atom) :-
  atom_chars(Separator, SeparatorChars),
  term_to_atom_chars(Head, HeadChars),
  atomic_list_concat_collect_chars(Tail, SeparatorChars, [HeadChars], CharLists),
  atomic_list_concat_flatten_fast(CharLists, Chars),
  atom_chars(Atom, Chars).

% Collect character lists in reverse order
atomic_list_concat_collect_chars([], _, Acc, Acc).
atomic_list_concat_collect_chars([Head|Tail], SeparatorChars, Acc, CharLists) :-
  term_to_atom_chars(Head, HeadChars),
  atomic_list_concat_collect_chars(Tail, SeparatorChars, [HeadChars, SeparatorChars|Acc], CharLists).

% Flatten by computing total length first, then building result
atomic_list_concat_flatten_fast(CharLists, Chars) :-
  atomic_list_concat_flatten_reverse(CharLists, [], Chars).

% Flatten from right to left (reversed list) to avoid O(n²) appends
% This is efficient because we append small character lists (one element at a time)
% to a growing accumulator. Since append(SmallList, BigList, Result) is O(|SmallList|),
% and each character list is small (length of one element's representation),
% the total complexity is O(total_chars) which is optimal.
atomic_list_concat_flatten_reverse([], Acc, Acc).
atomic_list_concat_flatten_reverse([Chars|Rest], Acc, Result) :-
  append(Chars, Acc, NewAcc),
  atomic_list_concat_flatten_reverse(Rest, NewAcc, Result).

atomic_list_concat_must_be_list(List, Context) :-
  ( var(List)
  -> throw(error(instantiation_error, Context))
  ; atomic_list_concat_is_partial_list(List)
  -> throw(error(instantiation_error, Context))
  ; atomic_list_concat_is_list(List)
  -> true
  ; throw(error(type_error(list, List), Context))
  ).

atomic_list_concat_is_list([]).
atomic_list_concat_is_list([_|Tail]) :-
  nonvar(Tail),
  atomic_list_concat_is_list(Tail).

atomic_list_concat_is_partial_list([_|Tail]) :-
  ( var(Tail)
  ; nonvar(Tail),
    atomic_list_concat_is_partial_list(Tail)
  ).

% memberchk(?Elem, +List) is semidet.
%
% Succeeds if Elem is a member of List. This is a deterministic predicate
% that commits to the first unification and does not leave a choice point.
% Useful when List is ground and you only need to check membership once.
memberchk(X, [X|_]) :-
  !.
memberchk(X, [_|T]) :-
  memberchk(X, T).
