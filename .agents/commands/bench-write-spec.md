---
description: Turn a reviewed decision source into a build spec with user stories, engineering seams, and testing decisions. Use before any build past the lighter-path threshold in .bench/BENCH.md.
---

# /bench-write-spec — lock the seams before the loop runs

## Entry orientation

Charge `bench-craft-spec` on one authorized decision source. It owns the process: explore, find seams, and
synthesize, with no interview. It owns the discipline: stories, acceptance coverage map, edge inventory,
**Won't handle**, hostile-input checklist, fences, and review rubric. It yields `specs/<feature>/spec.md` plus tickets.

## Exit handoff

The spec carries `Status: staged` (staged → implemented at the green gate → promote-then-delete on merge). Stop for sign-off.
Then recommend a fresh mid-tier build session on one retained integration source. Review its frozen base and tip, and hand the accepted source to `bench worktree land`.

## Entry contract

Accept exactly one of three decision sources: a ready compiled map, the
reviewer-confirmed current conversation, or a named reviewed artifact. No
unnamed memory, unreviewed note, or fourth override authorizes a draft. Top-level `decisions/` holds pre-spec
working maps; compiled maps live under `specs/<slug>/decisions/`.

- **Ready compiled map.** Validate it. Then move (do not copy) the source map and any map-owned assets from top-level
  `decisions/` into `specs/<slug>/decisions/`. A map-owned asset stays in `decisions/assets/`. Then update every reference to the moved paths in the same green change. A re-run reads the already-compiled spec-local map; it never recreates a top-level copy.
- **Reviewer-confirmed current conversation.** Close every load-bearing product or scope fork here, dated. Do
  not manufacture a map to restate it.
- **Named reviewed artifact.** Name it. It holds settled decisions, not unresolved prompts.

Record exactly one `Decision source:` line in the spec. For a map-backed source, re-read and re-verify every structured `## Sources`
entry before you choose seams. Disclose what you could not re-read, and consume them in place without copying a research manifest into the spec. Ask at most two late clarification
questions, one at a time, each with a recommended answer; route a dependency
tree or multi-session fog to `$bench-shape-idea`.

## Who runs this phase

The session holding the decision source authors the spec and tickets at whatever tier it runs. A fresh session
builds after ticket approval. Spec authoring owns engineering seams, deep-versus-thin design, tests,
acceptance coverage, hostile-input attachment, and gate attachment; shaping sources constrain behavior, scope,
compatibility, or a reviewer-chosen seam.

## Process

1. **Author.** Charge `bench-craft-spec` (and `bench-craft-domain` for terms) on the decision source; read the
   enforcement surface before you lock rows that touch it. Run `craft-spec`'s reader sweep before that lock.
   Give each story group its resolved model and effort from `craft-line`. Write `specs/<feature>/spec.md` from
   `craft-spec`'s template and run `bench coverage --check`. The stale-command-reference sweep remains fail-closed across staged specs.

   A spec that ships a phase declares it on one `Introduces commands:` line, valid in its own directory while staged.
   When no hostile-input checklist class covers a surface, quarry the seams library and propose a tuned profile addition. Apply `craft-spec`'s named
   `Bootstrap authority before execution` rule.
2. **Retire superseded work by promotion then deletion.** Leave no `Superseded by` marker: promote durable
   decisions, delete the old spec under a `spec-retire: <name>` commit, repair references. The same
   promote-then-delete commit removes the spec's `ROADMAP.md` row and that row's `roadmap/FT<n>.md` detail file.
   Whole-folder retirement removes the compiled maps and map-owned assets, plus tickets.
3. **Slice, then review once.** Where the stories partition into disjoint package or fence sets, could a narrower
   capability ship on its own gate? Apply `craft-spec`'s named `Bootstrap authority before execution` rule. The ticket graph splits where a consumer branch lands green alone.

   Charge `craft-tickets` and write the breakdown under `specs/<slug>/tickets/`. Carry
   its numbered title, `Blocked by:`, and delivered outcome list into the approval table. The spec and tickets
   receive one sign-off.

   One review round covers the spec-and-tickets pair with the reviewer-named model, and `craft-tickets`' granularity/edges/merge-split quiz is its approval step. The round applies `craft-spec`'s review rubric.

   The round declares its iteration cap before the first charge. The author folds partials left after the round and names them in the verification log.
   `--reviewer <tier> [effort]` overrides the round's delegate. The tier resolves through the invoking harness's
   own `.bench/lines.env` column and runs same-family through its native agent surface. For example,
   `--reviewer top high` under Codex resolves `BENCH_CODEX_TOP` at high. A model id is an invocation error.

   At close, write `Verification log: <n> iteration(s) to accept — <note>` into the spec. When the round takes more
   than one iteration to accept, append one `capture/learnings.md` entry. The entry names the stage that missed,
   what review caught, why it was missed, and the proposed rule change.
