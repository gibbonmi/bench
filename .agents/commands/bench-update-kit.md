---
description: Upstream sync — pull Matt Pocock's skills + kunchenguid's tooling, diff against what Bench already incorporates, and propose what to adopt. Runs each candidate through the craft-synthesis discipline; records current decisions in the provenance table and notable user-facing changes in the changelog. Proposes, never auto-merges. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-update-kit — sync the kit against upstream

## Entry orientation

This is the upstream maintenance phase. It compares Bench against the current
external source repos. It tags each candidate change as upstream-sourced. It
runs every tagged change through the synthesis discipline before it proposes
any adoption.

## Exit handoff

Report the pulled refs, each candidate's classification, and which changes
passed the synthesis loops. Report what record would change if the reviewer
accepts the proposal. `/bench-final-check` is the recommended next command
after you apply accepted changes. No build-phase command follows a rejected
proposal.

Bring *others'* improvements into the kit. This is the external, version-driven
input, the counterpart to the internal path, where `/bench-drain` drains your
own usage journal into roadmap items built under `craft-synthesis`. It
**proposes. You merge.**

## 1. Pull and diff (the input)

- Pull `mattpocock/skills` and kunchenguid's repos (`axi`, `gnhf`, `treehouse`,
  `firstmate`, `no-mistakes`) at HEAD. Record the refs you pulled.
- Read the baseline: the provenance table in `README.md` names what Bench
  currently maps, folds, recommends, or rejects. Git owns superseded decisions.
- Sort upstream against the baseline into three lists. **New** names a
  technique absent from the baseline. **Changed** names a mapped or folded
  item whose source moved. **Gone** names something Bench uses that upstream
  deprecated. Tag every candidate `upstream`.

## 2. Run the synthesis discipline

Hand the candidates to `craft-synthesis`. It owns the discipline end to end:
respect closed decisions, assess, run the three quality loops, and propose
rather than merge. Do not restate it here.

## 3. Record (the output)

Apply and record only after the loops pass and I sign off. Update the
provenance table so it stays the complete current baseline. Add a concise
typed entry under `CHANGELOG.md`'s Unreleased section only when the adoption
changes user-visible behavior. The approved commit owns the historical record.

The merge is mine; never auto-apply.
