# Adoption proposal for the first quality milestone

Status: proposal for reviewer review. This file is not an active policy.

This file proposes the first delivery commitment for the Bench kit. The kit adopts it only after the delivery commitment prerequisite publishes and the reviewer gives explicit direction.
The build that wrote this file did not write `.bench/commitment.json`. It ran no `plan` and no `approve` against this repository.

## Decision source

Decision 14 in `specs/roadmap-delivery-commitment/decisions/roadmap-delivery-commitment/tickets/14.md` sets the milestone.
The milestone covers FT376, FT373, and FT349, in that order.
The September 29 survey rows stay uncommitted intake. This includes the staged FT358 spec.
Any finding outside these three owners stays uncommitted.

## Source identities

The three owner files have the same bytes on `main` at `3347fdbdca6cb69726ce67dedfcf7299bbdbb55e` and on the integration tip `4107d7513cae91d3be48451631dd62a2e7cc4b1a`.

| Order | Outcome | Source path | Git blob | Plan source identity |
|---|---|---|---|---|
| 1 | FT376 | `roadmap/FT376.md` | `2aa9ea813d4312f6dc4f4814c0331d135992324e` | `sha256:2d5a52f6bf9e764cf1511395b052a8447160c4680300cf1a2494d7ac300df4f4` |
| 2 | FT373 | `roadmap/FT373.md` | `a54205e4c79764786c2f89fc73200c5f812445ab` | `sha256:64628d70448b2d9fe32ccdc2466af3a69ad192d019d767954fa4acb8f86c9938` |
| 3 | FT349 | `roadmap/FT349.md` | `a35a5e55cd72432f2017731d50dcfb3a905c4aa0` | `sha256:f4f89dbf1f4de267d1877ad2a3eb9ba7c410bebc7df46214d3fc573e6b51a268` |

The plan source identity is the value that `bench commitment plan` checks.
The `identity` cell of a roadmap row in `bench commitment inventory` is a different digest. Do not copy that cell into the plan input. Copy the value that a plan refusal names, or use the value in this table.

## Verified remaining obligations

Each owner keeps its complete obligation. No part of the three obligations has shipped.

FT376 remains open.
Only the drain commits `c9a9f4e0` and `4dd5e178` of 2026-10-03 changed its owner file.
`bench-craft-spec` has a reader sweep and a cheapest-wrong-result rule. It has no pre-review trace of a pin row operator, an unexported entry read, or a derived expectation. It also has no table of old and new rules for a consolidation spec.
No spec, ticket, or landing names FT376 outside this feature's own spec files.

FT373 remains open.
Only the drain commits `c9a9f4e0` and `4dd5e178` changed its owner file.
No gate phase, script, or Go package runs the `modernize` analyzers. No spec or ticket folder delivers the check.
The current recommended sequence names an FT373 light-path ticket, but no such ticket exists in `specs/`.
The module lists `golang.org/x/tools` v0.49.0, so the dependency precondition still holds.

FT349 remains open.
Only the drain commit `d49d9af6` of 2026-09-27 changed its owner file.
The ticket check floor exists only as guidance in `bench-implement-spec`. No preflight or lane code derives it from a `Writes:` line.

`internal/coverage` still exempts a seam cell that holds `planned`. No checkpoint refuses a leftover `planned` citation.
No spec, ticket, or landing names FT349 outside this feature's own spec files and drain records.

## Ordered remaining outcomes

1. FT376: before the first review charge, the spec stage traces each pinned check, entry read, and derived expectation to its grader.
2. FT373: a gate check refuses production code that re-implements a kept standard-library function.
3. FT349: the ticket checkpoint enforces the check floor and refuses a leftover `planned` citation.

The order follows decision 14. FT376 has the class feature, and FT373 and FT349 have the class fix.
The general rule puts fixes before features, but the explicit approved order governs. No literal dependency joins the three outcomes.

## Plan input

Save this document to a file outside the checkout. The commands below call that file `<file>`.

```json
{
  "version": 1,
  "milestones": [
    {
      "id": "quality-1",
      "outcomes": [
        {
          "id": "FT376",
          "criteria": [
            {"id": "FT376.trace", "text": "Before the first review charge, the spec stage traces each pin row operator, each unexported entry read, and each derived expectation to its grader."},
            {"id": "FT376.consolidation", "text": "A spec that consolidates repeated rules maps each changed old-rule and new-rule cell to an acceptance row, a flagged addition, or a Won't handle line."}
          ],
          "sources": [{"id": "FT376", "path": "roadmap/FT376.md", "identity": "sha256:2d5a52f6bf9e764cf1511395b052a8447160c4680300cf1a2494d7ac300df4f4"}]
        },
        {
          "id": "FT373",
          "criteria": [
            {"id": "FT373.check", "text": "A gate check refuses production code that re-implements slices.Contains, slices.ContainsFunc, a maps.Copy loop, min, max, strings.Cut, strings.CutPrefix, or strings.CutSuffix."},
            {"id": "FT373.sites", "text": "No production site in the kept categories remains when the gate turns the check on."}
          ],
          "sources": [{"id": "FT373", "path": "roadmap/FT373.md", "identity": "sha256:64628d70448b2d9fe32ccdc2466af3a69ad192d019d767954fa4acb8f86c9938"}]
        },
        {
          "id": "FT349",
          "criteria": [
            {"id": "FT349.floor", "text": "The ticket checkpoint refuses a commit whose checks omit a package that the ticket Writes line derives."},
            {"id": "FT349.planned", "text": "The ticket checkpoint refuses a planned citation for a row that the committed ticket covers."}
          ],
          "sources": [{"id": "FT349", "path": "roadmap/FT349.md", "identity": "sha256:f4f89dbf1f4de267d1877ad2a3eb9ba7c410bebc7df46214d3fc573e6b51a268"}]
        }
      ]
    }
  ],
  "active_milestone": "quality-1"
}
```

The input has no `continuations` key. The section below gives the reason.
No outcome names a deliverable yet. A later plan adds each staged spec or ticket folder as a deliverable of its outcome.

A disposable clone of the integration tip replayed this input. `bench commitment plan` returned plan `sha256:9f4981f343b46a8dd558494ada7c54f5c2aea11aff2f705dbb00649389a20786` with predecessor `absent`. The effects were `added` and `activated` for each outcome. A second plan returned the same identity.
An approval in a planning worktree of that clone staged the policy and projected the recommended sequence to the three outcomes.

## Source order and plan identity

The plan identity depends on the traversal order of the policy sources. This input fixes that order: FT376, then FT373, then FT349, with one source for each outcome.
A canonical source order is still an open reviewer decision. Make that decision before the first real approval. Otherwise, a later change to the traversal order can change the plan identity while every bound fact stays the same.

## Continuations

The proposal lists no continuation.
On 2026-10-04, `bench commitment inventory` shows 37 active runs. None of them needs a continuation, for these reasons:

- `dc-integration` builds this prerequisite. Its landing releases the run before adoption.
- The 27 `dc-c<n>-*` runs are review sessions of this build. A review session publishes no production scope.
- `jev-shaping` is a research and shaping run. It writes planning paths, and a planning path needs no continuation.
- The 8 `tt-cmp-*` runs are model comparison trials. They hold no approved deliverable.

Run the inventory again just before the approval. List a run only if it is still active and finishes approved production scope. Each listed run needs its `assignment` and `request` cells and its `scope` paths.

## Replayable commands

Run these commands after the prerequisite publishes and after the reviewer directs adoption. Use the newly installed Bench version.

1. Read the current runs and sources: `bench commitment inventory`.
2. Compare each source identity in the table with the plan refusal or plan result. If an owner file changed, verify the obligation again before you continue.
3. Create the owned planning worktree: `bench worktree create --request <request> --label quality-1-adoption`.
4. In that worktree, validate the transition: `bench commitment plan --input <file>`.
5. In that worktree, approve the exact plan: `bench commitment approve --plan <id> --decision <reference> --delayed none --removed none`. Use the plan identity from step 4 and the reviewer's decision reference.
6. Commit the staged `.bench/commitment.json` and `ROADMAP.md`, then land the planning worktree through `bench worktree land`.
