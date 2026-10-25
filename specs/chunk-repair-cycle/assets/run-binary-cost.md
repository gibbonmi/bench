# The run-binary cost of a Bench child of exec

Covers: RC76

## Result

The net saving for one test-only repair cycle is positive. The build does not stop RC-C7.

The saving is small. The ticket formula gives 4.6 seconds. The formula with the child's own source check gives 1.5 seconds. The cycle target is 11 to 14 minutes.

## Subject and machine

The measurement ran in the integration worktree at source tip `7ec672d9a1e05701a63ab1039adecfc22777e7a4`. The chunk base is `ef387b3f9a6f99831772abe809428d6d14b19f3a`.

The host is WSL2 on Linux 6.6 with 16 processors and Go 1.25.14. The Go build cache was warm. A cold build cache was not measured.

The throwaway tests were Go test files. A Go build overlay added them to the packages at run time. No file was written into the tree, and nothing was committed.

## Measured values

All times are wall times. A median covers every sample except the first.

| Value | Call that ran | Samples | First | Minimum | Median | Maximum |
|---|---|---|---|---|---|---|
| Digest cost | `freshness.Digest(root)` on 611 build inputs | 50 | 199.6 ms | 193.2 ms | 204.5 ms | 220.7 ms |
| Selection cost today | `runbinary.Own(ctx, root)` with `BENCH_RUN_BINARY` unset | 10 | 965.8 ms | 958.7 ms | 982.2 ms | 1010.2 ms |
| Child check of an inherited binary | `runbinary.Factory{}.Inherit(root, path)` | 30 | 236.5 ms | 225.9 ms | 233.7 ms | 254.5 ms |
| Seal and byte check only | `freshness.VerifyExecutable(path)` | 30 | 37.1 ms | 32.5 ms | 34.4 ms | 36.9 ms |

The commands were these. Each one ran through `bench worktree exec FT445-build -- sh -c '...'` with `T12_ROOT="$PWD"` and `-overlay` set to the overlay file.

```
go test -overlay overlay.json -count=1 -run TestT12DigestTiming -v ./internal/freshness
go test -overlay overlay.json -count=1 -run TestT12OwnTiming -v ./internal/runbinary
go test -overlay overlay.json -count=1 -run TestT12InheritTiming -v ./internal/runbinary
```

The digest test times 50 calls of `freshness.Digest` in a row. The selection test times 10 calls of `runbinary.Own` and closes each selection after the timer stops. The inherit test owns one selection, then times 30 calls of each check.

## What the values mean

`freshness.Digest` runs `go list -deps ./cmd/bench` and hashes the files it lists (`internal/freshness/freshness_buildinputs.go:50-75`, `internal/freshness/freshness.go:35-61`). The list holds no `_test.go` file. A test-only repair therefore keeps a built `dist/bench` fresh.

The exec check of ticket 13 is `freshness.VerifyAt` (`internal/freshness/freshness_verify.go:23-40`). It pays one digest and one seal check. The inherit row measures that same chain.

Exec strips `BENCH_RUN_BINARY` today (`internal/worktree/exec.go:244-252`). A child that selects then calls `runbinary.Own` through one of three routes:

- the lane of `bench commit` (`internal/gate/lane.go:277-297`);
- `bench test` and `bench probe` (`internal/testreport/command.go:49-50`);
- `bench gate` (`internal/gate/run_transaction.go:32-40`).

`Own` runs `scripts/go-build.sh`, then validates the new binary against the source digest (`internal/runbinary/runbinary.go:116-150`). The selection row includes that validation.

A child that does not select pays nothing for a selection. These children are `bench record`, `bench preflight`, and `bench spec`. A reused binary adds the exec check to such a child and no saving.

A reused binary has one cost in the child. The child's `Inherit` checks the source digest again (`internal/runbinary/runbinary.go:155-183`). A reuse therefore pays two checks and avoids one build and one check.

## Exec calls of one test-only repair cycle

The walk follows `.agents/commands/bench-implement-spec.md` at this base. It covers a one-file test-only repair with a biting coordinator probe and a confirmation record. Each step runs in the integration worktree as `bench worktree exec <target> -- bench ...`.

| Step | Child of exec | Selects a binary | Line in the command |
|---|---|---|---|
| 1 | `bench record assignment` | no | 68 |
| 2 | `bench commit` (assignment plan commit) | yes, lane | 68 |
| 3 | `bench preflight build <slug> --charge --ticket ...` | no | 21 |
| 4 | `bench preflight evidence <id> --check-current` | no | 29 |
| 5 | `bench test --package ...` (focused check) | yes | 25, 48 |
| 6 | `bench test --check ...` (one check that the ticket owes) | yes | 68 |
| 7 | `bench commit` (repair commit) | yes, lane | 48, 68 |
| 8 | `bench record excerpt` (run of step 5) | no | 66 |
| 9 | `bench record excerpt` (run of step 6) | no | 66 |
| 10 | `bench record ticket` | no | 68 |
| 11 | `bench preflight build <slug> --base ... --source-tip ...` | no | 52 |
| 12 | `bench probe` (coordinator probe) | yes | 66 |
| 13 | `bench record confirmation` | no | 58 of the repair policy |
| 14 | `bench commit` (record commit) | yes, lane | 66 |
| 15 | `bench gate --checkpoint ... --prepare` (prepared run) | yes | 66 |
| 16 | `bench gate --checkpoint ...` | yes | 68 |

The count is 16 `bench` children of exec. Eight of them select a binary: three commits, two tests, one probe, and two gate runs. Eight do not select a binary.

The walk makes four assumptions:

- The repaired test file sits inside the `Writes:` fence, so no `bench spec widen` call runs.
- The repair changes no production file, so no `bench worktree build --install` call runs.
- The cycle runs one check that the ticket owes.
- `bench handoff` runs once for each chunk checkpoint, so the walk leaves it out.

An added check adds one selecting child. An added record call adds one non-selecting child.

## Net saving for each repair

The ticket formula is the selection cost times the count of selecting children, less the digest cost times the count of `bench` children.

| Formula | Computation | Result |
|---|---|---|
| Ticket formula | 982.2 ms x 8 - 204.5 ms x 16 | +4.58 s |
| With the child's own check | 7 x (982.2 - 233.7 - 233.7) ms - 1 x 233.7 ms - 8 x 233.7 ms | +1.50 s |

The second row treats each exec check as the 233.7 ms inherit chain. Each selecting child that gains pays a second check of 233.7 ms in `Inherit`. The row credits no saving to the probe and charges it one exec check. Section "Facts for ticket 13" gives the reason.

Both results are positive. The net saving of each repair is 4.6 seconds at most and 1.5 seconds at least.

## Stop decision

The net saving is positive in both formulas. The ticket rule asks only for the sign. RC-C7 continues with ticket 13.

The saving is under 1 percent of the 11-minute lower bound of the cycle target. Ticket 13 does not carry that target alone. RC81 grades the target at the final reconciliation.

## Wrapper route for an inherited `BENCH_RUN_BINARY`

The route is as the spec states. These are the lines read.

- `bin/bench.sh:298-313`. When the variable is present, even empty, the wrapper checks four facts. The value is absolute. The file is regular and executable. The file is not a symbolic link. The value equals its physical path. Then the wrapper runs the value with `exec "$bin" "$@"`. It checks no source digest. An empty value fails the absolute check, so exec must omit the variable and must never set it empty.
- `internal/runbinary/runbinary.go:155-183`. In the Go child, `ReuseOrOwn` calls `Select`. A present value goes to `Inherit` and then `Validate`. With a named source root, `canonicalVerify` calls `freshness.VerifyAt`, which checks the source digest. With no source root, it calls `freshness.VerifyExecutable`, which checks only the seal pair.
- `internal/systemtest/owner_environment_test.go:36-39`. The wrapper's own binary search runs only when the variable is absent.

The answer is confirmed. The wrapper checks the shape only. The Go child checks the source digest.

## Facts for ticket 13

These facts come from reading code. No run confirmed them.

- A reused binary runs every verb, not only the verbs that select. The wrapper executes the value at line 313 for any verb. A verb that does not select, such as `bench record ticket`, therefore runs the worktree's `dist/bench`. The exec check is the only source check for that verb.
- `bench probe` mutates production source. `HeldRuns.selection` refuses an inherited binary that no longer validates (`internal/testreport/held_runs.go:61-88`). Today the same probe builds a fresh binary for each mutant. With a reused binary, the first mutant run can refuse. Ticket 13 must keep `bench probe` working under exec.
- Eight of the 16 children do not select. They pay the exec check and gain nothing. A check inside the child would avoid that cost.
