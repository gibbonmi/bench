# Refuse a cd into the Bench pool path

Blocked by: 01-name-the-leading-operator-in-the-follow-on-refusal.md
Writes: internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go, cmd/bench/main.go, .bench/hooks/block-bench-follow-on.sh, internal/guards/guards_test.go, internal/systemtest/bench_follow_on_test.go

## What to build

`bench worktree exec "<label>" -- <command>` is the only command form for a
Bench worktree, and `AGENTS.md` says so. On 2026-08-30 a write delegate ran
`cd <pool path>` on almost every call, with that rule in its charge. Prose
did not stop the habit, so the guard does.

The follow-on guard verb `guardBenchFollowOn` in `cmd/bench/main.go` already
reads every Bash envelope. It already tests the command text for the pool
prefix `poolkey.Pools(home)` before it resolves any root, in
`recordFollowOn`. Add a second denial on the same path. A simple command
with the command word `cd` and an absolute path argument under
`poolkey.Pools(home)` is the denied shape. When any simple command in the
stream has that shape, the verb refuses with exit 2 and its own message:

`BLOCKED: a Bench worktree runs through bench worktree exec "<label>" -- <command>; never cd into the pool path. target=<path>`

Put the classification in `internal/benchguard` beside `Classify`, as its own
function that takes the pools prefix as an argument. Do not resolve the path
against the file system: this is an honest-mistake layer, and a variable or a
relative `cd` inside a wrapper string stays allowed. A `cd` inside a
`bench worktree exec` child, for example `sh -c 'cd sub && go test ./...'`,
stays allowed because its target is relative. The follow-on verdict keeps
its precedence: the verb tests the `cd` denial first, because a `cd` into
the pool followed by a Bench call has two faults and the `cd` is the cause.

The hook file `.bench/hooks/block-bench-follow-on.sh` advertises what it
denies in its `# denies:` header line. Widen that line to name both denials,
so the advertisement and the enforcement agree. `internal/guards/guards_test.go`
pins the header text; update the pin to the new line. Add one real-path case
to `internal/systemtest/bench_follow_on_test.go`. It runs the hook script on
a `cd <pool path>` envelope and asserts exit 2 with the new message. Mirror
the existing cases there.

## Acceptance

- [ ] A Bash envelope whose command is `cd <pools>/<pool>/<assignment>` exits 2 from the hook with the new message and `target=` naming the path.
- [ ] A Bash envelope whose command is `cd <pools>/<pool>/<assignment> && go test ./...` exits 2 with the new message, not the follow-on message.
- [ ] `bench worktree exec "x" -- sh -c 'cd sub && go test ./...'` exits 0 from the hook.
- [ ] `cd /tmp` and `cd "$W"` exit 0 from the hook.
- [ ] The hook's `# denies:` header names the `cd` denial, and `internal/guards/guards_test.go` pins the new header text.
- [ ] A new table test in `internal/benchguard/benchguard_test.go` covers the four rows above, and the delegate records the first row red before the fix.
- [ ] `go test ./internal/benchguard/ ./internal/guards/ -parallel 2` passes, and the new system case passes under `-tags=system` with `BENCH_KIT` and `BENCH_RUN_BINARY` set the way the existing cases run.
