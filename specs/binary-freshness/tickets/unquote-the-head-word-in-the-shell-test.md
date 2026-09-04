# Unquote the head word in the shell test

Blocked by: none
Writes: .bench/lib/resolve-bench.sh, internal/conformance/guard_classifier_table_test.go, internal/systemtest/bench_follow_on_test.go, internal/benchguard/benchguard_test.go
Covers: BF23, BF25

## What to build

This ticket repairs review finding C1.

Under a stale binary the shim classifies a Bench call through the shell word
test. That test runs `set -- $1` under `set -f`, which never removes a quote.
`benchguard.InvokesBench` runs a real lexer, which does. So the two
classifiers disagree, and a quoted head passes with a warning. A Bench verb
then runs against the stale binary, which story 23 forbids.

The Coverage axis measured nine forms and found eight disagreements:
`"bench" gate`, `'bench' gate`, `\bench gate`, `( bench gate )`,
`$(bench gate)`, `be"nch" gate`, `ls; "bench" gate`, and
`env "X=1" bench help`.

The reviewer decided the route on 2026-09-04. Teach the shell word test to
remove a surrounding quote and a leading backslash from the head word. Then
add one shared-table row per quoting form. Keep every row
resolver-independent, which is what the existing rows obey.

A form the shell test genuinely cannot reach is a Won't-handle cut, not a
forced row. Report any such form rather than making the table agree by
weakening the Go side.

## Acceptance

- [ ] The shared table holds one row per measured quoting form.
- [ ] The shell word test and `benchguard.InvokesBench` agree on every row.
- [ ] A stale core refuses `'bench' gate` at exit 2 with the rebuild sentence.
- [ ] Self-probe: remove the unquoting, and report which rows red on which side.
