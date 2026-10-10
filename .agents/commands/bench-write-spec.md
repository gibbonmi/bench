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
Then recommend the approved implementation line for fresh ticket authors on one integration source. Review each frozen chunk before its successor, reconcile the final source, and hand it to `bench worktree land`. A staged spec is planning work: `bench preflight build <slug> --plan-only` validates it, and its delivery waits for `bench commitment start`.

## Entry contract

Accept exactly one of three decision sources: a ready compiled map, the
reviewer-confirmed current conversation, or a named reviewed artifact. No
unnamed memory, unreviewed note, or fourth override authorizes a draft. Top-level `decisions/` holds pre-spec
working maps; compiled maps live under `specs/<slug>/decisions/`.

- **Ready compiled map.** Validate it.
  Move the topic folder, its tickets and assets included, into the spec's decisions folder as one unit.
  Move the index file out with its topic folder, from top-level `decisions/` into `specs/<slug>/decisions/`; do not copy either one. Then update every reference to the moved paths in the same green change. A re-run reads the already-compiled spec-local map; it never recreates a top-level copy.
- **Reviewer-confirmed current conversation.** Close every load-bearing product or scope fork here, dated. Do
  not manufacture a map to restate it.
- **Named reviewed artifact.** Name it. It holds settled decisions, not unresolved prompts. A resolved ticket answer and its asset in a multi-item shaping map also count, even while a sibling ticket stays open. The phase copies only that ticket's asset, and the map stays in place.

Record exactly one `Decision source:` line in the spec. For a map-backed source, re-read and re-verify every structured `## Sources`
entry before you choose seams. Disclose what you could not re-read, and consume them in place without copying a research manifest into the spec. Ask at most two late clarification
questions, one at a time, each with a recommended answer; route a dependency
tree or multi-session fog to `$bench-shape-idea`.

## Who runs this phase

The invoking session first fans out read-only research delegates over the one decision source. Charge them under `craft-delegate`, `craft-line`, and `craft-research`, which owns the round, the readers, and the citations. Give them these independent questions:

- the enforcement surface that the rows touch;
- `docs/known-issues.md`, when the repository has it;
- `craft-spec`'s reader sweep;
- the code sites, tests, fixtures, and anchors of each decision-source item.

Each delegate returns its findings and its not-verified list under `craft-research`. The research delegates write nothing. The spec is the durable output of the research run, and its rows cite the findings.

After every research delegate returns, the invoking session forks the current conversation and charges that fork to author the spec and its chunk table only. The charge gives the fork the findings as compressed inputs, and the fork does not repeat their reads. The fork inherits the invoking session's line. Omit a model override. Keep all spec-and-ticket writes in the phase worktree. The invoking session coordinates the phase, slices the tickets, and presents the completed approval table to the reviewer.

After ticket approval, each ticket goes to a fresh author session under `.bench/BENCH.md`'s authorship rule. Spec authoring owns engineering seams, deep-versus-thin design, tests,
acceptance coverage, hostile-input attachment, and gate attachment; shaping sources constrain behavior, scope,
compatibility, or a reviewer-chosen seam.

## Process

1. **Author.** Charge `bench-craft-spec` (and `bench-craft-domain` for terms) on the decision source and the research findings. Lock no row before the findings cover its enforcement surface and `craft-spec`'s reader sweep. Plan the probes and the check floor from the findings, including any `docs/known-issues.md` finding. Read only what the findings leave open, before you lock the rows.

   Apply `craft-spec`'s implementation-line recommendation contract through `craft-line`. Write `specs/<feature>/spec.md` from the template and run `bench coverage --check`. The stale-command-reference sweep remains fail-closed across staged specs.

   A spec that ships a phase declares it on one `Introduces commands:` line, valid in its own directory while staged.
   When no hostile-input checklist class covers a surface, quarry the seams library and propose a tuned profile addition. Apply `craft-spec`'s named
   `Bootstrap authority before execution` rule. A spec that inherits a closed decision about a bench signal quotes that signal's current value in its first status update.
2. **Retire superseded work by promotion then deletion.** Leave no `Superseded by` marker: promote durable
   decisions, delete the old spec under a `spec-retire: <name>` commit, repair references. The same
   promote-then-delete commit removes the spec's `ROADMAP.md` row and that row's `roadmap/FT<n>.md` detail file. Whole-folder retirement removes the compiled topic folders, their tickets and assets included. It removes each compiled index file too.
3. **Review the spec alone.** After the spec fork returns, one independent review round grades the spec alone under `craft-spec`'s review rubric, its chunk table included.
   The round runs on `craft-line`'s spec review line. It returns accept or reject, and slicing starts only after it accepts. The fence section stays empty until the slicing writes it.
4. **Slice one chunk at a time.** Where the stories partition into disjoint package or fence sets, could a narrower
   capability ship on its own gate? Apply `craft-spec`'s named `Bootstrap authority before execution` rule. The ticket graph splits where a consumer branch lands green alone.

   After the spec review accepts, one research round charges two read-only readers for each chunk, with every chunk at the same time. Each reader returns each likely fence expansion under the slicing checks' `Fence-expansion evidence` classes. The invoking session merges each pair into one research table: the agreed set and the frontier. If the spec review changes the chunk table, the research reruns for each changed chunk.

   The slicer runs on `craft-line`'s ticket slicing line and slices the tickets one chunk at a time, in the order of the accepted chunk table. For each chunk, the slicer verifies the research table. Charge `craft-tickets` and write the chunk's tickets under `specs/<slug>/tickets/`. Each `Writes:` line keeps every research path under the slicing checks' keep rule. Compare each ticket's `Writes:` line with `bench preflight build <slug> --propose-writes` under the slicing checks.

   If the slicing session passes about 300k tokens of context, it compacts or restarts at a chunk boundary, never inside a chunk. It carries a ledger. The ledger holds the accepted tickets, the research table, the deferrals, the ordering edges, the spec flags, and the paths that wait for the fence.

   After the last chunk, one review round on `craft-line`'s ticket review line grades the whole ticket set under `craft-tickets`' granularity, edge, and merge-or-split quiz. Its iteration cap is 2. The round returns accept or reject.

   After the ticket review accepts, write the fence section. `bench preflight build <slug> --plan-only` then grades the whole ticket graph. Carry
   its numbered title, `Blocked by:`, and delivered outcome list into the approval table. The spec and tickets
   receive one sign-off.

   Every review round follows one loop. The round declares its iteration cap before the first charge. The invoking session repairs a rejected spec or ticket set, and the same round grades the repair alone. The invoking session folds the findings of each round and names them in the verification log.

   At close, write `Verification log: <n> iteration(s) to accept — <note>` into the spec. The count `<n>` sums the iterations of the spec round and of the ticket round. The note names each round that took more than one iteration.
   When a round takes more than one iteration to accept, append one `capture/learnings.md` entry for that round. The entry names the stage that missed,
   what review caught, why it was missed, and the proposed rule change.
