# Expose the complete hook-health record

Blocked by: none
Ownership fence: `internal/adopt/link_hook.go`, `internal/adopt/link_hook_test.go`
Integration surfaces: hook-health record→render-doctor-hook-health.md, repair-stale-hook-with-doctor.md, render-guards-hook-health.md, signal-stale-hook-status.md, converge-symlinked-hook-lifecycle.md, and converge-live-hook-installation.md; existing `PrePushMarker` consumers→`internal/adopt/unlink.go` + H7
Contracts: the complete hook-health record crosses `internal/adopt/link_hook.go`→doctor, guards, and status, asserted by H1 against the real classifier; hook classification and expected bytes cross `internal/adopt/link_hook.go`→the link transaction, asserted by H4 against the installed hook
Closure: H1/complete-record, H2/live-branch, H2/baked-branch, H3/live-provenance, H3/baked-provenance, H3/fallback-provenance, H4/template-currency, H5/not-applicable-states, H6/no-remote-git, H7/state-membership, E1/literal-hook-path, E2/special-file-mode, E3/trailing-newline, E4/empty-versus-absent, E5/malformed-prefix, E6/missing-git

## What to build

Expose one hook-health computation that classifies the effective hook path before reading it and returns the existing state, path, effective branch, provenance, fallback provenance, and three-valued currency. Recover the baked token from installed bytes for template comparison, keep ambient inspection offline, and leave the four existing state meanings unchanged.

## Acceptance

- [ ] [H1] One exported computation returns state, hook path, effective branch, provenance, fallback provenance, and currency for a managed hook.
- [ ] [H2] Effective branch follows live `origin/HEAD` and falls back to the token baked in the installed hook without reinstalling.
- [ ] [H3] Provenance distinguishes live resolution, a baked real branch, and the bare fallback token.
- [ ] [H4] Currency compares exact installed bytes with the embedded template substituted by the token recovered from those installed bytes.
- [ ] [H5] Absent, foreign, and diverted hooks report not-applicable currency rather than stale.
- [ ] [H6] Hook-health inspection invokes no remote-contacting git subcommand.
- [ ] [H7] The four existing `PrePushState` values retain their meanings and membership.
- [ ] [E1] A `core.hooksPath` containing spaces and glob characters is read literally.
- [ ] [E2] FIFOs, devices, and sockets at the hook path are rejected by file mode before any read.
- [ ] [E3] A managed hook missing only its final newline reports stale.
- [ ] [E4] An empty hook is foreign and remains distinct from an absent hook.
- [ ] [E5] A marker-bearing hook with a malformed pre-token prefix reports stale without panicking.
- [ ] [E6] A repository with no `git` on `PATH` degrades without panic.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| H1/complete-record | omit currency from the returned record | hook-health unit test | plant a managed hook, call the exported computation, run `go test ./internal/adopt`, expect the whole-record mismatch |
| H2/live-branch | ignore a resolvable `origin/HEAD` | hook-health unit test | set `origin/HEAD`, run the focused test, expect the branch mismatch |
| H2/baked-branch | reuse the live branch after `origin/HEAD` is removed | hook-health unit test | unset the ref without reinstalling, run the focused test, expect the baked-branch mismatch |
| H3/live-provenance | label a live resolution baked | hook-health unit test | set `origin/HEAD`, run the provenance case, expect the provenance mismatch |
| H3/baked-provenance | label a baked real branch live | hook-health unit test | remove `origin/HEAD`, run the baked case, expect the provenance mismatch |
| H3/fallback-provenance | collapse the bare fallback into an ordinary baked branch | hook-health unit test | plant a fallback-token hook, run the fallback case, expect the fallback flag mismatch |
| H4/template-currency | substitute the current repository branch instead of the installed token | hook-health unit test | plant a current hook baked for a different branch, run the currency case, expect false-stale detection |
| H5/not-applicable-states | report stale for a foreign hook | hook-health unit test | plant absent, foreign, and diverted fixtures, run the table, expect the foreign currency mismatch |
| H6/no-remote-git | invoke `ls-remote` during inspection | PATH-recorder test | replace git on `PATH` with the recorder, inspect health, expect the forbidden-subcommand failure |
| H7/state-membership | add a fifth state for stale | existing classifier contract | change the state set, run `go test ./internal/adopt`, expect the state-membership failure |
| E1/literal-hook-path | clean or glob-expand `core.hooksPath` | hook-health unit test | use `hooks [one]`, inspect health, expect the wrong-path failure |
| E2/special-file-mode | read a FIFO before classifying its mode | bounded special-file test | plant a FIFO with no writer, run the focused test under its timeout, expect the mode refusal rather than a hang |
| E3/trailing-newline | trim bytes before comparison | hook-health unit test | remove only the final newline, inspect health, expect current-versus-stale mismatch |
| E4/empty-versus-absent | classify every read error or empty read as absent | hook-health unit test | compare empty and absent fixtures, expect distinct-state failure |
| E5/malformed-prefix | treat any marker-bearing bytes as extractable and current | hook-health unit test | corrupt the prefix before the token, inspect health, expect stale rather than current |
| E6/missing-git | assume git output is always available | stripped-PATH unit test | remove git from `PATH`, inspect health, expect a returned record rather than panic |
