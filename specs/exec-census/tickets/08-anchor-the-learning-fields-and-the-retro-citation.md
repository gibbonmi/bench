# Anchor the learning fields and the retro citation

Blocked by: 07-add-the-census-duty-and-the-charge-line.md
Writes: .agents/commands/bench-final-check.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go

## What to build

Review finding, Spec axis. Row EC27 names the four learning fields and the
retro citation under the anchors registry test. The registry anchors only
the duty sentence and the zero close. An edit can drop the field shape or the
retro citation silently.

Reflow the prose so each anchored sentence sits on one physical line, and
register two more `Require` needles in the `AfterImplementSpec` group beside
the existing census rows:

- `Its --what lists each verb head with its count. Its --right names the Bench form per head, or none. Its --rule proposes the verb or the help change.` with `--what`, `--right`, `none`, and `--rule` as code spans, on one line.
- `A spec retro cites the landing's census entry under ### Bench CLI with its Feeds: line.` with `### Bench CLI` and `Feeds:` as code spans, on one line.

Each needle carries its own diagnostic. The removal test in
`registry_data_test.go` gains one row per needle. The prose stays within the
file's paragraph and sentence rules, so the reflow keeps the paragraph at six
sentences or fewer.

## Acceptance

- [ ] A copy of the final-check command without the fields sentence makes the registry report its diagnostic. (EC27)
- [ ] A copy of the final-check command without the retro citation makes the registry report its diagnostic. (EC27)
- [ ] The live-tree conformance check stays green.
