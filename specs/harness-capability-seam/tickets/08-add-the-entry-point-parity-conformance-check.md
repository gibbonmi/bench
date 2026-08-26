# Add the entry-point-parity conformance check

Blocked by: none
Writes: internal/conformance/entry_point_parity_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, tests/canary/entry-point-parity/ (new), projects/benchkit.md

## What to build

A new conformance check, `entry-point-parity`, holds every shim, adapter,
front door, and the CI script to one registry name. It runs at the dev tier
over the root.

The check owns one table. Each entry maps a shim or an adapter basename to a
registry name and a canned input. The check enumerates the shims and the
adapters from disk, and an entry outside the table is a diagnostic. So a root
with `.bench/hooks/extra.sh` outside the table names `extra.sh`.

For each runtime row the check runs the entry through the stub wrapper with
`BENCH_COMMAND_OBSERVE=1`, and it runs the direct verb with the same input.
It requires the observed `command-registry:<name>` line on stderr, equal exit
codes, and the direct stdout as a suffix of the shim's stdout. The suffix
rule holds because `session-start.sh` prints one advisory line first. The
check reads only the id line on stderr and ignores a shim's own warnings. A
shim run with no wrapper on PATH is not graded.

Three rows are static. The `worktree-lifecycle.sh` row names `worktree-hook`
and runs nothing, because a create needs a live pool. The CI row grades the
exec line in `scripts/release-preflight.sh` and the run line in the workflow.
The front-door row grades the exact verb in `.agents/commands/bench.md`.

Every registry command with `internalInventory` is reached by a row, or it
carries an exemption reason in the check's own table. A command that neither
reaches nor exempts is a diagnostic naming the command.

The check joins the enumerated family everywhere the family already appears:

- the `conformanceChecks` map in `internal/conformance/checks_test.go`, bound as `checkEntryPointParity` at `registry.Dev` and `registry.SubjectRoot`.
- the check list in `internal/conformance/registry/registry.go`.
- the conformance table in `projects/benchkit.md`.
- one canary family directory, `tests/canary/entry-point-parity/`.

The check reads no `SKILL.md`, so it does not join `hostileSkillReaders`.
This ticket writes `internal/conformance/checks_test.go`, the registry, and
the profile's conformance table, which ticket 07 also writes. A coordinator
therefore serializes this ticket with ticket 07.

## Acceptance

- [ ] A root with `.bench/hooks/extra.sh` outside the table yields a diagnostic naming `extra.sh`. (covers HC37)
- [ ] `block-dangerous-git.sh` with a benign envelope through the stub wrapper prints `command-registry:guard-git` on stderr. (covers HC38)
- [ ] `stop.sh` with a `stop_hook_active` envelope exits like the direct `stop-verdict` and its stdout ends with the direct stdout. (covers HC39)
- [ ] `bin/bench.sh` with no argument prints the same bytes as `bin/bench.sh status --route`. (covers HC40)
- [ ] A `scripts/release-preflight.sh` whose exec line names `release-preflight-2` yields a diagnostic. (covers HC41)
- [ ] A `.agents/commands/bench.md` that names `bench status --routes` yields a diagnostic. (covers HC42)
- [ ] A registry command with `internalInventory` that no row reaches and no exemption names yields a diagnostic naming the command. (covers HC43)
- [ ] The `entry-point-parity` canary fixture turns the check red, and the restored fixture turns it green. (covers HC44)
- [ ] `session-start.sh` in a temp repo prints `command-registry:session-inspect` and its stdout ends with the direct stdout. (covers HC48)
- [ ] Each of the three adapters under `BENCH_MODEL=mid` prints `command-registry:resolve-model`. (covers HC49)
- [ ] The check reports no diagnostic over the live root.
