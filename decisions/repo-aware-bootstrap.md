# Repo-aware one-command bootstrap (FT76)

## Destination

One public entry point — `npx --yes redbench@<immutable-version> setup`, with
`bench setup` as the installed equivalent — explores a target repo, presents the
smallest Bench adoption plan, confirms it, and writes it transactionally,
composing the existing link/init/setup seams. No second asset installer, no
second source for shared rules. Sources: `RR:A-02`, `RR:A-03`, `RR:A-04`,
`RR:S-09`.

## #1: How does FT76 relate to FT84's transactional lifecycle?

Blocked by: none
Type: Grill

### Question

FT76's "transactionally seeds" overlaps FT84 (staging, preflight, atomic
promote/rollback, idempotent re-runs, downgrade reconciliation) almost
clause-for-clause. Depend on it, absorb it, or duplicate minimally?

### Answer

Depend on FT84 and build FT84 first. FT84 owns write-transaction semantics; the
bootstrap is a composition layer over those seams and adds no second
transaction implementation. (Resolved 2026-07-21.)

## #2: Where is the seam between deterministic CLI and harness conversation?

Blocked by: none
Type: Grill

### Question

Does `bench setup` own explore → present → confirm → write end-to-end, or is it
a non-interactive seeder with all confirmation in the harness conversation?

### Answer

The CLI owns everything mechanical: repo inspection, instruction-file proposal,
asset seeding, and TTY confirmations. The harness conversation
(`/bench-setup-repo`) owns judgment content only — the real gate command,
profile seams and lines, CONTEXT.md. The one-shell-command promise must hold on
any harness or none. (Resolved 2026-07-21.)

## #3: Which instruction file receives the managed Bench block?

Blocked by: none
Type: Grill

### Question

AGENTS.md, CLAUDE.md, or both — and is CLAUDE.md conditioned on detected Claude
Code presence?

### Answer

Always converge, unconditionally, to the kit's own shape: AGENTS.md is
canonical and gains (or is created with) the marker-owned Bench block;
CLAUDE.md is created or extended with only the marker-owned import lines
(`@AGENTS.md`, `@.bench/BENCH.md`). No harness-detection gate — a preserved
CLAUDE.md without the imports is exactly the red doctor state setup must never
leave behind. (Resolved 2026-07-21.)

## #4: Which inspection facts are inferred silently vs confirmed?

Blocked by: none
Type: Grill

### Question

The bootstrap inspects Git state, remotes, harness files, language/build/test
signals, gates, and project docs. Which of those may it act on without asking,
which land in the plan preview, and which become one-at-a-time questions?

### Answer

Nothing is acted on silently. Every inferred fact appears in the plan preview
with its consequence ("Go module detected → proposing `go test ./...` gate").
Only genuine ambiguities — two build systems, an existing foreign gate,
conflicting instruction files — become one-at-a-time questions. One preview +
one confirm is the default interaction. (Resolved 2026-07-21.)

## #5: What is a setup-level conflict, and what does the CLI do with it?

Blocked by: #1
Type: Grill

### Question

Rollback of a failed write is FT84's. What remains here: which pre-existing
content counts as a project-owned conflict (modified managed assets, foreign
gate, competing instruction files), and does the CLI stop, report-and-skip, or
merge?

### Answer

Never merge, never overwrite. Conflicts are named in the plan preview before
any write; on confirm the CLI writes everything non-conflicting, preserves each
conflicting asset as project-owned, and exits with a distinct partial status
plus a machine-readable conflict list — the same posture FT84 gives modified
managed assets and unlink residuals. Resolving ownership is the reviewer's
call. (Resolved 2026-07-21, blanket approval.)

## #6: What repo-shape matrix does v1 advertise?

Blocked by: none
Type: Grill

### Question

Monorepos, repos with no Git yet, offline installs, no-global-install
environments, fresh clones of already-linked repos — which are advertised v1
targets and which are explicit non-goals?

### Answer

v1 advertises: single-repo root adoption in an existing Git repository; empty
and established repos; paths with spaces; pre-existing `AGENTS.md`/`CLAUDE.md`,
hooks, and settings; monorepos as root-level adoption only; no-global-install
environments (npx + repo-local executable path is the primary route); offline
runs from the packed artifact; and idempotent re-runs, including on a fresh
clone of an already-linked repo. Non-goals: repos without Git — setup stops
with the exact `git init` instruction rather than initializing; per-package
monorepo profiles. (Resolved 2026-07-21, blanket approval.)

## #7: What does doctor assert after setup?

Blocked by: #3
Type: Grill

### Question

The red/green contract setup hands to `bench doctor`: the
preserved-CLAUDE.md-without-imports red, the false all-harness success to
prevent, and what a green post-setup report must actually verify.

### Answer

Doctor asserts per harness surface, never one aggregate: the marker-owned
block present in `AGENTS.md`; the import lines present in `CLAUDE.md` (a
preserved `CLAUDE.md` without them is red, never folded into an all-harness
green); the gate executable present and runnable; the project profile present;
the repo-local `bench` resolvable; and every emitted pointer valid against the
installed artifact. `bench setup` ends by running doctor and printing that
honest report, the reload instruction, and the exact next action. (Resolved
2026-07-21, blanket approval.)

## #8: What can `npx redbench setup` do before any local binary exists?

Blocked by: none
Type: Research

### Question

The cold path runs from the wrapper with no pinned binary and possibly no kit
tree (today `link`/`init` refuse outside a real kit source tree). What does the
packed artifact actually provide on a cold `npx` run — postinstall build,
`bench repair` interaction, immutable-version pinning — and what must change so
`setup` can run there? Registry facts live in git history
(`decisions/governed-offline-release-publication-research.md`, retired
2026-07-21).

### Answer

— (open)

## Not yet specified

- Gate/profile template inference: how language/build/test signals map to the
  proposed smallest gate.
- Re-run and downgrade UX beyond the FT84 semantics the bootstrap inherits.
- `bench setup` output shape (plan preview format, AXI/TOON conformance).
- README one-command claim wording and the reload instruction at exit.

## Out of scope

- A second asset installer or a second source for shared rules (destination
  constraint, restated to stop drift).
- FT84's transaction internals — owned by that row, built first (#1).
- Redesigning the `/bench-setup-repo` interview content — this map only moves
  its mechanical preflight into the CLI.
- Host-level environment management (global npm installs, PATH edits beyond the
  repo-local executable path).
