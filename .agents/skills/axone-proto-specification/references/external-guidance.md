# External Guidance

This reference distills the external rules that matter most for this repo's proto work.

## Protocol Buffers

Primary sources:

- [Proto best practices](https://protobuf.dev/best-practices/dos-donts/)
- [Proto language guide](https://protobuf.dev/programming-guides/proto2/)
- [Proto3 specification](https://protobuf.dev/reference/protobuf/proto3-spec/)
- [Field presence](https://protobuf.dev/programming-guides/field_presence/)

Rules to apply:

- Never reuse field numbers
- When deleting a field, reserve its number and name
- Put frequently populated fields in the low field-number range when practical
- Make field presence an explicit design decision, especially for scalar values where “unset” and “default value” differ semantically
- Be conservative with enums: the first value must be zero, and evolution behavior matters for older clients

## Buf

Primary sources:

- [Breaking change detection](https://buf.build/docs/breaking/)
- [Breaking rules and categories](https://buf.build/docs/breaking/rules/)

Rules to apply:

- Treat Buf breaking checks as design feedback, not as a post-hoc formatting step
- Default to compatibility-preserving changes unless the breaking change is intentional and called out
- Remember that this repo uses `FILE` breaking mode, which is stricter than pure wire compatibility in some cases

## Cosmos SDK

Primary source:

- [Cosmos SDK encoding and protobuf guidance](https://docs.cosmos.network/v0.50/learn/advanced/encoding)

Rules to apply:

- Use Cosmos annotations when they carry semantic meaning, not mechanically
- For signer identity, use `cosmos.msg.v1.signer`
- For query safety, use `cosmos.query.v1.module_query_safe` when the RPC is truly side-effect free
- For address and decimal-like values, follow existing Cosmos scalar and customtype patterns from the touched module
- If an API involves interfaces or `Any`, follow Cosmos interface annotation guidance rather than inventing local conventions

## gRPC

Primary source:

- [gRPC status codes](https://grpc.io/docs/guides/status-codes/)

Rules to apply:

- Keep RPC boundaries explicit and domain-specific
- Use dedicated request and response messages
- Design validation and state errors so they can map cleanly to gRPC status semantics
- Distinguish invalid input from state-dependent failures when designing error behavior

## How to use these rules here

- Start from repo conventions first
- Use external rules to reject unsafe schema changes
- If repo-local precedent conflicts with generic best practice, preserve wire compatibility and module consistency first, then document the tradeoff
