---
name: craft-review
description: The three-axis review judgment — Standards, Spec, Coverage — and the citation standard a finding must meet. Use whenever reviewing a diff; the /bench-review-implementation phase, a delegate's returned work, a PR, or a self-review before commit all charge their review from here. Advisory, never the oracle.
index: reviewing a diff / what a finding must cite
---

# Reviewing on three axes

Review catches what the gate can't see, and it is advisory: review says whether
the work is *good*, the gate says whether it is *done*, the reviewer decides
whether it ships. A finding never overrides the gate, and a clean review never
substitutes for a green one.

## The axes stay separate

Code can pass one axis and fail another — the right thing built against the
conventions, clean conventions around the wrong thing, a correct happy path with
open edges — and merging the axes lets one mask another. Findings are reported
under separate headings, never reranked into a single list. Each axis ends with
its count and its worst issue.

- **Standards** — every place the diff violates a documented convention: the
  working agreement, the shared platform rules, the project profile, any
  conventions docs. Knowledge duplication is a Standards finding — two
  derivations of one fact is the code-standard defect review exists to grade.
  Separate hard violations from judgment calls. Skip anything the gate already
  enforces; double-reporting what tooling caught is noise.
- **Spec** — three hunts: requirements missing or partial; behavior nobody asked
  for (scope creep); requirements implemented but wrong. "Implemented but wrong"
  includes calls the diff introduces into APIs, flags, or config keys that don't
  exist — a symbol the dependency doesn't export is a finding one grep or
  `--help` run confirms. When the spec carries
  an acceptance coverage map, audit every row — a missing, partial,
  falsely-classified, or unclosed mapped behavior is a Spec finding.
- **Coverage** — the adversarial pass: name concrete inputs or states that would
  break the diff and that no acceptance row or existing test exercises. Walk the
  edge classes (the generic classes live in `/bench-write-spec`'s
  edge-inventory step) plus the profile's hostile-input checklist when one
  exists — the two together are the full inventory. An edge the spec
  explicitly marked won't-handle is not a finding; an edge nobody decided is.

## What a finding must cite

A finding without a citation is an opinion.

- Standards: the rule, named or quoted precisely.
- Spec: the spec line, quoted. For implemented-but-wrong behavior a traced
  execution is also a valid citation — the concrete input, the value the diff
  produces, and the expectation it violates — so a happy-path logic bug is
  reportable even when no spec line enumerates it.
- Coverage: the input or state, the expected break, and the row or test that
  should exist.

```
Spec: story 4 asks "adapter refuses an unbound BENCH_MODEL"; the diff guards
claude and codex but not opencode — no guard call in the opencode adapter.
```
Good — quotes the story, names the object, and the gap is checkable in one look.

```
The error handling in the adapters feels incomplete and could be more robust.
```
Bad — no citation and no object, so nothing can be checked, fixed, or vetoed.

## Refute before you report

Before a finding lands, try to kill it with the repo: grep for the test the
Coverage finding claims is missing, re-read the convention the Standards finding
cites, run the command the Spec finding says is broken. A finding the repo
refutes in one command wastes the reviewer's attention and teaches them to skim
the rest. "No findings" is a real result — state what was examined, not a bare
LGTM.
