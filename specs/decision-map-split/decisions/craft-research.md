# Research

Status: ready

## Destination

Define a model-invoked `craft-research` skill and consolidate the research
portions of FT99, FT106, FT125, FT304, and FT231. It fans out mutually
independent primary-source reading to read-only subagents, then returns one
verified cited Markdown artifact without duplicating research policy. The map
also fixes the decision-map storage that those research artifacts live in.

## Notes

## Decisions so far

- [Where does research policy live in Bench today, and what must remain single-sourced?](craft-research/tickets/1.md): Resolved 2026-08-02.
- [What contract does the upstream research skill actually establish?](craft-research/tickets/2.md): Resolved 2026-08-02.
- [What has made Bench's existing research fan-outs trustworthy or costly?](craft-research/tickets/3.md): Resolved 2026-08-02.
- [When should `craft-research` fire?](craft-research/tickets/4.md): Reach-for-it-anytime model-invoked guidance (reviewer, 2026-08-02).
- [When and how widely should research fan out?](craft-research/tickets/5.md): Adaptive, round-based fan-out (reviewer delegated to the recommended design, 2026-08-02).
- [What artifact should a research run leave behind?](craft-research/tickets/6.md): One coordinator-authored durable Markdown output per research run, keyed by topic (reviewer delegated to the recommended design.
- [Where is the boundary between research and reviewer-owned judgment?](craft-research/tickets/7.md): Read-side research only (reviewer, 2026-08-02).
- [Which current instructions should become references to `craft-research`?](craft-research/tickets/8.md): Single-source migration (reviewer delegated to the recommended design, 2026-08-02):
- [What quality must a durable research report provide?](craft-research/tickets/9.md): Resolved 2026-09-06.
- [How does this map group related roadmap work?](craft-research/tickets/10.md): Resolved 2026-09-06.
- [What storage boundary keeps the decision map authoritative?](craft-research/tickets/11.md): Resolved 2026-09-06 (reviewer).
- [Which focused reader discovers ready work without another state owner?](craft-research/tickets/12.md): Resolved 2026-09-06 (reviewer).
- [How do FT106's artifact home and freshness projection resolve?](craft-research/tickets/13.md): Resolved 2026-09-06 (reviewer).
- [How should the work compose for specification and measurement?](craft-research/tickets/14.md): Resolved 2026-09-06 (reviewer).

## Not yet specified

## Spec-writer discretion

- The gist line grammar in `Decisions so far`, provided each line links one
  resolved ticket file and the map lane can parse it.
- The ticket file's exact heading text and the `stale` and `ready` row column
  names, provided `bench maps` keeps one row shape.
- The migration script's language and location, provided it prints its target
  list before it applies.
- Exact skill headings and the compact contrastive example required by
  `craft-skills`, provided the example demonstrates independent fan-out versus a
  dependent question that stays serial.
- Exact research-asset heading names and filename slug normalization, provided
  the decided content and location precedence remain intact.
- Exact citation notation for local and external primary sources, provided a
  cold reader can resolve every material claim to the cited evidence.

## Out of scope

- Changing the four decision-ticket types.
- A tracker-backed map on GitLab or any issue tracker; the repository holds
  the map.
- An `Owner` field or any per-ticket claim; the worktree lease is the claim.
- A body reader such as `bench maps show`; the ticket file is the slice.
- A gate red on a stale Source; freshness stays advisory.
- A dual-shape parser; only the split shape parses after the migration.
- Implementing a general-purpose knowledge base, citation database, or web-search CLI.
- Replacing `craft-delegate`, `craft-line`, or harness-native subagent controls.
- Folding formal `/bench-review-implementation` axis review into generic research.
- FT231's three-arm measurement harness.

## Sources

- Path: `specs/decision-map-split/decisions/craft-research/assets/craft-research-research.md`
  Supports: #1 through #3 and the factual premises for #4 through #8. Three read-only research delegations ran 2026-08-02, with upstream sources re-read and local claims spot-checked by the coordinator.
  Drift: re-verify if research, delegation, line-routing, map-source, skill-index, assessment, or artifact-lifecycle guidance changes, or if the cited upstream research contracts move. Re-resolve the asset's line citations before `/bench-write-spec` reads this map if FT164 has landed.
- Path: `specs/decision-map-split/decisions/craft-research/assets/research-workflow-assessment.md`
  Supports: #9 through #14. It compares the local format references with the settled research map and current roadmap and map owners.
  Drift: re-verify when the reference report set, research policy, roadmap rows, map reader, or status projection changes.
- URL: https://github.com/mattpocock/skills/blob/main/skills/engineering/wayfinder/SKILL.md
  Supports: #11 through #13, retrieved 2026-09-06. It supplies the map-as-index, frontier, claim, and asset-link rules the reviewer chose or rejected.
  Drift: re-verify if the upstream skill's storage, ticket, or frontier rules change.
