---
description: Configure this repo for Bench — run `bench setup`, then refine the gate, the project profile (seams + lines + design repo), and optionally seed CONTEXT.md. Run once per repo. This is one-time setup, not a workflow phase.
disable-model-invocation: true
---

# /bench-setup-repo — configure this repo for Bench

## Entry orientation

This is the one-time setup phase. It first runs `bench setup`, one command
that inspects the repo, previews every inferred fact, and converges the
managed assets. That command also scaffolds a proposed gate and a starter
profile. Then it fills in the parts that are specific to *this* repo and
that inference cannot supply: the real gate command, the profile's seams and lines, and
the domain language. It ends with a refined `.bench/gate.sh`, a refined
`projects/<name>.md` profile, and optionally `CONTEXT.md`.

## Exit handoff

Report what setup wrote or confirmed, whether the gate ran green, and the
configured profile path. Recommend `/bench-shape-idea` when unresolved product
fog remains. Recommend `/bench-write-spec` when the first build is already clear.
Recommend `/bench-debug` when a concrete bug prompted setup.

## 0. Run `bench setup`

Run `bench setup` first. It is one command. It inspects the repo, prints a plan
preview of every inferred fact with its consequence, and asks only genuine
ambiguities. It then transactionally converges the managed assets, `AGENTS.md` /
`CLAUDE.md`, a proposed gate, and a starter profile. It closes by running
`bench doctor` and naming the exact next action.

- If `bench setup` names this command (`/bench-setup-repo`) as the next action,
  pick up from here with the repo-specific interview below.
- If it reports a conflict, stop and surface it — ownership resolution is the
  reviewer's call, not something to route around.
- If it names a different next action (a gate-configuration step for a
  zero-signal repo, for instance), do that first. Then return here for the
  judgment content.

From here on, this command is prompt-driven, not a blind script: explore,
present what you found, confirm with me, then write. Assume I might not
remember what a term means. Explain each one briefly before you ask.

## 1. Explore

Read the repo first. Do not assume. Inspect, quietly:

- `git remote -v` — is there a remote? GitHub, GitLab, or local-only?
- the stack: `package.json` / `pyproject.toml` / `Cargo.toml` / `go.mod` — what
  language, test runner, type checker, linter are actually present?
- `.bench/gate.sh` — what did `bench setup` propose or scaffold as the stub?
- `projects/<name>.md` — what did `bench setup` seed as the starter profile?
- **existing Pocock structure** — `CONTEXT.md`, `docs/adr/`, `docs/agents/` (an
  `issue-tracker.md` / `domain.md` from `setup-matt-pocock-skills`). If present,
  this repo is migrating from Pocock's skills: reuse all of it. Read `CONTEXT.md`
  and the ADRs as the existing mental model. Treat `docs/agents/*` as already-made
  decisions and do not re-ask what they answer. Bench only adds the gate, the
  profile, and the lines on top.
- what kind of project this is — a UI app, an agent-facing CLI, a library/service —
  since that determines which project-specific oracles the gate needs
- any UI surface, and whether a separate design source (repo, package, path) exists

## 2. Present findings, then ask one decision at a time

Summarize what's present and what's missing. Then walk the three sections **one at
a time**: present a section, get my answer, and move on. Do not dump all three at
once.

### Section A — the gate (the oracle)

> The gate is the only thing that can call a shift "done." It's a script at
> `.bench/gate.sh` that exits zero when the repo is in a shippable state.
> Everything else in Bench enforces it. Nothing overrides it. This is the
> load-bearing choice — if the gate is weak, the whole system is weak.

`bench setup` already proposed the gate command from a small detection table,
or on zero signal scaffolded a fail-closed stub for it. Confirm it is right,
then refine it. The real command is usually richer than the inferred one, e.g.:

- Python: `mypy <pkg> && pytest -q && ruff check <pkg>`
- Node: `pnpm -s typecheck && pnpm -s test && pnpm -s lint`

Then ask about project-specific oracles, by what kind of project this is:

- **UI app** → add a `design-conformance` check that lints the UI against the
  project's design source: no raw hex or hardcoded spacing, and components
  composed from the canonical inventory rather than restyled. Note a
  screenshot/visual loop as part of UI done-ness, if there is one.
- **Agent-facing CLI** → add a conformance check against its output contract
  (e.g. structured output, stable schemas, correct exit codes). If it wraps
  another tool, add a regression check against that baseline too.
- **Library / service** → usually just the stack checks, plus any invariant I
  name (a public-API snapshot test, a perf budget).

Confirm whether each check actually runs today. A check may name a contract I
have not built yet; the conformance ones often do. Write that check into the
gate **commented out** with a `TODO`, so the gate stays green until I implement
it. Do not ship a gate that is red by construction.

Write every check the `craft-gate` skill's way. Prove it bites. Attribute its
failure. Run the real path. Choose its fail posture out loud.

### Section B — the profile (seams + lines + design source)

> The profile at `projects/<name>.md` is what a cold session reads to learn this
> repo's shape: the **seams** (where tests attach and where you should compose
> rather than invent), the **lines** (default model + effort per kind of task),
> the gate command, and — for UI repos — the design-source path. It's how a fresh
> agent skips the archaeology.

Ask for these, with defaults proposed from exploration:

- the 3–5 major seams — infer candidates from the codebase's module structure and
  propose them; a seam is wherever a stable interface already separates concerns
- the domain hostile-input checklist — propose 5–8 edge classes from the stack,
  seeded from the matching sections of the kit's hostile-input library
  (`.agents/skills/bench-craft-seams/references/hostile-input-library.md`:
  shell CLI, HTTP API, web UI, background jobs), tuned to this project. It
  lives in the profile. `/bench-write-spec`'s edge inventory reads it when it
  maps acceptance coverage
- the design-source location if there's UI (submodule / package / path)
- the line defaults, written as a **harness × tier binding**:
  - First discover candidate ids, but do not let a harness assign the tiers.
    Run `bench models` for multi-source advisory discovery (Codex catalog,
    OpenAI API, Anthropic API, and unavailable/manual hints). Read the
    harness's own model list too, when useful (`claude --help` / the
    `/model` picker, `codex debug models --bundled`).
  - Ask me, for each harness this repo actually runs, which opaque safe
    model-id tokens bind **cheap / mid / top**. The tier is the only
    identity the harnesses share, so a harness nobody runs here leaves its
    column unbound rather than borrowing another's ids.
  - Confirm the binding and record it with the date. Write it in the
    profile's `Lines` prose as a harness × tier table, and machine-readably
    in `.bench/lines.env`, one `BENCH_<HARNESS>_<TIER>` key per cell.
  - The hooks and shift adapters enforce the line only through `lines.env`.
    Discovery is not validation, and a repo without `lines.env` stays
    unrouted.
  - Then route by tier: cheap plus low effort for plumbing at a known seam,
    top plus high effort only for the genuinely uncertain seam.

Refine the starter `projects/<name>.md` that `bench setup` seeded, using the
example profiles in the kit as a template for anything it left blank.

### Section C — domain language (CONTEXT.md), optional

> `CONTEXT.md` at the repo root is the ubiquitous-language list a cold session
> reads first. It stops the agent from inventing synonyms for your domain:
> pick one term per concept and list the ones to avoid. Optional, but it pays
> for itself fast.

Offer to seed it now, with a short `craft-grill` pass over the core nouns, or
leave it for later. `/bench-shape-idea` and `/bench-write-spec` can create it
lazily when terms resolve. Do not force it.

## 3. Confirm, then write

Show me drafts of `.bench/gate.sh`, `projects/<name>.md`, and (if chosen)
`CONTEXT.md` before writing. Let me edit. Then write them, make
`.bench/gate.sh` executable, and verify it runs (`bench gate`). If it errors
for a reason other than real failing checks, fix the wiring before you declare
done.

## 4. Done

Tell me what's now configured. Tell me that the working commands
(`/bench-shape-idea`, `/bench-write-spec`, `/bench-debug`,
`/bench-implement-spec`, `/bench-review-implementation`, `/bench-final-check`)
and `bench shift` will read from these files. If you plan headless runs,
note that `bench shift` needs `BENCH_AGENT` pointed at an adapter in
`.bench/adapters/`; an env var setup does not write it. Note that I can edit
`.bench/gate.sh` and the profile directly later. `bench setup` is safe to
re-run; it converges and reports rather than start over. A re-run of
`/bench-setup-repo` revisits the judgment content; it does not undo
existing edits.
