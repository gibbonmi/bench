## Standards

3 findings. Worst: High hard violation.

1. **High — hard violation:** `internal/preflight/requirements.json:23`–`:25`
   owns producer owner/schema facts, but
   `internal/preflight/release_requirements.go:253`–`:267` re-encodes the same
   owner set and mapping. Registry validation is also split between
   `internal/preflight/types.go:138`–`:160` and
   `internal/preflight/release_requirements.go:53`–`:109`. This violates
   `AGENTS.md`'s one-source-per-fact rule for production policy, parsers, and
   executable registries.
2. **High — hard violation:** the component-manifest wire schema is authored by
   the production generator at `scripts/build-release-evidence.mjs:141`–`:151`
   and independently re-declared by the production verifier at
   `internal/preflight/artifact_evidence.go:29`–`:50` and `:187`–`:268`. The
   `AGENTS.md` independence exception applies only to test expectations, not two
   production implementations of the same schema.
3. **Medium — hard violation:** `bin/bench.sh:292` advertises the user-visible
   `--profile public|bank` release-preflight contract, but
   `CHANGELOG.md:7`–`:18` has no Unreleased entry. AGENTS.md requires kit edits
   to follow `craft-synthesis`, whose adoption rule requires a concise typed
   changelog entry for user-visible behavior.

## Spec

4 findings. Worst: P1.

1. **P1:** hosted release can never reach publication. The spec says focused
   runs are diagnostic and hosted callers remain red only until required owner
   records exist (`specs/governed-release-evidence.md:99`–`:104`), but
   `.github/workflows/release.yml:60` invokes focused publish and
   `internal/preflight/command.go:148`–`:153` always rejects it with usage exit
   2; the publish job depends on that smoke job at workflow line 63.
2. **P1:** unreadable identity inputs replace the prior trusted generation.
   The spec requires unreadable or unsafe inputs to preserve prior evidence
   (`specs/governed-release-evidence.md:131`–`:136`), but an unreadable
   `package.json`, `go.mod`, or HEAD takes `internal/preflight/command.go:48`–`:55`
   through `PromoteEvidence`, which replaces the directory with phase records
   and `manifest.json` while omitting its prior index and checksums
   (`internal/preflight/artifact_evidence.go:315`–`:353`).
3. **P2:** required reruns fail on supported macOS. Coverage row 15 requires
   cross-verdict rerun idempotency
   (`specs/governed-release-evidence.md:223`–`:224`), but replacement calls the
   directory-exchange primitive at `internal/preflight/artifact_evidence.go:341`–`:350`,
   which unconditionally rejects Darwin at `:363`–`:366`.
4. **P2:** the release index omits Node and npm toolchain evidence. The spec
   requires identity/toolchains in the index
   (`specs/governed-release-evidence.md:123`–`:130`, `:217`), while
   `internal/preflight/release_evidence.go:39`–`:62` records only one Go
   toolchain string and repository search finds no observed Node/npm version
   capture in the release-evidence path.

## Coverage

1 finding. Worst: High.

1. **High:** concurrent first artifact builds targeting the same absent output
   can both pass `scripts/build-artifacts.sh:55`; one `mv` creates the output and
   the other nests its artifact directory beneath it at `:64`, while both exit
   zero. The mapped rerun row covers sequential and abandoned-stage reruns
   (`specs/governed-release-evidence.md:224`), and the artifact contract invokes
   one builder at a time (`internal/contract/surface/artifact/artifact_test.go:35`–`:73`).
   Add a concurrent same-destination row and black-box test requiring one exact
   five-tarball generation or a fail-closed loser.
