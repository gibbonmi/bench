# Multi-harness line binding

Status: ready

## Destination

Make the tier→model binding symmetric across harnesses so each harness reads,
reports, and recommends its own model family. Today `BENCH_TIER_*` holds Codex
gpt ids as the privileged canonical and `BENCH_ALIAS_*` is a Claude-only
translation, so every advisory surface — the SessionStart line report, the
profile "Lines" section, and the agent's own recommendations — leads with gpt
ids even inside a Claude Code session. After this change no family is canonical:
the tier (top/mid/cheap) is the only abstract identity, and every reader resolves
it to the current harness's family. Enforcement already runs the right family per
harness (adapters + the Agent-line hook); this is the *advisory* layer catching up.

Scope confirmed with reviewer: the bite is presentation only — what the agent
recommends, the profile framing, and the SessionStart report. No run is
mis-launching today.

## #1: Data model — one canonical family, or symmetric peers?

Blocked by: none
Type: Grill

### Question
`BENCH_TIER_*` is named for the abstract tier but holds Codex ids, structurally
privileging Codex as "the real model." Fix the presentation only, or restructure
the data model so no family is canonical?

### Answer
**Symmetric peers.** The tier (top/mid/cheap) becomes the only canonical identity;
every harness family is a peer binding, none privileged. Bigger blast radius
(lines.env schema, Go core, conformance check, all three adapters) but it's the
only shape that stops a family being "the real id." Presentation-only was rejected
because prose would still call gpt "the tier."

## #2: lines.env schema — key shape and harness set

Blocked by: none
Type: Grill

### Question
How are the per-harness-per-tier cells keyed, and is the harness set open or fixed?

### Answer
**`BENCH_<HARNESS>_<TIER>`, fixed validated set.** e.g. `BENCH_CODEX_TOP`,
`BENCH_CLAUDE_MID`. Groups by harness. The core validates against a known harness
list (codex/claude/opencode); adding a harness is a deliberate code change, not a
free-form key. Chosen over `BENCH_TIER_<TIER>_<HARNESS>` open-set for tighter
validation.

## #3: How does the resolver/hook learn the current harness?

Blocked by: none
Type: Grill

### Question
Today the harness is implied by output-shape flags (`--alias`→claude,
`--provider-model`→opencode, bare→codex). Keep that, or make harness explicit?

### Answer
**Explicit `--harness <name>`.** Each adapter and the Agent-line hook passes
`--harness codex|claude|opencode`; output shape follows from the harness's column.
The old shape-flags (`--alias`, `--provider-model`) are retired — they conflated
"output shape" with "which harness," the exact conflation being removed.

## #4: Agent-line enforcement scope — own-family-only or any bound tier?

Blocked by: none
Type: Grill

### Question
In Claude Code, a delegation names e.g. `opus`. Should the hook (`--harness claude`)
accept only `BENCH_CLAUDE_*` values, or any bound cell across the matrix?

### Answer
**Any bound tier value.** Enforcement stays permissive — the honest-mistake guard
keeps today's "is it on a bound tier" semantics, now matrix-wide. A model matching
no cell still denies (exit 2), and the deny message uses `--harness` to recommend
the *current* harness's family. **Steering to the native family is advisory, not
enforced** — realized in recommendation/report/skill (#8), not the verdict.

## #5: Docs + conformance for a 3×3 matrix

Blocked by: none
Type: Grill

### Question
Nine concrete ids. Keep listing them in the profile prose (cross-checked), or defer
to lines.env as the one source?

### Answer
**Full matrix in prose, cross-checked.** The profile "Lines" renders the matrix as
a table; `checkLineBinding` cross-checks every cell against lines.env, extending
today's prose↔env check. Completeness rule: for any *declared* harness, all three
tiers must be present and be safe tokens. Deferring to lines.env was rejected —
reviewer wants the binding human-readable in the profile.

## #6: Which harnesses does this repo bind now?

Blocked by: none
Type: Grill

### Question
The schema knows three harnesses. Populate all three, or only those in use?

### Answer
**codex + claude only.** OpenCode stays unbound until adopted; its adapter refuses
to launch (fail-closed) until `BENCH_OPENCODE_*` is added, and `bench doctor` says
so. The completeness check requires all tiers only for *declared* harnesses, so an
absent opencode column is not a gate failure. Don't invent opencode ids not being run.

## #7: SessionStart line report shape

Blocked by: none
Type: Grill

### Question
`check-agent-line --describe` is where "claude code looks to use gpt models" bit.
With `--harness claude`, what does it print?

### Answer
**Current-harness column primary.** The report leads with the asking harness's
family (Claude: `top=fable mid=opus cheap=sonnet`), full matrix behind a verbose
flag. Directly fixes the bite — a Claude session sees Claude models first, no gpt
row in the default report.

## #8: BENCH_MODEL contract for a shift

Blocked by: none
Type: Grill

### Question
Adapters now pass `--harness`. A concrete-id `BENCH_MODEL` + `--harness claude`
is incoherent (which family wins?). What does BENCH_MODEL name?

### Answer
**A tier token (top/mid/cheap).** `BENCH_MODEL=cheap bench shift …` → codex adapter
resolves `gpt-5.6-luna`, claude adapter resolves `sonnet`. The reviewer picks a
tier; the harness picks the family. Fully realizes "tier is the abstract identity"
and removes any family/harness mismatch at launch. The opencode `--provider-model`
projection is retired — opencode gets explicit `BENCH_OPENCODE_*` ids (full
provider/model strings) when bound, not a projection off a canonical id.

## Migration (resolved under #1)

Hard cut, no dual-read shim. `bench doctor` detects the retired
`BENCH_TIER_*`/`BENCH_ALIAS_*` keys and reports the exact rewrite to
`BENCH_CODEX_*`/`BENCH_CLAUDE_*`. The core reads only the new keys. (Auto-rewrite
vs report-and-offer for a reviewer-owned file: spec-writer's call — lean
report-and-offer, see Handoff item 9.)

## The recommendation rule (resolved under #4/#7)

The actual ask — "recommend claude models in claude code." Since enforcement is
permissive, this lives in the advisory layer: `craft-line` gains a rule to declare
the line / recommend in the *current harness's* family; the SessionStart report
(#7) reinforces it; the tier memory is updated. No family is named in the abstract
rule — it resolves per harness.

## Not yet specified

## Spec-writer discretion

- Exact deny-message wording and the verbose-flag name for the full-matrix report.
- Whether `bench models` discovery output should also become harness-scoped.

## Sources

## Out of scope
