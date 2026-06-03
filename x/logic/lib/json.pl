% json.pl
% JSON helpers backed by the codec VFS device.

:- consult('/v1/lib/dev.pl').
:- consult('/v1/lib/error.pl').

%! json_prolog(?Json, ?Term) is det.
%
% Relates JSON text with its canonical Prolog representation.
%
% Json is text: an atom, a list of characters, or a list of character codes.
%
% The canonical representation for Term is:
% - JSON objects are represented as `json(NameValueList)`;
% - JSON arrays are represented as Prolog lists;
% - JSON strings are represented as atoms;
% - JSON numbers are represented as numbers;
% - JSON booleans and null are represented as `@(true)`, `@(false)`, and `@(null)`.
json_prolog(Json, Term) :-
  ( nonvar(Json)
  -> with_context(json_prolog/2, string_bytes(Json, JsonBytes, text)),
     json_decode_bytes(json_prolog/2, JsonBytes, Term)
  ; nonvar(Term)
  -> json_encode_term(json_prolog/2, Term, Json)
  ;  throw(error(instantiation_error, json_prolog/2))
  ).

%! json_read(+Stream, ?Term) is det.
%
% Reads JSON text from Stream and unifies Term with its canonical Prolog representation.
json_read(Stream, Term) :-
  with_context(json_read/2, read_string(Stream, _, Json)),
  with_context(json_read/2, string_bytes(Json, JsonBytes, text)),
  json_decode_bytes(json_read/2, JsonBytes, Term).

%! json_write(+Stream, +Term) is det.
%
% Writes Term as JSON text to Stream.
json_write(Stream, Term) :-
  json_encode_term(json_write/2, Term, Json),
  json_write_atom(Stream, Json).

json_decode_bytes(Context, JsonBytes, Term) :-
  with_context(Context, string_bytes(JsonChars, JsonBytes, text)),
  json_codec_call([d,e,c,o,d,e,'\n'], JsonChars, Response),
  json_response(Context, Response, Term).

json_encode_term(Context, Term, Json) :-
  with_context(Context, term_to_atom(Term, TermAtom)),
  atom_chars(TermAtom, TermChars),
  append(TermChars, ['.'], PayloadChars),
  json_codec_call([e,n,c,o,d,e,'\n'], PayloadChars, Response),
  json_response(Context, Response, Json).

json_codec_call(CommandChars, PayloadChars, Response) :-
  dev_call(
    '/v1/dev/codec/json',
    text,
    json_write_request(CommandChars, PayloadChars),
    json_read_response(Response)
  ).

json_write_request(Stream, CommandChars, PayloadChars) :-
  json_put_chars(Stream, CommandChars),
  json_put_chars(Stream, PayloadChars).

json_read_response(Stream, Response) :-
  read_term(Stream, Response, []).

json_response(_, ok(Value), Value) :-
  !.
json_response(Context, error(Error), _) :-
  !,
  throw(error(Error, Context)).
json_response(Context, _, _) :-
  throw(error(system_error, Context)).

json_write_atom(Stream, Atom) :-
  atom_chars(Atom, Chars),
  json_put_chars(Stream, Chars).

json_put_chars(_, []).
json_put_chars(Stream, [Char | Rest]) :-
  put_char(Stream, Char),
  json_put_chars(Stream, Rest).
