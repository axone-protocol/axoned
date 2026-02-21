% chain.pl
% Chain-related predicates for querying chain information such as block headers and Comet block info.

% comet_info(?CometInfo) is det.
%
% Unifies CometInfo with the current Comet block info dict exposed by the VFS.
comet_info(CometInfo) :-
  open('/v1/sys/comet/@', read, Stream, [type(text)]),
  read_term(Stream, CometInfo, []),
  close(Stream).
