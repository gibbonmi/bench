# Capture owned ideas

Blocked by: Project occurrence ledgers and migrate recurrence

## What to build

`bench idea` accepts paired owner and incident metadata through the shared argument
grammar, validates the owner against the trusted current roadmap, and appends the
canonical token in the same write as the ordinary idea text.

## Acceptance

- [ ] Valid metadata appends one final canonical token while preserving the existing
  parked stdout and the no-flag form byte-for-byte.
- [ ] Missing, repeated, empty, or malformed metadata exits 2 without mutation;
  unknown, retired, or structurally unverifiable owners exit 1 without mutation.
- [ ] Multi-word and leading-dash text, hostile repository paths, nested and linked
  CLI invocation, glob characters, and newline-less inboxes preserve the complete
  text and write only the repository-root inbox.
- [ ] Incident boundary and control-byte cases use the shared occurrence grammar,
  and every refusal leaves `IDEAS.md` byte-identical.
