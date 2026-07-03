---
name: craft-seams
description: Vocabulary and principles for placing test seams and designing deep modules. Use whenever deciding where a test should attach, whether to introduce an abstraction, how to make code testable, or when /bench-write-spec is picking seams. Reach for this any time someone asks "where should this test live", "where does the interface go", or "is this the right boundary to test" — even on small changes.
index: placing a test / designing an interface
---

# Seams

A **seam** is a place where you can change behavior without editing in that place
— the location where a module's interface lives, and where both callers and tests
cross. Where to put the seam is a design decision, separate from what goes behind
it. Get it right and tests are cheap, stable, and survive refactors. Get it wrong
and you test internals, the suite breaks on every change, and TDD becomes a tax.

## Design deep modules

A **deep module** hides a lot of behavior behind a small interface. The interface
is everything a caller must know: signature, invariants, ordering, error modes,
config, performance. Aim for leverage (callers learn a little, get a lot) and
locality (changes concentrate in one place).

- Can I drop a method? Simplify a parameter? Hide more inside?
- **Deletion test:** imagine deleting the module. If complexity vanishes, it was a
  pass-through. If complexity reappears across callers, it earned its place.
- **One adapter is a hypothetical seam; two adapters is a real one.** Don't add a
  seam until something actually varies across it.

## The interface is the test surface

Callers and tests cross the same seam. Test through the interface, exercising
observable behavior — not private internals. If you want to test *past* the
interface, the module is the wrong shape; fix the shape, don't reach around it.

```
new Scheduler(clock).nextRun()
```
Good — time is a constructor parameter, so a test injects a fixed clock through
the same seam callers use.

```
class Scheduler { nextRun() { return Date.now() + this.interval } }
```
Bad — the dependency is hardwired, so a test must patch `Date.now` behind the
interface, binding itself to an internal.

This is why `/bench-write-spec` picks seams before `/bench-implement-spec` runs: a test attached to a
well-placed seam checks behavior the user cares about, so an agent can't satisfy
it by over-fitting to incidental implementation detail. A test attached to an
internal is something the agent can game.

## Picking the seam for a feature

- Prefer an existing seam to a new one.
- Use the **highest** seam that still exercises the real behavior — *and at which
  the failure modes are still observable*. If an error path can't go red from the
  high seam, add one lower coverage row or make the failure observable there;
  don't drop the case for the sake of seam height. Higher seams = fewer, more
  stable tests = more of the implementation free to change.
- The ideal number of new seams for a feature is one. Justify any beyond that.

Avoid the word "boundary" (overloaded). Say **seam** or **interface**.

The edge classes a domain's inputs actually present — per surface: shell CLI,
HTTP API, web UI, background jobs — are templated in
`references/hostile-input-library.md`. Quarry it when an edge inventory or a new
project profile needs a hostile-input checklist.

## Design it twice (for the genuinely uncertain seam)

Your first interface idea is rarely the best. When a seam is the uncertain one —
the one you declared a high-effort line for — don't settle on the first shape:
spawn 3+ parallel sub-agents, each designing a **radically different** interface
for the same module, then compare and recommend. The paste-ready sub-agent
briefs and the comparison method are in `references/design-it-twice.md`. That
is what a high-effort line buys at the uncertain seam: more interfaces
considered before one is chosen. At a known seam, skip it — one design is fine.

## Splitting when the structure gate fires

`bench structure` flags structural debt when a file outgrows its line budget or a
directory collects too many source files. Split along responsibility, never
line count, and never fragment to dodge the limit. A genuinely deep module can
instead
earn a per-path grant: propose a line in `.bench/structure.budgets`
(`<path> <max>`, trailing `/` for a directory) for the reviewer to approve —
the file is reviewer-owned; never edit it yourself. The full splitting method —
the deletion test per split, crowded directories, the deep pass — is in
`references/structure-gate-splitting.md`.
