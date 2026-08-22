# Split the always-loaded core from the reference and move the anchors

Blocked by: 01c-register-the-check-and-exclusions.md
Writes: .bench/BENCH.md, .bench/BENCH-reference.md, internal/anchors/registry_data.go, internal/conformance/docs_workflow_helpers_test.go, tests/canary/workflow-guidance-anchors/ (new fixture directories), projects/benchkit.md
Line: `opus` / medium under Claude Code, `gpt-5.6-terra` / medium under Codex — the disposition table is exact and every move is gate-graded.

## What to build

`.bench/BENCH.md` keeps the sections the spec's disposition table keeps and moves the rest to `.bench/BENCH-reference.md` as the table says. The `CLI Inventory` heading becomes `The Bench CLI contract`, and the contract paragraph ends with the pair: "`bench gate` is valid. `bench gate 2>&1 | tail -20` is not valid." The core keeps one sentence that points a reviewed spec-backed build at `bench worktree land` and at the reference. The core keeps every shared-rule marker, the string `BENCH-reference.md`, and the `bench help` inventory sentence. The `bench` bullet merges into the reference intro and keeps its `CONTEXT.md` and profile pointers.

The profile sentence that says the core keeps category-level operational guidance changes to name the reference, as does the reference's own `Command Notes` sentence.

The five moved registry rows change `File:` to the reference. The core's retained-integration-source row is deleted, and the integration-source anchor count moves from 12 to 11. One new row requires the pair in the core.

Seven new rows require the moved units that carried no anchor, one per unit, with the needles the spec names. Thirteen `Forbid` rows on the core prove absence, one per moved needle, with the needles the spec names. Two more `Forbid` rows, one per file, forbid "progressive loading". Each new or moved row gets one `files/`-form fixture. The two `.bench/` exclusion rows stay until ticket 02b.

## Acceptance

- [ ] The core holds exactly the kept sections, and every moved passage has a needle found at the reference (covers PD1).
- [ ] `bench anchors` on each file lists the moved rows at the reference and the kept rows at the core (covers PD2).
- [ ] Each new or moved row has a fixture that bites and clears (covers PD3).
- [ ] The core contains the pair sentence and a row requires it (covers PD4).
- [ ] No `## ` heading appears in both files (covers PD5).
- [ ] The eleven shared-rule markers stay in the core and in neither AGENTS.md nor README.md (covers PD6).
- [ ] The core contains `BENCH-reference.md` and CLAUDE.md does not import the reference (covers PD7).
- [ ] Both files carry a `Forbid` row for "progressive loading" with a fixture that bites (covers PD44).
