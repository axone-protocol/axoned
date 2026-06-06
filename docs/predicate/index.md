---
title: Predicates
sidebar_label: Introduction
sidebar_position: 0
---

Predicates are the relations available to AXONE Prolog programs. They are the building blocks used to describe facts, express rules, transform data, inspect the execution context, and call AXONE-provided capabilities from inside the Prolog VM.

This reference focuses on the predicate layer itself. It does not explain the module API, transaction flow, or query messages; those concerns belong to the protocol documentation. Here, the important question is what a Prolog program can call once it is evaluated by AXONE.

AXONE exposes several kinds of predicates:

- standard Prolog-compatible predicates for unification, terms, lists, arithmetic, streams, control flow, and dynamic clauses;
- reusable library predicates loaded with `consult/1`, such as helpers for chain context, bank state, codecs, cryptography, DID values, URIs, and WASM queries;
- lower-level predicates that interact with AXONE-managed capabilities through the Prolog VM.

For example, a program can load the chain library and define a rule in terms of the current block header:

```prolog
:- consult('/v1/lib/chain.pl').

current_chain_height(Height) :-
  header_info(Header),
  Height = Header.height.
```

The same pattern applies to other AXONE libraries: load the library that provides the relation you need, then compose that predicate with your own facts and rules.

## Execution Model

Predicate evaluation in AXONE is deterministic. Programs are evaluated inside the AXONE Prolog VM, with host capabilities exposed through controlled AXONE surfaces rather than arbitrary filesystem or network access.

Predicates cannot mutate chain state. Even predicates that read chain-backed data, verify signatures, transform encoded values, or call supported devices are constrained to deterministic evaluation.

Execution is also bounded. Predicate calls run under the limits enforced by the logic environment, including gas consumption, result limits, and captured output limits. Programs should therefore keep relations focused, prefer explicit bounds when enumerating solutions, and use library predicates that match the intended capability instead of rebuilding low-level access patterns unnecessarily.

## Using This Reference

Each generated predicate page documents the predicate indicator, signature, description, and examples when available. Library predicates also show the module to load before use.

Use this introduction for the mental model, the generated pages for predicate-level details, and the [Virtual File System](./vfs.md) page when you need to understand the `/v1/...` paths that back some AXONE capabilities.
