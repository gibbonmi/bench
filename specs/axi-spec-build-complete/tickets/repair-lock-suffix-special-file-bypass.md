# Repair the lock-suffix special-file bypass

Blocked by: none
Ownership fence: `internal/specbuild/state.go`, `internal/specbuild/runs_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: LS1/lock-suffix-checked-after-special-file-classification

## What to build

Close the accepted Spec and Coverage findings (P2, C2) from the Terra/xhigh
review of candidate `5ae4029d6207fd1bc1f75b85612d2c4492baae68`:
`Service.Runs`'s entry loop (`internal/specbuild/state.go` around lines
319-325) checks the `.json.lock` filename suffix *before* the symlink and
non-regular checks:

    if strings.HasSuffix(name, ".json.lock") {
        diagnostics = append(diagnostics, RunDiagnostic{Entry: name, State: "lock_skipped"})
        continue
    }
    if entry.Type()&os.ModeSymlink != 0 { ... }
    if !entry.Type().IsRegular() { ... }

A directory entry named `*.json.lock` is therefore classified `lock_skipped`
— the state owner's own expected, harmless companion file — even when it is
actually a symlink or a FIFO, never reaching the special-file checks that
would otherwise diagnose it. This violates the spec's edge inventory
requirement that the family-home enumerator "refuses non-regular and
symlinked entries into named diagnostic rows" (`specs/axi-spec-build-complete/spec.md`
around lines 26, 91) regardless of filename — a name pattern is not an
identity check. `internal/specbuild/runs_test.go`'s existing coverage
(around lines 75-102, 113-229) exercises a regular `.json.lock` companion and
special `*.json` entries separately, but no entry that is both
lock-suffixed and special, so this exact bypass is untested.

Fix: reorder the classification so the symlink and non-regular checks run
*before* the `.json.lock` suffix check — only an entry that is both a regular
file *and* `.json.lock`-suffixed reaches `lock_skipped`; a symlink or
non-regular entry with that suffix is diagnosed `symlink` or `non_regular` as
it would be for any other name. Do not change the relative order of the
symlink check versus the non-regular check, or any check after the suffix
check (`foreign`, the read, `malformed`, `nondigest_name`, healthy) — this is
a pure reorder of the lock-suffix branch to after the two special-file
branches, nothing else.

New coverage in `internal/specbuild/runs_test.go`: construct a symlink named
`<digest>.json.lock` (pointing anywhere; its target is irrelevant since the
symlink check fires first) alongside a real `.json.lock` companion for a
healthy retained run, and require the symlink surfaces as its own `symlink`
diagnostic while the real lock companion still classifies `lock_skipped` and
the healthy run still renders. Add a FIFO case (reusing `syscallMkfifo`, which
already skips the test when FIFOs are unavailable on the platform) named
`<digest>.json.lock` the same way, requiring `non_regular`.

## Acceptance

- [ ] [LS1] (covers local) (P2, C2) a directory entry named `*.json.lock`
  that is actually a symlink or a FIFO is diagnosed `symlink` or
  `non_regular` respectively, never `lock_skipped`; a genuine regular-file
  `.json.lock` companion still classifies `lock_skipped` exactly as before.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LS1/lock-suffix-checked-after-special-file-classification | move the `.json.lock` suffix check back before the symlink/non-regular checks | focused `Runs()` test | list a retained-state directory holding a symlinked and a FIFO entry each named `<digest>.json.lock`, alongside a genuine regular `.json.lock` companion for a healthy run, and require the symlink and FIFO surface their own special-file diagnostics rather than `lock_skipped` |
