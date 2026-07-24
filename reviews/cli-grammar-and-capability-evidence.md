# Review — cli-grammar-and-capability-evidence

Diff reviewed: `9732ebe..fde0b86` on `main` (FT87 slice 3, 78 files, +2792/−278).
Three axes, run as parallel read-only delegates.

Rounds 1 and 2 closed the hard Standards violations, the empty-string commit path,
and the spec's stale skip-evidence transport. Round 3 (`6b44e26`, `6164844`) closed
the remaining defects: the swallowed skip-log read error, the stale `bench models`
help comment, the comment-register violations in the two AXI tests, the
`ideaGrammar` doc comment, import grouping in six files, the profile's missing
`BENCH_REQUIRE_CAPABILITIES` knob and three conformance checks, row 20's stale
red-signal text, `gateEnv`'s unasserted skip-log isolation, the deleted-named-directory
commit bug, and `TestDeadline`'s untested clamp and silent overflow.

What remains below is not a defect list. Every entry is a decision that is the
reviewer's to make — an exemption to bless or a behavior to name — and each one
was left deliberately rather than fixed.

## Standards

- **Three hand-rolled help parsers stand outside the grammar.** `internal/spec/spec.go:232`,
  `internal/worktree/list.go:18`, and `cmd/bench/main.go:328` each parse help themselves,
  blessed permanently by `subcommand_routing_test.go:42`'s `whyNested` note ("each leaf
  owns its own grammar"). Either a scoped follow-up folds them into `usage.Parse`, or the
  exemption becomes explicit and bounded. `cmd/bench/main.go`'s worktree dispatch is
  already a closed exemption; the other two are not.

## Spec

- **The routing registry exempts six adopt subcommands for a reason nothing grades.**
  `internal/conformance/subcommand_routing_test.go:78-89` exempts `setup`, `link`, `init`,
  `doctor`, `unlink`, and `upgrade` under `whyNested` ("dispatches a subcommand tree rather
  than a flat argv"), but `internal/adopt/adopt.go:15-27` dispatches each to a leaf and
  `internal/adopt/doctor.go:181-189` is a hand-rolled flat `switch` whose default makes
  `bench doctor -h` exit 2. The fail-closed check the spec demanded is opened by an
  exemption reason that is factually wrong for at least `doctor`. Fails stories 2 and 4.
  Fixing it means either routing those six or narrowing the exemption to a reason that
  holds.

- **`bench spec implemented -- -x.md` is still rejected.** `internal/spec/spec.go:229`
  (`specArg`) has no `--` case, so the leading-dash path stays inexpressible on a
  subcommand whose sole argument is a path — exactly the inexpressibility story 5 exists
  to close. Same call as the hand-rolled parsers above; they are one decision if you route
  `spec` through the grammar.

## Coverage

- **`bench commit -m x .` can never succeed.** `internal/commit/commit.go:254-259` —
  `underAny` compares against `"./"` while porcelain paths are `internal/x.go`, so every
  changed file is listed as an offender. `.` is the obvious spelling of "everything I
  changed" and no test names it. Whether `.` should be accepted at all is a policy call:
  the block-check exists to make a commit name what it stages, and `.` names nothing.

- **A value flag consumes a flag-shaped token with no look-ahead.**
  `internal/usage/parse.go:83-88` — `bench commit --spec -m "msg" file` sets
  `specSlug = "-m"` and fails later at spec resolution; `bench commit -m --spec x file`
  commits with the message `--spec`. Both are getopt's documented behavior, so the choice
  is to keep it and say so in the grammar's contract, or return `MissingArg` when the next
  token looks like a declared flag. No case in `parse_test.go` gives a value flag a
  `-`-prefixed value either way.

- **`Render` interpolates `Reason` verbatim.** `internal/capability/capability.go:108,110` —
  a reason containing `\n` produces a two-line record that `readSkipTally`'s
  `strings.Split(..., "\n")` counts as two skips (`internal/gate/capability_skips.go:73`),
  and a reason pushing the line past 4096 bytes breaks the single-`Write` atomicity the
  design rests on (`capability.go:159-164` states the bound; nothing enforces it). Both
  inputs are authored by test code in this repo rather than by a user, which is why this
  was left: whether an in-repo-authored input earns validation is your call.

Note: the Edge inventory's FIFO dismissal ("git lists only tracked and untracked
regular paths") is wrong — `git status --untracked-files=all` does list an untracked
FIFO — but git refuses it at `add`, so the failure is loud. Flagged, not a finding.
