# Retro: parallel-session-landings

## Outcome

Published to `main` as `3947bcca` ("feat: parallel session landings"), parents
`105d243d` (destination base) and `1968620d` (reviewed source tip), tree
`5f6abb26` carrying `Status: implemented`; project-green advanced to the
published commit and the source assignment `3ae1acc9` released clean. The
reviewed source pair was `be5ec93e..1968620d`: eight original tickets, three
dogfood repairs, nine first-round review repairs, the PL15/PL16 public race
journey, seven second-round review repairs, and one composition-conflict
relocation. Destination-side enablement commits: the staged-spec fence sync,
the landing-allowance declaration, and two capture pins. Both merged specs
(`parallel-session-landings`, `roadmap-progressive-index`) retired.

## Gate-stage timings

Composed-tree landing gate: ~84 s wall (log start 09:46:55Z, green evidence
recorded 09:48:19Z), packages parallel; slowest packages `internal/worktree`
68.9 s, `internal/publication` 30.1 s, `internal/conformance` 16.3 s, race and
system phases within the same window. Interactive `bench commit` runs during
repairs held the same shape (~1.5–2 min wall each including race/system).

## Ticket-versus-spec-slice and delegate performance

Ticket-sized charges performed well: the seven-ticket repair batch (one opus
delegate, serial gate-then-commit) landed 7/7 with no second gate attempt; the
PL15/PL16 rework delegate (opus) returned green in one pass and correctly
reported a load-bearing charge deviation (the fingerprint recheck fires before
CAS at the public seam) instead of forcing the charged assertion. The three
axis reviews (fable, high) each returned verified, citation-backed findings
with no overlap misses. No delegate was handed a whole spec slice this
session; the retained-integration-source pattern kept every charge ticket-sized.

## Coordinator catches

- The in-flight race test's design (a second landing as winner) contradicted
  the destination-scoped gate lock the spec itself blesses; recharged with a
  plain-git winner before delegating.
- The first conflict-resolution attempt (absorbing destination files into the
  source) was wrong — it entangled FT198's conformance graph and would break
  fence authorization; replaced by a semantics-preserving block relocation
  after the reviewer redirected to /bench-debug.
- `prospective authorization refused: infrastructure` after a fully green gate
  was traced with a throwaway `InspectTree` probe to the undeclared
  `BENCH_HOME` closure variable (now a profile cold-session note).
- Self-catches recorded in learnings: bare `--help` on a worktree verb opened
  an interactive subshell; `dist/bench` was rebuilt with plain `go build`
  against the profile note (repaired via `scripts/go-build.sh`).

## Repair attribution

Rounds observed this session (first-round review targets, second-round review
targets, and landing repairs), attributed to the ticket whose surface they
repaired; prior-session rounds are recorded only where the handoff retained
evidence.

| Ticket | Repair rounds | Cause per round |
|---|---|---|
| compose-reviewed-source-trees | 2 | delegate-error (duplicated porcelain parser); spec-row (CAS-vs-infrastructure classification) |
| land-reviewed-sources-atomically | 2 | delegate-error (grammar-derived counts duplicated); spec-row (movement refusals lacked named rerun action) |
| resume-published-landings | 3 | spec-row (destructive destination states unchecked); shaping-ambiguity (resume landing-identity authority needed a reviewer choice); spec-row (declared ignored allowance missing on resume) |
| reauthorize-retained-assignments | 1 | delegate-error (independent request-digest derivations) |
| expose-explicit-source-review | 3 | spec-row (explicit-base movement drift unrevalidated); delegate-error (resume spec-path re-derived); other (literal-roots enumeration edge) |
| route-workflow-through-integration-sources | 3 | delegate-error (PL15/PL16 journey asserted an impossible concurrent-landing winner); delegate-error (pasted land-invocation harness); tree-drift (both branches inserted at one anchor; relocation commit `1968620d`) |
| restore-ft198-doctrine-commit | 1 | none |
| finish-and-land-ft198-source | 3 | other (explicit-base dogfood repair); other (gate-log fingerprint repair); other (abbreviated-resume-base repair) |

## Agent-experience improvements

### Bench CLI

- Landing refusals name only a class (`undeclared ignored residue`,
  `composition conflict: textual`, `staged spec bytes`); each cost a manual
  enumeration pass. Emitting the offending paths (bounded) in the refusal
  detail would cut most of this session's landing loop iterations.
- `prospective authorization refused: infrastructure` after a printed green
  gate gave no reason; surfacing the subject-closure reason (`declared
  environment unavailable`) in the refusal would have made the BENCH_HOME
  gotcha a one-shot fix. Friction: three failed land attempts plus a probe.
- A worktree verb given unknown flags opens an interactive subshell instead of
  failing usage-red; with the parked stale-dist discovery idea, this belongs in
  one ergonomics pass.

### Skills

- `bench-implement-spec --full` could state the landing invocation
  prerequisites (source-built binary path, `BENCH_HOME`, clean destination,
  declared allowance) where it names the final landing, sparing the next
  session the discovery loop now pinned in the profile.

### Process

- The ephemeral-token design worked as specified: two token losses at session
  boundaries were both recovered by `reauthorize` with zero persisted secrets.
- The reviewer's /bench-debug redirect mid-conflict was the right brake: the
  repro-loop discipline (merge-tree probe, dangling commit-tree candidate,
  composed-tree test venue) found the minimal fix the ad-hoc path was walking
  away from.
