# Pass the chunk base to the commit chain

Blocked by: none
Writes: internal/commit/commit.go, internal/commit/commit_test.go, cmd/bench/commit_chain.go, cmd/bench/commit_chain_test.go, cmd/bench/help_inventory_test.go, internal/commit/chain_grammar_test.go, internal/commit/grammar.go (new)
Covers: BC38, BC39, BC40

## What to build

`bench commit` takes `--base <commit>` only together with `--preflight-build <slug>`. The commit outcome carries the base, and the chain passes `--base <commit>` and `--source-tip <published>` to the build preflight step.
In `commitChainCommand`, the preflight arguments gain the base through the line `args = append(args, evidencecmd.FlagBase, outcome.Base)`, because the probe-06 command omits that exact text.
A call that names `--base` without `--preflight-build` prints the usage line at exit 2 and publishes nothing.

The usage line, the help row, and the chain help name `--base <commit>` in the `--preflight-build` group. The root help inventory test and the commit help test follow the new line.

The commit opener names its seam through `otelrecord.SeamCommit`: `otelCommitSeam` in `internal/commit/commit.go` takes that constant, so the cost renderer and the producer read one seam name.

## Acceptance

- [ ] `bench commit --preflight-build s --base B` passes `--base B` and `--source-tip <published>` to the preflight step (BC38).
- [ ] `bench commit --base B -m m p` without `--preflight-build` prints the usage line at exit 2 and publishes nothing (BC39).
- [ ] `bench commit --help` names `--base <commit>` in the `--preflight-build` group (BC40).
- [ ] `otelCommitSeam` is `otelrecord.SeamCommit` (BC-C2 review, Standards N1).
