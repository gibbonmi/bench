# Add a maintainer-run per-run reclaim verb

Blocked by: extract-terminal-run-ref-enumeration.md
Ownership fence: `internal/specbuild`, `internal/spec/build.go`, `cmd/bench/specbuild.go`, `bin/bench.sh`
Contracts: the classified ref inventory crosses `internal/specbuild`→`cmd/bench/specbuild.go`'s renderer, asserted by RM4 against the real enumeration's output over every classification the inventory can carry; the reclaim grammar crosses `internal/spec/build.go`→`cmd/bench/specbuild.go` dispatch, asserted by RM5
Assumptions: `extract-terminal-run-ref-enumeration.md` has landed, so the classified inventory exists and is the only source of what a run no longer needs; the plan/apply shape and fingerprint discipline copy `abandon`'s rather than inventing a grammar; one slug per invocation and no set-selecting predicate; claims re-derived from the tree at pickup

## What to build

`bench spec build reclaim <slug>` plans the provisional residue of one
spec-build run, and `--apply <fingerprint>` deletes it on the exact planned
fingerprint. It exists because promotion did not reclaim its own provisional refs
until recently: runs that promoted before that left assignment branches,
candidate refs, and checkpoint refs behind, release compacted the intent rows
that named the branches, and nothing enumerates the residue.

The verb consumes the classified inventory from the previous ticket — it does not
carry a second derivation of what a run no longer needs. It deletes only the
reclaimable class. Refs classified retained, unclassified, or ambiguous are
reported in the plan and left intact, so a maintainer sees the whole picture and
the tool acts on only the part it can prove dead.

Shape it on `abandon`: a plan carrying a fingerprint over its own facts, an
`--apply` that refuses an empty or drifted fingerprint and demands a fresh plan,
and a TOON receipt. One slug per invocation, for the same reason discard takes
one ref — a maintainer reclaiming residue is making a judgment per run, not
authorizing a sweep.

The grammar table, the help suffix, the operation list, the dispatch arm, the
service interface, and the launcher's help block all need the new row. The
guidance prose that advertises the operation count is a separate ticket.

## Acceptance

- [ ] [RM1] `reclaim` with no `--apply` plans without mutating any ref, and the plan carries a fingerprint over its own facts.
- [ ] [RM2] `--apply` with the exact planned fingerprint deletes the reclaimable refs of a terminal run and leaves every other class intact.
- [ ] [RM3] `--apply` with an empty, malformed, or drifted fingerprint refuses and mutates nothing, naming the fresh-plan action.
- [ ] [RM4] the receipt reports the retained, unclassified, and ambiguous classes alongside the reclaimable one, so a maintainer sees what was left and why.
- [ ] [RM5] `reclaim` is rejected with the usual exit-2 usage error when given no slug, more than one slug, or an unknown flag.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RM1 | delete the reclaimable refs during the plan call | the plan-is-read-only test | move the delete loop into the planner, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the unchanged-refs assertion to fail |
| RM2 | delete every ref the inventory returned rather than the reclaimable class | the apply-scope test | drop the class filter, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the surviving-retained assertion to fail |
| RM3 | accept any non-empty fingerprint | the drifted-fingerprint refusal test | drop the equality comparison, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the refusal assertion to fail |
| RM4 | render only the reclaimable count in the receipt | the receipt-classes test | drop the other three columns, run `go test ./internal/specbuild -run Reclaim -timeout 180s`, expect the reported-classes assertion to fail |
| RM5 | give `reclaim` a grammar with no positional bound | the reclaim usage test | set the grammar's max args to unbounded, run `go test ./internal/spec -timeout 120s`, expect the exit-2 assertion to fail for two slugs |
