# AXI query disclosure implementation re-review

Candidate: `c2542a97..543d8551`

Raw findings: Standards 3, Spec 2, Coverage 3. Accepted repair targets: 6.

## Standards

Count: 3. Worst: worktree orphan state is derived twice and paired to paths by index.

1. **Low — auto-fix.** `internal/worktree/list.go` independently classifies orphan rows while building `orphanPaths` and again in `actionsForRows`; carry the clean path with the row or derive the action at the owning seam.
2. **Low — auto-fix.** `internal/coverage/coverage_test.go`'s `wantTable` reconstructs the production action derivation; replace mirrored derivation with independent literal expectations.
3. **Low — auto-fix.** `capture/session-handoff.md` still pins the first review before any repair; rewrite it at this review boundary.

## Spec

Count: 2. Worst: `worktree list` changes accepted and rejected argv outside QD2's named bare-`help` delta.

1. **Medium — auto-fix.** `internal/worktree/list.go` delegates multi-token and hostile-token cases to `usage.Parse`, changing `--help extra`, `-h extra`, `-- x`, and the empty token from their pre-change usage/2 bytes. Restore exact parity outside the three single-token help spellings and make the checked-in argv fixture cover each equivalence class.
2. **Low — auto-fix.** `.agents/skills/bench-craft-cli/SKILL.md` says coverage emits a repair retry "after refusal"; the retry is appended by the successful extraction query for a repairable map. Describe the shipped transition precisely.

No-op: alternate modes such as hook-driven `guards --brief` are not agent discovery turns in the approved command inventory; the surface tickets scope disclosure to default extraction, with `coverage --check` itself the recommended check action. Do not widen those operational modes.

## Coverage

Count: 3. Worst: TOON-permitted control-bearing paths can turn a previously successful query into an action-rendering refusal.

1. **Medium — auto-fix.** `internal/axi/action.go` rejects all control runes although QD6 and the hostile-environment clause require control-bearing known values to survive through the existing escaping rules. Permit the TOON-supported control-byte partition in known arguments and prove exact public rendering without weakening refusals for unsupported controls.
2. **Medium — auto-fix.** The worktree argv fixture omits the multi-token, separator, and empty-token partitions named by the Spec finding; repair with the same target.
3. **Low — auto-fix.** `internal/learnings` lacks a checked-in pre-disclosure malformed/refusal response pair, leaving its primary bytes and exit unguarded.
