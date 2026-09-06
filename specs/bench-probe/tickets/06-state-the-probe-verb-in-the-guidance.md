# State the probe verb in the guidance

Blocked by: 04-add-the-bench-probe-verb.md
Writes: .bench/BENCH-reference.md, CONTEXT.md, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route
Covers: PB36, PB37

## What to build

Verify the premise first. Read the Command Notes section of .bench/BENCH-reference.md.
Run `bench anchors .bench/BENCH-reference.md` and `bench anchors CONTEXT.md`, and keep
every required needle in place. Read the **mutation probe** entry in CONTEXT.md.

Add one paragraph to the Command Notes after the `bench consumers` paragraph. It
names `bench probe`, the sequence it runs, and the four verdicts with their exit
codes. It names the preserved copy under `$BENCH_HOME/probe/<repo-key>/` and the
`bench worktree exec <target> -- bench probe ...` form for a worktree. Write it in
ASD-STE100, and keep the forbid needles out of the text.

Add the glossary entry **probe verdict** after **mutation probe**: the one word
`bench probe` prints for a run, `bit`, `silent`, `invalid`, or `restore-failed`. Its
Avoid list holds `test result`, `probe outcome`, and `mutation score`.

Run `bench gate-prose . -- .bench/BENCH-reference.md CONTEXT.md` and
`bench test --check docs-currency-workflow`. The fixture directories in the Writes line
are closure headroom, and they take no edit.

## Acceptance

- [ ] The Command Notes paragraph names the verb, the four verdicts with their exits, the preserved copy, and the exec form.
- [ ] `CONTEXT.md` holds **probe verdict** with its three Avoid words.
- [ ] `bench gate-prose` and `bench test --check docs-currency-workflow` pass over the two files.
