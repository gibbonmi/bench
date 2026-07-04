# go-doctor-link-port

Status: implemented

## Problem

Three of the kit's adoption subcommands still carry all their logic in sourced
shell. `bench-link.sh` (373 lines) is the highest-stakes mutator in the kit: it
builds an install plan, runs a preflight that must refuse before touching a single
reviewer-owned file, rewrites the fence-aware managed block in `AGENTS.md`,
generates the `bench:managed-pre-push` hook, and writes the `rel<TAB>sha256`
link manifest that is the ownership test for every managed file. `bench-doctor.sh`
(200) reports the PATH shim's health across four states and atomically repairs it,
and must never execute the contents of a hostile shim it reads. `bench-init.sh`
(105) scaffolds the sentinel-red gate, the canary harness, and the learnings
journal. All three are sourced into `bin/bench.sh` and dispatched as shell
functions.

The strangler port has reached these three (slice 6 of `decisions/go-rewrite.md`).
Their contracts are a strong black-box net — `gate-link-contracts.sh` (13 link
cases plus the init scaffolding stories) and `gate-doctor-contracts.sh` (11
sandboxed doctor blocks, ~14 cases) all drive the CLI as a subprocess — but the logic they
guard is shell, so the fence-aware marker parse, the manifest encoding, the
adapter-symlink target computation, and the shim readback are all
shell-untestable except through a full CLI round-trip. Two new capabilities the
parent map committed to this slice have no home yet: `link` does not stamp the
installed kit version, so consumer repos accumulate unstamped manifests, and
`doctor` cannot detect binary↔asset version skew. And the planted
`.bench/bin/bench.sh` a consumer receives still sources these files, so it
half-works at adoption subcommands it has no asset tree to satisfy.

## Solution

The adoption logic moves into one new Go package, `internal/adopt`, behind three
binary subcommands (`link`, `init`, `doctor`), mirroring how `gitguard` absorbed
the git-guard analyzer into one deep package split by concern. The three `.sh`
files are **deleted** — unlike the hook shims (which the harness invokes by path),
nothing invokes these by path; `bin/bench.sh` routes all three through the
strangler dispatch to the binary. After this slice `bin/` keeps only the
dispatcher (`run_gate`), `bench-worktree.sh` (slice 7), and `bench-postinstall.sh`
as shell.

The binary is the deep unit: plan build, preflight, fence-aware marker rewrite,
pre-push emission, manifest read/write/stamp, the scaffolds, and the shim template
plus manager-owned target-dir selection. `bin/bench.sh` stays a thin router with
one new responsibility — it passes its resolved kit directory to the binary
through the existing `BENCH_KIT` env override, and it refuses adoption subcommands
when it is the *planted* wrapper rather than the real kit.

Two new behaviors land, both owned by this slice per the map. `link` stamps the
installed kit version into the manifest as a comment-prefixed `#kit<TAB><version>`
row that can never parse as a file row; a manifest without that row reads as
skew-unknown. `doctor` gains a skew check: it compares the stamped version against
`bench version` and reports skew, warning (never failing) when the stamp is
absent. The session-start half of the skew surface stays out — it belongs to the
already-landed slice-5 hooks spec.

The two contract fragments are the port-parity net and run with their assertions
intact, re-pointed only where the map requires: the safe-fresh link contract drops
its `[ -f .bench/bin/bench-link.sh ]` planted-file assertion (the file no longer
exists to plant) and gains its inverse plus the stamp, skew, and
planted-refuse rows. New `go test` tables cover what the shell never could: the
marker/fence parser edges, the manifest parse (stamp row, absent file, duplicate
rel), the adapter symlink-target computation, and the doctor dir selection and
shim content/readback round-trip.

## User stories

1. As a reviewer adopting the kit, I want `bench link` ported to `internal/adopt`
   with every observable behavior preserved — mode `copy`/`symlink` (ephemeral
   kit path downgrades `symlink`→`copy` with a note), the plan build over the kit
   asset tree, a preflight that refuses on every conflict class — the 7 classes:
   malformed managed markers (an unclosed marker-bearing fence is a sub-case of
   this class, not a separate one), foreign pre-push, missing kit asset, symlink
   parent, non-directory parent, modified managed file, project-owned file — with
   **nothing mutated before all preflight passes**, the fence-aware `AGENTS.md`
   managed-block rewrite, the `CLAUDE.md` retrofit-or-leave, the plan install
   (copy/symlink/adapter-symlink), the atomic single-writer `rel<TAB>sha256`
   manifest, and the `bench:managed-pre-push` generation honoring `core.hooksPath`
   and the worktree `.git`-file layout — so that the kit's highest-stakes mutator
   runs from the Go core with byte-identical safety. Line: claude-opus-4-8 / high.
   This is the kit's most dangerous mutator — a silent divergence in preflight
   ordering or a conflict class writes a reviewer-owned file the shell version
   would have refused — and although the 13 link contracts grade the observable
   behavior, the refuse-before-mutate guarantee is the safety surface, so it takes
   the mid tier at high effort.

2. As a reviewer bootstrapping a linked repo, I want `bench init` ported to
   `internal/adopt` — scaffolding the shared canary runner, the sentinel-red
   `.bench/gate.sh`, the seed canary fixture, and the learnings journal, each
   `[[ ! -e ]]`-guarded so a second `init` never clobbers a configured file — so
   that the scaffold step runs from the binary with its idempotency intact. Line:
   claude-sonnet-5 / medium. This is a mechanical port of byte-exact scaffold
   content whose every case (the gate red-until-configured sentinel, the
   second-init no-clobber) is pinned by an existing contract, so the gate fully
   grades it and the cheap tier fits; medium effort because the scaffolded gate
   text is load-bearing and a divergence flips a fresh gate's red-until-configured
   guarantee.

3. As a reviewer keeping the PATH shim healthy, I want `bench doctor` ported to
   `internal/adopt` — the four report states (healthy exit 0; stale, foreign, and
   missing each exit 1) reading the shim back by its `# bench-target:` comment and
   **never executing its contents**, and `--fix` selecting the first writable
   non-manager-owned PATH dir (fallback `~/.local/bin`), writing a
   `%q`-quoted shim that targets the resolved wrapper, idempotent when already
   current, refusing a foreign file byte-identically, and atomic via temp+rename —
   so that shim health and repair run from the Go core. Line: claude-opus-4-8 /
   high. Doctor mutates a file on the user's PATH and reads a potentially hostile
   shim, so the no-execute-on-report and atomic-refuse-foreign guarantees are
   safety-critical; the 11 doctor contracts grade the states, but the safety
   posture takes the mid tier at high effort.

4. As a reviewer auditing a linked repo, I want `bench link` to stamp the
   installed kit version into the manifest as a leading `#kit<TAB><version>` row
   that the file-row parser skips (so it can never be read as a managed file), the
   manifest keeping exactly one writer — so that every freshly linked repo records
   which kit version installed it. Line: claude-opus-4-8 / medium. This is a new
   behavior on the ownership-test surface where a stamp row that parsed as a file
   row would corrupt the manifest, but the encoding is contained and gate-graded
   by a new stamp contract row plus a Go manifest-parse table, so it takes the mid
   tier at medium effort.

5. As a reviewer running a possibly-skewed kit, I want `bench doctor` to compare
   the stamped kit version against `bench version` and report version skew,
   treating an absent stamp (a pre-stamp manifest) as skew-unknown — a warning
   that never changes the exit code — so that a binary↔asset mismatch surfaces at
   the doctor the map designated as its detector. Line: claude-opus-4-8 / high.
   The warn-never-fail posture on the skew-unknown branch is the whole guarantee —
   a skew check that failed on a legitimate pre-stamp manifest would brick doctor
   for every not-yet-relinked repo — and while a contract grades the two branches
   by exit and output presence, the exact skew-message wording it cannot, so it
   takes the mid tier at high effort.

6. As a reviewer in a consumer repo, I want the planted `.bench/bin/bench.sh` to
   refuse `link`, `init`, and `doctor` with a pointer at the real kit rather than
   half-working — the kit wrapper passing its resolved kit dir and its own resolved
   path to the binary via `BENCH_KIT`/`BENCH_WRAPPER`, and the dispatcher detecting
   the planted case (no source asset tree beside it) and refusing before routing —
   so that adoption subcommands run
   only through the kit-installed wrapper that can actually reach the assets. Line:
   claude-opus-4-8 / high. The wrapper→binary kit-dir crossing is the map's
   gate-blind uncertainty surface (it is exercised for real only under a global
   npm install), and a planted wrapper that half-ran `link` from a `.bench/` dir
   with no assets is a worse failure than a clean refusal, so it takes the mid
   tier at high effort.

7. As the strangler dispatch and the gate, I want `bin/bench.sh` to route `link`,
   `init`, and `doctor` through the adopt guard to the binary, the three sourced
   `.sh` files deleted with their source lines removed (`bench-worktree.sh` stays,
   slice 7), and both contract fragments re-pointed — the safe-fresh link
   contract's `[ -f .bench/bin/bench-link.sh ]` assertion replaced by its
   not-planted inverse, and the stamp/skew/planted-refuse rows added — with no
   existing assertion weakened and no dangling reference left in contracts,
   README, or link/package — so that the port leaves the dispatcher, the gate
   load, and the stale-reference sweep green. Line: claude-opus-4-8 / high.
   Touching the dispatcher and the gate contracts is the worst defect class in
   this kit (`craft-gate`) — re-pointing a contract without weakening its
   assertion, and deleting three sourced files without a dangling reference — so
   it takes the mid tier at high effort.

## Implementation decisions

- **One package, `internal/adopt`, split by concern (decision — map-silent).**
  The map left the split ("one package vs link/init/doctor apart") to the spec.
  One package, following the `gitguard` precedent of a single deep package split
  into files by pipeline stage: `link.go` (plan build, preflight conflict
  classification, install), `doctor.go` (states, fix, shim template + manager-owned
  dir selection), `init.go` (scaffolds), `manifest.go` (the `rel<TAB>sha256` read/
  write plus the `#kit` stamp), `marker.go` (the fence-aware `AGENTS.md` marker
  parse and rewrite), and `adopt.go` (the public dispatch of the three subcommands
  plus the shared kit-dir/repo-root resolution). Three separate packages would
  duplicate the manifest and marker types that all three subcommands share; the
  shared adoption substrate (kit-dir crossing, manifest, marker parse) is the
  reason one package is correct. Filesystem and git truth are injected at the
  package boundary (a `gitguard.Checker`-style seam) so preflight classification
  and dir selection are unit-testable without mutating a real tree.

- **The wrapper→binary context crossing rides env — both halves (decision —
  map-silent; env, not a flag).** Map item 2 requires two things to cross
  explicitly: the **kit dir** (the asset-source root) and the **wrapper's resolved
  path** (the shim target doctor's `--fix` writes). Both go through env, not flags,
  preserving the required `BENCH_KIT` override:
  - **Kit dir → `BENCH_KIT`.** The wrapper resolves its kit dir (`kit_dir` today)
    and exports `BENCH_KIT`; the binary reads it as the asset-source root. This is
    the one env var already honored by `bench-link.sh`
    (`kit="${BENCH_KIT:-$(kit_dir)}"`), so the override is preserved by
    construction.
  - **Wrapper resolved path → `BENCH_WRAPPER`.** The binary cannot derive which
    wrapper invoked it, and doctor's shim must target the *resolved wrapper*, never
    the platform binary (map item 2, watch-out #9 — the stale-shim failure class).
    So the wrapper resolves its own real path (`resolve_script_path`, symlink-
    walked, no `readlink -f`) and exports `BENCH_WRAPPER` before `route_binary`;
    doctor's shim-content builder uses it as the `%q`-quoted target. Absent (a
    direct binary call outside the wrapper) → doctor falls back to its own argv[0]
    resolution, degrading rather than writing a shim to nowhere.
  No new flag surface is added on either half.

- **The manifest stamp is a leading `#kit<TAB><version>` comment row (decision —
  map-silent encoding).** The map required the stamp to "never parse as a file
  row." The manifest is keyed by rel path (`$1 == rel`), and every real rel is a
  concrete `.bench/`/`.claude/`/`.agents/` path, so a `#kit` key can never collide;
  the manifest parser additionally skips any `#`-prefixed line, so the stamp is
  inert to the ownership test. `doctor` reads the `#kit` row for its skew check.
  A manifest with no `#kit` row (linked before this slice) reads as skew-unknown —
  `doctor` warns and never fails.

- **The planted-dispatcher refusal is a shell guard on the asset tree (decision —
  map-silent mechanism).** The map required the planted `.bench/bin/bench.sh` to
  refuse adoption subcommands (map #3), keeping the refusal in the dispatcher, not
  the binary. The planted wrapper's `kit_dir` is `.bench/`, which carries
  `dist/bench`, `hooks/`, `lib/`, and `BENCH.md` but **not** the source asset tree
  (`.agents/commands`, `AGENTS.md` live at the repo root of a real kit, and inside
  the npm package for a global install). So the dispatcher guards `link`/`init`/
  `doctor` with a cheap asset-tree sentinel (`[ -d "$kit/.agents/commands" ] &&
  [ -f "$kit/AGENTS.md" ]`): present → `route_binary`; absent → refuse
  non-zero with a pointer at the real kit, mutating nothing. The kit's own
  wrapper and a global npm install both satisfy the sentinel and route through.

- **Strangler dispatch and deletion in one change.** `bin/bench.sh` changes
  `link)`/`init)`/`doctor)` from sourced shell functions to the adopt guard →
  `route_binary "$@"`, and drops the three `. "$BENCH_BIN_DIR/bench-{link,init,
  doctor}.sh"` source lines (the `bench-worktree.sh` source stays — slice 7).
  `bin/bench-link.sh`, `bin/bench-init.sh`, and `bin/bench-doctor.sh` are deleted.
  Because `link`'s plan copies whatever `bin/*` remains, deleting the three files
  from the kit automatically stops them being planted into a consumer's
  `.bench/bin/` — the contract's planted-file assertion flips to its inverse. All
  in one change so the gate load, the contract fragments, and the docs
  stale-reference sweep stay green together.

- **Gate resolution and the tree hash stay in their one home.** Unchanged by this
  slice: gate resolution (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) stays in
  `bin/bench.sh`, and `bench version` remains the one version source doctor's skew
  check reads. The pre-push hook *text* is now emitted by the Go binary, byte-
  unchanged from the shell heredoc; `bench-postinstall.sh` and the generated
  pre-push stay shell (map #4).

## Testing decisions

- **What a good test is here.** Acceptance drives each subcommand end-to-end
  through the CLI as a subprocess — run `bench link`/`init`/`doctor` against a
  sandbox tree and assert stdout, exit code, and the resulting filesystem — never
  Go internals. The two existing contract fragments already do exactly this and
  are the port-parity net; they run with their assertions intact, re-pointed only
  where the map requires (the dropped planted-file assertion, the new
  stamp/skew/planted-refuse rows). Go table tests are additional, at the pure-
  function seam in `internal/adopt`, where the shell-untestability tax on the
  marker parse, the manifest encoding, the adapter-target computation, and the
  shim readback finally retires.

- **Seams.** Two, the fewest that exercise the real behavior:
  - The **adoption contract seam** — the two shell fragments driving the CLI
    subprocess against a sandbox tree (`gate-link-contracts.sh` for link + init,
    `gate-doctor-contracts.sh` for doctor, both the parity net). Prior art: the
    fragments themselves, unchanged in shape.
  - The **`go test` unit seam** — table tests beside `internal/adopt` at the
    injected filesystem/git boundary, for the parsers and computations the shell
    never reached. Prior art: `internal/gitguard` (the `Checker`-injected
    classification tables), `internal/maps`.

- **Gate command:** `bench gate` (the project gate, whose Go layer already runs
  `go build`/`go vet`/`go test ./...` and the cross-compile matrix, so the new
  package and tests are graded with no new wiring), plus the parent map's explicit
  per-slice done rule `go build ./... && go vet ./... && go test ./...`. Done per
  slice = gate green and those three green.

### Seam diagram — adoption contract seam (CLI subprocess → binary)

    trigger: contract fragment runs `bench <link|init|doctor>` against a sandbox tree
        │
        ▼
    mode + kit dir  ──▶ [ bench link → internal/adopt: plan build ]
      (+ BENCH_KIT)     [   preflight (7 conflict classes) — mutate NOTHING until all pass ]
                        [   AGENTS.md fence-aware block rewrite ▸ plan install ▸ manifest+#kit stamp ▸ pre-push ]
                                                                          ──▶ exit 0 + "linked …" / exit 1 naming the conflict, nothing written
    (no args)       ──▶ [ bench init → internal/adopt: [[ ! -e ]]-guarded scaffolds (gate/canary/learnings) ] ──▶ exit 0
    PATH sandbox    ──▶ [ bench doctor → internal/adopt: readback state (never exec shim) ; --fix atomic temp+rename ]
      env + BENCH_WRAPPER [   --fix writes a %q-quoted shim targeting the resolved wrapper (BENCH_WRAPPER) ]
                                                                          ──▶ exit 0 healthy/current / exit 1 stale|foreign|missing|refuse-foreign
    stamped manifest ─▶ [ bench doctor: #kit row vs `bench version` → skew | absent → skew-unknown warn ] ──▶ exit unchanged
    planted wrapper  ─▶ [ .bench/bin/bench.sh: no asset tree → refuse link/init/doctor, pointer ]          ──▶ exit ≠0, nothing written
        ◀ tests attach here: gate-link-contracts.sh (13 link cases + init stories) and gate-doctor-contracts.sh
          (11 sandboxed blocks) assert stdout/exit/filesystem before and after the port unchanged; the dropped
          planted-file assertion flips to its inverse, and the stamp / skew / planted-refuse rows are new-red.

### Seam diagram — `go test` unit seam (parsers + computations)

    trigger: gate Go layer runs `go test ./...`
        │
        ▼
    AGENTS.md bytes  ──▶ [ internal/adopt.<marker parse / fence-aware rewrite> ] ──▶ block region | malformed | unclosed-fence
    manifest bytes   ──▶ [ internal/adopt.<manifest parse>                     ] ──▶ hash | "" ; #kit row skipped ; last-wins
    (rel, dest)      ──▶ [ internal/adopt.<adapter symlink target>             ] ──▶ ../…-relative target
    PATH dirs + env  ──▶ [ internal/adopt.<doctor dir selection>              ] ──▶ first writable non-manager dir | ~/.local/bin
    shim bytes       ──▶ [ internal/adopt.<shim content / readback>           ] ──▶ target path (never executed)
        ◀ tests attach here: table tests pin marker edges (reversed markers, fenced example, unclosed fence, no
          trailing newline), manifest edges (#kit stamp row, absent file, duplicate rel → last-wins, symlink
          fingerprint), the adapter-target `../` computation, dir selection over manager-owned dirs, and the shim
          `%q` round-trip. Red before the package exists → does not compile.

### Acceptance coverage map

Per-item granularity is stated where a behavior quantifies over a set (each
conflict class, each report state, each parser edge).

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | plan/preflight/install parity — each of the 7 conflict classes refuses with nothing mutated; fresh link + relink yield exactly one managed block; manifest is `rel<TAB>sha256` single-writer | adoption contract | already covered (parity net): the 13 `gate-link-contracts.sh` cases (safe-fresh, existing-AGENTS, malformed-marker, conflict-no-manifest, modified-managed, metachar, worktree, hooksPath, default-branch, fenced, hooksPath-conflict, unclosed-fence) assert this today and stay green across the port | any divergence in preflight ordering or a conflict class writes a file the shell refused → one of the 13 unchanged cases goes red |
| 1 | marker/fence parse edges — reversed markers, fenced example preserved, unclosed fence refused, no-trailing-newline — each edge | go test unit | `internal/adopt` marker table before the parse fn exists → does not compile | a mis-ported fence walk rewrites text inside a documentation fence or mistrusts a valid block |
| 2 | scaffolds gate/canary/learnings; sentinel keeps the fresh gate red; second `init` clobbers nothing | adoption contract | already covered (parity net): the init scaffolding stories in `gate-link-contracts.sh` (fresh gate red-until-configured, sentinel-removed→green, second-init no-clobber) assert this today and stay green | a scaffold whose gate text drops the sentinel ships a green-by-default gate → the red-until-configured story goes red |
| 3 | four report states (healthy exit 0; stale/foreign/missing exit 1); report never executes shim contents; `--fix` is idempotent, atomic, refuses foreign byte-identically | adoption contract | already covered (parity net): the 11 `gate-doctor-contracts.sh` blocks (report, fix-write, spaced-target, idempotency, foreign-refuse, fallback, path-notice, stale-target, arg-passthrough, postinstall, session-start) assert this today and stay green | a port that exits 0 on stale/foreign, executes a hostile shim, or writes non-atomically trips one of the 11 unchanged blocks |
| 3 | doctor dir selection over manager-owned dirs; shim `%q` round-trip on a spaced/glob target | go test unit | `internal/adopt` dir-selection + shim table before the fns exist → does not compile | pins the manager-owned exclusion and quoting the black-box spaced-target case does not exhaust; a mis-ported `%q` corrupts a spaced target path |
| 4 | link writes a `#kit<TAB><version>` row; no file-row parses it; manifest keeps one writer | adoption contract + go test unit | a new link contract asserting the `#kit` row present after link and every managed rel still resolving its hash → red until link stamps; Go manifest table (`#kit` row skipped, duplicate rel last-wins) → no compile | catches a stamp that corrupts the ownership test (a `#kit` row read as a file, or a second writer) |
| 5 | stamped version ≠ `bench version` → doctor reports skew; absent stamp → skew-unknown warning, exit unchanged | adoption contract | a new doctor contract: a manifest stamped with a mismatched version → doctor output names skew; a pre-stamp manifest → warning, same exit → red until doctor reads the `#kit` row | catches a skew check that stays silent on real skew, or fails on a legitimate pre-stamp manifest (bricking not-yet-relinked repos) |
| 6 | planted `.bench/bin/bench.sh` refuses `link`/`init`/`doctor` with a pointer, mutating nothing; the kit wrapper routes through | adoption contract | a new contract: in a linked sandbox, `.bench/bin/bench.sh link` (and init/doctor) exits non-zero naming the real kit and writing nothing → red until the dispatcher asset-tree guard exists | catches a planted wrapper that half-runs `link` from a `.bench/` dir with no asset tree (mutating reviewer files against an empty source) |
| 7 | three `.sh` files deleted; safe-fresh contract's `[ -f .bench/bin/bench-link.sh ]` flipped to not-planted; no dangling reference in contracts, README, link/package | gate load + docs stale-reference sweep | the sweep against a tree still naming a deleted file → red; a dispatcher still sourcing a deleted `.sh` → gate load red; the un-flipped planted-file assertion → red once the file is gone | the conformance layer fails when a deleted file is still referenced, sourced, or asserted-present |

### Edge inventory

Walked per behavior against the profile's shell-CLI hostile-input checklist and
the map's item-6 owners; each resolved as a coverage row above or a **Won't
handle** line here.

- **paths/dirs with spaces, globs, or metachars** — covered: the link metachar
  kit-path contract and the doctor spaced-target `%q` contract (story 1/3); Go
  exec passes argv, never a shell string.
- **hand-edited file with no trailing newline** — covered: the doctor
  foreign-refuse contract recognizes a Bench-marked shim with no trailing newline
  (story 3), and the marker parser table walks the no-trailing-newline block
  (story 1).
- **absent file vs present-but-empty file** — covered: doctor refuses a
  present-but-empty file distinctly from missing (story 3); the manifest parse
  table distinguishes absent-file (no hash) from a present empty manifest (story 4).
- **unquoted multi-word / glob args** — covered: the doctor arg-passthrough
  contract passes `'a b' '*' c` through the shim verbatim (story 3).
- **required tool missing from PATH (no global bench, no `readlink -f`)** —
  covered: doctor resolves its own path without `readlink -f` (`resolve_script_path`
  stays), and a missing core binary hits `route_binary`'s exit-127 rim (story 3/6).
- **invocation through a symlink** — safe by construction: `resolve_script_path`
  already walks the symlink chain, and preflight's `has_symlink_parent` refuses a
  symlinked parent dir (covered by the safe-fresh contract, story 1).
- **SIGINT mid-mutation** — covered by construction: doctor's fix is temp+rename
  and the manifest is temp+`mv` (atomic single writer), so an interrupt leaves the
  old or the whole-new file, never a partial (stories 1, 3).
- **re-run idempotency** — covered: the relink contract (one marker), the
  second-init no-clobber story, and the doctor `--fix` idempotency block (stories
  1, 2, 3).
- **unclosed fence / malformed markers** — covered: the unclosed-fence and
  malformed-marker link contracts plus the marker parser table (story 1).
- **foreign pre-push / non-managed `.husky/pre-push`** — covered: the conflict and
  hooksPath-conflict contracts (story 1).
- **unset `origin/HEAD` (default branch)** — covered: the default-branch
  resolution contract (story 1).
- **a manifest row hand-crafted to look like a stamp, or a duplicate rel** —
  covered: the manifest parse table skips every `#`-prefixed line and applies
  last-assignment-wins on a duplicate rel (story 4).
- **cwd deeper than the repo root** — safe by construction: the binary resolves
  the root via `git rev-parse` and the gate runs from the root via `bin/bench.sh`,
  as today.
- **unwritable target dir (map item-6 owner)** — covered: doctor's dir selection
  picks the *first writable* non-manager-owned PATH dir and falls back to creating
  `~/.local/bin`; the fallback and path-notice doctor contracts (story 3) exercise
  a PATH with no writable dir. Walked here explicitly so the item-6 owner is a
  decision on the page, not inferred from "first writable."
- **Won't handle: the real global `npm i -g` kit-dir crossing** (hoisted platform
  package, no repo checkout) — gate-blind, a manual smoke per release (map Handoff
  #5/#7); the asset-tree sentinel makes a mis-resolved kit dir refuse cleanly
  rather than half-link, and the spec escalates per `craft-line` if resolution
  differs from the dev layout.
- **Won't handle: PATH-probe behavior beyond the sandboxed doctor cases** —
  gate-blind, as today (map Handoff #5); the 11 doctor blocks fix the sandboxed
  matrix.
- **Won't handle: the npm optional-deps lockfile edge** — open from the parent map
  (`decisions/go-rewrite.md`, asset #4), not this slice.
- **Won't handle: a kit dir whose asset tree is present-but-torn** (passes the
  sentinel yet is internally incomplete) — parity: the sentinel is a cheap
  planted-vs-real check, not a full integrity scan; a torn install is surfaced
  separately by the manifest and the new skew check, not by refusing the link.

## Out of scope

- **Slice 7 — worktree + shift loop ported to Go** (needs its flagged contract
  backfill first). A separate slice with its own spec and its own uncertainty flag
  in the parent map; `bench-worktree.sh` stays shell after this slice by design.
  Estimate to build later: one spec-sized session plus the backfill (~25–35 edits,
  ~5–10 gate runs).

- **Slice 8 — gate fragments → `go test`** (the gate ports last, under the canary
  layer's watch). A separate slice with its own spec. Estimate: one-plus
  spec-sized session (~30–40 edits, ~8–12 gate runs).

- **The session-start half of the version-skew surface** — a distinct capability
  that belongs to the already-landed slice-5 hooks spec's surface, not this one
  (map #2 splits the stamp's two halves deliberately). Not deferrable *within* this
  slice — it is a different hook's concern; recorded here only to name the seam.
  Estimate if ever revisited: ~4–6 edits, ~2 gate runs.
