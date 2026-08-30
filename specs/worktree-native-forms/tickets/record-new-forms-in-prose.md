# Record the new forms in the changelog and the glossary

Blocked by: add-worktree-build-verb.md, add-create-from-sibling-start.md, run-system-suite-as-named-check.md, fill-preflight-next-column.md, name-the-gate-form-in-both-inventories.md
Writes: CHANGELOG.md, CONTEXT.md

## What to build

A reader opens the release notes and the glossary and finds this slice's forms.
`CHANGELOG.md` records the four added forms and the replaced gate line. The four
forms are `bench worktree build <target>`, `bench test --check system`,
`bench worktree create --from <target>`, and the preflight `next` column.
`CONTEXT.md` defines the term **worktree build**, and it names the system check.

This ticket writes prose only, so it lands last and it describes the landed
behavior. The two stories carry no coverage row, and the review round grades the
prose. The reviewer reads the changelog against the five landed behaviors above,
and the glossary against the build verb and the system check. Every sentence
obeys ASD-STE100.

## Acceptance

- [ ] `CHANGELOG.md` names `bench worktree build <target>` and its `dist/bench` output.
- [ ] `CHANGELOG.md` names `bench test --check system` and `bench worktree create --from <target>`.
- [ ] `CHANGELOG.md` names the preflight `next` column and the `base-current` remedy.
- [ ] `CHANGELOG.md` records that `bench worktree exec <target> -- bench gate` replaces the raw gate line.
- [ ] `CONTEXT.md` defines **worktree build** and names `bench test --check system`.
- [ ] `bench gate-prose` stays green for both edited files.
