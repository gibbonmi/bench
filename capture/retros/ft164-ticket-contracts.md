# Retro: ft164-ticket-contracts

Promoted 2026-08-03, commit `83c630e`, candidate `cbdcf7d`, run terminal.
Line: opus/high for every write delegate and every review (reviewer override of
the spec's fable/high — opus followability was the build's purpose); fable
coordinated and verified. Ten tickets (seven planned, three repair), three
composed review rounds (opus ×3) plus one `gpt-5.6-sol` falsification pass.

## What worked

- **The contract carried the build.** Ten tickets written in the shape the spec
  teaches all assigned through the live parser with zero hand-normalization —
  including the terminal dogfood: a fresh-context author given only the
  template section wrote a ticket that `bench spec build assign` accepted
  first try (run `dogfood-ft164`, retired after the probe).
- **Self-probes earned their cost immediately.** The export delegate's charged
  mutation exposed a message-blind hostile-ticket test the suite would have
  kept green forever; the missing refusal test exists because the charge said
  run the mutation, not reason about it.
- **The falsification pass out-found the same-family review.** Four of its five
  findings were enforcement blind spots (unpinned AB2 row, half-pinned junction
  rule, substring-satisfiable per-ID check, normalized-away EOF case) that the
  opus review had graded implemented. Cross-family refutation is worth its
  price on gate-touching builds.
- **The build dogfooded its own failure stories.** The land ticket's fence
  missed `registry_test.go` — the registry-tracing gap story 5 teaches — and
  two coordinator probes were vacuous (fence-fallback swallow; wrong `### Red
  mutations` occurrence) exactly per the probe-site/kind rules being landed.

## What missed

- **Tip movement under an active run cost one abandoned run.** A capture commit
  before the first checkpoint hit the zero-checkpoint recomposition defect
  (empty patch, unrecoverable); the run was rebuilt from snapshotted diffs.
  Rule now in learnings: freeze the tree from `start` to `promote`; the
  `--full` handoff write belongs before start or after promote.
- **Two kit defects surfaced and are parked for the Codex hotfix:** `start`'s
  no-exact-green remediation names a command that cannot succeed from a
  reduced tip, and zero-checkpoint recomposition feeds git an empty patch
  (both in `capture/IDEAS.md`; the reviewer holds the prompt).
- **The section resolver is fence-blind**, which forced `###` headings inside
  the taught template/examples plus a write-`##`-in-real-files prose line; the
  cold dogfood author flagged that contradiction as their one hesitation. A
  fence-aware `markdownH2Sections` would retire the workaround.
- **The spec's own enumeration drifted mid-build** (24→27 needles after repair
  rounds); trued up in round 3. An enumeration a build extends needs its
  update named as part of the repair-round definition.

## Open reviewer-veto surface (not applied)

Duplicate re-derivation sentence (fold to a pointer?); the `###` heading
workaround; `ParseTicket` doc comment naming its consumer; the inventory's
grouping label for two classify-section needles; row 7.2's moved-artifact
clause vacuously satisfied for the held probe-kind entry; the anchor test file
already over its `structure.budgets` grant before this build.
