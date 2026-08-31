# Stat every consumer-payload allowlist source

Blocked by: none
Writes: internal/conformance/package_shipped_surface_test.go
Covers: CE14, CE15

## What to build

The shipped-surface conformance check derives a `package.json` exclusion from each
kit-only allowlist row. It never proves that a row's source path exists, so a
misspelled row passes and excludes nothing.

The check now stats the source path of every allowlist row, on both audiences. A
tree row must name a directory. A file row must name a regular file. A row that
fails either test reds the check. The red names the row's source and the row's
audience, so the repair starts at the row.

The fixture that proves the red mirrors the misspelling in `package.json`
files[]. The derived-exclusion loop then stays silent, and only the new stat
check reds. The existence check lives beside the derived-exclusion loop. It uses the canonical
payload reader that the loop already calls.

## Acceptance

- [ ] CE14 — an allowlist row whose source path is absent reds the shipped-surface
      check.
- [ ] CE14 — the misspelled-row fixture mirrors the misspelling in `package.json`
      files[], so the delivered exclusion loop stays silent.
- [ ] CE14 — a tree row that names a regular file reds the check, and a file row
      that names a directory reds the check.
- [ ] CE14 — the check stats the rows of both audiences.
- [ ] CE15 — the absence red names the row's source and the row's audience.
- [ ] `bench gate` stays green over the live tree.

## Delegate charge

You work in the Bench repo on the `citation-execution-proof` spec. Read
`specs/citation-execution-proof/spec.md` first. Then read
`internal/conformance/package_shipped_surface_test.go` in full. Also read the
`kitpayload` payload reader that this file calls.

Add a loop beside the derived-exclusion loop. Stat the source path of every
allowlist row, for both audiences. Require a directory for a row whose `Tree`
field is true. Require a regular file for every other row. Append one diagnostic
per failing row. Name the row's source and audience in that diagnostic.

Reuse the rows the canonical reader already returned. Do not open the payload a
second time.

Add `TestAllowlistSourceExists`. Drive the check over a fixture root with a
misspelled row. Mirror the misspelling in the fixture's `package.json` files[],
so only the new diagnostic appears. Assert the exact diagnostic text. Confirm the live tree stays green.

Run only `bench worktree exec <label> -- go test ./internal/conformance/`. Do not
commit. Do not edit the spec.
