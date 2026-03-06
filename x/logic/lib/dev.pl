% dev.pl
% Generic helpers for transactional VFS devices.

:- consult('/v1/lib/error.pl').

%! dev_call(+Path, +Type, :WriteGoal, :ReadGoal) is det.
%
% Executes a transactional device call following a half-duplex protocol.
%
% ## Overview
%
% A device is a special type of file in the virtual filesystem that implements
% a transactional request-response protocol. Unlike regular files, devices follow
% a strict half-duplex communication pattern with three distinct phases:
%
% 1. **Request phase**: Write operations accumulate bytes into a request buffer
% 2. **Commit phase**: First read operation commits the accumulated request
% 3. **Response phase**: Subsequent reads stream the response until EOF
%
% Once committed (after the first read), the device transitions to read-only mode
% and rejects any further write attempts with a permission error.
%
% ## Protocol Flow
%
% ```
% 1. open device stream in read_write mode
% 2. run WriteGoal(Stream)  ← builds request
% 3. run ReadGoal(Stream)   ← commits & reads response  
% 4. close stream
% ```
%
% The commit operation is device-specific and executes the actual transaction
% (e.g., a smart contract query, a database call, etc.). Most devices require
% at least one write before the first read; reading without writing typically
% fails with an `invalid_request` error.
%
% ## Arguments
%
% - `Path`: atom representing the device path in the VFS (e.g., '/v1/dev/wasm/...')
% - `Type`: stream type, either `text` or `binary`
% - `WriteGoal`: callable that receives Stream as first argument to build the request
% - `ReadGoal`: callable that receives Stream as first argument to read the response
%
% ## Usage Notes
%
% **⚠️ Advanced Feature**: This predicate provides low-level access to transactional
% devices in the virtual filesystem of the Prolog VM. 
% For most use cases, prefer higher-level predicates like `wasm_query/3` for smart 
% contracts, which provide simpler interfaces.
%
% Use `dev_call/4` only when you need:
% - Direct control over the device protocol
% - Custom request/response handling
% - Integration with devices that don't have specialized predicates
%
% Transactional devices include:
% - Codecs and transforms (bech32, base64, etc.) exposed as devices
% - WASM smart contract interactions (prefer `wasm_query/3` for common cases)
% - Other transactional operations as needed
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
