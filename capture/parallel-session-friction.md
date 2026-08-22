# Commit/promote friction observed 2026-08-03 (injected-interface-junctions build)

Grill input for a shape-idea on multi-session parallelism. Every item below
was hit live in one session sharing a tree with one human writer. Each item
names the observed refusal and the design question it raises.

1. **Whole-tree cleanliness is a global writer lock.** `bench spec build
   start`/`assign` require a clean working checkout, and `bench commit` fails
   closed on any dirty path outside its named set. The reviewer's in-flight
   drafts (`ROADMAP.md`, an unstaged spec) blocked an unrelated one-file ticket
   commit. The only exits were a mixed-authorship commit naming the reviewer's
   files (needed explicit authorization) or waiting. Question: what is the
   smallest unit that must be exclusively locked for a green landing — the
   tree, or the diff plus its gate inputs?

2. **Gate verdicts bind to the whole tree, not the published diff.** A
   path-scoped commit grades the tree as it stands, so green-for-graded-tree ≠
   green-for-published-commit whenever another writer's files are present.
   Today resolved by refusal (item 1). Parallel sessions need per-diff
   grading — e.g. gate an ephemeral composed tree the way `promote` already
   gates its prospective tree.

3. **Any tip movement forces every active run through promote-recompose.** An
   unrelated capture/ticket commit moved `main`; `assign` then refused with
   "requires recomposition: bench spec build promote", and recomposition
   discards a held review. With N sessions committing, every landing forces
   every active run through recompose and re-review — serialization plus
   review churn. Question: can recomposition be scoped to actual overlap
   (fence intersection) instead of any ancestry change?

4. **Lifecycle mutations only run from the primary checkout.** `checkpoint`
   invoked from inside an assignment worktree refused ("spec build spec does
   not belong to working checkout"). All coordinators funnel through one
   checkout for every lifecycle mutation — a structural single-coordinator
   assumption.

5. **Run state is one locked file with coordinator-assembled evidence.**
   `.git/bench/specbuild/<slug-digest>.json` plus hand-built receipts (tree
   hash via throwaway index, ticket digest, assumption digests — schema
   reverse-engineered from `checkpoint.go`, again). Single-writer fine;
   parallel coordinators contend on the lock and each repays the schema cost.
   The parked receipt-skeleton emitter becomes load-bearing, not ergonomic.

6. **`promote` conflates recompose with publish.** Recomposing onto a moved
   tip is routine maintenance, but it ships inside the publish-authority
   command. So the harness permission layer (correctly protective of publish)
   also blocks the maintenance half. A human had to run a promote whose only
   effect was recomposition. Splitting recompose from publish would let
   sessions maintain their runs without publish authority.

7. **Ticket files must be committed to main before assign** ("ticket no longer
   matches committed subject"). This makes ticket authoring re-enter the
   item-1 tree lock and the item-3 tip-movement cascade. The cycle observed
   live: commit ticket → tip moves → recompose required → human-gated
   promote.

Ambient: session start reported 48 open assignments and dozens of preserved
recovery refs — assignment sprawl is already multi-session shaped; the
lifecycle above is not.
