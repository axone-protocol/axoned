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

member(X, [X|_]).
member(X, [_|Xs]) :- member(X, Xs).

select(E, [E|Xs], Xs).
select(E, [X|Xs], [X|Ys]) :-
  select(E, Xs, Ys).

maplist(_Cont_1, []).
maplist(Cont_1, [E1|E1s]) :-
  call(Cont_1, E1),
  maplist(Cont_1, E1s).

maplist(_Cont_2, [], []).
maplist(Cont_2, [E1|E1s], [E2|E2s]) :-
  call(Cont_2, E1, E2),
  maplist(Cont_2, E1s, E2s).

maplist(_Cont_3, [], [], []).
maplist(Cont_3, [E1|E1s], [E2|E2s], [E3|E3s]) :-
  call(Cont_3, E1, E2, E3),
  maplist(Cont_3, E1s, E2s, E3s).

maplist(_Cont_4, [], [], [], []).
maplist(Cont_4, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s]) :-
  call(Cont_4, E1, E2, E3, E4),
  maplist(Cont_4, E1s, E2s, E3s, E4s).

maplist(_Cont_5, [], [], [], [], []).
maplist(Cont_5, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s], [E5|E5s]) :-
  call(Cont_5, E1, E2, E3, E4, E5),
  maplist(Cont_5, E1s, E2s, E3s, E4s, E5s).

maplist(_Cont_6, [], [], [], [], [], []).
maplist(Cont_6, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s], [E5|E5s], [E6|E6s]) :-
  call(Cont_6, E1, E2, E3, E4, E5, E6),
  maplist(Cont_6, E1s, E2s, E3s, E4s, E5s, E6s).

maplist(_Cont_7, [], [], [], [], [], [], []).
maplist(Cont_7, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s], [E5|E5s], [E6|E6s], [E7|E7s]) :-
  call(Cont_7, E1, E2, E3, E4, E5, E6, E7),
  maplist(Cont_7, E1s, E2s, E3s, E4s, E5s, E6s, E7s).

source_files(Files) :- bagof(File, source_file(File), Files).
