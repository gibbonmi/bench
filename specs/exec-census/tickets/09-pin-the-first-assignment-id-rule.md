# Pin the first assignment id rule

Blocked by: 03-name-the-real-verb-head.md
Writes: internal/census/census_test.go, specs/exec-census/spec.md

## What to build

Review finding, Coverage axis. A text that names two assignment ids records
one line under the first id and none under the second. No row or edge line
decides that. Story 2 counts one call as one record, so the first id
owns the record. This ticket pins that rule; the code does not change.

Add one line to the spec's edge inventory: a text that names two assignment
ids records once, under the first id in the text. Add one census unit test
that proves it.

## Acceptance

- [ ] A text `cat <pool>/<o>-<a>/x <pool>/<o>-<b>/y` appends one record under `<a>` and none under `<b>`. (EC02)
- [ ] The edge inventory names the rule.
