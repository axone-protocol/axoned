% Operators

:-(op(1200, xfx, [:-, -->])).
:-(op(1200, fx, [:-, ?-])).
:-(op(1105, xfy, '|')).
:-(op(1100, xfy, ;)).
:-(op(1050, xfy, ->)).
:-(op(1000, xfy, ',')).
:-(op(900, fy, \+)).
:-(op(700, xfx, [=, \=])).
:-(op(700, xfx, [==, \==, @<, @=<, @>, @>=])).
:-(op(700, xfx, =..)).
:-(op(700, xfx, [is, =:=, =\=, <, =<, >, >=])).
:-(op(600, xfy, :)).
:-(op(500, yfx, [+, -, /\, \/])).
:-(op(400, yfx, [*, /, //, div, rem, mod, <<, >>])).
:-(op(200, xfx, **)).
:-(op(200, xfy, ^)).
:-(op(200, fy, [+, -, \])).

% Control constructs

true.

fail :- \+true.

! :- !.

P, Q :- call((P, Q)).

If -> Then; _ :- If, !, Then.
_ -> _; Else :- !, Else.

P; Q :- call((P; Q)).

If -> Then :- If, !, Then.

% Term unification

X \= Y :- \+(X = Y).

% Type testing

atomic(X) :-
  nonvar(X),
  \+compound(X).

nonvar(X) :- \+var(X).

number(X) :- float(X).
number(X) :- integer(X).

callable(X) :- atom(X).
callable(X) :- compound(X).

ground(X) :- term_variables(X, []).

% Term comparison

X @=< Y :- compare(=, X, Y).
X @=< Y :- compare(<, X, Y).

X == Y :- compare(=, X, Y).

X \== Y :- \+(X == Y).

X @< Y :- compare(<, X, Y).

X @> Y :- compare(>, X, Y).

X @>= Y :- compare(>, X, Y).
X @>= Y :- compare(=, X, Y).

% Logic and control

once(P) :- P, !.

false :- fail.

% Atomic term processing

% Implementation defined hooks

halt :- halt(0).

% Consult

[H|T] :- consult([H|T]).

% Definite clause grammar

phrase(GRBody, S0) :- phrase(GRBody, S0, []).

% Prolog prologue

source_files(Files) :- bagof(File, source_file(File), Files).
