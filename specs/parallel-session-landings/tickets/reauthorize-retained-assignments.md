# Reauthorize a retained assignment whose request token was lost

Blocked by: expose-explicit-source-review.md
Writes: internal/intent, internal/worktree, internal/usage, cmd/bench/main.go, bin/bench.sh, internal/systemtest

## What to build

Add the reviewer-held recovery operation `bench worktree reauthorize`. It accepts
an exact assignment ID, new opaque request token, approved base, current source
tip, and owned worktree path. Under the ledger lock, verify the assignment,
owner marker, path, branch, tip, inclusive recorded-start ancestry, and approved-
base ancestry, then compare-and-swap only the request digest against the prior
digest read in that same operation. Disclose the recorded start and approved
base even when they differ.

Both ancestry proofs consume the shared source-range owner extended by
expose-explicit-source-review.md, which is why that ticket blocks this one. Add
no ancestor parser to intent or worktree.

Do not discover assignments by path, add an assignment field or schema version,
or blind-write the ledger. Command parsing must keep every flag value separate
from the required positional path. This is a complete independently green
recovery command before landing porcelain consumes the replacement token.

## Acceptance

- [ ] Each reauthorize value flag, including `--assignment`, is proved unable to
      masquerade as the positional path, with `--` path controls (covers PL27).
- [ ] Exact verification plus one expected-old digest CAS replaces the token;
      start equal to tip is accepted, start need not equal approved base,
      concurrent movement refuses, and the retained worktree bytes are unchanged
      (covers PL29).
