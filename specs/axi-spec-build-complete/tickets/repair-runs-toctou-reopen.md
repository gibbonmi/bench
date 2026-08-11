# Repair the retained-run reopen-by-path timing window

Blocked by: none
Ownership fence: `internal/specbuild/state.go`, `internal/specbuild/runs_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: TC1/stat-immediately-before-each-read

## What to build

Close the accepted Coverage finding (C1) from the Terra/xhigh review of
candidate `399dca908c7b1e1a4162eb7625e497dfb6786750`: `Service.Runs`
(`internal/specbuild/state.go:303-360`) classifies each entry from the
`os.ReadDir` snapshot taken on line 308 — `entry.Type()&os.ModeSymlink` on line
323, `entry.Type().IsRegular()` on line 327 — and then reopens the path by name
for `os.ReadFile(filepath.Join(directory, name))` on line 335. Nothing
re-establishes between those steps that the path still holds what the snapshot
described, so a path swapped to a symlink, FIFO, device, or socket in that
window is read through rather than diagnosed. The FIFO case is the sharp one:
`os.ReadFile` on a FIFO with no writer blocks, and the family home is the
ambient orientation call, so the hang is exactly the failure
`projects/benchkit.md:129-131` names — "special files in any discovered path
... must be rejected before reading so neither static inspection nor an ambient
command can block" — alongside `:132-134`'s dangling-link case, where a plain
read reports not-found and a reader that does not stat first turns a broken
link into an authoritative empty state. The spec's edge inventory
(`spec.md:91`) states the requirement for this one new discovered-path reader:
it "stats before reading, refuses non-regular and symlinked entries into named
diagnostic rows".

`os.Lstat` the entry's full path immediately before `os.ReadFile` and classify
from that result, not from the `ReadDir` snapshot: a symlink diagnoses
`symlink`, any other non-regular mode diagnoses `non_regular`, and an `Lstat`
error diagnoses `unreadable`, so the diagnostic row reports what was actually
at the path at read time. The existing snapshot checks stay as the cheap first
pass. Reuse the three existing class names rather than minting a
timing-window class: the observable fact is identical — the path held a symlink
or a non-regular file when the reader looked — and a fourth name would split one
fact across two vocabularies, leaving every consumer of the diagnostic states
(`internal/specbuild/runs_test.go:97`, the family-home renderer, the
old-to-new fixtures) to learn a distinction that describes when the swap
happened rather than what was found.

The race is deterministic to drive without inventing a test double. `Service`
already carries the `fault func(string) error` hook (`lifecycle.go:89`) reached
through `faultAt` (`checkpoint.go:24-29`), nil in production and used at named
points across the package (`assign/state`, `checkpoint/ref`, `reclaim/apply`,
`refresh/preserve`). Add one named point in `Runs` between classification and
read; the test's callback matches that point, replaces the entry with a FIFO or
a symlink, and returns nil, so the read arrives at a swapped path with no
concurrency. `internal/specbuild/runs_test.go` already builds both fixtures
(`syscallMkfifo` at `:108-110`, the symlink and special-file cases at `:50`),
so the new case composes what is there.

## Acceptance

- [ ] [TC1] (covers local) (C1) an entry that classifies as a healthy regular
  `*.json` record from the directory listing but is replaced by a symlink or a
  FIFO before the read is diagnosed under the class that names what is at the
  path at read time — `symlink` or `non_regular` — and is never read through,
  never blocks, and never counts among the healthy runs; unswapped healthy,
  hostile, and capped entries render exactly as they do today, and the fault
  hook stays nil in production.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TC1/stat-immediately-before-each-read | drop the pre-read `os.Lstat` and let `os.ReadFile` trust the `os.ReadDir` snapshot | focused `Runs()` test | list a retained-state directory holding one healthy record, swap that entry to a FIFO and, in a second case, to a symlink at the named fault point between classification and read, and require the swapped entry to surface its `non_regular` or `symlink` diagnostic row with the call returning rather than blocking, while the other retained runs stay healthy |
