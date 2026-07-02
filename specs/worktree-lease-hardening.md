# Worktree lease staleness + TOCTOU hardening

## Problem
The pool's lease protocol is check-then-create: `worktree_acquire` scans for
directories without a lease file, then writes one. Two concurrent acquires can
both pass the check and share a worktree (TOCTOU). And a lease left by a killed
process (SIGKILL, crash — anything the INT/TERM trap can't catch) never
expires: the worktree is leased forever and the pool grows a zombie per crash.

## Solution
Leases become atomic and owned. A claim is an O_EXCL (noclobber) create that
records `<pid> <iso8601>`; the loser of a race moves to the next candidate. A
lease whose recorded pid is no longer alive is stale and reclaimable — the
reclaim itself is atomic (mv-aside, then a fresh claim), so two reclaimers
cannot both win. Release only removes a lease this process owns.

## User stories
1. As a shift/worktree user, I want two concurrent acquires to never share a
   worktree, so parallel shifts cannot corrupt each other's iteration commits.
   Line: claude-fable-5 (inline, session model) / medium. Concurrency protocol
   in the pool — the one seam in this batch where the answer is easy to get
   subtly wrong.
2. As a shift/worktree user, I want a lease whose owner is dead to be reclaimed
   on the next acquire, so a crashed run costs nothing and the pool stays warm
   instead of growing zombies.
   Line: claude-fable-5 (inline, session model) / medium. Same protocol.
3. As a shift/worktree user, I want release to remove only a lease it owns, so
   a reclaimed worktree's new lease can't be deleted by the old owner's
   deferred cleanup trap.
   Line: claude-fable-5 (inline, session model) / medium. Same protocol.

## Implementation decisions
- Claim = `set -C` noclobber create of the lease file with `<pid> <utc-time>`
  content; exists → read pid: alive → held; dead → atomic reclaim (mv the stale
  lease aside — only one mover wins — then claim fresh).
- A lease with unreadable/empty content (legacy format, or a crash mid-write)
  counts as stale only when older than a minute by mtime — fresh-and-empty is a
  writer mid-claim, not a zombie.
- Release reads the lease pid first: own pid (or unreadable) → clean and remove;
  another live pid → leave the worktree entirely alone (it was reclaimed).
- The dirty-tree skip in the scan stays (cheap preference for clean trees);
  correctness now rests on the claim, not the scan.
- `bench status`'s leased-pool signal keeps its semantics: lease present =
  in-flight; staleness is acquire's business, not the dashboard's.

## Testing decisions
- Good test = drive `bench worktree` (with the recording `$SHELL` pattern the
  existing lease/reuse contract uses) against planted lease states; assert
  which pool path each invocation got.
- Seam: the `bench` CLI (`worktree` subcommand) — the same seam the existing
  lease/reuse contract tests. No new seams.
- Gate command: `.bench/gate.sh`.

### Seam diagram

    trigger: bench worktree / bench shift (concurrent invocations)
        │
        ▼
    pool dirs + lease files ──▶ [ worktree_acquire:        ] ──▶ exclusive worktree path
    owner pid liveness      ──▶ [   scan → atomic claim →  ] ──▶ lease file <pid> <time>
                                [   stale reclaim          ]
                                     ◀ tests attach here: plant dead/live/legacy
                                       leases, run bench worktree, assert the
                                       path reused vs newly created

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | two overlapping `bench worktree` runs get distinct paths | CLI contract (rendezvous shell: each waits until both recorded) | fails under the old non-atomic claim when both scan before either leases | the shared-worktree corruption is the failure mode |
| 2 | lease with a dead pid → next acquire reuses that path | CLI contract | fails while a dead lease still reads as held (new pool dir created instead) | zombie leases are exactly this observable |
| 2 | empty legacy lease older than a minute → reclaimed | CLI contract (planted empty lease, backdated mtime) | fails while empty means held-forever | crash-mid-write and legacy leases must age out |
| 3 | lease with a *live* foreign pid → acquire creates a new path, lease untouched | CLI contract | fails if liveness check or claim atomicity is wrong | stealing a live lease is story 1's bug through another door |
| edge | release after reclaim: old owner's release leaves the new lease in place | covered by story-3 mechanics (release guard reads pid) — asserted in the same contract | lease file still present after foreign release | deferred trap cleanup must not unlease the new owner |

### Edge inventory
- interrupted/partial state → stories 2 (dead pid, empty lease) are the
  interrupt cases; the existing interrupt-cleanup contract still covers the
  trap path.
- re-run idempotency → existing lease/reuse contract (unchanged) proves
  release→acquire reuses the same path.
- hostile environment → **Won't handle:** pool on NFS or shared between hosts —
  pid liveness is per-host; the pool lives under `$HOME/.bench` and is
  documented as local.
- pid recycling (dead pid reused by an unrelated process) → **Won't handle:**
  the window is a same-host pid collision against a short-lived pool lease;
  the consequence is a warm worktree wrongly read as held (safe direction —
  never shared).
- boundary: lease exactly at the one-minute mtime edge → **Won't handle:**
  either verdict is safe (held = new dir created; stale = reclaimed after
  owner definitely gone).

## Out of scope
- Cross-host pooling or lock servers — different capability; not priced.
- A `bench worktree --sweep` that garbage-collects stale pool entries in bulk —
  separate command surface; ~6 edits, ~3 gate runs.
