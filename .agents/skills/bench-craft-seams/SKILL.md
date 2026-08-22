---
name: craft-seams
description: Vocabulary and principles for placing test seams and designing deep modules. Use whenever someone asks "where should this test live", "where does the interface go", or "is this the right boundary to test", deciding whether to introduce an abstraction, or when /bench-write-spec is picking seams — even on small changes.
index: placing a test / designing an interface
---

# Seams

A **seam** is a place where you can change behavior without editing in that
place. It is the location where a module's interface lives, and where both
callers and tests cross. Where to put the seam is a design decision, separate
from what goes behind it. Get it right and tests are cheap, stable, and
survive refactors. Get it wrong and you test internals, the suite breaks on
every change, and TDD becomes a tax.

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

This is why `/bench-write-spec` picks seams before `/bench-implement-spec`
runs. A test attached to a well-placed seam checks behavior the user cares
about, so an agent can't satisfy it by over-fitting to incidental
implementation detail. A test attached to an internal is something the agent
can game.

## Picking the seam for a feature

- Prefer an existing seam to a new one.
- Before you place a test in a package that owns a registry or an inventory,
  read that package's own admission rule. Such a package carries a seam
  contract, not only a location, so a test that compiles and passes there can
  still be inadmissible.
- Use the **highest** seam that still exercises the real behavior — *and at which
  the failure modes are still observable*. If an error path can't go red from
  the high seam, add one lower coverage row, or make the failure observable
  there. Don't drop the case for the sake of seam height. Higher seams mean
  fewer, more stable tests, and more of the implementation stays free to change.
- The ideal number of new seams for a feature is one. Justify any beyond that.
- When the seam crosses a dependency, classify it — in-process,
  local-substitutable, remote-owned, or true-external — and take the test
  strategy from `references/dependency-categories.md`.

Avoid the word "boundary" (overloaded). Say **seam** or **interface**.

Before exploring the tree for where a seam belongs, declare a small read budget
— a number of files. If you spend it without finding traction, stop guessing
and reroute through `bench outline`. Report the budget spent alongside the
reroute; a silent, unbounded read loop is not an accepted way to find a seam.

The edge classes a domain's inputs actually present — per surface: shell CLI,
HTTP API, web UI, background jobs — are templated in
`references/hostile-input-library.md`. Quarry it when an edge inventory or a new
project profile needs a hostile-input checklist.

## Design it twice (for the genuinely uncertain seam)

Your first interface idea is rarely the best. When a seam is the uncertain one
— the one you declared a high-effort line for — don't settle on the first
shape. Spawn 3+ parallel sub-agents, each designing a **radically different**
interface for the same module, then compare and recommend. The paste-ready
sub-agent briefs and the comparison method are in
`references/design-it-twice.md`. That is what a high-effort line buys at the
uncertain seam: more interfaces considered before one is chosen. At a known
seam, skip it — one design is fine.

## Splitting when the structure gate fires

`bench structure` flags structural debt when a file outgrows its line budget
or a directory collects too many source files. Split along responsibility,
never line count, and never fragment to dodge the limit. Before choosing split or grant, check both the file-length budget and the directory's file-count headroom: a split is free only when the directory has room.

A
genuinely deep module can instead earn a per-path grant. Propose a line in
`.bench/structure.budgets` (`<path> <max>`, trailing `/` for a directory) for
the reviewer to approve. The file is reviewer-owned; never edit it yourself.
The full splitting method — the deletion test per split, crowded directories,
the deep pass — is in `references/structure-gate-splitting.md`.
