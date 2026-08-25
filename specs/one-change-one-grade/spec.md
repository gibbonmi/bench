# One change, one full grade

Status: implemented

Roadmap: FT215

Decision source: ready compiled map `specs/one-change-one-grade/decisions/one-change-one-grade.md`, resolved 2026-08-25

Verification log: 1 iteration to accept — the round found the phase-table pin claim, a gerund in OG26, ticket 1's lone row, and the undrafted invariant sentence; all four folded

## Problem

One landed change pays the whole-project gate at least twice. The worktree
`bench commit` grades the composed snapshot, and the landing grades the same
change again with the spec flip applied. The two runs never share evidence,
because the landing keys its evidence to `main`'s runner identity. Each run
costs 100–130 s, and the last build paid ten runs for four landings.

A red prose bound costs a full run too. The prose check runs only inside the
gate's `test` phase, so a 27-word sentence surfaces after 100 s and not after
one.

## Solution

A worktree `bench commit` runs the fast lane and not the whole-project gate.
The lane is a declared check list: gofmt, the prose check on the named
Markdown, `go vet`, and `go build ./...`. It runs in a private checkout of the
composed snapshot and takes seconds. A lane pass publishes onto the worktree
branch. A lane fail refuses the commit and names the check.

The landing stays the one whole-project gate. It grades the composed tree with
the flip applied, and it accepts only a green gate verdict. The lane writes no
gate verdict and no evidence. Its record carries its own record class, so no
reader mistakes a lane pass for green.

The kit root carries a built-in lane. A linked repo declares a lane in its
phase manifest, and a repo with no declared lane keeps today's full-gate
commit. The shared rules, the glossary, and the profile state the lane.

## User stories

### The worktree commit runs the fast lane

Line: opus / medium.

The lane changes the landing owner's authority in a worktree and adds gate
code, so the profile's cached routing for gate logic applies.

1. As an operator, I want a worktree `bench commit` to run the fast lane and not the gate, so that a ticket commit costs seconds.
2. As an operator, I want the lane to grade a private checkout of the composed snapshot, so that it grades what I commit.
3. As an operator, I want a lane pass to publish onto the worktree branch, so that the ticket lands on the integration source.
4. As an operator, I want a lane fail to refuse the commit and name the failing check, so that I repair the right thing.
5. As an operator, I want the prose check to grade only the Markdown the commit names, so that a violation elsewhere does not block me.
6. As an operator, I want the prose check to obey `.bench/prose-exclusions`, so that the lane and the gate agree on authored Markdown.
7. As an operator, I want gofmt to rewrite the named Go files before the lane runs, so that a formatting slip costs no lane run.
8. As an operator, I want `bench commit --dry-run` to run the lane and commit nothing, so that I can preview the outcome.
9. As an operator, I want the lane's output to name its outcome as `pass` or `fail`, so that a lane pass never reads as green.
10. As an operator, I want the lane's Bench-owned checks to use the gate's run binary selection, so that they grade with the tree's own code.

### The landing stays the one whole-project gate

Line: opus / medium.

These stories guard the oracle, and the profile routes gate logic mid.

11. As a reviewer, I want the landing to run the gate on the flipped composed tree, so that `main` receives only graded trees.
12. As a reviewer, I want the reviewed landing to accept only a green gate verdict, so that a lane pass never publishes onto `main`.
13. As a reviewer, I want the lane to write no gate verdict and no evidence, so that no consumer reads a lane pass as green.
14. As a reviewer, I want the lane record to carry its own record class, so that the gate reader names it instead of misreading it.
15. As a reviewer, I want a lane record to be never reusable green, so that a misplaced record cannot authorize a landing.
16. As an operator, I want `bench status` to keep the last gate verdict after a lane commit, so that the dashboard states the gate's truth.
17. As an operator, I want the shift and the stop hook to keep the gate as their oracle, so that unattended work keeps its oracle.

### A project declares its lane

Line: opus / medium.

The declaration extends the phase manifest, which the gate owns.

18. As a kit maintainer, I want the kit root to carry a built-in lane, so that this repo needs no declaration file.
19. As a linked-repo owner, I want to declare a lane in the phase manifest, so that my worktree commits run my own checks.
20. As a linked-repo owner, I want a repo with no declared lane to keep the full-gate commit, so that the lane is opt-in.
21. As a linked-repo owner, I want a malformed lane declaration to refuse the commit and name the defect, so that a typo never passes silently.
22. As a reviewer, I want the lane record to name the tree, the lane identity, and the outcome, so that a reader knows what passed.

### The rules and the glossary state the lane

Line: fable / high.

Shared platform rules and the glossary steer every session, so the leverage
override routes them top.

23. As a fresh session, I want invariant 4 to require a lane pass at a worktree commit, so that I skip the gate there.
24. As a fresh session, I want the reference's landing shape to name the lane, so that the reference and the verb agree.
25. As a fresh session, I want `CONTEXT.md` to define `fast lane` and `lane record` with Avoid lists, so that no session says "mini gate".
26. As a fresh session, I want the profile to carry the kit's lane table beside the phase table, so that it advertises the lane's argv.
27. As a reviewer, I want an ADR to record the lane-versus-gate authority split, so that the decided state survives without history.
28. As a reviewer, I want the anchor registry to pin the new invariant-4 sentence, so that an edit that drops it turns the gate red.
29. As a reviewer, I want the CHANGELOG to name the lane, so that the release note carries the behavior change.

### Reviewed exclusions

Line: sonnet / low.

One pinned argv keeps a closed ruling closed.

30. As a reviewer, I want no gate path that selects packages from a diff, so that the rulings in four maps stay closed.
31. As a reviewer, I want the `internal/worktree` test floor left to its own spec, so that this spec stays one capability.

## Implementation decisions

**The lane is a phase table.** A lane is an ordered list of checks with the
phase manifest's entry schema: name, argv, env, needs, optional, dir. The
gate's phase runner executes it in a private checkout of the composed snapshot
under the gate timeout. The kit root's lane is built in, beside the built-in
phase table: `gofmt` through `bench gate-go gofmt`, `prose` through a new
`bench gate-prose` plumbing verb, `vet` through `go vet ./...`, and `build`
through `go build ./...`. A linked repo declares a `lane` array in its phase
manifest. A manifest without a `lane` array, and a root without a manifest
that is not the kit root, declare no lane. A malformed `lane` entry refuses the
run and names the defect, as a malformed phase does today.

**The prose check takes the named Markdown.** The built-in prose entry carries
a placeholder for the named Markdown paths, resolved when the lane runs. A
manifest lane may use the same placeholder. `bench gate-prose <root> [--] [path...]`
grades the named paths through the same walker, classifier, and exclusion list
as the whole-tree prose check. Exit 0 is clean, 1 is findings, 2 is a usage
error. The whole-tree grader composes the same per-subject grader, so the rule
has one source.

**The landing owner is built with its authority.** `bench commit` builds the
owner with the lane authority when the root declares a lane, and with the gate
authority otherwise. The reviewed landing always builds it with the gate
authority and publishes only on a green gate verdict. The authorization result
gains two lane kinds, pass and fail. The owner publishes on green or on a lane
pass, and each construction names the kinds it accepts. `--dry-run` runs the
same authority and stops before publication.

**The lane record is its own record class.** The lane writes one record to a
lane file in the worktree's own Git dir. It never writes the gate cache or the
evidence store. Its fields are `schema`, `tree`, `lane`, `outcome`,
`run_binary`, and `recorded_at`. `outcome` is `pass` or `fail`.

The record
class registers as `lane record` with its own validator. The gate reader names
it and refuses reuse by class, not by name suffix. No consumer
reads the lane file in this spec.

**The run binary.** The lane's Bench-owned checks run through the same run
binary selection the gate's phase table uses. Every executable hop is the
gate's chain: the installed broker composes the snapshot in private storage,
selects the run binary from the composed tree, and records its digest in the
lane record. A lane pass authenticates nothing for the landing, which
re-derives its own chain.

**Guidance.** Invariant 4 in the shared rules gains one exact sentence after
"Commit on green, never on red". The sentence is:

> Green is the landing's whole-project gate, and a worktree commit requires a
> lane pass, not a gate run.

The reviewer vetoes that sentence at sign-off. The reference's
landing shape names the lane at the worktree commit. The glossary defines
`fast lane` and `lane record` with Avoid lists. The profile carries the kit's
lane table beside its phase table. A new conformance check compares that table
with the kit's built-in lane argv.

One ADR records the authority split. The
anchor registry pins the new sentence. The CHANGELOG names the lane.

## Testing decisions

The highest seam that shows each failure is the verb's own output:
`commit.Command` stdout and stderr on a fixture repository, and the landing
journey's stdout. The fixture is the existing commit landing fixture with a
lane declared in its phase manifest and a controllable gate script. The lane
table, the manifest loader, the record class, and the reader are package tests
in the gate package. The prose grader is a unit test. The guidance rows are
live-tree conformance checks. The gate's `test` phase observes all of it.

### Seam diagram

    trigger: bench commit [--dry-run] in a worktree
        │
        ▼
    named paths ──▶ [ attribute, gofmt, compose snapshot ] ──▶ tree
                                                                 │
                     lane declared? ──no──▶ [ gate authority ] ──▶ green | refused
                          │ yes
                          ▼
                    [ lane authority: private checkout, run lane checks ] ──▶ lane record
                          │ pass                          │ fail
                          ▼                               ▼
                    [ publish onto worktree branch ]   refused{check,diagnostic}
                      ◀ tests attach here: fixture repo with a manifest lane, assert stdout and the branch ref

    trigger: bench worktree land
        │
        ▼
    source, destination ──▶ [ compose, flip ] ──▶ [ gate authority ] ──▶ green ──▶ publish onto main
                                                  ◀ tests attach here: the landing journey asserts the gate ran

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| OG01 | 1 | With a declared lane, `bench commit` prints `lane{outcome=pass,checks=...}` and no `phase test:` line. | commit.Command stdout | A commit that still runs the gate prints the phase lines. |
| OG02 | 2 | An unnamed dirty file with a compile error does not fail the lane's build check. | commit.Command stdout | A lane that grades the working tree fails on the unnamed file. |
| OG03 | 3 | After a lane pass, the worktree branch ref points at a commit whose tree equals the composed snapshot. | commit.Command and `git rev-parse` | A lane that composes but does not publish leaves the ref unchanged. |
| OG04 | 4 | A lane check that exits nonzero makes `bench commit` exit 1 and print `lane{outcome=fail,check=<name>}` with the check's first diagnostic line. | commit.Command stdout and exit | A commit that swallows the check's exit publishes on a red check. |
| OG05 | 4 | A lane fail leaves the worktree branch ref unchanged. | `git rev-parse` after commit.Command | A commit that publishes before it reads the outcome lands a red tree. |
| OG06 | 5 | A named Markdown file with a 27-word sentence fails the lane naming the file and the line. | commit.Command stdout | A lane without the prose check passes the sentence. |
| OG07 | 5 | A 27-word sentence in an unnamed Markdown file does not fail the lane. | commit.Command stdout | A whole-tree prose walk fails on the unnamed file. |
| OG08 | 6 | A named Markdown path listed in `.bench/prose-exclusions` produces no finding. | prose per-subject grader unit | A grader that skips the exclusion list grades the excluded file. |
| OG09 | 7 | The lane's gofmt check passes on a named Go file that was misformatted before the commit. | commit.Command stdout | A lane that runs before the rewrite fails gofmt on the file the commit rewrites. |
| OG10 | 8 | `bench commit --dry-run` with a declared lane prints the lane record line and no `phase` line. | commit.Command stdout | A dry run that keeps the gate authority runs the gate. |
| OG11 | 8 | `bench commit --dry-run` after a lane pass leaves the branch ref unchanged. | `git rev-parse` after commit.Command | A dry run that publishes is a commit. |
| OG12 | 9 | The lane's stdout contains `outcome=pass` or `outcome=fail` and never the token `green`. | commit.Command stdout | Reusing the gate's `gate: green` line makes a lane pass read as green. |
| OG13 | 10 | The kit root's lane argv for gofmt and prose names the run binary token. | gate lane table unit | A lane that calls the installed `bench` by name grades with the broker's code. |
| OG14 | 11 | The gate tally is absent after a lane-pass commit and exactly `g` after the landing, which publishes. | landing journey stdout | A landing that trusts the lane pass publishes with no gate run. |
| OG15 | 12 | The reviewed landing owner refuses an authority result of kind lane pass. | landing.Owner unit | An owner that accepts every non-failure kind publishes on a lane pass. |
| OG16 | 13 | After a lane run, the gate cache holds no record for the composed tree. | gate package test on the Git dir | A lane that reuses the gate's verdict writer records a verdict. |
| OG17 | 13 | After a lane run, the evidence store holds no record for the composed tree. | gate package test on the Git dir | A lane that writes evidence lets the landing reuse it. |
| OG18 | 14 | The record-class registry contains `lane record` with fields `lane, outcome, recorded_at, run_binary, schema, tree`. | verdict registry guard test | A class outside the registry is a JSON shape no reader names. |
| OG19 | 15 | `Inspect` on a lane record with outcome `pass` answers `ReusableGreen=false` with a reason naming the lane class. | gate Inspect unit | A suffix-based narrowing misses a class not named `partial verdict`. |
| OG20 | 16 | After a lane commit, `bench status` reports the gate row from the last gate verdict. | status.GateVerdict unit | A status that reads the lane file shows a lane pass in the gate row. |
| OG21 | 18 | At the kit root, the lane table is exactly gofmt, prose, vet, build with the profile's argv. | gate lane table unit | A table missing a check passes a commit that check would fail. |
| OG22 | 19 | A phase manifest with a `lane` array yields a lane table of exactly those checks. | gate manifest unit | A loader that ignores the manifest lane runs the built-in one. |
| OG23 | 20 | A fixture repo with a gate script and no lane leaves the tally exactly `g`, prints no `lane{` line, and publishes. | commit.Command stdout | A commit that skips the gate whenever no lane is declared publishes ungraded. |
| OG24 | 21 | A `lane` entry with an empty argv makes `bench commit` exit 1 naming the defect. | commit.Command stderr and exit | A loader that drops malformed entries grades nothing and passes. |
| OG25 | 22 | The lane record names the composed tree hash, the lane identity, and the outcome. | lane record read | A record without the tree cannot be tied to what passed. |
| OG26 | 23, 28 | The anchor registry requires the invariant-4 sentence, and a removal of that sentence turns the anchor check red. | anchors EvaluateGroup unit | A sentence no anchor pins drifts silently. |
| OG27 | 24 | The reference's landing section contains `lane`. | anchors RequireInSection | A reference that names only the gate contradicts the verb. |
| OG28 | 25 | `CONTEXT.md` defines `fast lane` and `lane record`, each with an Avoid list. | glossary read | The terms stay undefined and the code invents synonyms. |
| OG29 | 26 | The profile's lane table rows equal the kit's built-in lane argv. | profile table conformance check | A table that drifts from code advertises a lane the kit does not run. |
| OG30 | 27 | An ADR with Status accepted records the authority split and names no file path. | ADR read | Without the ADR the split lives only in history. |
| OG31 | 29 | `CHANGELOG.md` names the fast lane under the unreleased heading. | changelog read | A release note without the lane hides the behavior change. |
| OG32 | 30 | The kit root's `test` phase argv stays `go test -count=1 ./...`. | gate phase table unit | A package-selecting test phase changes the argv. |
| OG33 | 5 | `bench gate-prose <root> -- <file>` exits 1 and names the file and the line for a 27-word sentence. | gate-prose command unit | A verb that exits 0 on findings lets the lane pass the sentence. |

Not covered: story 17 — no behavior changes, and the existing stop-hook tests hold.
Not covered: story 31 — an exclusion with no behavior; `specs/worktree-test-floor/decisions/worktree-test-floor.md` owns it.

### Edge inventory

- A named Markdown path with spaces or glob characters is graded by its exact name.
- A named path with control bytes appears sanitized in the refusal, as today.
- A named Markdown path the commit deletes is absent from the snapshot and is not graded.
- A named path that is a FIFO, a device, or a socket fails the lane naming the path.
- A named Markdown file without a trailing newline is graded like any other.
- An empty named Markdown file produces no finding.
- A symbolic link to a Markdown file is not followed and is not graded.
- A commit that names no Markdown runs the prose check with an empty list, which passes.
- A commit that names only Markdown still runs vet and build, because the lane is declared and not derived.
- Two worktrees run lanes at once, and each writes its own lane file in its own Git dir.
- An interrupted or timed-out lane writes no record and publishes nothing.
- A lane fail under `--dry-run` exits 1 and names the check.
- A linked repo whose manifest declares a lane but whose root has no `go.mod` runs exactly the argv it declares.

**Won't handle** lane outcome reuse across runs — every lane runs, and a run of seconds needs no cache.

**Won't handle** a lane in the primary checkout — `bench commit` refuses the primary checkout unchanged.

**Won't handle** a `bench status` row for the lane — the gate row keeps the gate's truth, and a lane row is a separate capability.

**Won't handle** a lane for the shift loop — the shift keeps the whole-project gate as its oracle.

## Ownership fences

- `internal/commit/`
- `internal/landing/`
- `internal/gate/`
- `internal/prose/`
- `internal/status/status_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/conformance/`
- `internal/worktree/land_journey_test.go`
- `cmd/bench/main.go`
- `cmd/bench/command_registry_test.go`
- `.bench/BENCH.md`
- `.bench/BENCH-reference.md`
- `CONTEXT.md`
- `projects/benchkit.md`
- `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md`
- `CHANGELOG.md`

The tickets run serially on one source. The commit ticket calls the lane
authority and the record that the gate ticket lands. The guidance ticket
advertises the argv the gate ticket pins.

## Out of scope

- A `bench status` row that projects the lane record: 6 edits, 1 gate run.
- Lane outcome reuse keyed to the tree: 8 edits, 2 gate runs.
- The `internal/worktree` test floor, owned by `specs/worktree-test-floor/decisions/worktree-test-floor.md`: 30 edits, 4 gate runs.
- The ticket-local evidence machinery of `spec-build-review-gate-cadence`, which waits on FT173 and FT130.
- A gate path that selects packages from a diff. Four maps rule it unsound, and no estimate applies.
- Reopening ADR 0016's evidence key. The map's #2 chose option D.

## Further notes

The shaping that produced this spec landed on 2026-08-25 as `8a798296`. Its
four Markdown files paid two full runs of about 100 s each, one at the
worktree commit and one at the landing. Under this spec the same change pays
one lane run of seconds and one full run at the landing.

The lane's checks are declared, never derived from the diff. Only the prose
check's input list follows the named paths, and that list narrows a check's
subject without narrowing what the landing proves.
