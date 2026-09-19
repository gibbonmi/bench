# 7. Name the running executable in the unknown-check refusal

Blocked by: 1-split-named-check-owner.md
Writes: internal/testreport/
Covers: TP24, TP25, TP26, TP54

## What to build

Chunk: TP-C2.

The unknown-check refusal prints `unknown check: <name>`, then `executable: <path>`, then `seal: <value>`, then the check list, at exit 2.
One package variable supplies the absolute path of the running executable. `freshness.SealDigests` supplies the source digest.

The `seal` value is `unsealed` when the seal is unreadable. A control character in the path or in the caller's check name prints escaped through `sanitize.Controls`.
Move `TestUnknownNamedCheckReportsOperandAndInventory` to a new test file and rewrite its whole-output expectation.

Use the glossary term **running executable** in each comment and message.

## Acceptance

- [ ] The whole refusal output equals the four parts in that order.
- [ ] A seal file beside the executable gives `seal:` with its `sources` value.
- [ ] No seal file gives `seal: unsealed` at exit 2.
- [ ] A check name that holds U+0001 prints escaped, with no raw U+0001 byte.
