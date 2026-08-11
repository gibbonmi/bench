# Repair: fence-grammar paren depth, recursive tickets, unconditional ESC refusal

Blocked by: repair-compose-changed-files.md
Ownership fence: `internal/preflight/`
Integration surfaces: fence tokenizer + tickets enumeration + path rendering→`internal/preflight/gather.go`, `internal/preflight/command.go` (repaired in place, exercised by every RG row)
Contracts: none crosses — all three repairs are internal to the package's existing seams, asserted by the RG rows through the existing CLI harness
Closure: RG1/wrapped-paren-never-authorizes, RG1/entry-after-closed-paren-authorizes, RG2/subdir-row-owned, RG2/subdir-phantom-detected, RG3/fenced-esc-path-reds

## What to build

Three review findings against the spec's literal predicates.
P1 (blocking): the fence tokenizer resets paren depth at each line boundary,
so a parenthetical wrapped across lines authorizes its backticked tokens —
fails open in the anchored-grammar hostile class. Track depth across the whole
`## Ownership fences` section: a token inside an open paren is never an
authorization even on a continuation line, and a legitimate entry on a line
after the paren closed still authorizes.
P2: predicates 3 and 4 say "under `specs/<slug>/tickets/`", but enumeration
skips directories — recurse, keeping the lstat special-file refusal at every
level.
P3: PF7 says a changed path carrying ESC exits 1 with the
unrepresentable-TOON-cell error, but the check only fires when the path
reaches a red row's detail cell — a fenced ESC path exits 0. Refuse
unrepresentable changed paths unconditionally before verdict rendering, and
update the test whose comment records the narrowed reading.

## Acceptance

- [ ] [RG1] (covers local) a fences section whose parenthetical wraps across lines never authorizes its backticked tokens (the out-of-fence path is red), and an entry on a line after a closed paren authorizes normally — both asserted through the CLI harness.
- [ ] [RG2] (covers local) a declared row cited only in `tickets/sub/x.md` is owned, and a phantom own-tag token there is detected — both modes.
- [ ] [RG3] (covers local) a changed path carrying ESC exits 1 with the unrepresentable-TOON-cell error even when the path is inside a fence and every check would be green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RG1/wrapped-paren-never-authorizes | reset paren depth per line (revert to the shipped behavior) | the wrapped-parenthetical contract test | seed the two-line parenthetical with a tracked change under its token, run review, expect the false-green failure |
| RG1/entry-after-closed-paren-authorizes | never reset depth once opened | the entry-after-paren contract test | seed a closed paren then a real entry, run, expect the false-red failure |
| RG2/subdir-row-owned | skip directories during enumeration (shipped behavior) | the subdir-citation contract test | seed the only citation in `tickets/sub/x.md`, run review, expect the rows-owned false-red failure |
| RG2/subdir-phantom-detected | enumerate subdirs but drop their token scan | the subdir-phantom contract test | seed an own-tag phantom token (the spec's tag followed by 99) only in `tickets/sub/x.md`, run, expect the missed-red failure |
| RG3/fenced-esc-path-reds | check representability only in red-row details (shipped behavior) | the fenced-ESC contract test | seed an authorized ESC path on a green tree, run, expect the exit-0 failure |
