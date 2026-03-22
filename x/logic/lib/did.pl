% did.pl
% DID Core helpers.

:- consult('/v1/lib/error.pl').

%! did_components(+DID:atom, -Parsed) is det.
%
% Also supports the reverse mode:
% did_components(-DID:atom, +Parsed) is det.
%
% Parse or reconstruct a DID / DID URL compliant with W3C DID Core.
%
% Parsed = did(Method, MethodSpecificId, Path, Query, Fragment)
%
% where:
% - Method is an atom.
% - MethodSpecificId is an atom.
% - Path is an atom including its leading `/` when present, otherwise left unbound.
% - Query is a raw atom without its leading `?` when present, otherwise left unbound.
% - Fragment is a raw atom without its leading `#` when present, otherwise left unbound.
%
% Components are preserved raw. No percent-decoding or URI normalization is performed.
did_components(DID, Parsed) :-
  nonvar(DID),
  !,
  with_context(did_components/2, must_be(atom, DID)),
  did_parse(DID, Parsed).
did_components(DID, Parsed) :-
  nonvar(Parsed),
  !,
  did_build(DID, Parsed).
did_components(_, _) :-
  throw(error(instantiation_error, did_components/2)).

did_parse(DID, Parsed) :-
  atom_chars(DID, Chars),
  ( did_parse_chars(Chars, Method, MethodSpecificId, Path, Query, Fragment)
  -> did_unify_parsed(Parsed, Method, MethodSpecificId, Path, Query, Fragment)
  ;  throw(error(domain_error(encoding(did), DID), did_components/2))
  ).

did_build(DID, Parsed) :-
  did_parse_term(Parsed, Method, MethodSpecificId, Path, Query, Fragment),
  atom_chars(Method, MethodChars),
  did_validate_method_chars(Method, MethodChars),
  atom_chars(MethodSpecificId, MethodSpecificIdChars),
  did_validate_method_specific_id_chars(MethodSpecificId, MethodSpecificIdChars),
  did_build_path_chars(Path, PathChars),
  did_build_query_chars(Query, QueryChars),
  did_build_fragment_chars(Fragment, FragmentChars),
  did_concat_chars(
    [['d', 'i', 'd', ':'], MethodChars, [':'], MethodSpecificIdChars, PathChars, QueryChars, FragmentChars],
    DIDChars
  ),
  atom_chars(DID, DIDChars).

did_unify_parsed(Parsed, Method, MethodSpecificId, Path, Query, Fragment) :-
  var(Parsed),
  !,
  did_optional_term(Path, PathTerm),
  did_optional_term(Query, QueryTerm),
  did_optional_term(Fragment, FragmentTerm),
  Parsed = did(Method, MethodSpecificId, PathTerm, QueryTerm, FragmentTerm).
did_unify_parsed(Parsed, Method, MethodSpecificId, Path, Query, Fragment) :-
  Parsed = did(ParsedMethod, ParsedMethodSpecificId, ParsedPath, ParsedQuery, ParsedFragment),
  ParsedMethod = Method,
  ParsedMethodSpecificId = MethodSpecificId,
  did_match_optional_term(ParsedPath, Path),
  did_match_optional_term(ParsedQuery, Query),
  did_match_optional_term(ParsedFragment, Fragment).

did_parse_term(Parsed, Method, MethodSpecificId, Path, Query, Fragment) :-
  Parsed = did(Method, MethodSpecificId, Path, Query, Fragment),
  !,
  with_context(did_components/2, must_be(nonvar, Method)),
  with_context(did_components/2, must_be(atom, Method)),
  with_context(did_components/2, must_be(nonvar, MethodSpecificId)),
  with_context(did_components/2, must_be(atom, MethodSpecificId)).
did_parse_term(Parsed, _, _, _, _, _) :-
  throw(error(type_error(did, Parsed), did_components/2)).

did_optional_term(absent, _).
did_optional_term(present(Value), Value).

did_match_optional_term(Term, absent) :-
  var(Term).
did_match_optional_term(Term, present(Value)) :-
  Term = Value.

did_build_path_chars(Path, []) :-
  var(Path),
  !.
did_build_path_chars(Path, PathChars) :-
  with_context(did_components/2, must_be(atom, Path)),
  atom_chars(Path, PathChars),
  ( PathChars = ['/' | Tail]
  -> did_validate_path_chars(Path, Tail)
  ;  throw(error(domain_error(encoding(did), Path), did_components/2))
  ).

did_build_query_chars(Query, []) :-
  var(Query),
  !.
did_build_query_chars(Query, [QuestionMark | QueryChars]) :-
  did_question_mark_char(QuestionMark),
  with_context(did_components/2, must_be(atom, Query)),
  atom_chars(Query, QueryChars),
  did_validate_query_or_fragment_chars(Query, QueryChars).

did_build_fragment_chars(Fragment, []) :-
  var(Fragment),
  !.
did_build_fragment_chars(Fragment, [Hash | FragmentChars]) :-
  did_hash_char(Hash),
  with_context(did_components/2, must_be(atom, Fragment)),
  atom_chars(Fragment, FragmentChars),
  did_validate_query_or_fragment_chars(Fragment, FragmentChars).

did_validate_method_chars(Method, MethodChars) :-
  ( did_method_chars(MethodChars)
  -> true
  ;  throw(error(domain_error(encoding(did), Method), did_components/2))
  ).

did_validate_method_specific_id_chars(MethodSpecificId, MethodSpecificIdChars) :-
  ( did_method_specific_id_chars(MethodSpecificIdChars)
  -> true
  ;  throw(error(domain_error(encoding(did), MethodSpecificId), did_components/2))
  ).

did_validate_path_chars(Path, PathChars) :-
  ( did_path_chars_raw(PathChars)
  -> true
  ;  throw(error(domain_error(encoding(did), Path), did_components/2))
  ).

did_validate_query_or_fragment_chars(Value, Chars) :-
  ( did_query_or_fragment_chars_raw(Chars)
  -> true
  ;  throw(error(domain_error(encoding(did), Value), did_components/2))
  ).

did_parse_chars(['d', 'i', 'd', ':' | Rest], Method, MethodSpecificId, Path, Query, Fragment) :-
  did_split_required(Rest, ':', MethodChars, AfterMethod),
  MethodChars \= [],
  did_method_chars(MethodChars),
  did_take_until_delimiters(AfterMethod, [47, 63, 35], MethodSpecificIdChars, Suffix),
  MethodSpecificIdChars \= [],
  did_method_specific_id_chars(MethodSpecificIdChars),
  did_parse_suffix(Suffix, Path, Query, Fragment),
  atom_chars(Method, MethodChars),
  atom_chars(MethodSpecificId, MethodSpecificIdChars).

did_parse_suffix([], absent, absent, absent).
did_parse_suffix(['/' | Rest], Path, Query, Fragment) :-
  did_take_until_delimiters(Rest, [63, 35], PathTailChars, AfterPath),
  did_path_chars_raw(PathTailChars),
  atom_chars(PathAtom, ['/' | PathTailChars]),
  Path = present(PathAtom),
  did_parse_query_and_fragment(AfterPath, Query, Fragment).
did_parse_suffix(Rest, absent, Query, Fragment) :-
  did_parse_query_and_fragment(Rest, Query, Fragment).

did_parse_query_and_fragment([], absent, absent).
did_parse_query_and_fragment([QuestionMark | Rest], Query, Fragment) :-
  did_question_mark_char(QuestionMark),
  did_take_until_delimiters(Rest, [35], QueryChars, AfterQuery),
  did_query_or_fragment_chars_raw(QueryChars),
  atom_chars(QueryAtom, QueryChars),
  Query = present(QueryAtom),
  did_parse_fragment(AfterQuery, Fragment).
did_parse_query_and_fragment(Rest, absent, Fragment) :-
  did_parse_fragment(Rest, Fragment).

did_parse_fragment([], absent).
did_parse_fragment([Hash | FragmentChars], Fragment) :-
  did_hash_char(Hash),
  did_query_or_fragment_chars_raw(FragmentChars),
  atom_chars(FragmentAtom, FragmentChars),
  Fragment = present(FragmentAtom).

did_split_required([Sep | Rest], Sep, [], Rest) :-
  !.
did_split_required([Char | Rest], Sep, [Char | Prefix], Suffix) :-
  did_split_required(Rest, Sep, Prefix, Suffix).

did_take_until_delimiters([], _, [], []).
did_take_until_delimiters([Char | Rest], DelimiterCodes, [], [Char | Rest]) :-
  char_code(Char, Code),
  memberchk(Code, DelimiterCodes),
  !.
did_take_until_delimiters([Char | Rest], DelimiterCodes, [Char | Prefix], Suffix) :-
  did_take_until_delimiters(Rest, DelimiterCodes, Prefix, Suffix).

did_concat_chars([], []).
did_concat_chars([Chars | Rest], AllChars) :-
  did_concat_chars(Rest, RestChars),
  did_copy_chars(Chars, AllChars, RestChars).

did_copy_chars([], Tail, Tail).
did_copy_chars([Char | Rest], [Char | Tail], End) :-
  did_copy_chars(Rest, Tail, End).

did_method_chars([Char | Rest]) :-
  did_method_char(Char),
  did_method_chars_rest(Rest).

did_method_chars_rest([]).
did_method_chars_rest([Char | Rest]) :-
  did_method_char(Char),
  did_method_chars_rest(Rest).

did_method_char(Char) :-
  char_code(Char, Code),
  Code >= 0'a,
  Code =< 0'z.
did_method_char(Char) :-
  did_digit_char(Char).

did_method_specific_id_chars(Chars) :-
  did_method_specific_id_segment(Chars, Rest),
  did_method_specific_id_chars_rest(Rest).

did_method_specific_id_chars_rest([]).
did_method_specific_id_chars_rest([':' | Rest]) :-
  did_method_specific_id_segment(Rest, Next),
  did_method_specific_id_chars_rest(Next).

did_method_specific_id_segment(Chars, Rest) :-
  did_method_specific_id_unit(Chars, Next),
  did_method_specific_id_segment_rest(Next, Rest).

did_method_specific_id_segment_rest(Chars, Rest) :-
  did_method_specific_id_unit(Chars, Next),
  did_method_specific_id_segment_rest(Next, Rest).
did_method_specific_id_segment_rest(Rest, Rest) :-
  Rest = []
; Rest = [':' | _].

did_method_specific_id_unit(['%', Hi, Lo | Rest], Rest) :-
  did_hex_char(Hi),
  did_hex_char(Lo).
did_method_specific_id_unit([Char | Rest], Rest) :-
  did_method_specific_id_plain_char(Char).

did_method_specific_id_plain_char(Char) :-
  did_alpha_char(Char).
did_method_specific_id_plain_char(Char) :-
  did_digit_char(Char).
did_method_specific_id_plain_char('.').
did_method_specific_id_plain_char('-').
did_method_specific_id_plain_char('_').

did_path_chars_raw([]).
did_path_chars_raw(Chars) :-
  did_path_unit(Chars, Rest),
  did_path_chars_raw(Rest).

did_path_unit(['%', Hi, Lo | Rest], Rest) :-
  did_hex_char(Hi),
  did_hex_char(Lo).
did_path_unit(['/' | Rest], Rest).
did_path_unit([Char | Rest], Rest) :-
  did_uri_pchar_plain_char(Char).

did_query_or_fragment_chars_raw([]).
did_query_or_fragment_chars_raw(Chars) :-
  did_query_or_fragment_unit(Chars, Rest),
  did_query_or_fragment_chars_raw(Rest).

did_query_or_fragment_unit(['%', Hi, Lo | Rest], Rest) :-
  did_hex_char(Hi),
  did_hex_char(Lo).
did_query_or_fragment_unit(['/' | Rest], Rest).
did_query_or_fragment_unit([QuestionMark | Rest], Rest) :-
  did_question_mark_char(QuestionMark).
did_query_or_fragment_unit([Char | Rest], Rest) :-
  did_uri_pchar_plain_char(Char).

did_uri_pchar_plain_char(Char) :-
  did_uri_unreserved_char(Char).
did_uri_pchar_plain_char(Char) :-
  did_uri_sub_delim_char(Char).
did_uri_pchar_plain_char(':').
did_uri_pchar_plain_char('@').

did_uri_unreserved_char(Char) :-
  did_alpha_char(Char).
did_uri_unreserved_char(Char) :-
  did_digit_char(Char).
did_uri_unreserved_char('-').
did_uri_unreserved_char('.').
did_uri_unreserved_char('_').
did_uri_unreserved_char('~').

did_uri_sub_delim_char('!').
did_uri_sub_delim_char('$').
did_uri_sub_delim_char('&').
did_uri_sub_delim_char('''').
did_uri_sub_delim_char('(').
did_uri_sub_delim_char(')').
did_uri_sub_delim_char('*').
did_uri_sub_delim_char('+').
did_uri_sub_delim_char(',').
did_uri_sub_delim_char(';').
did_uri_sub_delim_char('=').

did_alpha_char(Char) :-
  did_lower_alpha_char(Char).
did_alpha_char(Char) :-
  did_upper_alpha_char(Char).

did_lower_alpha_char(Char) :-
  char_code(Char, Code),
  Code >= 0'a,
  Code =< 0'z.

did_upper_alpha_char(Char) :-
  char_code(Char, Code),
  Code >= 0'A,
  Code =< 0'Z.

did_digit_char(Char) :-
  char_code(Char, Code),
  Code >= 0'0,
  Code =< 0'9.

did_hex_char(Char) :-
  did_digit_char(Char).
did_hex_char(Char) :-
  char_code(Char, Code),
  Code >= 0'a,
  Code =< 0'f.
did_hex_char(Char) :-
  char_code(Char, Code),
  Code >= 0'A,
  Code =< 0'F.

did_question_mark_char(Char) :-
  char_code(Char, 63).

did_hash_char(Char) :-
  char_code(Char, 35).
