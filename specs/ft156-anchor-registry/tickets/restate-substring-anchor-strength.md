# Restate substring anchor strength

Blocked by: expose-anchor-query.md, strip-comments-from-required-anchors.md
Ownership fence: `internal/anchors`, `projects/benchkit.md`, `ROADMAP.md`, `specs/ft156-anchor-registry/spec.md`, `specs/ft156-anchor-registry/decisions/ft156-anchor-registry.md`, `specs/ft187-communication-surface-cut/spec.md`
Contracts: the substring matcher's actual guarantee crosses `internal/anchors`→the kit profile, roadmap, and staged-spec claims, asserted by PS1 through an enumerated prose sweep of every named surface

## What to build

Sweep registry comments and diagnostics, the kit profile, the roadmap, and staged specs. Restate every overclaim so substring anchors promise placement and presence only, preserving the closed deferral of stronger paraphrase-resistant matching.

## Acceptance

- [ ] [PS1] No enumerated live kit claim says substring anchors prevent or close paraphrase evasion; every surviving claim states only the placement-and-presence guarantee the matcher implements.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PS1 | restore one paraphrase-prevention overclaim in a swept surface | the independent semantic review | enumerate matcher prose, profile, roadmap, and staged-spec matches; cite the restored claim as contradicting the real substring matcher |
