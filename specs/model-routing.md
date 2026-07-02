# Model + effort routing (the line rubric and its enforcement)

Decision map: `decisions/model-routing.md` (closed — all 9 tickets answered).

## Problem

Invariant #2 requires declaring a line (model / effort / cap) before long runs,
but nothing tells a session *how* to pick the tier, and nothing checks that a
delegation actually carries the declared model. Routing is ad hoc judgment;
headless `bench shift` runs silently use whatever model the harness defaults to.

## Solution

A `craft-line` skill encodes a three-signal decision table that picks the
*starting* tier+effort, and an escalation ladder that corrects it on gate
feedback. Tier→model bindings become machine-readable in `.bench/lines.env`,
consumed by two enforcement surfaces: a Claude Code PreToolUse hook on the
`Agent` tool (denies delegations to unbound models, with a reason the model
self-corrects on), and the `BENCH_AGENT` adapters (map `BENCH_MODEL` to the
harness's model flag; refuse undeclared headless delegation in routed repos).
Effort is enforceable on no surface; it stays declaration discipline carried by
the skill.

## User stories

1. As the reviewer, I want a `craft-line` skill holding the decision table —
   spec precision × seam uncertainty × gate coverage → tier + effort as a joint
   output, with weak gate coverage bumping one tier — so every session routes by
   the same signals.
2. As the reviewer, I want the skill to encode the escalation ladder: first red
   retries same-tier with the gate output as added guidance; second red at the
   same tier escalates one tier; delegate-reported seam uncertainty escalates
   immediately; **any bump to the top tier pauses and asks me** unless the
   project's `Lines` section grants a standing opt-out; every move is reported
   in one line.
3. As the reviewer, I want the skill to state how tiers resolve: the abstract
   cheap/mid/top roles bind per project via `.bench/lines.env`, narrated in
   `projects/<name>.md` `Lines`.
4. As the reviewer, I want `.bench/lines.env` in this repo with the decided
   rotation — top=`claude-fable-5`, mid=`claude-opus-4-8`, cheap=`claude-sonnet-4-6`
   — so the binding the map decided is what enforcement reads.
5. As the reviewer, I want `projects/benchkit.md` `Lines` rewritten to the new
   rotation, citing `.bench/lines.env` as the binding source and `craft-line` as
   the rubric, keeping the cached per-work-type routings, recording that no
   top-tier opt-out is granted, and keeping the Sonnet 5 revisit note.
6. As the reviewer, I want `/bench-implement-spec`'s "Open with the line" step to
   reference `craft-line` so the declaration is rubric-driven, not restated prose.
7. As the orchestrating agent, I want a PreToolUse hook on the `Agent` matcher
   that reads the call's resolved model from stdin and denies with a reason when
   it isn't a bound tier model, so I self-correct without interrupting the
   reviewer.
8. As a user of an unrouted repo, I want that hook to allow everything when
   `.bench/lines.env` is absent or unparseable (warn on stderr), so enforcement
   arms only where a binding exists and a broken binding never bricks delegation.
9. As a shift operator, I want the three reference adapters to pass
   `BENCH_MODEL` to the harness's model flag when set, to validate it against
   `.bench/lines.env` membership when the file exists, and to refuse (nonzero,
   clear error) when the file exists but `BENCH_MODEL` is unset — so headless
   delegation always carries an explicit, bound line in routed repos.
10. As a user of an unrouted repo, I want adapters without `.bench/lines.env`
    to behave exactly as today (pass-through, no flag) — and an explicitly set
    `BENCH_MODEL` still passes through.
11. As the reviewer, I want gate checks with canary fixtures covering: the
    settings.json `Agent`-matcher wiring, the hook's allow/deny/no-op behavior,
    `.bench/lines.env` validity (three tiers, model-id shaped values), the
    adapter mapping/refusal behavior, and the `/bench-implement-spec` skill
    reference — so drift turns the gate red.
12. As a non-Claude harness user, I want the AGENTS.md skills index to gain the
    `craft-line` row (trigger: declaring a line / choosing a delegate's model)
    so manual-load harnesses find the rubric.
13. As a kit consumer, I want BENCH.md's adapter-contract section to document
    `BENCH_MODEL` (and the lines.env arming rule) so the contract stays the
    single description of how shift delegates.

## Implementation decisions

- **Single machine-readable binding:** `.bench/lines.env` with exactly
  `BENCH_TIER_TOP`, `BENCH_TIER_MID`, `BENCH_TIER_CHEAP` (exact model ids).
  Sourced by the hook and adapters; prose in `Lines` cites it. Per-repo
  authored (here, by this build; elsewhere by `/bench-setup-repo` later — see
  Out of scope).
- **Hook:** new `.bench/hooks/check-agent-line.sh`, wired in
  `.claude/settings.json` under PreToolUse matcher `"Agent"` (Claude Code only;
  Codex has no Agent tool — its delegation flows through adapters). Reads stdin
  JSON, compares the resolved model field against the three bindings by exact
  string. Deny = JSON permissionDecision deny + reason naming the bound models.
  **Fail-open** on missing file, malformed JSON, or missing model field
  (stderr warning) — availability beats strictness; the gate, not the hook, is
  the oracle.
- **Adapters stay one-liners plus a guard block**; `BENCH_MODEL` maps to
  `claude -p --model`, `codex exec -m`, `opencode run --model` (exact flags
  verified against each CLI during build; a harness with no model flag fails
  the guard rather than silently dropping the line). `BENCH_EFFORT` is **not**
  added — no harness exposes an effort flag; a decorative var repeats the
  `BENCH_MAX_TOKENS` mistake this repo already removed.
- **Escalation approval is guidance, not hook logic** — the hook can't tell an
  escalation from a first pick; the ask-before-top rule lives in the skill and
  the `Lines` opt-out.
- **Codex `.codex/hooks.json` unchanged** — nothing to enforce there.
- No transcript grepping for declared lines (brittle; the deterministic check
  is binding membership).

## Testing decisions

- A good test here exercises the enforcement scripts' external behavior — JSON
  stdin in, decision/exit code out — at the gate-check seam, matching how this
  repo tests everything (gate checks + `tests/canary/` fixtures that prove each
  check bites).
- Seams: (1) the gate check invoking `check-agent-line.sh` with fixture stdin;
  (2) the gate check invoking each adapter with a stubbed harness binary on
  PATH (asserting argv / refusal); (3) gate structural checks on settings.json
  wiring, lines.env shape, and the command-prose anchor. Prior art:
  `codex-hooks-broken`, `command-handoff-anchor`, `dangling-index` canaries.
- Gate command: `bench gate` (`.bench/gate.sh`).

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 7 | hook denies stdin with unbound model, reason names bound models | gate check: hook + deny fixture stdin | check fails before hook exists (no file / wrong output) | deny is the enforcement; absent or wrong output = no gate |
| 7 | hook allows stdin with a bound model | gate check: hook + allow fixture stdin | same check, allow assertion red pre-build | an over-eager hook that denies bound models bricks routing |
| 8 | hook allows + warns when lines.env absent | gate check: hook run in temp dir without lines.env | assertion red pre-build | fail-closed here would brick every unrouted repo |
| 8 (edge) | hook allows on malformed stdin JSON / missing model field | gate check: hook + malformed fixtures | assertion red pre-build | malformed input must not turn into a denial |
| 7 (edge) | hook allows a declared alias, denies an undeclared one | gate check: hook + alias fixtures | observed live: hook denied `opus` and broke in-session delegation | the Agent tool only speaks aliases; without this row routing blocks itself |
| 4, 11 | lines.env exists, parses, three tiers with model-id-shaped values | gate check; canary `lines-env-broken` (planted bad fixture goes red) | gate red on fixture | binding drift silently disarms both surfaces |
| 7, 11 | settings.json wires `Agent` matcher → check-agent-line.sh | gate check; canary `agent-hook-unwired` | gate red on fixture | hook file without wiring is decorative |
| 9 | claude adapter passes `--model "$BENCH_MODEL"` to stubbed binary | gate check: run adapter with PATH stub echoing argv | assertion red pre-build | silent flag-drop = undeclared delegation resumes |
| 9 | adapter refuses when lines.env exists and BENCH_MODEL unset | gate check: stubbed run, assert nonzero + message | assertion red pre-build | the refusal is the headless backstop |
| 9 (edge) | adapter refuses BENCH_MODEL not in lines.env | gate check: stubbed run with unbound id | assertion red pre-build | Codex/OpenCode have no hook layer; membership must check here |
| 10 | adapter without lines.env passes through unchanged | gate check: stubbed run in temp dir | assertion red pre-build | regression risk to every existing bench-shift user |
| 6 | implement-spec command anchors a craft-line reference | gate anchor check; canary `line-anchor-missing` | gate red on fixture | prose references rot invisibly without an anchor |
| 12 | AGENTS.md index row ↔ skill on disk | already covered — existing index-sync check + `dangling-index` canary | already covered | existing check |
| 1–3 | skill content (table, ladder, ask-before-top) | not TDD-able — prose guidance; frontmatter/index gate-checked, content reviewed by `/bench-review-implementation` | n/a | n/a |
| 13 | BENCH.md adapter-contract mentions BENCH_MODEL | gate anchor check (same style as story 6) | gate red pre-build | contract doc drifting from adapter behavior |

### Edge inventory

Walked per behavior (error path, empty, boundary, malformed, partial state,
re-run, hostile env):

- Hook: malformed JSON, missing model field, absent lines.env → coverage rows
  above (all fail-open).
- Hook: lines.env present but a tier var empty → treated as unparseable →
  fail-open row; gate check makes the file itself red so it can't persist.
- Adapter: BENCH_MODEL set, lines.env absent → pass-through row (story 10).
- Adapter: model id needing quoting → quoted expansion asserted by the argv
  stub row.
- Re-run idempotency: hooks/adapters are stateless per call — nothing to walk.
- Hook: the Agent tool addresses models by alias only (`opus`, `fable`), so
  `lines.env` may declare `BENCH_ALIAS_TOP/MID/CHEAP`; a declared alias allows,
  an undeclared alias (deliberately: bare `sonnet` → excluded Sonnet 5) denies
  → coverage rows below. Adapters stay id-only — headless runs can pass full
  ids, and Claude Code aliases would mistranslate to other harnesses.
  (Discovered live: the original "exact ids only" cut denied every legitimate
  in-session delegation.)
- **Won't handle:** OpenCode hook layer (none exists; the adapter guard is its
  enforcement).
- **Won't handle:** effort enforcement anywhere (no surface exposes it; skill
  discipline only).
- **Won't handle:** in-conversation "was a line declared" transcript check
  (brittle; membership check is the deterministic core).

## Out of scope

- **`bench models` writes/refreshes `.bench/lines.env`** — its own small CLI
  feature (today it only lists models); ~40 min agent time.
- **`/bench-setup-repo` interview authors lines.env for linked repos** — setup-phase
  capability, touches the interview script; ~30 min.
- **Escalation audit log** (append tier moves to a file for routing-quality
  review) — separate observability capability; ~1 h.
- **Sonnet 5 mid-seat rebinding** — parked on ROADMAP.md with a 2026-09-01
  trigger; a `Lines`+`lines.env` edit when taken.
