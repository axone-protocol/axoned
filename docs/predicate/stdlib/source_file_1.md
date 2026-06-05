---
sidebar_position: 77
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# source_file/1

## Module

Built-in predicate.

## Description

True when File is one of the Prolog source files loaded in the current
interpreter.

## Signature

```text
source_file(?File) is nondet
```

## Examples

### Match a loaded source file

This scenario demonstrates checking whether a source file has been loaded.

Here are the steps of the scenario:

- **Given** the query:

```  prolog
consult('/v1/lib/lists.pl'),
source_file('/v1/lib/lists.pl').
```

- **When** the query is run
- **Then** the answer we get is:

```  yaml
height: 42
gas_used: 4301
answer:
  has_more: false
  variables:
  results:
  - substitutions:
```
