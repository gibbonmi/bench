# Accept the source-tip pair in bench diff

Blocked by: none
Writes: internal/diff/diff.go, internal/diff/range.go, internal/diff tests

## What to build

`bench diff --base <commit> --source-tip <commit>` renders the immutable
base-to-tip review view: revision, aggregate, files, and with `--full` the log
and the patch body. This is the same explicit pair `bench preflight review`
accepts, so a review of a frozen pair never trusts the checkout for the tip.
`--source-tip` without `--base` is a grammar error. `--source-tip` with
`--commit` is a grammar error. The rendered view omits live-checkout facts,
as the `--commit` view already does.

## Acceptance

- [ ] `bench diff --base <b> --source-tip <t>` exits 0 and reports the base-to-tip files and aggregate.
- [ ] The output does not change when the checkout moves after the tip.
- [ ] `--source-tip` alone, and `--source-tip` with `--commit`, each exit 2 with the usage line.
- [ ] An unresolvable tip, and a base that is not an ancestor of the tip, each report a structured error.
