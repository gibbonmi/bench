# Project occurrence ledgers and migrate recurrence

Blocked by: none

## What to build

Roadmap context parses the shared occurrence grammar and each row's one canonical
ledger, derives recurrence facts from accepted keys, renders the schema-3 trust and
row fields, and migrates the nine legacy heading claims to their approved key sets.

## Acceptance

- [x] A missing ledger yields count zero, while one sorted unique ledger yields its
  exact key string and derived count.
- [x] Empty, malformed, duplicate, unsorted, or repeated ledgers emit
  `malformed-ledger` and make `sequence_trusted` false.
- [x] The nine approved rows migrate to their enumerated counts and no heading or
  second field stores a numeric recurrence tally.
- [x] Incident keys enforce the 1–64 ASCII-byte grammar, including CRLF and
  missing-final-newline documents, without collapsing absent and empty source states.
