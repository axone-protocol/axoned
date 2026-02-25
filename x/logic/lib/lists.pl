% lists.pl
% List processing predicates inspired by SWI-Prolog's library(lists).

%! member(?Elem, ?List) is nondet.
%
% True if Elem unifies with an element of List.
member(Elem, [Elem|_]).
member(Elem, [_|Tail]) :-
  member(Elem, Tail).

%! select(?Elem, ?List1, ?List2) is nondet.
%
% True when List2 is List1 with one occurrence of Elem removed.
select(Elem, [Elem|Tail], Tail).
select(Elem, [Head|Tail], [Head|Rest]) :-
  select(Elem, Tail, Rest).

%! append(+ListOfLists, ?List) is det.
%
% Concatenates a list of lists into a single list.
append([], []).
append([List|Lists], Concatenated) :-
  append(List, Rest, Concatenated),
  append(Lists, Rest).

%! prefix(?Prefix, +List) is nondet.
%
% True if Prefix is a prefix of List.
prefix(Prefix, List) :-
  append(Prefix, _, List).

%! suffix(?Suffix, +List) is nondet.
%
% True if Suffix is a suffix of List.
suffix(Suffix, List) :-
  append(_, Suffix, List).

%! sublist(?SubList, +List) is nondet.
%
% True if SubList is a contiguous sublist of List.
sublist(SubList, List) :-
  append(_, Rest, List),
  append(SubList, _, Rest).
