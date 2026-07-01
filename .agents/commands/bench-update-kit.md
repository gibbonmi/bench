---
description: Upstream sync — pull Matt Pocock's skills + kunchenguid's tooling, diff against what Bench already incorporates, and propose what to adopt. Runs each candidate through the craft-synthesis discipline; records adoptions in the provenance table + CHANGELOG. Proposes, never auto-merges. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-update-kit — sync the kit against upstream

## Entry orientation

This is the upstream maintenance phase. It compares Bench against the current
external source repos, tags candidate changes as upstream-sourced, and runs them
through the synthesis discipline before proposing any adoption.

## Exit handoff

Close by reporting the pulled refs, each candidate's classification, which changes
passed the synthesis loops, and what record would change if the reviewer accepts
the proposal. The recommended next command is `/bench-final-check` after accepted
changes are applied; otherwise no build-phase command follows a rejected proposal.

Bring *others'* improvements into the kit. This is the external, version-driven
input — the counterpart to `/bench-integrate-learnings`, which drains your own usage journal. It
**proposes; you merge.**

## 1. Pull and diff (the input)

- Pull `mattpocock/skills` and kunchenguid's repos (`axi`, `gnhf`, `treehouse`,
  `firstmate`, `no-mistakes`) at HEAD. Record the refs you pulled.
- Read the baseline: the provenance table in `README.md` and prior `CHANGELOG.md`
  entries — what Bench already maps, folds, recommends, or rejected.
- Diff upstream against the baseline into three lists: **new** (techniques not in the
  baseline), **changed** (a mapped/folded item whose source moved), **gone**
  (something Bench uses that was deprecated upstream). Tag every candidate `upstream`.

## 2. Run the synthesis discipline

Hand the candidates to `craft-synthesis`; it owns the discipline end to end
(respect closed decisions, assess, the three quality loops, propose-don't-merge).
Don't restate it here.

## 3. Record (the output)

Only after the loops pass and I've signed off, apply and record: update the provenance
table and append a `CHANGELOG.md` entry (date, refs pulled, what was adopted, what was
rejected and why). That entry is the baseline the next `/bench-update-kit` reads.

The merge is mine; never auto-apply.
