# Bench capability-realization audit evidence

Audit subject: commit `58d966e2f92f7f37eba07b6215e8eef45371b72d` on branch `audit/sol`.

The initial exact-source `bench diff` snapshot reported one untracked input file,
`audit/bench-relentless-capability-realization-grill.md`, and no staged, unstaged,
or committed movement from `main`. Audit outputs are intentionally excluded from
that initial snapshot.

## Safety boundary

- Allowed: read-only repository inspection, local exact-source builds, local tests,
  and isolated disposable fixtures removed after observation.
- Prohibited: commits, pushes, releases, deployments, secret reads, destructive
  operations, weakening checks, and external mutation.
- External research is read-only and revisions are pinned where the source permits.

## Observation log

| ID | Observation | Evidence class |
|---|---|---|
| E001 | The PATH wrapper prints its static command inventory for `bench help`, but `bench commands --brief`, `bench status`, and `bench version` all exit 127 with `no pinned binary for this platform` before an exact-source local build. | OBSERVED |
| E002 | The documented exact-source build command produces `bench 0.2.0 (linux/amd64)` without installing or fetching a platform binary. | OBSERVED |
| E003 | The initial exact-source status reports one dirty path (the audit charge), one out-of-pool worktree, 62 structure issues, and seven unresolved decision maps. | OBSERVED |
| E004 | The exact-source binary rejects `bench help` with exit 2 while the PATH wrapper accepts it; `bench commands --brief` returns only `version`, `commands --brief`, and `status`. | OBSERVED |
| E005 | Direct repository state after the local build was branch `audit/sol`, HEAD `58d966e2f92f7f37eba07b6215e8eef45371b72d`, no staged or tracked diff, and only `audit/` untracked. Origin is `https://github.com/gibbonmi/bench.git`. | OBSERVED |
| E006 | `bench test` passed 67 packages with 10 capability skips, but under the documented audit environment it also left ten untracked `worktrees/001-*` fixture directories in the subject repository. They were moved to the desktop trash, after which status returned to only `?? audit/`. | OBSERVED |
| E007 | Focused executable tests passed for exact-tree gate reuse, red/timeout invalidation, moved-subject refusal, ownership fences, uncited/phantom acceptance rows, coverage membership, and shift staging-failure recovery. | OBSERVED |
| E008 | `bench gate --fresh` ran gofmt, vet, test, race, system, and shellcheck; it finished green in 59.8 seconds with seven capability skips and recorded exact tree `29f339a0133fae9d7e0ecb6afe6c648554a3f224`. A second run reused that verdict. | OBSERVED |
| E009 | The same green gate coexists with seven unresolved maps (one invalid source path), 62 structure issues, and a dirty audit path. Those are status/advisory domains, not members of this repository's gate oracle. | OBSERVED |
| E010 | `bench review`, `bench final-check`, and `bench claims` each reject as unknown subcommands with exit 2. Review and Final Check exist only as model-read command documents; no general claim record or claim-state command exists. | OBSERVED |
| E011 | `capture/session-handoff.md` says the repository is `~/workspace/bench` on `main` at `9616919a` with a completed roadmap drain, while the actual worktree is `bench-audit-sol` on `audit/sol` at `58d966e`. `bench status` reports no handoff lag because staleness is measured only by commits since the handoff file was last written. | OBSERVED |
| E012 | In isolated clones, `bench handoff --harness codex|claude` refreshed repository pins and preserved the entire stale reviewer-owned `State` body verbatim. The derived next command was `bench link` for both; a manual `--next /bench-debug` override is intentionally not translated. | OBSERVED |
| E013 | `bench canary` reports 233 fixture bindings, and `TestEveryRetainedFixtureBitesThroughRegisteredOwner` passes. Many named fixtures guard literal prompt/document anchors; this proves conformance checks detect their planted mutations, not that a model enacts the prose. | OBSERVED |
| E014 | Project guidance totals 31,456 words across skills, commands, AGENTS, BENCH, and profile; the mandated cold set `AGENTS.md` + `.bench/BENCH.md` + `projects/benchkit.md` + `CONTEXT.md` is roughly 6,000 words before task-specific guidance. The five listed cold artifacts including reference and handoff total 8,404 words. | OBSERVED |
| E015 | All 28 project skill files and all 11 command documents were read. The current Codex session exposed the 17 craft/prototype skills but none of the 11 `bench-*` phase adapters even though their files exist; Claude has symlink adapters. | OBSERVED |
| E016 | No live `specs/*/spec.md` exists, but about forty ticket files remain. `bench coverage` on a ticket-only slug reports `spec not found`, while the spec status reader ignores directories without `spec.md`. | OBSERVED |
| E017 | Review guidance says to commit transient pickup state and also says the review phase runs no gate. The only safe public commit command (`bench commit`) runs the gate. It also says “the gate and you” decide done although the canonical invariant says only the gate does. | OBSERVED |
| E018 | Status/handoff routing selects the first syntactically invocable action and skips higher-severity prose actions. On the audit board it can bypass dirty-tree, worktree, and structure advice to reach a later phase command; `TestFirstInvocable` confirms that rule. | OBSERVED |
| E019 | The full root wrapper with no arguments prints a 40-plus-command inventory and exits 0; the exact binary with no arguments says `no subcommand` and exits 2. Neither routes user intent. | OBSERVED |
| E020 | The current release workflow verifies and authorizes, then publishes packages with raw `npm publish`; the runbook requires the resumable `bench release submit/promote` state machine. The repository's prior assessment records the same release blocker. | OBSERVED |
| E021 | Current upstream snapshots inspected: `mattpocock/skills@9c9f36ccd3995266cd675468af71639c8dde1ec5`, `kunchenguid/no-mistakes@6859d1e827f5ab2592a4703d3bab8734a38c9aa5`, and `kunchenguid/axi@408a6536625e5b05e5c56e6c4a04fe83e1f510a5`. Bench's documented v1.1 adoption maps to `mattpocock/skills` tag `v1.1.0` commit `d574778f94cf620fcc8ce741584093bc650a61d3` and Bench commit `5e3c0ba98500a904be47673be975f8770cb33d0d`; later roadmap comparisons cite upstream `84fdeffd12f2ee307994d1eb6feb48173b6e0502`. | OBSERVED |
| E022 | Bench's README names no-mistakes and AXI as sources but records no exact imported upstream revisions for them, so exact historical source-to-integration attribution cannot be reconstructed from the repository. | OBSERVED |
| E023 | After all controlled experiments, the two handoff clones, upstream clone, isolated test home, exact-source binaries, ignored `dist/bench`, two audit gate logs, shared gate evidence, and last-gate record were moved to the desktop trash. The final guarded Git status was branch `audit/sol` with only `?? audit/`; staged and tracked diff lists were empty. | OBSERVED |

## Controlled-experiment matrix

| Experiment | Execution | Result | Classification |
|---|---|---|---|
| Entry routing | Ran wrapper and binary with no args/help; ran status and handoff selection tests. | Inventory and state signals exist, but no canonical intent router distinguishes new, ready, failing, interrupted, unverified, and completed work. | CONFIRMED GAP |
| Gate differentiation | Ran all 233 registered planted-reason owners and focused gate outcome tests. | Gate catches registered contract mutations and invalidates exact evidence on red, timeout, or subject movement. | CONFIRMED |
| Final-check differentiation | Invoked `bench final-check`; inspected and exercised status housekeeping owners. | No executable Final Check exists. Status detects retirement/orphan/dirty cases, but the prompt phase adds no independently testable control. | UNREPRODUCIBLE AS CLAIMED |
| Review independence | Invoked `bench review`; inspected review protocol and current no-mistakes isolation design. | Bench has prompt-level fresh-axis instructions, no executable review runner or retained clean-review evidence; controlled with/without-narrative trials were unavailable without authoring new orchestration. | UNREPRODUCIBLE |
| Stale evidence | Ran exact gate reuse/invalidation tests; compared live handoff to tree and handoff-age tests. | Gate evidence freshness is strong; handoff semantic freshness is weak and can be false at the same commit. | MIXED |
| Claim bypass | Invoked `bench claims`; searched command and state owners. | There is no general claim system to bypass. Semantic assertions can be written directly into handoff `State` and prompt-produced artifacts without evidence references. | CONFIRMED ABSENCE |
| Ownership violation | Ran build/review out-of-fence, uncited-row, phantom-row tests. | Detection is deterministic at preflight; before preflight it remains a workflow obligation. | CONFIRMED, LATE-BOUND |
| Requirement loss | Ran coverage membership/violation tests and inspected ticket-only live state. | Mapped rows are checked, but no live specs mean current ticket artifacts have no executable coverage query. Semantic adequacy remains review-only. | MIXED |
| Context recovery | Rewrote handoffs in two isolated clones and compared pins/state. | Pins recover; arbitrary state survives even when contradicted by the tree. No bounded active-state schema or coverage/attempt history exists. | PARTIAL |
| Model swap | Rendered Claude and Codex handoffs and ran harness-prefix conformance tests. | Prefix translation is single-sourced for derived phase actions; durable task state is unchanged. This session did not expose Codex phase adapters. | PARTIAL |
| Failure inheritance | Ran `TestLoopStagingFaultPreservesAndSplitsEvidence`. | Shift preserves committed-progress/recovery evidence; ordinary prompt phases have no equivalent failure record. | PARTIAL |
| Diagnosing comparison | Compared current Bench command to current upstream skill and the required no-skill/default variants. | Textual mechanism preservation is observable; multi-trial outcome comparison was not available in this audit and no causal claim is made. | DESIGNED, NOT EXECUTED |

## Research and upstream pins

Primary research URLs and the exact upstream comparison links are catalogued in
sections E and X of `report.md`. Research claims are labeled `RESEARCH`; the
DeepSeek/J-Space report is treated as single-run community evidence, not as a
universal result.
