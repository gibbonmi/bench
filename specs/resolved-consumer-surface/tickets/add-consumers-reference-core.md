# Add the consumers reference core

Blocked by: none
Writes: internal/consumers/ (new), go.mod, go.sum (new)

## What to build

The `internal/consumers` analysis core and its loader seam. The dependency
`golang.org/x/tools` lands here. The core consumes typed packages and
returns reference rows keyed to origin objects, sorted by file, line, then
column. Each row names its innermost enclosing named declaration.

Core tests type-check fixture source in process with `go/parser` and
`go/types`. One focused loader test drives the real `go list` path against
a minimal fixture module; it is the single subprocess site. The ticket's
tests render rows through `internal/toon`, so CS13's byte comparison runs
here.

## Acceptance

- [ ] CS1: a qualified query over the typed fixture emits every planted
      reference row.
- [ ] CS4: each row names its innermost enclosing named declaration.
- [ ] CS13: an alias-spelled query and an origin-spelled query emit
      byte-identical tables.
