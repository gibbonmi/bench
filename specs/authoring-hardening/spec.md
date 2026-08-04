# Authoring hardening

Status: staged

Decision source: reviewer-confirmed conversation 2026-08-04 — six proposals hardening `craft-tickets` and `craft-spec` from the recovery-discard build's evidence; the reviewer chose the fence/red-mutation agreement check as the enforcement-first item, sized the ticket-breadth rule as a split signal rather than a cap, confirmed retiring the ticket grammar's `Assumptions:` field (including the `decisions/parallel-session-landings.md` sealed-field-list edit) after a mid-tier consult, and directed one bundled batch with a ticket per proposal.

## Problem

An eight-ticket build finished with every ticket green and still cost three avoidable rounds, because the ticket and spec grammars let four defect shapes through at authoring time:

- A ticket's red-mutation command can name a package its ownership fence cannot write. The delegate obeys the row, breaches the fence, and the checkpoint refusal costs a full repair round — the contradiction was visible in the ticket file the whole time, and nothing checks the two fields against each other.
- A `-run` selector can match zero tests. `go test` still exits zero with `ok`, so a red-mutation probe can report evidence from a command that never executed a test.
- Every mutation the tickets described was control-flow. Reverting an authorizing constant or dropping a field from a hash preimage is invisible to that entire class, so two safety properties shipped correct but unasserted until coordinator probes caught them.
- Nothing at spec or ticket time flags breadth. The wide end carried a risk signal, but narrower tickets also required repairs, so narrowness proved nothing; the spec itself bundled two capabilities with no shared seam, so the narrower one could not ship on its own gate.

Separately, the ticket grammar's `Assumptions:` field pools two clause kinds that are not assumptions — tree-checkable preconditions redundant with `Blocked by:`, and restatements of the standing verify rule — and its machinery digests a line whose per-clause structure the comma-splitting parser never actually produces, duplicating a fact `TicketDigest` already seals.

## Solution

Ticket authoring gains one enforced agreement and three authored disciplines; spec authoring gains one split signal; the `Assumptions:` field leaves the grammar and its machinery. A contradictory ticket is refused before a delegate is ever charged with it, and the remaining hazards are named in the skills that govern the artifacts, at the exact decision points where the recovery-discard build went wrong.

## User stories

1. **A ticket whose red-mutation command contradicts its fence is refused before it costs a round.** For every `go test` package argument in a red-mutation row's operation-sequence cell, the row must have somewhere its red can live: a fence entry that is a package prefix of the probed package, a fenced `_test.go` file whose directory is exactly the probed package, or an owner cell naming a backticked `_test.go` path that already exists in that package — the pre-existing-owner escape for regression rows that author nothing. A ticket failing the rule is refused at `bench spec build assign`, a staged spec carrying one turns the gate's conformance phase red, and a canary fixture proves the check bites. `craft-tickets` documents the owner-cell convention beside the field it extends.
   `Line: opus / medium.` Gate and conformance logic routes mid effort per the profile; the model half is uncached and flagged for confirmation at sign-off.

2. **A ticket author sees breadth as a split signal.** `craft-tickets` names a fence spanning more than two packages as a split candidate demanding one line of justification in the breakdown — a smell, not a cap, because an honest tracer (a verb, its grammar row, its launcher line) can legitimately cross three. Breadth counts the distinct directories the fence entries resolve to, a file entry counting its parent, so a file-scoped fence cannot dodge the count.
   `Line: fable / high.` Skill prose rides the leverage override: it compounds through every session that loads it.

3. **An authorizing value gets at least one input mutation.** `craft-tickets`' red-mutation guidance requires that a ticket touching a value that authorizes an action — a fingerprint, digest, or token — carry at least one mutation row that changes the value's inputs (revert the constant, drop a hashed field) rather than its code path.
   `Line: fable / high.` Same leverage override; the rule is one paragraph beside the existing mutate-the-subject rule.

4. **A probe that matches zero tests is a failed probe.** `craft-tickets` requires red-mutation evidence to carry the matched-test count per probe, makes zero a hard failure of the operation sequence, and tells authors to verify `-run` selectors against real top-level test names — subtests of one registered function do not match a bare selector.
   `Line: fable / high.` Same leverage override; this is the hung-run hazard's silent sibling.

5. **A spec author sees a disjoint story partition as a split signal.** `craft-spec`'s story-sizing guidance gains a partition check: when stories partition into disjoint package sets connected by no shared seam or contract, surface the split to the reviewer at spec time. A deliberate bundle can still be chosen — but chosen, not defaulted.
   `Line: fable / high.` Same leverage override, one paragraph in the sizing section.

6. **The spec-build machinery stops carrying the `Assumptions:` field.** `ParseTicket` ignores a legacy `Assumptions:` line; the ticket, assignment, checkpoint, and integration structures drop the field, its digests, and its comparisons. Persisted records written before retirement still load, and a ticket file carrying the legacy line still assigns, checkpoints, and integrates with the line ignored.
   `Line: opus / medium.` Mechanical removal graded by compilation and the existing lifecycle tests.

7. **The grammar surface and the records stop advertising the field.** The `craft-tickets` template, field bullet, and taught example drop the line — the `example-agreement` check holds example and parser in agreement through the change — and the two sealed-field lists in `decisions/parallel-session-landings.md` are edited to match (reviewer-approved ADR edit, 2026-08-04). Genuinely unverifiable-at-authoring-time claims, which the evidence build produced zero of, live in the ticket's What-to-build prose.
   `Line: fable / high.` Guidance and ADR prose; the leverage override applies.

## Implementation decisions

**One validation function, two consumers.** The fence/red-mutation agreement rule lives in `internal/specbuild` beside `ParseTicket`, which already owns both fields. `assign` refuses on it, and the conformance sweep grades staged specs through the same exported function — a second derivation of the rule would be the duplicated-knowledge defect the code standard forbids. The extraction takes the package arguments of `go test` invocations in the red-mutation section; other command tokens are ignored.

**The agreement rule is ownable-red, decided per row.** Extraction reads only the operation-sequence cell of each red-mutation row — a `go test` in the mutation column or in ticket prose is not a probe — and takes every package argument of every `go test` invocation it finds there. A probed package passes when the fence holds a segment-aware package prefix of it, or a fenced `_test.go` file whose directory is exactly the probed package (a test file in a subdirectory is a different Go package and does not count), or the row's owner cell names a backticked `_test.go` path that exists in the tree with its directory exactly the probed package. The recovery-discard evidence fixed both refusal shapes: a fence naming only a non-test file inside the probed package (the `internal/spec/build.go` case) is refused, while a fence naming the package's own `_test.go` file (the `runtime_worktree_test.go` case) passes. Segment-aware means `internal/spec` neither covers nor is covered by `internal/specbuild`.

**Wildcards are refused.** A `./...` package argument in a red-mutation row is unscopeable — no fence can own it — and the focused-probe discipline already forbids it in spirit; the rule makes it a refusal.

**A staged violation found at landing is repaired in the sweep's own change.** If a spec is still staged when the sweep lands and carries a violating ticket — today that is `specs/recovery-discard/tickets/add-spec-build-reclaim.md`'s RM5, should that run end abandoned rather than promoted — the sweep's ticket repairs the artifact in the same green change rather than landing a red gate on a pre-existing file.

**Retirement moves every anchor in one green change.** The grammar fact is advertised in five places that red independently: the `craft-tickets` template and taught example, the `example-agreement` suite's independently-authored expectation literals, the template anchors in the fixture-bite and docs-workflow conformance suites, and the sealed-field lists in `decisions/parallel-session-landings.md`. The surface story owns all five, because a red gate between any pair is not a landable intermediate state. The machinery half also reaches `internal/contract/runtime`'s spec-build receipts, which construct assumption digests today.

**The sweep grades staged specs only.** An implemented spec's tickets are history awaiting promote-then-delete retirement, and grading them would red the gate on artifacts no build will ever charge again. Status is read from the spec's own `Status:` line.

**Retirement keeps legacy artifacts readable.** `ParseTicket` skips a legacy `Assumptions:` line rather than refusing it, so staged tickets written under the old grammar still parse; persisted assignment records and checkpoint receipts carrying assumption data load under the new structures because unknown JSON fields are ignored by construction. No migration pass, no version bump.

**The prose lands inside the sections that already own the decision points.** Story 2 in `craft-tickets`' breakdown-drafting step, stories 3 and 4 in its red-mutations field bullet, story 5 in `craft-spec`'s story-sizing section. No new sections, no new skills.

**Sequencing: build starts after the recovery-discard run promotes.** Stories 1 and 6 edit `internal/specbuild`, which that run's candidate rewrites, and its active records carry assumption digests the retirement removes. The spec stages now; implementation waits for the promotion.

**This spec is itself a reviewer-chosen bundle.** By story 5's own partition rule, stories 1 and 6 (`internal/specbuild`), stories 2–4 (`craft-tickets`), story 5 (`craft-spec`), and story 7 (grammar surface and ADR) would surface as a split candidate. The reviewer directed one batch with a ticket per proposal (2026-08-04); recording the choice here is the rule working as intended.

## Testing decisions

- A good test drives the real parser over a real ticket file: write the ticket markdown to a temp spec directory, run `ParseTicket`, and assert the refusal or acceptance — never a hand-built `Ticket` struct standing in for the parse. Prior art: the `example-agreement` suite, which parses authored markdown through the real function.
- Seams receiving tests: the validation beside `ParseTicket` in `internal/specbuild` (prior art: the assign and checkpoint fixtures), the conformance sweep in `internal/conformance` (prior art: the docs-currency and example-agreement checks), and the canary fixture family in `tests/canary/` (prior art: the existing conformance-owned fixtures).
- The gate observes this through the conformance phase (the sweep and `example-agreement`), the canary phase (the bite fixture), and the specbuild package tests; no new gate phase.

### Seam diagram

    trigger: `bench spec build assign` charges a ticket, or the gate's conformance phase sweeps staged specs
        │
        ▼
    ticket file  ──▶  [ ParseTicket + fence/mutation agreement validation ]  ──▶  refusal, or a validated Ticket
        │                                                                              │
        │                                                                              ▼
        └──────────▶  [ conformance sweep over specs with Status: staged ]  ──▶  gate verdict
                          ◀ tests attach here: real ticket markdown written to a
                            temp spec dir, parsed by the real function; the canary
                            drives the sweep against a deliberately-broken fixture

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | assign refuses a ticket whose red-mutation command names a package the fence cannot own a test in | internal/specbuild | not yet observed — no validation exists to fail | the recovery-discard RM5 shape charges a delegate today and costs a repair round |
| 1 | a ticket whose fence names the probed package's own `_test.go` file stays assignable | internal/specbuild | cannot start red — assign accepts this shape today; regression row pinning the legitimate shape | an over-strict rule would refuse the legitimate real-producer ticket shape and block honest file-fenced tickets |
| 1 | a regression row whose owner cell names an existing `_test.go` in the probed package passes without fence coverage | internal/specbuild | not yet observed | without the escape the rule reds honest production-only tickets that rely on pre-existing coverage |
| 1 | a `./...` package argument in a red-mutation row is refused | internal/specbuild | not yet observed | a wildcard probe is unscopeable, so the agreement rule cannot bind it |
| 1 | the conformance sweep reds a staged spec carrying a contradictory ticket and skips implemented specs | internal/conformance | not yet observed | a contradiction that waits for assign wastes the build's round; grading retired history would red the gate on dead artifacts |
| 1 | the canary fixture proves the sweep bites with its targeted substring | tests/canary | not yet observed | a check rotted into always-pass is invisible without its fixture |
| edge of 1 | a probe of `./internal/spec` is not covered by a fence entry under `internal/specbuild`, and the reverse | internal/specbuild | not yet observed | naive string-prefix matching passes both, which is the cheapest wrong containment |
| edge of 1 | a fenced `_test.go` in a subdirectory of the probed package does not satisfy the rule | internal/specbuild | not yet observed | path containment accepts a file that is a different Go package, which is the cheapest wrong depth reading |
| edge of 1 | `go test ./a ./b` extracts both package arguments | internal/specbuild | not yet observed | first-token extraction passes every single-package row and never sees the second package |
| edge of 1 | a `go test` mention outside the operation-sequence cell is not extracted | internal/specbuild | not yet observed | extracting from prose or the mutation column would red spuriously and teach authors to strip context |
| 2 | craft-tickets names a fence spanning more than two packages as a split candidate with one line of justification | .agents/skills prose | not TDD-able — guidance prose graded at review | both four-directory tickets drew repairs, and narrower tickets did too, so width is only a wide-end risk signal and narrowness proves nothing |
| 3 | craft-tickets requires one input mutation for a ticket touching an authorizing value | .agents/skills prose | not TDD-able — guidance prose graded at review | every mutation in the evidence build was control-flow, and both unasserted properties were input-shaped |
| 4 | craft-tickets makes zero matched tests a hard probe failure and requires the matched count in evidence | .agents/skills prose | not TDD-able — guidance prose graded at review | `-run Recovery` matched nothing and exited `ok`; only the delegate's care caught it |
| 5 | craft-spec surfaces a disjoint story partition as a split signal at spec time | .agents/skills prose | not TDD-able — guidance prose graded at review | the evidence spec bundled two capabilities sharing only a theme, and the narrower one could not ship alone |
| 6 | a ticket carrying a legacy `Assumptions:` line parses, assigns, checkpoints, and integrates with the line ignored | internal/specbuild | cannot start red — the lifecycle already succeeds on such tickets; asserted as a regression row pinning legacy tolerance | retirement that refused legacy lines would strand every staged ticket written under the old grammar |
| 6 | a pre-retirement persisted record carrying assumption data reloads and its lifecycle continues | internal/specbuild | cannot start red — unknown JSON fields are ignored by construction; asserted as a regression row across a fresh process | the process-boundary class is where silent record breakage would land |
| 6 | the runtime spec-build receipts no longer construct assumption digests | internal/contract/runtime | already enforced — the machinery removal reds the runtime fixture's receipt expectations until they move in the same change | a receipt shape asserted against removed machinery is unlandable, so the fixture forces the paired edit |
| 7 | the template, the taught example, and every anchored expectation drop the field in one change | internal/conformance | already enforced — the template edit alone reds the example-agreement, fixture-bite, and docs-workflow anchors until their independently-authored expectations move with it | a red between anchors is unlandable, which is the paired-change discipline those checks exist to force |
| 7 | the sealed-field lists in `decisions/parallel-session-landings.md` no longer name the field | decisions/ | not TDD-able — ADR prose; reviewer-approved edit graded at review | a stale sealed-field list re-teaches the retired grammar to the next parallel session |

**Degenerate implementation.** The cheapest wrong story 1 validates containment by naive string prefix, so `internal/spec` covers `internal/specbuild` and both directions pass; the segment-boundary edge row goes red on it. The cheapest wrong sweep hardcodes a green verdict; the canary fixture row goes red on it. The cheapest wrong story 6 deletes the prose while leaving the machinery — the gate cannot see that drift because the machinery is green either way, so the review axis grades it against this spec's decision, stated here so the review knows to look.

**Composition degenerate.** Stories 1's two consumers span `internal/specbuild` and `internal/conformance`: the composition degenerate is a sweep that carries its own subtly different copy of the agreement rule, green against fixtures that happen to agree. The canary fixture drives the sweep through the real gate entry, and the review grades the single-function decision; a second derivation is the duplicated-knowledge defect by name.

### Edge inventory

- **Error path** — the refusal rows above (contradictory package, wildcard, sweep red).
- **Empty/absent input** — a ticket with no red-mutation section or no `go test` argument is vacuously green for story 1 (the rule binds commands that exist); a spec directory with no tickets sweeps green. Both cannot start red and are asserted as regression cases inside the story-1 rows' suites. An owner cell naming a `_test.go` path that does *not* exist falls back to the fence rule rather than crashing, same suites.
- **Boundary values** — the segment-aware containment row and the subdirectory-depth row (edge of 1); a fence entry exactly equal to the probed package is containment, covered by the first story-1 row's suite.
- **Malformed input** — the `./...` row, the multi-package row, and the outside-the-cell row above; a `go test` invocation with flags before the package argument (`go test -count=1 ./pkg`) extracts the package, asserted in the story-1 suites.
- **Interrupted or partial state** — **Won't handle:** the validation and the sweep are pure reads with no state to interrupt; nothing mutates.
- **Re-run idempotency** — **Won't handle:** read-only checks over tracked files; same input, same verdict, no persisted state.
- **Process-boundary lifecycle** — story 6's persisted-record row drives a fresh reload rather than reusing in-memory structures.
- **Hostile environment** — fence entries and package arguments with spaces or glob characters pass through the existing backtick-trimmed string comparison with no new file I/O introduced; control bytes in ticket text stay owned by `ParseTicket`'s existing read posture, which this spec does not widen. **Won't handle** beyond that existing posture: the sweep introduces no new path resolution beyond the spec directory walk the docs checks already perform.

## Out of scope

- **Static verification that `-run` selectors match real test functions.** Selectors in a fresh ticket reference tests the build has not written yet, so the check cannot run at authoring time; a post-build sweep is a separate capability. Estimate: 8 edits, 4 gate runs.
- **Rebinding the "NEVER assume, always verify" rule.** Owned by FT187 and its staged spec; the mid-tier consult (2026-08-04) confirmed no ordering dependency with the retirement in either direction.
- **Sweeping implemented specs' tickets.** They retire with their spec under promote-then-delete; grading them adds gate noise for artifacts no build will charge. Estimate if reconsidered: 2 edits, 1 gate run.
- **The learnings drain.** The five open journal entries behind this spec stay open; the reviewer deferred `/bench-what-next` explicitly (2026-08-04).
