# Bench platform assessment — 2026-08-13

Assessment target: `main` at `287bfbf9`, plus the dirty assessment artifacts
named below. This refreshes the 2026-08-12 assessment after the canary-ownership,
serial-gate, and worktree-fixture work. Six read-only area sweeps
covered adoption and packaging, workflow guidance, enforcement, CLI/Go core,
gate authority and records, and live operational state. Load-bearing claims
marked ✓ were independently reverified against source or live output.

Severity in this report:

- **high** — an invariant or advertised guarantee is not actually held;
- **med** — a real defect or reachable unowned state;
- **low** — friction, drift risk, or hygiene.

## Decision

**NO-GO for public release or external deployment.** The July critical
worktree, gate-cache, shift-recovery, packaging, and unknown-hook failures are
substantially repaired. Canary inventory and planted-reason proof now have
separate truthful owners, and the source tree's full gate is green. Release is
still blocked because the tag workflow bypasses the governed, resumable
publication lifecycle.

Current inventory: **1 high, 7 med, and 4 low** findings. Live defect-shaped
roadmap work is prioritized in `capture/FIXES.md`; features and decision-only
opportunities remain solely on the roadmap.

## Previous assessment reconciliation

Each prior ranked backlog item was checked against the current tree rather than
commit subjects.

| Prior rank | Current | Reverified evidence |
| ---: | --- | --- |
| 1. Ownership-safe worktree cleanup | **FIXED ✓** | Automatic cleanup requires verified ownership, landed state, valid recovery metadata, and clean content (`internal/worktree/classifier.go:251-293`); foreign-worktree preservation is covered at `internal/worktree/resume_test.go:19-69`. |
| 2. Gate-cache oracle identity and atomic records | **FIXED, coverage partial ✓** | Oracle identity is part of exact reuse (`internal/gate/subject.go:63-110`; `internal/gate/verdict.go:250-269`) and records replace atomically/fail closed (`internal/gate/verdict.go:146-190`). The old A-green/B-red black-box repro is not retained as a current regression. |
| 3. Shift failure recovery and truthful completion | **FIXED ✓** | Stage errors propagate (`internal/shift/shift.go:55-75`); adapter and commit failures select recovery evidence instead of success/teardown (`internal/shift/loop.go:202-206,251-289,316-318`). |
| 4. Static guard metadata and aggregate deadline | **FIXED ✓** | Guard descriptions are parsed, not executed (`internal/guards/guards.go:103-153`), under one bounded scan (`internal/guards/guards.go:164-200`). |
| 5. Clean package staging and runtime selection | **FIXED ✓** | Release staging omits source `prepare` (`scripts/build-release-evidence.mjs:221`); the launcher chooses the selected platform package before development `dist` (`bin/bench.sh:164-180`). |
| 6. Gate-backed release and security/package evidence | **PARTIAL ✓** | Native verification and publish preflight are ordered before publication (`.github/workflows/release.yml:12-47`), but the publish jobs still call raw `npm publish` instead of `bench release` (`.github/workflows/release.yml:49-73`). |
| 7. Installed shim/adoption routing and fresh-clone runtime | **FIXED ✓** | Adoption commands route to the installed target before local-wrapper preference (`internal/adopt/doctor.go:110-142`); platform resolution precedes development fallback (`bin/bench.sh:164-180`). |
| 8. One-command bootstrap | **PARTIAL ✓** | `bench setup` now seeds a project-local profile transactionally (`internal/adopt/setup.go:289-360`), but no public package/tag exists and README presents installed harness commands before installation (`README.md:9,176-194`). |
| 9. Transactional relink/upgrade/unlink | **PARTIAL ✓** | Link stages and promotes atomically (`internal/adopt/link_transaction.go:164-190,370-399`) and partial unlink exits nonzero (`internal/adopt/unlink.go:55-68`); uninstall guidance can still delete a foreign executable. |
| 10. Consumer/maintainer payload split | **FIXED ✓** | One payload manifest marks kit-only surfaces (`.bench/consumer-payload.json:10-15`) and link consumes only consumer rows (`internal/adopt/link.go:75-95`). |
| 11. Coverage, read-state, default-branch, and intent correctness | **PARTIAL ✓** | Exact story membership and missing-map rules landed (`internal/coverage/coverage.go:75-140`); `ResolvedDefault` is the single branch owner (`internal/git/default_branch.go:7-35`); the full concurrent intent-lock posture was not re-reproduced. |
| 12. Offline/reproducible release evidence | **FIXED, native limits remain ✓** | The release plan declares offline archives (`scripts/release-plan.json:3-18`) and offline smoke installs them (`scripts/smoke-offline.sh:210-270`). macOS/ARM and live registry execution were unavailable. |
| 13. Output/resource/CLI consistency and stale docs | **PARTIAL ✓** | Shared grammar and bounded model discovery landed (`internal/models/models.go:21-40,81-172`); workflow/document contradictions and narrow-reader work remain under FT89, FT106, and FT125. |

## Findings by area

### Adoption and packaging

**H-A1 — Release automation bypasses its governed publisher.** The runbook
requires `bench release submit` and `promote` with live digest verification and
resumption (`docs/release-runbook.md:35-78`), while the tag workflow directly
loops over `npm publish` and then publishes the wrapper (`.github/workflows/release.yml:49-73`). A partial immutable publication therefore bypasses the state machine in
`internal/publication/statemachine.go`. **Tracked: FT142.** ✓

**M-A1 — No supported public entry point exists.** There are no repository
tags, and README truthfully says `redbench` is not published
(`README.md:176-194`). This is a release-readiness blocker, not a false claim.
✓

**M-A2 — First-hour guidance remains out of order and split on profiles.** The
README begins with harness commands before installation (`README.md:9,176`),
and the legacy `init` path still tells consumers to consult a kit-side project
example (`internal/adopt/init.go:70-82`) even though `setup` can emit the
project-local profile. ✓

**M-A3 — Uninstall guidance can delete a foreign executable.** Doctor refuses
to overwrite an unmarked shim (`internal/adopt/doctor.go:304-318`) but its
status output still prints `npm uninstall -g redbench && rm -f <path>` for the
resolved path (`internal/adopt/doctor.go:189-237`); README prints the same raw
deletion (`README.md:232-239`). ✓

**L-A1 — Source-pack and release-wrapper semantics differ.** `package.json`
runs `prepare` for an ordinary source `npm pack` (`package.json:14-21`), while
release staging deliberately constructs another wrapper shape
(`scripts/build-release-evidence.mjs:221`). Both can be legitimate, but the
distinction needs an explicit contract so a source tarball is not mistaken for
the release artifact. ✓

### Workflow commands and skills

**M-W1 — The skills index still has two derivations and neither parses YAML.**
The shell generator extracts frontmatter independently
(`.bench/skills-index.sh:23-81`) from the Go conformance implementation
(`internal/conformance/skills_index_checks_test.go:209-270`). The existing bite
tests make this visible but do not satisfy one-source-per-fact. **Tracked:
FT89.** ✓

**M-W2 — Four workflow contracts remain mutually inconsistent.** Design-it-
twice supplies bare briefs while `craft-delegate` requires a complete charge
(`.agents/skills/bench-craft-seams/references/design-it-twice.md:3-23`;
`.agents/skills/bench-craft-delegate/SKILL.md:35-47`); `craft-synthesis`
recognizes only upstream/learnings origins despite assessment-driven work
(`.agents/skills/bench-craft-synthesis/SKILL.md:9-16`); `CONTEXT.md:18` says
token cap while `craft-line` says iteration cap; README's uncertain Regroup
example routes to shift contrary to the cheap-and-observable rule
(`README.md:356-369`; `.agents/commands/bench-write-spec.md:16-21`). **Tracked:
FT89/FT102.** ✓

**L-W1 — Reference and inventory drift remains.** `internal/lines/lines.go:11`
cites the removed `.bench/lib/lines-env.sh`; the AXI guidance inventory omits
some current query surfaces; active guidance still depends on a mutable AXI
URL (`.agents/skills/bench-craft-cli/SKILL.md:7-12,71-82`). **Tracked:
FT89/FT106.** ✓

### Enforcement and oracle

Planted-reason ownership is no longer a finding. Every retained kit fixture has
direct mutation-specific proof through its owning in-process check or owning
package coverage; `bench canary` remains inventory-only, and linked repositories
own proof in their native tests. Ambient status also projects the retained exact
gate subject and identifies staleness rather than reporting the verdict
unavailable. ✓

**L-E1 — Capability evidence is explicit but incomplete on this host.** The
green gate reported five skips, including one FIFO and three privilege skips;
the guard FIFO regression is capability-skipped on Windows
(`internal/guards/guards_test.go:333-356`). Strict capability execution and
native runners remain required release evidence. ✓

### CLI and Go core

**M-C1 — Shift objectives still cross unnecessary durability surfaces.** The
objective uses a private 0600 scratch file (`internal/shift/loop.go:180-186`),
but is also printed in the banner and interpolated into durable commit subjects
(`internal/shift/loop.go:225-227,276-278`). This is the surviving C-08 data-
handling remainder. ✓

### Live operational state

**M-O1 — One active decision map is structurally invalid.** `bench maps`
reports `spec-build-review-gate-cadence` pointing at removed spec-build code.
Eight decision rows remain unresolved. This is operational reconciliation work, not
authority to rewrite the maps during assessment. ✓

**L-O1 — Structure debt is broad but not uniformly architectural.** `bench
structure` reports 50 issues, including crowded `internal/adopt` (23 files),
`internal/gate` (21), and `internal/worktree` (30). File length alone does not
pass the deletion test; only verified shallow modules or leaky seams should
become deepening candidates. **Tracked generally: FT108/FT186 and the
architecture report produced with this assessment.** ✓

## Ranked improvement backlog

| Rank | Work | Evidence required | Rough agent time |
| ---: | --- | --- | ---: |
| 0 | FT198: shape a progressively loaded roadmap | One canonical detail owner; migration and history preserve status; index completeness is mechanically checked | 0.5–1 day decision |
| 1 | FT189: bound the upstream worktree-list hang | The public command refuses malformed administration state or terminates within the declared bound | 0.5–1 day |
| 2 | FT142: route CI publication through `bench release` | Partial publish resumes with live digest checks; wrapper cannot lead platforms | 1–2 days |
| 3 | Close A2/A3/A10 adoption residuals | Packed first-hour setup; emitted local profile; marker-verified shim removal | 1–2 days |
| 4 | FT133/FT174: harden coverage and ticket ownership | Every acceptance row and ticket edge resolves to one accountable producer and authorized fence | 1–2 days |
| 5 | FT89/FT102: single-source skill indexing and reconcile workflow contracts | One generator/verifier owner; conformance mutations for each closed contradiction | 1–2 days |
| 6 | Close C-08 objective exposure | Objective identifier, stdout, scratch, and commit durability follow one documented policy | 0.5–1 day |
| 7 | FT162/FT185: unify terminal subject and structured gate evidence | Green/red runs project the same subject and outcome without reconstruction | 1 day |
| 8 | Reconcile the invalid map through maintenance | `bench maps` has no invalid rows and no vanished source paths are invented | 0.25–0.5 day |

The 24 live fix-class rows are prioritized in `capture/FIXES.md`. FT198 remains
the overall roadmap's current reviewer-decision entry and is intentionally not
in that fix-only view; no spec is staged.

## Verification notes

Reverified ✓:

- `GOCACHE=/tmp/bench-go-cache npm_config_cache=/tmp/bench-npm-cache bench gate`
  completed green in 37.8 seconds: gofmt, vet, test, race, system, and shellcheck
  all green; exact log `.logs/gate-20260812T224413.050969588Z-1689120.jsonl`.
- Focused packages passed: `internal/adopt`, `releaseevidence`,
  `releasepreflight`, `gate`, `commit`, `guards`, `canary`, `conformance`,
  `worktree`, `shift`, `coverage`, `intent`, `git`, `outline`, `models`, and
  `commit`.
- The FT153 build gives every retained kit fixture direct planted-reason proof,
  keeps `bench canary` inventory-only, and leaves linked proof to native tests.
- Decision #27 retired the serial-gate build because the serial baseline already
  meets its destination. FT175 was retired rather than built.
- `TestListCommandPublicRowsAndDisclosure` and its completed-assignment sibling
  passed 50 deterministic leading-zero-ID repetitions before FT203 retired.
- `bench structure` returns 50 issues; `bench maps` returns 10 rows, one invalid;
  no active staged spec remains.
- Current tree state during close: assessment-owned `ASSESSMENT.md` and
  `capture/FIXES.md`; registered assignment worktrees were not edited, cleaned,
  released, staged, or committed.

Coverage limits:

- Nothing was published; no live npm registry, GitHub Actions execution,
  publisher credentials, macOS/ARM, Windows, musl, or bank-profile release was
  exercised.
- The old gate A-green/B-red, concurrent intent-reclaimer, and every
  capability-strict platform probe were not rerun.
- No architecture candidate has been selected or authorized for refactoring;
  the 50 structure findings remain inventory, not proof of a better seam.
- Linked-consumer planted-reason proof was not rerun during this refresh; its
  native-test ownership is contract-derived from the landed FT153 work.
