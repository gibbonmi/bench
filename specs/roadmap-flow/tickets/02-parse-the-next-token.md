# Parse the Next token in the board grammar

Blocked by: none
Writes: internal/roadmap/tree.go, internal/roadmap/tree_validation.go, internal/roadmap/tree_test.go, internal/roadmap/tree_helpers_test.go, internal/roadmap/roadmap.go, internal/conformance/roadmap_detail_integrity_test.go, tests/canary/roadmap-detail-integrity/

## What to build

`roadmap.ParseDocument` reads one `Next:` line per detail file, anchored at
column zero and outside a fenced code block, and validates its value against
the exported ordered token set `shape`, `spec`, `ticket`, `decide`, `kit-edit`.
Four fault classes arrive: missing line, unknown token, unanchored line, and
duplicate line. A row under a `## ` heading that contains `Parked` is exempt.
The unknown-token, unanchored-line, and duplicate-line classes reach the gate
through `roadmap-detail-integrity` now, each with one canary fixture. The
missing-line class exists in the parser but reaches the gate only in ticket 05,
so the live board stays green. Spec group B, rows RF11, RF12, RF13, RF15, RF17, RF18,
RF19, RF30, RF31, RF32.

## Acceptance

- [ ] A detail file with no `Next:` line yields one diagnostic that names the file path and the missing line.
- [ ] `Next: refactor` yields a diagnostic naming the rejected value; each of the five tokens yields no diagnostic.
- [ ] A row under `## Parked and scheduled work` with no `Next:` line yields no diagnostic; the same row under the features section yields one.
- [ ] A `Next: shape` inside a fenced code block yields the missing-line diagnostic.
- [ ] A line indented by one ASCII space, a line that starts with U+00A0, and a token wrapped onto the next line each yield the unanchored-line diagnostic naming the line.
- [ ] Two `Next:` lines yield the duplicate diagnostic naming the second line.
- [ ] A final `Next: spec` line with no trailing newline yields no diagnostic.
- [ ] The canary fixtures for the three live classes make the owner check red, and the fixture inventory records each class.
