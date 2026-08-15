# skills-index-reader

Status: implemented

Decision source: decisions/deepening-2026-08.md (compiled ready map, reviewer-resolved 2026-08-15; tickets #4 and #13 govern this spec; the shell-script disposition was closed in conversation 2026-08-15 — delete `.bench/skills-index.sh`, the regenerator becomes `bench skills-index --write`)

Verification log: spec 3 + tickets 0 iteration(s) to accept — loop 1 (Claude `opus`, medium, read-only) returned 7, 2, then 1 blocking findings, all folded (second-parser guard row SI6, `bin/bench.sh` label and fence, module-owned `Command` for routing, SI5 literal and gate posture, ticket-1 ownership of the `checkSkillsIndexGenerateVerify` removal, `tier_test.go` registration, prose-budget clause); the round-4 acceptance pass and loop 2 were stopped by the reviewer ("no more reviews", 2026-08-15) — acceptance is by reviewer decision, not a clean re-review, and the ticket breakdown is unreviewed. Reviewer sign-off 2026-08-15: spec, two tickets, and fences approved as presented, including the map #4 exit-test widening.

## Problem

The skills index in `.bench/BENCH-reference.md` is derived twice. `.bench/skills-index.sh`
reads each skill's `index:`/`index-note:` frontmatter with awk, reads
`.bench/consumer-payload.json` with a line-shape awk match to mark kit-only skills, and
checks or rewrites the marker block; `internal/conformance` re-implements the same
frontmatter read (`frontmatterField`), the same kit-only marking (`kitOnlySkillSources`),
the same line rendering and the same attributed diagnostics (`checkSkillsIndex`), and
then probes the shell script's generate/verify contract in a temp tree
(`checkSkillsIndexGenerateVerify`). ASSESSMENT M-W1 and roadmap FT89 name this: two
parsers, one fact. A change to the line shape or the frontmatter rule must land twice,
and the shell reader has already had to grow its own parse-failure guard for a JSON
key order the Go reader handles for free.

## Solution

One shipped Go module owns the skills index: reading a skill's frontmatter fields, the
kit-only marking from the consumer-payload allowlist, rendering one index line per
indexed skill, checking the committed block with the same attributed diagnostics the
gate emits today, and regenerating it idempotently. Two consumers invoke it: the
conformance check `checkSkillsIndex` (in-process, so the gate's oracle is the module) and
a new operational verb `bench skills-index [--check|--write]`, which replaces
`.bench/skills-index.sh` as the operator's regenerator. The script is deleted and every
reference to it re-points to the verb. The checked-in index block stays the
independently authored expectation (ADR 0006): four of the fixtures under
`tests/canary/skills-index-command-adapters/` keep proving that an unindexed, dangling,
field-less, or stale-worded skill goes red through the real check. Rows are data read by
the renderer and the checker (ADR 0012's declarative-registry shape); no second
enumeration of the line format survives. FT89's roadmap row stays in place — this is its
single-sourcing slice, not its whole.

## User stories

1. **One reader behind the gate.** As the gate, I grade the skills index through one
   module: `checkSkillsIndex` calls it and carries no frontmatter, allowlist, or
   line-rendering logic of its own; the conformance `frontmatterField` helper is the
   module's field reader; the generate/verify contract (red on an empty block, green
   right after write, idempotent second write, kit-only marker generated from an
   allowlist whose row orders `audience` before `source`) is graded in-process against
   the module. Every diagnostic string the canaries pin survives byte-for-byte except
   the regenerate hint, which names the verb. Line: opus / medium. Gate and conformance
   logic take mid effort; the oracle's correctness matters more than speed.
2. **One regenerator for the operator.** As an operator, I run
   `bench skills-index --write` to regenerate the block and `bench skills-index`
   (`--check`, the default) to see the same attributed drift lines the gate would emit;
   `.bench/skills-index.sh` no longer exists, `bin/bench.sh` labels the verb, and every
   document that named the script names the verb. Line: sonnet / medium. CLI plumbing at a known seam with the routing registry,
   inventory doc, and stale-reference sweeps all gate-observable.

The decision map's ticket #13 still reads "the script and the conformance check both
invoke it"; the build appends one line under its Answer recording the 2026-08-15
supersession (script deleted, regenerator `bench skills-index --write`) so a cold reader
does not re-derive the dead shape.

Story partition: story 1 lives in `internal/skillsindex` + `internal/conformance`; story
2 in `cmd/bench` + docs + the deletion. They share the module seam, so this is one spec
with two serial tickets (expand, then contract), not a split.

## Implementation decisions

- New package `internal/skillsindex`, linked into the binary. Its interface: read one
  frontmatter field from a `SKILL.md` (first value inside the leading `---` fence, key
  followed by `:`, value trimmed; nothing after the closing fence counts); enumerate the
  index entries of a root as data — skill name, trigger, optional note, kit-only flag —
  in alphabetical directory order, skipping a skill whose name has a same-named
  `.agents/commands/<name>.md` (a command adapter); render one entry to its line
  (`- <trigger> → \`.agents/skills/<name>/SKILL.md\`[ (kit-only)][ + <note>]`); check a
  root's `.bench/BENCH-reference.md` block between the `bench:skills-index` markers and
  return the attributed diagnostics; write the block in place, restoring mode 0644.
- Kit-only marking reads `.bench/consumer-payload.json` through the root package's
  `PayloadRow`/`PayloadKitOnlyPrefixes` — the payload reader that already exists — so
  the module adds no JSON parsing of its own. An absent allowlist marks nothing. An
  allowlist that fails to parse: `--write` refuses with the literal
  `.bench/consumer-payload.json unreadable: kit-only marking unresolved (write refused)`
  and leaves the reference file untouched — the shell reader's fail-closed posture — while
  `check` (the gate's `checkSkillsIndex` and `--check`) keeps today's posture and marks
  nothing, so the gate gains no diagnostic the decision source did not ask for.
- Diagnostics keep today's exact strings. Four carry the regenerate hint
  `(regenerate: bench skills-index --write)`: `skills index entry for '<n>' drifted from
  its frontmatter`, `skills index missing entry for skill '<n>'`, `skills index entry
  '<n>' has no indexed .agents/skills/<n> on disk`, `skills index block drifted from
  generated form`. Three carry none: `skill '<n>' missing index: frontmatter (the skills
  index is generated)`, `.bench/BENCH-reference.md missing (skills index unverifiable)`,
  `.bench/BENCH-reference.md skills-index markers missing (bench:skills-index)`.
- `bench skills-index [--check|--write]`: the module owns the command entry —
  `skillsindex.Command(args) (string, int)` parses its `usage.Grammar` (`--check`,
  `--write`, `--help`) and resolves the root with `git.Root()`, the `structure.Command`
  pattern — and `cmd/bench` registers it through `outputCommand`, AXI-exempt with the
  mutation reason (`--write` changes a tracked file; not an approved query, so no
  `help[]` envelope). `--check` (default) prints each diagnostic on its own stdout line,
  exit 1 on any, exit 0 clean; `--write` regenerates and exits 0, or prints the blocking
  diagnostics and exits 1 without touching the file; usage exits 2; outside a repo the
  standard not-in-repo refusal. The verb lands in one commit with: the `bin/bench.sh`
  case label (`skills-index) route_porcelain "$@" ;;`) and its help line — the cold-pickup
  CLI sweep derives the known command set from those labels, so a documented verb without
  one is red — the subcommand-routing registry row (routed to `internal/skillsindex`), and
  the `bench skills-index` token in `.bench/BENCH.md`'s CLI inventory (Oracle group) —
  inserted into the existing Oracle bullet, which sits at its 180-line prose budget; if
  rewrapping adds a line, the same commit bumps that budget row in `projects/benchkit.md`.
  Diagnostics are returned in
  skill-alphabetical order, then the unattributed block line, so the verb and the check
  print the same sequence — replacing today's map-iteration order, in which the
  `missing index:` diagnostics lead as a block and the rest are nondeterministic.
- Ticket 1 (module + collapse) also removes `checkSkillsIndexGenerateVerify` (it lives in
  `skills_index_checks_test.go` and carries all three literals SI6 bans; its contract
  moves into the module's tests, SI3) and its call site in `docs_workflow_checks_test.go`,
  so ticket 1 is green on its own while the script still exists and `checkShellSyntax`
  still lints it. The four regenerate-hint literals (in `checkSkillsIndex` today, in the
  module after the collapse) name `bench skills-index --write` from ticket 1 — one commit
  before the verb routes; that is gate-safe because the cold-pickup reverse sweep reads
  only markdown, never Go source.
- Ticket 2 deletes `.bench/skills-index.sh`; the shell-syntax pattern list drops the
  path; `.bench/BENCH-reference.md`, `projects/benchkit.md`, and the "generator" comment
  on `kitOnlySkillSources` name the verb. The canary registry's `ShellSources` entries for
  the four fixtures keep naming the deleted script (that grader tolerates a missing
  retired twin, and the rows record which shell twin retired) — the build leaves them.
- Exit-test rule (map #4): every pre-existing test passes with test logic unmodified.
  Map #4 permits only renames; this spec widens that to the mechanical edits below because
  the conformance *implementations* live in `_test.go` files, so collapsing a parser is a
  test-file edit by location, not an assertion change. The permitted edits are —
  `checkSkillsIndex`, `kitOnlySkillSources`, and `frontmatterField` collapsing to calls
  into the module, and `markerBlock` (whose only caller was `checkSkillsIndex`) deleted; the removal of
  `checkSkillsIndexGenerateVerify` and its call site; the `checkShellSyntax` pattern entry;
  the subcommand-routing registry row; the SI6 guard test plus its one
  `classifiedLiveTreeTests` entry in `internal/conformance/tier_test.go` (the guard reads
  the live kit tree, exactly as its prior art `TestFixtureBiteProofArchitecture` does).
  A changed assertion reverts the move.

## Testing decisions

- A good test drives the module through its public seam on a temp root it builds
  (skills, allowlist, reference file) and observes lines, diagnostics, file bytes, and
  mode — never the internals of the fence scan.
- Seams: the module seam (`internal/skillsindex` package tests, new); the conformance
  seam (`checkSkillsIndex` graded by the four canaries, prior art
  `internal/conformance/fixture_bite_test.go`); the CLI seam (`cmd/bench` dispatch table
  tests and the subcommand-routing check, prior art `internal/conformance/subcommand_routing_test.go`).
- The gate observes the feature through the `test` phase (module + conformance) and
  through conformance's docs sweeps (command inventory listing, stale command
  references).

### Seam diagram

    trigger: gate `test` phase (conformance) · operator `bench skills-index [--check|--write]`
        │
        ▼
    root ──▶ [ internal/skillsindex: frontmatter field · entries (data) · render · check · write ] ──▶ diagnostics | rewritten block
                  ◀ tests attach here: temp root with .agents/skills/*, .bench/consumer-payload.json,
                    .bench/BENCH-reference.md; assert lines, diagnostics, bytes, mode
        │
        ├──▶ conformance checkSkillsIndex → canary fixtures (already red-capable)
        └──▶ cmd/bench dispatch → stdout lines + exit code

### Acceptance coverage map
| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| SI1 | 1 | on a root with two indexed skills, one command-adapter skill, and one kit-only skill declared with `audience` before `source`, entries render exactly `- <trigger> → \`.agents/skills/<name>/SKILL.md\` (kit-only)` for the withheld one and no line for the adapter, alphabetical | module seam | new test; red before the module exists (compile) then behavioral red asserted with a literal expected block | pins rendering, adapter skipping, and order-independent kit-only marking — the fact the shell reader once dropped |
| SI2 | 1 | frontmatter field read: first value wins, a key after the closing fence is not read, `index-note` absent yields empty | module seam | new test with a literal fixture | pins the fence rule both old readers agreed on |
| SI3 | 1 | check on a reference file whose block is empty returns `skills index missing entry for skill '<n>' (regenerate: bench skills-index --write)` per skill; write then check returns none; a second write leaves bytes identical and mode 0644; on a root mixing a missing-`index:` skill `alpha`, a drifted committed entry for `beta`, and a committed entry for absent `gamma`, check returns exactly the three attributed lines in that (alphabetical) order and no `block drifted` line | module seam | new test; the generate/verify contract moved in-process | the contract `checkSkillsIndexGenerateVerify` graded via the shell script keeps an owner, and the deterministic order the verb and the gate share is observed |
| SI4 | 1 | check attributes drift: a committed line whose trigger differs → `drifted from its frontmatter`; a committed line for a skill not on disk → `has no indexed .agents/skills/<n> on disk`; a skill with no `index:` → `missing index: frontmatter (the skills index is generated)`; a non-attributable difference → `block drifted from generated form` | conformance seam | already covered — canaries `unindexed-skill`, `dangling-index`, `missing-index-field`, `stale-index-wording` under `tests/canary/skills-index-command-adapters/`, graded by `internal/conformance/fixture_bite_test.go` | the checked-in expectation is independent of the module (ADR 0006); if the move changes a string, the fixture goes red |
| SI5 | 1 | with `.bench/consumer-payload.json` containing `{not json`, `Write` returns the error `.bench/consumer-payload.json unreadable: kit-only marking unresolved (write refused)` and the reference file's bytes are unchanged; `Check` on the same root returns no allowlist diagnostic (today's posture) | module seam | new test | the fail-closed posture the shell reader carried must not be lost in the move, and the gate must not gain a diagnostic the source did not ask for |
| SI6 | 1 | a conformance guard test (`Test*`, registered in `classifiedLiveTreeTests`; prior art `TestFixtureBiteProofArchitecture`'s `go/ast` self-parse) parses `internal/conformance/skills_index_checks_test.go` and `checks_test.go` and fails if the marker literal `<!-- bench:skills-index:start -->`, the string `consumer-payload.json`, or the line-format literal `\u2192 \`.agents/skills/` appears in either file outside the guard's own `ast.FuncDecl` (hoisting a literal to a package-level const or helper trips it) | conformance seam (new guard test) | red today — all three appear in `skills_index_checks_test.go` (`checkSkillsIndex`, `kitOnlySkillSources`, `checkSkillsIndexGenerateVerify`); `checks_test.go` carries none and is guarded forward-only | the cheapest wrong story 1 adds the module and leaves the old parsers in place; this is the row that sees it |
| SI7 | 2 | in a temp repo (`gittest.RepoOnBranch` + chdir, the `command_registry_test.go` pattern) with one indexed skill and an empty block, `bench skills-index` prints exactly `skills index missing entry for skill '<n>' (regenerate: bench skills-index --write)` and exits 1; `--write` exits 0 and a following `--check` prints nothing and exits 0; `--help` prints usage on stdout, exit 0; `--bogus` exits 2 | CLI seam (`cmd/bench` dispatch test) | new test | the verb is the operator's only regenerator once the script is gone |
| SI8 | 2 | `.bench/skills-index.sh` does not exist; `rg --hidden -g '!.git' skills-index.sh` matches only `internal/conformance/registry_test.go` (the retired-twin rows), `decisions/`, `ASSESSMENT.md`, `specs/`, plus `CHANGELOG.md` and `capture/` where the landing writes them | source (`rg --hidden`) | not TDD-able — deletion + sweep; verified by review | a surviving reference is a dead pointer to a regenerator that no longer runs |
| SI9 | 2 | `bin/bench.sh` carries the `skills-index)` case label and help line, `.bench/BENCH.md` lists `bench skills-index`, and the subcommand-routing registry routes it to `internal/skillsindex` (whose `Command` calls `usage.Parse`) — all in one commit | conformance docs sweeps | already covered — `checkColdPickupCLILists` is red on a documented `bench skills-index` with no `bin/bench.sh` label and on a labelled command missing from the inventory; `checkSubcommandRouting` is red on an unregistered dispatch name or a routed package that never reaches `usage.Parse` | an unlisted or unlabelled verb is the 2026-08-14 CLI-inventory learning recurring |
| SI10 | 2 | every pre-existing test passes with test logic unmodified — the permitted mechanical edits are exactly the exit-test list in Implementation decisions; a changed assertion reverts the move (map #4) | ordinary `test` phase + `git diff -- '*_test.go'` | exit test | the refactor's only oracle is that nothing observable changed |

### Edge inventory
- error path — SI3/SI5/SI7 (missing markers and missing reference file keep today's diagnostics; unparseable allowlist refuses).
- empty/absent input — no skills on disk: block renders empty and check passes (SI3's pre-write state with zero skills) — row folded into SI3's fixture as a zero-skill sub-case; absent allowlist: no kit-only marking (SI1 asserts the marker only where declared).
- boundary values — one skill, one adapter (SI1).
- malformed input — SKILL.md with no frontmatter fence: `index` reads empty → `missing index:` diagnostic (SI4 `missing-index-field`).
- interrupted or partial state — write is temp-file + rename today; **Won't handle** a crash mid-rename — the reference file is tracked in git and regenerated by the next `--write`; nothing downstream reads a torn block before the gate reruns.
- re-run idempotency — SI3.
- process-boundary lifecycle — SI7 drives the built verb; the module's file bytes are what the next process reads.
- hostile environment — **Won't handle** a skill directory name that is not `[a-z0-9-]`; today's line regex ignores it and the sweep would flag the entry as unattributed drift, unchanged by this move.

## Ownership fences
- `internal/skillsindex/`
- `internal/conformance/skills_index_checks_test.go` (including the SI6 guard test)
- `internal/conformance/checks_test.go` (the `frontmatterField` body only)
- `internal/conformance/docs_workflow_checks_test.go` (the `checkSkillsIndexGenerateVerify` call only)
- `internal/conformance/validity_checks_test.go` (the shell-syntax pattern list only)
- `internal/conformance/subcommand_routing_test.go`
- `internal/conformance/tier_test.go` (the `classifiedLiveTreeTests` entry only)
- `decisions/deepening-2026-08.md` (one superseded-note line under ticket #13's Answer)
- `cmd/bench/main.go`, `cmd/bench/main_test.go`, `cmd/bench/command_registry_test.go`
- `bin/bench.sh` (the case label and help line only)
- `.bench/BENCH.md`, `.bench/BENCH-reference.md`, `projects/benchkit.md`, `CHANGELOG.md`
- `.bench/skills-index.sh` (deletion)
- `specs/skills-index-reader/`
- `capture/retros/skills-index-reader.md`, `capture/agent-performance/` (final-check's retro and scorecards)

## Out of scope
- Real YAML frontmatter parsing and validation (FT89's remaining slice) — a separate capability; the module keeps the fence-and-prefix rule both readers use today. Roughly 4 edits, 3 gate runs once a YAML reader is chosen (dependency decision per AGENTS.md).
- Deriving other inventories (commands, plumbing list) from one implementation (FT89) — separate capability, ~6 edits, 3 gate runs.
- Any change to the index line format or to which skills are indexed.
