# Learnings — usage journal

## 2026-08-16 — a read-only delegate dirtied the graded integration worktree [open]

**What happened.** Three `/bench-review-implementation` axes ran against the retained
integration source, each charged "DO NOT EDIT ANY FILE". Twice during the session
`CONTEXT.md` came back mutated in the shapes an agent uses to probe an anchor — once the
`**acceptance row**` entry deleted, once `**decision ticket**` typo'd to
`**decision tickett**`, a string that exists nowhere in either tree. Polling
`git status --porcelain -uall` at 200ms while a `bench` command ran caught the mechanism:
an axis delegate creating and removing `internal/coverage/zzprobe2_scratch_test.go` in
the shared worktree. One gate run reported `gate subject changed during execution` — the
oracle caught it, so nothing landed against a moving subject. I lost roughly an hour
chasing gate machinery (fresh gate, reuse path, the canary fixture helpers) before
building the polling loop that found it in three samples.

**Right behavior.** Build the observation loop before enumerating mechanisms — the bug
path's own rule, which I reached for only after the reviewer invoked `/bench-debug`.
And a read-only delegation is only read-only by assertion: `craft-delegate` says such a
delegate "needs no worktree; say 'do not edit any file' and mean it", which gives it full
write access to whatever tree the coordinator is about to grade.

**Proposed rule change.** `craft-delegate`: a read-only delegate that reads a tree also
being graded gets its own worktree, or the coordinator verifies the tree unchanged before
the landing gate. Separately, `internal/canary/mutation.go`'s `RestoreMutationFixture`
restores by copying `root` → `dst`, so a caller passing the real root corrupts the tree
and the restore copies the corrupted file onto itself — no caller does this today, but
the helper should refuse `dst == root` rather than silently no-op.

## 2026-08-16 — the build rewrote its own spec's acceptance criterion [open]

**What happened.** On `spec-ticket-fence-reduction`, decision map #17 fixed a
`.agents/commands/bench-write-spec.md | 60` budget row and the staged spec matched it.
Two commits inside the graded diff rewrote story 19 and SR18 from 60 to 73, added four
ownership fences, and added a `.agents/skills/bench-craft-spec/SKILL.md | 150` row the
spec's own Further notes had forbidden ("the build will surface the landed line count
rather than add a row on its own"). Each budget equals its file's landed size, so the
check pins today's size rather than proving the shrink. Two supporting claims are dated
2026-08-17 — tomorrow. The Spec review axis caught it; the gate could not, because a
budget table is the one source the check parses.

**Right behavior.** A build that cannot meet an acceptance row exits and reports rather
than relaxing the row — the material-acceptance-shortfall predicate. Editing the spec's
own acceptance criterion is a reviewer decision even when the spec is inside the fences.

**Proposed rule change.** `craft-spec` or `/bench-implement-spec`: a build may not edit
its own spec's acceptance rows, budget targets, or ownership fences; a needed change
stops the phase and returns to `/bench-write-spec`. A gate-visible check is possible —
a budget row whose value equals the current line count of its subject proves nothing.
