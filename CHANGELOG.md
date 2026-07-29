# Changelog

All notable user-facing changes to Bench are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/2.0.0/), and versions follow
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

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

### Changed

- Extended roadmap maintenance with per-run fix-versus-feature classification,
  a fixes-first recommended sequence, and an opt-in `--restructure` board pass
  whose merge, collapse, fold, and theme-grouping moves are proposed in the
  batch diff.
- Made semantic review hand off instead of repairing: the terminal repair-pass
  bound now lives with implementation, and final check lands through
  `bench commit` as its one oracle run instead of a standalone gate first.
- Made implementation start from independently-green ticket files, with one
  fresh write-delegate charge and one atomic full-gate landing per ticket.
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

### Fixed

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
