# Fast-lane input facts for FT215

Produced 2026-08-27 by one read-only mid-tier delegate. The coordinator
spot-checked the cited attribution, lane, and parser producers at `97bb035`.

## Current boundary

`bench commit` resolves the paths that the caller names. The resolver
normalizes, deduplicates, sorts, and contains each path inside the repository.
It accepts files, safe directories, deletions, and symlinks. It refuses the
repository root, special files, and unsafe descendants
(`internal/landing/attribution.go:22-106`).

The built-in fast lane always runs `gofmt`, `prose`, `vet`, and `build`.
The first two use the selected Bench binary. The Go tool commands inspect
`./...`. The lane omits the test phase and writes only a lane record
(`internal/gate/lane.go:31-74`).

Before the lane, `bench commit` formats named Go files. A named directory
includes Go descendants for that formatting step. The prose lane receives
only directly named paths whose names end in `.md`; it does not expand a
named directory (`internal/commit/commit.go:79-103,135-145,214-272`).

## Producer-derived partitions

| attributed path | current observer | safe narrowing fact | unresolved choice |
|---|---|---|---|
| ordinary Go source | gofmt, vet, build | all three Go checks can observe it | package scope stays spec-owned |
| Go module or build metadata | vet, build | Go can observe it | exact metadata class |
| named embed input | vet, build | its importing Go package can observe it | one embed registry or derived scan |
| ordinary Markdown | prose | Go checks have no shown producer relation | whether prose alone authorizes the lane |
| `.bench/prose-exclusions` | prose | prose always loads it | it must select prose |
| decision map | prose plus map parser elsewhere | map discovery is location-specific | whether to add focused validation |
| spec or ticket | prose plus coverage or preflight elsewhere | parser ownership is location-specific | which focused command joins the lane |
| roadmap board | prose plus roadmap validator elsewhere | board paths have a dedicated input source | whether to add focused validation |
| tracked retrospective | prose plus retrospective validator elsewhere | the owner reads direct Markdown children | whether to add focused validation |
| ignored local capture | no committable producer is established | the files stay local and ignored | exclude or prove composition |
| shell, JSON, HTML, or generic asset | usually none | only explicit embed inputs affect Go | complete-lane fallback |
| deletion | observers grade the composed absence | the old path supplies a class | classification and fallback |
| rename | deletion plus addition | both paths can supply classes | require both sides or fall back |
| named directory | composed descendants | the raw operand has no single class | expand the composed diff |
| symlink | no safe focused observer is established | attribution does not traverse it | complete lane or refusal |
| mixed known paths | union of their checks | the union preserves each known observer | union rule |
| any unknown regular path | no complete focused relation | current lane is the conservative baseline | complete lane or refusal |
| special or unreadable path | none | attribution refuses it | no lane decision remains |

## Named input owners

The tree has four embed inputs:

- `.bench/consumer-payload.json`
- `internal/adopt/prepush.sh`
- `internal/releaseevidence/registry.json`
- `internal/releaseevidence/requirements.json`

Their producers are explicit
(`consumer_payload.go:22`; `internal/adopt/link_hook.go:19`;
`internal/releaseevidence/types.go:80-83`).

Decision maps, specs, tickets, roadmap rows, and retros each have a focused
parser or validator. Those owners do not currently participate in the lane
(`internal/maps/schema.go:23-89`; `internal/spec/spec.go:110-127`;
`internal/preflight/gather.go:107-139,326-366`;
`internal/roadmap/tree.go:39-43`; `internal/retros/retros.go:30-43`).

The conformance input registry distinguishes Go source, decision documents,
the roadmap board, capture retros, and several named files. Catch-all remains
the default. These labels are evidence about input families, not an executable
lane selector (`internal/conformance/registry/registry.go:77-103,115-153`).

## Safe classification boundary

The classifier can use the composed changed paths, not only the raw command
operands. This resolves directory descendants and represents a rename as its
deletion and addition.

A positive class can omit a check only when a producer proves that the check
cannot observe the class. Any unknown path must retain the current complete
lane unless the reviewer chooses refusal.

## Reviewer decisions

- Ticket #10 chooses the known classes, mixed-path union, and unknown fallback.
- Ticket #10 also chooses whether a direct Markdown path runs prose alone.
- The spec must decide whether focused document validators join the lane.
- The spec must preserve attribution refusals for special and unreadable paths.
