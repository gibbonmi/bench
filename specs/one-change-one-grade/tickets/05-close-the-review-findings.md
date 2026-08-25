# Close the review findings

Line: sonnet / low.

Blocked by: 04-state-the-lane-in-the-rules.md
Writes: internal/commit/lane_test.go, internal/gate/lane_record_test.go, specs/one-change-one-grade/spec.md, reviews/one-change-one-grade.md (delete)

## What to build

Two decided edges gain a test, and two coverage rows say what their tests read.
A lane fail under `--dry-run` exits 1, prints `lane{outcome=fail,check=<name>}`,
and publishes nothing. A lane whose check outlives the gate timeout returns an
error and writes no lane record. Rows OG14 and OG23 name the gate script's tally
instead of the literal `gate: green`, because the fixture's script route prints
no Bench-owned verdict line. The green fix commit deletes the review pickup file.

## Acceptance

- [ ] A lane fail under `--dry-run` exits 1, names the check, and leaves the branch ref unchanged.
- [ ] A lane that exceeds `gateTimeout` returns an error and leaves no lane file in the Git dir.
- [ ] OG14 and OG23 read as the tests assert them.
