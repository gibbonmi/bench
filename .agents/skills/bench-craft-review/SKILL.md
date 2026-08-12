---
name: craft-review
description: The three-axis review judgment — Standards, Spec, Coverage — and the citation standard a finding must meet. Use whenever reviewing a diff; the /bench-review-implementation phase, a delegate's returned work, a PR, or a self-review before commit all charge their review from here. Advisory, never the oracle.
index: reviewing a diff / what a finding must cite
---

# Reviewing on three axes

Review catches what the gate can't see, and it is advisory: review says whether
the work is *good*, the gate says whether it is *done*, the reviewer decides
whether it ships. A finding never overrides the gate, and a clean review never
substitutes for a green one. Every finding also carries exactly one
repair-routing disposition — `no-op`, `auto-fix`, or `ask-user` — defined and
enforced by `/bench-review-implementation`.

## Re-derive, then compare

Every axis derives its facts from the current primary source *before* it
compares the candidate against them. A declaration-only confirmation — trusting
what the ticket or the diff's own commit message claims — is incomplete: an
axis that never re-reads its source cannot catch a claim the source refutes.
Every finding cites its derivation, not a recollection of it. The three axes
run in parallel fresh contexts so one axis's derivation cannot seed another's.
Before any candidate-controlled execution (a script, a test, a tool the diff
itself introduces), ask what authenticates the verifier: a candidate's own
proof of correctness is not evidence until something outside it confirms it.

Review also treats a compiled map's defaulted decisions as authoritative unless
the spec explicitly overrides them — a claimed repair is checked against both
its coverage row and the applicable defaulted-decision table.

## The axes stay separate

Code can pass one axis and fail another — the right thing built against the
conventions, clean conventions around the wrong thing, a correct happy path with
open edges — and merging the axes lets one mask another. Findings are reported
under separate headings, never reranked into a single list. Each axis ends with
its count and its worst issue.

- **Standards** — independently reread the current working agreement, the shared
  platform rules, the project profile, and any conventions docs, then hunt every
  place the diff violates one. Knowledge duplication is a Standards finding — two
  derivations of one fact is the code-standard defect review exists to grade.
  When the diff landed as delegate slices, hunt duplication at fence boundaries
  specifically — a shared primitive no slice owned arrives derived once per fence
  (the ownership fence lives in `craft-spec`; `craft-tickets` owns the
  independently-green tracer grouping and the advisory `Writes:` disjointness
  note). A false disjointness claim across two tickets' `Writes:` notes is a
  Standards finding. The smell baseline below rides this axis.
  Comment prose is graded against `craft-comments` — a comment that narrates
  the change, cites provenance, or argues its own correctness is a Standards
  finding. Separate hard violations from judgment calls; skip what the gate
  already enforces.
- **Spec** — the approved spec drives the behavior; quote the applicable spec
  line rather than trusting what a ticket or commit message claims was built.
  Three hunts: requirements missing or partial; behavior nobody asked for (scope
  creep); requirements implemented but wrong, including calls the diff
  introduces into APIs, flags, or config keys that don't exist — a symbol the
  dependency doesn't export is a finding one grep or `--help` run confirms.
  When the spec carries an acceptance coverage map, audit every row — a
  missing, partial, falsely-classified, or unclosed mapped behavior is a Spec
  finding.
- **Coverage** — before hunting, independently enumerate two things from the
  approved spec and the tree, not from the diff's own account: the
  producer-derived input family the change actually consumes, and the
  spec-authorized write set it's allowed to touch. Then run the adversarial
  pass: name concrete inputs or states that would break the diff and that no
  acceptance row or existing test exercises. Walk the edge classes (the generic
  classes live in `bench-craft-spec`'s edge inventory) plus the profile's
  hostile-input checklist when one exists — the two together are the full
  inventory. An edge the spec explicitly marked won't-handle is not a finding;
  an edge nobody decided is.

## The smell baseline

The Standards axis carries a baseline beneath the documented conventions:
Fowler's classic smells (*Refactoring* ch. 3): Mysterious Name, Duplicated Code,
Long Function, Long Parameter List, Global Data, Mutable Data, Feature Envy,
Data Clumps, Primitive Obsession, Repeated Switches, Shotgun Surgery, Divergent
Change, Lazy Element, Speculative Generality, Message Chains, Middle Man,
Refused Bequest. When a hunt needs more than the name, each smell's gloss and
its fix live in `references/smell-baseline.md`.

A documented repo standard overrides the baseline wherever they disagree, and a
smell is always a judgment call, never a hard violation — file baseline
findings under judgment calls.

## What a finding must cite

A finding without a citation is an opinion.

- Standards: the rule, named or quoted precisely; a baseline finding names the
  smell and quotes the hunk.
- Spec: the spec line, quoted. For implemented-but-wrong behavior a traced
  execution is also a valid citation — input, produced value, and violated
  expectation — so a happy-path logic bug is reportable with no spec line for it.
- Coverage: the input or state, the expected break, and the row or test that
  should exist.

One rule cuts across the axes: a universal claim without an enumeration is a
sample. A finding quantified over a set — "no caller", "every fixture", "all N
are safe" — cites the enumeration of that set (the grep over all of it, the
per-member run), not one measured member extended to the rest. A sampled claim
says so and names the unmeasured remainder.

```
Spec: story 4 asks "adapter refuses an unbound BENCH_MODEL"; the diff guards
claude and codex but not opencode — no guard call in the opencode adapter.
```
Good — quotes the story, names the object, and the gap is checkable in one
look. Bad, by contrast: "The error handling feels incomplete and could be
more robust" — no citation and no object, so nothing can be checked or fixed.

## Refute before you report

Before a finding lands, try to kill it with the repo: grep for the test the
Coverage finding claims is missing, re-read the convention the Standards finding
cites, run the command the Spec finding says is broken. A finding the repo
refutes in one command wastes the reviewer's attention and teaches them to skim
the rest. "No findings" is a real result — state what was examined, not a LGTM.
