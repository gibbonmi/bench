# Structure gate-prose findings

Blocked by: improve-gate-prose-help-and-findings.md
Writes: internal/prose/subject.go, internal/prose/prose.go, internal/prose/prose_test.go, internal/gate/gate_prose.go, internal/gate/gate_prose_test.go
Covers: none

## What to build

Expose the named prose-grade result as structured data and render the
gate-prose output from those fields. Keep diagnostic-string rendering owned by
the prose package instead of parsing that rendered protocol in the gate.

## Acceptance

- [ ] Gate-prose obtains path, line, rule, and sentence from structured prose
      results.
- [ ] The gate does not parse prose diagnostic strings.
- [ ] Existing public diagnostics, including line number and offending
      sentence, remain stable.
