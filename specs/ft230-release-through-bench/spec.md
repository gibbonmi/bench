# FT230 — the release workflow publishes only through `bench release submit`

Status: staged

Decision source: named reviewed artifact — `roadmap/FT230.md` (drained from the 2026-08 capability audit, item A7 in `docs/audits/2026-08-bench-capability/results-fable-high/action-items.yaml`), plus the reviewer's 2026-08-20 adapter-seam confirmation (`--adapter` flag, default `fixture`).

Verification log: 2 iteration(s) to accept — round 1 returned six blocking findings (the release workflow's step-name byte contracts in `native_workflow_test.go` and the `preflight-publish-order-bypassed` canary; no live ordering test in `internal/publication`; the missing `NODE_AUTH_TOKEN` on the submit step; unscoped download/upload anchors; the runbook's unattended-first-publication contradiction) plus four advisories (env twin removed, the credential-argv row dropped as not red-capable, hop-2 wording, story 19 repointed); all folded, ordering assertion moved into T1 so it exists before T2 retires the step-name check.

## Problem

`.github/workflows/release.yml` publishes with raw `npm publish` in two steps,
bypassing the resumable, digest-verified state machine the release runbook
requires. A partial publication cannot be resumed or rolled back through Bench.
Worse, the CLI cannot currently reach the real registry at all: `bench release
submit/promote/rollback` construct the fixture registry unconditionally, so the
`NPMCLIRegistry` adapter exists but is dead code, and the conformance check
`checkReleaseWorkflow` asserts the presence of the raw `npm publish` anchor —
the enforcement requires the defect. The raw steps also skip the approved-set
digest verification entirely: the publish job never downloads the preflight
evidence it would need to verify anything.

## Solution

The publish job of the release workflow invokes `bench release submit` and
nothing else. The CLI gains an explicit registry-adapter selection — `--adapter
npm|fixture`, default `fixture` — so a real publish is an explicit opt-in and
every existing hermetic caller is unchanged.
The npm adapter carries the publish arguments the raw steps carried (`--access
public`, and provenance on an explicit `--provenance` opt-in). The conformance
check flips: a raw `npm publish` in the release workflow is red, and the
`bench release submit` invocation anchor is required. Promotion stays
operator-run with the reviewer present, exactly as the runbook says; CI never
promotes. The runbook's presence rule is reconciled by amendment: only the
reviewer cuts and pushes a release tag, so the tag push is the attended act
and the CI submit is its mechanical arm (reviewer decision, 2026-08-20).

## User stories

### Group A — the CLI reaches real npm only on explicit opt-in

Line: opus / high. The publication trust path is release-critical core wiring;
mid tier at high effort per the cached routing for gate-adjacent core work.

1. As a release operator, I want `bench release submit --adapter npm` to drive the real npm CLI adapter, so that submission can reach the public registry.
2. As a gate author, I want the adapter to default to `fixture` when unspecified, so that every existing caller and hermetic contract test is byte-identical in behavior.
3. As a release operator, I want an unknown `--adapter` value refused with usage (exit 2), so that a typo cannot silently select a registry.
4. As a release operator, I want the npm adapter to publish with `--access public`, so that the scoped `@redbench/*` packages publish exactly as the raw steps did.
5. As a maintainer, I want provenance attached only on an explicit `--provenance` opt-in, so that CI gets provenance and an operator's first-publication laptop path still works.
6. As a release operator, I want `--path staged --adapter npm` refused up front with a named refusal, so that a staged run cannot die in the middle of a publication.
7. As a security reviewer, I want the adapter to take credentials only from the ambient npm configuration, so that Bench never holds, passes, or records a token.
8. As a release operator, I want `promote` and `rollback` to honor the same `--adapter` selection, so that the whole lifecycle addresses one registry.

### Group B — the workflow publishes through Bench

Line: opus / medium. YAML and evidence plumbing over settled core behavior.

9. As a maintainer, I want the publish job to run `bench release submit --path first --adapter npm --provenance`, so that publication is resumable and digest-verified.
10. As a maintainer, I want both raw `npm publish` steps gone from the workflow, so that no path bypasses the state machine.
11. As a maintainer, I want the publish job to download the publish-preflight evidence into `dist/preflight` before submit, so that the approved-set verification has its `release-index.json` and `SHA256SUMS`.
12. As a maintainer, I want the publish job to build the `bench` binary from the tag's own checkout, so that the publisher is the audited tree, not a fetched artifact.
13. As a maintainer, I want `dist/publication/publication-record.json` uploaded as a run artifact even on failure, so that a partial publication is diagnosable and resumable.
14. As a reviewer, I want no `bench release promote` invocation in CI, so that moving `latest` keeps requiring my presence.
15. As a maintainer, I want the platform-first, wrapper-last order asserted at the state machine's record, so that ordering survives the retirement of the workflow's step-name check.

### Group C — the contract flips and the docs agree

Line: opus / medium. Conformance anchors and runbook prose on the pattern the
neighboring checks already use.

16. As a gate owner, I want `checkReleaseWorkflow` to emit a diagnostic when the release workflow contains a raw `npm publish`, so that the bypass cannot return silently.
17. As a gate owner, I want the check to require the `bench release submit` invocation anchor, so that removing the swap is red rather than silent.
18. As a gate owner, I want a mutation bite test for both new diagnostics, so that the flipped check is proven red-capable.
19. As an operator, I want the runbook's first-publication section to name the exact CI submit invocation and the tag-push presence rule, so that the runbook and the workflow agree.

## Implementation decisions

- **Adapter selection is a CLI concern, not a state-machine concern.**
  `parseCommandArgs` gains `--adapter` with values `npm` and `fixture`; the
  flag defaults to `fixture` (no environment twin — the workflow passes the
  flag explicitly). `runSubmit`, `runPromote`, and `runRollback` construct
  `NewNPMCLIRegistry(registryBase)` or `NewFixtureRegistry(registryBase)` from
  that one resolved value. `RunFirstPublication`, `RunStagedPublication`,
  `RunPromotion`, and `RunRollback` keep their `Registry` interface parameter
  unchanged — the seam already exists; this change only stops hardcoding one
  side of it.
- **The npm adapter carries the arguments the raw steps carried.**
  `NPMCLIRegistry` gains `Access` and `Provenance` fields; `Publish` appends
  `--access <v>` when `Access` is set and `--provenance` when `Provenance` is
  true. The CLI sets `Access: "public"` for the `public` profile and
  `Provenance` from a new `--provenance` flag (default off). Credentials stay
  ambient (npm config / `NODE_AUTH_TOKEN`); Bench adds no token flag and the
  publication record keeps carrying no credential material.
- **The staged path refuses the npm adapter before any side effect.**
  `runSubmit` with `--path staged --adapter npm` returns a structured refusal
  naming the unimplemented capability and the working alternative
  (`--path first`), before the release lock and before any registry call.
- **The workflow's publish job becomes one submit invocation.** It checks out
  the tag, sets up Go and Node, builds `dist/bench` from the checkout,
  downloads `release-artifacts` into `dist/artifacts` and
  `publish-preflight-evidence` into `dist/preflight`, then runs
  `dist/bench release submit --version "${GITHUB_REF_NAME#v}" --profile public
  --path first --adapter npm --provenance --registry
  https://registry.npmjs.org` with `env: NODE_AUTH_TOKEN:
  ${{ secrets.NPM_TOKEN }}` on the submit step (setup-node's generated
  `.npmrc` interpolates it; without it every publish returns 401). An
  `if: always()` step uploads `dist/publication/` as the `publication-record`
  artifact. No promote, no rollback, no raw npm anywhere in the workflow.
- **The conformance contract flips in the same change as the workflow.** In
  `checkReleaseWorkflow`, the anchor map's "does not publish to npm" and "does
  not publish with provenance" entries are replaced by: a diagnostic when the
  text contains `npm publish` (presence is red), a required anchor on the
  `release submit` invocation including `--adapter npm` and `--provenance`,
  a required anchor on the `publish-preflight-evidence` download, a required
  anchor on the `publication-record` upload, and a diagnostic when the text
  contains `release promote`. The download and upload anchors are scoped to
  `workflowJob(text, "publish")`, never the whole file — the authorize job
  already uploads `publish-preflight-evidence`, so a file-wide
  `strings.Contains` is the cheapest wrong implementation. The workflow edit
  and the check flip cannot land separately: each is red against the other's
  old state.
- **Two step-name byte contracts retire with the raw steps.** The
  platform-first/wrapper-last diagnostic that indexes `name: Publish platform
  packages` / `name: Publish wrapper` in the release workflow retires — its
  concern moves to a new record-level ordering assertion over the state
  machine's package plan in `internal/publication`. The
  `preflight-publish-order-bypassed` canary fixture, which anchors the exact
  bytes of the wrapper step, retires with it: the fixture directory, its
  registration, and its row in the gate-pipeline fixture inventory go in the
  same change. Both retirements are contract amendments recorded in the
  commit.
- **The runbook amendment is part of the change, not a side note.** The
  first-publication section records the presence rule the reviewer decided
  2026-08-20: only the reviewer cuts and pushes a release tag, the tag push is
  the attended act, and the CI submit is its mechanical arm. The same section
  names the exact CI invocation. Promotion's interactive requirement is
  untouched.
- **Usage surfaces update together.** The `usageLine` in
  `internal/publication/command.go`, the inventory suffix in
  `cmd/bench/main.go`, and the help golden in `cmd/bench/main_test.go` all gain
  `--adapter npm|fixture` and `--provenance`; one grammar, three projections,
  all in the adapter ticket.
- **Bootstrap authority.** The trust root is the pushed git tag and GitHub's
  intra-run artifact chain. Hop 1: the tag triggers the workflow; every job
  checks out that ref. Hop 2: `verify` builds artifacts and evidence from that
  checkout; `authorize` verifies the *downloaded* artifacts against the
  preflight it runs from the same checkout — the tarball trust across jobs is
  GitHub's intra-run artifact chain, stated as such. Hop 3: the publish job compiles
  `bench` from the same checkout (never a downloaded binary), so the verifier
  is authenticated by the checkout, not by itself. Hop 4: `bench release
  submit` authenticates every artifact against `dist/preflight/SHA256SUMS` and
  `release-index.json` digests before any publish call, and refuses drift.
  Hop 5: the live registry bytes are re-read and their integrity compared to
  the approved local tarball before the next package advances. No hop trusts a
  path, record, or executable to authenticate itself.

## Testing decisions

- A good test drives the public command surface (`publication.Command`) or the
  adapter with a stub `npm` executable on `PATH` that records its argv, and
  observes exit codes, structured output, and the recorded invocations — never
  internal constructor identity.
- Seams that receive tests: the command argument seam
  (`internal/publication/command.go` — prior art: the existing usage-error
  tests in the package), the npm adapter's process seam
  (`internal/publication/npm_registry_test.go` — prior art: its existing stub
  tooling), and the conformance check seam
  (`internal/conformance/workflow_checks_test.go` — prior art:
  `TestNativeWorkflowEvidenceEdgeBites`).
- The gate observes the feature through `TestRootConformance` grading the live
  `release.yml` (dev tier, in the ordinary test phase since 2026-08-18) and
  through the ordinary package tests.

### Seam diagram

    tag push (vX.Y.Z)
        │
        ▼
    release.yml publish job ──▶ [ dist/bench release submit --adapter npm ] ──▶ npm CLI ──▶ registry
        │                            │        ◀ tests attach here: stub `npm` on PATH records argv;
        │                            │          Command(args) exit code + TOON output asserted
        │                            ▼
        │                    dist/publication/publication-record.json ──▶ uploaded run artifact
        ▼
    [ checkReleaseWorkflow ] reads .github/workflows/release.yml
        ◀ tests attach here: mutated workflow text in a temp root; diagnostics asserted

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| R1 | 1 | `submit --adapter npm` against a stub `npm` on PATH invokes `npm publish` with the approved tarball | publication command + stub npm | a hardcoded fixture registry never spawns npm, so the recorded argv stays empty and the row is red |
| R2 | 2 | `submit` without `--adapter` drives the fixture registry and spawns no `npm` process | publication command + stub npm | an accidental npm default spawns the stub, whose non-empty argv log reds the row |
| R3 | 3 | `submit --adapter bogus` exits 2 with a usage line naming `--adapter` | publication command | silent fallback to either adapter returns 0 or 1, never the usage exit |
| R4 | 4 | the recorded `npm publish` argv contains `--access public` under the public profile | npm adapter stub seam | an adapter that drops the flag publishes scoped packages private and the argv assertion reds |
| R5 | 5 | the recorded argv contains `--provenance` only when `--provenance` was passed | npm adapter stub seam | always-on provenance breaks the laptop path; always-off breaks CI — either constant fails one of the two assertions |
| R6 | 6 | `submit --path staged --adapter npm` refuses before the release lock and spawns no process | publication command + stub npm | the current mid-flight `StageSubmit` error only fires after publication starts, so the zero-invocations assertion reds it |
| R7 | 8 | `promote --adapter npm` and `rollback --adapter npm` drive the stub npm, and without the flag drive the fixture | publication command + stub npm | selection wired only into submit leaves the other two hardcoded and their argv logs empty |
| R8 | 10, 16, 18 | `checkReleaseWorkflow` on a workflow containing `npm publish` emits the raw-publish diagnostic | conformance bite test | the current check requires that anchor instead, so the old contract fails this row outright |
| R9 | 9, 17, 18 | `checkReleaseWorkflow` on a workflow missing the `release submit --adapter npm` anchor emits the missing-submit diagnostic | conformance bite test | a check that never flipped has no such diagnostic to emit |
| R10 | 11 | `checkReleaseWorkflow` reds when `workflowJob(text, "publish")` lacks the `publish-preflight-evidence` download | conformance bite test | a file-wide contains-check stays green off the authorize job's own upload of the same artifact name |
| R11 | 13 | `checkReleaseWorkflow` reds when `workflowJob(text, "publish")` lacks the `publication-record` upload | conformance bite test | without the job-scoped anchor a partial publication leaves no retrievable evidence |
| R12 | 14 | `checkReleaseWorkflow` reds when the workflow text contains `release promote` | conformance bite test | a promote step added to CI would otherwise pass every other anchor |
| R13 | 15 | a new assertion over the state machine's package plan proves all four platform transitions precede the wrapper transition | new ordering test in `internal/publication` | the step-name ordering diagnostic retires with the raw steps, and no live test in the package asserts order today — without this row the ordering loses its only enforcement |

Not covered: story 7 — no plausible implementation of this diff adds a
credential argument, so no seam goes red on it; the ambient-credential rule is
graded by review, not by a test.
Not covered: story 12 — the `bench`-built-from-checkout step is asserted only
as the submit-invocation anchor's surrounding job text (R9's anchor includes
the `dist/bench` invocation path); a mechanical proof that CI compiled rather
than downloaded the binary needs workflow execution, which the gate cannot do.
Not covered: story 19 — runbook prose agreement is graded by the review round,
not by a gate check; phrase-grepping operator prose demonstrates no reliable
bite (the FT106 discipline).

### Edge inventory

- Missing `dist/preflight` at submit: refused by `VerifyApprovedSet` with the
  existing structured error (existing tests own it).
- Interrupted submit: record `in_progress`, `next_action` resumes (existing
  state-machine tests own it).
- Live-integrity mismatch: terminal failure, `resolve-integrity-mismatch`
  (existing tests own it).
- `npm` binary absent with `--adapter npm`: the adapter's structured
  `npm ... failed` error surfaces as unsatisfied release intent, exit 1.
- Hostile version or tag string: reaches `exec.CommandContext` argv, never a
  shell; the tarball filename comes from `packageFileName`, not the input.
- **Won't handle:** staged OIDC submission on the npm adapter — the up-front
  refusal (story 6) survives in scope; the implementation is a separate priced
  capability below.
- **Won't handle:** CI-driven promotion — the runbook requires reviewer
  presence; the operator-run `bench release promote` path survives.
- **Won't handle:** authority pre-checks before publish — the runbook's
  precondition 1 records that registry probes prove nothing; the manual npm-UI
  confirmation survives.
- **Won't handle:** absent `NODE_AUTH_TOKEN`/OIDC in CI — npm itself fails,
  submit exits 1 with the adapter's structured error; a Bench-side auth probe
  is the excluded authority pre-check above.

## Ownership fences

- `internal/publication/command.go`
- `internal/publication/npm_registry.go`
- `internal/publication/registry.go`
- `internal/publication/` test files (new and existing `*_test.go`)
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `.github/workflows/release.yml`
- `internal/conformance/workflow_checks_test.go`
- `internal/conformance/native_workflow_test.go`
- `internal/conformance/registry_test.go`
- `internal/conformance/tier_test.go`
- `tests/canary/package-core-guard/preflight-publish-order-bypassed/`
- `tests/canary/compliance-hardening/mutable-workflow-action/MUTATE.json`
- `decisions/assets/gate-pipeline-fixture-inventory.md`
- `docs/release-runbook.md`
- `specs/ft230-release-through-bench/`

## Out of scope

- Staged OIDC trusted-publishing implemented in the npm adapter
  (`StageSubmit`/`Approve` plus the runbook's interactive approval seam):
  ~7 edits, 3 gate runs.
- A repeatable live-registry rehearsal script against a local registry
  (audit A7's validation note asks for one rehearsal before the next release;
  the operator can run it manually via `--registry <local>` today): ~4 edits,
  2 gate runs.
- Structured output for `bench release` joining the FT185 gate-output schema:
  owned by FT185, 0 edits here.

## Further notes

The audit's acceptance criterion "submit/promote invoked" reads as the command
pair being the only publication authority, not as CI invoking promote; the
runbook's reviewer-presence rule decides that promote stays out of CI, and
story 14 pins it. Deployment is NO-GO today, so nothing here performs a
release; the workflow change is proven by conformance anchors and the fixture
path, and the audit's live-registry rehearsal remains an operator step before
the next release.
