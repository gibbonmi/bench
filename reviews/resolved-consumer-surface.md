# Review pickup: resolved-consumer-surface

Frozen pair: base `976b26e69b8f8bcc1f850570d0667f205a780e20`, tip `108dde57f35efdf96910f36a39fc8f85214ba257`. Reviewer: `codex exec gpt-5.6-sol high`, one axis per delegate. Raw findings: Standards 8, Spec 13, Coverage 7. Repair targets after collapse: 10 auto-fix, 6 ask-user.

## Standards

Count 8. Worst: the deleted-declaration match keys on directory plus spelling and bypasses `declKey`.

- `internal/consumers/blast.go` `blastAggregate` and `blastFileCount` duplicate `aggregate` and `countFiles` in `command.go`; rule: one source per fact. `auto-fix` (R9).
- `internal/outline/outline.go` `baseName` re-implements `path.Base`; smell: Lazy Element. `auto-fix` (R9).
- `internal/consumers/consumers.go` `relPath` comment describes an out-of-root fallback that does not exist; rule: a comment describes the current code. `auto-fix` (R9).
- `internal/consumers/resolve.go` `origin` comment states the alias and generic claim twice and argues its own correctness. The `internal/consumers/core_test.go` generic test comment claims a mutation that is not observable. Rule: `craft-comments`. `auto-fix` (R9).
- `internal/outline/outline_kinds_test.go` `oldOutline` comment narrates history ("before the kind vocabulary grew"). `auto-fix` (R9).
- `internal/consumers/blast.go` `topLevelDecls` re-derives the declaration-name forms `declName` owns, and `lookupDecl` re-derives `Resolve`'s grammar. `auto-fix` (R10): route the name through `declName`.
- `internal/consumers/blast.go` `tipDeclNames` keys directory plus spelling, not `declKey`. `ask-user` (A1).
- `.bench/BENCH-reference.md` `consumers` paragraph restates help and the review doc. `no-op`: the ticket `wire-review-blast-step` asks for the note.

## Spec

Count 13. Worst: blast is not a pure frozen-pair derivation. A dirty checkout at the tip changes rows, and a deletion-only edit inside a surviving declaration yields no rows.

- "enumerate consumers of every touched symbol": `parseSpan` gives a zero-count tip span an empty range, so a deletion-only hunk never touches its declaration. `auto-fix` (R1).
- "Blast rows derived only from the frozen pair": with tip at HEAD and a dirty checkout, the loader reads working-tree bytes. `auto-fix` (R3): refuse a dirty checkout in `--changed` mode.
- "an identical base and tip emit the definitive empty blast table": the command loads packages before it looks at the path set. `auto-fix` (R4): an empty changed-Go set skips the load.
- "a symbol in a non-Go file refused with the language named": a repository the Go loader cannot load refuses before the sweep. `auto-fix` (R5): on a no-module load failure, run the sweep first.
- "one kind `fixture` row per file under a `testdata/` segment": binary, oversized, and nonregular fixtures are skipped before classification. `auto-fix` (R6).
- "A control byte in a git-sourced path drops only its own row and sets `truncated=true`": a poisoned path reaches TOON and fails the whole response. `auto-fix` (R2).
- "A queried interface lists its satisfying types": an empty interface emits nothing. `ask-user` (A2).
- "each candidates row carries the exact qualified re-query spelling", for two modules that differ only in a dotted domain. `load(root, "./...")` loads one module, so the case is unreachable. `no-op`.
- Citation renders as `citation[1]{...}` where the spec writes `citation{...}`. `ask-user` (A3): spec amendment.
- A `--source-tip` that is not the checkout's HEAD refuses. `ask-user` (A4): spec amendment.
- Four real-loader tests plus the AXI envelope cases where the spec names one subprocess site. `ask-user` (A5): spec amendment to "one loader seam".
- The loader drops files outside the root; the spec is silent. `ask-user` (A5): spec amendment.
- `go.sum` was not new. `no-op`: spec prose only.
- Help Order 23 ties with `outline`. `no-op`: the sort is stable.

## Coverage

Count 7. Worst: the deletion-only edit above (collapsed into R1).

- Git C-quotes a patch header path with a control byte; `diffPath` does not unquote it, so the declaration is missed and `truncated` stays false. `auto-fix` (R2).
- Citation `cmd` joins argv with plain spaces; a ref such as `feature;echo` loses its token boundary. `auto-fix` (R8): render through the existing shell-quoting owner in `internal/axi`.
- The non-Go sweep reads with `os.ReadFile` and no `Lstat`; a live symlink or a FIFO is followed or blocks. `auto-fix` (R7): use outline's read policy.
- A generic type whose method set satisfies the queried interface for every `T` emits no `implements` row. `ask-user` (A6).
- A queried declaration under a tracked `vendor/` refuses because `./...` excludes vendored packages. `ask-user` (A6).
- The candidates spelling for a dotted path segment: collapsed into the Spec no-op above.
- Blast walk: `bench.commandRegistry`, `outline.Command`, and `outline.Symbols` consumers outside the diff examined; no unlisted-consumer miss.

## Ask-user set

- A1: `tipDeclNames` directory key. Recommendation: keep; the base side is untyped, and a package is a directory there; reword the `declKey` comment so it does not claim every comparison.
- A2: empty interface. Recommendation: keep the exclusion and state it in help; every named type satisfies `any`, and that flood serves no review.
- A3: `citation[1]{...}`. Recommendation: accept; amend the spec line to the one-row table.
- A4: non-HEAD tip refusal. Recommendation: accept; the review phase runs at the retained source's HEAD; amend the spec and note the refusal in help.
- A5: loader seam count and the out-of-root drop. Recommendation: accept; amend the testing decision to "one loader seam, tested through its real path" and add the out-of-root line to the implementation decisions.
- A6: generic implementers and `vendor/`. Recommendation: Won't handle lines in the spec; state both limits in help.
