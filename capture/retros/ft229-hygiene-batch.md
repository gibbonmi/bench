# Retro — ft229-hygiene-batch

## Outcome

Landed as `c8a3fad2` on `main`, gate green, spec flipped to `Status: implemented`
by that published commit. Reviewed source pair: base `364f34fa`, tip `34fa3046`,
16 commits, composed and gated once on the destination. Source released and
removed.

All 46 stories and all 34 acceptance rows are graded, including `H35`, which the
review round added. Story 9 keeps its named `Not covered` exception. The 37
tickets-only folders are gone and `bench status` renders no residue row.

Two fail-opens at the enforcement boundary were found and closed **after** the
build reported itself complete, both introduced by this spec's own narrowing of
the degraded guard rim. Neither was reachable by the gate, because neither had a
row. That is the headline result of this build, not the seven features.

## Gate-stage timings

From the landing gate's own run record (`.logs/gate-20260820T103114Z`):

| phase | elapsed |
|---|---|
| gofmt | 0.06 s |
| vet | 1.2 s |
| test | 55.5 s |
| race | 4.3 s |
| system | 3.4 s |
| shellcheck | 0.4 s |

The `test` phase is 87% of the wall clock and carries the conformance registry,
so every prose anchor, census rule, and grammar check is paid inside it. Sixteen
gate runs landed this build; at roughly 65 s each that is about 17 minutes of
oracle time, which is the real cost of serial ticket commits on one source.

## Ticket-versus-spec-slice and delegate performance

Ten of eleven tickets were built by write delegates in the prior session, one
per isolated worktree, each landing as a single commit on the retained
integration source. Their per-ticket repair rounds are not retained anywhere I
can read, so the attribution table below records only what the commit graph and
this session's review can observe.

The eleventh ticket — the 37-folder residue deletion — was run inline by the
coordinator. It is a bulk tree edit with an exact printed target list and no
seam, so a delegate would have added a verification hop without adding safety.

Review axes ran three parallel read-only delegates at `opus`/medium, then one
combined-axis delegate over the 12-file repair diff. The combined delegate was
the right proportion for that size and found three real defects, one of which
the coordinator had already missed.

## Coordinator catches

- A delegate reported the rim's escaped-separator fail-open. The coordinator
  reproduced it independently before accepting it, and separately confirmed that
  Go's `encoding/json` emits that exact escaping by default — the difference
  between a contrived case and a producer-shaped one.
- A delegate claimed `filepath.EvalSymlinks` was dead code and the coordinator's
  first mutation probe said otherwise. The probe was wrong: removing the call
  broke compilation, so the "red" was a build error, not a caught behavior
  change. The delegate was right. A mutation probe that reds must be checked for
  *why* it reds.
- The coordinator's own operator-run harness reported every case as allowed. The
  harness was wrong — `$?` was read after a command substitution, so it captured
  `basename`'s exit code. The finding underneath it was still real.
- A delegate's claim that `outline_meta`'s `emitted_symbols` described an absent
  table was refuted by measurement: it equals the dirs table's sum exactly.
  Likewise the claim that a control-byte path renders the same table as an empty
  repository — `outline_dirs[0]` versus `outline_dirs[1]` with a `.` row.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| close-the-light-path-ticket-on-landing | 0 | none |
| count-tickets-only-folders-in-status | 1 | other |
| delete-the-landed-tickets-only-residue | 0 | none |
| correct-the-stale-diagnostics-and-pin-them | 1 | other |
| summarize-the-bare-outline-form | 1 | spec-row |
| close-the-guard-deny-gaps | 1 | other |
| narrow-the-degraded-guard-rim | 2 | spec-row, spec-row |
| name-the-build-command-on-a-cold-session | 0 | none |
| warn-before-removing-the-live-binary | 1 | spec-row |
| prune-gate-run-logs-by-count | 2 | spec-row, tree-drift |
| verify-a-frozen-source-tip-in-preflight | 0 | none |

Seven of the ten rounds are `spec-row`: the coverage map did not ask for the
behavior the repair added. The rim's two rounds share one cause — rows inherited
from the surface being replaced, which by construction cannot mention the surface
replacing it.

## Agent-experience improvements

### Bench CLI

- **`bench worktree path` emits `~/...`.** No shell expands a tilde inside
  quotes, so the verb's output cannot be used directly and every caller rebuilds
  the path from `$HOME`. The coordinator did exactly that for most of this
  session. Emit the resolved absolute path. Parked.
- **No way to gate a composed snapshot without committing.** `bench commit` is
  the only surface that gates the composed path set, and it commits on green as
  its whole purpose. Diagnosing which phase reds therefore means either
  committing junk — which happened, landing a commit literally named `probe` —
  or grading a different subject with a working-tree `go test`, which passed
  while the composed snapshot was red. `bench commit --dry-run` closes it. Parked.
- **The guard reads a heredoc body as command text.** Writing a file whose
  contents mention a destructive git command is refused. This cost two blocked
  writes while capturing this retro's own source material. Parked.
- **`bench worktree exec` is the right verb and nothing durable says so.** A
  parallel session independently specified this while this build ran.

### Skills

- `craft-review`'s citation standard held: every axis cited files and lines, and
  three delegate claims were falsifiable enough for the coordinator to overturn
  two of them by measurement. That is the standard working, not failing.
- The mutation-probe discipline needs one added clause: a probe that reds must be
  confirmed to red for a behavioral reason rather than a compile error. Two
  probes in this session reported a kill that was not one.

### Process

- **The handoff was materially wrong and the tree corrected it.** It said the
  spec was staged and nothing implemented; ten tickets had in fact landed on a
  retained source. `bench status` reporting the handoff's age is what made the
  disagreement visible. The rule that the tree wins earned its place here.
- **Review caught what the gate structurally could not.** Both fail-opens sat in
  a surface with zero assertions, so no amount of gate rigor would have found
  them. The three-axis round is not a formality on an enforcement-boundary diff.
- **Serial ticket commits on one source cost real oracle time** — about 17
  minutes of gate for this build. Worth it here, because the guard and the gate
  are both in the diff; worth questioning on a build with no enforcement surface.
