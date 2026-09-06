# Name bench probe as the coordinator's probe form

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md
Covers: none

## What to build

The probe verb shipped on 2026-09-06 with this grammar:

```
bench probe <file> (--swap <old> --with <new> | --omit <old>) (--package <expr> [--run <go-regex>] | --check <name>) [--full]
```

The verb preserves the subject under the Bench home and applies one exact swap
or omission. It runs one focused test or named check. Then it restores the
subject and proves the restore byte-exact. It prints `bit`, `silent`, `invalid`, or `restore-failed`.
The delegate skill still teaches the copy-aside sequence as the probe form.

Make `bench probe` the one probe form in `craft-delegate`. The `## Isolation`
paragraph in `SKILL.md` keeps the `git stash` ban and names `bench probe` as the
substitute for a probe. The copy-aside sequence stays only for a non-probe edit
that must test the committed version.

The `## Probes` rules in `references/delegation-discipline.md` say three things. A coordinator probe runs
through `bench probe`. A probe the verb reports as `invalid` proves nothing and
is replaced before a verdict is read. A charge names `bench probe` as the
delegate's self-probe form, so the copy-aside sequence leaves the charges.

## Acceptance

- [ ] `.agents/skills/bench-craft-delegate/SKILL.md` names `bench probe` as the probe form under `## Isolation`, and the `git stash` ban stays.
- [ ] `.agents/skills/bench-craft-delegate/references/delegation-discipline.md` `## Probes` names `bench probe` as the coordinator's probe form and states that an `invalid` verdict is replaced before a verdict is read.
- [ ] The charge rule says a self-probe names `bench probe`, so no charge carries the copy-aside sequence.
- [ ] `bench gate-prose . -- <both files>` exits 0.
- [ ] `bench test --check guidance-prose-budgets` and `bench test --check kit-compliance` exit 0 over the worktree.
