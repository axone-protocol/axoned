% apply.pl
% Goal-application predicates inspired by SWI-Prolog's library(apply).

%! maplist(:Goal, ?List1) is semidet.
%
% Applies Goal to each element of List1.
maplist(_Goal, []).
maplist(Goal, [E1|E1s]) :-
  call(Goal, E1),
  maplist(Goal, E1s).

%! maplist(:Goal, ?List1, ?List2) is semidet.
%
% Applies Goal to pairs of elements from List1 and List2.
maplist(_Goal, [], []).
maplist(Goal, [E1|E1s], [E2|E2s]) :-
  call(Goal, E1, E2),
  maplist(Goal, E1s, E2s).

%! maplist(:Goal, ?List1, ?List2, ?List3) is semidet.
%
% Applies Goal to triples of elements from List1, List2 and List3.
maplist(_Goal, [], [], []).
maplist(Goal, [E1|E1s], [E2|E2s], [E3|E3s]) :-
  call(Goal, E1, E2, E3),
  maplist(Goal, E1s, E2s, E3s).

%! maplist(:Goal, ?List1, ?List2, ?List3, ?List4) is semidet.
%
% Applies Goal to 4-tuples of elements from the 4 lists.
maplist(_Goal, [], [], [], []).
maplist(Goal, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s]) :-
  call(Goal, E1, E2, E3, E4),
  maplist(Goal, E1s, E2s, E3s, E4s).

%! maplist(:Goal, ?List1, ?List2, ?List3, ?List4, ?List5) is semidet.
%
% Applies Goal to 5-tuples of elements from the 5 lists.
maplist(_Goal, [], [], [], [], []).
maplist(Goal, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s], [E5|E5s]) :-
  call(Goal, E1, E2, E3, E4, E5),
  maplist(Goal, E1s, E2s, E3s, E4s, E5s).

%! maplist(:Goal, ?List1, ?List2, ?List3, ?List4, ?List5, ?List6) is semidet.
%
% Applies Goal to 6-tuples of elements from the 6 lists.
maplist(_Goal, [], [], [], [], [], []).
maplist(Goal, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s], [E5|E5s], [E6|E6s]) :-
  call(Goal, E1, E2, E3, E4, E5, E6),
  maplist(Goal, E1s, E2s, E3s, E4s, E5s, E6s).

%! maplist(:Goal, ?List1, ?List2, ?List3, ?List4, ?List5, ?List6, ?List7) is semidet.
%
% Applies Goal to 7-tuples of elements from the 7 lists.
maplist(_Goal, [], [], [], [], [], [], []).
maplist(Goal, [E1|E1s], [E2|E2s], [E3|E3s], [E4|E4s], [E5|E5s], [E6|E6s], [E7|E7s]) :-
  call(Goal, E1, E2, E3, E4, E5, E6, E7),
  maplist(Goal, E1s, E2s, E3s, E4s, E5s, E6s, E7s).

%! foldl(:Goal, +List, +V0, -V) is det.
%
% Left-folds List using Goal.
% Goal is called as call(Goal, Elem, Acc0, Acc1).
foldl(_Goal, [], V, V).
foldl(Goal, [Elem|Elems], V0, V) :-
  call(Goal, Elem, V0, V1),
  foldl(Goal, Elems, V1, V).

%! foldl(:Goal, +List1, +List2, +V0, -V) is det.
%
% Left-folds two lists in lockstep using Goal.
% Goal is called as call(Goal, Elem1, Elem2, Acc0, Acc1).
foldl(_Goal, [], [], V, V).
foldl(Goal, [Elem1|Elems1], [Elem2|Elems2], V0, V) :-
  call(Goal, Elem1, Elem2, V0, V1),
  foldl(Goal, Elems1, Elems2, V1, V).

%! foldl(:Goal, +List1, +List2, +List3, +V0, -V) is det.
%
% Left-folds three lists in lockstep using Goal.
% Goal is called as call(Goal, Elem1, Elem2, Elem3, Acc0, Acc1).
foldl(_Goal, [], [], [], V, V).
foldl(Goal, [Elem1|Elems1], [Elem2|Elems2], [Elem3|Elems3], V0, V) :-
  call(Goal, Elem1, Elem2, Elem3, V0, V1),
  foldl(Goal, Elems1, Elems2, Elems3, V1, V).

%! foldl(:Goal, +List1, +List2, +List3, +List4, +V0, -V) is det.
%
% Left-folds four lists in lockstep using Goal.
% Goal is called as call(Goal, Elem1, Elem2, Elem3, Elem4, Acc0, Acc1).
foldl(_Goal, [], [], [], [], V, V).
foldl(Goal, [Elem1|Elems1], [Elem2|Elems2], [Elem3|Elems3], [Elem4|Elems4], V0, V) :-
  call(Goal, Elem1, Elem2, Elem3, Elem4, V0, V1),
  foldl(Goal, Elems1, Elems2, Elems3, Elems4, V1, V).
