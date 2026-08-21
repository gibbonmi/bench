# Enrich evidenced refusals through one emitter

Blocked by: none
Writes: internal/worktree/land.go, internal/worktree/land_test.go, internal/worktree/reauthorize.go, internal/worktree/reauthorize_test.go, internal/worktree/classifier.go, internal/worktree/ownership.go
Line: opus / high — this ticket owns the cross-cutting field grammar three siblings consume.

## What to build

A typed refusal value carries `detail` plus optional `observed`, `wanted`, and
`next`; one renderer produces the field text for every enriched site. Land
keeps its stdout `refused{detail=…,observed=…,wanted=…,next=…}` record (absent
fields omitted); release and reauthorize keep their existing
`bench worktree <verb>: ` stderr prefix and exit 1 and append the same field
grammar after their detail text. No field carries a raw control byte. Bounded
path listings render as a TOON path table after the record or message line
(prior art: the cleanup plan's ignored-paths preview) under the classifier's
existing entry limit, with the true total stated; path bytes never enter the
`k=v` grammar. `next=` values follow the orphan-line quoting precedent. The
enriched sites gain their values: the destination-not-clean refusal names
offending paths (switch the site to the `-z` porcelain parser already used by
the resume path); ignored-residue refusals name residue paths; the
reauthorize base-ancestry refusal names the addressed assignment's recorded
start as `wanted`. Contract shared with siblings: the abbreviation ticket, the
recovery-hint ticket, and the incomplete-exit ticket all render through this
grammar.

## Acceptance

- [ ] The destination-not-clean refusal output names offending paths in the bounded table with the true total (covers LR16).
- [ ] The ignored-residue refusal output names residue paths in the bounded table with the true total (covers LR17).
- [ ] Reauthorize with a non-ancestor --base and an otherwise-proven assignment names the recorded start in its stderr refusal fields (covers LR18).
- [ ] A refusal whose observed path carries a newline, ESC, and a comma emits no raw control byte and keeps the record line and the listing table structurally intact (covers LR19).
