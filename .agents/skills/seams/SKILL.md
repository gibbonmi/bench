---
name: seams
description: Vocabulary and principles for placing test seams and designing deep modules. Use whenever deciding where a test should attach, whether to introduce an abstraction, how to make code testable, or when /spec is picking seams. Reach for this any time the question is "where does the interface go" or "is this the right boundary to test" — even on small changes.
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

This is why `/spec` picks seams before `/build` runs: a test attached to a
well-placed seam checks behavior the user cares about, so an agent can't satisfy
it by over-fitting to incidental implementation detail. A test attached to an
internal is something the agent can game.

## Picking the seam for a feature

- Prefer an existing seam to a new one.
- Use the **highest** seam that still exercises the real behavior. Higher seams =
  fewer, more stable tests = more of the implementation free to change.
- The ideal number of new seams for a feature is one. Justify any beyond that.

Avoid the word "boundary" (overloaded). Say **seam** or **interface**.

## Design it twice (for the genuinely uncertain seam)

Your first interface idea is rarely the best. When a seam is the uncertain one —
the one you declared a high-effort line for — don't settle on the first shape.
Spawn 3+ parallel sub-agents, each designing a **radically different** interface
for the same module under a different constraint:

- minimize the interface — 1–3 entry points, maximum leverage each
- maximize flexibility — many use cases, room to extend
- optimize the common caller — make the default case trivial
- ports and adapters — if the seam crosses a real dependency

Each returns its interface (with invariants, ordering, error modes), a usage
example, what it hides, and where its leverage is thin. Then compare them on the
three things that matter — depth (leverage at the interface), locality (where
change concentrates), seam placement — and give an opinionated recommendation, or
a hybrid if elements combine well. A menu isn't the deliverable; a strong read is.

This is what a high-effort line buys at the uncertain seam: not more code, but more
interfaces considered before one is chosen. At a known seam, skip it — one design
is fine.

## Splitting when the structure gate fires

`bench structure` fails the gate when a file outgrows its line budget or a directory
collects too many source files. That's the signal to modularize *now*, mid-shift —
not to defer it. But how you split decides whether you've reduced debt or just moved
it:

- **Split along responsibility, not line count.** Find the seam already latent in
  the file — the cluster of functions with one reason to change — and lift it into
  its own module behind a small interface. The line count drops as a side effect of
  a real boundary, never as the goal.
- **Apply the deletion test before every split.** Would extracting this cluster
  *concentrate* complexity behind one interface, or just *move* it to another file?
  Only "concentrates" is a real deepening; "just moves it" is fragmentation wearing
  a refactor's clothes.

Bench's **deep pass** is this skill applied broadly — run `bench structure` (tighten
the budgets if needed) and work the deletion test across every flagged module — and
it's fully self-contained; nothing outside the kit is required. The deletion test is
borrowed from Pocock's `improve-codebase-architecture`, which Bench does *not* bundle:
if you happen to have his skills installed, `/improve-codebase-architecture` is a
richer external version of the same review (organic exploration, a visual before/after
HTML report). It's an optional upgrade, never a dependency.
- **Never fragment to beat the number.** Slicing a cohesive file into `part_a` /
  `part_b` to dodge the limit is worse than the long file: it scatters one concept
  across files with no interface between them. If you can't name the responsibility
  the split isolates, don't split — the file may legitimately be one deep module, so
  raise its budget in the gate instead.
- **A crowded directory is an ungrouped module.** Thirty files in one dir means
  several modules are hiding in a flat namespace. Group related files into a package
  with a clear entry point; the package *is* the seam.

The deep-module rule applies throughout: prefer a few files with real interfaces
over many shallow ones. The gate measures the symptom; this is the cure.