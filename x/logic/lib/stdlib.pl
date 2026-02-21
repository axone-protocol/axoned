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
