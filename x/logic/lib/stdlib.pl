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
  (E = at; E = past).

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

% put_char(+Char) is det.
%
% Writes Char to the current output stream.
put_char(Char) :-
  current_output(S),
  put_char(S, Char).

% put_code(+Code) is det.
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

% nl is det.
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

% write_term(+Term, +Options) is det.
%
% Writes Term to the current output stream with Options.
write_term(Term, Options) :-
  current_output(S),
  write_term(S, Term, Options).

% write(+Term) is det.
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

% writeq(+Term) is det.
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

% write_canonical(+Term) is det.
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
