# injected-interface-junctions

Status: staged

Decision source: `capture/audits/injected-interface-composition.md` — the
reviewed injected-interface composition audit (2026-08-03), authored under the
reviewer's standing go-with-recommendation directive of the same date; the
audit table and triage are the post-hoc veto surface. Falsification review:
one mid-tier pass, verdict PASS-WITH-FIXES, all findings folded in.

## Problem

FT181 shipped a story that was green at the specbuild seam and broken
end-to-end because the package's tests drove a fake `WorktreeOwner` while the
real `internal/worktree` planner refused the same state. The audit found the
same gap class live at four more seams: a release owner whose test stub cannot
fail, an abandon fake covering exactly the decayed shapes FT181 found broken,
a canary test runner that re-implements the production runner's resolution
semantics, and guard probes with opposite fail-safe polarities that no test
composes. `ShapeUnknown` — the classifier branch behind one FT181 coordinator
catch — appears in zero test files. Nothing prevents the class from recurring:
the audit was a hand-run sweep.

## Solution

Junction tests: for each P1 gap, at least one test drives the real producer
through the consuming surface, so a softened half turns the gate red. A new
conformance check pins the audit's result — every injected port names its
real-producer test or its written exemption — with a canary fixture proving
the check bites, so the sweep never re-runs by hand.

A junction test that goes red against *unmodified* production is a production
finding: the implementing delegate stops and surfaces it — never adjusts the
expectation or the production side to force green. That rule is the point of
the spec.

## User stories

1. **ReleaseOwner junction.** As the reviewer, I can trust specbuild's
   integrate/release paths because at least one Service-level test wires the
   real `worktree.ReleaseProvisional` — one success and one refusal
   (request/assignment mismatch) observed through the Service surface.
   Line: `sonnet` / low / ~3 iterations. Exact spec, known shape,
   gate-covered.
2. **AbandonOwner decayed-family junction.** As the reviewer, I can trust the
   decayed/husk/unreadable abandon paths because the shapes now graded by the
   synthetic `decayedOwner` fingerprint are also driven through the real
   planner with privilege-free fixtures.
   Line: `opus` / medium / ~4 iterations. The test shape is known
   (re-drive existing cases through the real planner) and a composition
   divergence surfaces as an observed red handled by the Solution's
   stop-and-surface rule; mid buys judgment for that surfacing, not seam
   discovery.
3. **Canary runner junction.** As the reviewer, I can trust the canary sweep's
   gate resolution because a test drives `SweepTier` with the real
   `defaultRunner` against one gate-owned hermetic fixture with a planted
   trivial bash gate, pinning cwd-relative resolution (the exit-127 field
   symptom) with real bash. The package's tests run nested at `GOMAXPROCS=2`
   (profile cold-session note); the fixture stays single and small.
   Line: `sonnet` / low / ~3 iterations. Exact spec, known shape,
   gate-covered.
4. **Gitguard checker junction.** As the reviewer, I can trust the git guard's
   verdicts because `Classify` runs with the real
   `Checker{git.RefResolves, git.BranchExists}` against a temp repo —
   composing both probes' resolved answers with the clobber logic — and once
   with a PATH-front stub `git` that sleeps past the 2s probe bound, composing
   both timeout polarities. The probes resolve in the process working
   directory, so these tests `t.Chdir` and must not run parallel.
   Line: `sonnet` / low / ~3 iterations. Exact spec, known shape,
   gate-covered.
5. **ShapeUnknown fixtures.** As the reviewer, I know an undecidable stat is
   classified fatally, because deterministic privilege-free fixtures pin all
   three `ShapeUnknown` return sites, and one specbuild test composes the
   real classifier's `ShapeUnknown` into the ownership refusal.
   Line: `sonnet` / low / ~3 iterations. Fixture designs are pinned in this
   spec.
6. **Injected-interface conformance check.** As the reviewer, I never re-run
   this audit by hand: a dev-tier conformance check reads a single-source
   registry of injected ports — each naming its real-producer test or its
   exemption reason — verifies the named tests exist where the registry says,
   fails on an unregistered port, and ships with a canary fixture proving the
   red. Includes the compile-time wiring pins for all three
   silently-downgradable owner capabilities.
   Line: `opus` / medium / ~4 iterations. Matches the cached
   gate/conformance-logic routing; the check's own canary covers the output.

## Implementation decisions

- Junction tests live in the **consumer's** package (internal/specbuild,
  internal/canary, internal/gitguard) next to the fake-driven tests they
  complement; classifier fixtures live in internal/worktree. No new packages,
  no new interfaces, no production behavior change.
- Existing fakes stay. The junction tests are additive; per-fence fixture
  tests remain the cheap fast path. Where a fake's in-code note claims the
  real contract lives elsewhere (the `decayedOwner` comment), the note is
  updated to name the junction test that now composes it.
- The conformance registry is a Go table in the conformance package's
  fashion: one row per injected port — package, port name, disposition
  (`real-producer test <file>:<TestName>` or `exempt: <reason>`). A **port**
  is any port-shaped value — an interface, a func type, or a struct of func
  fields — injected through a constructor **or an exported-function
  parameter**; the derivation must see `canary.Runner` and
  `gitguard.Checker`, which are func-shaped parameters, not
  constructor-injected interfaces. The check derives the port inventory from
  the named packages' source rather than a second hand-list, so one side is
  enforcement and the other advertisement per craft-gate's one-source rule.
  Registry rows for the priced-out ports carry the audit's citations (runtime
  contract suite for the gate owners and `buildService`; the no-egress
  exemption for `NPMCLIRegistry`).
- The check registers as a `registry.Check` row (tier `Dev`, subject and a
  valid `InputSource` — go-source), positioned per the registry's
  order-is-contract rule; it emits **four distinct failure messages**:
  unregistered port, registry-named test missing, empty exemption reason,
  zero-derived inventory for a package known to declare ports (fail closed).
  The named-test-exists verification is a **tripwire** (it catches deletion,
  not a junction test gutted back to a fake) and is labeled as such in the
  check's doc; the behavior half of the defense is the junction tests
  themselves plus the canary.
- Compile-time wiring pins in cmd/bench:
  `var _ specbuild.PromotionGateOwner = productionGateOwner{}`,
  `var _ specbuild.ReleaseOwner = productionWorktreeOwner{}`,
  `var _ specbuild.AbandonOwner = productionWorktreeOwner{}` — the
  `releaseOwnerFrom`/`abandonOwnerFrom`/promotion type assertions all
  silently downgrade capability when a method is lost.
- The discarded classifier error at the specbuild precondition consumer stays
  as-is (behavior change priced out below); story 5 asserts the refusal it
  currently produces.

## Testing decisions

- A good test drives the **real producer through the consuming surface** and
  asserts the consumer-observable outcome (refusal string, verdict, sweep
  result) — never the fake's bookkeeping.
- Prior art: `realOwner` (real `worktree.Create` behind the specbuild seam),
  `abandonOwner` (real planner delegation), `greenGate` (real git subprocess
  in tests), `requireUnreadableMetadata` (root-guarded chmod fixture with
  `t.Cleanup` permission restore — the mode-0 classifier fixture restores the
  same way or `t.TempDir` cleanup fails).
- Gate seam: `go test` over the touched packages plus the conformance suite;
  the new check joins the dev tier and the canary sweep.

### Seam diagram

    trigger: go test (dev tier) / bench gate conformance phase
        │
        ▼
    fixture repo/state ──▶ [ consuming surface (Service / Classify / SweepTier / conformance check) ]
                                │ wired with the REAL producer (ReleaseProvisional, planner,
                                │ defaultRunner, git probes, port inventory)
                                ▼
                          consumer-observable verdict (refusal, classification, sweep result, check red)
                      ◀ tests attach at the consumer surface; the real producer is below it

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Service release path completes against real `ReleaseProvisional` when evidence is durable | specbuild Service test with a real-release owner | not classically TDD-able (production presumed working); composition-degenerate probe: relax the producer's request/assignment/path mismatch guard (`ReleaseProvisional`'s retained-checkout refusal) and the refusal row below must go red while the fake-driven tests stay green — probe observed and reverted at build time | a bare-nil release owner passes every existing test; only the real producer refuses bad evidence |
| 1 | Service surfaces the real refusal on request/assignment mismatch | same seam, mismatch fixture | not classically TDD-able; asserts the exact refusal string (`provisional release request, assignment, or path mismatch; checkout retained`) which no fake emits — red proof via the same producer-side probe | pins that refusals propagate through the Service instead of being swallowed |
| 2 | Decayed-directory and husk abandon paths hold with the real planner | specbuild abandon tests re-driven via `abandonOwner` (real) with dir-without-.git / non-dir fixtures | not classically TDD-able; composition-degenerate probe: soften the real planner's shape policy (admit a shape `ClassifyPathShape` refuses) — the fake-driven twin stays green, the junction test goes red — probe observed and reverted at build time | this is FT181's exact failure shape; fake-admitted states must be shown real-admitted |
| 2 | Unreadable-metadata abandon holds with the real planner | same seam, `requireUnreadableMetadata` fixture (root-guard skip) | same producer-side probe class | the fake's synthetic fingerprint cannot diverge from the planner; the real path can |
| 3 | `SweepTier` with `defaultRunner` resolves a relative gate against the fixture cwd | canary test, one gate-owned hermetic fixture, planted trivial bash gate | not classically TDD-able; degenerate probe: point the fixture at a missing gate path and assert the exit-127 diagnostic surfaces through the sweep — probe observed at build time | `resolvingRunner` re-implements this resolution; only the real runner proves bash agrees |
| 4 | `Classify` composed with real probes blocks checkout of an unresolvable ref and permits a resolvable one | gitguard test, temp git repo, `t.Chdir`, no parallel | observed red on first write: no existing test constructs the real `Checker` — write the verdict assertion before wiring the repo state and watch it fail | the constant fakes pin one corner of the polarity matrix; the real repo pins the composition |
| 4 | Timeout polarity composes: a hung git blocks both checkout (ref unresolvable) and forced creation (branch presumed present) | gitguard test, PATH-front stub `git` sleeping past the 2s bound, `t.Chdir`, no parallel | observed red on first write: no test composes the timeout defaults through `Classify` today | the two probes fail safe in opposite directions; only composition shows both land on "block" |
| 5 | `ClassifyPathShape` returns `ShapeUnknown` + error at all three sites | worktree classifier test: ENOTDIR (file-as-parent), ELOOP (self-symlink), EACCES (chmod-0 dir, root-guarded, `t.Cleanup` restore) | observed red on first write: `ShapeUnknown` appears in zero test files, so each site assertion is unwitnessed — write it red per craft-tdd (assert before the fixture exists, or mutate the site) | three return sites, zero fixtures; a regression to `ShapeAbsent` would silently convert refusal into absence |
| 5 | Real classifier's `ShapeUnknown` produces the specbuild ownership refusal | specbuild precondition test: replace `assigned.Path` with a self-symlink (ELOOP) — the ENOTDIR construction needs the pool parent turned into a file, priced out as too heavy | observed red on first write: no test composes this today | the consumer discards the error; only the enum path guards the refusal, so it must be witnessed |
| 6 | Check red on an injected port with no registry row | conformance check + canary fixture | planned red, recorded when observed: run the check before the registry row lands (classic TDD red), retained forever by the canary fixture | an unregistered port is the audit gap recurring; the canary keeps the proof alive |
| 6 | Check red when a registry row names a test that does not exist (tripwire: catches deletion, not decay) | conformance check unit fixture | planned red, recorded when observed during check authorship | a deleted junction test must not leave a green advertisement |
| 6 | Exempt rows require a non-empty reason; zero-derived inventory fails closed | conformance check unit fixtures | planned reds, recorded when observed during check authorship | an empty exemption is an unregistered port with better manners; silent zero-inventory is de-enforcement |
| 6 | `productionGateOwner` and `productionWorktreeOwner` satisfy `PromotionGateOwner`/`ReleaseOwner`/`AbandonOwner` at compile time | cmd/bench compile-time pins | not TDD-able as a test (build-time assertion); degenerate probe: remove one method locally, observe the build red, revert | the runtime type assertions downgrade silently; the pins make method loss a build failure |

### Edge inventory

- Junction test red against unmodified production → row-adjacent rule in
  Solution: stop and surface as a production finding; never force green.
- Root-executing environment defeats chmod fixtures → root-guard skip, prior
  art `requireUnreadableMetadata`; skip carries its reason (skip-ownership
  check applies).
- Platform: ENOTDIR/ELOOP fixtures are POSIX-stable; repo already assumes
  Linux CI (existing FIFO fixture prior art).
- Registry drift: port added without row (canary-covered row), test renamed
  (existence row), exemption emptied (reason row), zero-derived inventory
  (fail-closed row).
- **Won't handle:** driving `Sweep`'s full fixture inventory with the real
  runner on every gate run — one hermetic junction fixture suffices; the
  sweep's breadth stays on the fast fake.
- **Won't handle:** deterministic per-site pinning inside the *specbuild*
  consumer for classifier sites one and three — the classifier-level rows pin
  the sites; the consumer row pins the enum path once via ELOOP.

## Out of scope

- Propagating the classifier error the specbuild precondition consumer
  currently discards into its refusal diagnostics — behavior change, separate
  capability (~3 edits, 2 gate runs).
- The fixture-and-seam inventory generator and receipt-skeleton helper —
  parked in `capture/IDEAS.md`, awaiting the `/bench-what-next` drain.
- Registry coverage of ports in packages outside the audited inventory
  (future packages) — the derivation rule covers the audited set including
  func-shaped ports; a new package binds in when it first declares one
  (~2 edits, 1 gate run).
