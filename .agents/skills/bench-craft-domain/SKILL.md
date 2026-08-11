---
name: craft-domain
description: The domain-modeling companion — canonical terms with Avoid lists, concrete scenarios at concept edges, producer-derived equivalence partitions, code-versus-claim checks, and glossary-only CONTEXT.md upkeep. Use during grilling, shaping, or spec authoring, or whenever a term feels vague, overloaded, or contradicted by the code.
index: pinning domain terms / enumerating concept-edge scenarios
---

# Domain modeling

A domain model is the shared vocabulary plus the concrete scenarios that pin
what each term means at its edges. Getting it right upstream is what makes
completeness cheap downstream: a spec row or a review can only enumerate what
the vocabulary can name. This is a companion craft skill — not a phase, not a
gate. It fires from `craft-grill`, `/bench-shape-idea`, and
`/bench-write-spec`; every other phase consumes `CONTEXT.md` ambiently.

## Canonical terms and Avoid lists

Give every load-bearing concept one canonical name and challenge the vague or
overloaded ones on sight — a term that covers two meanings, or two terms that
cover one, is a defect to surface now, not later. Record each canonical term
with an Avoid list: the synonyms and near-misses sessions must not drift into.
Complete when every load-bearing noun in the work at hand resolves to one
canonical term or to an open question put to the reviewer.

> **worktree** — an isolated Git checkout leased for one assignment. Not
> "sandbox", not "branch copy" — worktree.

Good — one name, a boundary, and the drift words named.

> **worktree** — a separate place where work happens.

Bad — no edge and no Avoid list, so every session re-invents the meaning.

## Concrete scenarios at concept edges

Where two concepts meet, or where a term's coverage ends, write a concrete
scenario: exact inputs and the observable outcome, never a category. "Handles
bad input" names a family; "an empty input list exits 2 and prints the input
path" names one behavior. Complete when each edge in the work at hand carries
a scenario sharp enough to become an acceptance row unchanged.

## Producer-derived equivalence partitions

Derive input families from what the real producer can emit, not from
imagination. Read the producer — the caller, the CLI, the file format, the
API — and take each output shape it can actually produce as one partition.
Imagined inputs pad the inventory while missing the shapes that occur;
producer-derived partitions bound completeness because the producer's range
is finite and readable. Complete when every partition traces to a producer
surface you read, and the producer's whole range is accounted for.

## Code-versus-claim comparison

Check what the code actually does against what the glossary, spec, or
conversation claims it does. On a conflict, resolve it explicitly: fix the
code, fix the term, or put the fork to the reviewer as a decision — never
leave both meanings live, because each survives into different downstream
artifacts and they drift apart.

## Inline glossary maintenance

As terms resolve, update `CONTEXT.md` in the same pass. The file stays
glossary-only — terms, definitions, and Avoid lists; workflow prose lives
elsewhere. A term resolved in conversation but not written into `CONTEXT.md`
is not resolved: the artifact is the source of truth, not the chat.

## Hard-to-reverse decisions

A hard-to-reverse architectural decision surfaced here routes to `craft-adr`,
which owns the recording format. This skill carries no ADR format of its own.
