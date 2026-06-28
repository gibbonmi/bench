# Project: benchkit

The Bench kit itself — a harness-agnostic agent-development workflow shipped as the
npm package `benchkit`. It is not an application: it is shell + markdown + JSON that
other repos consume. The deliverable is the `bench` CLI (`bin/bench.sh`), the working
agreement (`AGENTS.md`), and the `.claude/` skills, commands, and hooks. Because the
artifacts are plain files, the kit must work identically under Claude Code, Codex, and
any other AGENTS.md harness — that portability is the product.

## Seams (test here; everything else is free to change)

- **The gate contract** (`.bench/gate.sh` / `bench gate`). The oracle surface. Everything
  in Bench routes through its exit code: 0 = shippable, non-zero = not done. The
  highest seam — if the gate is weak, the whole system is weak. Test by feeding it a
  conformant tree (green) and a broken one (red); never by trusting a reading of the
  diff.
- **The `bench` CLI subcommands** (`gate`, `worktree`, `shift`, `init`). The
  operational shell surface. Stable command names and exit codes are the contract;
  the implementation behind each is free to change. Keep gate resolution
  (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect) in one place.
- **The kit content surface** (`.claude/skills/*/SKILL.md`, `.claude/commands/*.md`).
  Harness-facing content. The contract is structural: every skill carries YAML
  frontmatter (name + description) and follows progressive disclosure; every command
  is a phase the index names. The gate's conformance layer enforces the index ⇆ disk
  sync.
- **AGENTS.md** — the single canonical working agreement. `CLAUDE.md` only imports it;
  never duplicate content there. The four invariants and the skills/commands index
  live here, and the gate checks that the index matches what's on disk.

No UI. There is **no design source** for this repo.

## Gate (`.bench/gate.sh`)

```
.bench/gate.sh
```

The oracle for a kit, in layers (all green today):

1. **Parse + validity** — `bash -n` on `bin/bench.sh` and the hooks; JSON parse on
   `package.json` and `.claude/settings.json`.
2. **Structure** — every `SKILL.md` carries frontmatter; every `package.json`
   `files[]` path resolves (npm-pack integrity).
3. **Kit conformance** — the AGENTS.md index stays in sync with disk, both
   directions: every skill dir is referenced, every indexed skill exists, every
   command file is referenced as `/name`. This is the check that silently rots and
   breaks harnesses; it is the analog of gl-axi's `axi-conformance`.
4. **shellcheck** — best-effort, runs only when installed (`-S warning`). Not a hard
   dependency; upgrades the shell lint automatically once present.

The gate file lives outside `package.json` `files[]`, so it never ships to consumers.

## Lines (model + effort routing)

- **Skill / command / doc authoring** → **top model, high effort**. Prose that shapes
  agent behavior is the genuinely uncertain, high-leverage seam; the
  `writing-great-skills` and `adr` skills apply. Spend here.
- **`bench` CLI shell plumbing** → cheap model, low–medium effort at the known seam.
  Mechanical once the gate-resolution and worktree-pool shapes exist.
- **Gate / conformance logic** → mid effort. Correctness of the oracle matters more
  than speed — a wrong gate is the worst class of bug in a kit whose whole premise is
  "the gate is the oracle."

## Notes for cold sessions

- Read `AGENTS.md` first — the four invariants and the working agreement override
  convenience. `CLAUDE.md` is just a one-line import of it; edit `AGENTS.md`, not it.
- `CONTEXT.md` pins the ubiquitous language (gate, oracle, shift, seam, line, …). Use
  those terms exactly; don't invent synonyms.
- The kit's portability across harnesses is a closed decision. Anything Claude-only
  (the Stop hook, the PreToolUse git guard) is an *extra* interactive layer on top of
  the harness-independent substrate (the `bench shift` loop + the git `pre-push`
  hook) — never the only thing enforcing an invariant.
- The two `projects/*.md` example profiles (gl-axi, regroup) are shipped templates,
  not live projects. This file (`benchkit.md`) is the profile for this repo.
