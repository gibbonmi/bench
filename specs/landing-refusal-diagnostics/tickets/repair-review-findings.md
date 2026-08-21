# Repair landing-refusal review findings

Writes (advisory): internal/worktree/worktree.go, internal/worktree/ownership.go, internal/worktree/land.go, internal/worktree/land_test.go, reviews/landing-refusal-diagnostics.md
Line: gpt-5.6-sol / high / ~2 iterations

## Review predicates

- Every retained-release `next=` uses `--request <request>` and never echoes the caller token; the unsafe-path assignment-pointer form, output channel, and exit status stay unchanged.
- The unmatched-request helper does not carry a six-bare-string/data-clump signature: a small recovery-context value groups the optional target, detail, base, and tip while opaque caller-token lookup stays in `intent` and introduces no behavioral seam.
