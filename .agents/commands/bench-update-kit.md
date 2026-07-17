---
description: Upstream sync — pull Matt Pocock's skills + kunchenguid's tooling, diff against what Bench already incorporates, and propose what to adopt. Runs each candidate through the craft-synthesis discipline; records current decisions in the provenance table and notable user-facing changes in the changelog. Proposes, never auto-merges. Maintenance, not a workflow phase.
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
input — the counterpart to the internal path, where `/bench-what-next` drains your
own usage journal into roadmap items built under `craft-synthesis`. It
**proposes; you merge.**

## 1. Pull and diff (the input)

- Pull `mattpocock/skills` and kunchenguid's repos (`axi`, `gnhf`, `treehouse`,
  `firstmate`, `no-mistakes`) at HEAD. Record the refs you pulled.
- Read the baseline: the provenance table in `README.md` — what Bench currently
  maps, folds, recommends, or rejects. Git owns superseded decisions.
- Diff upstream against the baseline into three lists: **new** (techniques not in the
  baseline), **changed** (a mapped/folded item whose source moved), **gone**
  (something Bench uses that was deprecated upstream). Tag every candidate `upstream`.

## 2. Run the synthesis discipline

Hand the candidates to `craft-synthesis`; it owns the discipline end to end
(respect closed decisions, assess, the three quality loops, propose-don't-merge).
Don't restate it here.

## 3. Record (the output)

Only after the loops pass and I've signed off, apply and record: update the provenance
table so it remains the complete current baseline. Add a concise typed entry under
`CHANGELOG.md`'s Unreleased section only when the adoption changes user-visible
behavior. The approved commit owns the historical record.

The merge is mine; never auto-apply.
