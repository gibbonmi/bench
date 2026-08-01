# Retro: reduced-gate-phase-set

Landed `6f3486a`, 2026-08-01. Ten tickets, 26 acceptance-coverage rows, five tests
beyond the map, a three-axis review, three review repairs, and two defects the gate
caught that nothing else did.

## What the build proved about its own design

The spec's central bet was that excludability should be **proven by construction**
rather than asserted: run the excludable phases against a tree with the declared paths
absent, so a hidden dependency fails loudly. That bet paid twice on the first real
composed gate.

It found that the **canary reads live `specs/` content** — 49 decision-map fixtures
mutating a real spec file. A plain "skip these paths" exemption would have shipped green
and quietly stopped exercising the decision-map checker. Instead the gate refused, named
the file, and the fixtures now own their material.

It also found the enforcement predicate itself was wrong, twice, both times mine:

| version | rule | why it failed |
|---|---|---|
| 1 | red on any structured skip | host `fifo`/`privilege` skips false-red every gate |
| 2 | red on any `kind=environment` skip | 146 environment skips are normal on a full run |
| 3 | red when the missing path is a scope member | correct |

Version 2 is the instructive one. I checked that *capability*-kind skips occur normally
and would false-red, and never checked the environment count — which was printed on the
same line of the same output I was reading. Reasoning about a taxonomy in the abstract
beat reading the number sitting next to the one I had already read.

## What the process got right

**Coordinator probes earned their place.** Each delegate's own probes tended to aim at
behavior it had already tested. The independent probes found what those missed: an
uncovered `Status == "green"` conjunct in the prep-release refusal, and confirmation that
several guards were live in both directions rather than vacuously true. The discipline
that mattered was choosing a *different kind* of mutation than the delegate used.

**Fence-stops worked.** Three times a delegate hit a wall and reported instead of
widening — a second `ReducedScope()` caller in another process, a fixture that had been
relying on the fail-open bug, and a helper change needing an out-of-fence emitter. Each
time the thing outside the fence was real, and each stop produced a better-scoped follow-up
than a silent breach would have.

**Two delegate claims did not survive checking**, and the checking is what mattered more
than the outcome: one mutation red that reproduced only after a rebuild (I was wrong, its
evidence was right), and one canary change that stopped twice mid-work and needed its
answers pulled apart by hand. Verify against the tree, not the claim — including when the
claim is a correction to you.

## What the process got wrong

**Provisional cadence is fast and its recovery paths are not.** Six tickets integrated
with zero full gates, five delegates ran concurrently, and the wall-clock win was real.
But two necessary commits to `main` mid-run deadlocked the lifecycle permanently: every
operation routes to `promote` for recomposition, and `promote` checks a clean review and
released assignments *before* the recompose branch — conditions a mid-repair run cannot
meet. `abandon` is blocked by the same gate, so the run could not even be retired. The
build finished as light-path work with the candidate applied as a patch.

**I did not check the precondition order.** When the first commit triggered recomposition,
I verified that `recompose.go` fast-forwards from `run.Base` and concluded a moved tip was
safe. It is safe *at that step*. I never read the checks above it, then committed again
and compounded it. The lesson is narrow and mechanical: when a tool says "do X to
recover", read what guards X before assuming X is reachable.

**Stale binaries cost three round trips.** A mutation probe against a contract test is a
silent no-op unless `dist/bench` is rebuilt — and the failure is asymmetric: a stale run
yields PASS, never FAIL. A probe that reds is trustworthy; a probe that passes may have
tested nothing. Same family as the seal refusal on the light-path commit and the stale
gopls diagnostics that cried wolf four times.

## Decisions worth remembering

- `specs/` joined the allowlist on the reviewer's ruling, and conformance still grades
  spec format on every reduced run — the exemption never made specs ungraded.
- Reduction and stripping apply only to the kit's own root. A linked repo keeps today's
  behavior until FT144's two-audience work gives it a declaration.
- The spec's literal "capabilities required" wording was wrong and all three review axes
  agreed the implementation should deviate. Following it would have rebuilt the silent
  green the spec was written to prevent.

## Left behind, deliberately

- A permanently `active` spec-build run record for this slug, with six registered
  worktrees and its provisional refs, retirable by no sanctioned operation.
- Five recovery refs from the earlier cadence build that `bench worktree recovery` refuses
  because it cannot prove their payloads landed. The refusal is correct; the refs are
  inert.
- Two conformance checks reachable only through `conformance-suite`'s whole-package run,
  by test-name prefix rather than registry registration.
