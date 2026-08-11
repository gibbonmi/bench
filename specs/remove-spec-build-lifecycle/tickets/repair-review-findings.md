# Repair the composed-review findings

Blocked by: repair-prose-and-advertisements.md

## What to build

Close the three-axis review's blocking findings S1 and T1 and the C1/C2
coverage gaps. Restore the two anchor rows deleted while their prose
survives: `projects/benchkit.md`'s "Both dogfood traces use the public
porcelain:" needle, and a re-needled row over the surviving
falsification-offer paragraph in `bench-implement-spec.md`. Delete the
stranded provisional/abandon chain in `internal/worktree/path.go`
(`PlanAbandon`, `ApplyAbandon`, `ReleaseProvisional`,
`planProvisionalRelease`, `validateProvisionalEvidence`,
`compactProvisionalAssignment`, `abandonReceipt`, `planLeftoverEntry`,
`planRemovedCheckout`) with its now-subjectless tests — verifying no
non-test caller per symbol before each deletion, and leaving the explicit
`clean --apply` preserve path untouched (its fate is an open reviewer
decision). Add one table-driven test asserting `spec build` and
`worktree recovery` answer their unknown-argument errors and that
`worktree --help` omits `recovery`. Add an anchor row pinning the
CHANGELOG removal entry.

## Acceptance

- [ ] The two anchor rows exist again and bite (needles match live prose).
- [ ] The listed provisional/abandon symbols and their orphaned tests are
      gone; every deleted symbol had no non-test caller at deletion time.
- [ ] The removed-grammar test exists and is green; re-adding `recovery` to
      the worktree help string turns it red (demonstrate once, restore).
- [ ] An anchor row requires the CHANGELOG removal entry.
- [ ] Full `go test ./...` green.
