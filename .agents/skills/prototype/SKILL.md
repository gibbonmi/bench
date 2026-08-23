---
name: prototype
description: Build a disposable prototype to answer one named question with running code. Use when the reviewer says "prototype this", "spike it", or "try an approach", or when a Prototype decision ticket from /bench-shape-idea needs a reviewer choice made concrete.
index: spiking a disposable prototype
---

# Disposable prototypes

A prototype exists to answer one question, then disappear. Its deliverable is
the verdict, never the code.

## Name the question first

Before writing anything, state the question the prototype answers, in one
sentence, to the reviewer or in the decision ticket that asked for it. One
question per prototype — two questions means two prototypes, because a verdict
must map to exactly one question to be usable.

```
Can the TOON parser stream a 50MB file under 200ms without loading it whole?
```
Good — named, falsifiable, and the prototype's verdict answers it directly.

```
Explore how the parser might handle large files.
```
Bad — no question, so nothing tells you when the prototype is done or what
verdict to record.

## Build rules

- **Trivial to run.** One command or one pasted snippet reproduces the result.
  If running it needs setup steps, the setup is part of the artifact and the
  prototype is already too heavy for its question.
- **State stays in memory** — unless persistence *is* the question. A map, a
  slice, a hardcoded fixture. Wiring a store answers a question nobody asked
  and makes the artifact feel too valuable to discard.
- **Surface the relevant state.** Print or display the state the question
  turns on, visibly, every run — the verdict is read off that surface, not
  inferred from silence.
- Production standards do not apply. Skip error handling, naming polish, and
  tests except where the question itself needs them.

## Record the verdict, then discard

The verdict is the answer plus the one observation that decided it. Record it
where the question was asked: the conversation with the reviewer, or the
Prototype decision ticket. Then delete the artifact. Write prototypes in the
session scratch directory or an isolated worktree, so discard is one removal.
Confirm with `git status` that nothing leaked into the tree.

A prototype that survives becomes production code that skipped review, spec,
and the gate — discard is what licenses its speed. There is no keep route:
promoting the *learning* means building fresh through the normal workflow,
with the verdict as input.
