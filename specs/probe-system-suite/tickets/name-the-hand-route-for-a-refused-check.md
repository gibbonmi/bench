# Name the hand route for a refused probe check

Blocked by: none
Writes: internal/testreport/selection_facts.go, internal/testreport/selection_facts_test.go, internal/probe/probe.go, internal/probe/refusal_test.go
Covers: none

## What to build

`bench probe` refuses `--check prose` and `--check system`. The approved `bench-probe` spec
keeps `bench test` as the caller for those two checks after a hand mutation. FT168 owns a
probe that runs them. The refusal and the help do not name the hand route, so a caller must
find the route alone.

The focused-run owner holds one route string. The probe refusal adds the route after its
hint. The probe help adds the route as a fourth note. Both surfaces read the one string.

## Acceptance

- [ ] `bench probe <file> --omit <old> --check system` prints the refusal and the hand route, and exits 1.
- [ ] `bench probe <file> --omit <old> --check prose` prints the same line, and exits 1.
- [ ] `bench probe --help` prints the hand route as its fourth note.
- [ ] The refusal and the help read the route from one string in `internal/testreport`.
- [ ] The comments on the refusal state the decided reason, not a claim about the report.
