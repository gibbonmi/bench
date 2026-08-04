# Record the two new verbs in the changelog

Blocked by: none
Ownership fence: `CHANGELOG.md`
Contracts: the shipped invocation grammar of both new verbs crosses `bin/bench.sh` and `internal/usage/worktree.go`→`CHANGELOG.md` and is asserted by CL1 by reading the entry against the shipped usage strings rather than against this build's description of them
Assumptions: `CHANGELOG.md` already exists and carries a typed entry format this ticket follows rather than introduces; both verbs shipped in this build are user-visible CLI surface

## What to build

`craft-synthesis` requires one concise typed `CHANGELOG.md` entry for user-visible
behavior, and this build adds two operator-facing verbs plus a deliberate behavior change
to an existing one. None of it is recorded.

Write entries covering:

- `bench worktree recovery <ref> --discard <fingerprint>` — retires a preserved recovery
  payload the landedness proof does not accept, one ref per invocation, exact fingerprint
  required.
- `bench spec build reclaim <slug>` with `--apply <fingerprint>` — maintainer-run plan and
  apply over one terminal run's leftover provisional refs, deleting only what it can prove
  dead.
- The behavior change to `bench worktree recovery --apply`: supplying a fingerprint with a
  plan that authorizes no action now refuses instead of returning silent success. A caller
  who passed a fingerprint asked for an action, and exit-zero read as "the work is gone"
  when it was not.

Match the file's existing entry shape and register exactly — read it before writing. Keep
each entry to what an operator needs to know that they did not know before; this is a
changelog, not a summary of the build.

## Acceptance

- [ ] [CL1] both new verbs appear with their invocation grammar as shipped, verified against `bin/bench.sh` rather than from memory.
- [ ] [CL2] the `--apply` silent-success removal is recorded as a behavior change, not as a new feature.
- [ ] [CL3] the entries follow the file's existing typed format, with no new heading style or section introduced.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CL1 | write the discard grammar as `--discard` with no fingerprint argument | the reviewer, reading the entry against the shipped usage string | compare the entry to `usage.WorktreeRecovery` in `internal/usage/worktree.go` and to the `bin/bench.sh` help block, expect the documented grammar to contradict the shipped one |
| CL2 | file the `--apply` change under added features | the reviewer, reading the entry | move the line, re-read it against the spec's implementation decision that it is a deliberate behaviour change, expect the entry to misdescribe an existing verb as a new one |
| CL3 | introduce a new heading level for these entries | the reviewer, reading the diff | add the heading, compare against the surrounding entries, expect the file to carry two formats for one fact |
