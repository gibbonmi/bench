# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-08-02 — parked staging defect did not reproduce as diagnosed  [open]

- **What happened:** The parked `bench commit` deleted-directory staging failure
  was issued as a light-path fix with a prescribed shape (`git add -A` for named
  directories). That shape was already implemented and contract-pinned; the
  retire repro committed green. The field error reproduced byte-for-byte only
  under a held `.git/index.lock` — consistent with the concurrent second
  session — and was misdiagnosed because staging discarded git's stderr. I
  landed the stderr-relay fix instead, flagged in the ticket, and left
  retry-on-held-lock as a reviewer decision.
- **Right behavior:** Reproduce a parked defect against the current binary
  before building its prescribed fix; when the diagnosis falsifies, fix the
  observable that misled it and surface the behavior question rather than
  landing a no-op.
- **Proposed rule change:** none — but decide whether `bench commit` staging
  should retry briefly on a held index.lock, since a post-gate staging failure
  costs a full green run and two sessions legitimately share this repository.

## 2026-08-02 — concurrent landing discarded this session's dirty paths  [open]

- **What happened:** While this session held finished, uncommitted work in the
  main checkout (waiting to serialize behind the pcgs session), the pcgs
  session's landing of `5d67654` left the tree fully clean: this session's
  modified and untracked paths were gone afterward, with no stash and no
  recovery ref. The work was restored from this session's context, but only
  because the full content happened to be in memory.
- **Right behavior:** A landing path that meets another writer's dirty paths
  should refuse (as `bench commit`'s block-check does) or set them aside
  recoverably — never leave them silently discarded. Two live sessions sharing
  the primary checkout also suggests side-work belongs in a worktree even for
  main-branch landings, serialized by explicit turn-taking.
- **Proposed rule change:** none proposed — reviewer should identify what the
  pcgs session ran to clean the tree and decide whether that surface needs a
  fail-closed or set-aside posture.

- 2026-08-02  Delegate backup files collided across a shared session scratchpad: a stale
  `mine/` directory from one write-delegate was swept up by a later delegate's restore glob
  and clobbered four out-of-fence files (caught by its return-time `git status` check and
  repaired). Right behavior: a delegate's transient backups live inside its own worktree
  under a unique name, and restore commands name exact files, never globs. Proposed rule:
  craft-delegate's isolation section names worktree-local backup paths alongside the
  existing cp-aside/git-show-restore pattern.
