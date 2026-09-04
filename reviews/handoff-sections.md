# Review pickup: handoff-sections

Frozen pair: base `6fcd0882cff0315b90de6878fb5ab868aec7de19`, tip `51d17f52cb55a71695a117af8d8ca7e92d21af9a`.
Raw findings: 14 (Standards 6, Spec 4, Coverage 4). Repair targets after collapse: 6.

## Standards

Count: 6. Worst: the Shape text still states the retired file-age rule.

- S1 `auto-fix` — `internal/handoff/text.go:35-38` keeps "computes its age from the file's last write" while `AGENTS.md:78-81` and the spec date each section by its branch. Rule: one source per fact. Same fix as Spec 2.
- S2 `auto-fix` — `internal/adopt/init.go:76-80` spells `capture/session-handoff.md` as a literal beside the new `handoffdoc.DocumentPath`. Rule: one source per fact.
- S3 `ask-user` — `internal/handoffdoc/store.go:12-19` mirrors the intent ledger's 2s and 10ms deadline literals from `internal/intent/intent.go:240,270`. The collapse needs an export from `internal/intent`, outside the fence. Reported, not repaired.
- S4 `auto-fix` — `handoffdoc.WriteSection` (`store.go:71`) has no production caller; only tests call it. Speculative Generality. Non-blocking; reported, not repaired.
- S5 `no-op` — `internal/handoff/sections_test.go:73-79` hand-rolls a repo the `gittest.RepoOnBranch` helper owns. Non-blocking.
- S6 `auto-fix` — `internal/adopt/adopt_test.go:515-516` comment says all Init tests use the preamble; three do not. Non-blocking; reported, not repaired.

## Spec

Count: 4. Worst: the status row's unresolved-section advisories are unasked-for and untested.

- Sp1 `ask-user` — `internal/status/handoff.go:148-161,167-174,214-238` emits "N sections unresolved" and per-section reason rows. Spec line 109 asks only to name the behind section. Decision under the batch approval: cut the rows; a dead section is residue the retirement removes. Recorded in the spec for veto.
- Sp2 `auto-fix` — HS5 partial: the Shape's age sentence describes the old rule. Same fix as S1.
- Sp3 `ask-user` — the legacy migration (`internal/handoffdoc/document.go:285-330`) has no spec sentence or row, and its ticket claims `Covers: HS7` falsely. Decision: add story 27 and row HS26; the ticket cites HS26.
- Sp4 `no-op` — the safe-unlink lock reclaim exceeds spec line 102. It is ticket-authorized and tested. Recorded as a flagged addition in Further notes.

## Coverage

Count: 4. Worst: an unterminated fence in one section's State swallows every later section.

- C1 `ask-user` — `internal/handoffdoc/document.go:251-257` carries fence state across sections. Probed: an open fence in section A's State absorbs section B and Shape, and B's writer re-appends a duplicate. Decision: `bench handoff` refuses a State with an unterminated fence and prints the line. `Parse` refuses an unterminated fence at end of file with the file and line. New rows HS27 and HS28.
- C2 `ask-user` — an unfenced `## ` line in State bricks every later verb (`h.md:20: heading "## Open questions" is not ...`). Decision: `bench handoff` refuses a State line that opens a level-two heading outside a fence and prints the line. New row HS29.
- C3 `auto-fix` — `internal/worktree/lifecycle.go:451` discards the removal error, so a broken document or a held lock leaves the section pinned in silence. Decision: the retirement prints the removal error and keeps its verdict. New row HS30.
- C4 `auto-fix` — `internal/handoff/state_scan.go:41` scans fenced lines, and an ambiguous 7-hex prefix reads as prose because both `cat-file` probes fail. Decision: the scan skips fenced lines, and an ambiguous abbreviation refuses with its own reason. New rows HS31 and HS32.
