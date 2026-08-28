# Add the blast mode over a frozen pair

Blocked by: emit-citation-row.md
Writes: internal/consumers/

## What to build

`bench consumers --changed --base <b> --source-tip <t>` enumerates the
consumers of every declaration the diff touched, at the tip. The derivation
is a pure function of the hunk list and the tip packages, and the git
queries stay at the rim through `internal/git`. Rows are
`blast[N]{changed_symbol,file,line,touched}`. The `touched` value is true
when the consumer file sits in the diff's changed-file set. A deleted
declaration emits one `blast_deleted[N]{changed_symbol,base_file,base_line}`
row and no consumer rows.

An identical base and tip emit the definitive empty table. The grammar
mirrors `bench test`: `--source-tip` without `--base` exits 2. The response
ends with one per-symbol `--full` help action for each changed symbol with
an untouched consumer, deduplicated in stable order, plus the citation row.

## Acceptance

- [ ] BL1: `--changed --full` over the fixture pair emits the consumers of
      each touched declaration.
- [ ] BL2: a declaration the diff deleted emits one `blast_deleted` row and
      no consumer rows.
- [ ] BL3: an identical base and tip emit the definitive empty blast table.
- [ ] BL4: two blast runs over one frozen pair are byte-equal.
- [ ] BL5: `--source-tip` without `--base` exits 2 with usage.
- [ ] BL6: a consumer row inside the diff's file set emits `touched=true`
      and one outside emits `touched=false`.
- [ ] BL7: a blast result offers one per-symbol `--full` help action for
      each symbol with an untouched consumer.
