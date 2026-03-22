[//]: # (This file is auto-generated. Please do not modify it yourself.)

# Logic Module

<a name="top"></a>

## Function

The `logic` module lets users execute Prolog queries against the current blockchain state and against reusable Prolog
sources, without deploying a smart contract and without mutating chain state.

From a user perspective, the module combines four things:

- the [Axone Prolog VM](https://github.com/axone-protocol/prolog), adapted for deterministic blockchain execution;
- a versioned virtual file system (VFS) exposed under `/v1`;
- a set of host-provided Prolog libraries mounted in that VFS;
- an on-chain registry for immutable user programs that can later be loaded as libraries.

This page focuses on how to use the module. Predicate-by-predicate documentation is maintained elsewhere.

## What You Can Do With This Module

In practice, the module is useful when you want to:

- evaluate a Prolog query against an inline program;
- query chain data through Prolog instead of stitching multiple RPC calls together client-side;
- reuse built-in Axone libraries such as chain helpers, bank helpers, codec helpers, and WASM query helpers;
- publish a Prolog source on-chain once and load it later from other queries with `consult/1`;
- compose reusable rule libraries that remain immutable and content-addressed.

The module is intentionally query-oriented. It is designed to help users express logic, data traversal, and rule-based
reasoning in a deterministic environment.

## User Model

The module revolves around two main user actions:

1. `Ask`: execute a query.
2. `StoreProgram`: publish a reusable Prolog source on-chain.

A typical flow is:

1. write a small Prolog program or library;
2. optionally publish it with `StoreProgram`;
3. run `Ask` with an inline program and a query;
4. load built-in libraries or previously published programs with `consult/1`.

## Programs, Queries, and Execution

A program is the Prolog source compiled before the query runs. It may contain:

- facts;
- rules;
- directives such as `:- consult('/v1/lib/chain.pl').`;
- helper predicates that structure the final query.

The query is the goal executed after the program has been loaded.

The `Ask` endpoint accepts:

```text
{
  string program
  string query
  uint64 limit
}
```

Notes:

- `program` is optional;
- `query` is the goal to solve;
- `limit` is the requested number of solutions, capped by the module parameters.

Example:

```text
{
  program: "parent(john, mary).",
  query: "parent(X, mary).",
  limit: 1
}
```

The module first compiles `program`, then evaluates `query` in the resulting Prolog environment.

## Query Semantics

`Ask` is a blockchain query, not a transaction:

- it has no side effect on chain state;
- it does not publish code;
- it does not write to the blockchain;
- it is still bounded by metering, gas, and module limits.

This matters because the module can expose rich capabilities through the VFS, but those capabilities are still
constrained to deterministic, query-safe operations. For example, smart-contract interaction is exposed as query-only
devices, not as transaction submission.

## Query Response

`Ask` returns:

- `height`: the block height used for evaluation;
- `gas_used`: the SDK gas consumed by the evaluation;
- `answer`: the logical result set;
- `user_output`: text written by the program to the current output stream.

The `answer` field contains:

- `variables`: the variables found in the query;
- `results`: the returned solutions;
- `has_more`: whether more solutions existed beyond the returned slice.

Each result may contain:

- `substitutions`: bindings such as `X = john`;
- `error`: an error message for that solution when relevant.

Substitution values are returned as Prolog terms encoded as strings. For example, atoms, lists, compounds, pairs, and
dicts are returned in Prolog syntax rather than being remapped into ad hoc JSON structures.

## User Output

The default Prolog output stream is captured by the module and returned in `user_output`.

This is useful for:

- debugging a query;
- emitting an explanatory text rendering;
- producing a human-readable report alongside logical substitutions.

The output is bounded by `max_user_output_size`. When the limit is reached, the oldest bytes are discarded and only
the last bytes are kept.

## Virtual File System (VFS)

A key concept of the module is the VFS exposed to the Prolog VM through `open/4` and `consult/1`.

Prolog code does not access the node's real file system. It only sees the capabilities explicitly mounted by the
host under a canonical namespace rooted at `/v1`.

This makes the module easier to understand:

- if a resource is visible through the VFS, Prolog can use it;
- if it is not mounted in the VFS, it does not exist from the program's perspective.

The namespace follows a stable layout:

- `/v1/lib`: immutable host-provided Prolog libraries;
- `/v1/run`: invocation-scoped runtime snapshots;
- `/v1/var/lib`: persistent chain-backed data exposed as files;
- `/v1/dev`: transactional device-like capabilities.

In several mounts, the special path segment `@` denotes the canonical aggregate representation of a resource.

## `/v1/lib`: Built-In Prolog Libraries

The module ships with reusable Prolog libraries that can be loaded explicitly with `consult/1`.

Available libraries include:

- `/v1/lib/apply.pl`: higher-order helpers such as `maplist/2..7` and `foldl/4..7`;
- `/v1/lib/bank.pl`: helpers for account balances;
- `/v1/lib/bech32.pl`: Bech32 encode/decode helpers;
- `/v1/lib/chain.pl`: helpers exposing header and Comet block information;
- `/v1/lib/crypto.pl`: small crypto-adjacent helpers such as hexadecimal conversions;
- `/v1/lib/dev.pl`: low-level helpers for transactional VFS devices;
- `/v1/lib/error.pl`: error and type-checking helpers;
- `/v1/lib/lists.pl`: list processing helpers;
- `/v1/lib/type.pl`: type predicates used by the other libraries;
- `/v1/lib/wasm.pl`: helpers for CosmWasm smart queries.

Example:

```prolog
:- consult('/v1/lib/chain.pl').

current_chain(ChainID) :-
  header_info(Header),
  ChainID = Header.chain_id.
```

These libraries are regular Prolog sources. They are not magical special forms: users load them explicitly and can
combine them with their own predicates.

## `/v1/run`: Runtime Snapshots

`/v1/run` exposes information tied to the current query execution context.

The standard mounts are:

- `/v1/run/header`
- `/v1/run/comet`

Examples of readable paths include:

- `/v1/run/header/@`
- `/v1/run/header/height`
- `/v1/run/header/hash`
- `/v1/run/header/time`
- `/v1/run/header/chain_id`
- `/v1/run/header/app_hash`
- `/v1/run/comet/@`
- `/v1/run/comet/validators_hash`
- `/v1/run/comet/proposer_address`
- `/v1/run/comet/evidence`
- `/v1/run/comet/last_commit`

These files are read-only snapshots. They are intended for logical inspection, not mutation.

Example:

```prolog
:- consult('/v1/lib/chain.pl').

block_height(Height) :-
  header_info(Header),
  Height = Header.height.
```

## `/v1/var/lib`: Persistent Chain-Backed Data

`/v1/var/lib` exposes data backed by blockchain state, still through file-like paths.

The standard mounts currently include:

- `/v1/var/lib/bank`
- `/v1/var/lib/logic/users`

### Bank Data

Bank balances are exposed under paths such as:

- `/v1/var/lib/bank/<address>/balances/@`
- `/v1/var/lib/bank/<address>/spendable/@`
- `/v1/var/lib/bank/<address>/locked/@`

The `bank.pl` library wraps those paths into higher-level predicates:

- `bank_balances/2`
- `bank_spendable_balances/2`
- `bank_locked_balances/2`

Example:

```prolog
:- consult('/v1/lib/bank.pl').

rich_account(Address) :-
  bank_spendable_balances(Address, Balances),
  member(uaxone-Amount, Balances),
  Amount > 1000000.
```

### Published User Programs

User-published programs are exposed under:

```text
/v1/var/lib/logic/users/<publisher>/programs/<program_id>.pl
```

This is the bridge between on-chain publication and program reuse inside `Ask`.

A stored source can therefore be loaded like any other Prolog file:

```prolog
:- consult('/v1/var/lib/logic/users/axone1.../programs/<program_id>.pl').
```

## `/v1/dev`: Transactional Devices

`/v1/dev` exposes interactive capabilities through file-like request/response devices.

This is where the Axone Prolog VM extension for `read_write` streams becomes important. In Axone, `open/4` can open a
resource in `read_write` mode, which allows half-duplex transactional I/O:

1. write a request;
2. trigger the device on first read;
3. read back the response.

This pattern is used for deterministic host capabilities that are easier to model as devices than as plain predicates.

Standard device mounts include:

- `/v1/dev/codec/<codec>`
- `/v1/dev/wasm/<contract_address>/query`

In most cases, users should prefer the provided helper libraries rather than interact with devices directly.

### Codec Devices

Codec devices are mounted under `/v1/dev/codec`. For example, the Bech32 helper library delegates encoding and
decoding to the `bech32` codec device.

Users typically do not open those devices themselves. Instead they use predicates such as:

- `bech32_address/2`
- `hex_bytes/2`

### WASM Query Devices

CosmWasm smart queries are exposed through:

```text
/v1/dev/wasm/<contract_address>/query
```

The recommended user-facing helper is `wasm_query/3` from `/v1/lib/wasm.pl`.

Example:

```prolog
:- consult('/v1/lib/wasm.pl').

contract_query(Address, RequestBytes, ResponseBytes) :-
  wasm_query(Address, RequestBytes, ResponseBytes).
```

The device is query-only. It lets users ask a contract for data, not execute transactions.

## `consult/1` and Reuse

`consult/1` is the main mechanism for composing logic from multiple sources.

Users can consult:

- built-in libraries from `/v1/lib`;
- chain-backed published programs from `/v1/var/lib/logic/users/...`;
- lists of files, just as in regular Prolog.

This is the idiomatic way to build reusable logic on top of the module. In other words, `StoreProgram` does not create
a new query endpoint by itself; it publishes a source file that later queries can load with `consult/1`.

## Publishing Programs With `StoreProgram`

`StoreProgram` is the transaction endpoint of the module.

It validates a Prolog source and stores it as an immutable program artifact.

The important properties are:

- the identifier is content-addressed: `program_id = sha256(source)`;
- the same source always yields the same `program_id`;
- the artifact is immutable;
- publication is recorded per publisher address;
- the endpoint is idempotent for repeated publication of the same source.

This means two different publishers may publish the same source:

- the artifact is shared by content;
- each publisher still gets their own publication path in `/v1/var/lib/logic/users/<publisher>/programs/...`.

The response contains:

- `program_id`: the SHA-256 digest of the source, encoded as lowercase hexadecimal.

## Program Lifecycle

From a user perspective, a stored program behaves like an immutable published library:

1. write the source;
2. publish it with `StoreProgram`;
3. keep the returned `program_id`;
4. load it later through its VFS path;
5. compose it with inline logic inside `Ask`.

Example usage pattern:

```prolog
:- consult('/v1/lib/bank.pl').
:- consult('/v1/var/lib/logic/users/axone1.../programs/<program_id>.pl').

eligible(Address) :-
  shared_rule(Address),
  bank_spendable_balances(Address, Balances),
  member(uaxone-Amount, Balances),
  Amount >= 500000.
```

If the source changes, the hash changes too. Updating a library therefore means publishing a new source, which yields a
new immutable `program_id`.

## Discovering Published Programs

The module exposes query endpoints to inspect published artifacts:

- `Program`: fetch metadata for a `program_id`;
- `ProgramSource`: fetch the original source for a `program_id`;
- `Programs`: list stored programs;
- `ProgramsByPublisher`: list the programs published by a given address.

The metadata includes:

- `program_id`;
- `created_at`;
- `source_size`;

Publisher-scoped results also include publication metadata such as `published_at`.

## Execution Limits

The module is intentionally constrained. The main parameters users should understand are:

- `max_size`: maximum accepted size, in bytes, for the user-supplied program and query in `Ask`; it also bounds source
  size when publishing a new program;
- `max_result_count`: maximum number of solutions that can be requested;
- `max_user_output_size`: maximum number of bytes retained in `user_output`;
- `max_variables`: maximum number of variables the interpreter may allocate.

Those limits exist because Prolog evaluation may involve recursion, unification, backtracking, file I/O through the
VFS, and device interactions. Users should expect the module to reject overly large or overly expensive evaluations.

## Gas and Metering

Even though `Ask` is a query, execution is still metered.

The gas policy translates VM activity into SDK gas with separate coefficients for:

- compute;
- memory;
- unification;
- I/O.

In practice, gas usage grows with:

- the size of the supplied program and query;
- the amount of backtracking required;
- the number and size of terms being unified or copied;
- VFS and device I/O.

Users should therefore:

- keep programs focused;
- bind variables early when possible;
- request only the number of solutions they need;
- avoid broad search spaces when a more selective query can be written.

## Determinism and Safety

The module is built on top of the Axone Prolog VM, which is adapted for blockchain use. From a user standpoint, the
important consequences are:

- execution is deterministic;
- resource consumption is metered;
- variable allocation can be capped;
- the visible environment is restricted to the mounted VFS;
- there is no arbitrary operating-system file access;
- there is no ability to submit transactions or mutate chain state from Prolog.

This is what makes the module suitable for chain queries while still supporting rich logic, reusable libraries, and
controlled access to chain capabilities.

## When To Use Inline Programs vs Stored Programs

Use an inline program when:

- the logic is short-lived;
- the query is one-off;
- you do not need reuse by other users or clients.

Use `StoreProgram` when:

- the source should be reused across many queries;
- the source should be referenced as a library through `consult/1`;
- you want an immutable, content-addressed artifact;
- you want other users or applications to rely on the exact same source.

## Summary

The `logic` module should be understood as a user-facing Prolog query engine for blockchain data, backed by:

- a deterministic Axone Prolog VM;
- a capability-oriented VFS;
- built-in Prolog libraries;
- query-only devices for advanced integrations;
- immutable, on-chain published user programs that can be reused as libraries.

If you want to understand how to use the module, the central ideas are simple:

- write Prolog;
- load libraries with `consult/1`;
- inspect chain capabilities through `/v1`;
- publish reusable sources with `StoreProgram`;
- execute logic with `Ask`.

## API

<a name="logic/v1beta3/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### logic/v1beta3/params.proto

<a name="logic.v1beta3.GasPolicy"></a>

#### GasPolicy

GasPolicy defines the coefficients used to translate VM metering units into SDK gas.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `compute_coeff` | [uint64](#uint64) |  | compute_coeff applies to Instruction, ArithNode, and CompareStep VM meter kinds. If set to 0, the value considered is 1. |
| `memory_coeff` | [uint64](#uint64) |  | memory_coeff applies to CopyNode and ListCell VM meter kinds. If set to 0, the value considered is 1. |
| `unify_coeff` | [uint64](#uint64) |  | unify_coeff applies to UnifyStep VM meter kind. If set to 0, the value considered is 1. |
| `io_coeff` | [uint64](#uint64) |  | io_coeff applies to the total size in bytes of the user-supplied program and query sources, and to codec device I/O buffered by the VFS. If set to 0, the value considered is 1. |

<a name="logic.v1beta3.Limits"></a>

#### Limits

Limits defines the limits of the logic module.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_size` | [uint64](#uint64) |  | max_size specifies the maximum total size, in bytes, accepted for the user-supplied program and query sources. A value of 0 means that there is no limit on the total source size. |
| `max_result_count` | [uint64](#uint64) |  | max_result_count specifies the maximum number of results that can be requested for a query. A value of 0 means that there is no limit on the number of results. |
| `max_user_output_size` | [uint64](#uint64) |  | max_user_output_size specifies the maximum number of bytes to keep in the user output. If the user output exceeds this size, the interpreter will overwrite the oldest bytes with the new ones to keep the size constant. A value of 0 means the user output is disabled. |
| `max_variables` | [uint64](#uint64) |  | max_variables specifies the maximum number of variables that can be create by the interpreter. A value of 0 means that there is no limit on the number of variables. |

<a name="logic.v1beta3.Params"></a>

#### Params

Params defines all the configuration parameters of the "logic" module.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `limits` | [Limits](#logic.v1beta3.Limits) |  | Limits defines the limits of the logic module. The limits are used to prevent the interpreter from running for too long. If the interpreter runs for too long, the execution will be aborted. |
| `gas_policy` | [GasPolicy](#logic.v1beta3.GasPolicy) |  | GasPolicy defines the coefficients used to translate VM metering units into SDK gas. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta3/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### logic/v1beta3/types.proto

<a name="logic.v1beta3.Answer"></a>

#### Answer

Answer represents the answer to a logic query.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `has_more` | [bool](#bool) |  | has_more specifies if there are more solutions than the ones returned. |
| `variables` | [string](#string) | repeated | variables represent all the variables in the query. |
| `results` | [Result](#logic.v1beta3.Result) | repeated | results represent all the results of the query. |

<a name="logic.v1beta3.ProgramMetadata"></a>

#### ProgramMetadata

ProgramMetadata represents the metadata of a stored program.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program_id` | [string](#string) |  | program_id is the SHA-256 hash of the program source encoded as lowercase hexadecimal. |
| `created_at` | [int64](#int64) |  | created_at is the block timestamp (Unix seconds) of artifact creation. |
| `source_size` | [uint64](#uint64) |  | source_size is the source size in bytes. |

<a name="logic.v1beta3.ProgramPublication"></a>

#### ProgramPublication

ProgramPublication represents the publication metadata of a program by a publisher.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `published_at` | [int64](#int64) |  | published_at is the block timestamp (Unix seconds) of publication for this publisher. |

<a name="logic.v1beta3.Result"></a>

#### Result

Result represents the result of a query.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `error` | [string](#string) |  | error specifies the error message if the query caused an error. |
| `substitutions` | [Substitution](#logic.v1beta3.Substitution) | repeated | substitutions represent all the substitutions made to the variables in the query to obtain the answer. |

<a name="logic.v1beta3.StoredProgram"></a>

#### StoredProgram

StoredProgram represents a program source with its storage metadata.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  | source is the original Prolog source. |
| `created_at` | [int64](#int64) |  | created_at is the block timestamp (Unix seconds) of artifact creation. |
| `source_size` | [uint64](#uint64) |  | source_size is the source size in bytes. |

<a name="logic.v1beta3.Substitution"></a>

#### Substitution

Substitution represents a substitution made to the variables in the query to obtain the answer.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `variable` | [string](#string) |  | variable is the name of the variable. |
| `expression` | [string](#string) |  | expression is the value substituted for the variable, represented directly as a Prolog term (e.g., atom, number, compound). |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta3/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### logic/v1beta3/genesis.proto

<a name="logic.v1beta3.GenesisProgramPublication"></a>

#### GenesisProgramPublication

GenesisProgramPublication associates a publisher and program_id with publication metadata.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher` | [string](#string) |  | publisher is the bech32 account address that published the program. |
| `program_id` | [string](#string) |  | program_id is the SHA-256 hash of the published program source encoded as lowercase hexadecimal. |
| `publication` | [ProgramPublication](#logic.v1beta3.ProgramPublication) |  | publication is the publication metadata for this publisher/program pair. |

<a name="logic.v1beta3.GenesisState"></a>

#### GenesisState

GenesisState defines the logic module's genesis state.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta3.Params) |  | The state parameters for the logic module. |
| `stored_programs` | [GenesisStoredProgram](#logic.v1beta3.GenesisStoredProgram) | repeated | stored_programs are the canonical immutable program artifacts keyed by program_id. |
| `program_publications` | [GenesisProgramPublication](#logic.v1beta3.GenesisProgramPublication) | repeated | program_publications are the user-scoped immutable publication views pointing to stored artifacts. |

<a name="logic.v1beta3.GenesisStoredProgram"></a>

#### GenesisStoredProgram

GenesisStoredProgram associates a program_id with its canonical stored artifact.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program_id` | [string](#string) |  | program_id is the SHA-256 hash of the program source encoded as lowercase hexadecimal. |
| `program` | [StoredProgram](#logic.v1beta3.StoredProgram) |  | program is the canonical immutable stored program artifact. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

 [//]: # (end services)

<a name="logic/v1beta3/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### logic/v1beta3/query.proto

<a name="logic.v1beta3.PublishedProgram"></a>

#### PublishedProgram

PublishedProgram represents a publisher-scoped program view.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [ProgramMetadata](#logic.v1beta3.ProgramMetadata) |  | program is the metadata of the stored program. |
| `publication` | [ProgramPublication](#logic.v1beta3.ProgramPublication) |  | publication is the publication metadata for this publisher/program pair. |

<a name="logic.v1beta3.QueryAskRequest"></a>

#### QueryAskRequest

QueryAskRequest is request type for the QueryService/Ask RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [string](#string) |  | program is the logic program to be queried. |
| `query` | [string](#string) |  | query is the query string to be executed. |
| `limit` | [uint64](#uint64) |  | limit specifies the maximum number of solutions to be returned. This field is governed by max_result_count, which defines the upper limit of results that may be requested per query. If this field is not explicitly set, a default value of 1 is applied. |

<a name="logic.v1beta3.QueryAskResponse"></a>

#### QueryAskResponse

QueryAskResponse is response type for the QueryService/Ask RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [uint64](#uint64) |  | height is the block height at which the query was executed. |
| `gas_used` | [uint64](#uint64) |  | gas_used is the amount of gas used to execute the query. |
| `answer` | [Answer](#logic.v1beta3.Answer) |  | answer is the answer to the query. |
| `user_output` | [string](#string) |  | user_output is the output of the query execution, if any. the length of the output is limited by the max_user_output_size parameter. |

<a name="logic.v1beta3.QueryParamsRequest"></a>

#### QueryParamsRequest

QueryParamsRequest is request type for the QueryService/Params RPC method.

<a name="logic.v1beta3.QueryParamsResponse"></a>

#### QueryParamsResponse

QueryParamsResponse is response type for the QueryService/Params RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#logic.v1beta3.Params) |  | params holds all the parameters of this module. |

<a name="logic.v1beta3.QueryProgramRequest"></a>

#### QueryProgramRequest

QueryProgramRequest is request type for the QueryService/Program RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program_id` | [string](#string) |  | program_id is the immutable identifier of the stored program. |

<a name="logic.v1beta3.QueryProgramResponse"></a>

#### QueryProgramResponse

QueryProgramResponse is response type for the QueryService/Program RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program` | [ProgramMetadata](#logic.v1beta3.ProgramMetadata) |  | program is the metadata of the stored program. |

<a name="logic.v1beta3.QueryProgramSourceRequest"></a>

#### QueryProgramSourceRequest

QueryProgramSourceRequest is request type for the QueryService/ProgramSource RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program_id` | [string](#string) |  | program_id is the immutable identifier of the stored program. |

<a name="logic.v1beta3.QueryProgramSourceResponse"></a>

#### QueryProgramSourceResponse

QueryProgramSourceResponse is response type for the QueryService/ProgramSource RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source` | [string](#string) |  | source is the original Prolog source of the stored program. |

<a name="logic.v1beta3.QueryProgramsByPublisherRequest"></a>

#### QueryProgramsByPublisherRequest

QueryProgramsByPublisherRequest is request type for the QueryService/ProgramsByPublisher RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher` | [string](#string) |  | publisher is the bech32 account address that published the programs. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |

<a name="logic.v1beta3.QueryProgramsByPublisherResponse"></a>

#### QueryProgramsByPublisherResponse

QueryProgramsByPublisherResponse is response type for the QueryService/ProgramsByPublisher RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `programs` | [PublishedProgram](#logic.v1beta3.PublishedProgram) | repeated | programs is the list of programs published by the requested publisher. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |

<a name="logic.v1beta3.QueryProgramsRequest"></a>

#### QueryProgramsRequest

QueryProgramsRequest is request type for the QueryService/Programs RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |

<a name="logic.v1beta3.QueryProgramsResponse"></a>

#### QueryProgramsResponse

QueryProgramsResponse is response type for the QueryService/Programs RPC method.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `programs` | [ProgramMetadata](#logic.v1beta3.ProgramMetadata) | repeated | programs is the metadata list of stored programs. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="logic.v1beta3.QueryService"></a>

#### QueryService

QueryService defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#logic.v1beta3.QueryParamsRequest) | [QueryParamsResponse](#logic.v1beta3.QueryParamsResponse) | Params queries all parameters for the logic module. | GET | `/axone-protocol/axoned/logic/params` |
| `Ask` | [QueryAskRequest](#logic.v1beta3.QueryAskRequest) | [QueryAskResponse](#logic.v1beta3.QueryAskResponse) | Ask executes a logic query and returns the solutions found. Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee is charged for this, but the execution is constrained by the current limits configured in the module. | GET | `/axone-protocol/axoned/logic/ask` |
| `Program` | [QueryProgramRequest](#logic.v1beta3.QueryProgramRequest) | [QueryProgramResponse](#logic.v1beta3.QueryProgramResponse) | Program queries the metadata of a stored program by its immutable identifier. | GET | `/axone-protocol/axoned/logic/programs/{program_id}` |
| `ProgramSource` | [QueryProgramSourceRequest](#logic.v1beta3.QueryProgramSourceRequest) | [QueryProgramSourceResponse](#logic.v1beta3.QueryProgramSourceResponse) | ProgramSource queries the source of a stored program by its immutable identifier. | GET | `/axone-protocol/axoned/logic/programs/{program_id}/source` |
| `Programs` | [QueryProgramsRequest](#logic.v1beta3.QueryProgramsRequest) | [QueryProgramsResponse](#logic.v1beta3.QueryProgramsResponse) | Programs lists stored programs. | GET | `/axone-protocol/axoned/logic/programs` |
| `ProgramsByPublisher` | [QueryProgramsByPublisherRequest](#logic.v1beta3.QueryProgramsByPublisherRequest) | [QueryProgramsByPublisherResponse](#logic.v1beta3.QueryProgramsByPublisherResponse) | ProgramsByPublisher lists stored programs published by a given publisher. | GET | `/axone-protocol/axoned/logic/publishers/{publisher}/programs` |

 [//]: # (end services)

<a name="logic/v1beta3/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

### logic/v1beta3/tx.proto

<a name="logic.v1beta3.MsgStoreProgram"></a>

#### MsgStoreProgram

MsgStoreProgram defines a Msg for storing a Prolog program source as a user library.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher` | [string](#string) |  | publisher is the bech32 account address publishing the program artifact. After publication, this exact address is used as the `<identity>` path segment in the logic module virtual file system path `/v1/var/lib/logic/users/<identity>/programs/<program_id>.pl`. This is the path that Prolog code can load through `consult/1`. |
| `source` | [string](#string) |  | source is the Prolog program source to parse and store. |

<a name="logic.v1beta3.MsgStoreProgramResponse"></a>

#### MsgStoreProgramResponse

MsgStoreProgramResponse defines the response for executing a MsgStoreProgram.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `program_id` | [string](#string) |  | program_id is the SHA-256 hash of the program source (lowercase hexadecimal). After publication, this exact identifier is used as the `<program_id>` path segment in the logic module virtual file system path `/v1/var/lib/logic/users/<identity>/programs/<program_id>.pl`. This is the path that Prolog code can load through `consult/1`. |

<a name="logic.v1beta3.MsgUpdateParams"></a>

#### MsgUpdateParams

MsgUpdateParams defines a Msg for updating the x/logic module parameters.

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#logic.v1beta3.Params) |  | params defines the x/logic parameters to update. NOTE: All parameters must be supplied. |

<a name="logic.v1beta3.MsgUpdateParamsResponse"></a>

#### MsgUpdateParamsResponse

MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.

 [//]: # (end messages)

 [//]: # (end enums)

 [//]: # (end HasExtensions)

<a name="logic.v1beta3.MsgService"></a>

#### MsgService

MsgService defines the transaction service for the logic module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#logic.v1beta3.MsgUpdateParams) | [MsgUpdateParamsResponse](#logic.v1beta3.MsgUpdateParamsResponse) | UpdateParams defined a governance operation for updating the x/logic module parameters. The authority is hard-coded to the Cosmos SDK x/gov module account |   |   |
| `StoreProgram` | [MsgStoreProgram](#logic.v1beta3.MsgStoreProgram) | [MsgStoreProgramResponse](#logic.v1beta3.MsgStoreProgramResponse) | StoreProgram validates a Prolog user library source and stores its canonical artifact if needed. Artifact identity is content-addressed: `program_id = sha256(source)`. The endpoint is idempotent for the same publisher + same source, and also when different publishers submit the same source. After a successful call, the published program is exposed through the logic module virtual file system at the immutable path `/v1/var/lib/logic/users/<identity>/programs/<program_id>.pl`. This path is intended to be loaded from Prolog, for example with `consult/1`. |   |   |

 [//]: # (end services)

### Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
