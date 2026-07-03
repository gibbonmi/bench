# shim-autoinstall

Status: implemented

## Problem

A global `npm install -g benchkit` under a Node version manager (nvm, asdf, fnm,
volta) drops the `bench` executable only into the version manager's
version-specific bin. Only nvm-loaded **interactive** shells have that dir on
PATH — login shells and `bash -c "bench ..."` get `command not found`. Harness
hooks and any non-interactive invocation break. The user has no signal about why
and no one-line fix.

## Solution

Ship a plain-shell wrapper (`exec "<resolved-cli>" "$@"`) into a stable PATH dir
outside manager-owned territory, so every shell finds `bench`. One mutator writes
it: `bench doctor --fix`. npm postinstall calls that same code path under strict
guards; the SessionStart hook only detects and advises. `bench doctor` (report
only) tells the user when the shim is missing, stale, or blocked, and prints the
machine-exact fix and removal lines. A version-manager move strands the shim, but
the wrapper fails **loud** (exit 127, remedy message) rather than silently running
a stale copy, and one `doctor --fix` repairs it.

## User stories

1. As a user who installed `benchkit` globally under a version manager, I want a
   plain-shell `bench` shim on a stable PATH dir, so login shells and `bash -c`
   resolve `bench` the same as an interactive shell.
   Line: claude-sonnet-4-6 / medium. The outcome is the sum of the doctor
   behaviors below, so it inherits their routing rather than adding new logic.

2. As a user, I want `bench doctor` (no args) to report shim health — exit 0 when
   healthy, exit 1 when a fixable issue exists — and print the machine-exact fix
   and removal lines, mutating nothing.
   Line: claude-sonnet-4-6 / medium. Report/classify branching over platform and
   version-manager variance is real logic, but every state is a gate-observable
   black-box assertion.

3. As a user, I want `bench doctor --fix` to write the shim and announce every
   mutation (dir created, file written) with its absolute path, exiting 0 on
   success.
   Line: claude-sonnet-4-6 / medium. The write path (target selection, atomic
   write, announcement) is the deep unit and each step is assertable in the
   sandbox.

4. As a user, I want `bench doctor --fix` to be idempotent — a second run over an
   already-correct shim is an announced no-op that still exits 0.
   Line: claude-sonnet-4-6 / low. Re-run detection is a single marker-match branch
   at an already-built seam.

5. As a user, I want `bench doctor --fix` to refuse a foreign (non-Bench) file at
   the target path — report it, exit 1, leave it byte-identical — never clobber.
   Line: claude-sonnet-4-6 / medium. Getting the refuse-vs-write decision and the
   no-mutation guarantee right is where a wrong read does real damage.

6. As a user, I want the shim placed in the first writable PATH dir that is
   outside manager-owned territory (version-manager trees, Homebrew prefixes,
   system prefixes), creating `~/.local/bin` only as a fallback when no PATH dir
   qualifies.
   Line: claude-sonnet-4-6 / medium. The exclusion list plus the writability
   filter is the classification core the whole feature rests on.

7. As a user whose chosen dir is a fallback not on PATH, I want `bench doctor
   --fix` to print the exact one-line PATH addition for my shell and stop — rc
   files are never edited.
   Line: claude-sonnet-4-6 / low. One conditional branch emitting a fixed string
   for the resolved shell.

8. As a user whose node version moved (a version-manager switch relocated the
   package), I want the stale shim to fail loud — exit 127 with a remedy message
   on stderr — not silently run a stale copy.
   Line: claude-sonnet-4-6 / low. The wrapper is a fixed generated template; the
   assertion drives it with a missing target.

9. As a user, I want the shim to pass multi-word and glob-containing arguments
   through to the real CLI intact.
   Line: claude-sonnet-4-6 / low. `exec "<target>" "$@"` with correct quoting;
   assertion drives it with hostile args.

10. As a user running a global install through npm, I want postinstall to
    auto-run the fix under strict guards — `npm_config_global` truthy **and** the
    package root has no `.git` — relaying the fix's announcements.
    Line: claude-sonnet-4-6 / medium. The guard conjunction is the safety
    boundary; each permutation must be pinned.

11. As a developer working from a git checkout or an `npm link`, I want
    postinstall to skip the mutation (the `.git` guard covers both).
    Line: claude-sonnet-4-6 / low. One guard branch, asserted with a planted
    `.git`.

12. As a user installing under pnpm, yarn, or bun (no `npm_config_global`), I want
    postinstall to fall through to a single advice line and exit 0.
    Line: claude-sonnet-4-6 / low. Env-absent branch to the advice path.

13. As any user, I want postinstall to never fail a global install — it exits 0 on
    every path, including probe failure and write failure.
    Line: claude-sonnet-4-6 / medium. The always-exit-0 invariant must hold across
    every guard and failure permutation; it is the one thing that must not regress.

14. As a user opening a cold session where `bench` resolves only by path, or where
    the resolved CLI fails, I want the SessionStart advisory to append "run `bench
    doctor --fix`" to its existing line.
    Line: claude-sonnet-4-6 / low. A one-line append to the existing advisory
    branch; thin pass-through.

15. As a user uninstalling, I want the README uninstall note to carry the generic
    removal pair and `bench doctor`'s report to print the machine-exact removal
    pair, so the orphaned shim is easy to delete.
    Line: claude-sonnet-4-6 / low. Doc line plus one report line built from the
    already-resolved shim path.

## Implementation decisions

- **New subcommand `bench doctor [--fix]`** lives in `bin/bench-doctor.sh`,
  sourced by `bin/bench.sh` alongside the other `bench-*.sh` includes and
  dispatched from the `case` in `bin/bench.sh` (`doctor) shift; doctor "$@" ;;`).
  It is the **deep unit**: it hides platform and version-manager variance (probe
  order, the exclusion list, macOS BSD-userland vs Linux, shell detection) behind
  a report/fix interface. Postinstall, the wrapper, and the session-start
  extension are thin pass-throughs with no seam of their own.

- **Doctor must run with no `bench` on PATH and no GNU `readlink -f`** — that is
  the feature's own premise. It resolves the running CLI's real path with the
  existing `resolve_script_path` helper (portable symlink walk already in
  `bin/bench.sh`), never `readlink -f`.

- **Target-dir selection.** First writable dir already on PATH that is *not*
  manager-owned. Manager-owned (writable but wrong): version-manager trees
  (`$NVM_DIR`/`~/.nvm`, `$ASDF_DATA_DIR`/`~/.asdf`, fnm dirs, `$VOLTA_HOME`/
  `~/.volta`), Homebrew prefixes (`$HOMEBREW_PREFIX`, `/opt/homebrew`,
  `/usr/local`), and system prefixes (`/usr`, `/bin`, `/sbin`, `/opt`). If none
  qualifies, create `~/.local/bin` as the fallback. Every mutation (dir created,
  file written) is announced with its absolute path.

- **The wrapper is a generated artifact**, never hand-edited: a marker comment
  identifying it as Bench-written, then `exec "<target>" "$@"`, where `<target>`
  is the CLI's own resolved real path captured at fix time (a static wrapper, not
  a symlink, not a runtime prober). Missing target → `bench moved — run
  \`bench doctor --fix\` or reinstall` on stderr, exit 127.

- **Atomic write.** `--fix` writes to a temp file in the target dir then `mv`s it
  into place, so a SIGINT mid-fix leaves either the old state or the complete new
  shim — never a partial file.

- **Collision policy.** `--fix` writes only when the target path is absent or the
  existing file carries the Bench marker. A foreign file (no marker) is reported
  and refused with exit 1 and left byte-identical. Present-but-empty is classified
  as foreign (no marker) and refused; absent is written.

- **PATH-off case.** When the chosen dir is not on PATH (the fallback), `--fix`
  prints the exact one-line PATH addition for the user's shell (`$SHELL`
  basename → the right rc syntax) and stops; it never edits rc files. Debian-family
  self-heals `~/.local/bin` at next login; elsewhere the user pastes.

- **Postinstall** is a new thin file (e.g. `bin/bench-postinstall.sh`) wired into
  `package.json` `scripts.postinstall` and added to `files[]` so it ships. It
  invokes `doctor --fix` **package-relative** (it cannot assume `bench` is on
  PATH) only when `npm_config_global` is truthy AND the package root contains no
  `.git`. On any other condition — env var absent, `.git` present, probe failure,
  write failure — it prints one advice line ("run `bench doctor --fix`") and
  exits 0. **It never exits nonzero**; the shim is a convenience and must not fail
  a global install. `npm_config_global` is set to `"true"` for `npm install -g`
  across npm v6–v10+ and is absent under pnpm/yarn/bun; the fail-to-advice posture
  makes a wrong read non-destructive either way.

- **SessionStart extension only.** `.bench/hooks/session-start.sh` already prints
  an advisory when `bench` resolves by path or the resolved CLI fails. This spec
  appends "run `bench doctor --fix`" to that existing branch — the hook never
  mutates (a silent hook write to `~/.local/bin` would violate the
  announce-mutations posture).

- **README uninstall note** carries the generic pair
  (`npm uninstall -g benchkit && rm -f "$(command -v bench)"` — after npm's
  symlink is gone, `command -v bench` resolves to the orphaned shim). `bench
  doctor`'s report prints the machine-exact removal pair since it knows the real
  shim path. No lifecycle cleanup hook: `--ignore-scripts` installs would miss it,
  and a delete hook writing outside npm's tree is worse than the inert,
  self-describing orphan.

- **Doctor scope is the shim only.** The parked wiring-drift doctor checks
  (roadmap) keep their evidence bar and stay parked.

## Testing decisions

- **What a good test is here:** drive `bench doctor`, the generated wrapper, and
  the postinstall script as **subprocesses** under a fabricated `HOME`/`PATH`/env
  sandbox and observe exit code, stdout/stderr, and filesystem effects. Never
  assert against internal functions — the CLI-subprocess boundary is the external
  behavior.

- **Which seam:** one seam family, the **CLI-subprocess sandbox** — exactly the
  runtime-contracts pattern already used by `.bench/gate-runtime-contracts.sh`
  (fresh `mktemp -d`, fabricated env, invoke `bash "$root/bin/..."`, assert). This
  is the highest seam that exercises the real behavior; all four units (doctor
  report, doctor `--fix` + generated wrapper, postinstall, session-start
  advisory) are driven through it. Prior art: the `bench idea/roadmap` and gate
  cwd contracts in `gate-runtime-contracts.sh`, and the `--space-path`/`--no-repo`
  fixture options in `gate-contract-runner.sh`.

- **New gate fragment:** `.bench/gate-doctor-contracts.sh`, sourced by
  `.bench/gate.sh` in the runtime-contracts group, built on the shared `contract`
  harness. One **red-by-construction canary** per new check under `tests/canary/`.

- **Gate command:** `.bench/gate.sh` (the project gate).

- **Gate-blind spot (declared):** the contracts fake npm's env; a genuine
  `npm i -g` lifecycle run is not exercised by the gate. One **manual smoke on a
  real global install** before release. This is the only acceptance not
  gate-observable.

### Seam diagram

    trigger: `npm i -g` postinstall  ·  `bench doctor [--fix]`  ·  the generated shim  ·  SessionStart hook
        │
        ▼
    fabricated HOME/PATH/env sandbox
    (mktemp -d; PATH = planted dirs; SHELL set; readlink faked)
        │
        ▼
    args + env  ──▶  [ bin/bench-doctor.sh      ]  ──▶  exit code (0/1/127)
    planted PATH ──▶  [ bin/bench-postinstall.sh ]  ──▶  stdout announcements / advice
    planted files──▶  [ generated wrapper        ]  ──▶  stderr remedy
    (marker/foreign) [ session-start.sh advisory ]  ──▶  filesystem (shim written / refused / untouched)
                          ◀ tests attach here: invoke as a subprocess in the
                            sandbox; assert exit code + stdout/stderr substrings +
                            shim bytes + that refused/skipped paths are unchanged

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 2 | `doctor` exits 0 when the shim is healthy, 1 when missing/stale/foreign; mutates nothing | doctor subprocess | `bash bin/bench-doctor.sh` in a sandbox with no shim → expect exit 1, no file written (fails: no such file) | exit code and no-mutation are the report contract; a report that writes or returns the wrong code fails the assertion |
| 2, 15 | `doctor` prints the machine-exact fix line and removal pair | doctor subprocess | grep report output for the resolved shim path in the removal pair (fails: absent) | a drifted or generic path in the removal line fails the substring match |
| 3, 6 | `--fix` writes a wrapper carrying the marker + resolved target into the first non-manager writable PATH dir, announcing the path | doctor `--fix` subprocess | plant a writable nvm dir + a plain dir on PATH; run `--fix`; assert shim lands in the plain dir with marker+target and the path is announced (fails: no file) | picking the manager dir, omitting the marker, or silent write each fails a distinct assertion |
| 4 | second `--fix` over a correct shim is an announced no-op, exit 0 | doctor `--fix` subprocess | run `--fix` twice; assert shim bytes unchanged and second run announces no-op (fails: file absent first run) | a re-write or a nonzero second exit fails the idempotency assertion |
| 5 | `--fix` refuses a foreign file at the target: exit 1, byte-identical | doctor `--fix` subprocess | plant a marker-less file at the target; run `--fix`; assert exit 1 and byte-identical (fails: file absent) | a clobber changes the bytes; a wrong exit code fails the refuse assertion |
| 6 (fallback) | with only manager dirs on PATH, `--fix` creates `~/.local/bin`, announces the dir creation, writes there | doctor `--fix` subprocess | PATH = only a planted nvm dir; run `--fix`; assert `$HOME/.local/bin/bench` exists and dir-creation announced (fails: absent) | failing to fall back leaves no shim; a silent mkdir fails the announcement assertion |
| 7 | when the chosen dir is off PATH, `--fix` prints the exact PATH-addition line for the shell and edits no rc file | doctor `--fix` subprocess | fallback case with `SHELL=/bin/bash`; assert output contains the `~/.local/bin` PATH export and no rc file changed (fails: absent line) | a missing paste-line or an edited rc file fails the assertion |
| 8 | wrapper with a missing target exits 127 with the remedy substring on stderr | wrapper subprocess | generate a shim, delete the target, run the shim; expect exit 127 + `bench moved` on stderr (fails: shim absent) | a silent success or a different code fails the loud-failure contract |
| 9 | wrapper passes multi-word and glob args through intact | wrapper subprocess | run the shim with `'a b' '*'` against a stub target that echoes `"$@"`; assert args intact (fails: shim absent) | dropped or resplit args fail the passthrough assertion |
| 10 | postinstall invokes `--fix` only when `npm_config_global` truthy AND no `.git` | postinstall subprocess | `npm_config_global=true`, no `.git`; run postinstall; assert shim written + announcements relayed (fails: script absent) | a missing invocation leaves no shim; the guard conjunction is pinned per permutation |
| 11 | postinstall skips mutation when `.git` is present | postinstall subprocess | plant `.git`, `npm_config_global=true`; run; assert no shim written, advice printed, exit 0 (fails: script absent) | a mutation under the dev-checkout guard fails the no-write assertion |
| 12 | postinstall with `npm_config_global` unset prints one advice line, exit 0 | postinstall subprocess | unset the var; run; assert advice line + exit 0, no shim (fails: script absent) | a mutation or nonzero exit under a non-npm installer fails the fall-through assertion |
| 13 | postinstall exits 0 on every guard and failure permutation | postinstall subprocess | force a write failure (read-only target) with guards satisfied; assert exit 0 + advice (fails: script absent) | any nonzero exit fails the never-fail-the-install invariant |
| 14 | SessionStart appends "run `bench doctor --fix`" when `bench` resolves by path or the CLI fails | session-start subprocess | drive session-start.sh with `bench` resolvable only by path; assert output contains the doctor advice (fails: substring absent) | a missing append leaves the cold session with no repair pointer |

Enumeration for the quantified rows: story 13's "every permutation" = the cross of
{`npm_config_global` truthy/unset} × {`.git` present/absent} × {write succeeds /
write refused} — the postinstall contract asserts exit 0 for each. Story 2's
states = {healthy, missing, stale, foreign}, one exit-code assertion each. Each
listed check gets its own red-by-construction canary fixture; the canary set below
enumerates them.

### Edge inventory

Walked per behavior against the profile's shell-CLI hostile-input checklist:

- **paths/dirs with spaces or globs** → covered: doctor probe + wrapper quoting,
  exercised via the `--space-path` fixture option (story 3/9 assertions run under a
  spaced parent).
- **hand-edited file, no trailing newline** → covered: marker detection must match
  a marker line with no trailing newline (story 5 foreign-vs-marker classification
  asserts a newline-less Bench-marked file is still recognized, not refused).
- **absent vs present-but-empty target** → covered: absent → write (story 3),
  present-but-empty → no marker → foreign, refuse (story 5).
- **unquoted multi-word args** → covered: story 9 (`"$@"` passthrough).
- **required tool missing (no global `bench`, no `readlink -f`)** → covered:
  doctor and postinstall run with PATH stripped of `bench` and resolve via the
  portable symlink walk, not `readlink -f` (story 2/10 sandbox strips PATH).
- **invocation through a symlink** → covered: target resolves to the CLI's own
  real path; the existing runtime contract already fakes `readlink` (reuse it).
- **SIGINT mid-fix** → covered by construction: atomic temp-file + `mv`; the
  assertion checks no `.tmp` residue remains after `--fix` and the shim appears
  whole (a partial write is the failure this prevents).
- **re-run idempotency** → covered: story 4 (announced no-op).
- **cwd deeper than repo root / outside any repo** → covered: doctor is
  HOME/PATH-scoped; a coverage assertion runs `doctor` with cwd in a non-repo
  `mktemp -d` (the `--no-repo` fixture) and expects a normal report.
- **hostile environment** (manager-owned dir writable but wrong) → covered:
  story 6 excludes it via the manager-owned list even when writable.

**Won't handle:**
- **A path containing a literal newline** — the one shape the porcelain/`$@`
  handling across the kit already misreads (documented in `bin/bench.sh`); out of
  band for a PATH dir and consistent with the rest of the CLI.
- **rc-file editing to put the fallback dir on PATH** — a closed decision (#6 in
  the map): the paste-line is printed, rc files are never touched.
- **Uninstall lifecycle cleanup of the orphaned shim** — a closed decision (#5):
  the orphan is inert and self-describing (127 message); removal is handed to the
  user twice (README + doctor report).

### Canary fixtures (one red-by-construction per new check)

Each plants a broken kit tree so a rotted-to-always-pass check goes red with its
targeted substring, mirroring `tests/canary/`'s existing fixtures:

- `doctor-foreign-clobbered` — a `--fix` that overwrites a foreign file instead of
  refusing (story 5 check).
- `doctor-stale-silent` — a wrapper that exits 0 on a missing target instead of
  127 (story 8 check).
- `wrapper-args-dropped` — a wrapper using `$*`/`$1` instead of `"$@"` (story 9).
- `postinstall-nonzero-exit` — a postinstall that exits nonzero on write failure
  (story 13 invariant).
- `postinstall-guard-bypassed` — a postinstall that mutates with `.git` present
  (story 10/11 guard).
- `doctor-manager-dir-chosen` — a `--fix` that writes into a manager-owned dir
  (story 6 exclusion).
- `session-start-advice-dropped` — the advisory branch missing the doctor pointer
  (story 14).

## Out of scope

- **Wiring-drift doctor checks** (detecting a linked repo whose `.bench/` has
  drifted from the kit) — a separate capability with its own evidence bar, already
  parked on the roadmap; `bench doctor` here ships shim scope only. Estimate to
  build later: ~6 edits, ~4 gate runs.
- **Windows/PowerShell shim** — a distinct platform capability; the package is
  `os: [darwin, linux]` today. Estimate: ~5 edits, ~3 gate runs.
