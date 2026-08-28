# Add the consumers command surface and registration

Blocked by: add-consumers-reference-core.md
Writes: internal/consumers/, cmd/bench/, bin/bench.sh, internal/conformance/, projects/benchkit.md, .agents/skills/bench-craft-cli/SKILL.md

## What to build

`bench consumers <qualified-symbol>` works end to end. The response is the
`consumers[N]{file,line,via,enclosing}` table, the meta table, and the
`help[0]` envelope. The default is complete up to the 200-row cap. Over the
cap, the default emits per-package aggregates, `truncated=true`, and one
`--full` action carrying the symbol; `--full` emits every row. The help
text carries the promise clause and the soundness clause.

Registration lands in all four graded places in this one commit. The places
are `commandRegistry` with the AXI disposition, `approvedAXIQueries`,
`subcommandRouting`, and both approved-query tables (the `craft-cli` skill
and the project profile). The wrapper gains a `consumers` porcelain route
beside the other read verbs; nothing grades that line, so no row exists.

## Acceptance

- [ ] CS6: a matched symbol with zero references emits the definitive empty
      consumers table.
- [ ] CS14: a symbol with 3 planted references emits 3 rows and `rows=3` in
      meta.
- [ ] CS15: an unknown flag exits 2 with usage on stdout.
- [ ] CS16: the help text carries the soundness clause verbatim.
- [ ] CS19: a terminal symbol result ends with the empty `help[0]` envelope.
- [ ] CS20: a symbol with more references than the cap emits per-package
      aggregates, `truncated=true`, and one `--full` action.
- [ ] CS21: `--full` emits every planted row past the cap.
- [ ] CS23: `--help`, `-h`, and bare `help` each print usage on stdout at
      exit 0.
- [ ] CS24: the help text carries the identifies-edges promise clause
      verbatim.
- [ ] AD1: `bench help` lists `consumers` with its one-line promise.
- [ ] AD3: `axi-query-registry` accepts the `consumers` disposition.
