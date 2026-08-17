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

- [x] AU1: unknown landedness or assignment state retains as `uncertain`.
- [x] AU2: registration without the verified owner-and-assignment join retains as `foreign`.
- [x] AU3: malformed assignment or mismatched recovery metadata retains as `malformed`.
- [x] AU4: an unexpected lock surviving explicit planning retains as `unexpected-lock`.
- [x] AU5: a live lease wins as `live-lease` even with another explicit refusal.
- [x] AU6: an aged unlanded assignment with ignored residue retains as `ignored`.
- [x] AU7: a young active unlanded assignment retains as `active`.
- [x] AU8: an active landed assignment retains as `landed` under the current override.
- [x] AU9: an aged active unlanded assignment retains as `orphaned`.
- [x] AU10: a cleanup-pending assignment with an unlanded branch retains as `unmerged`.
- [x] AU11: a cleanup-pending landed assignment needing preservation retains as `dirty`.
- [x] AU12: a verified cleanup-pending landed clean assignment projects `remove` with no reason code.
- [x] AU13: that assignment with declared bounded ignored output projects `discard-remove`.
