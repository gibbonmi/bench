# Close the chain: reconcile FT100 and FT179, changelog, permanent exclusions, handoff

Blocked by: 01c-register-the-check-and-exclusions.md, 01b-pin-thresholds-terms-and-profile.md, 02-split-bench-core-from-reference.md, 02b-rewrite-core-and-reference.md, 03-rewrite-skills-batch-1.md, 04-rewrite-skills-batch-2.md, 05-rewrite-skills-batch-3.md, 06-rewrite-skills-batch-4.md, 07-rewrite-skills-batch-5.md, 08-rewrite-commands-batch-a.md, 09-rewrite-commands-batch-b.md, 10-rewrite-thin-adapters.md, 11-rewrite-root-docs-batch-a.md, 12-rewrite-root-docs-batch-b.md, 13-rewrite-root-docs-batch-c.md, 14-rewrite-roadmap-bodies.md, 15-rewrite-decisions-batch-1.md, 16-rewrite-decisions-batch-2.md, 17-rewrite-decisions-batch-3.md, 18-rewrite-decisions-batch-4.md, 19-rewrite-specs-tickets-capture.md, 20a-rewrite-comments-worktree-source.md, 20b-rewrite-comments-worktree-tests.md, 21a-rewrite-comments-conformance-a-to-l.md, 21b-rewrite-comments-conformance-m-to-z.md, 22a-rewrite-comments-gate.md, 22b-rewrite-comments-adopt.md, 23a-rewrite-comments-publication-preflight.md, 23b-rewrite-comments-roadmap-roadmapflow.md, 24-rewrite-comments-shift-status-skillsindex-lines.md, 25a-rewrite-comments-coverage-git-spec.md, 25b-rewrite-comments-gitguard-env-structure-handoff.md, 26-rewrite-comments-cmd-diff-bounds-and-peers.md, 27a-rewrite-comments-landing-freshness-intent.md, 27b-rewrite-comments-systemtest-maps-releaseevidence.md, 27c-rewrite-comments-small-packages-a.md, 27d-rewrite-comments-small-packages-b.md, 28-rewrite-shell-comments.md, 28b-rewrite-embedded-templates.md
Writes: roadmap/FT100.md, roadmap/FT179.md, CHANGELOG.md, .bench/prose-exclusions (new in ticket 01c), internal/conformance/prose_mechanics_test.go (new in ticket 01c), capture/session-handoff.md
Line: `sonnet` / low under Claude Code, `gpt-5.6-luna` / low under Codex — the texts are decided and the delegate transcribes them.

## What to build

`roadmap/FT100.md` gains one paragraph. It says this spec delivered the core-versus-reference partition with anchors moved and the STE form across the surface. It says the demonstrated-delta audit, the cut line, and the budget measure remain. It names FT231 as the literal block and FT89 as the recommended precursor, as the dependency tables state. The row's mention of a 175-line bound becomes a pointer to the profile's budget table, which holds the current number.

`roadmap/FT179.md` gains one paragraph. It says this spec rewrote every explanatory comment under the current register, deleted no-value comments, and removed provenance tags on edited lines. It says the high-stakes surface documentation, the `craft-comments` rule additions, and the `craft-review` update remain. Each of the two rows gains one `Occurrence:` line naming `asd-ste100-progressive-disclosure`. Both rows keep their heading line, `Next:` line, and ledger byte-identical.

`CHANGELOG.md` gains one `Unreleased` entry in ASD-STE100. The entry names four things:

- the new `prose-mechanics` conformance check
- the `.bench/prose-exclusions` file
- the split of `.bench/BENCH.md` into an always-loaded core and reference material
- the STE rewrite of the kit's guidance and comments

`.bench/prose-exclusions` holds exactly the four permanent rows. The live-tree subset test narrows its approved set to those four rows in the same commit. `capture/session-handoff.md` is rewritten in full for the review phase.

## Acceptance

- [ ] FT100 and FT179 hold the decided paragraphs and one new `Occurrence:` line each, with their grammar lines byte-identical (covers PD37, PD39).
- [ ] `.bench/prose-exclusions` holds exactly the four permanent rows, and the subset test's approved set equals them (covers PD32).
- [ ] `CHANGELOG.md` holds the entry, and the handoff is rewritten (covers PD38).
- [ ] The live-tree test passes over the whole tree (covers PD27).
