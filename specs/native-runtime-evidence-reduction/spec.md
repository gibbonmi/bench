# Native runtime evidence reduction

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-28. Three forks closed. The `artifacts` job builds one generation. The job publishes one upload. The `native-proof` matrix drops the Darwin targets, and the `smoke` matrix keeps them.

Verification log: 2 iteration(s) to accept — the first round blocked on three findings. Two consumers of the shipped-target count sat outside the stories, and the first ticket was not independently green.

## Problem

The `native runtime` workflow costs more than the evidence it returns.

The `artifacts` job runs four complete artifact builds. It builds two isolated generations, and `scripts/build-artifacts.sh` already clones its own second source and builds again inside each generation. Each build produces four platform binaries. The job then publishes four uploads, and the `evidence` job downloads six.

The `native-proof` job also runs on two macOS runners. Those runners have never returned a green proof. The Darwin branch of `scripts/native-proof.sh` asserts that `nm -a` reports no symbols. A loadable Mach-O cannot satisfy that assertion, because dyld resolves every libSystem call through the undefined-import table.

## Solution

The workflow keeps the claims it can pay for, and it drops the rest.

The `artifacts` job builds one generation and publishes one upload. The byte-for-byte artifact comparison survives, because `scripts/build-artifacts.sh` performs it inside that single build and writes `dist/reproducibility.json`. The workflow loses only its cross-checkout comparison of the final release evidence.

The release plan states which targets carry a native proof. The `native-proof` matrix reads that proven list, so no macOS runner starts. The `smoke` matrix keeps the full shipped list, so the macOS binary still runs on macOS. Every consumer that counts proofs counts against the proven list.

## User stories

### Group A — the artifacts job is shorter

Line: `opus` / medium. The change edits two workflows and their handoff paths, so it needs care rather than depth.

1. As a release engineer, I want the `artifacts` job to build one generation, so that the job costs two complete builds instead of four.
2. As a release engineer, I want the `artifacts` job to publish one upload, so that the job hands the next job a single object.
3. As a release engineer, I want `dist/reproducibility.json` to travel inside that upload, so that the evidence job still reads the comparison record.
4. As a release engineer, I want the byte-for-byte artifact comparison kept, so that the published reproducibility claim stays true.
5. As a release engineer, I want the `evidence` job to download one generation, so that it does not clone or finalize a second source.
6. As a release engineer, I want `dist/workflow-reproducibility.json` removed from every reference, so that no consumer reads a file that nothing writes.
7. As a release engineer, I want the `evidence` job to keep its preflight verification, so that the finalized evidence is still graded.
8. As a release engineer, I want the `smoke` job to download from the new upload root, so that its run finds the tarballs.
9. As a release engineer, I want the tag release jobs to download from the new upload root, so that publication still finds the tarballs.
26. As a release engineer, I want the `native-proof` job to build one generation, so that it stops reading an upload that nothing publishes.

### Group B — shipped targets and proven targets are separate facts

Line: `opus` / medium. The change adds one field and derives every count from it.

10. As a release engineer, I want `scripts/release-plan.json` to state per target whether that target carries a native proof, so that one file holds both facts.
11. As a release engineer, I want the `native-proof` matrix derived from the proven targets, so that no macOS runner starts.
12. As a release engineer, I want the `smoke` matrix derived from the shipped targets, so that the shipped macOS binary still runs on macOS.
13. As a release engineer, I want `scripts/aggregate-native-proofs.sh` to require exactly the proven proof set, so that a missing or extra proof stays red.
14. As a release engineer, I want `scripts/native-proof.sh` to refuse an unproven target, so that a stray call cannot mint a proof the aggregator would reject.
15. As a release engineer, I want finalization to count native proofs against the proven targets, so that it demands no proof that never runs.
16. As a maintainer, I want the Darwin branch of `scripts/native-proof.sh` removed, so that no unreachable and unverifiable assertion stays in the tree.
17. As a maintainer, I want the proof script to refuse a target it cannot verify, so that no plan edit mints an unverified proof.
18. As a release consumer, I want `scripts/verify-release-artifact.mjs` to count proofs against the proven targets, so that the offline verification of a published artifact still passes.
19. As a reviewer, I want the release evidence probe to synthesize proofs for the proven targets, so that it matches a real run.

### Group C — the gate grades the new shape

Line: `opus` / high. The group edits the oracle, and every changed check needs a demonstrated red.

20. As a reviewer, I want the workflow conformance check to grade the one-generation shape, so that a return to two generations fails the gate.
21. As a reviewer, I want the conformance check to grade the proof matrix source, so that reading the shipped matrix fails the gate.
22. As a reviewer, I want each changed conformance check to carry a mutation that turns it red, so that the check provably bites.
23. As a reviewer, I want the canary fixtures re-anchored on lines that survive, so that the guard suite stays green and still bites.
24. As a teammate with no history, I want the release docs to state the current claim, so that I do not read a retired promise.
25. As a reviewer, I want the two retired claims named in the docs, so that the reduction reads as a decision.

## Implementation decisions

**The plan carries one new field.** Each target row in `scripts/release-plan.json` gains `native_proof`, a required boolean. `readReleasePlan` in `release-plan.mjs` rejects a target that omits it or gives a non-boolean. It also rejects a plan whose proven set is empty, because an empty set makes the aggregator vacuous.

**The Go plan struct accepts the field in the same change.** `internal/releaseevidence` decodes `normalized-json` with unknown fields disallowed. The new field therefore reddens that package the moment the plan carries it, so the plan change and the Go struct field land together. The Go consumer reads the field later.

**Two derived views, one projection.** `release-plan.mjs` keeps `targets`, `target`, and `matrix-json` over the shipped list. It adds `proof-targets`, `proof-matrix-json`, and `proof-target` over the proven list. Both matrix commands project the same five transport fields, so `native_proof` never reaches a GitHub matrix variable. One internal helper produces both views from one filter, so the projection is written once.

**The proof script refuses an unproven target.** `scripts/native-proof.sh` calls `proof-target` in place of `target`. That command exits non-zero for an unproven target, and the script's existing failure message covers the case.

**The Darwin branch goes, and the script fails closed.** With no proven Darwin target the branch is unreachable. The script keeps the Linux branch whole.

The script refuses any proven target whose operating system has no platform predicate. Without that refusal it would emit a green strip status for a target nothing verified. A future Darwin proof is therefore a plan-data change plus a re-authored strip assertion, written against real macOS `nm -a` output.

**Every proof count derives from the proven list.** Three consumers count proofs today, and all three move together. `internal/releaseevidence` finalization compares against the proven-target count. `scripts/verify-release-artifact.mjs` compares its `native_proofs` length against the proven-target count. The release evidence probe in `internal/conformance` synthesizes a proof per proven target. Each keeps its per-operating-system clauses unchanged.

**The workflow hands over one object.** The `artifacts` job uploads one artifact whose `path` names `dist/artifacts` and `dist/reproducibility.json`. The upload action resolves `dist` as the common root, so the archive root moves up one level.

Every consumer therefore downloads into `dist`, not into `dist/artifacts`. There are five. Three sit in `native-runtime.yml`: the `native-proof`, `evidence`, and `smoke` jobs. Two sit in `release.yml`: the `authorize` and `publish` jobs.

**The `native-proof` job loses its second generation too.** It clones a second source, runs `scripts/native-proof.sh` twice, and publishes a second proof upload. Each of those reads or writes a retired artifact. So the clone, the second download, the second proof run, and the second upload all go.

**`workflow-reproducibility.json` is retired.** No Go evidence reads it. Its only producer is the `evidence` job's `compare-artifacts.sh` call, and that call goes with the second generation. `scripts/compare-artifacts.sh` itself stays, because `build-artifacts.sh` still calls it. `reproducibility.builds` stays at 2, because the inner build still performs two builds.

## Testing decisions

A good test here reads the workflow text and the plan data, and it grades the shape a release run would take. The build cannot run a GitHub matrix, so the workflow checks stay text checks, and the canary guard supplies their red.

Three seams already exist. `internal/conformance` owns the workflow text checks, their mutation proofs, and the release evidence probe. `tests/canary/package-core-guard` owns the fixture-driven reds. `internal/releaseevidence` owns the proof-count behavior.

Two seams are new. The plan reader receives unit tests beside the existing plan tests. The aggregator and the proof script receive an execution test that runs each script against a temporary proof directory. That test hosts in `internal/conformance`, which is already an ownership fence and already executes release scripts.

The gate seam is `bench gate`, which runs the conformance package, the canary guard, and the release evidence package.

### Seam diagram

    trigger: bench gate
        │
        ▼
    workflow text  ──▶  [ checkNativeRuntimeWorkflow ]  ──▶  diagnostics
                            ◀ tests attach here: a fixture mutates one line, the check names it

    plan JSON      ──▶  [ readReleasePlan / proof views ]  ──▶  proven target list
                            ◀ tests attach here: a unit test feeds a plan and reads the view

    proof set      ──▶  [ inspectNativeProofs / verifier ]  ──▶  release evidence
                            ◀ tests attach here: a probe writes proofs and reads the verdict

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| A1 | 1 | The `artifacts` job contains no second-source clone | `internal/conformance` workflow check | A restored clone returns the four-build cost, and the check names it |
| A2 | 2 | The `artifacts` job contains exactly one upload step | `internal/conformance` workflow check | A second upload step returns the split handoff this spec removes |
| A3 | 3 | The single upload names `dist/reproducibility.json` | `tests/canary` fixture | A dropped record leaves finalization with no comparison to read |
| A4 | 4 | `build-artifacts.sh` still writes `dist/reproducibility.json` | `internal/releaseevidence` probe | A build that skips the inner comparison makes the published claim false |
| A5 | 5 | The `evidence` job contains no second-source clone | `internal/conformance` workflow check | A restored clone re-adds the finalization this spec retires |
| A6 | 6 | No tracked file under `.github`, `scripts`, `internal`, `tests`, or `docs` names `workflow-reproducibility.json` | `internal/conformance` tracked-file sweep | A surviving reference points at a file nothing writes |
| A7 | 7 | The `evidence` job still runs `release-preflight.sh --mode verify` | `internal/conformance` workflow check | A dropped preflight ships evidence that nothing graded |
| A8 | 8 | The `smoke` job downloads the one artifact into `dist` | `internal/conformance` workflow check | A download into `dist/artifacts` nests the tarballs one level too deep |
| A9 | 9 | The `authorize` and `publish` jobs download the one artifact into `dist` | `internal/conformance` workflow check | A stale path makes every tag release fail after the upload root moves |
| A10 | 26 | The `native-proof` job contains no second-source clone | `internal/conformance` workflow check | A second generation reads an upload that nothing publishes |
| A11 | 26 | The `native-proof` job downloads the one artifact into `dist` | `internal/conformance` workflow check | A download into `dist/artifacts` nests the tarballs one level too deep |
| B1 | 10 | `readReleasePlan` rejects a target with no boolean `native_proof` field | plan reader unit test | A silent default hides which targets carry a proof |
| B2 | 10 | `readReleasePlan` rejects a plan whose proven set is empty | plan reader unit test | An empty set makes the aggregator pass with no proof at all |
| B3 | 11 | `proof-matrix-json` omits every unproven target | plan reader unit test | A leaked Darwin row starts the macOS runner this spec removes |
| B4 | 11 | The `native-proof` job reads `proof-matrix-json` | `internal/conformance` workflow check | A job reading the shipped matrix starts the macOS runner anyway |
| B5 | 12 | The `smoke` job reads `matrix-json` | `internal/conformance` workflow check | A smoke job on the proven matrix stops running the macOS binary |
| B6 | 13 | The aggregator fails when a proven proof file is absent | script execution test | A missing proof would otherwise publish an unproven target |
| B7 | 13 | The aggregator fails when the directory holds an unproven proof file | script execution test | An extra file signals a matrix and plan that disagree |
| B8 | 14 | `native-proof.sh` exits non-zero for an unproven target | script execution test | A stray call would mint a proof the aggregator later rejects |
| B9 | 15 | Finalization succeeds with proofs for the proven targets only | `internal/releaseevidence` probe | A count against shipped targets blocks every release |
| B10 | 15 | Finalization fails when one proven target has no proof | `internal/releaseevidence` probe | A relaxed count would accept a release with no Linux proof |
| B11 | 16 | `native-proof.sh` contains no Darwin branch | script execution test | A surviving branch keeps an assertion no run can pass |
| B12 | 17 | The proof script refuses a proven target whose operating system has no platform predicate | script execution test | An unverified target would otherwise mint a green strip status |
| B13 | 18 | `verify-release-artifact.mjs` accepts an index holding proofs for the proven targets only | clause extraction, reviewer-graded | A count against shipped targets fails every offline verification |
| B14 | 19 | The release evidence probe writes a proof per proven target | `internal/conformance` probe | A probe on the shipped list grades a set no run produces |
| C1 | 20, 21 | Each changed conformance check names its own diagnostic | `internal/conformance` mutation test | A check with no red is advertisement, not enforcement |
| C2 | 22 | Every mutation in the changed checks turns the check red | `internal/conformance` mutation test | An anchor that no longer occurs makes the mutation a no-op |
| C3 | 23 | Every canary fixture under `package-core-guard` still bites | canary fixture-bite test | A stale anchor silently retires a guard |
| C4 | 24, 25 | The release docs state the current claim and name both retired claims | reviewer reading | A stale doc tells a cold reader we prove something we do not |

Two rows are graded by the reviewer, not by the gate.

Row C4 covers prose accuracy. It has no red-capable seam here, and a keyword sweep would bind the docs to one wording.

Row B13 covers the consumer-side verifier. Nothing the gate runs executes `scripts/verify-release-artifact.mjs`. Its only callers are the two smoke scripts, and those run in the workflow.

The row's evidence is a clause extraction. The build evaluated the live count expression against a real release index. It went green at the proven count and red at the shipped count. The change under it is one token, and the per-proof predicate is untouched. A real execution seam for that verifier is out of scope and parked.

### Edge inventory

- A plan with every target unproven. `readReleasePlan` rejects it (row B2).
- A proof directory holding a file for an unproven target. The aggregator rejects it (row B7).
- A proof directory holding a stray non-proof file. The existing exact-set check rejects it, and that check now compares against the proven list.
- A `native-proof.sh` call naming a target absent from the plan. The existing matrix check rejects it.
- A release index whose proof count matches the shipped list. The verifier rejects it (row B13).
- An upload whose common root resolves higher than `dist`. Rows A8 and A9 pin every download path against that root.

**Won't handle:** a corrected Darwin strip assertion. The branch is unreachable while no Darwin target is proven. A re-proved Darwin target re-authors that assertion against real macOS `nm -a` output.

**Won't handle:** the cross-checkout comparison of final release evidence — it retires with the second generation, and the in-script artifact comparison remains for every in-scope consumer.

**Won't handle:** a flag that restores the second generation — no in-scope caller wants it, and a dormant second path would drift.

**Won't handle:** reducing the `smoke` matrix — the reviewer kept it so that one job still runs the shipped macOS binary on macOS.

## Ownership fences

- `.github/workflows/native-runtime.yml`
- `.github/workflows/release.yml`
- `scripts/release-plan.json`
- `scripts/release-plan.mjs`
- `scripts/native-proof.sh`
- `scripts/aggregate-native-proofs.sh`
- `scripts/verify-release-artifact.mjs`
- `internal/conformance/`
- `internal/releaseevidence/`
- `tests/canary/package-core-guard/`
- `docs/release-runbook.md`
- `docs/field-guide.html`
- `projects/benchkit.md`
- `ROADMAP.md`
- `specs/native-runtime-evidence-reduction/`

## Out of scope

- **Restore a native proof for the Darwin targets.** It needs a macOS strip assertion authored against real `nm -a` output and one green macOS run. 6 edits, 3 gate runs.
- **Replace the retired cross-checkout evidence comparison.** A cheaper form could compare final evidence between two runs rather than two checkouts. 8 edits, 4 gate runs.
- **Reduce the `smoke` matrix.** 4 edits, 2 gate runs.

## Further notes

The reduction retires two published claims, and the docs must say so plainly. Bench no longer proves that two independent checkouts finalize identical release evidence. Bench no longer proves the shipped macOS binaries on a macOS runner through `native-proof`; the `smoke` job remains the macOS execution evidence.

This spec also closes the open Darwin strip-assertion defect. The assertion leaves the tree with its branch, so no separate fix is needed.
