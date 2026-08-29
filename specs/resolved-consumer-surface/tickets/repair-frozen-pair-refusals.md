# Grade the frozen-pair refusals of the blast mode

Blocked by: add-blast-mode.md
Writes: internal/consumers/

## What to build

The blast mode positions its rows in the checkout, so the checkout must
hold the exact pair the answer names. Two conditions break that promise.
Each one refuses on stdout at exit 1 with a remedy, and emits no result
block and no citation row.

A dirty checkout refuses. Its rows would come from working-tree bytes that
no commit froze, and a review recomputation cannot byte-match them. A
`--source-tip` that is not the checkout's HEAD also refuses, and that
refusal names both commits. Enumeration at the wrong tip grades a pair the
agent did not ask for.

## Acceptance

- [ ] BL8: a dirty checkout in `--changed` mode emits a structured stdout
      refusal at exit 1 naming the remedy.
- [ ] BL9: a `--source-tip` that is not the checkout's HEAD emits a
      structured stdout refusal at exit 1 naming both commits.
