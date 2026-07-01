---
description: Configure this repo for Bench — fill in the gate, the project profile (seams + lines + design repo), and optionally seed CONTEXT.md. Run once per repo, after `bench link` and `bench init`. This is one-time setup, not a workflow phase.
disable-model-invocation: true
---

# /bench-setup-repo — configure this repo for Bench

`bench link` wired the kit in and `bench init` scaffolded an empty `.bench/gate.sh`.
This fills in the parts that are specific to *this* repo and can't be hardcoded:
the real gate, the profile, and the domain language. It's prompt-driven, not a
script — explore, present what you found, confirm with me, then write. Assume I
might not remember what a term means; explain each one briefly before asking.

## 1. Explore

Read the repo first; don't assume. Check, quietly:

- `git remote -v` — is there a remote? GitHub, GitLab, or local-only?
- the stack: `package.json` / `pyproject.toml` / `Cargo.toml` / `go.mod` — what
  language, test runner, type checker, linter are actually present?
- `.bench/gate.sh` — already filled, or still the scaffold stub?
- `projects/<name>.md` — does a profile already exist for this repo?
- **existing Pocock structure** — `CONTEXT.md`, `docs/adr/`, `docs/agents/` (an
  `issue-tracker.md` / `domain.md` from `setup-matt-pocock-skills`). If present,
  this repo is migrating from Pocock's skills: reuse all of it. Read `CONTEXT.md`
  and the ADRs as the existing mental model, treat `docs/agents/*` as already-made
  decisions, and don't re-ask what they answer. Bench only adds the gate, the
  profile, and the lines on top.
- what kind of project this is — a UI app, an agent-facing CLI, a library/service —
  since that determines which project-specific oracles the gate needs
- any UI surface, and whether a separate design source (repo, package, path) exists

## 2. Present findings, then ask one decision at a time

Summarize what's present and what's missing. Then walk the three sections **one at
a time** — present a section, get my answer, move on. Don't dump all three at once.

### Section A — the gate (the oracle)

> The gate is the only thing that can call a shift "done." It's a script at
> `.bench/gate.sh` that exits zero when the repo is in a shippable state. Everything
> else in Bench enforces it; nothing overrides it. This is the load-bearing
> choice — if the gate is weak, the whole system is weak.

Propose a gate command from what you found in the stack, e.g.:

- Python: `mypy <pkg> && pytest -q && ruff check <pkg>`
- Node: `pnpm -s typecheck && pnpm -s test && pnpm -s lint`

Then ask about project-specific oracles, by what kind of project this is:

- **UI app** → add a `design-conformance` check that lints the UI against the
  project's design source (no raw hex or hardcoded spacing, components composed from
  the canonical inventory rather than restyled). If there's a screenshot/visual loop,
  note it as part of UI done-ness.
- **Agent-facing CLI** → add a conformance check against its output contract (e.g.
  structured output, stable schemas, correct exit codes) and, if it wraps another
  tool, a regression check against that baseline.
- **Library / service** → usually just the stack checks, plus any invariant I name
  (a public-API snapshot test, a perf budget).

Confirm whether each check actually runs today. If a check is a contract I haven't
built yet (the conformance ones often are), write it into the gate **commented
out** with a `TODO`, so the gate stays green until I implement it — don't ship a
gate that's red by construction.

### Section B — the profile (seams + lines + design source)

> The profile at `projects/<name>.md` is what a cold session reads to learn this
> repo's shape: the **seams** (where tests attach and where you should compose
> rather than invent), the **lines** (default model + effort per kind of task),
> the gate command, and — for UI repos — the design-source path. It's how a fresh
> agent skips the archaeology.

Ask for, with defaults proposed from exploration:

- the 3–5 major seams — infer candidates from the codebase's module structure and
  propose them; the seams are wherever a stable interface already separates concerns
- the design-source location if there's UI (submodule / package / path)
- the line defaults, written as **tier → model bindings**. First discover what's
  actually available — don't hardcode model names. Run `bench models` (it queries
  the Anthropic model list when a key is set); if the harness uses subscription auth
  and there's no key, read the harness's own list (`claude --help` / the `/model`
  picker, `codex --help`). Bind **cheap / mid / top** to the cheapest, a middle, and
  the most capable available model, confirm with me, and record the binding with the
  date and harness. Then the routing by tier: cheap + low effort for plumbing at a
  known seam; top + high effort only for the genuinely uncertain seam.

Write `projects/<name>.md` from the example profiles in the kit as a template.

### Section C — domain language (CONTEXT.md), optional

> `CONTEXT.md` at the repo root is the ubiquitous-language list a cold session
> reads first — it stops the agent inventing synonyms for your domain (pick one
> term per concept and list the ones to avoid). Optional, but it pays for itself
> fast.

Offer to seed it now (a short `craft-grill` pass over the core nouns) or to leave it for
later — `/bench-shape-idea` and `/bench-write-spec` can create it lazily when terms get resolved. Don't force
it.

## 3. Confirm, then write

Show me drafts of `.bench/gate.sh`, `projects/<name>.md`, and (if chosen) `CONTEXT.md`
before writing. Let me edit. Then write them, make `.bench/gate.sh` executable, and
verify it runs (`bench gate`) — if it errors for a reason other than real failing
checks, fix the wiring before declaring done.

## 4. Done

Tell me what's now configured and that the working commands (`/bench-shape-idea`, `/bench-write-spec`,
`/bench-debug`, `/bench-implement-spec`, `/bench-review-implementation`, `/bench-final-check`) and `bench shift` will read from these
files. Note that I can edit `.bench/gate.sh` and the profile directly later — re-running
`/bench-setup-repo` is only for starting over.
