# Build the parallel census and the static pins

Line: opus / medium.

Blocked by: none
Writes: internal/worktree/parallel_census_test.go (new), internal/worktree/effect_census_test.go

## What to build

The census parses the package's `_test.go` files with the Go AST. It builds the
call graph over test-file functions. A top-level test is serial when its body,
or any test-file function it reaches, calls `bindEnv` or `chdir`. Every other
top-level test is eligible.

The census reports each eligible test that lacks `t.Parallel()`. It also reports
each serial test that carries the call. Every report names the file and the
line. The census lives beside the harness effect census and shares its file
walk, so it skips a special file as that walk does.

Unit tests drive the census over synthetic file sets in a temporary directory. A
planted omission turns the census red without an edit to the live tree. The
live-tree report lands with the last marks ticket, so this ticket adds no
live-tree census test.

Add the count pin and the package-clause pin here. Keep the harness effect
census as it is, and pin its refusal.

## Acceptance

- [ ] WF02 pins the file and the line that the census reports for a synthetic eligible test without `t.Parallel()`.
- [ ] WF03 pins a synthetic test as serial when the helper it calls calls `bindEnv`.
- [ ] WF04 pins a synthetic test as serial when a subtest closure calls `chdir`.
- [ ] WF05 pins the pair that the census reports for a synthetic serial test that calls `t.Parallel()`.
- [ ] WF12 pins at least 334 top-level tests in the package.
- [ ] WF13 pins the harness effect census on an `exec.Command` outside the harness files.
- [ ] WF14 pins `package worktree` in every `internal/worktree` test file.
