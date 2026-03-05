% dev.pl
% Generic helpers for transactional VFS devices.

:- consult('/v1/lib/error.pl').

%! dev_call(+Path, +Type, :WriteGoal, :ReadGoal) is det.
%
% Executes a transactional device call:
% 1. open device stream in read_write mode
% 2. run WriteGoal(Stream)
% 3. run ReadGoal(Stream)
% 4. close stream
%
% The transactional commit is device-defined and is typically triggered by the
% first read operation performed inside ReadGoal/1.
dev_call(Path, Type, WriteGoal, ReadGoal) :-
  setup_call_cleanup(
    dev_open(Path, Type, Stream),
    (
      apply_stream(WriteGoal, Stream),
      apply_stream(ReadGoal, Stream)
    ),
    close(Stream)
  ).

% apply_stream(+Goal, +Stream) is det.
%
% Calls Goal with Stream prepended as first argument.
%
% Examples:
% - apply_stream(write_bytes(Bytes), S) => write_bytes(S, Bytes)
% - apply_stream(read_bytes(Bytes),  S) => read_bytes(S, Bytes)
apply_stream(Goal, Stream) :-
  must_be(callable, Goal),
  Goal =.. [Functor | Args],
  GoalWithStream =.. [Functor, Stream | Args],
  call(GoalWithStream).

% dev_open(+Path, +Type, -Stream) is det.
%
% Opens a transactional device stream in read_write mode.
%
% Type must be one of:
% - text
% - binary
dev_open(Path, Type, Stream) :-
  must_be(atom, Type),
  must_be(oneof([text, binary]), Type),
  open(Path, read_write, Stream, [type(Type)]).

% dev_write_bytes(+Stream, +Bytes) is det.
%
% Writes all bytes from Bytes (a list of integers in [0,255]) to Stream.
dev_write_bytes(Stream, Bytes) :-
  must_be(list(byte), Bytes),
  dev_write_bytes_(Stream, Bytes).

dev_write_bytes_(_, []).
dev_write_bytes_(Stream, [Byte | Rest]) :-
  put_byte(Stream, Byte),
  dev_write_bytes_(Stream, Rest).

% dev_read_bytes(+Stream, -Bytes) is det.
%
% Reads all bytes from Stream until EOF and unifies them with Bytes.
dev_read_bytes(Stream, Bytes) :-
  get_byte(Stream, Byte),
  ( Byte =:= -1 ->
      Bytes = []
  ; Bytes = [Byte | Rest],
    dev_read_bytes(Stream, Rest)
  ).
