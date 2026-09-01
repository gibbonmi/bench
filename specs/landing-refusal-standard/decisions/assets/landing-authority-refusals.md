# Refusal-face audit (FT169 research)

Source: read-only Terra delegate, 2026-09-01. Coordinator spot checks:
`internal/worktree/land_identity.go:41` and `:114` pass an empty `next`
to `checkoutClean`; `internal/worktree/merge.go:302` supplies one; no
production file reads `MERGE_HEAD` — all resolve.

Standard under audit: a refusal names the failed path, the recovery state,
and the exact next command.

## Per-occurrence grades

| # | Occurrence | Class | Site | Missing element |
|---|---|---|---|---|
| 1 | lost request token | fixed | `internal/worktree/worktree.go:594` | none — names `bench worktree reauthorize ...` |
| 2 | fence refusal without path | partial | `internal/worktree/land_identity.go:179` | paths ride inside the detail string via `fmt.Errorf`; no path table, no next |
| 3 | six-refusal chain | partial | `internal/worktree/land.go:230-258` | groups batch, but each group short-circuits at its first failure; 4 of 6 faces carry no next |
| 4 | stale seal + residue | partial | `internal/worktree/land_identity.go:48` | seal half fixed (`freshness_verify.go:79-80`); the residue half names paths but no declare-or-clean route |
| 5 | dirty destination, moved paths | partial | `internal/worktree/land_identity.go:41` | paths listed; `checkoutClean` gets `next=""` |
| 6 | backslash paste | fixed | `internal/worktree/land_refusal.go:99-104` | none — `ShellQuote` plus the exec pointer fallback |
| 7 | byte-identical IDEAS.md | fixed | `internal/landing/composition.go:57-68` | none — union rule removed the cause |
| 8 | exec routing variable | partial | `bin/bench.sh:398-401` | names the variable, not the exact command; does not say the landing runs outside `worktree exec` |
| 9 | conflict vs pending merge | open | `internal/worktree/land_refusal.go:61-67` | one repair line for both states; `git merge --continue` never named; no `MERGE_HEAD` reader exists |
| 10 | spec retire primary refusal | partial | `internal/usage/worktree.go:53-55` | next names `bench worktree create` only, not the follow-on retire route |
| 11 | `--spec` path form | open | `internal/worktree/land_identity.go:130,141` | no parse-seam normalization; `gather.go:243` matches the raw argument by equality |
| 12 | late `--base` refusal | partial | `internal/worktree/land.go:248-252` | preflight grades the base against the source tip only, so the review still passes on a base the landing refuses; no next |
| 13 | short-tip refusal form | partial | `internal/worktree/land_identity.go:105` | short hex now expands; the mismatch message states no accepted form and no next |

## Structural notes

- One formatter exists and supports the full standard:
  `refusal{detail, observed, wanted, next, paths}` in
  `internal/worktree/classifier.go:23-65`. The `next` field is optional,
  and 4 of 29 sites set it. The gap is structural, not accidental.
- Three formatter families disagree on the standard: `refusal{}` with
  optional next, `toon.Errorf(kind, hint)` with a mandatory hint, and bare
  `fmt.Errorf` or `echo`. Every refusal that meets the full standard
  comes from the mandatory-hint family or from one composed site.
- The strongest fixes removed a cause instead of wording a message
  (occurrences 1, 6, 7). The remaining partials all need the message to
  carry the route.
- A path detail loses its type at a layer boundary. In
  `land_identity.go:170-185`, the code folds the authorization error into
  a string, so the computed paths print inside the sentence, not the path
  table.
- Adjacent correctness fault, outside the 13: `spec.Resolve`
  (`internal/spec/spec.go:177`) reads an explicit path against the process
  working directory, not the passed base. The landing passes
  `base = a.Worktree`, so a relative path form reads the primary
  checkout's file.

## Read and not read

The delegate read the worktree land/refusal/classifier code, the landing
composition, preflight decision and gather, spec, usage, freshness,
runbinary, toon, and `bin/bench.sh:360-455`. It did not read the commit
package, ownership beyond line 463, the lifecycle code, or the tests
beyond two excerpts. It ran no tests and made no edits.
