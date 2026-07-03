# Shim auto-install — node-independent `bench` on PATH in every shell

Graduated from the roadmap 2026-07-03 (grilled on the top tier per the parked
note). Root cause on record: a global npm install puts `bench` only in the
version manager's version-specific bin, which only nvm-loaded interactive
shells see — login shells and `bash -c` get "command not found". The fix is a
plain-shell wrapper in a stable PATH dir. All tickets were resolved in the
bootstrap grill.

## #1: Who is allowed to write the shim?

Type: Grill

### Question
Where does the mutation live — npm postinstall, a doctor command, a session
hook, or some mix?

### Answer
One mutator: `bench doctor --fix`. npm postinstall *calls it* under strict
guards (#3) and relays its announcements; the SessionStart hook never mutates —
it only detects a missing/broken shim and prints "run `bench doctor --fix`".
Writing outside npm's tree is the invasive step, so it goes through the same
audited code path as the explicit command, and a session hook that silently
writes to `~/.local/bin` would violate the announce-mutations posture.
`bench doctor` ships with shim scope only; the parked wiring-drift doctor
checks (roadmap) keep their evidence bar and stay parked.

## #2: Where does the shim go?

Type: Grill

### Question
Which directory receives the wrapper, and how are wrong-but-writable dirs
excluded?

### Answer
First writable dir already on PATH that is outside manager-owned territory;
create `~/.local/bin` only as a fallback when no dir qualifies. Manager-owned
territory (writable but wrong): version-manager trees (`$NVM_DIR`/`~/.nvm`,
`$ASDF_DATA_DIR`/`~/.asdf`, fnm's dirs, `$VOLTA_HOME`/`~/.volta`), Homebrew
prefixes (`$HOMEBREW_PREFIX`, `/opt/homebrew`, `/usr/local`), and system
prefixes (`/usr`, `/bin`, `/sbin`, `/opt`) — writability filters most of these
naturally; the list catches the writable exceptions. Every mutation (dir
created, file written) is announced with its path.

## #3: When may postinstall invoke the fix, and what is its failure posture?

Type: Grill

### Question
Which guards gate the postinstall call, and what happens on ambiguity or
failure?

### Answer
Invoke `doctor --fix` only when `npm_config_global` is truthy **and** the
package root contains no `.git` (one guard skips both the dev checkout and
`npm link`ed checkouts). On any other condition — env var absent (pnpm, yarn,
bun), probe failure, write failure — print one advice line ("run
`bench doctor --fix`") and exit 0. Postinstall never exits nonzero: the shim is
a convenience and must not fail a global install. Non-npm installers fall
through to the SessionStart advice path.

## #4: What is the shim, and how does it survive version-manager moves?

Type: Grill

### Question
Static pointer or runtime prober, given an nvm node-version switch relocates
the global package and strands the shim?

### Answer
A static wrapper script (not a symlink, not a prober): marker comment
identifying it as Bench-written, then `exec "<target>" "$@"` where the target
is the running CLI's own resolved real path at fix time. If the target is gone
it prints "bench moved — run `bench doctor --fix` or reinstall" to stderr and
exits 127. Staleness is caught by `doctor` and the SessionStart advisory, and
the fix is one command. A runtime prober globbing nvm/asdf/fnm/volta trees can
silently pick a stale copy after an upgrade — a loud self-describing break
beats a quiet wrong answer; a wrapper beats a symlink because a dangling
symlink can't explain itself.

## #5: What happens on uninstall?

Type: Grill

### Question
`npm uninstall -g` won't remove files it didn't install, so the wrapper is
orphaned. Preuninstall cleanup?

### Answer
No lifecycle cleanup — `--ignore-scripts` installs would miss it, and a
delete hook writing outside npm's tree is worse than the disease. The orphan is
inert and self-describing (#4's 127 message). The removal line is handed to the
user twice: the README uninstall note carries the generic pair
(`npm uninstall -g benchkit && rm -f "$(command -v bench)"` — after npm's
symlink is gone, `command -v bench` resolves to exactly the orphaned shim), and
`bench doctor`'s report prints the machine-exact removal pair, since it knows
the real shim path.

## #6: What is `bench doctor`'s surface contract?

Type: Grill

### Question
Exit codes, collision policy, and SessionStart wiring — resolved inline as
trivially decidable.

### Answer
`bench doctor` is report-only: exit 0 healthy, 1 when a fixable issue is found;
prints the machine-exact fix and removal lines. `--fix` announces every
mutation and exits 0 on success; a second run is an announced no-op. Collision:
`--fix` writes only when the path is absent or the existing file carries the
Bench marker; a foreign file is reported and refused (exit 1) — never
clobbered. When the chosen dir is not on PATH (fallback case), `--fix` prints
the exact one-line PATH addition for the user's shell and stops — rc files are
never edited (self-heals at next login on Debian-family; user-owned paste
elsewhere). SessionStart extends its existing advisory line only: when `bench`
resolves by path or the resolved CLI fails, append the doctor advice.

## Handoff

1. **Module boundaries.** `bench doctor [--fix]` — new subcommand (own
   `bin/bench-doctor.sh` per the existing dispatch pattern): probe, classify,
   report, write. The postinstall script — new thin file wired into
   `package.json` `scripts` + `files[]`: guards, package-relative call, advice
   line. The wrapper — generated artifact, never hand-edited. `session-start.sh`
   — advisory extension only. README — uninstall note.
2. **Contracts.** `doctor`: no args, report to stdout, exit 0 healthy / 1
   fixable, mutates nothing. `doctor --fix`: writes shim (creates fallback dir
   if needed), announces each mutation with its path, refuses foreign files
   (exit 1), prints PATH paste-line when the dir is off PATH, exit 0 on
   success, idempotent. Postinstall: env-guarded, always exit 0, one advice
   line on skip/fail. Wrapper: `exec` target with `"$@"`; missing target →
   stderr remedy + exit 127; carries the Bench marker + target path.
3. **Deep vs thin.** `doctor` is the deep unit — it hides platform and
   version-manager variance (probe order, exclusion list, macOS/Linux
   differences) behind the report/fix interface. Postinstall, wrapper, and the
   session-start extension are thin pass-throughs with no seam of their own.
4. **Black-box assertables.** Under a fabricated `HOME`/`PATH` sandbox:
   `doctor` exit codes for healthy / missing / stale / foreign states; `--fix`
   writes a wrapper containing marker + expected target and announces the path;
   `--fix` on a foreign file exits 1 and leaves it byte-identical; second
   `--fix` is a no-op; wrapper with missing target exits 127 with the remedy
   substring; wrapper passes multi-word args through intact; postinstall exits
   0 under every guard permutation (var unset, `.git` present, write refused).
5. **Gate attachment.** The runtime-contracts family (tmp-sandbox pattern
   already used by `gate-runtime-contracts.sh`), one red-by-construction canary
   per new check. Gate-blind spot: a genuine `npm i -g` lifecycle run (the
   contracts fake npm's env) — one manual smoke on a real install before
   release.
6. **Hostile-input owners.** Paths with spaces/globs → doctor's probe and the
   wrapper's quoting. No-trailing-newline hand-edited file → doctor's marker
   detection. Absent vs present-but-empty shim path → doctor classification
   (absent → write; empty → foreign, refuse). Unquoted multi-word args →
   wrapper `"$@"`. Required tool missing → the feature's own premise: doctor
   and postinstall must run with no `bench` on PATH and no `readlink -f`
   (package-relative / self-resolved invocation). Symlink invocation → target
   resolution of the CLI's own real path (existing runtime contract already
   fakes `readlink`). SIGINT mid-fix → atomic write (temp file + `mv`). Re-run
   idempotency → announced no-op. cwd deeper than repo root → doctor is
   HOME/PATH-scoped and must work outside any repo.
7. **Uncertainty flags.** `npm_config_global` behavior across npm majors (and
   its absence under pnpm/yarn/bun) — spec-writer verifies against current npm
   docs; the fail-to-advice posture makes a wrong read non-destructive either
   way.
8. **Rejected alternatives.** Editing shell rc files (un-uninstallable,
   fights dotfile managers); preuninstall cleanup (unreliable under
   `--ignore-scripts`, worse than the orphan); symlink shim (dangling links
   can't explain themselves); runtime-probing wrapper (can silently resolve a
   stale copy); SessionStart auto-fix (silent mutation from a hook); shim
   logic duplicated in postinstall (second mutator, drifts from doctor).
9. **Domain watch-outs.** An nvm node-version switch strands the shim until
   the next doctor/SessionStart look — the wrapper's own 127 message is the
   floor. macOS ships BSD userland (no `readlink -f`; `~/.local/bin` not on
   PATH by default). Debian-family adds `~/.local/bin` to PATH only at next
   login after creation. Homebrew prefixes are user-writable on macOS but are
   another manager's territory. The orphaned wrapper after uninstall is a
   documented, inert leftover.

Dependency order: n/a — single spec; internal build order doctor `--fix` core →
postinstall caller → session-start advisory + README note.
