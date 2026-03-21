---
name: axone-proto-specification
description: Design and evolve protobuf APIs for this repository. Use when creating or changing files under proto/, or when shaping the Go implementation that follows from a proto contract. Focuses on API specification first, with repo-specific conventions plus Cosmos SDK, gRPC, and protobuf compatibility rules.
compatibility: Requires go, make, docker, and local access to the repository. Read reference files in this skill when changing wire-visible schemas.
---

# Axone Proto Specification

Use this skill when changing protobuf schemas or designing module APIs that will be expressed in `proto/`.

This skill is about specification first. In this repo, the `.proto` contract is the source of truth; generated code, gateway routes, docs, and most implementation structure follow from it.

## Read order

1. Read [repo conventions](references/repo-conventions.md).
2. Read [external guidance](references/external-guidance.md) if the change affects fields, service shape, compatibility, or error semantics.
3. Then design the schema before editing code.

## Design rules

1. Start from the domain model and invariants, not from the current keeper or handler code.
2. Treat every field number, field name, service name, and HTTP path as part of a long-lived public contract.
3. Preserve the local convention of the touched module unless you are intentionally migrating the whole module. Do not normalize `logic` and `mint` styles accidentally.
4. Prefer small, explicit RPCs with dedicated request and response messages.
5. Use Cosmos-specific annotations only when they express real semantics such as signer identity, safe queries, or scalar meaning.
6. Design pagination, field presence, and reserved slots deliberately. Do not leave them to later cleanup.
7. Write comments as part of the spec, not as afterthoughts. In this repo, comments are linted and flow into generated docs.

## What to change first

- Edit `proto/<module>/<version>/*.proto`
- Keep module docs aligned through `proto/<module>/docs.yaml` when the user-facing description changes
- Let generation produce Go code and markdown docs
- Then adapt the Go implementation to the contract rather than mutating the contract to match incidental implementation details

## Compatibility rules

- Never reuse a field number.
- When removing a field, reserve both its number and name.
- Do not change wire-visible meaning casually: type, cardinality, field meaning, enum behavior, and path semantics all matter.
- If a field needs presence semantics, decide that explicitly before choosing the representation.
- If an RPC is side-effect free, keep it query-shaped and mark it accordingly.
- If an RPC mutates state or is authority-gated, express signer and authority semantics explicitly in the schema.

## Validation

- Run `make proto`
- If you changed Go implementation too, run `make lint-go` and `make test-go`
- Inspect generated diffs instead of trusting generation blindly

## Output expectation

When using this skill, the agent should be able to explain:

- why the schema shape fits the module domain
- why the chosen field numbers and names are safe to evolve
- why the annotations and HTTP mapping match Cosmos and gRPC semantics
- what generated and implementation files must change downstream
