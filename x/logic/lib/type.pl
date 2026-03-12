% type.pl
% Type-checking predicates.

known_type(any).
known_type(atom).
known_type(atomic).
known_type(boolean).
known_type(callable).
known_type(char).
known_type(chars).
known_type(code).
known_type(codes).
known_type(byte).
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

%! has_type(+Type, @Term) is semidet.
%
% Succeeds when Term satisfies Type without throwing.
% Fails when Type is known but Term does not match it.
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
has_type(byte, X) :- integer(X), X >= 0, X =< 255.
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
  number(L), number(U),
  ( integer(L), integer(U) ->
      integer(X)
  ; number(X)
  ),
  X >= L,
  X =< U.

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
