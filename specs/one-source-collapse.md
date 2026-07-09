# one-source-collapse

Status: implemented

Spec source: FT50 roadmap row + `ASSESSMENT.md` item 9 (reviewer-directed batch
drain — no decision map; every defaulted decision below is flagged **[default]**
for post-hoc veto).

## Problem

Three one-source violations and one stale comment survive in the kit, each a
silent-break risk: `bench guards` matches the pre-push marker with a copied
literal (change the marker and guards misreports a managed hook as unmanaged);
the gate-cache filename `bench-last-gate` is a copied literal in the verdict
writer and reader (rename one side and the status row and commit-reuse protocol
die with no error); the git guard's wrapper search order is deliberately inlined
against the shared resolver lib with nothing detecting drift; and
`git.DefaultBranch`'s doc comment claims to mirror a shell function that no
longer exists.

## Solution

Collapse each duplicated fact to one source where sharing is safe, and gate the
one duplication that is deliberate: export the marker const for guards to
import, share a filename const between writer and reader, add a conformance
check that reds when the guard's inlined search order drifts from the shared
lib, and delete the stale mirror clause.

## User stories

1. As a kit maintainer, I want `bench guards` to match the pre-push marker via
   the exported adopt-package constant instead of its copied literal, so that a
   marker change cannot silently make guards misreport a managed hook.
   Line: claude-sonnet-5 / low. This is a mechanical constant export at a known
   seam whose behavior is already black-box covered by the AXI guards contract.

2. As a kit maintainer, I want the `bench-last-gate` cache filename to be one
   shared constant used by both the verdict writer (`gate.Record`) and reader
   (`status.GateVerdict`), so that renaming one side cannot silently break the
   gate-status and commit-reuse protocol.
   Line: claude-sonnet-5 / low. This is a mechanical constant introduction with
   the on-disk protocol already pinned by the runtime contract suite.

3. As a reviewer, I want a conformance check that goes red when the git guard's
   deliberately-inlined wrapper search order drifts from the shared resolver
   lib, so that the fail-closed inline stays safe without sourcing the lib.
   Line: claude-sonnet-5 / medium. The check shape is fully specified here, but
   it is gate logic, and the profile's cached routing gives gate and conformance
   logic mid effort because a wrong oracle is the worst class of kit bug.

4. As a reviewer, I want the drift check proven to bite — a canary fixture plus
   a bite-test entry per the existing registry convention — so that a check
   rotted into an always-pass fails the canary layer instead of lying green.
   Line: claude-sonnet-5 / medium. Same gate-logic effort row as story 3; the
   fixture and registry mechanics are established convention.

5. As a code reader, I want the stale "Go mirror of bench.sh's
   `default_branch`" clause deleted from `git.DefaultBranch`'s doc comment, so
   that the comment names only surfaces that exist.
   Line: claude-sonnet-5 / low. This is a prose-only deletion with no runtime
   surface, graded by review rather than the gate.

## Implementation decisions

- **Marker const (story 1).** Export the existing unexported marker constant in
  the adopt package in place (no move); the guards package imports it. Verified
  cycle-free: adopt imports only the git package and stdlib, and nothing in
  adopt imports guards. The const's doc comment gains guards as a named
  consumer, keeping its "one source" claim truthful.
- **Cache-filename const (story 2).** One exported filename constant in the git
  package **[default]** (the assessment's suggested home; both gate and status
  already import git, so no new edge in the import graph). Writer and reader
  compose it with their existing absolute-git-dir resolution; only the filename
  is shared, not a path helper.
- **Black-box tests keep their literals [default].** The runtime contract
  tests, the AXI guards contract, and the stophook test helper deliberately
  keep the literal filename and marker strings: they pin the on-disk protocol
  from outside the binary, so a future rename of either constant goes red in
  the contract suite instead of passing tautologically.
- **Drift check (story 3).** A new conformance check in the root-conformance
  suite, package-core-guard family **[default]**. It extracts the ordered
  resolver-candidate tokens (repo `.bench/bin` wrapper, kit `bin` wrapper,
  PATH fallback) from the guard shim's inlined resolver and from the shared
  resolver lib, and reds on any order mismatch. Anchor honesty: a missing
  file, an unfindable resolver function, or an absent token is itself red —
  never a vacuous pass — so renaming the resolver cannot amputate the check.
- **Sharing the resolver stays rejected.** The guard's inline is a settled
  fail-closed posture (sourcing the lib adds a fail-open mode); the drift
  check is the chosen mitigation, not a step toward sharing. The shim's
  inline-rationale comment currently ends by predicting a future collapse
  "in the slice that ports stop.sh" — stop.sh already sources the lib, so
  that sentence is stale; replace it with a pointer to the drift check
  **[default]**.
- **Comment deletion (story 5).** Remove only the mirror clause; the rest of
  the doc comment (one source read by both diff and status) remains true and
  stays.
- **Roadmap row.** FT50's row gains the spec path and flips its Next to the
  implement phase, keeping `bench status`'s spec cross-check honest
  **[default]**.

## Testing decisions

- A good test here exercises external behavior at an existing seam: guards
  output over a real linked repo, the gate cache file on disk, conformance
  diagnostics over a fixture tree. No test reaches into package internals.
- Stories 1–2 are structural collapses tied by the compiler; their behavior is
  already covered (AXI guards contract, runtime gate/status contracts). Only
  stories 3–4 get TDD: the bite test goes red against a drifted fixture before
  the check exists.
- Prior art: `internal/conformance/fixture_bite_test.go` (bite pattern),
  `canaryFixtureRegistry` (fixture registration), the AXI guards contract and
  runtime gate/status contract suites (existing black-box coverage).
- Gate: `.bench/gate.sh` (the project gate).

### Seam diagram

Drift check (stories 3–4), attached at the root-conformance suite:

    trigger: bench gate → conformance phase (root-conformance suite)
        │
        ▼
    guard shim (inlined resolver)  ──▶  [ resolver-order      ]  ──▶  diagnostics:
    shared resolver lib            ──▶  [   drift check       ]       red on drift or
                                                                      missing anchor
        ◀ tests attach here: bite test materializes a drifted canary
          fixture and asserts the diagnostic; an amputated fixture
          (resolver renamed) must also red

Gate-cache protocol (story 2), already covered — shown for the veto pass:

    trigger: gate run (Stop hook / bench gate) · reader: bench status
        │
        ▼
    verdict + tree hash  ──▶  [ gate.Record ]  ──▶  <git-dir>/bench-last-gate
                                                          │
                              [ status.GateVerdict ]  ◀───┘  ──▶  status row
        ◀ tests attach here: runtime contracts drive the built binary in a
          fixture repo and assert the literal cache path and status rows

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | a link-installed pre-push hook reports as managed; a foreign hook as unmanaged | AXI guards contract (built binary, throwaway linked repo) | already covered — existing managed / unmanaged / not-installed row assertions | if guards matched any literal other than the installed template's marker, the managed row would read unmanaged and the contract reds |
| 1 | guards compiles against the exported const | Go build/vet/test in the conformance core check | already covered — compile is the tie | renaming the const without updating guards fails the build |
| 2 | writer and reader agree on the cache filename on disk | runtime contract (built binary, fixture repo) | already covered — gate run asserts the literal cache path exists; status tests seed the literal path and assert the row | a filename drift on either side breaks the literal-path assertions |
| 3 | check reds when the guard's inlined order drifts from the lib | conformance bite test over a drifted canary fixture | new bite test, written red-first against the fixture | a check that never fires cannot produce the diagnostic the fixture is built to trip |
| 3 | check reds when an anchor is missing (resolver renamed or file absent) | conformance check unit test over an amputated fixture | new test, written red-first | a naive substring scan passes vacuously when the resolver is renamed; this row forces missing-anchor-is-red |
| 3 | check is green on the real tree | root-conformance suite inside the gate | already covered — the gate runs the full suite on the kit tree | a false-positive check turns the shipping gate red immediately |
| 4 | the fixture is registered and its family materializes | canary fixture registry + family check | already covered — unregistered fixtures fail materialization; empty families red the family check | the bite test cannot run against an unregistered fixture |
| 5 | stale mirror clause removed | none | not TDD-able — prose-only deletion, graded by review | no runtime surface exists to observe |

### Edge inventory

Walked against the profile's shell-CLI hostile-input checklist. The only new
input-consuming surface is the drift check, which reads two kit-owned files by
root-relative path inside the conformance process — no shell, no user input,
no persistent state.

- absent file vs present-but-empty — coverage row (missing-anchor red covers
  both: no file and no token are the same red).
- hand-edited file lacking a trailing newline — **Won't handle:** token-order
  extraction is substring-based and newline-insensitive.
- paths with spaces/globs, control bytes, unquoted multi-word args — **Won't
  handle:** the check takes no user input; both paths are fixed kit-relative
  constants.
- required tool missing from PATH / symlink invocation / cwd below root —
  **Won't handle:** the check runs in-process against the conformance root the
  harness resolves; no external tools invoked.
- SIGINT mid-loop / re-run idempotency — **Won't handle:** read-only check, no
  state to leak; idempotent by construction.
- invocation through every shipped surface — **Won't handle:** the check's
  only surface is the conformance suite, which every gate entry path already
  routes through.
- amputated-caller guard for stories 1–2 — the compiler is the guard: an
  orphaned const or a missed call site fails the build (coverage rows above).

## Out of scope

Empty — every assessment finding adjacent to this batch is already its own
roadmap row (FT51–FT53), and the stop.sh-port slice the old shim comment
predicted has no remaining work: stop.sh already sources the shared lib.
