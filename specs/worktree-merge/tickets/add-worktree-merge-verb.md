# Add the worktree merge verb

Blocked by: compose-merge-in-landing-owner.md
Writes: internal/worktree/merge.go (new), internal/worktree/merge_test.go (new), internal/usage/worktree.go, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go

## What to build

`bench worktree merge --from <commit|target> <target>` runs end to end on the
happy paths. The grammar constant joins the worktree usage block. The registry
gains the help row and the nested dispatch, and the `--help` description row
adds `merge` to its enumeration. The command selects the target's
assignment record, checks the active state, and validates the creation bundle.
Every refusal renders as the landing's `refused{...}` record on stdout.

`--from` resolves in two lookups: an assignment target, then `<from>^{commit}` required
to be an ancestor of the default branch tip. A sibling contributes its branch
tip.

The verb takes the target's fingerprint and resolves the
root's declared fast lane for the target worktree. A root with no lane keeps
the whole-project gate, as `bench commit` does. The verb calls the owner's
merge operation with the derived subject, reconciles the checkout to the
published commit, and prints the `merged{...}` record.

The subject is `merge: compose <from-spelling> <8-char commit> into <label>`.
The record is
`merged{worktree=<assignment id>,from=<commit>,kind=<kind>,previous_tip=<commit>,tip=<commit>,tree=<tree>}`.
The refusal and boundary tickets add the remaining predicates on this file.

## Acceptance

- [ ] WM1: on a diverged target, `merge --from <commit>` publishes a commit
      whose tree equals the merge-tree of the pair, and the checkout HEAD
      equals that commit.
- [ ] WM2: a target whose tip is an ancestor of `--from` prints
      `kind=fast-forward`, and the branch tip equals the incoming commit with
      no new commit object.
- [ ] WM3: a target that already contains `--from` prints `kind=current` at
      exit 0, and the tip, the checkout, and the lane record are unchanged.
- [ ] WM5: the published commit subject is
      `merge: compose <from-spelling> <8-char commit> into <label>`.
- [ ] WM6: `--from <sibling label>` publishes a commit whose second parent is
      the sibling's branch tip.
- [ ] WM7: stdout is one `merged{worktree=,from=,kind=,previous_tip=,tip=,tree=}`
      record with the exact values.
- [ ] WM8: the target operand accepts the label, the id, an 8–12 character
      prefix, and the absolute path, and an ambiguous prefix refuses naming
      both ids.
- [ ] WM18: a declared lane runs on the composed tree before the ref moves,
      and stdout carries `lane{outcome=pass,checks=...}`.
- [ ] WM25: `bench help` lists the merge row, and `bench worktree --help`
      prints the merge grammar.
- [ ] WM32: a missing `--from`, an empty `--from`, or a second positional
      exits 2 with usage on stderr.
