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

The cold path already works for `link`/`init`/`doctor`; `setup` only needs a
dispatch route plus the Go-core command. Evidence, each claim owned by a
runnable gate contract:

- The packed wrapper is a valid kit: `files[]` ships `.agents/`,
  `.bench/BENCH.md`, and the adapter tree, so the `adoption_route` kit check
  passes from an installed artifact — proven end-to-end by the artifact
  contract, which installs the wrapper+native tarballs offline into a
  spaced prefix and runs `link`, `doctor --fix`, and an idempotent relink
  (`internal/contract/surface/artifact/artifact_test.go:135`).
- A cold install gets its binary from the platform package: the staged wrapper
  carries `optionalDependencies` for the exact target matrix
  (`internal/contract/surface/package_test.go:172`), with the pinned
  `binary-pins.json` repair path as fallback
  (`bin/bench-repair-binary.mjs:14`, manifest generated by
  `scripts/build-artifacts.sh:115`). No implicit fetch: absent both, the
  wrapper names `bench repair` and exits 127 (`bin/bench.sh` route_binary).
- `npx redbench@<version> setup` runs the `redbench` bin (npx selects the bin
  matching the package name; both `bench` and `redbench` are declared) and
  registry versions are immutable — the pin the destination's
  `@<immutable-version>` relies on (retired publication research, git
  history).
- What must change: add `setup` to the wrapper dispatch as an adoption route
  and implement the Go-core command; nothing structural in the artifact.
  Post-setup resolution must be repo-local (`.bench/bin/bench.sh`, already
  installed by link — `internal/contract/surface/link_test.go:67`) because the
  npx cache copy is ephemeral.
- Domain fact: public npm release is NO-GO today, so until first publish the
  only working cold entry is the git-dep form (`npx github:gibbonmi/bench#main`),
  which runs `prepare` and needs a Go toolchain — exactly what
  `/bench-setup-repo` documents now. (Resolved 2026-07-21.)

## Handoff

1. **Module boundaries.** (a) Wrapper dispatch: `setup` as a new adoption
   route in `bin/bench.sh`, thin forwarder. (b) Go-core setup command:
   inspection, plan preview, ambiguity questions, and write orchestration —
   composing the existing link/init seams and FT84's transaction seams.
   (c) Instruction-file convergence: the marker-owned AGENTS.md block and
   CLAUDE.md import lines, owned by the existing link asset writer. (d)
   Doctor: post-setup per-harness assertions extending the existing command.
   (e) `/bench-setup-repo`: judgment content only; its mechanical preflight
   moves into the CLI.
2. **Contracts.** `bench setup`: exit 0 on full convergence; distinct partial
   status with a machine-readable conflict list (#5); plan preview on stdout
   before any write; prompts only for genuine ambiguities (#4); ends by
   running doctor and printing the reload instruction and exact next action
   (#7). Doctor: per-harness red/green rows, never one aggregate.
3. **Deep vs thin.** The Go-core setup command is the deep module (inspection
   → proposal → transactional write). The wrapper route, the doctor rows, and
   the slimmed `/bench-setup-repo` preflight are thin; no seam of their own.
4. **Black-box assertables.** Exit codes; preview text; target-repo file and
   git state (exactly one managed marker pair in AGENTS.md, import lines in
   CLAUDE.md, gate.sh, profile, link manifest, repo-local launcher); doctor
   output — all assertable in surface/artifact contracts as today's link
   contracts do.
5. **Gate attachment.** New setup contracts in `internal/contract/surface/`
   plus a packed cold-run leg in the artifact contract; the gate already runs
   all contract packages as one phase. No seam is gate-invisible.
6. **Hostile-input owners.** Spaced/metachar paths and offline installs — the
   artifact contract; modified/conflicting managed assets — link conflict
   contracts plus FT84 semantics; foreign gates and competing instruction
   files — the setup preview/partial-status contract; preserved CLAUDE.md
   without imports — the doctor contract.
7. **Uncertainty flags.** Non-TTY invocation (CI, harness-spawned): whether
   setup fails closed, auto-confirms a clean plan, or takes a flag — not
   settled; spec-writer escalates. Re-run on an already-linked repo:
   converge-and-report vs preview-only — not settled. Gate/profile inference
   detail remains fog (below).
8. **Rejected alternatives.** Absorbing FT84 (#1); a non-interactive-only CLI
   (#2); harness-conditioned CLAUDE.md (#3); silent inference (#4); merging
   conflicts (#5); `git init` on gitless repos and per-package monorepo
   profiles (#6); a second asset installer or second source for shared rules
   (destination).
9. **Domain watch-outs.** Public npm release is NO-GO, so the
   `npx redbench@<version>` entry activates only after first publish; the
   git-dep form needs a Go toolchain. The npx cache kit copy is ephemeral —
   everything durable must land repo-local. Registry versions are immutable;
   the entry point must always name a pinned version.

Dependency order: FT84 transactional lifecycle → Go-core setup + wrapper
route → doctor assertions + `/bench-setup-repo` slimming.

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
