# repo-aware-setup (FT76)

Status: staged

Compiled from the closed `decisions/repo-aware-bootstrap.md` map. Four
defaulted decisions — the two Handoff uncertainty flags (non-TTY invocation,
re-run posture) plus two the map left unspecified (the gate-inference table,
the porcelain non-AXI output shape) — are settled in this spec and flagged for
reviewer veto in Implementation decisions. Reviewer routing
2026-07-22: the build runs on cheap-tier (`sonnet`) subagents only; story
lines below name Claude Code aliases, with effort set for the cheap tier's
headroom and any deviation from the profile's cached routings named in-story.
Escalation follows the `craft-line` ladder: a stalled story may bump to `opus`
with a declared line, and any bump to `fable` pauses for reviewer approval.

## Problem

Adopting Bench today takes a documented multi-step path: install the kit, run
`bench link`, run `bench init`, then hand-configure the gate and profile through
a harness conversation. A new user has no single entry point that inspects their
repository, proposes the smallest adoption plan, and writes it safely. The
existing pieces — link's transactional asset lifecycle (FT84), init, doctor,
and the `/bench-setup-repo` interview — exist but are not composed.

## Solution

One public entry point — `bench setup` (and, from the packed artifact, the same
route through the wrapper) — that explores the target repository, presents every
inferred fact and its consequence in a plan preview, asks only genuine
ambiguities one at a time, then transactionally seeds the managed assets,
instruction-file convergence, gate, and profile by composing the existing
link/init seams and FT84's transaction. It ends by running doctor and printing
the honest per-harness report, the reload instruction, and the exact next
action. No second asset installer; no second source for shared rules.

## User stories

1. As a new user, I want `setup` dispatched as an adoption route in the wrapper
   so that one command reaches the Go-core implementation from every shipped
   surface. Line: sonnet / low. This is thin dispatch plumbing at a known
   seam, matching the profile's cached cheap routing for CLI shell plumbing
   exactly.
2. As a user, I want the plan preview to list every inferred fact with its
   consequence ("Go module detected → proposing `go test ./...` gate") on
   stdout before any write, so that nothing is acted on silently. Line:
   sonnet / high. The inspection-to-proposal core is the deep module of this
   feature, so the cheap tier gets its highest effort even though every
   behavior here lands in a contract row.
3. As a user, I want genuine ambiguities — two build systems, a foreign gate,
   competing instruction files — to become one-at-a-time questions, with
   everything else resolved by one preview plus one confirm, so that setup asks
   only what it cannot infer. Line: sonnet / high. Question sequencing decides
   the interaction contract and the gate cannot grade prompt quality, so the
   cheap tier runs at high effort and the review phase grades the prompt text.
4. As a user with pre-existing content, I want conflicts never merged and never
   overwritten: the preview names each conflict, a confirmed run writes
   everything non-conflicting, preserves each conflicting asset as
   project-owned, and exits with the link lifecycle's partial status plus its
   machine-readable conflict list. Line: sonnet / high. Conflict semantics are
   the fail-closed heart of the feature, and high effort compensates the cheap
   tier where a subtle overwrite would be the worst possible defect.
5. As a user on any harness, I want instruction files converged
   unconditionally to the kit's shape — the marker-owned Bench block in
   `AGENTS.md`, only the marker-owned import lines in `CLAUDE.md`, surrounding
   project content preserved — so that no harness-detection gate exists. Line:
   sonnet / medium. The existing marker writer owns the mechanics and every
   matrix cell is gate-observable, which is exactly the fully-covered work the
   cheap tier is for.
6. As a user in a directory without Git, I want setup to stop with the exact
   `git init` instruction and write nothing, so that setup never initializes a
   repository for me. Line: sonnet / low. This is a single fail-closed branch
   with a fully gate-observable contract, matching the cached cheap routing.
7. As a CI or harness-spawned (non-TTY) caller, I want a bare `bench setup` to
   fail closed with no writes while naming `--yes` and `--plan`; `--plan`
   prints the preview and exits zero without writing, and `--yes` auto-confirms
   only an ambiguity-free plan, refusing with the open questions named
   otherwise. Line: sonnet / medium. The flag matrix is precisely specified
   and fully contract-observable, so the cheap tier at medium effort covers
   it.
8. As a returning user, I want a re-run on an already-linked repository to
   converge and report — idempotent and downgrade-aware through FT84's relink
   semantics, with an unchanged tree reported as already converged — so that
   setup is safe to repeat, including on a fresh clone of a linked repo. Line:
   sonnet / medium. The FT84 transaction owns the hard part and the re-entry
   contract is pinned by observable cells, so medium effort suffices.
9. As a user, I want the proposed gate derived from a small deterministic
   detection table, with multiple matches becoming one question and a
   zero-signal repository receiving a fail-closed gate stub (a gate that exits
   non-zero naming its configuration step) and a red doctor row, so that setup
   never fabricates a green oracle. Line: sonnet / high. The table is
   mechanical but the zero-signal posture guards the oracle, and
   oracle-adjacent work gets extra effort when it runs on the cheap tier.
10. As a user finishing setup, I want the command to end by running doctor and
    printing that honest report, the harness reload instruction, and the exact
    next action, so that a cold user knows precisely what to do next. Line:
    sonnet / medium. The exit contract is string-assertable end to end, with
    the doctor-verdict tracking pinned by the red and green fixtures.
11. As a user on any harness, I want doctor to assert per-harness rows — the
    `AGENTS.md` marker block, the `CLAUDE.md` import lines (a preserved
    `CLAUDE.md` without them is red, never folded into an aggregate green), the
    gate present and executable, the profile present, the repo-local `bench`
    resolvable, and every setup-emitted pointer valid — so that a false
    all-harness success is impossible. Line: sonnet / medium. Every doctor row
    is pinned by its own fixture, so the cheap tier at the profile's mid
    effort for gate-adjacent logic covers it.
12. As a cold new user, I want the packed artifact to carry the setup route so
    that an offline install into a spaced prefix can run setup end to end, with
    everything durable landing repo-local because the npx cache copy is
    ephemeral. Line: sonnet / high. The artifact leg is the acceptance seam of
    the whole feature and meets the most hostile environment matrix, so it
    gets the highest effort the cheap tier can spend.
13. As a reviewer, I want `/bench-setup-repo` slimmed to judgment content only
    — the real gate command, profile seams and lines, `CONTEXT.md` — with its
    mechanical preflight moved into the CLI, so that one source owns each half.
    Line: sonnet / high. This is command prose the profile's leverage override
    would send to the top tier; it lands on sonnet only under the reviewer's
    2026-07-22 all-cheap routing, the deviation is named here, and the review
    phase must grade the prose the gate cannot.
14. As a new user reading the README, I want the quickstart to state the
    one-command claim in its currently-true form — the git-dependency `npx`
    form requiring a Go toolchain until first publish, with the pinned
    `npx redbench@<version>` form documented as activating after first publish
    — so that the entry point never overclaims. Line: sonnet / medium. This is
    short prose whose truthfulness matters more than its polish, deviating
    from the top-tier doc routing under the same reviewer decision, with the
    review phase as the honest grader.

## Implementation decisions

- **Wrapper dispatch.** `setup` joins the wrapper's adoption routes exactly as
  `link`/`init`/`doctor` do, forwarding to the Go core's adopt dispatcher. No
  logic in the wrapper.
- **One deep module.** The Go-core setup command owns inspection → proposal →
  preview → confirm → write orchestration inside the adopt package. It composes
  the existing init, link, and FT84 transaction seams; it adds no second
  transaction implementation and no second asset installer (map destination
  constraint, restated to stop drift).
- **Inspection facts.** Git presence, instruction files, existing gate, harness
  files, and build-system signals. Nothing is acted on silently: every inferred
  fact appears in the preview with its consequence. Only genuine ambiguities
  become questions; one preview + one confirm is the default interaction.
- **Gate inference table (defaulted decision — map left this as fog; veto
  here).** v1 detects: `go.mod` → `go test ./...`; `package.json` with a
  `test` script → `npm test`; `Cargo.toml` → `cargo test`; a `Makefile` with a
  `test` target → `make test`. Multiple matches → one question. Zero matches →
  a fail-closed gate stub that exits non-zero naming its configuration step,
  plus a red doctor row; a fabricated green gate is never written.
- **Non-TTY posture (settles Handoff uncertainty flag #1 — veto here).**
  Prompting requires a TTY on stdin. Non-TTY without flags: fail closed, no
  writes, distinct exit status naming `--yes` and `--plan`. `--plan`: preview
  only, exit 0, no writes. `--yes`: auto-confirm an ambiguity-free plan;
  refuse (no writes) naming each open question otherwise. Conflicts are not
  ambiguities: a `--yes` run over a conflicted but ambiguity-free plan proceeds
  exactly as an interactive confirm would — non-conflicting assets written,
  conflicts preserved, partial exit status with the conflict list — because
  nothing is ever overwritten and the status surfaces what a human didn't see.
- **Re-run posture (settles Handoff uncertainty flag #2 — veto here).**
  Converge-and-report. A re-run previews the delta (or reports already
  converged), then converges through FT84's relink reconciliation, including
  downgrade handling. Preview-only re-runs are what `--plan` is for.
- **Instruction-file convergence.** Always, unconditionally, to the kit's own
  shape via the existing marker writer: `AGENTS.md` gains (or is created with)
  the marker-owned Bench block; `CLAUDE.md` is created or extended with only
  the marker-owned import lines. No harness-detection gate. Unclosed fences and
  malformed markers remain conflicts, exactly as the marker scanner rules
  today.
- **Conflict semantics.** Reuse the link lifecycle's verdict vocabulary and
  partial posture verbatim — same conflict block shape, same exit-status
  meaning. Resolving ownership stays the reviewer's call.
- **Prompt I/O seam.** Confirm/question I/O is injected (reader/writer) so
  sequencing is unit-testable at the adopt package seam without a pty
  dependency; black-box contracts drive the non-TTY matrix. The production
  wiring keeps the real-TTY glue to one constructor line binding stdin, which
  stays untested absent a pty library — a named residual risk. Adopting a pty
  library instead would be a new dependency outside the precedent shape — a
  reviewer decision this spec does not make.
- **Exit contract.** Setup ends by invoking doctor, then prints the reload
  instruction and the exact next action (`/bench-setup-repo` conversation for
  judgment content, or the gate-configuration step when the stub was written).
  Setup exits zero only when doctor reports green; a red doctor row (the
  zero-signal stub path, a preserved `CLAUDE.md` conflict) makes setup exit
  with the partial status while still printing the full report and next
  action — the exit code tracks the doctor verdict, never a fabricated print.
- **Doctor.** Per-harness red/green rows, never one aggregate; rows enumerated
  in story 11. Extends the existing doctor report renderer.
- **Output shape (defaulted decision — map left this unspecified; veto
  here).** Setup is adoption porcelain, following link's precedent:
  human-readable preview text, the machine-readable conflict block on partial.
  It is not an AXI query command and emits no TOON tables.
- **`/bench-setup-repo` slimming.** The command keeps judgment content only;
  its mechanical preflight moves into the CLI. The command edit lands in the
  same diff as the CLI change, keeping the stale-command-reference sweep green.

## Testing decisions

- A good test here runs the built command against a throwaway fixture repo and
  asserts exit status, stdout text, and resulting file/git state — external
  behavior at the shipped surface, never adopt-package internals. Prior art:
  the link, doctor, and unlink surface contracts and the artifact contract's
  offline install-and-run legs.
- Seams tested: (1) the `bench setup` surface contract, (2) the packed-artifact
  cold-run leg, (3) the doctor per-harness rows, plus one unit seam for prompt
  sequencing with injected I/O. New contract files respect the 400-line
  structure budget by splitting along fixture families.
- Gate: `.bench/gate.sh` (the project gate; the new contracts join the existing
  contract phase).

### Seam diagram

Seam 1 — setup surface contract:

    trigger: user/CI runs `bench setup [--plan|--yes]` in a fixture repo
        │
        ▼
    repo state (git, files)  ──▶  [ adopt setup: inspect → propose   ]  ──▶  preview text, prompts
    kit assets (BENCH_KIT)   ──▶  [ → confirm → FT84 transaction      ]  ──▶  exit status
    flags, stdin TTY-ness    ──▶  [ → doctor → next action            ]  ──▶  converged repo tree
                      ◀ tests attach here: run built binary in throwaway repos;
                        assert stdout, exit code, and file/git state

Seam 2 — packed-artifact cold run:

    trigger: offline `npm install` of staged tarballs into a spaced prefix
        │
        ▼
    wrapper+native tarballs  ──▶  [ installed wrapper → setup route ]  ──▶  converged target repo
    empty target repo        ──▶  [ (same code path as seam 1)      ]  ──▶  repo-local launcher
                      ◀ tests attach here: artifact contract installs, runs
                        setup, asserts durable state landed repo-local

Seam 3 — doctor per-harness rows:

    trigger: `bench doctor` after setup (or against a broken fixture)
        │
        ▼
    repo file state  ──▶  [ doctor row evaluators ]  ──▶  per-harness red/green rows, exit status
                      ◀ tests attach here: fixture repos per row state;
                        assert row text and exit code

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench setup` reaches the Go core through the wrapper and converges an empty repo | setup surface contract | new contract running `bench setup --yes` in a fresh git repo (and once from a subdirectory) asserting exit 0 and each converged asset by name — the `AGENTS.md` marker pair, the `CLAUDE.md` import lines, `gate.sh`, the profile, the link manifest, and the repo-local launcher — red before the route exists | an unimplemented route, a mis-anchored cwd, or a partial-tree writer each fail at least one enumerated asset assertion |
| 2 | preview lists every inferred fact with consequence, before any write | setup surface contract | contract asserting the preview text across two fixtures whose facts and consequences differ, and an untouched tree when the run stops at preview, red before implementation | a write-first implementation fails the untouched-tree assertion and a hardcoded preview string fails one of the two differing fixtures |
| 3 | ambiguity-free plans need one confirm; ambiguities are asked one at a time | prompt-sequencing unit seam + setup surface contract | unit test with injected I/O asserting question order and single-confirm flow; contract with a two-build-system fixture under `--yes` refusing with the open question named — both red first | a prompt-free or ask-everything implementation fails the sequencing assertions; silent ambiguity resolution fails the refusal contract |
| 4 | conflicts preserved, non-conflicting written, partial status + conflict block emitted | setup surface contract | contract with a foreign gate plus a modified managed asset asserting partial exit, conflict block content, preserved conflicting files, and the non-conflicting assets present in the same run, red first | merge-or-overwrite destroys preserved content, a write-nothing stop-on-conflict fails the non-conflicting-written assertion, and a full-success exit hides the conflict |
| 5 | instruction files converge across the file-state matrix with project content preserved | setup surface contract | contract matrix: absent/empty/no-block `AGENTS.md`, absent/empty/no-import `CLAUDE.md`, missing trailing newline, unclosed fence → conflict; red first | each degenerate writer (overwrite, append-duplicate, fence-blind) fails at least one matrix cell |
| 6 | gitless directory: stop, name `git init`, write nothing | setup surface contract | contract in a plain directory asserting non-zero exit, the `git init` instruction, and no files created, red first | an implementation that initializes or seeds anyway fails the no-writes assertion |
| 7 | non-TTY matrix: bare fail-closed naming flags; `--plan` previews at exit 0; `--yes` refuses ambiguous plans and proceeds to partial on conflicted ambiguity-free plans | setup surface contract | contract driving all four invocations non-interactively — bare, `--plan`, `--yes` over an ambiguous fixture, `--yes` over a conflicted ambiguity-free fixture — asserting exits, write behavior, and named flags, red first | hanging on a prompt, writing without confirmation, auto-answering questions, or overwriting conflicts under `--yes` each fail their cell |
| 8 | re-run converges and reports: unchanged trees report already converged, deltas reconcile through FT84 relink | setup surface contract | contract running setup twice asserting second exit 0, already-converged report, and identical tree hash, plus a delta cell where a managed asset was modified or removed between runs asserting relink reconciliation, red first | a non-idempotent second run changes the tree hash, and an existence-check degenerate ("`.bench` exists → exit 0") fails the delta-reconciliation cell |
| 9 | detection table proposes the right gate and writes it; zero-signal writes a fail-closed stub and a red doctor row | setup surface contract | contract fixtures for every table cell (`go.mod`, `package.json` test script, `Cargo.toml`, `Makefile` test target) asserting the proposed command lands in the written `gate.sh`; a zero-signal fixture asserting the stub exits non-zero, doctor goes red, and setup exits partial, red first | a table hardcoding one ecosystem fails the other cells, and a fabricated default gate would execute green and fail the stub-exits-non-zero assertion |
| 10 | setup ends with the doctor report, reload instruction, and exact next action, exit tracking the doctor verdict | setup surface contract | contract asserting all three trailing output elements after a converged green run at exit 0, and partial exit on the zero-signal red-doctor fixture, red first | a fabricated doctor print cannot track the verdict across the green and red fixtures, and skipping the handoff text fails the output assertion |
| 11 | doctor asserts per-harness rows; preserved `CLAUDE.md` without imports is red | doctor per-harness rows | contract fixtures per row state, including the imports-stripped `CLAUDE.md` asserting a red row and non-zero exit, red first | an aggregate-green doctor passes the broken fixture and fails the red-row assertion |
| 12 | packed artifact runs setup offline from a spaced prefix; durable state lands repo-local | packed-artifact cold run | new artifact-contract leg installing the staged tarballs offline and running setup, asserting convergence and the repo-local launcher, red first | a route missing from the packed dispatch, or state left in the npx cache copy, fails the leg |
| 13 | `/bench-setup-repo` keeps judgment content; mechanical preflight lives only in the CLI | conformance sweep (existing) | stale-command-reference sweep already covers dangling references (already covered); the judgment/mechanical split is prose — not TDD-able, graded at review | reference rot is the automatable failure and the sweep already bites on it |
| 14 | README states the currently-true one-command entry | none (prose) | not TDD-able — wording truthfulness is graded at `/bench-review-implementation` against the artifact contract's actual behavior | no mechanical oracle exists for claim accuracy; review is the honest owner |

### Edge inventory

- Paths with spaces/glob characters — covered by the artifact leg's spaced
  prefix (row 12) and existing wrapper routing contracts.
- Absent vs present-but-empty instruction files — both cells asserted in the
  row-5 matrix.
- Hand-edited file lacking a trailing newline — a row-5 matrix cell.
- Unclosed code fence around markers — conflict, a row-5 matrix cell reusing
  the marker scanner's existing rule.
- Special file (FIFO) at a fixed inspection path — treated as a conflict, not
  read; asserted as a cell in the row-4 conflict contract.
- cwd deeper than repo root — setup anchors at the repository root; asserted
  as a cell in the row-1 contract (run from a subdirectory).
- **Won't handle:** real-TTY confirm wiring — one constructor line binding
  stdin, untestable without a pty library, which is a dependency decision this
  spec leaves to the reviewer; the residual risk is named in Implementation
  decisions.
- Re-run idempotency (second setup, fresh clone of linked repo) — rows 8, 12.
- No global `bench` on PATH — the artifact leg runs through the installed
  wrapper and repo-local launcher only (row 12).
- Non-TTY stdin — row 7.
- **Won't handle:** control bytes in git-sourced text — the preview renders
  fixed fact strings and repo-relative paths, not free git text; the existing
  sanitize layer owns any future free-text surface.
- **Won't handle:** interrupt mid-write — crash rollback is FT84's transaction
  contract, already covered by the lifecycle contracts setup composes.
- **Won't handle:** invocation through a symlinked wrapper — path resolution is
  the wrapper's existing routing contract, unchanged by adding a route.
- **Won't handle:** missing git binary on PATH — same environment floor as
  link/init today; setup adds no new dependency on it.

## Out of scope

- **Per-package monorepo profiles** — a different inspection and profile model,
  not more of this one (map #6 non-goal) — 40 edits, 15 gate runs.
- **Post-publish README flip to `npx redbench@<version>`** — blocked on first
  npm publish, a separate release capability — 3 edits, 2 gate runs.
- **Host-level environment management** (global npm installs, PATH edits beyond
  the repo-local executable path) — map out-of-scope; ongoing non-goal.
- **`/bench-setup-repo` interview content redesign** — the map moves only the
  mechanical preflight; the judgment interview is a separate authoring effort —
  10 edits, 3 gate runs.
