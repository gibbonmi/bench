# Changelog

All notable user-facing changes to Bench are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/2.0.0/), and versions follow
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- The `internal/worktree` test suite now runs its eligible tests in parallel
  under a census. The census parses the package test files and turns the gate
  red when an eligible test omits `t.Parallel()`, or when a test that binds
  environment or changes directory calls it. The `clean`, `reclaim`, `list`, and
  `resume-clean` verbs now take their repository root at the command boundary.
  Their grammar and their output do not change.
- A worktree `bench commit` now runs the fast lane instead of the whole-project
  gate: gofmt, the prose check on the named Markdown, `go vet`, and `go build`.
  A lane pass publishes onto the worktree branch and writes a lane record, never
  a gate verdict. `bench worktree land` stays the one whole-project gate, and a
  repo with no declared lane keeps the full-gate commit.
- Lifecycle identity refusals now name the one component that failed: the request
  token, the assignment state, the assignment path, the owner marker, the worktree
  registration, or the Bench lock. `bench worktree exec` and `bench worktree path`
  now print the resolver's reason instead of one sentence for all of them.
- Claude and Codex now run the shared follow-on guard before Bash calls. It refuses
  outer shell composition after a Bench command, and `bench guards` reports the
  manifest-derived harness wiring.
- `bench worktree land` now runs under one installed promotion broker for the
  complete command. The wrapper selects that broker from the installation
  manifest, and it refuses `BENCH_KIT`, `BENCH_RUN_BINARY`, and `BENCH_WRAPPER`
  for this command. Repository code thus cannot authorize its own publication,
  and the landing no longer rebuilds and re-runs itself.
- The promotion broker composes the prospective tree in private temporary
  storage. It builds the gate executable from that exact tree. It takes the
  phase schedule from the landing destination, so a candidate phase manifest
  cannot omit the checks that grade it. The broker alone accepts the gate
  evidence and updates the destination ref. A red gate leaves the destination
  ref and the project-green marker unchanged, and no temporary tree or
  executable stays after any outcome. A landing whose reviewed diff changes the
  broker source names the `bench repair` or release install step that publishes
  the new broker.
- `bench commit` now formats only changed Go files inside its named paths before
  it composes and runs the unchanged full gate. Dry runs remain non-mutating.
- The `prose-mechanics` check now rejects unsupported, malformed, or unreadable
  learning journals and reports the repair location. `capture/learnings.md` is no
  longer tracked, so each journal remains local.

### Fixed

- `bench worktree land` now treats only `capture/IDEAS.md`,
  `capture/learnings.md`, and `capture/session-handoff.md` as destination-local
  ignored state. The notes survive landing and resume, while foreign ignored
  residue still refuses with a path-specific diagnostic.
- The `internal/worktree` test suite now selects one Bench executable per run,
  routes every repository, descendant, environment, and directory effect
  through one serial journey harness, and tests landing, lifecycle, and
  reclaim policy in three pure packages. The parent package span fell from a
  125.790-second median to a 56.898-second median on the reference WSL host.
  Raw runs and demand counts are recorded in
  `specs/worktree-test-latency/evidence/demand-reduction.md`. The gate driver,
  `-count=1`, and public command behavior are unchanged, and the spec adds no
  `t.Parallel` and no scheduler.
- `bench commit` now refuses the primary checkout and directs users to
  `bench worktree create`. This makes the worktree-only phase rule executable
  at Bench's publication boundary.
- The source wrapper now recovers an executable Go path from a bounded clean
  Bash login when the harness PATH is partial. It preserves the harness PATH and
  gives an actionable refusal when no safe Go path is available.
- Every existing `ROADMAP.md` row outside `## Parked and scheduled work` now
  carries a `Next: <token>` line in its detail file (`roadmap/FT<n>.md`) —
  `shape`, `spec`, `ticket`, `decide`, or `kit-edit` — under the position-anchored
  `Next:` row grammar `bench roadmap --flow` and the `Feeds:` retro marker both
  read. The missing-line fault class is now enforced: the `row-next-grammar`
  and `roadmap-detail-integrity` checks red a row with no `Next:` line and no
  parked-section exemption, closing the gap between the check and the board it
  grades.
- `bench worktree land` now exits 3 when publication succeeded but marker,
  checkout reconciliation, or source release remains incomplete, and its terminal
  record includes a token-safe resume invocation. Pre-publication refusals, usage
  errors, and completed releases retain exits 1, 2, and 0 respectively.
- The `/bench-debug` (`$bench-debug`) phase restores the upstream discipline a
  silent compression lost: Phase 1 points at a shipped ten-entry
  loop-construction menu (`.agents/skills/bench-debug/references/loop-constructions.md`),
  the two hard stop-gates return verbatim ("No red-capable command, no
  Phase 2"; "Do not proceed until you have reproduced and minimised"),
  "Tighten the loop" returns as a named step, the Phase 1 completion criteria,
  Phase 2 confirmations, and Phase 6 close-out return as checkbox forms, and
  Phase 2 gains the reproduction-economics rule — a green proxy only narrows a
  hypothesis. Every Bench-specific addition survives unchanged.
- Every Bench phase now carries a reviewed invocation policy graded on both harness
  surfaces: one table records, per phase, whether the Claude model may reach for the
  command on its own and whether Codex may invoke the adapter implicitly, and the
  gate reds when either surface drifts from it, when a command file has no policy
  row, or when a row names no command file. `$bench-debug` is now implicitly
  invocable on Codex — a reported symptom routes to the bug path without the
  operator naming the phase — and its adapter description carries the symptom
  trigger instead of an explicit-only clause. The inert `disable-model-invocation`
  key leaves every Codex adapter, where `agents/openai.yaml` is the real policy
  surface.
- Bare wrapper and binary invocations now render the `bench status --route`
  recommendation, `/bench` and `$bench` follow that route, help spellings show the
  command inventory, and capture maintenance is now `/bench-drain` (`$bench-drain`)
  with the former phase name retained as a one-release alias.
- The roadmap is now a split board: `ROADMAP.md` keeps its section prose and exactly
  one heading line per row with no body, and each row's body, `Occurrence:` ledger,
  and `Sources:` line move to `roadmap/FT<n>.md`, so a cold read costs the index
  alone and a row edit opens one file. `bench roadmap --context` lists `roadmap/` in
  its sources block and renders a missing detail owner, an orphaned row file, an
  inline body, a heading mismatch, an unrecognized file, a duplicate row ID, and a
  wrapped heading as parse failures over a malformed index; the new Dev-tier
  `roadmap-detail-integrity` conformance check reds the gate on any of them. The
  schema stays 4, and the inline-body layout is no longer accepted.
- The destructive-Git guard no longer intercepts `git stash`; stash lifecycle is
  available for isolating unrelated working-tree state while other destructive Git
  operations remain guarded.
- `bench coverage --check` now refuses a coverage row that references more than four
  stories or states more than one predicate (a `;` outside backticks), and a declared
  story that no row references unless the spec carries a
  `Not covered: story <n> — <reason>` line under the map.
- A spec's acceptance coverage map now carries four columns — story, behavior, seam, and
  why it catches the failure — behind the optional leading row ID; the retired `red
  signal` column's job belongs to the why-it-catches clause. `bench coverage <spec>` now
  projects `rows[N]{story,behavior,seam}` instead of `rows[N]{story,seam,red_signal}`, so
  task seeding reads a column every schema has. The previous six-column header still
  parses.
- `/bench-write-spec` now runs one review round over the spec-and-tickets pair rather than
  two, with the ticket-granularity quiz as its approval step, and `--reviewer <tier>
  [effort]` now takes a tier only — resolved through the invoking harness's own
  `.bench/lines.env` column. A bound model id is an invocation error.

### Fixed

- SessionStart now diagnoses a Go repository whose harness PATH lost the
  toolchain effects promised by its initialization marker. It asks the clean
  Bash login under a two-second process-group bound, validates but never runs
  the reported executable, and prints a shell-quoted PATH prepend while keeping
  session startup informational.
- Workflow: repair re-reviews now block only on accepted repair predicates and
  repair-induced changes, so unrelated findings cannot restart full discovery.
- An adopted repository can now run its own gate green. `bench setup` seeds
  `.bench/gate-inputs.json` declaring `BENCH_HOME` and `HOME`, so the scaffolded
  gate's wrapper call survives the closed gate environment instead of dying on an
  unbound variable, and the scaffolded gate skips inventory validation until
  `tests/canary` exists. An existing `.bench/gate-inputs.json` is left untouched
  and stays out of the link manifest. The sentinel still keeps a fresh stub red
  until the operator removes it.
- `bench link` now adopts only same-file converged adapter symlinks on first link.
- `bench skills-index` now refuses hostile or malformed producer files instead of
  acting on them: a symlinked, special, oversized, or non-UTF-8 `SKILL.md`,
  `.bench/BENCH-reference.md`, or `.bench/consumer-payload.json` is named and
  refused, a skill directory with no `SKILL.md` is reported even when a same-named
  command adapter exists, and a control rune can no longer forge an index line. This
  command keeps every skill under a repository path containing spaces or glob
  characters. `--write` leaves the reference untouched on every refusal, and an
  interrupted or failed replacement leaves no `.bench/.skills-index-*` scratch file
  behind. A missing or non-executable `git` is now reported as the missing tool it
  is, rather than as not being in a repository.
- `/bench-implement-spec` now creates or retains its integration worktree before
  running build preflight, so earlier local `main` commits cannot contaminate the
  spec-owned changed-path check.

### Removed

- Removed the provisional spec-build lifecycle wholesale: the `bench spec build`
  family (`start`, `assign`, `checkpoint`, `integrate`, `review`, `status`,
  `promote`, `abandon`, and the maintainer-run `reclaim`) and
  `bench worktree recovery` are no longer commands and answer with the standard
  unknown-subcommand structured error. `bench resume`'s reconcile now deletes
  the `refs/bench/specbuild/` and `refs/bench/recovery/` namespaces and purges
  lifecycle-typed assignments from the intent ledger at every session start,
  idempotently. Zero backwards compatibility — no shim, no migration tooling;
  reviewed spec-backed builds now keep serial green ticket commits on one
  retained integration source, review its explicit frozen pair, and publish it
  through `bench worktree land` with the spec transition in the gated tree.

### Added

- New `prose-mechanics` conformance check, with its approved exclusions in the
  new `.bench/prose-exclusions` file. `.bench/BENCH.md` now splits into an
  always-loaded core and separate reference material. The kit's guidance and
  comments now carry the ASD-STE100 register throughout.
- Added the `landed` assignment classification and `bench worktree clean --landed`
  plan/apply sweep for retiring landed worktrees as one fingerprinted set.
- New `bench skills-index [--check|--write]` command — the operator's one skills-index
  reader and regenerator; `.bench/skills-index.sh` is gone and every reference now
  names the verb.
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
