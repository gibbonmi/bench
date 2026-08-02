# Surface git's own error when staging fails

Blocked by: none

## What to build

When `bench commit`'s post-gate staging of a named path fails, the user sees
git's own diagnostic (e.g. the `index.lock` held-by-another-process message)
alongside the existing `error: staging %q failed` line, instead of a bare
`exit status 128`.

**Deviation from the issued direction, flagged per the prompt.** The prompt
prescribed deletion-aware staging (`git add -A -- <path>`) for a named
directory whose contents are all deleted. That shape is already implemented
(`internal/commit/commit.go` stages with `git add -A -- :(literal)<path>`,
FT87 slice 3) and already pinned by the contract
`deleted named directory commits its removals`; the retire-shaped repro —
including the exact 14-path shape of the 2026-08-02 field failure — commits
green through the current binary. What reproduces the field error
byte-for-byte (`error: staging "<dir>" failed: exit status 128` after a green
gate) is a held `.git/index.lock` — consistent with the concurrent pcgs
session the handoff warned to serialize against. The misdiagnosis was possible
only because staging runs `git add` with its stderr discarded, so the one
line git prints to name the real cause never reached the user. This ticket
fixes that opacity; whether staging should also retry briefly on a held
index.lock is a behavior decision left to the reviewer.

## Acceptance

- [x] A `bench commit` whose staging `git add` fails relays git's stderr to
  the caller's stderr (a held `.git/index.lock` shows git's
  another-process message, not only `exit status 128`).
- [x] The existing `error: staging %q failed` refusal line and exit 1 are
  unchanged, and a green staging path prints nothing new.
