# Review pickup: worktree-exec-comfort

Frozen base `72236666`, reviewed tip `1ccdb060`. Three axes ran on `opus` /
medium in fresh contexts. The repair tickets `repair-worktree-surfaces.md` and
`repair-guard-and-inventories.md` in `specs/worktree-exec-comfort/tickets/`
carry every accepted finding. The fix commit that closes them deletes this file.

## Standards

Count 7, worst issue: the `show` cancel path diverges from `exec` while its
comment asserts parity.

- `auto-fix` — `internal/worktree/exec.go:40` repeats the assignment regex of
  `internal/shellcommand/shellcommand.go:51`; one exported predicate serves both.
- `auto-fix` — `internal/worktree/show.go:38` claims the interrupt reports as
  `exec` does; `exec.CommandContext` kills with SIGKILL, so the code is 137.
  The comment drops the claim, and the spec records the cancel edge as
  Won't handle.
- `auto-fix` — `internal/worktree/identity_component_test.go:335` hardcodes the
  `next=` default that `nextList` in `identifier_operand_test.go` owns.
- `no-op` — `CONTEXT.md` restates the exec exception; story 43 asks for it.
- `auto-fix` — `benchguard.Classify` only wraps `Judge`, and only tests call it.
  The spec names `Classify` as the verdict function, so `Judge` takes that name
  and returns the `Verdict`.
- `auto-fix` — `CHANGELOG.md` gains an `Unreleased` entry; the fence gains the
  file.
- `no-op` — the spec's seam map keeps `benchguard.Classify` after the rename.

## Spec

Count 5, worst issue: X8's child-level behavior has no test.

- `ask-user` — X2 and X9 drive a test-local copy of the commit grammar.
  `internal/commit` imports `internal/usage`, so the real value is out of reach.
  The reviewer decides whether a row moves to `internal/commit`.
- `auto-fix` — X8: two `--env` values reach the child in one test.
- `auto-fix` — F11: the `next=` line always quotes the path, and `list` quotes
  only an unsafe path. Both surfaces use `axi.ShellQuote`, and the spec says
  "shell-quoted".
- `auto-fix` — S8: the assertion names the full `show` grammar.
- `auto-fix` — the `show` comment finding above.

## Coverage

Count 6, worst issue: a `'` in the worktree path breaks the `next=` line.

- `auto-fix` — a `'` in the path breaks `line()`; `axi.ShellQuote` fixes it,
  and an F9 variant pins it.
- `no-op` — the row's landed cell and the global landed flag differ only when
  the default branch cannot resolve. The release route is right for that record.
- `auto-fix` — the `show` operand takes no `lineSafe` check; the Edge inventory
  promises one.
- `auto-fix` — `bench worktree show --help` exits 2; the verb parses through
  `usage.Parse`, so the sole help spelling answers at exit 0.
- `no-op` — the `show` cancel code is recorded as Won't handle.
- `no-op` — the `worktree:` line prints the stored path raw; the ledger path is
  the verb's own output.
