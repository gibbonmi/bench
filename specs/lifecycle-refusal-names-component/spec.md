# Lifecycle refusals name the component that failed

Status: staged

Roadmap: FT224

Decision source: named reviewed artifact — `roadmap/FT224.md` as drained on 2026-08-25, with the reviewer's deferral of open questions to the author on 2026-08-25

Verification log: 1 iteration to accept — the round found the release site, the ambiguity sentence, and two rows charged to the wrong ticket; all three folded

## Problem

The worktree and landing verbs refuse an identity check with one string that
names several causes. `request, assignment, or path mismatch` covers a
request-token miss, an inactive assignment, and a wrong path. `owner marker
or assignment branch mismatch` covers five predicates. `target is not one
active Bench-owned worktree` hides the resolver's own reason.

An operator then falsifies each cause by hand against the intent ledger.
One wrong token costs a session of diagnosis.

Five faces the row listed have shipped since it opened. The release refusal
names its re-run route and never points at `--discard-ignored`. The
`spec retire` next-step names both board files. The landing serves a
tickets-only `--spec`. The conflict refusal carries its path table. And
`worktree exec` re-roots the child's kit, so the gate names its own skips.

This spec covers only the identity strings that stay collapsed.

## Solution

Every identity refusal names exactly one component: the check that failed.
Where the component has a sanctioned recovery, the refusal names the exact
command. The component set is one executable registry. A test walks that
registry and requires a producing fixture for each entry. So a new component
without a fixture turns the gate red, and a merged string cannot come back.

The refusal shape does not change. `refused{detail=…,observed=…,wanted=…,next=…}`
stays, and the path table stays. Only the `detail` text and the `next`
coverage change. The `worktree exec` and `worktree path` verbs print the
resolver's component sentence on stderr instead of the blanket sentence.

## User stories

### The landing names the identity component

Line: opus / low.

Each story is an exact-spec change at a known seam under a covering gate, so
the cheap mid-tier row applies.

1. As an operator, I want an unknown request token refused as a request-token miss, so that I do not falsify the other causes by hand.
2. As an operator, I want that refusal to name `bench worktree reauthorize` when one active assignment owns my path, so that I recover in one command.
3. As an operator, I want that refusal to name `bench worktree list` when no single active assignment owns my path, so that I find the token's owner.
4. As an operator, I want an inactive assignment refused with its observed state and the wanted state, so that I know the assignment closed.
5. As an operator, I want a path mismatch refused with the assignment's worktree and my target, so that I re-run against the right tree.
6. As an operator, I want an owner-marker mismatch named as the owner marker, so that I do not suspect the branch.
7. As an operator, I want a registration mismatch named as the worktree registration, so that I repair the Git registration and not the marker.
8. As an operator, I want a lock mismatch named as the Bench lock, so that I re-lock instead of re-creating.
9. As an operator, I want the resume landing to name the same components, so that a resumed landing reads like a first landing.
10. As an operator, I want every independent refusal reported in one preflight, so that one fix does not reveal the next refusal a run later.
19. As an operator, I want `bench worktree release` to name the same request component and recovery, so that release and landing read alike.

### The target verbs name the resolver's reason

Line: opus / low.

Each story replaces one blanket sentence with the component the resolver already computes.

11. As an operator, I want `worktree exec` on an unknown target to say that no assignment matches, so that I correct the label or id.
12. As an operator, I want `worktree exec` on an ambiguous label to say which ids collide, so that I pick the id.
13. As an operator, I want `worktree exec` on an inactive assignment to name its state, so that I do not retry a released tree.
14. As an operator, I want `worktree exec` on a broken creation bundle to name the failed bundle component, so that I repair the right thing.
15. As an operator, I want `worktree path` to print the same component sentence as `worktree exec`, so that two verbs do not describe one failure two ways.

### The property is enforced, not remembered

Line: opus / low.

The registry test is the fence that keeps a merged string from returning.

16. As a reviewer, I want the component set to be one executable registry, so that a component without a producing fixture turns the gate red.
17. As a reviewer, I want each producing site to name exactly one component, so that a refusal never lists alternatives for the operator to falsify.
18. As a reviewer, I want the glossary to define `identity component`, so that guidance and code use one term.

## Implementation decisions

The `internal/worktree` package owns one registry of identity components.
Each entry carries its name, its detail sentence, and whether it carries a
recovery command. The registry is a package-level slice that the producing
sites read and the registry test walks. The components are: `request`,
`assignment-state`, `assignment-path`, `owner-marker`, `registration`, and
`lock`.

One constructor builds a refusal from a component plus the observed and
wanted values. The landing preflight, the resume landing, the release verb,
and the creation bundle validator call that constructor. No site composes a
detail string from prose of its own. The release verb keeps its own
`; checkout retained` clause after the component sentence.

The creation bundle validator returns the component error instead of two
merged sentences. The target resolver returns it unchanged, and the
`worktree exec` and `worktree path` verbs print the resolver's error on
stderr after the verb prefix. The selector's unknown-target sentence stays
as it is. Its ambiguous-target sentence gains the colliding ids, because
today it names no id.

Precedence inside one bundle check is the declared registry order: request,
state, path, marker, registration, lock. A bundle check names the first
component that fails. The landing preflight still collects the independent
destination, assignment, and source refusals and prints each one.

The `request` refusal keeps today's recovery: a unique active assignment at
the target yields `bench worktree reauthorize` with that id; any other case
yields `bench worktree list`. Resume accepts `cleanup-pending` as an active
state, as it does today.

`CONTEXT.md` gains one glossary entry for `identity component`.

## Testing decisions

The highest seam that shows each failure is the verb's own output:
`LandCommand` and `ResumeLandCommand` stdout, and `ExecCommand` and the
path verb's stderr. The existing landing fixture plants a real assignment,
and each test mutates one identity dimension before the call. The registry
walk is a package test that maps each component name to the fixture that
proves it. The gate's `test` phase observes all of it.

### Seam diagram

    trigger: bench worktree land | land --resume | exec | path
        │
        ▼
    request, target ──▶ [ identity checks: request → state → path → marker → registration → lock ]
                              │ first failing component
                              ▼
                        [ refusal constructor ] ──▶ refused{detail=<component>,observed,wanted,next}
                              ◀ tests attach here: mutate one dimension, assert the component and next

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LR01 | 1 | An unknown request prints `refused{detail=request token matches no assignment` as its first field. | LandCommand stdout | The old merged string still passes a substring test on `mismatch`; this row pins the new sentence. |
| LR02 | 2 | With one active assignment at the target, the request refusal's `next=` names `bench worktree reauthorize --assignment <id>`. | LandCommand stdout | Dropping the recovery lookup leaves `next=` absent. |
| LR03 | 3 | With zero or two active assignments at the target, the request refusal's `next=` names `bench worktree list`. | LandCommand stdout | Today the refusal carries no `next=` in this case. |
| LR04 | 4 | An assignment in state `complete` prints `detail=assignment <id> is not active,observed=complete,wanted=active`. | LandCommand stdout | A shared state-or-path string cannot carry `observed=complete`. |
| LR05 | 5 | An assignment that owns another worktree prints `detail=assignment <id> owns another worktree,observed=<its path>,wanted=<target>`. | LandCommand stdout | A shared string cannot name the two paths. |
| LR06 | 6 | A rewritten owner marker prints `detail=owner marker does not match assignment <id>`. | LandCommand stdout | The old five-predicate string names the branch too. |
| LR07 | 7 | A re-pointed worktree registration prints `detail=worktree registration does not match assignment <id>`. | LandCommand stdout | A marker-only check reports the marker for a registration fault. |
| LR08 | 8 | An unlocked worktree prints `detail=Bench lock does not match assignment <id>`. | LandCommand stdout | A registration-only check reports the registration for a lock fault. |
| LR09 | 9 | The resume landing prints the same component sentence for each of the six mutations. | ResumeLandCommand stdout | Resume has its own copy of the merged string today. |
| LR10 | 10 | A dirty destination plus an unknown request prints two `refused{` lines in one run. | LandCommand stdout | A first-refusal-exits rewrite prints one line. |
| LR11 | 11 | `worktree exec` on an unknown target prints the selector's no-match sentence on stderr. | ExecCommand stderr | The blanket sentence hides the reason. |
| LR12 | 12 | `worktree exec` on a colliding label prints the selector's ambiguity sentence with the ids. | ExecCommand stderr | The blanket sentence hides the ids. |
| LR13 | 13 | `worktree exec` on a released assignment prints `assignment <id> is not active`. | ExecCommand stderr | The blanket sentence reads as an unknown target. |
| LR14 | 14 | `worktree exec` after an owner-marker rewrite prints the owner-marker component. | ExecCommand stderr | The old bundle validator merges marker and branch. |
| LR15 | 15 | `worktree path` prints byte-identical stderr to `worktree exec` for the same broken target. | path verb stderr | Two verbs with two strings drift. |
| LR16 | 16 | A registry entry without a producing fixture turns the registry test red. | package registry test | A hand-kept list of fixtures misses a new component silently. |
| LR17 | 17 | No production detail sentence in the package contains ` or ` between two component names. | package source test | A merged string reappears without a fixture noticing. |
| LR18 | 18 | `CONTEXT.md` defines `identity component` with an Avoid list. | glossary read | The term stays undefined and the code invents synonyms. |
| LR19 | 19 | `bench worktree release` with an unknown request prints the request-token sentence, `; checkout retained`, and the reauthorize `next=`. | ReleaseCommand stderr | The release site keeps its own copy of the merged string today. |

### Edge inventory

- The request token is empty, unsafe, or a digest-shaped string; the grammar and the digest rule refuse before identity, unchanged.
- The target is a label that collides, an id prefix that collides, or a path with control bytes; the selector refuses, unchanged.
- The assignment is `cleanup-pending` on a first landing, which refuses as not active, and on a resume, which proceeds.
- The owner marker file is absent, malformed, names another owner, or names another path; each names `owner-marker`.
- The registration is absent or names another branch; each names `registration`.
- The lock is absent or carries another reason; each names `lock`.
- Two components fail at once inside one bundle; the first in registry order is named.
- The destination refusal and an identity refusal fail at once; both print.

**Won't handle** naming every failing bundle component in one refusal — the first component names the check that blocked; the per-refusal preflight survives.

**Won't handle** the reauthorize verb's own identity proofs — that verb returns one sentence per proof today, and the landing refusal that routes to it survives.

**Won't handle** the selector's unknown and ambiguous sentences as registry entries — each names one cause today, and the exec and path verbs stop hiding them.

## Ownership fences

- `internal/worktree/classifier.go`
- `internal/worktree/worktree.go`
- `internal/worktree/land_identity.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/ownership.go`
- `internal/worktree/path.go`
- `internal/worktree/exec.go`
- `internal/worktree/identity_component.go`
- `internal/worktree/identity_component_test.go`
- `internal/worktree/land_reauthorization_test.go`
- `internal/worktree/land_resume_refusal_test.go`
- `internal/worktree/land_identity_test.go`
- `internal/worktree/identifier_operand_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/worktree_test.go`
- `cmd/bench/command_registry_test.go`
- `CONTEXT.md`
- `CHANGELOG.md`

The tickets run serially on one source, because ticket 02 calls the
constructor and registry that ticket 01 lands.

## Out of scope

- The landing's authority questions, the interrupted-landing recovery, and the board merge rule stay on FT169: at least 40 edits and 4 gate runs.
- A refusal that lists every failing component of one bundle is a separate rendering capability: at least 8 edits and 2 gate runs.
- The reauthorize verb's identity-proof refusals as registry entries: at least 6 edits and 1 gate run.
- A machine-readable component field in the `refused{` record for AXI consumers: at least 10 edits and 2 gate runs.

## Further notes

The five shipped faces the row listed were verified against the tree on
2026-08-25: `retainedReleaseError` names the re-run route, `roadmapRemainder`
names both board files, `TicketsOnlyFolder` serves the tickets-only close,
the conflict refusal carries `paths`, and `execEnv` re-points the child's
wrapper. The row's open question, whether the landing serves a tickets-only
spec, closed with that verification. The reviewer deferred the remaining
questions to the author on 2026-08-25.
