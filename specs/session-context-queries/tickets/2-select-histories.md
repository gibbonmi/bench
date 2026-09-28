# 2. Select bounded spec histories

Blocked by: none
Writes: internal/spec/history.go, internal/spec/history_test.go, internal/spec/history_command_test.go (new), internal/spec/history_selected.go (new), internal/spec/history_selected_test.go (new), internal/spec/testdata/selected-defaults.json (new), internal/spec/spec.go, internal/spec/spec_test.go, cmd/bench/main.go, cmd/bench/selected_queries_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand, CHANGELOG.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: QU4, QU5, QU6, QU7, QU8, QU19, QU20, QU21, QU22, QU23, QU24, QU25, QU27

## What to build

Add repeated spec selection with an explicit per-spec event limit.
Consume the existing history producer and render summaries plus selected events.
Keep the shared producer complete for the roadmap context reader.
Do not apply selection or limits inside `History`.
Preserve positional complete history and the current root command disposition.
The selected query retains all per-target outcomes.

Start the selected route only when `--spec` or `--limit` is present, before `specArg` parses the positional operand.
Keep the current expectations in `internal/spec/spec_test.go` unchanged; add a case only.

Port from the stream reference source, and verify each ported line against the current tree:

- `194b7dba:internal/spec/history_selected.go`
- `194b7dba:internal/spec/history_selected_test.go`
- `194b7dba:internal/spec/history_command_test.go`
- `194b7dba:internal/spec/testdata/selected-defaults.json`
- `194b7dba:internal/spec/history.go`
- `194b7dba:cmd/bench/main.go`
- `194b7dba:cmd/bench/help_inventory_test.go`

The stream calls `slugOf`, but `main` has `SlugOf` (`internal/spec/spec.go:496`).
`Command` returns three values on `main` (`internal/spec/spec.go:201`).
`main.go` now routes help through `Kind: commandHelp` (`cmd/bench/main.go:111`), so do not port the stream's `helpCommand` move.

## Acceptance

- [ ] Each spec retains exact matching, commit order, and retire/delete classification.
- [ ] Each history respects the explicit event limit and reports complete counts.
- [ ] Each omitted history names its exact complete-detail command, including for `x.md.md` and `.md`.
- [ ] One failed spec cannot erase successful or empty spec results.
- [ ] Existing positional invocations match the baseline input matrix.
- [ ] The shared history producer returns the baseline complete sequence after selected queries run.
- [ ] The serialized byte count is the byte length of the complete positional history table before projection.
- [ ] One unrepresentable commit subject cannot erase other selected spec results.
- [ ] Two 3-event specs at limit 2 print exactly 8 lines with no spill line through the bounded dispatcher.

## Headroom

The new selected-history file owns selection and projection after the complete producer returns.
The existing history entry retains its positional path.
The selected-history tests use a separate file to retain structural headroom.

## Verification supplement

The command fixture records positional baselines before production changes.
It normalizes only fixture commit identities and temporary repository paths.
Expected tables use the canonical TOON owner.
Named mutations independently prove baseline preservation, omission counts, bytes, and producer completeness.

The producer matrix includes empty, retire-only, delete-only, mixed, and duplicate-commit histories.
Exact slug matching and deterministic commit dates expose prefix and ordering defects.
The operand matrix includes aliases, repeated requests, spaces, quotes, Unicode, leading dashes, and shell-shaped text.
First, later, and repeated unsafe requests retain stable ordinals.

A real Git boundary wrapper fails one selected history while other targets use the real executable.
That wrapper binds `PATH` in the spec package, so the worktree parallel census stays outside the fence.
Representable subject controls remain valid; refused subjects outside the selected prefix still fail their own history.
