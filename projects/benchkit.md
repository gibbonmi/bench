# Project: benchkit

The Bench kit itself — a harness-agnostic agent-development workflow shipped as the
npm package `benchkit`. It is not an application: it is shell + markdown + JSON that
other repos consume. The deliverable is the `bench` CLI (`bin/bench.sh`), the working
agreement (`AGENTS.md`), the portable `.agents/` skills and commands, and harness
adapters that call shared `.bench/hooks/` scripts. Because the artifacts are plain
files, the kit must work identically under Claude Code, Codex, and any other
AGENTS.md harness — that portability is the product.

## Seams (test here; everything else is free to change)

- **The gate contract** (`.bench/gate.sh` / `bench gate`). The oracle surface. Everything
  in Bench routes through its exit code: 0 = shippable, non-zero = not done. The
  highest seam — if the gate is weak, the whole system is weak. Test by feeding it a
  conformant tree (green) and a broken one (red); never by trusting a reading of the
  diff.
- **The `bench` CLI subcommands** (`gate`, `worktree`, `shift`, `init`, `link`,
  `models`, `structure`, `idea`, `roadmap`, `status`). The operational shell surface.
  Stable command names and exit codes are the contract; the implementation behind each is
  free to change. Keep gate resolution (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) in
  one place.
- **The ambient dashboard** (`bench status`). The single deterministic renderer the
  SessionStart hook and the user both call: it ranks the signals that fire on a fixed
  severity ladder and leads with the next action. Reads gate state from the **gate cache**
  (`<git-dir>/bench-last-gate`, written by the Stop hook) — never a cold gate run. The
  contract (gate-tested): show-only-on-signal, a five-row budget, a stale-green that is
  not a clean bill, and a zero-severity roadmap footer that never leads.
- **The roadmap capture sink** (`bench idea` / `bench roadmap` → `ROADMAP.md`). The
  capture-and-forget surface: park an out-of-scope idea, commit to nothing, promote it
  later only via `/bench-shape-idea`. The contract (gate-tested in a throwaway repo):
  `idea` appends one dated line and creates the file; a no-arg `idea` errors without
  appending; `roadmap` reports empty when there's nothing parked. `ROADMAP.md` is
  per-consumer content — never in the kit's `package.json` `files[]`.
- **The kit content surface** (`.agents/skills/*/SKILL.md`, `.agents/commands/*.md`).
  Portable harness-facing content. The contract is structural: every skill carries YAML
  frontmatter (name + description) and follows progressive disclosure; every command
  is a phase the index names. Bench craft skill directories stay indexed in
  `AGENTS.md`, but their visible skill names use `craft-*` so `$bench` menus show
  only human-run phase adapters. Codex command-adapter skills are derived from
  `.agents/commands/` and documented in `.bench/BENCH.md`. The gate's conformance
  layer enforces those contracts so disk, docs, and adapters do not drift.
  `.claude/` paths are adapters, not a second source of truth.
- **The safe link contract** (`bench link`). The adoption surface. It must preserve
  project-owned `AGENTS.md` text, install only a managed Bench block plus Bench-owned
  assets, fail on same-named project-owned skills/commands/hooks, and be idempotent
  through the link manifest. It also installs the `.bench/bin/` local CLI set the
  shared hooks use when no global `bench` command is on PATH.
- **AGENTS.md** — the canonical working agreement for project-owned content. `CLAUDE.md`
  imports it (and `.bench/BENCH.md`); never duplicate content there. The four invariants
  and the communication rules are canonical in `.bench/BENCH.md`; the craft skill and
  command indexes live here. The gate checks those indexes, checks command-adapter skills
  against `.bench/BENCH.md`, and checks that the shared rules are not copied back into
  AGENTS.md.

No UI. There is **no design source** for this repo.

## Gate (`.bench/gate.sh`)

```
.bench/gate.sh
```

The oracle for a kit, in layers (all green today):

1. **Parse + validity** — `bash -n` on `bin/bench.sh` and the shared hooks; JSON
   parse on `package.json`, `.claude/settings.json`, and `.codex/hooks.json`. Plus
   two CLI invariants: the scripts the harness execs by path are executable in git,
   and the CLI names the `.sh` gate/done files that exist (an extensionless ref routes
   to auto-detect, not the oracle).
2. **Structure** — every `SKILL.md` carries frontmatter; every `package.json`
   `files[]` path resolves and the dry-run npm package includes required install
   assets while excluding local-only settings.
3. **Kit conformance** — the AGENTS.md index stays in sync with disk for craft
   skills and command references; every indexed skill exists, every craft skill
   exposes a `craft-*` visible name, every command file is referenced as `/name`, every
   command has an explicit Codex `$bench-*` adapter documented in `.bench/BENCH.md`,
   and the feature-build workflow keeps the acceptance-coverage anchors in
   `/bench-write-spec`, `craft-tdd`, `/bench-implement-spec`, and
   `/bench-review-implementation`. This is the check that silently rots and breaks
   harnesses; it is the analog of gl-axi's `axi-conformance`.
4. **shellcheck** — best-effort, runs only when installed (`-S warning`). Not a hard
   dependency; upgrades the shell lint automatically once present.
5. **Safe-link behavior** — the gate runs `bench link` against throwaway repos to
   prove fresh installs, existing `AGENTS.md` preservation, relink idempotence,
   same-name conflicts, modified-managed file protection, Codex/Claude hook adapters,
   shared hooks, and default copy mode.
6. **Canary (meta)** — the gate runs itself against deliberately-broken fixtures in
   `tests/canary/` and asserts each goes red with its targeted error substring. Proves
   the checks above still *bite*: a check rotted into an always-pass fails here. This
   is the gate guarding the gate. Fixtures hide dot-dirs behind a `dot-` prefix (e.g.
   `dot-claude`) so the harness doesn't load fixture skills as real ones; the canary
   restores them at run time. See `specs/canary.md`.

The gate file lives outside `package.json` `files[]`, so it never ships to consumers.

## Lines (model + effort routing)

**Tier → model** (this harness; refresh with `bench models`; set 2026-06-30):
cheap = Haiku 4.5 (`claude-haiku-4-5`) · mid = Sonnet 5 (`claude-sonnet-5`) · top =
Opus 4.8 (`claude-opus-4-8`). Once work is planned, default execution to **mid
(Sonnet 5)**; escalate to top only for the genuinely uncertain seam — no silent
escalation.

- **Skill / command / doc authoring** → **top model, high effort**. Prose that shapes
  agent behavior is the genuinely uncertain, high-leverage seam; the
  `craft-skills` and `craft-adr` skills apply. Spend here.
- **`bench` CLI shell plumbing** → cheap model, low–medium effort at the known seam.
  Mechanical once the gate-resolution and worktree-pool shapes exist.
- **Gate / conformance logic** → mid effort. Correctness of the oracle matters more
  than speed — a wrong gate is the worst class of bug in a kit whose whole premise is
  "the gate is the oracle."

## Notes for cold sessions

- Read `AGENTS.md` first — the working agreement. The four invariants and the
  communication rules are canonical in `.bench/BENCH.md` (AGENTS.md points there); read
  that too. `CLAUDE.md` imports both; edit `AGENTS.md` or `.bench/BENCH.md`, not it.
- `CONTEXT.md` pins the ubiquitous language (gate, oracle, shift, seam, line, …). Use
  those terms exactly; don't invent synonyms.
- The kit's portability across harnesses is a closed decision. Claude and Codex hook
  adapters are interactive layers on top of shared `.bench/hooks/` scripts and the
  harness-independent substrate (the `bench shift` loop + the git `pre-push` hook) —
  never the only thing enforcing an invariant.
- The two `projects/*.md` example profiles (gl-axi, regroup) are shipped templates,
  not live projects. This file (`benchkit.md`) is the profile for this repo.
