# Characterize automatic eligibility outcomes

Blocked by: 01-characterize-explicit-eligibility-outcomes.md
Writes: internal/worktree/eligibility_test.go, specs/worktree-cleanup-eligibility/tickets/02-characterize-automatic-eligibility-outcomes.md

## What to build

Extend the independent matrix with the stricter automatic-cleanup reading at its
existing planning seam. Every real-fixture case asserts one current final tuple
and any competing evidence that establishes the present override order, including
that the explicit-only branch-deletion assertion cannot affect automatic cleanup.
Expected tuples remain authored in the test rather than generated from production
rules.

## Acceptance

- [ ] AU1: unknown landedness or assignment state retains as `uncertain`.
- [ ] AU2: registration without the verified owner-and-assignment join retains as `foreign`.
- [ ] AU3: malformed assignment or mismatched recovery metadata retains as `malformed`.
- [ ] AU4: an unexpected lock surviving explicit planning retains as `unexpected-lock`.
- [ ] AU5: a live lease wins as `live-lease` even with another explicit refusal.
- [ ] AU6: an aged unlanded assignment with ignored residue retains as `ignored`.
- [ ] AU7: a young active unlanded assignment retains as `active`.
- [ ] AU8: an active landed assignment retains as `landed` under the current override.
- [ ] AU9: an aged active unlanded assignment retains as `orphaned`.
- [ ] AU10: a cleanup-pending assignment with an unlanded branch retains as `unmerged`.
- [ ] AU11: a cleanup-pending landed assignment needing preservation retains as `dirty`.
- [ ] AU12: a verified cleanup-pending landed clean assignment projects `remove` with no reason code.
- [ ] AU13: that assignment with declared bounded ignored output projects `discard-remove`.
