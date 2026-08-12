# Repair structured envelope proof

Blocked by: none
Writes: `cmd/bench/command_registry_test.go`, `internal/axi/axitest/`

## What to build

Replace the marker-substring envelope assertions (review finding R5, High) with
structured validation: a new whole-document decode helper in
`internal/axi/axitest` (the existing `RecoverHelpCommandArgv` slices from the
first `help[` and hard-requires one row, so it cannot serve this) decodes the
COMPLETE stdout with the official `toon-go` decoder and reports the document's
blocks in order. The registry envelope test consumes it for every member's
success and empty cases. No production byte changes.

## Acceptance

- [ ] [SE1] (covers QD3) every approved member's success and empty stdout
  decodes in full as structured TOON — nothing before, between, or after the
  primary table and the help block — replacing the `strings.Contains` marker
  checks at the success, empty, and envelope assertions.
- [ ] [SE2] (covers QD3) the `help[N]{cmd,why}` block is required to be
  schema-correct and in TERMINAL position; a fixture-level red is demonstrated
  for trailing non-TOON bytes after the help block and for a malformed help
  table, then removed.
- [ ] [SE3] (covers local) the helper lives at `internal/axi/axitest` and
  changes no existing helper's interface; all current consumers stay green.
