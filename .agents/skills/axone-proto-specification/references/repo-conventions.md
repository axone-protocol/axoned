# Repo Conventions

This reference captures the protobuf conventions that are specific to this repository.

## Source of truth

- Schemas live under `proto/`
- Generation is driven by `make proto`
- `make proto` runs formatting, lint, code generation, and proto docs generation
- `proto/buf.yaml` uses Buf breaking mode `FILE`

## Files to inspect before editing

- `proto/buf.yaml`
- `proto/buf.gen.gogo.yml`
- `proto/buf.gen.doc.yml`
- `proto/<module>/docs.yaml`
- one existing query file and one tx file in the same module

## Module-local style matters

The repo is not perfectly uniform between modules.

- `logic` uses `QueryService` and `MsgService`
- `mint` uses `Query` and `Msg`

Do not “fix” this by accident. Follow the convention already used inside the module you touch unless the change is a deliberate module-wide refactor.

## Existing patterns worth preserving

- `option go_package` points to `github.com/axone-protocol/axoned/x/<module>/types`
- Query RPCs commonly expose `google.api.http` GET routes
- Safe query RPCs in `logic` use `option (cosmos.query.v1.module_query_safe) = true`
- Transaction messages use `option (cosmos.msg.v1.signer) = "..."`
- Address fields use `(cosmos_proto.scalar) = "cosmos.AddressString"`
- Complex embedded messages often use `(gogoproto.nullable) = false`
- YAML tags are widely added via `(gogoproto.moretags)`
- Decimal-like values in `mint` use Cosmos scalar annotations plus gogoproto custom types

## Comments are part of the contract

Buf lint in this repo enables comment-related rules. Even though some comment rules are relaxed in `proto/buf.yaml`, the repo clearly expects meaningful comments on services, messages, and many fields. Write them as specification text, not filler.

## Docs flow from proto

- `make doc-proto` renders markdown docs from the proto files
- `proto/<module>/docs.yaml` injects module-level narrative into rendered docs
- If the schema meaning changes, update `docs.yaml` when the module narrative also changes

## Design heuristics specific to this repo

- Prefer domain terms that match AXONE concepts, not generic transport names
- Keep REST paths stable and explicit
- Make cross-module semantics obvious in comments when a field represents a path, content hash, authority, or blockchain-specific scalar
- In `logic`, some API shapes intentionally surface product semantics such as program identity and VFS paths. Preserve those semantics instead of abstracting them away

## Validation

- `make proto`
- Inspect generated `*.pb.go`, `*.pb.gw.go`, and `docs/proto/*.md`
- If proto changes require Go updates, follow with `make lint-go` and `make test-go`
