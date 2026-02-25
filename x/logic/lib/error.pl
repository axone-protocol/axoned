% error.pl
% Error-handling predicates inspired by SWI-Prolog's library(error).

%! must_be(+Type, @Term) is det.
%
% Succeeds when Term satisfies Type.
% Throws:
% - error(instantiation_error, must_be/2) when Term is insufficiently instantiated;
% - error(type_error(Type, Term), must_be/2) when Term has the wrong type;
% - error(existence_error(type, Type), must_be/2) when Type is unknown.
must_be(Type, Term) :-
  ( nonvar(Type),
    has_type(Type, Term)
  -> true
  ; nonvar(Type)
  -> is_not(Type, Term)
  ; throw(error(instantiation_error, must_be/2))
  ).

is_not(var, Term) :-
  !,
  throw(error(uninstantiation_error(Term), must_be/2)).
is_not(list, Term) :-
  !,
  ( var(Term)
  -> throw(error(instantiation_error, must_be/2))
  ; partial_list(Term)
  -> throw(error(instantiation_error, must_be/2))
  ; throw(error(type_error(list, Term), must_be/2))
  ).
is_not(list(Of), Term) :-
  !,
  not_a_list(list(Of), Term).
is_not(Type, Term) :-
  known_type(Type),
  !,
  ( var(Term)
  -> throw(error(instantiation_error, must_be/2))
  ; throw(error(type_error(Type, Term), must_be/2))
  ).
is_not(Type, _) :-
  throw(error(existence_error(type, Type), must_be/2)).

known_type(any).
known_type(atom).
known_type(atomic).
known_type(boolean).
known_type(callable).
known_type(char).
known_type(chars).
known_type(code).
known_type(codes).
known_type(compound).
known_type(constant).
known_type(float).
known_type(integer).
known_type(nonneg).
known_type(positive_integer).
known_type(negative_integer).
known_type(nonvar).
known_type(number).
known_type(oneof(_)).
known_type(pair).
known_type(var).
known_type(text).
known_type(list).
known_type(list(_)).
known_type(between(_, _)).

has_type(any, _).
has_type(atom, X) :- atom(X).
has_type(atomic, X) :- atomic(X).
has_type(boolean, X) :- X == true.
has_type(boolean, X) :- X == false.
has_type(callable, X) :- callable(X).
has_type(char, X) :- atom(X), atom_length(X, 1).
has_type(chars, X) :- proper_char_list(X).
has_type(code, X) :- integer(X), X >= 0, X =< 1114111.
has_type(codes, X) :- proper_code_list(X).
has_type(text, X) :- atom(X).
has_type(text, X) :- proper_char_list(X).
has_type(text, X) :- proper_code_list(X).
has_type(compound, X) :- compound(X).
has_type(constant, X) :- atomic(X).
has_type(float, X) :- float(X).
has_type(integer, X) :- integer(X).
has_type(nonneg, X) :- integer(X), X >= 0.
has_type(positive_integer, X) :- integer(X), X > 0.
has_type(negative_integer, X) :- integer(X), X < 0.
has_type(nonvar, X) :- nonvar(X).
has_type(number, X) :- number(X).
has_type(oneof(Choices), X) :- ground(X), memberchk(X, Choices).
has_type(pair, X) :- nonvar(X), X = _-_.
has_type(var, X) :- var(X).
has_type(list, X) :- proper_list(X).
has_type(list(Type), X) :- proper_list(X), element_types(X, Type).
has_type(between(L, U), X) :-
  integer(L), integer(U), integer(X),
  X >= L, X =< U.
has_type(between(L, U), X) :-
  number(L), number(U), number(X),
  X >= L, X =< U.

not_a_list(list(Of), X) :-
  ( var(X)
  -> throw(error(instantiation_error, must_be/2))
  ; partial_list(X)
  -> throw(error(instantiation_error, must_be/2))
  ; proper_list(X)
  -> ( nonvar(Of)
     -> element_is_not(X, Of)
     ; throw(error(instantiation_error, must_be/2))
     )
  ; throw(error(type_error(list(Of), X), must_be/2))
  ).

element_is_not([H|T], Of) :-
  has_type(Of, H),
  !,
  element_is_not(T, Of).
element_is_not([H|_], Of) :-
  !,
  is_not(Of, H).
element_is_not([], _).

element_types([], _).
element_types([H|T], Type) :-
  has_type(Type, H),
  element_types(T, Type).

proper_list([]).
proper_list([_|T]) :-
  nonvar(T),
  proper_list(T).

partial_list([_|T]) :-
  ( var(T)
  ; nonvar(T), partial_list(T)
  ).

proper_char_list([]).
proper_char_list([H|T]) :-
  has_type(char, H),
  nonvar(T),
  proper_char_list(T).

proper_code_list([]).
proper_code_list([H|T]) :-
  has_type(code, H),
  nonvar(T),
  proper_code_list(T).
