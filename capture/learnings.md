# Learnings — usage journal

- 2026-08-12  pocock-guidance-doctrine build: three coordinator mutation probes went
  vacuous because sed patterns spanned prose line wraps; the green looked like a pass
  until the mutated text was re-checked. Right behavior: after applying any probe
  mutation, verify it landed (rg the mutated text) before reading the check's verdict.
  Proposed rule: craft-delegate's verification list gains "confirm the mutation is
  present in the file before trusting its result".
- 2026-08-12  Twice ran the landing ff-merge inside the ticket worktree (merging the
  branch into itself, "Already up to date") instead of the main checkout; caught by
  the release refusal. Right behavior: the landing merge always runs in the primary
  checkout; worth a one-line note wherever the landing cadence is documented.
