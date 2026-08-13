# Changelog

All notable user-facing changes to Bench are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/2.0.0/), and versions follow
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- `bench link` now adopts only same-file converged adapter symlinks on first link.

### Removed

- Removed the provisional spec-build lifecycle wholesale: the `bench spec build`
  family (`start`, `assign`, `checkpoint`, `integrate`, `review`, `status`,
  `promote`, `abandon`, and the maintainer-run `reclaim`) and
  `bench worktree recovery` are no longer commands and answer with the standard
  unknown-subcommand structured error. `bench resume`'s reconcile now deletes
  the `refs/bench/specbuild/` and `refs/bench/recovery/` namespaces and purges
  lifecycle-typed assignments from the intent ledger at every session start,
  idempotently. Zero backwards compatibility — no shim, no migration tooling;
  reviewed spec-backed builds land tickets serially commit-on-green through
  path-scoped `bench commit`, with `--spec <slug>` on the final green landing
  commit marking the spec implemented.

### Added

- New `/bench-deepen` phase (Codex: `$bench-deepen`) — an architecture-deepening
  survey ported from Matt Pocock's `improve-codebase-architecture`: scans for
  shallow modules, presents deepening candidates as a self-contained HTML report,
  then grills the reviewer's pick. Scopes from a named direction, the latest
  `ASSESSMENT.md` (the `/bench-assess` conjunction), or commit-history hot spots;
  `/bench-assess` now routes architecture-shaped findings to it on exit.
- `bench diff` now renders one coherent AXI snapshot with revision and aggregate
  facts, all-files checkout inventory, whitespace verdicts, exact untracked bodies in
  `--full`, and executable follow-up help; it retries a moving live checkout once and
  then refuses a torn response.

- Added the standalone `prototype` skill (`.agents/skills/prototype`): a
  disposable prototype answers one named question, runs trivially, keeps state
  in memory unless persistence is the question, surfaces the relevant state,
  records the verdict where the question was asked, and is then discarded —
  no branch-retention route. The `.claude/skills` adapter check now admits a
  standalone skill symlink resolving to its own `.agents/skills` SKILL.md
  while still rejecting phase-adapter duplication and dangling links.
- Added on-demand reference leaves to `craft-tdd` (`references/tests.md` for
  the four properties of a good test, `references/mocking.md` for system-seam
  mocking and honest stubs) and to `craft-seams`
  (`references/dependency-categories.md` mapping the in-process,
  local-substitutable, remote-owned, and true-external dependency classes to
  their test strategies), and deepened `references/design-it-twice.md` with
  what makes candidate interfaces radically different.
- Added the `craft-domain` companion skill (`.agents/skills/bench-craft-domain`):
  canonical domain terms with Avoid lists, concept-edge scenarios,
  producer-derived equivalence partitions, code-versus-claim comparison, and
  glossary-only `CONTEXT.md` maintenance, charged from grilling, shaping, and
  spec authoring while hard-to-reverse decisions stay with `craft-adr`.
- Added `bench preflight review|build <slug>` to check at phase entry that a spec's
  artifacts agree with the tree, one verdict row per check and exit 0 only when
  every applicable check is green.
- Added `BENCH_RED_MUTATIONS_OPTIONAL=1`, a reviewer-set experiment flag under
  which `bench spec build assign` and refresh accept a modern ticket declaring
  neither a `Closure:` inventory nor `## Red mutations` rows; a declared
  inventory still validates in full, and the unset default is unchanged.
- Added `bench spec build assign --refresh <debug-receipt>`: an out-of-fence repair
  route that re-bases a blocked, uncheckpointed assignment onto the repaired
  candidate, preserving its attributed in-fence work byte-for-byte behind a durable
  preservation ref, with fail-closed refusals for forged receipts, out-of-fence
  payloads, replay conflicts, and candidate movement, and convergent re-entry after
  interruption.
- Added a compiled Go core and native platform packages with installed-binary
  smoke coverage and bounded repair.
- Added the ambient dashboard, standalone dashboard page, and AXI query surfaces
  for status, roadmap context, guards, diffs, coverage, learnings, and maps.
- Added recoverable gated shifts and identity-safe worktree lifecycle management.
- Added an authoritative release preflight that binds source, artifact, version,
  ancestry, and changelog evidence before publication.
- Added profile-governed release evidence, exact component inventories, and
  diagnostic focused release preflight runs.
- Added reproducible per-platform offline archives for direct, local npm, and
  internal-registry use, plus `bench commands --brief` for compact command discovery.
- Added `bench worktree recovery <ref> --discard <fingerprint>` to retire a preserved
  recovery payload the landedness proof does not accept, one exact ref and fingerprint
  per invocation.
- Added `bench spec build reclaim <slug> [--apply <fingerprint>]` for a maintainer to
  plan, then apply, the deletion of one terminal spec-build run's leftover provisional
  refs, deleting only what the plan can prove dead and reporting the rest.
- Added a fail-closed prose-budget conformance check with a canary mutation:
  `projects/benchkit.md` owns the one budget table, the check parses it rather
  than repeating its numbers, classifies every skill automatically, and
  refuses malformed or duplicate policy before reading.

### Changed

- Made `bench canary` report only validated canary inventory; retained kit fixtures
  are proven by ordinary mutation tests, and newly initialized linked repos leave
  planted-reason proof to project-native tests.
- Recast the AXI CLI guidance as a ten-principle per-surface contract with
  contextual result actions and a registry-checked approved-query inventory.
- Replaced one-question-at-a-time grilling with numbered frontier rounds: every
  question whose prerequisites are settled appears in the same round with a
  recommendation, then the skill waits and recomputes; `bench-shape-idea` uses
  the same frontier vocabulary. Light-path TDD now stops before its first test
  for reviewer seam confirmation, and ticket breakdown returns to an explicit
  reviewer round trip — a numbered title/`Blocked by:`/delivered-outcome list,
  iterated and approved before assignment — with the existing batch-approval
  AFK carve-out as the only no-round-trip exception.
- Slimmed `craft-tickets` to the independently-green tracer rule, `Blocked by:`,
  `What to build`, `Acceptance`, and the advisory `Writes:` note, dropping
  Contracts, Integration surfaces, Closure, covers annotations, red-mutation
  tables, handoff ledgers, and fence enforcement; `bench-implement-spec` is now
  the short orchestration pointer to that ticket shape, commit-on-green
  cadence, TDD, review, and final-check. `craft-delegate` shed receipt,
  lifecycle, duplicated mutation, and charge ceremony down to its safety core
  (fresh context, explicit line and bounded charge, worktree isolation,
  independent done-claim verification). `craft-line` kept the tier binding and
  ladder, but now classifies only reds the current diff owns, against the
  pinned inherited baseline and spec-predicted reds, before a retry or
  escalation.
- Made review re-derive-then-compare on all three axes: Coverage independently
  enumerates producer membership and the spec's authorized write set, Spec
  quotes the applicable spec line and drives from behavior, and Standards
  keeps the Fowler baseline, all in parallel fresh contexts. Every finding
  carries exactly one disposition (`no-op`, `auto-fix`, or `ask-user`);
  actionable findings are written to `reviews/<slug>.md` and committed before
  repair begins.
- Slimmed `.bench/BENCH.md` to its own 150-line budget with restated
  predicates in place of the removed prose, and added FT107's remaining
  operational rules to `AGENTS.md`: PID-or-sentinel waits, plan-before-apply
  destructive scripts, `rg --hidden` for repository-wide sweeps, and
  non-interactive Bench verb discovery through `bench commands --brief` or
  source.

- Made spec implementation route fence drift through repaired approval and require
  `craft-tickets`' complete handoff ledger before lifecycle start, with a
  section-sensitive canary for every spec-to-ticket handoff clause.
- Made `/bench-debug` distinguish the regression-test seam from the edit owner:
  enumerate relevant callers, fix a uniform invariant once at its narrowest shared
  owner, keep the test at the highest observable seam, and report paths with no
  honest shared owner as an architecture finding.
- Made spec authoring trace bootstrap authority from the raw OS entrypoint and
  require every executable hop to authenticate its successor before launch.
- Made a validated debug receipt's required fence the maximum repair-reslicing
  envelope: independently-green slicing yields one repair ticket or an ordered
  chain, the union of their ownership fences remains inside that envelope, and
  the original blocked assignment refreshes only after the terminal repair
  ticket lands.
- Replaced the fixture-driven dev gate with a branch-native architecture: one
  package-universe test driver, registry-derived race coverage, one bounded tagged
  system owner, direct mutation-to-check canaries, and immutable Go command and
  lifecycle decisions. Retired nested Go/canary drivers, repeated contract and
  conformance repositories, component partitions, and stripped-subject reruns.
- Made `bench commit` compose and authorize only its named paths, so unrelated
  concurrent working-copy and index state remains intact while the exact tree
  that can land receives the gate verdict. A destination advance now refuses
  publication without overwriting the winner; rerunning recomposes from its tip.
- Made `bench worktree recovery` accept one TOON quote layer around an exact
  fingerprint while retaining its lowercase 64-hex validation.
- Made agent-guidance Markdown edits under `.agents/` skip the Go toolchain
  gate components through the per-component input declarations, while the
  contract and canary components that consume the guidance tree still run.
- Made changes confined to active decision maps and their research assets eligible
  for the reduced gate run, alongside capture and compiled spec documents.
- Made decision maps situational, separated shaping decision tickets from
  independently-green implementation tickets, and moved engineering seams,
  tests, coverage, hostile inputs, and gate attachment into spec authoring.
- Changed gate and contract launches to refuse stale or indeterminate Bench
  binaries before execution and report one copy-paste rebuild command.
- Made Research decision tickets use concurrent read-only delegation when another
  unblocked ticket can proceed, with an inline fallback.
- Extended roadmap maintenance with per-run fix-versus-feature classification,
  a fixes-first recommended sequence, and an opt-in `--restructure` board pass
  whose merge, collapse, fold, and theme-grouping moves are proposed in the
  batch diff.
- Cut the shared communication rules to a smaller, conflict-free surface: the
  warm senior-colleague register still locates warmth in sentence construction
  and rules out self-criticism; one formatting rule separates parallel facts
  from cohesive prose; structured labels fire only during an invoked
  harness-native `/bench-*` or `$bench-*` phase; and the decision, line, and
  lighter-path rules state their actual scope.
- Made semantic review hand off instead of repairing: the terminal repair-pass
  bound now lives with implementation, and final check lands through
  `bench commit` as its one oracle run instead of a standalone gate first.
- Made implementation start from independently-green ticket files, with one
  fresh write-delegate charge and one atomic full-gate landing per ticket.
- Made reviewed spec builds checkpoint and integrate ownership-safe frontier
  tickets provisionally, review the exact composition, and gate once at promote.
  Light-path changes, `bench shift`, and ordinary `bench commit` keep
  commit-on-green cadence.
- Made final check capture a spec-backed implementation retro after green
  landing, with reviewed disposition and removal through roadmap maintenance.
- Co-located compiled decision maps and their owned assets with their specs while
  keeping top-level `decisions/` reserved for pre-spec decision work.
- Realigned the skill layer to Claude 5 harnesses: a new harness-echo rule keeps
  skills from restating what the harness carries, fork delegations are governed
  by the line, the review smell glosses moved behind progressive disclosure,
  grills use the harness's structured question surface, and auto-memory is
  routed apart from the learnings journal.
- Made managed asset link, relink, and unlink report partial lifecycle outcomes explicitly.
- Made binary repair explicit, manifest-pinned, resource-bounded, concurrency-safe,
  and manually prunable; unified user-facing output and package metadata on Bench.
- Bounded Go-side network, subprocess, read, and output work; made worktree
  refresh explicit and offline suppression observable across providers.
- Removed internal roadmap and learning history from release notes, leaving Git as
  the historical record, and standardized release headings for preflight checks.
- Made linking and relinking preserve project-owned instructions, reject managed
  asset conflicts, and behave idempotently across supported harnesses.
- Consolidated shared workflow rules, command phases, and hook behavior so Claude
  Code, Codex, and other `AGENTS.md` harnesses consume one portable kit.
- Parallelized independent gate phases while preserving one aggregate oracle
  verdict and canary coverage.
- Reduced gate time by building host-only artifact fixtures while retaining a
  two-target matrix-iteration contract.
- Made `bench worktree recovery --apply` refuse instead of silently exiting zero when
  a supplied fingerprint accompanies a plan that authorizes no action.

### Fixed

- Made `bench worktree recovery` refuse both verbs for a ref its plan cannot classify:
  only a verified payload the landedness proof refused plans the discard-eligible
  `discard` verdict, and a re-run converging on an interrupted verb records `discarded`
  rather than the proof-backed `retired`.
- Made prospective gate execution compile the exact unpublished Bench tree without
  requiring ignored binary artifacts from another checkout.
- Made spec-build integration release an exact live checkpoint payload after
  candidate advancement while retaining drifted assignments for safe re-entry.
- Made freshness source discovery ignore ambient VCS metadata, matching the
  repository's VCS-independent build and keeping gates runnable from nested
  temporary fixture roots.
- Kept session-handoff capture writes out of stale-gate and handoff-next signals, and made
  ambient dirty-path counts describe only the current checkout.
- Made `bench learnings` render a drained journal as empty and reject dated headings
  without a final `[open]` state across learnings and drain status.

- Made blocked branch-deletion refusals name the safe or force form the caller used.
- Made `bench guards` finish timed-out scan workers before returning, and kept
  its timeout races under the dev gate's race check.
- Made `bench gate` treat help and invalid arguments as usage without starting
  the oracle or replacing gate evidence.
- Included committed, staged, and tracked working-tree changes in branch-relative
  review diffs.
- Removed landed non-default branches during worktree and session cleanup while
  preserving checked-out and unique branch work.
- Preserved interrupted or failed shift work under durable recovery references.
- Prevented gate-cache and declared-input changes from silently authorizing stale
  verdicts.

## [0.2.0] - 2026-06-27

### Added

- Added `npx` distribution, safe repository linking, interactive setup, and
  upstream synthesis maintenance.

### Changed

- Made the kit application-agnostic and moved project-specific policy into project
  profiles.
- Generalized design-system consumption so committed artifacts, rather than one
  authoring harness, are the workflow dependency.

## [0.1.0] - 2026-06-27

### Added

- Initial release of the decision-mapping, specification, implementation, review,
  and debugging phases; seam-aware TDD and ADR guidance; the gate-as-oracle shift
  loop; isolated worktrees; and destructive-Git guards.
