# Learnings — usage journal

<!-- entries below -->

## [open] 2026-08-24 — bench commit cannot stage an index-only deletion of an ignored path

What happened: `git rm --cached capture/session-handoff.md` staged a deletion, but `bench commit -- <path>` reported "nothing to commit" because the path is now git-ignored. The workaround was a working-tree deletion (preserve copy first), which bench commit accepted.

Right behavior: bench commit should honor an already-staged change at a named path, ignored or not.

Proposed rule change: none yet; consider a bench commit fix if untracking recurs.


## 2026-08-24 — the landing pays two refusals to local capture state [open]

What happened: the FT238 landing refused twice. First, `capture/learnings.md`
is ignore-listed but still tracked, so its local edits dirtied the destination.
Second, `capture/session-handoff.md` is now untracked ignored residue that
`.bench/build-outputs.json` does not declare. Both files were set aside by
hand and restored after the landing.

Right behavior: land the untracking of both capture inboxes, which the ignore
rules already decide. Declare `capture/session-handoff.md` in the landing
allowance, or give the handoff its own allowance entry.

Proposed rule change: none; the decisions exist. Only the two file states lag.
