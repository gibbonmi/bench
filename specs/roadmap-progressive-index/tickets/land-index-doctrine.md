# Land the schema-4 index and the index-first doctrine atomically

Blocked by: add-context-envelope.md, add-row-selector.md
Writes: internal/roadmap, .agents/commands/bench-what-next.md, .claude/commands/bench-what-next.md, internal/conformance/recurrence_maintenance_contract_test.go, internal/conformance/docs_workflow_helpers_test.go, internal/anchors/registry_data.go, CONTEXT.md

## What to build

The atomic sequence the spec's Doctrine-anchors decision pins — none of
these lands green alone. The snapshot advertises schema 4: bodies
(`body`/`text`/`raw`) empty in index mode with true `*_bytes`, complete in
`--row` and `--full`; the `truncated` column leaves every block; `ideas` and
`learnings` gain their parsed `line` fields in every mode; every mode emits
the same sixteen-blocks-plus-help list. In the same diff,
`/bench-what-next` moves to index-first evidence (index once, `--row`
fetches, capture bodies read from the paths the index names), the recurrence
contract's schema and evidence anchors move to 4/index-first, the helpers'
occurrence expectation moves to the enumerated two-spelling invocation set
while its other two needles survive, and the retro-evidence anchor narrows
to the inventory rule with its `Diagnostic` byte-identical (the
`implementation-retro-drain-anchor` canary matches that literal). Every
oracle edit keeps its fail posture with its bite updated, never weakened.
`CONTEXT.md` gains the index/detail glossary entry.

## Acceptance

- [ ] Default mode: empty bodies with true byte counts across roadmap rows,
      ideas, learnings, retros; `line` fields present; schema 4; no
      `truncated` column; failure `raw` follows the same rule
      (covers PI1, PI2, PI3, PI4, PI5).
- [ ] `--full` emits every body complete at schema 4 (covers PI9); every
      mode enumerates the complete block list (covers PI18).
- [ ] The rewritten phase prose and all three pinning sites are green
      together, and reverting the prose alone or the schema alone turns the
      relevant check red — both reverts performed and their red output
      recorded in the hand-back, since this row cannot be verified from the
      final green tree (covers PI17).
