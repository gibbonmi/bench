# Canary phase cost (FT40)

The canary phase dominates the gate wall: 66 fixtures plus a vacuity baseline,
each materializing a temp repo and running the full inner gate. Every candidate
cut touches oracle semantics — the canary is the gate guarding the gate, so any
cost reduction must preserve "a rotted check fails the sweep loudly."

## Measured facts (2026-07-07, 16-core WSL2)

- Gate wall 130s; canary phase 128s of it (roadmap measurement). A solo
  `bench canary <root>` sweep re-measured green at 129.7s wall.
- One inner gate, solo, warm build cache: ~1.9s sequential — conformance 0.4s,
  contract 1.3s, shellcheck 0.2s.
- The sweep runs fixtures on `min(NumCPU, fixtures)` workers (16 here), yet the
  solo sweep burned only ~282s CPU over 130s wall — ≈2.2 cores busy. The
  workers are sleeping, not computing: see #1's answer for the mechanism.
- Outer-phase solo walls against the real kit: contract 5.8s; conformance and
  shellcheck under a second each. The gate wall is canary or nothing.
- FT39 corroboration (mechanism evidence, not the captured repro FT39's
  graduation requires): every gate run has canary spawning ~18 concurrent
  copies of the concurrent-acquire poll storm alongside the outer contract
  phase's own passing run — exactly the "under load" condition in which the
  batch session saw the outer test miss its fixed 60s window.
- Fixture → phase targeting is clean today: the 15 `behavior-owned` fixtures
  assert contract-phase diagnostics; the other 51 fixtures (7 families) assert
  conformance-phase diagnostics; none target shellcheck.
- The gate verdict cache is already whole-tree-keyed with a no-forged-verdict
  guarantee (`<git-dir>/bench-last-gate`); `bench commit` reuses a fresh green
  verdict for the identical tree. Any canary-specific skip only pays on trees
  that *changed* — i.e. it needs a canary-relevant *subset* key, and choosing
  that subset is where the risk lives.

## #1: Where does the 128s actually go, and how far do semantics-neutral fixes reach?

Blocked by: —
Type: Prototype

### Question
Effective parallelism of the sweep is ~1 despite 16 workers. Until the
contention source is identified, we cannot know whether the oracle-semantics
cuts are needed at all. Candidates to test: per-invocation `go test` cost
(package resolution + vet + binary staleness checks, ~66×2 invocations),
process/filesystem thrash from concurrent subprocess-heavy contract runs,
build-cache lock serialization. Candidate semantics-neutral fixes to prototype:
precompile the conformance and contract test binaries once per sweep
(`go test -c`) and exec them per fixture; tune worker count; anything the
profile suggests. Also capture the outer contract phase's solo wall — it is the
budget anchor for #2. If semantics-neutral fixes alone bring the sweep at or
under the slowest other phase, #3 collapses to "no semantics change".

### Answer
The wall is one bug-shaped test behavior, not contention.
`bench_worktree_concurrent-acquire_contract` (the FT39 test) polls a record
file on a fixed 60-second deadline and never checks whether its two
`bench worktree` children are still alive. Against 18 of the 66 fixture trees
(all 15 behavior-owned, plus missing-cli-inventory, stale-cli-doc-reference,
extensionless-gate-ref — trees plausible enough for the runtime suite to
engage), both children exit within milliseconds, no record can ever appear,
and the test sleeps out the full minute anyway. Eighteen 60s sleeps over 16
workers form two waves ≈ 128s; everything else in the sweep is ~2s of real
work per fixture and parallelizes cleanly.

Prototyped fix (semantics-neutral, ~10 lines in the test's poll loop): treat
"both runs exited with fewer than two records" as a terminal failure instead
of sleeping to the deadline. Measured: sweep 128s → **25.8s**; the passing
path is unchanged and stays green; no fixture EXPECT matches either the old
or new failure message, so no canary semantics move. The fix is correct
independent of canary — a run that has exited can never record — and it
removes most of the load FT39's flake needs. The prototype was reverted; the
production change is a contract-test edit that needs reviewer sign-off (#6).

The other semantics-neutral candidates are dead ends at these numbers:
precompiled test binaries and worker tuning address overheads of ~0.3s per
invocation, noise next to the 60s sleeps. Remaining post-fix composition:
~26s ≈ 66 fixtures × (conformance 0.2–0.4s + contract 1.3–4.6s under
contention + shellcheck 0.2s) / 16 workers. Outer contract solo wall is 5.8s,
so even post-fix the canary remains the gate's critical path — #2 and #3 stay
live, with 26s as the new baseline.

## #2: What is the budget — how fast is fast enough?

Blocked by: #6
Type: Grill

### Question
The cut list can't be sized without a target. Proposed principle: the canary
must come off the gate's critical path — its wall should be at or under the
slowest other phase, so the gate wall is set by real checks, not the
meta-check. Anchors are now measured (#1): outer contract 5.8s solo; canary
25.8s after the #6 fix lands. Strict reading of the principle means #3's cuts
proceed until the sweep is ≤ ~6s; the lenient alternative is "a ~26s gate is
fine, stop here." Reviewer call.

### Answer
Lenient (2026-07-07): the ~26s post-#6 gate is accepted. The canary stays the
gate's critical path and no further second is bought with oracle-semantics
risk. If the gate wall creeps later, scope-to-phase (the only fails-safe cut)
is the documented first resort. This resolves #3 to "no cuts" and retires the
contingent #4/#5.

## #3: Which oracle-semantics cuts are acceptable, if the budget demands more than #6 delivers?

Blocked by: #2
Type: Grill

### Question
Three candidates, combinable, each with a distinct risk profile. Post-#6
numbers: the sweep is ~26s of real work — conformance 0.2–0.4s + contract
1.3–4.6s + shellcheck 0.2s per fixture across 16 workers.

- **Scope each fixture's inner gate to the phase its EXPECT targets.** Fails
  safe: if a check migrates phases, the scoped inner gate goes green, the
  fixture reports "did not bite", the sweep goes red loudly. Win: drops the
  contract suite — the dominant per-fixture cost — from the 51
  conformance-targeting fixtures; rough estimate halves the sweep or better.
  Needs a phase-declaration mechanism (#4) and a per-phase vacuity baseline.
- **Skip the sweep when canary-relevant inputs are unchanged since last green.**
  Fails dangerous: a missed key input yields a silently vacuous canary — the
  exact failure class the canary exists to catch. Only pays on changed trees
  whose diff is outside the subset (doc-only edits are the common case).
- **Batch fixtures per inner run.** Cross-contamination both directions: one
  fixture's breakage can mask another's diagnostic (false "did not bite") or
  emit another's EXPECT substring (false pass). Feasibility unproven; choosing
  it spawns a prototype ticket.

### Answer
None (2026-07-07): the lenient budget (#2) is met by the #6 fail-fast alone.
The three candidates above stay on the page as the menu if the gate wall
creeps; scope-to-phase is the first resort because it fails safe. The
contingent design tickets (#4 phase declaration, #5 skip-cache keying) were
deleted with their branches.

## #6: Adopt the concurrent-acquire fail-fast, and does FT39 fold into this work?

Blocked by: —
Type: Grill

### Question
#1 proved the sweep's 128s is 18 fixtures sleeping out the FT39 test's fixed
60s deadline, and prototyped the fix: the poll loop treats "both `bench
worktree` runs exited with fewer than two records" as terminal failure.
Measured 128s → 25.8s, passing path unchanged. Two reviewer calls: (a) approve
the contract-test edit — it changes the oracle's failure latency, never its
verdict on a live subject, but test edits are sign-off territory; (b) decide
whether FT39 (parked pending a captured repro) folds into this build or stays
parked — the fail-fast removes most of the load its flake needs but does not
extend the window the outer test's legit red would come from.

### Answer
Both calls resolved (2026-07-07). (a) Approved: the contract-test fail-fast
edit ships as a production change in this build. (b) FT39 stays parked — its
roadmap row and graduation criterion (a captured red with the exact message
from a real `bench gate` run) remain in force; a deadline fix without that
repro would be speculative. Parking forecloses nothing: if the flake survives
the fail-fast, the captured red graduates FT39 and the fix is built then.

## Handoff

1. **Module boundaries.** One unit: the overlap wait loop inside the
   concurrent-acquire runtime contract test
   (`bench_worktree_concurrent-acquire_contract`). Inside: how the loop decides
   "keep waiting" vs "fail now". Outside and unchanged: the barrier design
   (record file → go-file release), the two-run spawn, the post-barrier
   assertions, the canary sweep, all production code.
2. **Contracts.** The test's observable interface is its verdict and failure
   message. New terminal condition: both spawned `bench worktree` runs have
   exited while the record file holds fewer than two lines → immediate fatal,
   with a message distinct from the existing timeout message. The 60s deadline
   stays as the backstop for live-but-stuck runs. Passing path unchanged.
3. **Deep vs thin.** n/a — a ~10-line edit inside one existing test; no new
   module or seam.
4. **Black-box assertables.** Gate stays green (the passing path is
   byte-identical); the canary sweep stays green with all 66 fixtures biting
   and the vacuity baseline clean; sweep wall drops from ~128s to ~26s
   (prototype-measured 25.8s).
5. **Gate attachment.** The gate itself exercises both paths: the outer
   contract phase runs the passing path every gate, and the canary sweep
   drives the new terminal branch ~18 times per sweep (broken fixture trees
   make both children exit early). The latency win has no gate assertion —
   verify the wall drop manually once after the change lands.
6. **Hostile-input owners.** n/a — test-only edit with no new input surface;
   the loop consumes only its own children's exit signals and its own record
   file. The profile checklist classes attach to CLI surfaces this change
   doesn't touch.
7. **Uncertainty flags.** None — the prototype settled mechanism and numbers.
   One implementation freedom: how the loop observes "both exited"
   (non-blocking drain of the existing done channel is the natural shape);
   spec-writer's choice, no escalation needed.
8. **Rejected alternatives.** Strict ~6s budget and all three oracle-semantics
   cuts (#3) — rejected under the lenient budget; scope-to-phase is the
   documented first resort if the wall creeps. Precompiled test binaries and
   worker-count tuning — measured dead ends (~0.3s effects vs 60s sleeps).
   FT39 fold-in (raised or event-keyed deadline) — stays parked pending its
   captured repro.
9. **Domain watch-outs.** The fail-fast is sound because the fixture shell
   records before it holds: a run that exits without recording can never
   record later. The condition is deliberately "both exited" — if one child
   crashes while the other holds, the record stays at one line and the 60s
   deadline backstop still catches it; that rare case is accepted. On the
   terminal branch the go-file release before failing is unnecessary (no
   straggler is alive) but harmless if kept.

Dependency order: n/a — single spec (the fail-fast edit).
