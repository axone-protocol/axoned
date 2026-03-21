---
name: axone-logic-predicate
description: "Add or change predicates in the AXONE logic module. Use when working on x/logic predicate behavior, Prolog libraries, predicate docs, VFS-backed logic capabilities, or feature scenarios. Follows the current architecture direction: new predicates should be written in Prolog, either pure Prolog or Prolog backed by the path-based logic VFS, not as new native Go predicates."
compatibility: Requires go, make, and docker for validation and doc generation.
---

# Axone Logic Predicate

Use this skill when changing predicate behavior in `x/logic`.

## Architectural direction

New predicates should not be introduced as native Go predicates in `x/logic/predicate` or registered through `x/logic/interpreter/registry.go`.

The target model is:

1. pure Prolog predicates implemented in `x/logic/lib/*.pl`, or
2. Prolog predicates implemented in `x/logic/lib/*.pl` and backed by the path-based logic VFS for host I/O.

For transactional endpoints, use the existing device-file model: path-based VFS access plus a half-duplex request/commit/response interaction pattern.

Treat `x/logic/predicate` and the interpreter registry as existing runtime primitives and legacy surface area. Change them only for maintenance, bug fixes, or truly unavoidable interpreter-level work. If you need to touch them for new functionality, call that out explicitly as an architectural exception.

## Choose the implementation shape

- If the behavior can be expressed by composing existing predicates and libraries, implement it as pure Prolog in `x/logic/lib/*.pl`.
- If the behavior needs host data or external capabilities, expose that capability through the VFS under `x/logic/fs/**`, then wrap it with a Prolog predicate in `x/logic/lib/*.pl`.
- If the behavior is an interactive transactional endpoint, expose it under `/v1/dev/...` and use the device helpers from `x/logic/lib/dev.pl`.
- Do not add a new Go predicate when a Prolog wrapper over `open/4`, `read_term/3`, `write_term/3`, or `dev_call/4` would solve the problem.

## VFS model

The logic VFS is the boundary between logical evaluation and host capabilities.

- `/v1/lib`: embedded Prolog libraries
- `/v1/run`: invocation-scoped runtime resources
- `/v1/var/lib`: persistent host-managed resources
- `/v1/dev`: interactive device-like capabilities

Choose the path family that matches the semantics:

- Snapshot or read-only runtime state: prefer `/v1/run/...`
- Persistent queryable resources: prefer `/v1/var/lib/...`
- Request-response endpoints: prefer `/v1/dev/...`

For `/v1/dev/...` endpoints, keep the protocol aligned with the existing half-duplex model:

1. writes build the request
2. the first read commits the transaction
3. subsequent reads stream the response

## Implementation patterns

### Pure Prolog predicate

- Add or update a library file in `x/logic/lib/*.pl`
- Validate inputs with `with_context/2` and `must_be/2` when appropriate
- Use `setup_call_cleanup/3` around stream access
- Write PlDoc comments because `make doc-predicate` depends on them

### VFS-backed predicate

- Implement the host-facing filesystem in `x/logic/fs/**`
- Mount it through the standard VFS if it is part of the canonical host surface
- Expose an ergonomic Prolog predicate in `x/logic/lib/*.pl`
- Prefer returning logical terms from Prolog wrappers, not leaking protocol details unless the predicate is intentionally low-level

### Transactional device-backed predicate

- Reuse `dev_call/4`, `dev_write_bytes/2`, and `dev_read_bytes/2` from `x/logic/lib/dev.pl`
- Encode protocol errors in-band when designing the device protocol
- Keep transport details inside the VFS device and Prolog wrapper; keep the public predicate ergonomic
- Follow the pattern used by `wasm_query/3` and codec-backed helpers

## Documentation contract

For Prolog predicates, documentation is not optional. The generated predicate docs depend on the source comments.

- Put a `%!` signature line immediately above the predicate definition, for example `%! wasm_query(+Address, +RequestBytes, -ResponseBytes) is det.`
- Continue the doc block with `%` comment lines directly above the predicate head.
- Keep the doc block attached to the predicate. The generator reads the contiguous `%` block and associates it with the next predicate head.
- Document the public predicate, not just helper predicates.
- Include enough description for the generated `## Description` section to stand on its own.
- Mention loading expectations when relevant, but do not duplicate the auto-generated `consult('/v1/lib/...')` section in prose.

## Feature contract

Every predicate change should come with feature coverage in `x/logic/tests/predicate/features`.

- Name the feature file after the predicate, using the generated doc naming convention: `name_arity.feature` for `name/arity`.
- Start with `Feature: name/arity`.
- Add scenarios for the main success path and important failure paths.
- Mark the scenarios that should appear in generated docs with `@great_for_documentation`.
- Treat features as both executable tests and documentation examples. Keep them readable and domain-oriented.
- When the predicate is provided by a library file under `/v1/lib`, include the explicit `consult('/v1/lib/...').` step in the scenario program unless the scenario is specifically testing availability before consult.

## Files to update

- Predicate API: `x/logic/lib/*.pl`
- Host capability surface: `x/logic/fs/**`
- Predicate feature coverage: `x/logic/tests/predicate/features/*.feature`
- Predicate docs generator inputs: PlDoc in `x/logic/lib/*.pl`

Only touch these for exceptions or maintenance:

- `x/logic/predicate/*.go`
- `x/logic/interpreter/registry.go`

## Validation

- Run `make doc-predicate`
- Run `make test-go`
- If you changed Go code under `x/logic/fs/**` or related plumbing, also run `make lint-go`

## Gotchas

- `make doc-predicate` still scans both Go predicates and Prolog predicates. That is a generator detail, not a design recommendation.
- New `.pl` files under `x/logic/lib` are embedded automatically and become available under `/v1/lib/...`.
- Feature files under `x/logic/tests/predicate/features` are both behavioral tests and documentation inputs. Only scenarios tagged `@great_for_documentation` are rendered into generated docs.
- `dev_call/4` is the preferred low-level abstraction for transactional endpoints. Build a higher-level domain predicate on top of it instead of exposing raw device usage directly to callers when possible.
- If a predicate should be auto-loaded as part of the interpreter bootstrap rather than explicitly consulted from `/v1/lib/...`, treat that as a special case and justify it before editing bootstrap files.

## Read first

- `x/logic/fs/std_fs.go`
- `x/logic/lib/dev.pl`
- `x/logic/lib/wasm.pl`
- `x/logic/tests/predicate/features/dev_call_4.feature`
- `scripts/generate_predicates_doc.go`
