# Retain all failure diagnostic lines

Blocked by: none
Writes: internal/testreport, CHANGELOG.md
Covers: none

## What to build

Retain each diagnostic line for failed tests and failed packages.
Make --full expose those lines in their original order within each result.
Keep runner-line filtering, control sanitization, and deterministic result order.
Keep the default preview contract.
Count each reported failing test identity once, regardless of its diagnostic line count.
Keep the existing suppression of silent parent failures.

Add a Fixed entry under a dedicated Full failure diagnostics changelog heading.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] Full output retains the first, middle, and last diagnostic lines.
- [ ] Full output retains multiple compiler diagnostic lines.
- [ ] One failed test with several lines counts as one failed test.
- [ ] Distinct failed tests retain their distinct count.
- [ ] Runner lines remain absent and control bytes remain sanitized.
- [ ] Default output retains its existing preview behavior.
