# Review pickup — one-source-collapse (c796de9)

Three-axis semantic review of the FT50 landing commit. This file is transient
pickup state: the `/bench-implement-spec` session that resolves it deletes it in
the same green fix commit.

## Standards

1 finding (judgment call, advisory).

- **Duplicated shim heredoc** — `internal/conformance/guard_resolver_drift_test.go`:
  the lib fixture is factored to one const (`canonicalResolveBenchLib`) but the
  shim fixture is pasted as a near-identical ~10-line heredoc in three tests
  (renamed-anchor, swapped-order, accepts-matching), differing only in function
  name or candidate order. A `writeShim(name, order)` helper collapses them the
  way the lib const already does. Cited standard: AGENTS.md one-source-per-fact
  (Fowler *Duplicated Code*); the "honest repetition" carve-out arguably applies,
  so advisory only.

Worst issue: the triplicated heredoc — advisory.

## Spec

0 findings. All five stories, every [default] decision, and all 8 coverage-map
rows traced to the diff; no scope creep, nothing missing.

## Coverage

4 findings. Verified by the aggregator against the code before persisting.

1. **(worst) Inserted/foreign candidate passes green — silent false pass.**
   `extractResolverCandidateOrder` (guard_resolver_drift_test.go) records only
   the first occurrence of exactly three fixed tokens and the check accepts on
   `len==3` + equal order. Breaking input: a shim loop gaining
   `"$HOME/.local/bin/bench.sh"` between the two known wrappers — the new path
   is not a token, and after repo-wrapper masking its `bin/bench.sh` tail
   aliases the kit-wrapper token, so extraction still yields
   `[repo, kit, path]` → green. A genuinely new search candidate — the exact
   drift class the check exists to catch — is invisible. No test inserts a 4th
   or substring-aliasing candidate.
2. **In-body comment defeats first-occurrence ordering — conditional false
   pass.** `extractResolverFunctionBody` keeps comment lines; a stale comment
   naming the paths in canonical order above a swapped operative loop yields
   green. Latent today (the real shim's comment sits above the signature,
   outside the extracted body); no test exercises a resolver body containing a
   comment.
3. **Multi-line inline group truncates the body — false red, fails loud.**
   Reformatting the shim's `{ printf ...; return 0; }` group across lines puts
   a bare `}` before the function's real close; the body truncates before
   `command -v bench` → missing-candidate red on a correct shim. Acceptable
   fail-closed posture; recorded for awareness.
4. **(minor) Brace-on-next-line signature reds loud** — `sigRe` requires
   `name() {` on one line. Same fail-loud class as 3.

Worst issue: #1 — the check passes green on the drift class it was built to
catch. Suggested fix shape for 1+2: extract candidates from the operative
lines (the `for candidate in ...` list and the `command -v bench` line) after
stripping comment lines, and red on any candidate-like token that is not one
of the three known ones. Findings 3–4 are accept-candidates (fail-loud).
