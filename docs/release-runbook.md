# Release publication runbook

This is the operator's procedure for publishing a Bench release to the public npm
registry. It exists because the hermetic gate proves only the fixture-exercisable
state machine: a live publish, real-registry integrity, name ownership, and
interactive 2FA approval are evidence only a present maintainer can produce. Every
step here is a manual precondition or a live-registry action; nothing in this file
is optional polish. Automation submits and verifies; only a person authorizes.

## Preconditions (all manual, all blocking)

1. **Publication authority is proven, never inferred.** The five identities are
   `redbench` and the four `@redbench/*` platform packages. Before the first
   publication of any of them, confirm name ownership and publish
   permission directly in the npm web UI under the publishing account. An `E404`
   from `npm view` is ambiguous: it cannot distinguish an available public name
   from an identity you cannot access. So registry probes prove nothing about
   authority. Do not proceed on a probe.
2. **Authorization comes from the authoritative preflight only.** A full (not
   focused) `bench release-preflight --mode publish --profile public|bank` run
   must be green for the selected profile. Every required producer record must
   be present: FT88 data handling, FT87 offline/network control, and for
   `bank` FT71 local events. A red or focused preflight never authorizes
   publication — there is no override path.
3. **Toolchain pins.** The staged path requires npm 11.15+ and Node 24+ on
   the operator machine. Older npm has no `stage` command; do not substitute
   direct publish for a staged release because the tool is missing.
4. **Reviewer presence.** First publication, staged approval or rejection, and
   final promotion each require the reviewer present. For the first publication
   the attended act is the reviewer's own tag push, which the workflow's submit
   job then carries out (see below). Staged approval and promotion are
   interactive and have no unattended path at all.

## First publication (identities do not exist yet)

Staging cannot create a new package identity, so the first release is a
reviewer-authorized direct publish driven by `bench release submit` against the
approved release directory:

1. Platform packages publish first, each under the version-specific candidate
   tag (never `latest`). Each live registry SHA-512 integrity is verified
   against the approved local tarball before the next package advances.
2. The wrapper publishes last and is verified the same way.
3. A candidate tag is public, not private: an exact version or the tag remains
   installable. Treat the release as live from the first successful publish.

The reviewer's presence is spent on the tag, not on the publish. Only the
reviewer cuts and pushes a release tag; that push is the attended act. The
release workflow's publish job is its mechanical arm, running exactly one
command from the tag's own checkout:

    dist/bench release submit --version "${GITHUB_REF_NAME#v}" --profile public \
      --path first --adapter npm --provenance --registry https://registry.npmjs.org

The job compiles that binary from the checked-out tag and downloads the
publish-preflight evidence the `authorize` job produced. It uploads
`dist/publication/` as the `publication-record` artifact even when the run
fails. So a partial publication is diagnosable and resumable. CI never
promotes.

## Live integrity re-verification (before any promotion)

Before `bench release promote`, re-verify the complete live set against the
approved local tarballs — every package, not a sample. The command performs this,
and the operator confirms the recorded registry integrities in
`publication-record.json` match the release-index-bound local integrities. A
single mismatch is terminal for this version: stop, roll back by deprecation,
and cut a new version. Bytes that drifted on the registry are never accepted
because the version number matches.

## Staged releases (identities exist)

1. CI or the operator submits every package with stage-only OIDC under the
   version-specific candidate tag (`bench release submit`); staging never makes
   a package live.
2. The command emits an exact `next_action` for the interactive step: the
   maintainer reviews the staged bytes and approves with 2FA. Platform packages
   go first, the wrapper only after every platform integrity verifies live.
   Approval is deliberately proof-of-presence; there is no non-interactive
   substitute, and OIDC cannot approve, view, or download a stage.
3. Before any approval, a rejected or failed stage is safe cleanup. After any
   package is live, failure follows the rollback procedure below.

## Promotion

`bench release promote` moves `latest` only after the complete live set
re-verifies: platform `latest` tags first, the wrapper `latest` last, so no
consumer can resolve a wrapper whose platform dependency is not yet default.

## Rollback (never unpublish)

A published version is immutable and its name/version is burned even after
unpublish, so automation never unpublishes and the operator does not either.
`bench release rollback` preserves the prior `latest`, removes the candidate
tags, and deprecates the bad version with a recovery message pointing at the
replacement. Recovery is always a new version.

## Proof scope

Bench ships more targets than it proves. The release plan in
`scripts/release-plan.json` gives each target a `native_proof` field, and that
field is the one source for the proven list. The shipped list and the proven
list are separate facts, and the plan is the one place that states which is
which.

The release ships four targets:

- `darwin/arm64`
- `darwin/x64`
- `linux/arm64`
- `linux/x64`

The release proves two of them with a native proof: `linux/arm64` and
`linux/x64`. Bench does not prove the two Darwin binaries with a native proof.
The `smoke` job still runs on macOS runners, so it remains the macOS execution
evidence for the shipped macOS binaries.

## Reproducibility

The release keeps a byte-for-byte comparison of the artifacts. The artifact
build makes a second independent build from the same source, compares the two
outputs, and writes the verdict to `dist/reproducibility.json`. The release
evidence reads that record, so the published reproducibility claim stays true.

Bench does not compare the finalized release evidence of two independent
checkouts. Bench retired that cross-checkout comparison, and no release
document promises it now.

## Evidence

Every run appends ordered transitions to the durable
`publication-record.json`, bound to the immutable release-index digest. The
operator confirms after promotion that the record carries the full transition
history and no credential material. The record is publication evidence; it never
rewrites the reproducible release index, `SHA256SUMS`, or package bytes.
