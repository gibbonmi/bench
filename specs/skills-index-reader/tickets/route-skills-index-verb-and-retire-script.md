# Route the skills-index verb and retire the shell script

Blocked by: extract-skills-index-module.md
Writes: internal/skillsindex, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go, bin/bench.sh, internal/conformance/subcommand_routing_test.go, internal/conformance/validity_checks_test.go, .bench/BENCH.md, .bench/BENCH-reference.md, projects/benchkit.md, CHANGELOG.md, decisions/deepening-2026-08.md, .bench/skills-index.sh

## What to build

`bench skills-index [--check|--write]` — `skillsindex.Command` with a `usage.Grammar`,
`git.Root()`, registered in `cmd/bench` through `outputCommand` with the AXI mutation
exemption — prints the module's diagnostics one per line (exit 1) or nothing (exit 0),
regenerates on `--write`, usage on 2. In the same commit: the `bin/bench.sh` case label
and help line, the subcommand-routing row (routed to `internal/skillsindex`), the
`bench skills-index` token in `.bench/BENCH.md`'s Oracle bullet (bumping the
`projects/benchkit.md` prose-budget row only if rewrap adds a line), `.bench/skills-index.sh`
deleted, its `checkShellSyntax` pattern dropped, `.bench/BENCH-reference.md`,
`projects/benchkit.md`, and the `kitOnlySkillSources` comment re-pointed to the verb, a
CHANGELOG line, and the one-line supersession note under map ticket #13. A thinner cut
strands the cold-pickup CLI sweep red (documented verb without a label, or label without an
inventory line), so these land together.

## Acceptance

- [ ] covers SI7: in a temp repo the verb prints exactly the missing-entry line and exits 1; `--write` then `--check` exit 0 with no output; `--help` exits 0; `--bogus` exits 2.
- [ ] covers SI8: the script is gone and `rg --hidden` finds only the permitted residue.
- [ ] covers SI9: `checkColdPickupCLILists` and `checkSubcommandRouting` are green with the label, inventory token, and routing row present.
- [ ] covers SI10: no pre-existing assertion changes.
