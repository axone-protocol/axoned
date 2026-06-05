# Virtual File System

The AXONE logic Virtual File System (VFS) is the capability surface exposed to
Prolog programs. Programs do not read the node filesystem directly. They resolve
versioned `/v1/...` paths mounted by the AXONE blockchain to load AXONE libraries,
inspect the current execution context, read chain-backed state, or call chain
services.

Most programs should use the predicates documented in this section. The VFS paths
matter when you are composing lower-level predicates, auditing what a wrapper
does, or integrating a capability that does not yet have a dedicated predicate.

## `/v1`

The path determines the protocol:

| Path family | Open mode | Stream type | What changes for callers |
| --- | --- | --- | --- |
| `/v1/lib/...` | `read` or `consult/1` | Text | Loads Prolog source provided by the AXONE blockchain. |
| `/v1/run/...` | `read` | Text | Reads immutable data for the current logic execution. |
| `/v1/var/lib/...` | `read` | Text | Reads deterministic chain-backed state exposed to logic. |
| `/v1/dev/...` | `read_write` | Text or binary | Writes a request, then reads the response from a transactional endpoint. |

Read-only paths usually return Prolog terms terminated by `.`. Always close the
stream you open:

```prolog
setup_call_cleanup(
  open('/v1/run/header/height', read, Stream, [type(text)]),
  read_term(Stream, Height, []),
  close(Stream)
).
```

Transactional paths under `/v1/dev/...` follow a half-duplex protocol. The first
read commits the request; after that the stream is response-only and further
writes fail. Prefer `dev_call/4` from `/v1/lib/dev.pl` when you need direct
device access:

```prolog
:- consult('/v1/lib/dev.pl').

call_device(Path, RequestBytes, ResponseBytes) :-
  dev_call(Path, binary, dev_write_bytes(RequestBytes), dev_read_bytes(ResponseBytes)).
```

Device calls are bounded by request and response size limits. Use the wrapper
predicates listed below when they exist; they encode the expected request shape
and error mapping for the target capability.

## `/v1/lib/*.pl`

AXONE embeds reusable Prolog libraries under `/v1/lib`. Load them with
`consult/1`; they are Prolog source files.

| Field | Value |
| --- | --- |
| Paths | `/v1/lib/*.pl` |
| Open mode | `read`, usually through `consult/1` |
| Stream type | Text |
| Response | Prolog source code |

Example:

```prolog
:- consult('/v1/lib/chain.pl').
:- consult('/v1/lib/bank.pl').
```

## `/v1/run/header`

The `/v1/run/header` paths expose the SDK header snapshot for the current logic
execution. Use them when the relation depends on the AXONE block context.

| Field | Value |
| --- | --- |
| Paths | `/v1/run/header/@`, `/v1/run/header/height`, `/v1/run/header/hash`, `/v1/run/header/time`, `/v1/run/header/chain_id`, `/v1/run/header/app_hash` |
| Open mode | `read` |
| Stream type | Text |
| Response | One Prolog term. `@` returns a `header{...}` dict; field paths return the selected field. |
| Recommended predicate | `header_info/1` from `/v1/lib/chain.pl` |

Example:

```prolog
:- consult('/v1/lib/chain.pl').

current_height(Height) :-
  header_info(Header),
  Height = Header.height.
```

## `/v1/run/comet`

The `/v1/run/comet` paths expose CometBFT block data attached to the current
execution, such as validator information, proposer address, evidence, and last
commit data.

| Field | Value |
| --- | --- |
| Paths | `/v1/run/comet/@`, `/v1/run/comet/validators_hash`, `/v1/run/comet/proposer_address`, `/v1/run/comet/evidence`, `/v1/run/comet/last_commit`, `/v1/run/comet/last_commit/round`, `/v1/run/comet/last_commit/votes` |
| Open mode | `read` |
| Stream type | Text |
| Response | One Prolog term. `@` returns a `comet{...}` dict; field paths return the selected field. |
| Recommended predicate | `comet_info/1` from `/v1/lib/chain.pl` |

## `/v1/run/source/files`

`/v1/run/source/files` lists the Prolog source files loaded in the current
interpreter.

| Field | Value |
| --- | --- |
| Path | `/v1/run/source/files` |
| Open mode | `read` |
| Stream type | Text |
| Response | One Prolog list of source file atoms |
| Recommended predicate | `source_file/1` |

## `/v1/var/lib/bank/<address>`

The `/v1/var/lib/bank/<address>` paths expose account balances from AXONE chain
state. `<address>` must be a valid account Bech32 address.

| Field | Value |
| --- | --- |
| Paths | `/v1/var/lib/bank/<address>/balances/@`, `/v1/var/lib/bank/<address>/spendable/@`, `/v1/var/lib/bank/<address>/locked/@` |
| Open mode | `read` |
| Stream type | Text |
| Response | A stream of `Denom-Amount` Prolog terms, one term per coin |
| Recommended predicates | `bank_balances/2`, `bank_spendable_balances/2`, `bank_locked_balances/2` from `/v1/lib/bank.pl` |

Amounts are integers when they fit in `int64`; larger amounts are atoms
preserving the full decimal value.

## `/v1/var/lib/logic/users/<publisher>/programs/<program_id>.pl`

Published logic programs are exposed as Prolog source files. This lets a program
load another program that has been published on AXONE chain state.

| Field | Value |
| --- | --- |
| Path | `/v1/var/lib/logic/users/<publisher>/programs/<program_id>.pl` |
| Open mode | `read`, usually through `consult/1` |
| Stream type | Text |
| Response | Prolog source code |
| Recommended predicate | `consult/1` |

`<publisher>` is the publisher account address. `<program_id>` is the hex
program identifier.

## `/v1/dev/codec/<codec>`

Codec paths transform data through codecs supported by the AXONE blockchain. They
are transactional endpoints: open in `read_write`, write the complete request,
then read the serialized response term.

| Field | Value |
| --- | --- |
| Path | `/v1/dev/codec/<codec>` |
| Open mode | `read_write` |
| Stream type | Text |
| Request | Command plus payload. Shape depends on the codec. |
| Response | One serialized Prolog term, usually `ok(Value)` or `error(Code)` |
| Recommended predicates | `bech32_address/2`, `json_prolog/2`, `json_read/2`, `json_write/2`, and text helpers such as `string_bytes/3` |

### `/v1/dev/codec/bech32`

The Bech32 codec converts between AXONE Bech32 atoms and `Hrp-Bytes` Prolog
pairs.

| Field | Value |
| --- | --- |
| Path | `/v1/dev/codec/bech32` |
| Request | `encode <hrp> <hex_bytes>` or `decode <bech32>` |
| Response | `ok(Bech32)`, `ok(Hrp-Bytes)`, or `error(Code)` |
| Recommended predicate | `bech32_address/2` from `/v1/lib/bech32.pl` |

### `/v1/dev/codec/json`

The JSON codec converts between JSON text and AXONE's canonical Prolog JSON
representation.

| Field | Value |
| --- | --- |
| Path | `/v1/dev/codec/json` |
| Request | `decode` followed by JSON text, or `encode` followed by a canonical Prolog JSON term |
| Response | `ok(Value)` or `error(Code)` |
| Recommended predicates | `json_prolog/2`, `json_read/2`, `json_write/2` from `/v1/lib/json.pl` |

Canonical JSON terms use `json(NameValueList)` for objects, lists for arrays,
atoms for strings, numbers for JSON numbers, and `@(true)`, `@(false)`, and
`@(null)` for booleans and null.

### `/v1/dev/codec/text`

The text codec converts between Prolog textual values and byte lists. It backs
the `string_bytes/3` predicate for encodings such as `text`, `utf8`, `octet`,
and `hex`.

| Field | Value |
| --- | --- |
| Path | `/v1/dev/codec/text` |
| Request | `encode` or `decode`, followed by a Prolog payload term |
| Response | `ok(Value)` or `error(Code)` |
| Recommended predicate | `string_bytes/3` |

## `/v1/dev/crypto/<algorithm>`

Crypto paths expose hashing and signature verification supported by the AXONE
blockchain.

| Field | Value |
| --- | --- |
| Path | `/v1/dev/crypto/<algorithm>` |
| Open mode | `read_write` |
| Stream type | Binary for hashes, text for signature verification |
| Request | Hash devices receive raw bytes. Signature devices receive `verify <pubkey_hex> <data_hex> <signature_hex>`. |
| Response | Hash devices return raw digest bytes. Signature devices return a Prolog term such as `ok(true)`, `ok(false)`, or `error(Code)`. |
| Recommended predicates | `crypto_data_hash/3`, `eddsa_verify/4`, `ecdsa_verify/4` from `/v1/lib/crypto.pl` |

### `/v1/dev/crypto/md5`

| Field | Value |
| --- | --- |
| Path | `/v1/dev/crypto/md5` |
| Operation | MD5 digest |
| Stream type | Binary |
| Request | Raw bytes |
| Response | Raw digest bytes |
| Recommended predicate | `crypto_data_hash/3` with `algorithm(md5)` |

### `/v1/dev/crypto/sha256`

| Field | Value |
| --- | --- |
| Path | `/v1/dev/crypto/sha256` |
| Operation | SHA-256 digest |
| Stream type | Binary |
| Request | Raw bytes |
| Response | Raw digest bytes |
| Recommended predicate | `crypto_data_hash/3` with `algorithm(sha256)` |

### `/v1/dev/crypto/sha512`

| Field | Value |
| --- | --- |
| Path | `/v1/dev/crypto/sha512` |
| Operation | SHA-512 digest |
| Stream type | Binary |
| Request | Raw bytes |
| Response | Raw digest bytes |
| Recommended predicate | `crypto_data_hash/3` with `algorithm(sha512)` |

### `/v1/dev/crypto/ed25519`

| Field | Value |
| --- | --- |
| Path | `/v1/dev/crypto/ed25519` |
| Operation | Ed25519 signature verification |
| Stream type | Text |
| Request | `verify <pubkey_hex> <data_hex> <signature_hex>` |
| Response | `ok(true)`, `ok(false)`, or `error(Code)` |
| Recommended predicate | `eddsa_verify/4` with `type(ed25519)` |

### `/v1/dev/crypto/secp256r1`

| Field | Value |
| --- | --- |
| Path | `/v1/dev/crypto/secp256r1` |
| Operation | ECDSA secp256r1 signature verification |
| Stream type | Text |
| Request | `verify <pubkey_hex> <data_hex> <signature_hex>` |
| Response | `ok(true)`, `ok(false)`, or `error(Code)` |
| Recommended predicate | `ecdsa_verify/4` with `type(secp256r1)` |

### `/v1/dev/crypto/secp256k1`

| Field | Value |
| --- | --- |
| Path | `/v1/dev/crypto/secp256k1` |
| Operation | ECDSA secp256k1 signature verification |
| Stream type | Text |
| Request | `verify <pubkey_hex> <data_hex> <signature_hex>` |
| Response | `ok(true)`, `ok(false)`, or `error(Code)` |
| Recommended predicate | `ecdsa_verify/4` with `type(secp256k1)` |

## `/v1/dev/wasm/<contract_address>/query`

The CosmWasm query path executes a smart query against a contract address visible
from the AXONE blockchain.

| Field | Value |
| --- | --- |
| Path | `/v1/dev/wasm/<contract_address>/query` |
| Open mode | `read_write` |
| Stream type | Binary |
| Request | Raw query bytes, typically UTF-8 JSON |
| Response | Raw contract response bytes |
| Recommended predicate | `wasm_query/3` from `/v1/lib/wasm.pl` |

This device is query-only. It does not execute transactions or mutate contract
state.
