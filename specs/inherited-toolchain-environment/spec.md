# inherited-toolchain-environment

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-22 — execute `$bench-write-spec` for FT242 in an isolated worktree, review with Terra/medium, and commit the accepted spec to `main`; behavior and scope are the drained `roadmap/FT242.md` decision

Verification log: 3 iteration(s) to accept — Terra/medium found the missing real-client evidence and unbounded login-shell discovery; the author folded a two-second process-group bound and manual Codex/CLI evidence gate, then added a descendant sentinel after correction review showed the first timeout row could not detect an orphaned child

## Problem

A harness process can inherit an initialization marker without the PATH effects
that marker promises. In the reproduced Codex client, `ENVMAN_LOAD=loaded` was
present and Go was absent from PATH even though the same user's clean login shell
resolved the installed Go executable. The kit then has two bad outcomes: a
Git-source Bench run cannot build its private exact-source executable, and an
installed Bench grading a Go repository can silently omit the built-in Go phases
because the phase table treats a missing toolchain as if the repository did not
ask for those phases.

## Solution

Make the two existing owners tell the truth independently. The built-in gate
phase table refuses a Go module when Go is absent from PATH, before any phase
runs, while repositories with no Go module and repositories declaring their own
phase manifest keep their current behavior. The informational SessionStart hook
recognizes the reproduced partial environment, asks the user's clean login shell
where Go resolves, and prints a shell-quoted recovery that prepends the discovered
directory to the existing PATH. It never executes the discovered Go binary or
blocks session startup. If no clean-login Go exists, it names the missing
requirement without inventing a user-specific path.

Canonical vocabulary: **environment closure** means the harness process carries
the toolchain effects implied by its inherited initialization state. Avoid
"loaded environment" and "working PATH": both hide the partial-propagation case.

## User stories

### Group A — a required Go toolchain cannot disappear from the gate

Line: `gpt-5.6-terra` / medium. The gate and manifest seams are known, but a
wrong fail posture would weaken the oracle.

1. As a maintainer of a Go repository using the built-in phase table, I want a
   missing Go executable to make phase-table construction red, so that no green
   verdict credits omitted Go checks.
2. As an operator, I want the refusal to name Go, PATH, and the graded repository,
   so that the missing prerequisite is distinguishable from a test failure.
3. As a maintainer of a non-Go repository, I want an absent `go.mod` to remain a
   valid no-Go shape, so that the kit does not impose Go on unrelated projects.
4. As a maintainer of a repository with its own phase manifest, I want that
   manifest to remain authoritative when Go is absent, so that the built-in
   table does not second-guess project-declared phases.

### Group B — SessionStart diagnoses partial propagation safely

Line: `gpt-5.6-sol` / high. The hook is cross-harness guidance and launches a
login shell, so leverage and bootstrap-authority reasoning override the known
shell seam.

5. As an agent in a Go repository, I want SessionStart to identify when Go is
   absent from the harness PATH but available to my clean login shell, so that I
   do not misdiagnose an installed toolchain as missing.
6. As an agent repairing that environment, I want a copy-paste command that
   prepends the discovered directory to the existing PATH, so that recovery
   cannot discard Codex-provided tools such as `rg`.
7. As a security reviewer, I want clean-login output treated only as diagnostic
   data and never executed by the hook, so that discovery does not grant
   executable authority.
8. As a session user, I want the diagnostic to stay informational and exit zero,
   so that a missing toolchain cannot prevent the session from opening.
9. As an agent whose harness PATH already resolves Go, I want no recovery warning,
   so that healthy sessions stay quiet.
10. As an agent whose machine has no clean-login Go, I want an honest missing-Go
    diagnosis without a fabricated recovery path, so that the hint cannot point
    at arbitrary home-directory bytes.
11. As a user opening a session outside a repository, I want the hook to remain
    completely silent, so that one repository never leaks ambient advice into
    another context.
12. As a maintainer, I want the same repository and toolchain fixture exercised
    in the partial harness shape and the healthy CLI shape, so that automated
    regression evidence compares only the environment effects that differ;
    actual Codex-client/CLI observations remain required for any portability claim.
13. As a session user, I want timed-out login initialization and every child it
    launched terminated, so that diagnosis cannot leave background work running
    after SessionStart continues.

## Implementation decisions

- **The phase table owns required-tool refusal.** `phaseTable` remains the one
  loader for manifest-backed and built-in schedules. Its built-in branch returns
  an error when a regular `go.mod` is present and `go` cannot be resolved from
  PATH; it does not materialize an environment skip or an empty Go phase set.
  Manifest-backed tables bypass this built-in requirement and let the ordinary
  runner report any command they explicitly declare.
- **Session inspection owns the recovery hint.** The existing SessionStart hook
  remains the cross-harness trigger and `internal/sessioninspect` remains its
  bounded Go owner. A first inspection phase checks only an in-repository regular
  `go.mod`, only when `exec.LookPath("go")` fails in the harness environment. It
  invokes a clean login shell with `ENVMAN_LOAD` removed, validates one line naming
  an absolute existing executable, and uses `sanitize.ShellQuote` to render the
  executable's directory in a command that prepends to literal `"$PATH"`.
  Invalid, multiline, or control-bearing discovery output produces an honest
  missing-Go diagnosis and no path-bearing recovery command.
- **Discovery is bounded by one registry fact.** `bounds.EnvironmentDiscoveryTimeout`
  is two seconds, inside the existing ten-second SessionStart aggregate. Discovery
  runs through `bounds.Run`, which owns process-group timeout and teardown. On
  timeout or nonzero exit, inspection prints the missing-Go diagnosis with no
  path-bearing recovery, continues the remaining SessionStart phases, and exits
  zero. System tests derive their outer deadline from the named bound.
- **Bootstrap authority before execution.** Trust chain: the harness invokes the
  repository-tracked SessionStart hook; the hook invokes the user's configured
  login shell, accepted here as user-owned configuration; the hook validates but
  does not execute the path that shell reports; the operator explicitly runs the
  printed recovery; the gate then independently resolves Go from the repaired
  PATH and executes it as the declared project toolchain. The login shell is a
  reviewer-visible trust assumption, not an authentication root for the reported
  executable.
- **One term, one profile rule.** `CONTEXT.md` gains `environment closure` and
  its Avoid list. The project profile's cold-session notes state the recovery
  rule once: prepend a recovered tool directory to the ambient PATH; never
  replace the rest of the harness toolchain. The existing hostile-input checklist
  already owns `required tool missing from PATH`; it is not duplicated.
- **No automatic home-directory search.** The hook asks the user's own clean
  login resolution and never probes conventional install locations. A failed or
  unsafe answer stays a diagnosis, not a guessed command.
- **No new command or module.** The feature deepens the existing phase-table and
  SessionStart seams. The hook remains informational and the gate remains the
  oracle.

## Testing decisions

- The gate unit seam constructs real temporary repositories and controls PATH.
  Tests call the existing phase-table loader, asserting its returned schedule or
  error; no internal collaborator is mocked.
- Unit tests at the session-inspection seam pin validation, quoting, and the
  two-second process-group timeout. The tagged system seam launches the real
  `.bench/hooks/session-start.sh` through the existing process ledger. A private
  HOME and shell initialization reproduce the inherited-marker/clean-login-PATH
  split, while a sibling launch uses the same repository and fake Go executable
  in the healthy CLI shape.
- Fixture parity is not portability evidence. Before phase close, the diagnostic
  ticket records one actual Codex-client invocation and one actual CLI invocation
  against the same WSL repository, Go executable, and initialization files. No
  portability claim is made unless both observations are present.
- The project gate observes the unit rows in `test` and the real-hook rows in
  `system`; shellcheck observes the edited hook when that capability is present.

### Seam diagram

    trigger: `bench gate` resolves a repository's phase table
        │
        ▼
    go.mod + PATH ──▶ [ phaseTable: manifest or built-in schedule ] ──▶ phases or refusal
                       ◀ unit tests attach here: real temp repository + controlled PATH

    trigger: harness invokes repository-tracked SessionStart hook
        │
        ▼
    repo + harness PATH ──▶ [ session-start.sh → bounded session-inspect ] ──▶ dashboard and optional recovery hint
                                              │
                                              └── bounded clean-login discovery only when required Go is absent
                       ◀ unit tests attach at session-inspect; system tests drive the real hook with private HOME/PATH

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| TE1 | 1 | a built-in phase-table load for a repository with a regular `go.mod` and PATH containing no `go` returns red before producing a schedule | phase table | the current loader returns a schedule with every Go phase silently absent, so the exact false green makes this row red |
| TE2 | 2 | the missing-tool refusal names `go`, `PATH`, and the graded repository path | phase table | a generic table error would leave the reproduced infrastructure failure indistinguishable from an invalid manifest or test red |
| TE3 | 3 | a repository with no `go.mod` and the same Go-less PATH returns its unchanged non-Go built-in phases without error | phase table | the cheapest fail-closed patch rejects every repository when Go is missing |
| TE4 | 4 | a repository with a valid phase manifest and the same Go-less PATH returns exactly its declared phases | phase table | a global preflight for Go would override the manifest even when the project declared no Go command |
| TE5 | 5 | the real SessionStart hook in a Go repository with `ENVMAN_LOAD=loaded`, no harness-PATH Go, and a clean-login Go prints a partial-environment diagnosis containing the discovered executable | system hook | today the hook either prints its ordinary dashboard/no-core hint or fails later, never identifying the marker-without-effects shape |
| TE6 | 6 | that diagnostic prints a copy-paste assignment that prepends the shell-quoted discovered directory to literal `"$PATH"` | system hook | a replacement assignment can fix Go while deleting the harness-provided `rg`, the second observed failure |
| TE7 | 7 | a marker executable placed at the discovered Go path is not invoked during SessionStart | system hook | the cheapest unsafe implementation runs `go version` or the discovered binary to validate it |
| TE8 | 8 | every partial or missing-toolchain SessionStart case exits zero | system hook | a diagnostic that propagates the login shell or resolver failure would block session startup |
| TE9 | 9 | with the same repository and Go executable already on harness PATH, SessionStart prints no partial-environment or recovery line | system hook | an unconditional login-shell probe or warning makes healthy CLI sessions noisy and changes the normal path |
| TE10 | 10 | when neither harness PATH nor the clean login resolves Go, SessionStart names missing Go but prints no path-bearing PATH assignment | system hook | guessing a conventional home-directory location would still print a plausible but unauthoritative command |
| TE11 | 11 | outside a repository, the real hook remains exit zero with empty stdout and stderr | system hook | moving toolchain discovery above the repository guard breaks the existing silence contract |
| TE13 | 5, 6 | a discovered executable under a directory containing spaces and glob characters is rendered as one shell-quoted PATH element | system hook | unquoted recovery output would split or expand before it repairs the environment |
| TE14 | 7, 10 | multiline, relative, nonexistent, or control-bearing clean-login output never appears inside a recovery command | system hook | treating arbitrary shell output as a trusted path creates a pasteable command from unauthenticated bytes |
| TE15 | 8 | a clean-login discovery process that outlives `bounds.EnvironmentDiscoveryTimeout` is killed as a process group, prints no recovery command, and SessionStart continues to exit zero | session inspection + system hook | an unbounded `bash -lc` can wedge SessionStart forever while every ordinary partial/missing case still passes |
| TE16 | 13 | a timed-out discovery shell that spawns a child whose trap writes a sentinel leaves no live child and no sentinel after SessionStart returns | session inspection + system hook | replacing the process-group kill with a direct or context-only parent kill still returns on time while the child survives |

Not covered: story 12 — real Codex-client/CLI WSL verification is phase-close
evidence, not a deterministic repository test; the diagnostic ticket requires
both observed invocations before any portability claim.

Cheapest wrong implementations and their reds: keep returning no Go phases →
TE1; reject every Go-less repository → TE3; preflight Go before manifest loading
→ TE4; replace PATH rather than prepend → TE6; execute the discovered binary →
TE7; warn on every session → TE9; guess a home path → TE10; interpolate raw
shell output → TE13 and TE14; wait forever on login initialization → TE15.
Kill only the login-shell parent → TE16.

### Edge inventory

- Required tool missing from PATH: main behavior, TE1 and TE5.
- `go.mod` absent versus present-but-empty: absence remains non-Go (TE3); a
  present regular file still declares the Go repository shape and requires Go
  (TE1's fixture includes the empty-file variant).
- Harness PATH already healthy versus partial propagation versus Go absent from
  both environments: TE9, TE5, and TE10.
- `ENVMAN_LOAD` absent, empty, or `loaded`: discovery is triggered by the
  observable missing-Go/clean-login-Go mismatch, while TE5 pins the reproduced
  `loaded` shape; the marker is evidence in the message, not the predicate's sole
  authority.
- Re-run idempotency: repeated SessionStart invocations print the same diagnosis
  and make no filesystem or environment mutation.
- Paths containing spaces or glob characters: TE13.
- Malformed discovery output: TE14 enumerates relative, nonexistent, multiline,
  and control-bearing values.
- Symlinked discovered executable: allowed as diagnostic data when it resolves to
  an existing executable; the hook never executes it, and the gate independently
  resolves the repaired PATH.
- Login shell exits nonzero: treat as no discovered Go, retain exit zero, and
  print no path-bearing recovery.
- Login shell exceeds `bounds.EnvironmentDiscoveryTimeout`: `bounds.Run` kills
  the whole process group at two seconds; TE15 pins continued exit-zero startup
  and absence of a recovery command, while TE16 plants a descendant sentinel to
  prove no shell-init child survives.
- **Won't handle:** shells other than Bash for clean-login reconstruction — the
  current Bench hook and reproduced WSL client are Bash surfaces; Bash remains an
  in-scope caller.
- **Won't handle:** automatically editing the harness environment or startup
  files — the hook is informational; the operator remains the surviving caller
  of the printed command.
- **Won't handle:** discovering arbitrary compilers named by project-owned phase
  manifests — their argv already reaches the runner and fails red when absent;
  the built-in table is the only surface that currently filters a required tool.
- **Won't handle:** installing Go or selecting among multiple Go versions — 0
  implementation edits and 0 gate runs here; toolchain management remains
  external.

## Ownership fences

- `specs/inherited-toolchain-environment/`
- `internal/gate/phases.go`
- `internal/gate/manifest.go`
- `internal/gate/branch_native_phases_test.go`
- `internal/sessioninspect/sessioninspect.go`
- `internal/sessioninspect/sessioninspect_test.go`
- `internal/bounds/bounds.go`
- `internal/bounds/bounds_test.go`
- `internal/systemtest/session_start_test.go`
- `CONTEXT.md`
- `projects/benchkit.md`
- `CHANGELOG.md`
- `capture/session-handoff.md`

## Out of scope

- Installing Go or choosing a Go distribution — separate external capability,
  0 repository edits and 0 gate runs.
- General environment reconstruction for Node, Python, or project-manifest
  commands — approximately 20 edits and 5 gate runs after a separate decision
  source enumerates their producers and trust assumptions.
- Persistently mutating Codex, WSL, Envman, or shell startup configuration —
  approximately 8 external edits and no repository gate; the hook remains
  diagnostic only.
- Changing the optional-capability posture for FIFO, privilege, or shellcheck
  skips — approximately 12 edits and 3 gate runs; those classes are not required
  source toolchains.

## Further notes

The reproduced commit failure (`gate-20260822T145702.664532004Z-8472`) reached
private executable selection and recorded `env: ‘go’: Permission denied`; after
prepending `/home/mgibs/.local/opt/go/bin` to the existing PATH, the identical
`bench commit` reached all six phases and landed green. This is supporting
evidence, not a user-specific path the implementation may search.

The reviewer-approved ticket frontier is two independently green, disjoint
slices with no blockers: `fail-built-in-go-table-closed.md` and
`diagnose-partial-session-environment.md`. The user's 2026-08-22 instruction to
run the phase, use Terra/medium review, and commit the accepted result to `main`
is the standing batch approval for this spec-and-tickets pair; the three-iteration
verification log is its post-hoc veto surface.
