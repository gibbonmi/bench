# Emit the replayable citation row

Blocked by: add-consumers-command-surface.md
Writes: internal/consumers/

## What to build

Every success response carries one `citation{sha,state,version,cmd,hash}`
row immediately before the help envelope, because the AXI contract pins
`help` as the terminal block. The hash is sha256 over all output bytes
before the citation row. The `version` value is the bench version, and
`state` reports `clean` or `dirty` from the checkout. The determinism probe
runs the command twice at one SHA and compares bytes.

## Acceptance

- [ ] CS10: every success response carries one
      `citation{sha,state,version,cmd,hash}` row before the terminal help
      envelope.
- [ ] CS11: a dirty checkout emits `state=dirty` in the citation row.
- [ ] CS12: two identical runs at one SHA produce byte-equal output.
- [ ] CS22: a recomputed sha256 over the bytes before the citation row
      equals the printed hash value.
