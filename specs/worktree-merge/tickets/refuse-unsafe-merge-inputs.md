# Refuse unsafe merge inputs

Blocked by: add-worktree-merge-verb.md
Writes: internal/worktree/merge.go, internal/worktree/merge_test.go

## What to build

Every refusal the verb makes before publication, each one predicate rendered
on stdout as the landing's `refused{...}` record with its path table at exit
1. A `--from` commit outside the default branch's history that is no sibling
tip refuses naming the commit, and no lane runs. A composition conflict
outside the capture rule table names every path and changes nothing, and a
settled capture path is disclosed on stderr.

A dirty target refuses before
composition. A dirty or detached sibling refuses naming `bench commit` at the
sibling or its branch. An ambiguous, self, or unknown `--from` refuses with
the exact spellings. A failed identity component refuses naming the component.
A checkout off the assignment branch refuses naming the branch ref, and a
tip-versus-HEAD mismatch refuses naming both commits.

No new identity component joins the registry. This ticket and `grade-the-publication-boundary.md` both edit the
command file, so the boundary ticket runs after this one lands.

## Acceptance

- [ ] WM9: a conflicting non-capture path prints `refused{...}` and the path
      table on stdout at exit 1, names every conflicting path, and leaves the
      tip, the checkout, and the lane record unchanged.
- [ ] WM10: a conflicted `capture/learnings.md` composes as the union and
      stderr carries `merge composition{resolved=capture/learnings.md:union}`.
- [ ] WM11: a target with one untracked path or one modified path refuses
      before composition, names the path, and runs no lane.
- [ ] WM12: a sibling with a modified path refuses naming `bench commit` at
      the sibling, and a detached sibling refuses naming its assignment branch
      ref.
- [ ] WM13: a `--from` equal to an assignment label that is also a branch name
      refuses as ambiguous, naming the assignment id and the full commit.
- [ ] WM14: a `--from` that resolves to the target's own assignment refuses.
- [ ] WM15: a `--from` that names nothing refuses at exit 1 naming the value.
- [ ] WM16: a target whose owner marker or lock fails refuses naming that
      component, a non-active state refuses naming the state, and a checkout
      off the assignment branch refuses naming the branch ref.
- [ ] WM17: a target whose branch tip is not its checkout HEAD refuses naming
      both commits.
- [ ] WM31: a conflicting path that holds a control byte renders escaped in
      the stdout refusal table at exit 1.
- [ ] WM33: an add/add `capture/learnings.md` whose two sides differ in file
      mode refuses naming the path.
- [ ] WM34: a `--from` commit that is not an ancestor of the default branch
      tip and is no sibling tip refuses at exit 1 naming the commit, and no
      lane runs.
