# non-claude-model-tiers - owner-defined tiers across harnesses

Status: implemented

Source: `decisions/non-claude-model-tiers.md`.

## Problem

Bench already routes work through cheap, mid, and top tiers, but the live
validation and discovery surfaces still read as Claude-shaped. A reviewer who
wants Codex, OpenAI, or provider-qualified ids cannot bind those ids in the same
owner-controlled surface without fighting the gate or letting the harness invent
its own tiers.

## Solution

Keep `.bench/lines.env` as the authoritative reviewer-owned binding, widen tier
values from `claude-*` to opaque safe model-id tokens, and make `bench models` a
multi-source advisory inventory. Discovery may show candidate ids and source
freshness, but it never assigns cheap, mid, or top and never validates the final
binding.

## User stories

1. As a reviewer, I want `.bench/lines.env` to accept safe non-Claude model ids
   and reject unsafe tier strings, so I can bind Codex, OpenAI, or provider
   qualified models without weakening shell safety. Line: claude-opus-4-8 /
   medium. This touches the oracle's structural validation at a known seam, and a
   wrong gate is expensive even when the grammar is precise.

2. As a routed agent, I want the Agent hook and the headless adapters to allow a
   bound non-Claude `BENCH_MODEL` and deny an unbound one by exact membership, so
   provider namespace does not matter once the reviewer has chosen the tiers.
   Line: claude-opus-4-8 / medium. The behavior is gate-observable, but it spans
   shared enforcement surfaces where accidental Claude coupling is the failure
   mode.

3. As a reviewer choosing tier ids, I want `bench models` to list candidate ids
   from Codex, OpenAI, and Anthropic with freshness or unavailable status, so I
   can make an informed owner decision without discovery becoming enforcement.
   Line: claude-opus-4-8 / medium. The command has a known CLI seam, but it
   coordinates subprocess, environment, HTTP, parsing, and output-posture paths.

4. As an agent reading setup and line guidance, I want Bench prose to say tier
   values are owner-defined opaque tokens and `bench models` is advisory, so a
   fresh session does not infer that Claude is the tier namespace. Line:
   claude-fable-5 / high. This is project guidance prose, and the profile's
   leverage override applies because a wrong instruction multiplies through
   future sessions while the gate can only check structure.

5. As a kit maintainer, I want canary and conformance fixtures that embed changed
   guidance to move with the source text, so stale-reference checks still prove
   the current source bites. Line: claude-sonnet-5 / low. Fixture sync is
   mechanical and the canary gate observes whether it stayed synchronized.

## Implementation decisions

- **Tier values are safe opaque tokens.** `BENCH_TIER_TOP`,
  `BENCH_TIER_MID`, and `BENCH_TIER_CHEAP` accept non-empty tokens matching
  `^[A-Za-z0-9][A-Za-z0-9._:/@+-]*$`. That admits examples such as `gpt-5.4`,
  `gpt-5.3-codex-spark`, and `openai/gpt-5`, while rejecting whitespace,
  control bytes, quotes, dollar expansion, command separators, pipes,
  ampersands, parens, globs, and tilde expansion. The gate does not ask any
  provider whether the token exists.

- **Aliases remain harness-specific bare names.** `BENCH_ALIAS_*` keeps the
  existing bare-alias contract for harnesses such as Claude Code whose Agent
  tool accepts aliases instead of exact model ids. Alias validation is not
  widened into the model-id grammar.

- **`internal/lines` stays the membership engine.** Parsing `.bench/lines.env`
  and answering "is this exact model id declared" remains provider-agnostic and
  should not learn provider syntax. Structural validation of the binding belongs
  to conformance because the map assigned it to the gate.

- **`bench models` is an advisory inventory, not a validator.** It exits 0 when
  any source is unavailable or malformed, reports that source as unavailable, and
  continues with the remaining sources. An id absent from discovery can still be
  a valid `.bench/lines.env` value when it matches the safe-token grammar.

- **Discovery sources have explicit freshness.** Codex discovery runs
  `codex debug models` when `codex` is on `PATH`, parses visible
  `models[].slug` values, and falls back to `codex debug models --bundled` when
  refresh fails. OpenAI discovery uses `OPENAI_API_KEY` with `GET /v1/models` and
  parses `data[].id`. Anthropic discovery keeps the existing `ANTHROPIC_API_KEY`
  and `/v1/models` behavior and parses `data[].id`. Live successful sources are
  labeled `live`, Codex bundled fallback is labeled `bundled`, and missing or
  failed sources are labeled `unavailable`.

- **Inventory output is structured for agents.** The command should emit a flat
  table of source rows and a flat table of model rows, for example
  `model_sources[N]{source,freshness,status,hint}:` and
  `models[N]{source,freshness,id}:`. Exact wording may follow the local TOON
  emitter, but tests pin the row fields, source names, freshness labels, and exit
  code.

- **Current Claude ids may remain this repo's chosen binding.** The profile and
  comments around `.bench/lines.env` can keep today's Claude values if the
  reviewer keeps that line, but they must describe those values as the current
  binding, not the allowed namespace. The profile drift check still compares the
  literal bound values against `.bench/lines.env`.

## Testing decisions

- A good test here drives the public seam and observes behavior the reviewer
  cares about: conformance diagnostics, hook or adapter allow-deny behavior,
  `bench models` stdout plus exit code, and stale-reference canaries. Parser tests
  are acceptable only where they front an uncontrollable source such as provider
  JSON or a Codex subprocess response.

- The prior art is the existing line-routing conformance family, the
  `internal/lines` exact-membership tests, the current `internal/models` parser
  tests, and the canary fixtures described in the project profile.

- Gate: the project gate, `bench gate`.

### Seam diagram

    trigger: `bench gate` runs line-routing conformance
        |
        v
    .bench/lines.env tokens  --->  [ line binding validation ]  --->  diagnostics or green
    projects/<name>.md prose  --->  [                         ]  --->  profile drift result
                         ^ tests attach here: conformance fixtures with accepted
                           non-Claude ids and rejected unsafe tier tokens

    trigger: Agent hook or adapter launch in a routed repo
        |
        v
    BENCH_MODEL / resolvedModel  --->  [ internal/lines verdict surface ]  --->  allow or deny
    .bench/lines.env            --->  [                                ]  --->  selected model id
                              ^ tests attach here: hook and adapter conformance
                                probes using at least one Codex-shaped binding

    trigger: reviewer runs `bench models`
        |
        v
    codex debug models JSON     --->  [ models advisory inventory ]  --->  source rows + model rows
    OpenAI /v1/models JSON      --->  [                           ]  --->  exit 0
    Anthropic /v1/models JSON   --->  [                           ]  --->  unavailable rows
                              ^ tests attach here: command/parser tests stub
                                subprocess, env, HTTP, malformed JSON, and output

    trigger: `bench gate` runs docs and canary conformance
        |
        v
    command / skill / profile text  --->  [ stale-reference and canary checks ]  --->  green or targeted red
                                  ^ tests attach here: existing fixture drift
                                    checks after changed source text is mirrored

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Safe non-Claude tier values such as `gpt-5.4`, `gpt-5.3-codex-spark`, and `openai/gpt-5` are accepted without diagnostics | line binding validation | conformance fixture with those three `BENCH_TIER_*` values is red today because the validator requires `claude-*` | rejects accidental provider-prefix validation while still exercising the real gate seam |
| 1 | Unsafe tier values are rejected with the offending key and value named | line binding validation | table-driven conformance cases for empty, whitespace, ESC, dollar expansion, backtick, semicolon, pipe, ampersand, glob, quote, and tilde tokens are red until the safe-token helper exists | prevents widening the grammar into a shell-unsafe free string and makes diagnostics actionable |
| 2 | The Agent hook allows a non-Claude bound resolved model and denies an unbound non-Claude resolved model | Agent hook line-routing conformance | routed fixture changes one tier to `gpt-5.4`; the allow case is red today under the Claude-only validator, while the deny case must stay red for an undeclared safe id | proves the hook reads owner membership from `lines.env` instead of assuming Claude names |
| 2 | Each shipped adapter accepts a bound non-Claude `BENCH_MODEL` and refuses an unbound non-Claude `BENCH_MODEL` in a routed repo | adapter line-routing conformance | routed adapter fixture uses `BENCH_MODEL=gpt-5.3-codex-spark`; the valid case is red today under the Claude-only gate, and the existing unbound denial remains the paired guard | catches adapter paths that bypass the shared verdict or silently reintroduce provider assumptions |
| 3 | `bench models` emits Codex visible slugs with source and freshness, and falls back to bundled freshness when refresh fails | models command stdout and exit code | command tests stub `codex debug models` JSON and a failing refresh plus `--bundled`; both are red today because the command ignores Codex | pins local Codex discovery and the required no-network fallback without needing a live Codex service |
| 3 | `bench models` emits OpenAI and Anthropic `data[].id` ids when keys are present, and unavailable source rows when keys are absent | models command stdout and exit code | fake HTTP and env tests assert OpenAI, Anthropic, and missing-key rows; OpenAI and unavailable rows are red today because output is Anthropic-only no-key prose | proves multi-source inventory and the exit-0 unavailable posture together |
| 3 | Malformed provider JSON, a failing Codex command, or a non-2xx live response produces an unavailable source row and exit 0 | models command stdout and exit code | stub malformed JSON, failed subprocess, and non-2xx response cases; red today because there is no per-source unavailable row contract | prevents discovery failure from becoming enforcement or a hard gate on tier ownership |
| edge of 3 | A provider model id containing control bytes or whitespace is not emitted as a model row | models command output renderer | parser or renderer test feeds a source id with ESC or whitespace; red until invalid ids are skipped or downgraded to an unavailable-source hint | keeps the structured output safe even when a live provider returns an unrenderable id |
| 4 | Setup, line, profile, and lines-env prose describe tier values as owner-defined opaque tokens and `bench models` as advisory discovery | docs conformance plus review | not fully TDD-able: existing stale-reference and anchor checks catch mirrored text drift, while semantic wording needs review against this spec | prose semantics are the behavior, so the spec gives the reviewer a precise sentence-level target |
| 5 | Canary and embedded-doc fixtures that quote changed guidance are updated with the source text and still fail for their planted reason | canary and docs conformance | `bench gate` canary run is red if a fixture still embeds stale text or if a stale-reference check stops biting | keeps the gate's defense against stale or weakened guidance intact |

### Edge inventory

- error path -> rows: unsafe tier diagnostics (1), unbound hook and adapter
  denial (2), unavailable model sources and malformed provider data (3), canary
  red posture (5).
- empty or absent input -> rows: empty tier values (1), missing API keys and
  missing `codex` command as unavailable source rows (3).
- boundary values -> rows: first-character and allowed-character token grammar
  cases in story 1; empty provider result sets remain a valid source row with no
  model rows.
- malformed input -> rows: control bytes, whitespace, and shell-dangerous tier
  values (1); malformed provider JSON and invalid provider ids (3).
- interrupted or partial state -> **Won't handle**: these commands and checks
  are read-only and do not create scratch state.
- re-run idempotency -> covered by the same read-only command and conformance
  seams; repeated `bench models` runs do not write files.
- hostile environment -> rows: missing `codex` command, missing API keys,
  failing subprocess, and non-2xx provider responses (3).
- paths and directory names with spaces or globs -> **Won't handle**: the
  feature does not introduce new path arguments, and existing root discovery
  owns invocation-path behavior.
- symlink invocation and cwd below the repo root -> **Won't handle**: no new
  root-resolution code is introduced, and existing CLI and conformance entry
  points still own that class.
- unquoted multi-word arguments -> covered by existing adapter prompt tests;
  this feature changes model id values only, and the safe-token grammar rejects
  spaces by design.

## Out of scope

- **Automatic tier assignment or `.bench/lines.env` write-back** - a separate
  reviewer-assistance capability that would need ranking, confirmation, and
  write-path safety; this feature is read-only discovery plus validation.
  Estimate: 6 edits, 3 gate runs.
- **Provider cost or quality ranking** - a separate recommendation system rather
  than owner-defined binding; it would need pricing data, freshness policy, and
  reviewer-visible uncertainty. Estimate: 7 edits, 3 gate runs.
- **Opencode-specific live discovery** - a separate provider integration once a
  local opencode catalog surface exists in this repo's environment; the current
  spec covers generic safe tokens and manual fallback. Estimate: 4 edits, 2 gate
  runs.
