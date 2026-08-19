# Guard the scaffolded inventory call on the inventory directory

Blocked by: none
Writes: internal/adopt/init.go, internal/adopt/doctor_rows.go, internal/adopt/setup_report.go, internal/adopt/adopt_test.go

## What to build

`scaffoldGate()` wraps its inventory call so a repository with no
`tests/canary` directory skips validation and one with the directory still
validates it:

    if [ -d "$root/tests/canary" ]; then
      "$bench" canary "$root" || err "canary inventory validation failed"
    fi

The path is rooted at `$root`, never relative. Directory existence is the whole
predicate: a present-but-empty directory falls through to `bench canary`, which
reports the empty inventory and reds the gate. The sentinel line and its marker
stay exactly as they are. `bench canary`, its messages, and `Inventory`'s callers
do not change. `bench init` keeps writing the same `scaffoldGate()` text.

The sentinel marker constant becomes exported from `internal/adopt` (the doctor
row and the setup remedy keep reading it) so ticket 03's journey can remove the
sentinel line by the one marker rather than a second literal.

`TestScaffoldGateUsesCanarySubcommand`'s canary-line expectation moves to the
guarded form; it keeps asserting the sentinel and the retired-API absences.

## Acceptance

- [ ] `scaffoldGate()` contains the guarded inventory call rooted at `$root`, and `bench init` writes that same text into a fresh repository (SG1).
- [ ] `scaffoldGate()` still carries the sentinel line and marker (SG2).
- [ ] `bench canary` on an absent and on an empty `tests/canary` still reports `canary fixture inventory is empty` and exits 1 — the existing inventory test stays green unchanged (SG3).
- [ ] the sentinel marker is readable from outside the package under one exported name, and the doctor row and setup remedy print that same value.
