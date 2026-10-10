# Route ticket and spec authors to the known issues

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/commands/bench-write-spec.md, docs/known-issues.md, CHANGELOG.md
Covers: none

## What to build

`docs/known-issues.md` says that a coordinator and a delegate read it, but no phase step makes them read it. So the witness-cap row did not reach the author of FT217 ticket 01, and that author ran 54 serial probes.

The implement-spec author charge names `docs/known-issues.md` by path. The write-spec author step reads that file before it plans the probes and the check floor. `.agents/commands/` is the one source of each phase command. `.claude/commands` is a link to that folder, and each Codex adapter skill only points to its command file.

`docs/known-issues.md` gains two rows in its existing shape:

- A bare `bench preflight build` grades `base-current` red between chunks after `main` moves. The workaround passes an explicit `--base` and `--source-tip`. FT342 owns the retire condition.
- `bench record verification` reads each excerpt from a file, and an unspilled response leaves no file. The workaround writes the excerpt to a file first. FT405 owns the retire condition.

The probe-cost row names one `bench test` call that repeats `--package` for the floor checks of several packages.

The rule that a charge names `docs/known-issues.md` by path and copies no row has one owner, the charge section of `craft-delegate`. This ticket does not state that rule a second time.

## Acceptance

- [ ] The author charge in `.agents/commands/bench-implement-spec.md` names `docs/known-issues.md` by path and copies no row from it.
- [ ] The author step in `.agents/commands/bench-write-spec.md` reads `docs/known-issues.md` before it plans the probes and the check floor.
- [ ] `docs/known-issues.md` has the `base-current` row with the FT342 retire condition and the excerpt-file row with the FT405 retire condition.
- [ ] The probe-cost row names one `bench test` call with a repeated `--package`.
- [ ] No other file states the charge rule from `craft-delegate` a second time.
- [ ] `bench test --check root-conformance` and `bench gate-prose` on each edited Markdown file stay green.
