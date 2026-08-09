# Changelog

All notable user-facing changes to Bench are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/2.0.0/), and versions follow
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

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

### Changed

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
