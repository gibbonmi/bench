---
name: craft-comments
description: How to write code comments — what earns a comment, what stays out, and the register they use. Use whenever writing or editing code that carries comments, docstrings, or doc comments, or when reviewing comment prose in a diff.
index: writing or reviewing code comments
---

# Writing comments

A comment addresses the teammate who just walked in — the next reader of
the *current* code. That reader has no memory of the conversation, spec, or
diff that produced it. Everything below follows from that reader. Write the
sentences in ASD-STE100 prose per `craft-spec`'s `references/ste-prose.md`.

## What earns a comment

First try to move the fact into code: rename the symbol, extract a function
whose name states it, introduce a type. What code can carry, code carries.
A comment earns its place only for what remains:

- **The why** — the constraint behind a non-obvious choice: the ordering that
  matters, the API quirk being worked around, the edge the simpler version got
  wrong.
- **The warning** — an invariant the next editor could silently break, at the
  point where they'd break it.
- **Doc comments on public symbols** use the language's own convention. Go
  starts a full sentence with the symbol's name; shell writes a one-line
  contract above the function.

A comment that restates the line below it carries nothing — delete it.

## The register

Use the timeless present tense; describe the current state of the code. The change
that produced the code is invisible to the next reader, so the comment never
mentions it: no narration ("now handles", "fixed to use"), no provenance
("per the spec", "as requested in review"), no argument for its own
correctness. That prose addresses the reviewer of the diff — it is PR-talk,
and it becomes noise the day the diff merges. What survives merge is only
what's true of the code as it stands.

Match the surrounding file: its comment density, its idiom, its doc-comment
shape. A sparse file stays sparse.

```go
// Retry once: the registry returns 503 during its nightly compaction window.
```
Good — a constraint the code can't show, stated timelessly.

```go
// Changed to retry here, which correctly handles the 503 case from review.
```
Bad — narrates the change and argues correctness to a reviewer; the next
reader learns nothing about when or why the retry matters.

## Aging

Comments rot when they duplicate what code already enforces. A *what* comment
needs an update in lockstep with every edit; a *why* comment stays true until
the constraint itself moves. A comment that states only the why stays cheap
to own. When you edit code under an existing comment, read the comment as
part of the diff. Update it or delete it, and never leave it describing the
code that was.
