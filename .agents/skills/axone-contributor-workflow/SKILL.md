---
name: axone-contributor-workflow
description: Contribute safely to the AXONE axoned repository. Use when changing Go code, proto files, generated docs, predicate logic, tests, CI workflows, release plumbing, or local chain behavior in this repo. Applies the repo's validation steps, generation commands, and contributor-specific gotchas before finalizing changes.
compatibility: Requires git, go, make, and docker. Node is only needed when validating generated markdown for MDX compatibility.
---

# Axone Contributor Workflow

Use this skill for routine contribution work in this repository. It is the default workflow for changes under `app/`, `cmd/`, `proto/`, `scripts/`, `x/**`, `docs/`, `Makefile`, and `.github/workflows/`.

## Rules

1. Read the touched files and the nearest build, generator, or CI entrypoint before editing.
2. Edit source files, not generated files, unless there is no generator.
3. Use `make` targets as the only validation interface.
4. When several `make` targets could apply, prefer the broader one.
5. Regenerate derived artifacts before finishing and inspect the final diff.

## Choose validation by change type

- Go code, imports, or dependencies: run `make lint-go` and `make test-go`.
- Proto files: run `make proto`.
- CLI command or Cobra flag changes: run `make doc-command`.
- Predicate, Prolog library, bootstrap, or predicate feature changes: run `make doc-predicate` and `make test-go`.
- Interface changes that affect mocks: run `make mock`.
- App wiring, local chain behavior, or upgrades: run `make build-go`.
  For local chain validation, use `make chain-init`, `make chain-start`, and `make chain-stop`.
  For upgrade work, use `make chain-upgrade FROM_VERSION=... TO_VERSION=...`.
- Generated markdown under `docs/`: run the relevant generation target from `Makefile`. CI will additionally check MDX compatibility.

## Gotchas

- Tests must follow the repo conventions: GoConvey is allowed, `testify`, `ginkgo`, and `gomega` are forbidden.
- `gofumpt` is mandatory, but the correct entrypoint is still `make lint-go`.
- CI regenerates command docs, proto docs, and predicate docs and fails on unexpected diffs.
- `make doc-command` rewrites machine-specific defaults in generated markdown. Do not hand-edit those outputs afterward.
- `make proto-gen` copies generated code back from `github.com/axone-protocol/axoned/*` into the repo root. Inspect the diff after regeneration.
- `make build-go-all` is host-limited and can fail on cross-compilation mismatches. Prefer `make build-go`.
- `make chain-stop` runs `killall axoned`.

## High-signal files to read first

- `Makefile`
- `.github/workflows/lint.yml`
- `.github/workflows/test.yml`
- `.golangci.yml`
- `README.md`

For subsystem-specific work, also read the closest source and generator entrypoints:

- Predicate docs: `scripts/generate_predicates_doc.go`
- Command docs: `scripts/generate_command_doc.go`
- Proto generation: `scripts/protocgen-code.sh` and `scripts/protocgen-doc.sh`
- Upgrade work: `app/upgrades.go`

## Final self-check

- The validation commands you ran match the change type.
- Any generated file diffs are expected and trace back to source edits.
- `git status --short` contains only intended files.
- If `go.mod` or `go.sum` changed, they were updated intentionally and are tidy.
