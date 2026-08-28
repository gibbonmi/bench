# Extend outline with the candidate inventory kinds

Blocked by: none
Writes: internal/outline/

## What to build

The outline kind vocabulary gains `helper`, `double`, and `fixture`. Helper
and double rows are pattern-table entries with an optional path predicate.
The helper forms are `_test.go` functions named with the `new`, `make`, or
`with` prefix before an upper-case letter. The double forms are names with
a case-insensitive `fake`, `stub`, `mock`, or `spy` prefix. Fixture rows
come from a walk-level path classifier, because `testdata/` files carry no
scanned extension. The byte delta lands with an old-to-new fixture pair in
the outline package tests, and the LOCATE promise line stays verbatim.

## Acceptance

- [ ] OI1: a `_test.go` function named with a declared helper prefix emits
      kind `helper`.
- [ ] OI2: a name with the fake, stub, mock, or spy prefix emits kind
      `double`.
- [ ] OI3: a file under a `testdata/` segment emits one fixture row carrying
      line 1 and its base name.
- [ ] OI4: the old-to-new fixture pair reds on any unreviewed outline byte
      delta.
- [ ] OI5: outline help keeps the LOCATE promise line verbatim.
