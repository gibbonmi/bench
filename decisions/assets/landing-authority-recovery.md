# Request token and recovery lifecycle (FT169 research)

Source: read-only Terra delegate, 2026-09-01. Coordinator spot checks:
`internal/worktree/identity_component.go:33`,
`internal/freshness/freshness_verify.go:80`,
`internal/worktree/land_identity.go:48` — all resolve.

## 1. The request token

The caller mints the token; Bench never does
(`internal/worktree/worktree.go:635`). `Create` refuses an empty token
(`internal/worktree/ownership.go:149`). Bench stores a SHA-256 digest as the
authorization identity and also persists the plain token beside it
(`internal/intent/assignment.go:456`, `:62-67`;
`internal/worktree/ownership.go:214`). The ledger is
`<git-common-dir>/bench-intent.json` (`internal/intent/intent.go:24`).

A lost token refuses with the detail
`request token matches no assignment`
(`internal/worktree/identity_component.go:33`,
`internal/worktree/land_refusal.go:106-118`). Two recovery routes exist.
When exactly one active assignment owns the worktree, the refusal names
`bench worktree reauthorize --assignment <id> --request <new-request>
--base <commit> --source-tip <commit> <path>`
(`internal/worktree/worktree.go:592-611`). Otherwise the next step is
`bench worktree list`, and the listing prints the persisted token
(`internal/worktree/list.go:192`). Reauthorization swaps the digest under a
lock after ownership proofs (`internal/worktree/reauthorize.go:79-116`).
A record from before the `request_token` field carries no token
(`internal/intent/assignment.go:66`); reauthorization is its only route.

## 2. The census record

`~/.bench/census/<pool-key>/<assignment-id>` records raw shell heads
(`internal/census/census.go:24`, `:218`). Retirement deletes the record,
not the landing: `executeCleanup` drops it on release, clean, and the
landing's own release step (`internal/worktree/lifecycle.go:425-441`). The
landing reads the count before the release step and prints `census heads`
in the landed record (`internal/worktree/land.go:259-272`). The assignment
record is separate, in the intent ledger
(`internal/intent/assignment.go:473-560`).

## 3. The stale executable seal

`freshness.Verify` compares stored and current source digests
(`internal/freshness/freshness_verify.go:14-31`). The message is
`bench binary "<path>" is untrusted: <cause>; rebuild with <action>` and
names the copy-paste rebuild command (`freshness_verify.go:79-86`). The
wrapper's landing route separately requires the installed promotion broker
and advises `bench doctor --fix` (`bin/bench.sh:364-376`).

## 4. Ignored residue at the destination

The first-run landing refuses with
`landing destination has undeclared ignored residue` plus a path table
(`internal/worktree/land_identity.go:44-51`). The allowance comes from
`.bench/build-outputs.json` plus runtime-ignored and local-capture paths
(`internal/worktree/build_outputs.go:73-87`). The resume path applies the
same allowance with a shorter detail and no table
(`internal/worktree/landingpolicy/landingpolicy.go:41-77`).

## 5. Incomplete-publication recovery

The destination CAS is the commit point. A failure after it prints a
terminal `landed{...,next=...}` record whose next is the exact
`bench worktree land --resume <published-commit> ...` continuation
(`internal/worktree/land.go:295-310`,
`internal/worktree/land_refusal.go:39-54`). The resume authenticates the
published commit by its three parents
(`internal/worktree/land_resume.go:222-252`). It reads the assignment in
`active` or
`cleanup-pending` state, and falls back to the cleanup receipt
(`internal/worktree/land_resume.go:92-105`, `:174-252`). The census file
survives an interrupted landing because the release step never ran.

A textual conflict fires before publication, so it is a refusal. Its next
names the hand merge, then `bench commit`, then the review, then the
re-run (`internal/worktree/land_refusal.go:56-94`). Nothing distinguishes
a committed conflict resolution from a pending `git merge --continue`:
`internal/commit/commit.go` reads no `MERGE_HEAD`. FT258 holds the
undecided `MERGE_HEAD` contract for `bench commit`; FT254 is blocked on it.

## Read and not read

The delegate read the worktree, intent, census, freshness, and runbinary
packages, `bin/bench.sh:340-390`, and `.bench/build-outputs.json`. It did
not read the test files, `internal/preflight`, `internal/landing`
internals, or the FT254/FT258 bodies beyond grep-matched lines.
