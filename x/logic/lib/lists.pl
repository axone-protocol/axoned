% lists.pl
% List processing predicates inspired by SWI-Prolog's library(lists).

% member(?Elem, ?List) is nondet.
%
% True if Elem unifies with an element of List.
member(Elem, [Elem|_]).
member(Elem, [_|Tail]) :-
  member(Elem, Tail).

% select(?Elem, ?List1, ?List2) is nondet.
%
% True when List2 is List1 with one occurrence of Elem removed.
select(Elem, [Elem|Tail], Tail).
select(Elem, [Head|Tail], [Head|Rest]) :-
  select(Elem, Tail, Rest).
