# Review pickup — ft234-pool-key-reclaim

Base `537cbf51` · reviewed tip `f4e34897` · three read-only axis delegates, `opus`.

Raw findings: 12 (Standards 4, Spec 4, Coverage 4).
De-duplicated repair targets: 7. Standards 4 and Coverage C3 name one fix.

## Standards

Count 4. Worst: S1 — the duplicated fingerprint-shape predicate, the exact
code-standard defect this axis exists to grade.

- **S1 · auto-fix.** `internal/worktree/pool_reclaim.go:299-302`
  `wellFormedFingerprint` and the inline check at
  `internal/worktree/worktree.go:264-265` independently derive the same fact —
  lowercase hex of a sha256 digest — in one package. Cites AGENTS.md "one source
  per fact". Coordinator-verified: byte-equivalent logic. Third and fourth
  variants exist in `internal/freshness/freshness.go:701` and
  `internal/gate/verdict.go:530`, pre-existing and outside this diff.
- **S2 · auto-fix.** `pool_reclaim.go` carries 75 comment lines over 413 (~18%)
  against a package running 1–8% (`clean.go` 5/362, `list.go` 1/184,
  `resume.go` 37/465). `pool_reclaim.go:319-325` and `:31-33` argue correctness
  to a diff reviewer. Cites `craft-comments`, "The register".
- **S3 · auto-fix.** `pool_reclaim.go:347-353` and `:389-391` each build the
  same `{key, verdict, reason}` projection for one schema. The field list is
  single-sourced; the projection that must track it is not.
- **S4 · ask-user, same target as C3.** Partial apply exits 0.

## Spec

Count 4 (3 no-op, 1 ask-user). Worst: P4, the silently swallowed pool-read error.

- **P1 · no-op.** AP3's amendment is honest and the row is satisfied; the axis
  re-derived the fingerprint argument independently and reached the build's
  conclusion. Minor schema smell: the amended seam cell now carries rationale
  rather than a seam name.
- **P2 · no-op.** SH1–SH3 name a predicate unit seam; the build tested them
  through the command, which is strictly stronger.
- **P3 · no-op.** `projects/benchkit.md` placement is correct — the appended
  lines land in the worktree-lifecycle bullet that already owns `create` and
  `land`. No dedicated paragraph is warranted; ticket 03's flag is answered.
- **P4 · ask-user.** `resume.go:355-357` silently converts an unreadable pool
  into "0 reclaimable" — the one state where RS3's guarantee that the ambient
  number cannot disagree with the verb stops holding, since the verb would exit
  1 with a structured error. The posture (don't fail a resume over an unrelated
  read) is right; the silence is the questionable half, and any fix is new
  resume-summary behavior with no spec row.

All 21 rows verdict delivered. No Won't-handle item was quietly built, the AXI
approved-children set is unchanged, and every written path is fenced — though
the axis notes the build widened its own fence to admit
`cmd/bench/main_test.go`, which the reviewer approved in conversation.

## Coverage

Count 4. Worst: C1 — a relative `gitdir:` pointer makes a live worktree read as
provably absent and eligible for permanent deletion.

- **C1 · auto-fix, must land before this ships. Coordinator-confirmed by direct
  experiment**, not by reading: a pool child whose `.git` reads
  `gitdir: <relative path>` naming a *live* repository classifies as `reclaim`
  with reason "every child points at an absent repository".
  `pool_reclaim.go:177` stats the pointer value as-is and `gitdirTarget:188`
  returns it unresolved, so a relative target resolves against the process CWD —
  which is the operator's repository, never the pool child. Git emits relative
  pointers under `worktree.useRelativePaths` or `git worktree add
  --relative-paths`. Every fixture writes an absolute target, so no test and no
  row exercises it. This contradicts the spec's own "Proven absence only"
  decision: CWD-relative resolution is not proof, so the repair is inside
  approved scope rather than a spec change.
- **C2 · auto-fix.** `gitdirTarget:194` applies `strings.TrimSpace` to the
  pointer value, so a repository path legitimately ending in a space reads as
  absent. Trim only `\r` and `\n`.
- **C3 · ask-user, same target as S4.** `pool_reclaim.go:338-341` returns a
  `retained` verdict when `RemoveAll` fails, and `:265` returns 0 regardless. A
  root-owned key fails with EACCES mid-loop and the operator's script reads exit
  0 as complete. No test asserts an exit code for a removal failure.
- **C4 · auto-fix.** The spec's Edge inventory claims a FIFO, device, or socket
  where a key, child, or `.git` belongs is retained unread, but no fixture uses
  `Mkfifo`. Deleting the `IsRegular` check at `:166` leaves the suite green.

Refuted, not reported: fingerprint collision (`canonicalParts` length-prefixes,
injective); uppercase and short fingerprints rejected; a dangling-symlink target
retains, which is the safe direction; NUL in a target returns EINVAL and
retains; the resume tests are non-vacuous.
