# Tighten the citation grammar for mentions and subtests

Blocked by: grade-citation-execution.md
Writes: internal/coverage/citations.go, internal/coverage/citations_test.go
Covers: CE8, CE9, CE10, CE11

## What to build

Two rules extend the one citation grammar in the coverage package. The
Blocked-by edge serializes writes to one file. It is not a semantic
prerequisite.

The mention rule grades a backticked `_test.go` path in a seam cell that carries
no parenthesized name list. Such a token is a violation, because a mention poses
as evidence. The rule reads the seam cell alone. A backticked test path in any
other cell stays ungraded, so honest prose in a why cell does not red. An empty
parenthesized list stays a non-citation, which is the delivered behavior.

The subtest rule resolves the segment after the slash in a cited name. The
segment must appear as a `t.Run` string literal in the cited file. A segment that
is absent is a violation. A file that calls `t.Run` with a non-literal name is
exempt from segment resolution, so a table-driven suite cannot false-red. The
parent function still resolves in both cases.

## Acceptance

- [ ] CE8 — a backticked test path in a seam cell without a name list is a
      violation.
- [ ] CE9 — a backticked test path in the why cell is not graded.
- [ ] CE10 — a cited subtest segment absent from the file's `t.Run` string
      literals is a violation.
- [ ] CE11 — a file with a non-literal `t.Run` name is exempt from segment
      resolution.
- [ ] `bench gate` stays green.

## Delegate charge

You work in the Bench repo on the `citation-execution-proof` spec. Read
`specs/citation-execution-proof/spec.md` first. Then read
`internal/coverage/citations.go` and `internal/coverage/citations_test.go` in
full. Include the execution arm a sibling ticket already landed. Extend the one
citation grammar in that file. Do not add a second parser.

Report a violation for a backticked `_test.go` path in the seam cell with no
parenthesized name list. Keep an empty parenthesized list a non-citation. Grade
the seam cell only. A backticked test path in the why cell stays ungraded.

Resolve a cited name's post-slash segment against the `t.Run` string literals of
the cited file. Report a violation when the segment is absent. Skip segment
resolution for a file that calls `t.Run` with a non-literal name. Keep the
parent-function resolution in both arms.

Add `TestMentionIsNotACitation` and `TestSubtestSegmentResolves` in
`internal/coverage`. Assert the exact violation text. Keep the delivered
resolution and execution messages unchanged.

Run only `bench worktree exec <label> -- go test ./internal/coverage/`. Do not
commit. Do not edit the spec.
