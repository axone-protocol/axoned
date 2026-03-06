% error.pl
% Error-handling predicates inspired by SWI-Prolog's library(error).

:- consult('/v1/lib/type.pl').

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
