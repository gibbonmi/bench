# Apply recurrence during maintenance

Blocked by: Project capture occurrences

## What to build

The reviewed maintenance phase consumes schema 3, refuses untrusted occurrence
evidence, records pending keys before removing their source units, and applies
recurrence only after every stronger prioritization input remains tied.

## Acceptance

- [x] The phase requires schema 3 and stops before batch mutation when
  `sequence_trusted` is false while leaving the complete context visible.
- [x] Every pending key is added before its source is removed, while an
  `already-recorded` source is removed without adding another key.
- [x] Severity, actionability, dependencies, and explicit pricing precede descending
  occurrence count; existing defect and cost rules apply only after that tie remains.
- [x] Conformance proves each required clause bites with its own targeted diagnostic.
