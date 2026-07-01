# Human-facing skill names — make Bench phases self-explanatory in Codex

> **GRILLED & RESOLVED (2026-06-30).** Rename the portable `.agents` phase surface to
> human-facing verb/object names, keeping command filenames and Codex adapter skill
> names identical so Claude, Codex, and other harnesses share one vocabulary.

## Grounding

- Bench now exposes workflow phases in Codex through explicit `$bench-*` adapter
  skills under `.agents/skills`.
- The canonical phase bodies still live in `.agents/commands/*.md`.
- Human usability is the priority for the adapter skill names: the user should know
  why to run a skill from its name in the skill selector.
- `bench` should stay spelled out. `bn-*` saves characters but loses namespace
  clarity and recognizability.

## #1: What are the human-facing adapter names?

Type: Grill

### Question
What should the explicit Codex adapter skills be called so a human can tell why to
run each one?

### Answer
Resolved from grill. Use verb/object names that describe the job the user is
starting:

| Current adapter | New human-facing name | Why run it |
| --- | --- | --- |
| `$bench-setup` | `$bench-setup-repo` | Configure this repo for Bench. |
| `$bench-ideate` | `$bench-shape-idea` | Start with a fuzzy idea and shape it into decisions before a spec. |
| `$bench-spec` | `$bench-write-spec` | Define the build target before implementation. |
| `$bench-diagnose` | `$bench-debug` | Reproduce and isolate a bug before fixing it. |
| `$bench-build` | `$bench-implement-spec` | Implement the approved spec. |
| `$bench-review` | `$bench-review-implementation` | Check the implementation against the spec and repo standards. |
| `$bench-qa` | `$bench-final-check` | Run the final external gate check after review fixes. |
| `$bench-update` | `$bench-update-kit` | Pull external improvements into the Bench kit. |
| `$bench-learn` | `$bench-integrate-learnings` | Promote lessons from `.bench/learnings.md` back into the kit. |

`$bench-final-check` is deliberately plain English rather than gate jargon because
it communicates finality to the reviewer.

## #2: Do we rename only the Codex adapter skills, or also the canonical command files?

Blocked by: #1
Type: Grill

### Question
The agreed names are human-facing skill names. Should the build rename only the
explicit Codex adapter skills and keep `.agents/commands/*.md` stable, or should it
also rename the canonical command files and Claude slash-command surface?

### Answer
Resolved from grill. Rename the canonical `.agents/commands/*.md` files and the
matching Codex adapter skills together. `.claude/commands` and `.claude/skills` are
adapters to `.agents`, so renaming the portable `.agents` surface gives Claude and
Codex the same human-facing vocabulary instead of maintaining a separate Codex-only
name layer.

## #3: If adapter names diverge from command filenames, where does the mapping live?

Blocked by: #2
Type: Research

### Question
If `$bench-shape-idea` delegates to `.agents/commands/bench-ideate.md`, the current
one-command-one-same-name-adapter gate contract must change. What is the smallest
gateable mapping surface: per-skill body references only, an explicit table in
`.bench/BENCH.md`, command frontmatter, or a separate manifest?

### Answer
Resolved from grill. No separate mapping surface. Keep the same basename across the
canonical command file, the Codex adapter skill, and harness-facing invocation. The
filename is the mapping: `.agents/commands/<name>.md` pairs with
`.agents/skills/<name>/SKILL.md`. This keeps movement between Claude, Codex, and
other harnesses frictionless.

## #4: What docs should name the human-facing adapter skills?

Blocked by: #2, #3
Type: Research

### Question
Once the mapping is chosen, identify every current-state surface that should mention
the human-facing names: `.bench/BENCH.md`, README, project profile, AGENTS pointer,
and the gate conformance comments.

### Answer
Resolved from grill. Update current-state surfaces that name the phase commands:
`.bench/BENCH.md`, `AGENTS.md`, `README.md`, `projects/benchkit.md`, the adapter
skill bodies/frontmatter, and `.bench/gate.sh` literals. Do not rewrite historical
`decisions/` content except this map if it needs a final status note.
